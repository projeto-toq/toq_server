# Guia de Processamento de Mídias — Listings em `PENDING_PHOTO_PROCESSING`

## 1. Contexto e Objetivo
- `listingmodel.StatusPendingPhotoProcessing` indica que a sessão de fotos terminou e o fotógrafo precisa subir mídias brutas para o bucket `toq-listing-medias`.
- Toda a orquestração backend usa `listingIdentityId` (não `listingId`). Esse identificador aparece em todas as chaves S3, registros de banco (`listing_media_assets`, `media_processing_jobs`) e payloads AWS.
- Este guia descreve o comportamento implementado em `internal/core/service/media_processing_service`, adapters S3/SQS/Step Functions e as Lambdas Go em `aws/lambdas/go_src`.
- Objetivo: alinhar frontend, backend e plataforma sobre contratos HTTP, estrutura dos jobs, observabilidade e como depurar falhas.

## 2. Fluxo ponta a ponta (alto nível)
1. **Manifesto de upload** – o cliente chama `POST /api/v2/listings/media/uploads` enviando `listingIdentityId` + lista de arquivos (`assetType`, `sequence`, `bytes`, `checksum`). O serviço valida status do listing e limites configurados.
2. **URLs pré-assinadas** – a resposta entrega URLs PUT com TTL (default 900s) e cabeçalhos obrigatórios (`Content-Type`, `x-amz-checksum-sha256`). O cliente envia diretamente para `toq-listing-medias/<listingIdentityId>/raw/...`.
3. **Processamento incremental** – `POST /uploads/process` move assets `PENDING_UPLOAD`/`FAILED` para `PROCESSING`, registra um job em `media_processing_jobs` e publica uma mensagem SQS `listing-media-processing-staging`.
4. **Pipeline AWS** – Step Functions `listing-media-processing-sm-staging` valida objetos, gera thumbnails (Lambda `listing-media-thumbnails-staging`), dispara MediaConvert para vídeos e consolida os resultados via Lambda `listing-media-consolidate-staging`.
5. **Callback** – a Lambda `listing-media-callback-staging` envia `POST /api/v2/listings/media/callback` com assinatura `X-Toq-Signature`. O serviço `HandleProcessingCallback` atualiza assets individuais e o job.
6. **Gestão** – o usuário pode listar, renomear (`POST /update`) ou excluir (`DELETE /delete`) assets `PROCESSED`/`FAILED` antes da finalização.
7. **Finalização** – `POST /uploads/complete` confirma que todos os assets estão `PROCESSED`, altera o listing para `StatusPendingOwnerApproval`, registra novo job e dispara Step Functions `listing-media-finalization-sm-staging` (Lambda `listing-media-zip-staging`).
8. **Distribuição** – `POST /media/download` gera URLs GET assinadas (TTL default 3600s) para cada asset ou para o ZIP `/<listingIdentityId>/processed/zip/listing-media.zip`.
9. **Aprovação do proprietário** – `POST /media/approve` permite que o owner aceite ou rejeite os materiais. Quando `LISTING_APPROVAL_ADMIN_REVIEW=true`, aprovações vão para `StatusPendingAdminReview`; caso contrário, avançam direto para `StatusReady`. Rejeições retornam o fluxo para `StatusRejectedByOwner`.

## 3. Modelos de domínio e persistência
- **MediaAsset (`internal/core/model/media_processing_model/media_asset.go`)**
	- Chaves: `(listingIdentityId, assetType, sequence)`.
	- Campos relevantes: `Status` (`PENDING_UPLOAD`, `PROCESSING`, `PROCESSED`, `FAILED`), `S3KeyRaw`, `S3KeyProcessed`, `Metadata` (JSON com `clientId`, `filename`, etc.).
- **MediaProcessingJob (`internal/core/model/media_processing_model/media_job.go`)**
	- Campos: `id`, `listingIdentityId`, `status` (`PENDING`, `RUNNING`, `SUCCEEDED`, `FAILED`), `provider` (`STEP_FUNCTIONS`, `STEP_FUNCTIONS_FINALIZATION`), `externalId` (ARN do Step Functions), `payload` (último callback), `lastError`, `callbackBody`.
	- `ApplyFinalizationPayload` guarda `zipBundles`, `assetsZipped`, `zipSizeBytes` e `unzippedSizeBytes` para o bundle final.
- **Persistência**
	- `media_processing_jobs.external_id` espelha `executionArn`.
	- `media_processing_jobs.callback_body` mantém o JSON bruto recebido do pipeline para auditoria.
	- `listing_media_assets` guarda tanto as chaves S3 quanto metadados usados nos presigns (sequência, título, etc.).

## 4. Endpoints HTTP (`/api/v2/listings/media`)
### 4.1 `POST /listings/media/uploads`
Request:
```json
{
	"listingIdentityId": 123,
	"files": [
		{
			"assetType": "PHOTO_HORIZONTAL",
			"sequence": 1,
			"filename": "IMG_2907.jpg",
			"contentType": "image/jpeg",
			"bytes": 5242880,
			"checksum": "sha256:4ee0c4...",
			"title": "Sala social",
			"metadata": {
				"clientId": "photo-1",
				"orientation": "LANDSCAPE"
			}
		}
	]
}
```
Response (`RequestUploadURLsResponse`):
```json
{
	"listingIdentityId": 123,
	"uploadUrlTtlSeconds": 900,
	"files": [
		{
			"assetType": "PHOTO_HORIZONTAL",
			"sequence": 1,
			"uploadUrl": "https://toq-listing-medias.s3.amazonaws.com/123/raw/photo/horizontal/horizontal-01-IMG_2907.jpg?...",
			"method": "PUT",
			"headers": {
				"Content-Type": "image/jpeg",
				"x-amz-checksum-sha256": "TuDE4w=="
			},
			"objectKey": "123/raw/photo/horizontal/horizontal-01-IMG_2907.jpg"
		}
	]
}
```

### 4.2 `POST /listings/media/uploads/process`
Body:
```json
{ "listingIdentityId": 123 }
```
Pré-condições: listing em `PENDING_PHOTO_PROCESSING` ou `REJECTED_BY_OWNER`. O serviço registra `media_processing_jobs` (status `PENDING`), marca assets como `PROCESSING` e envia `MediaProcessingJobMessage` para SQS.

### 4.3 `GET /listings/media`
Query params: `listingIdentityId` (obrigatório), `assetType`, `sequence`, `page`, `limit`, `sort` (`sequence|id`), `order` (`asc|desc`).
Response:
```json
{
	"data": [
		{
			"id": 42,
			"listingIdentityId": 123,
			"assetType": "PHOTO_VERTICAL",
			"sequence": 1,
			"status": "PROCESSED",
			"title": "Entrada",
			"metadata": {
				"clientId": "photo-3"
			},
			"s3KeyRaw": "123/raw/photo/vertical/vertical-03-IMG_2705.jpg",
			"s3KeyProcessed": "123/processed/photo/vertical/large/vertical-03-IMG_2705.jpg"
		}
	],
	"pagination": { "page": 1, "limit": 20, "total": 4 },
	"zipBundle": {
		"bundleKey": "123/processed/zip/listing-media.zip",
		"assetsCount": 42,
		"zipSizeBytes": 83886080,
		"estimatedExtractedBytes": 209715200,
		"completedAt": "2025-01-03T15:04:05Z"
	}
}
```

### 4.4 `POST /listings/media/update`
Body:
```json
{
	"listingIdentityId": 123,
	"assetType": "PHOTO_VERTICAL",
	"sequence": 1,
	"title": "Varanda gourmet",
	"metadata": {
		"tag": "capa"
	}
}
```

### 4.5 `DELETE /listings/media/delete`
Body:
```json
{
	"listingIdentityId": 123,
	"assetType": "PHOTO_VERTICAL",
	"sequence": 1
}
```
O serviço remove o registro via transação e, após commit, executa limpeza assíncrona no S3 para todas as chaves retornadas por `asset.GetAllS3Keys()`.

### 4.6 `POST /listings/media/download`
Body:
```json
{
	"listingIdentityId": 123,
	"requests": [
		{ "assetType": "PHOTO_VERTICAL", "sequence": 1, "resolution": "thumbnail" },
		{ "assetType": "PHOTO_HORIZONTAL", "sequence": 2, "resolution": "large" },
		{ "assetType": "VIDEO_VERTICAL", "sequence": 1, "resolution": "original" }
	]
}
```
Response: lista de URLs GET assinadas com `expiresIn` configurado (default 3600s).

### 4.7 `POST /listings/media/uploads/complete`
Body:
```json
{ "listingIdentityId": 123 }
```
Valida que todos os assets retornados por `ListAssets` estão com `Status=PROCESSED` e `s3_key_processed` não vazio (`ensureAssetsReadyForFinalization`). Cria um job `STEP_FUNCTIONS_FINALIZATION`, chama `StartMediaFinalization`, atualiza `media_processing_jobs` com `external_id` e move o listing para `StatusPendingOwnerApproval`.

### 4.8 `POST /listings/media/callback`
- Header obrigatório: `X-Toq-Signature = hex(hmac_sha256(CALLBACK_SECRET, raw_body))`.
- Payload (`MediaProcessingCallbackRequest`):
```json
{
	"executionArn": "arn:aws:states:us-east-1:058264253741:execution:listing-media-processing-sm-staging:process-28-20",
	"jobId": "20",
	"listingIdentityId": "28",
	"externalId": "arn:aws:states:us-east-1:058264253741:execution:listing-media-processing-sm-staging:process-28-20",
	"status": "SUCCEEDED",
	"provider": "STEP_FUNCTIONS",
	"traceparent": "00-2b63e64e71537bb0327788965465ed16-45348f2c8c2a34bf-01",
	"outputs": [
		{
			"rawKey": "28/raw/photo/horizontal/horizontal-01-IMG_2907.jpg",
			"processedKey": "28/processed/photo/horizontal/large/horizontal-01-IMG_2907.jpg",
			"thumbnailKey": "28/processed/photo/horizontal/thumbnail/horizontal-01-IMG_2907.jpg",
			"outputs": {
				"large_photo_horizontal": "28/processed/photo/horizontal/large/horizontal-01-IMG_2907.jpg"
			},
			"errorCode": "",
			"errorMessage": ""
		}
	],
	"assetsZipped": 0,
	"zipBundles": [],
	"zipSizeBytes": 0,
	"unzippedSizeBytes": 0,
	"failureReason": "",
	"error": null
}
```
Para o pipeline de zip, `provider = STEP_FUNCTIONS_FINALIZATION`, `status` pode ser `SUCCEEDED` ou `FINALIZATION_FAILED`, e `zipBundles` contém chaves como `"28/processed/zip/listing-media.zip"` acompanhadas de `zipSizeBytes` e `unzippedSizeBytes`.

### 4.9 `POST /listings/media/approve`
Body:
```json
{
	"listingIdentityId": 123,
	"approve": true
}
```
Regras:
- Apenas o owner autêntico pode aprovar/rejeitar.
- O listing precisa estar em `StatusPendingOwnerApproval`, caso contrário retorna 400.
- `LISTING_APPROVAL_ADMIN_REVIEW=true` move aprovações para `StatusPendingAdminReview`; quando `false`, o status final é `StatusReady`.
- Rejeições retornam o status para `StatusRejectedByOwner`, permitindo novos uploads/edições.

## 5. Orquestração AWS

### 5.1 Produção do job
`ProcessMedia` cria um `MediaProcessingJobMessage`:
```json
{
	"jobId": 21,
	"listingIdentityId": 28,
	"assets": [
		{ "key": "28/raw/photo/horizontal/horizontal-01-IMG_2907.jpg", "type": "PHOTO_HORIZONTAL" }
	],
	"retry": 0
}
```
Mensagem vai para `listing-media-processing-staging` (atributos SQS: `ListingIdentityId`, `JobId`, `RetryCount`, `Traceparent`). Retries usam `listing-media-processing-dlq-staging`.

### 5.2 Step Functions `listing-media-processing-sm-staging`
Estados:
1. **ValidateRawAssets** (`listing-media-validate-staging`) – HEAD nos objetos, injeta `hasVideos`, normaliza payload.
2. **ParallelProcessing**
	 - `GenerateThumbnails` (`listing-media-thumbnails-staging`) escreve `processed/photo/.../{thumbnail|small|medium|large}/`.
	 - `CheckVideoProcessing` -> `ProcessVideos` (AWS MediaConvert) quando `hasVideos=true`.
3. **ConsolidateResults** (`listing-media-consolidate-staging`) – reconcilia outputs, define `processedKey` (preferência `large > medium > small > thumbnail`) e `thumbnailKey`, propaga `errorCode` (`THUMBNAIL_PROCESSING_FAILED`, `VALIDATION_ERROR`, etc.).
4. **FinalizeAndCallback** (`listing-media-callback-dispatch-staging`) – envia payload assinado para o backend.
5. **ValidationFailed/ReportFailure** – mesmas Lambda de callback, com `status=VALIDATION_FAILED` ou `PROCESSING_FAILED`.

### 5.3 Finalização `listing-media-finalization-sm-staging`
1. **CreateZipBundle** (`listing-media-zip-staging`) – recebe `MediaFinalizationInput` com todos os `S3KeyProcessed`; escreve `/<listingIdentityId>/processed/zip/listing-media.zip`.
2. **FinalizeAndCallback** (`listing-media-callback-dispatch-staging`) – envia `status=SUCCEEDED`, `provider=STEP_FUNCTIONS_FINALIZATION`, `zipBundles` e `assetsZipped`.
3. **ReportFailure** – envia `status=FINALIZATION_FAILED` para o backend.

### 5.4 Lambdas
| Função | Descrição |
| --- | --- |
| `listing-media-validate-staging` | Confere existência dos objetos, checksum, constrói `traceparent`. |
| `listing-media-thumbnails-staging` | Usa `disintegration/imaging` para gerar tamanhos `thumbnail/small/medium/large` e corrigir EXIF. |
| `listing-media-zip-staging` | Consolida chaves processadas em um ZIP, garantindo nome `/<listingIdentityId>/processed/zip/listing-media.zip`. |
| `listing-media-consolidate-staging` | Monta `outputs[]`, define `processedKey`/`thumbnailKey`, agrega erros por asset. |
| `listing-media-callback-staging` | Recebe eventos (inclusive `body` vindo da Step Function) e faz POST para o backend com assinatura HMAC. |

## 6. Estrutura S3 (`toq-listing-medias`)
- **Raw uploads:** `/{listingIdentityId}/raw/{mediaTypeSegment}/{reference}-{filename}`  
	- `mediaTypeSegment` via `mediaTypePathSegment`: `photo/horizontal`, `photo/vertical`, `video/horizontal`, `project/doc`, etc.
	- `reference` deriva de `metadata.clientId` ou `sequence` (`horizontal-01`, `vertical-03`, ...).
- **Processed assets:** `/{listingIdentityId}/processed/{mediaTypeSegment}/{size}/{filename}`  
	- `size ∈ {thumbnail, small, medium, large, original}`.
	- Vídeos processados ficam em `video/{orientation}/original`.
- **ZIP bundles:** `/{listingIdentityId}/processed/zip/listing-media.zip`.
- **TTL padrão:** upload URLs 900s, download URLs 3600s (configuráveis via `env.yaml`).
- **Checksum:** `ListingMediaStorageAdapter` aceita SHA-256 em hex (`sha256:...`) ou Base64 e converte para o formato exigido pelo S3 (`x-amz-checksum-sha256`).

## 7. Observabilidade e logs
- Todas as operações passam por `utils.GenerateTracer` e propagam `traceparent` para SQS/Step Functions.
- Logs úteis:
	- `service.media.process.started` – confirma `job_id`, `assets_count`.
	- `service.media.callback.asset_lookup_error` – indica que não foi possível casar `rawKey` com um asset; geralmente erro de payload.
	- `service.media.complete.started_zip` – contém `execution_arn` da finalização.
- Banco:
	```sql
	SELECT id, listing_identity_id, status, external_id, started_at, completed_at, last_error
		FROM media_processing_jobs
	 WHERE listing_identity_id = 28
	 ORDER BY id DESC LIMIT 5;
	```
	`callback_body` mantém o JSON recebido (útil para reproducibilidade).
- Callback Lambda loga payloads com `asset_error_codes` quando `status != SUCCEEDED`.
- Métricas: counters/histograms expostos via Prometheus nos adapters instrumentados.

## 8. Estados e regras de negócio
- **Assets**
	- `PENDING_UPLOAD` → após presign; só `ProcessMedia` muda para `PROCESSING`.
	- `PROCESSING` → aguardando Step Functions; `HandleProcessingCallback` promove para `PROCESSED` ou `FAILED`.
	- `FAILED` → pode ser reprocessado via novo `POST /uploads/process` ou removido (`DELETE /delete`).
	- `PROCESSED` → necessário `s3_key_processed` válido para permitir finalização e download.
- **Listings**
	- `CompleteMedia` só aceita listings em `PENDING_PHOTO_PROCESSING` com todos os assets `PROCESSED`; após iniciar o Step Functions de zip, status muda para `StatusPendingOwnerApproval`.
- **Segurança**
	- Todos os endpoints (exceto `/media/callback`) exigem Bearer token, auth middleware + permission middleware.
	- Callback rejeita requisições sem `X-Toq-Signature` válido.
- **Limites configuráveis** (`env.yaml`)
	- `MaxFilesPerBatch`, `MaxFileBytes`, `MaxTotalBytes`, lista de MIME types permitidos, flag `AllowOwnerProjectUploads`.

## 9. Checklist operacional
1. **Uploads falhando** – verificar `listing_media_assets` se `s3_key_raw` foi preenchido e usar `aws s3 ls s3://toq-listing-medias/<listingIdentityId>/raw/`.
2. **Processamento parado** – procurar `service.media.process.started` com `listing_identity_id`; usar `aws sqs receive-message ...listing-media-processing-staging`.
3. **Callback com `THUMBNAIL_PROCESSING_FAILED`** – conferir objeto `rawKey` (precisa conter `raw/` no caminho), pois o consolidator deriva as demais chaves com base nesse padrão.
4. **ZIP não gerado** – usar `media_processing_jobs.external_id` para descrever a execução `listing-media-finalization-sm-staging` e checar logs da Lambda `listing-media-zip-staging`.
5. **Assinatura inválida** – garantir que `CALLBACK_SECRET` usado pela Lambda corresponde a `MEDIA_PROCESSING_CALLBACK_SECRET` configurado no backend.

Este documento reflete o comportamento atual do código. Alterações em contratos, fluxos ou infraestrutura devem ser atualizadas aqui antes de qualquer rollout.
