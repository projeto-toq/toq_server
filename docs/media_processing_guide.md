# Guia de Processamento de Mídias — Listings em Status `PENDING_PHOTO_PROCESSING`

## 1. Contexto e Objetivo
- Fluxo vigente: `listingmodel.StatusPendingPhotoProcessing` indica que a sessão de fotos foi concluída e o fotógrafo precisa subir mídias brutas (fotos/vídeos) para posterior tratamento.
- Objetivo deste guia: alinhar backend, frontend e plataforma sobre o processo completo de upload, processamento assíncrono e disponibilização de mídias finais (ZIP completo, mídias individuais e thumbnails), respeitando arquitetura hexagonal, segurança e escalabilidade do TOQ Server.
- Escopo: novas APIs REST (`/api/v2`), serviços core, adapters S3/SQS/Lambda/MediaConvert, modelos de domínio, monitoramento e contratos com o frontend.

## 2. Visão Geral do Fluxo (alto nível)
1. **Manifesto de upload**: frontend envia lista de arquivos que serão enviados (tipo, orientação, extensão, tamanho estimado, checksum). Backend valida permissão e status do listing.
2. **Emissão de URLs pré-assinadas**: backend retorna URLs PUT S3 com metadados obrigatórios. Fotógrafo envia diretamentre para o bucket dedicado `toq-listing-medias`, sempre dentro do diretório raiz `/{listingId}/`. Listings classificados como empreendimento em construção permitem que o **owner** autenticado (ou representante designado) gere URLs para projetos e renders sem depender do fotógrafo.
3. **Processamento Incremental**: após o upload de um ou mais arquivos, o frontend chama o endpoint de processamento (`/process`). O backend identifica os assets em estado `PENDING_UPLOAD`, atualiza para `PROCESSING` e dispara o pipeline assíncrono (Step Functions).
4. **Pipeline AWS**: Step Functions + Lambdas (Go) validam entradas, geram thumbnails, transcoding de vídeo e atualizam o status de cada asset individualmente.
5. **Feedback Visual**: O frontend pode consultar o status de cada asset e exibir thumbnails assim que estiverem prontos, sem esperar o lote todo.
6. **Gestão de Mídia**: O usuário pode reordenar, renomear ou excluir assets a qualquer momento antes da finalização.
7. **Finalização**: Quando satisfeito, o usuário solicita a finalização (`/complete`). O backend dispara o pipeline de geração de ZIP e avança o status do Listing para `StatusPendingOwnerApproval`.
8. **Distribuição**: APIs de download retornam URLs para o ZIP completo e mídias individuais.

## 3. Estrutura de Domínio e Persistência
- **Novos modelos core (`internal/core/model/media_processing_model`)**
  - `MediaAsset`: item individual (campos: `ID`, `ListingID`, `Type`, `Sequence`, `Status`, `S3KeyRaw`, `S3KeyProcessed`, `Title`, `Metadata`). A unicidade é garantida por `(ListingID, Type, Sequence)`.
  - `MediaProcessingJob`: metadados de execução assíncrona (`JobID`, `ListingID`, `ExternalID`, `Status`, `Payload`).

## 4. Endpoints HTTP (Left Adapter)
Todos sob `/api/v2/listings/media/*`.

### 4.1 `POST /api/v2/listings/media/uploads` — Solicitar URLs de upload
- **Objetivo**: cliente envia lista de arquivos para upload.
- **Request Body**
  ```json
  {
    "listingId": 123,
    "files": [
      {
        "clientId": "photo-1",
        "mediaType": "PHOTO_HORIZONTAL",
        "contentType": "image/jpeg",
        "bytes": 5000000,
        "checksum": "sha256:...",
        "sequence": 1
      }
    ]
  }
  ```
- **Response 201**
  ```json
  {
    "files": [
      {
        "clientId": "photo-1",
        "uploadUrl": "https://s3/...",
        "s3Key": "123/raw/photo/horizontal/uuid.jpg"
      }
    ]
  }
  ```

### 4.2 `POST /api/v2/listings/media/uploads/process` — Disparar Processamento
- **Objetivo**: Inicia o processamento dos arquivos que foram enviados (status `PENDING_UPLOAD`).
- **Request Body**
  ```json
  {
    "listingId": 123
  }
  ```
- **Response 202 (Accepted)**

### 4.3 `POST /api/v2/listings/media/uploads/complete` — Finalizar e Gerar ZIP
- **Objetivo**: Consolida as mídias, gera o ZIP final e avança o status do listing.
- **Request Body**
  ```json
  {
    "listingId": 123
  }
  ```

### 4.4 `POST /api/v2/listings/media/update` — Atualizar Metadados
- **Objetivo**: Atualizar título ou sequência de um asset.

### 4.5 `DELETE /api/v2/listings/media` — Remover Asset
- **Objetivo**: Remover um asset (arquivos e registro).
    "submittedAt": "2025-11-11T17:03:12Z",
    "processing": {
      "jobId": "arn:aws:states:...",
      "steps": [
        { "name": "validate-raw-assets", "status": "SUCCEEDED" },
        { "name": "image-thumbnails", "status": "RUNNING" },
        { "name": "video-transcode", "status": "PENDING" }
      ]
    },
    "availableDownloads": [
      {
        "type": "ZIP_FULL",
        "ready": false
      },
      {
        "type": "THUMBNAILS",
        "ready": true
      }
    ]
  }
  ```

### 4.4 `POST /api/v2/listings/media/downloads/query` — URLs de download
- **Request Body**
  ```json
  {
    "listingId": 123,
    "batchId": "UUID"
  }
  ```
- **Response 200**
  ```json
  {
    "batchId": "UUID",
    "generatedAt": "2025-11-11T18:22:00Z",
    "ttlSeconds": 3600,
    "downloads": [
      {
        "type": "ZIP_FULL",
        "url": "https://s3/...",
        "expiresAt": "2025-11-11T19:22:00Z"
      },
      {
        "type": "PHOTO_VERTICAL",
        "items": [
          {
            "label": "Sala - Vertical",
            "preview": "https://cdn/.../thumbnail/photo.jpg",
            "url": "https://s3/.../large/photo.jpg",
            "resolution": "4000x6000",
            "title": "Sala Social",
            "sequence": 1
          }
        ]
      }
    ]
  }
  ```

## 5. Serviços Core e Dependências (Hexagonal)
- **Service**: `internal/core/service/media_processing_service`
- **Ports**: `mediaprocessingrepository`, `storageport`, `queueport`.

## 6. AWS — Serviços Recomendados e Justificativas
- **Amazon S3**: armazenamento de mídias brutas e processadas.
  - Prefixos: `/{listingId}/raw/{mediaType}/` e `/{listingId}/processed/{mediaType}/{size}/`.
  - Tamanhos de imagem: `thumbnail` (200px), `small` (400px), `medium` (800px), `large` (1200px).
- **AWS SQS**: desacopla API HTTP do pipeline.
- **AWS Step Functions**: orquestra pipeline.
- **AWS Lambda (Go)**:
  - `validate`: Validação de manifesto.
  - `thumbnails`: Redimensionamento e correção de rotação (EXIF).
  - `zip`: Geração de pacotes para download.
  - `consolidate`: Agregação de resultados.
  - `callback`: Notificação ao backend.

### 6.1 Sequência detalhada do pipeline de processamento
1. **Ingestão e validação (Lambda #1)**
  - Disparada ao mudar `MediaBatch` para `RECEIVED`.
2. **Geração de thumbnails (Lambda #2)**
  - Usa `disintegration/imaging` (Go) para processar imagens.
  - Outputs: escreve versões `thumbnail`, `small`, `medium`, `large` no bucket.
3. **Transcodificação de vídeos (AWS MediaConvert)**
  - Step Functions cria job MediaConvert.
4. **Geração de ZIPs (Lambda #3)**
  - Lambda monta pacote ZIP com todas as mídias processadas.
  - Estrutura interna do ZIP é limpa (sem pastas de sistema).
5. **Atualização de metadados finais (Lambda #4)**
  - Consolida dados e persiste.
6. **Callback para backend**
  - Step Functions consolida saídas das branches, agrega erros de thumbnail/validação e envia payload padronizado para o backend.

### 6.2 Estrutura do Payload de Callback

O payload enviado pelo Step Functions para o endpoint `/api/v2/listings/media/callback` segue o contrato abaixo. Todos os identificadores são numéricos e já convertidos para string pelo backend apenas no limite HTTP. **Importante:** o dispatcher precisa enviar o header `X-Toq-Signature` calculado como `hex(hmac_sha256(MEDIA_PROCESSING_CALLBACK_SECRET, raw_body))`; caso contrário o callback é rejeitado com `401`.

```bash
BODY='{"status":"SUCCEEDED","jobId":123}'
SIGNATURE=$(printf '%s' "$BODY" | openssl dgst -sha256 -hmac "$MEDIA_PROCESSING_CALLBACK_SECRET" -binary | xxd -p -c 256)
curl -X POST "$CALLBACK_URL" \
  -H "Content-Type: application/json" \
  -H "X-Toq-Signature: $SIGNATURE" \
  -d "$BODY"
```

```json
{
  "jobId": 123456,
  "listingIdentityId": 987654,
  "provider": "STEP_FUNCTIONS",
  "status": "SUCCEEDED",
  "failureReason": "",
  "error": null,
  "traceparent": "00-<trace-id>-<span-id>-01",
  "outputs": [
    {
      "rawKey": "987654/raw/photo/horizontal/uuid.jpg",
      "processedKey": "987654/processed/photo/horizontal/large/uuid.jpg",
      "thumbnailKey": "987654/processed/photo/horizontal/thumbnail/uuid.jpg",
      "outputs": {
        "thumbnail_PHOTO_HORIZONTAL": "987654/processed/photo/horizontal/thumbnail/uuid.jpg"
      },
      "errorCode": "",
      "errorMessage": ""
    }
  ]
}
```

Em cenários de falha:

- `status` assume valores `VALIDATION_FAILED`, `PROCESSING_FAILED` ou `FAILED`.
- `error` é um objeto com as chaves `code` e `message`, ambos gerados a partir da exceção capturada na Step Function.
- `failureReason` replica o resumo de erro de alto nível.
- Quando a Lambda de thumbnails falha em itens específicos, cada output recebe `errorCode` (`THUMBNAIL_PROCESSING_FAILED`) e `errorMessage` com detalhes. Esses campos são persistidos no metadata do asset pelo backend.
- Se a validação falhou para um asset, o Consolidate marca `errorCode = VALIDATION_ERROR` e propaga a mensagem original.

> **Observabilidade**: A Lambda `listing-media-callback-staging` loga o resumo das falhas antes de invocar o backend, agrupando códigos de erro e preservando `traceparent`. Após o callback, o serviço core (`HandleProcessingCallback`) agrega os mesmos códigos ao `MediaProcessingJob.LastError` e incrementa contadores de falha por asset.

## 7. Convenções de Segurança e Compliance
- **Controle de acesso**: validar `PhotographerUserID`.
- **Prefixos S3**: `/{listingId}/raw/{mediaType}/` e `/{listingId}/processed/{mediaType}/{size}/`.
- **TTL de URLs**: upload 15 min, download 60 min.
- **Checksums**: exigir `x-amz-checksum-sha256`.

## 8. Orquestração Assíncrona & Estados
- **Estados do lote**: `PENDING_UPLOAD`, `RECEIVED`, `PROCESSING`, `READY`, `FAILED`.

## 9. Observabilidade
- **Tracing**: `utils.GenerateTracer`.
- **Métricas**: Prometheus counters e histograms.

## 10. Considerações para o Frontend
- Fazer upload direto via `fetch`/XHR com `PUT`.
- Respeitar headers obrigatórios.
- Exibir progresso através do endpoint de status.
- Para download final, usar URLs GET pré-assinadas.

---
Este documento é ponto de partida oficial para desenvolvimento paralelo frontend/backend. Qualquer alteração de contrato deve ser negociada e refletida aqui antes da implementação.
