# TOQ Server - Media Processing Infrastructure
## Recursos AWS Criados - Staging Environment
**Data:** 29 de Novembro de 2025  
**Região:** us-east-1  
**Ambiente:** staging

---

## 1. AWS KMS - Criptografia

### Chave KMS
- **Alias:** `alias/toq-media-processing-staging`
- **Key ID:** `2fd6d9e5-dc43-4275-ba77-79c99cad509c`
- **ARN:** `arn:aws:kms:us-east-1:058264253741:key/2fd6d9e5-dc43-4275-ba77-79c99cad509c`
- **Status:** Enabled
- **Auto-rotation:** Habilitada (anual)
- **Uso:** S3, SQS, todos os serviços do pipeline

---

## 2. Amazon S3 - Armazenamento

### Bucket Principal
- **Nome:** `toq-listing-medias`
- **Região:** us-east-1
- **Criptografia:** SSE-KMS (chave acima)
- **Versionamento:** Enabled
- **Block Public Access:** Todas opções ativadas
- **Logging:** Enabled → `toq-logs-staging`
- **Lifecycle Rule:** `raw/` → Glacier após 180 dias

### Estrutura de Pastas (Atualizado)
- **Raw (Upload):** `/{listingIdentityId}/raw/{mediaTypeSegment}/{reference}-{filename}`
  - `mediaTypeSegment` segue `mediaTypePathSegment` (`photo/horizontal`, `video/vertical`, `project/doc`, etc.).
  - O `reference` deriva de `metadata.clientId` ou `sequence` (`horizontal-01`, `vertical-03`), evitando datas no caminho.
- **Processed:** `/{listingIdentityId}/processed/{mediaTypeSegment}/{size}/{filename}`
  - Tamanhos suportados: `thumbnail` (200px), `small` (400px), `medium` (800px), `large` (1200px), `original` (fallback para vídeos/documentos).
- **Zip Bundles:** `/{listingIdentityId}/processed/zip/listing-media.zip`
  - O arquivo é sobrescrito a cada finalização bem-sucedida, mantendo sempre a versão mais recente do bundle ZIP disponível para download.

### Bucket de Logs
- **Nome:** `toq-logs-staging`
- **Região:** us-east-1
- **Propósito:** Armazenar access logs do bucket principal

---

## 3. Amazon SQS - Filas de Mensagens

### Fila Principal
- **Nome:** `listing-media-processing-staging`
- **ARN:** `arn:aws:sqs:us-east-1:058264253741:listing-media-processing-staging`
- **URL:** `https://sqs.us-east-1.amazonaws.com/058264253741/listing-media-processing-staging`
- **Tipo:** Standard
- **Visibility Timeout:** 60 segundos
- **Message Retention:** 4 dias (345600s)
- **Long Polling:** 20 segundos
- **Criptografia:** SSE-KMS
- **Redrive Policy:** 5 tentativas → DLQ
- **Payload publicado:**
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
  A mensagem inclui atributos SQS (`ListingIdentityId`, `JobId`, `RetryCount`, `Traceparent`) para facilitar rastreamento.

### Dead Letter Queue (DLQ)
- **Nome:** `listing-media-processing-dlq-staging`
- **ARN:** `arn:aws:sqs:us-east-1:058264253741:listing-media-processing-dlq-staging`
- **URL:** `https://sqs.us-east-1.amazonaws.com/058264253741/listing-media-processing-dlq-staging`
- **Message Retention:** 14 dias (1209600s)
- **Criptografia:** SSE-KMS

---

## 4. AWS IAM - Roles e Permissões

### 4.1 Backend Role
- **Nome:** `toq-media-processing-backend-staging`
- **ARN:** `arn:aws:iam::058264253741:role/toq-media-processing-backend-staging`
- **Assume Role:** EC2
- **Permissões:**
  - S3: PutObject, GetObject, DeleteObject (toq-listing-medias)
  - SQS: SendMessage, GetQueueAttributes
  - KMS: Encrypt, Decrypt, GenerateDataKey
  - Step Functions: StartExecution, DescribeExecution

### 4.2 Lambda Role
- **Nome:** `toq-media-processing-lambda-staging`
- **ARN:** `arn:aws:iam::058264253741:role/toq-media-processing-lambda-staging`
- **Assume Role:** lambda.amazonaws.com
- **Permissões:**
  - S3: Full access (toq-listing-medias)
  - SQS: SendMessage, ReceiveMessage, DeleteMessage
  - KMS: Encrypt, Decrypt, GenerateDataKey
  - Step Functions: SendTaskSuccess, SendTaskFailure
  - MediaConvert: CreateJob, GetJob
  - CloudWatch Logs, X-Ray

### 4.3 Step Functions Role
- **Nome:** `toq-media-processing-stepfunctions-staging`
- **ARN:** `arn:aws:iam::058264253741:role/toq-media-processing-stepfunctions-staging`
- **Assume Role:** states.amazonaws.com
- **Permissões:**
  - Lambda: InvokeFunction (todas as 5 Lambdas)
  - MediaConvert: CreateJob, GetJob
  - SQS: SendMessage
  - CloudWatch Logs, X-Ray

### 4.4 MediaConvert Role
- **Nome:** `toq-media-processing-mediaconvert-staging`
- **ARN:** `arn:aws:iam::058264253741:role/toq-media-processing-mediaconvert-staging`
- **Assume Role:** mediaconvert.amazonaws.com
- **Permissões:**
  - S3: GetObject, PutObject (toq-listing-medias)

---

## 5. Implementação das Lambdas (Go)

As funções Lambda foram migradas para Go (1.25) utilizando Arquitetura Hexagonal.

### Funções
1. **Validate (`listing-media-validate-staging`)**: Valida assets no S3 contra o manifesto.
2. **Thumbnails (`listing-media-thumbnails-staging`)**:
   - Gera thumbnails usando `disintegration/imaging`.
   - Corrige rotação baseada em EXIF automaticamente.
   - Gera tamanhos: `thumbnail` (200px), `small` (400px), `medium` (800px), `large` (1200px).
3. **Zip (`listing-media-zip-staging`)**:
  - Cria arquivo ZIP a partir das midias originais (`raw/*`), mantendo a saida em `processed/zip/`.
  - Limpa estrutura de pastas interna (remove prefixos de sistema).
4. **Consolidate (`listing-media-consolidate-staging`)**: Agrega resultados do processamento paralelo.
5. **Callback Dispatch (`listing-media-callback-staging`, exposta como `listing-media-callback-dispatch-staging`)**: Envia webhook de volta ao backend com assinatura `X-Toq-Signature` e preserva `traceparent`.

### Runtime
- **Runtime:** `provided.al2023`
- **Build:** Binários compilados estaticamente via `./scripts/build_lambdas.sh`.

---

## 6. AWS Step Functions - Orquestração

### 6.1 Processamento (`listing-media-processing-sm-staging`)
1. **ValidateRawAssets** (`listing-media-validate-staging`) — normaliza assets, garante que cada `rawKey` existe e injeta `hasVideos`, `traceparent` e erros de validação no payload.
2. **ParallelProcessing** — executa geração de thumbnails e, condicionalmente, job MediaConvert para vídeos.
3. **ConsolidateResults** — agrega resultados, atribui `errorCode` (`VALIDATION_ERROR`, `THUMBNAIL_PROCESSING_FAILED`, etc.) e repassa `traceparent`.
  - A Lambda `Consolidate` preenche `processedKey` escolhendo a melhor resolução disponível (`large` → `medium` → `small` → `thumbnail`) e sempre inclui `thumbnailKey` quando gerado.
4. **FinalizeAndCallback** — envia a estrutura unificada para `listing-media-callback-dispatch-staging`, que chama o backend.

### 6.2 Finalização (`listing-media-finalization-sm-staging`)
1. **CreateZipBundle** (`listing-media-zip-staging`) recebe `MediaFinalizationInput` com `rawKey` dos assets e grava `/<listingIdentityId>/processed/zip/listing-media.zip`.
2. **FinalizeAndCallback** — responde com `provider=STEP_FUNCTIONS_FINALIZATION`, `status=SUCCEEDED`, `zipBundles` e `assetsZipped`.
3. **ReportFailure** — em caso de erro, envia `status=FINALIZATION_FAILED` e detalha o motivo no callback.

> **Correlação Banco ↔ Step Functions (Nov/2025):** O serviço `ProcessMedia` registra o job e armazena o `executionArn` do workflow de processamento. `CompleteMedia` cria um novo job, dispara `listing-media-finalization-sm-staging` via `workflow.StartMediaFinalization` e atualiza `media_processing_jobs.external_id` com o ARN de finalização. Assim, qualquer `executionArn` encontrado em `service.media.complete.started_zip` pode ser rastreado no banco e vice-versa.

---

## 7. Salvaguardas Operacionais (Jan/2026)

- **Reconciliação de jobs travados:** worker interno roda a cada 5 minutos e marca como `FAILED` qualquer job `RUNNING` sem callback cujo `started_at` tenha ultrapassado `media_processing.workflow.stuck_job_timeout_minutes` (default 30). Assets em `PROCESSING` são igualmente marcados como `FAILED` para evitar pendência.
- **Pré-checagem de uploads:** antes de publicar no SQS, o backend agora faz `HeadObject` nos `rawKey` envolvidos; se o objeto não existir ou estiver inacessível, retorna erro de validação imediatamente (não enfileira nem altera status dos assets).
