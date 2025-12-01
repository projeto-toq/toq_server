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
- **Raw (Upload):** `/{listingId}/raw/{mediaType}/{uuid}.{ext}`
  - *Nota:* Segmentos de data (YYYY-MM-DD) foram removidos para simplificar a estrutura.
- **Processed (Thumbnails):** `/{listingId}/processed/{mediaType}/{size}/{uuid}.{ext}`
  - Tamanhos: `thumbnail` (200px), `small` (400px), `medium` (800px), `large` (1200px).
- **Zip Bundles:** `/{listingId}/processed/zip/{listingId}.zip`
  - Conteúdo interno do Zip é limpo (sem prefixos `processed/` ou datas).

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
   - Cria arquivo ZIP contendo todas as mídias processadas.
   - Limpa estrutura de pastas interna (remove prefixos de sistema).
4. **Consolidate (`listing-media-consolidate-staging`)**: Agrega resultados do processamento paralelo.
5. **Callback (`listing-media-callback-staging`)**: Envia webhook de volta ao backend.

### Runtime
- **Runtime:** `provided.al2023`
- **Build:** Binários compilados estaticamente via `./scripts/build_lambdas.sh`.
