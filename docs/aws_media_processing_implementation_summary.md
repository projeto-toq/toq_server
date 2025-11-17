# TOQ Server - Media Processing Infrastructure
## Recursos AWS Criados - Staging Environment
**Data:** 17 de Novembro de 2025  
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
  - Lambda: InvokeFunction (todas as 4 Lambdas)
  - MediaConvert: CreateJob, GetJob
  - SQS: SendMessage
  - CloudWatch Logs, X-Ray

### 4.4 MediaConvert Role
- **Nome:** `toq-media-processing-mediaconvert-staging`
- **ARN:** `arn:aws:iam::058264253741:role/toq-media-processing-mediaconvert-staging`
- **Assume Role:** mediaconvert.amazonaws.com
- **Permissões:**
  - S3: GetObject, PutObject (toq-listing-medias)
  - KMS: Encrypt, Decrypt, GenerateDataKey

---

## 5. AWS Lambda - Funções (Placeholder)

### 5.1 Validate
- **Nome:** `listing-media-validate-staging`
- **ARN:** `arn:aws:lambda:us-east-1:058264253741:function:listing-media-validate-staging`
- **Runtime:** Node.js 20.x
- **Memory:** 512 MB | **Timeout:** 60s
- **Propósito:** Validar assets no S3 (checksum, bytes, existência)

### 5.2 Thumbnails
- **Nome:** `listing-media-thumbnails-staging`
- **ARN:** `arn:aws:lambda:us-east-1:058264253741:function:listing-media-thumbnails-staging`
- **Runtime:** Node.js 20.x
- **Memory:** 2048 MB | **Timeout:** 300s
- **Propósito:** Gerar thumbnails com Sharp (layer necessária)

### 5.3 ZIP
- **Nome:** `listing-media-zip-staging`
- **ARN:** `arn:aws:lambda:us-east-1:058264253741:function:listing-media-zip-staging`
- **Runtime:** Python 3.12
- **Memory:** 3008 MB | **Timeout:** 900s | **Storage:** 2 GB
- **Propósito:** Compactar assets processados

### 5.4 Callback Dispatch
- **Nome:** `listing-media-callback-dispatch-staging`
- **ARN:** `arn:aws:lambda:us-east-1:058264253741:function:listing-media-callback-dispatch-staging`
- **Runtime:** provided.al2023 (custom)
- **Memory:** 256 MB | **Timeout:** 60s
- **Propósito:** Enviar callback ao backend

**Variáveis de Ambiente (todas):**
- `ENV=staging`
- `MEDIA_BUCKET=toq-listing-medias`
- `TRACE_HEADER_KEY=traceparent`
- `CALLBACK_URL=https://api-staging.toq.com/internal/media-processing/callback` (só callback)

---

## 6. AWS Step Functions - State Machine

### State Machine
- **Nome:** `listing-media-processing-sm-staging`
- **ARN:** `arn:aws:states:us-east-1:058264253741:stateMachine:listing-media-processing-sm-staging`
- **Tipo:** STANDARD
- **Status:** ACTIVE
- **Logging:** ALL (CloudWatch Logs)
- **Tracing:** X-Ray habilitado
- **Log Group:** `/aws/stepfunctions/listing-media-processing-sm-staging` (retenção 30 dias)

### Workflow States
1. **ValidateRawAssets** → Valida assets
2. **ParallelProcessing** →
   - Branch 1: GenerateThumbnails
   - Branch 2: CheckVideoProcessing → ProcessVideos (MediaConvert) ou Skip
3. **CreateZipBundle** → Compacta assets
4. **FinalizeAndCallback** → Notifica backend
5. **ProcessingSucceeded** / **ValidationFailed** / **ProcessingFailed**

---

## 7. AWS MediaConvert - Transcodificação

### Queue
- **Nome:** `toq-media-queue-staging`
- **ARN:** `arn:aws:mediaconvert:us-east-1:058264253741:queues/toq-media-queue-staging`
- **Status:** ACTIVE
- **Pricing:** ON_DEMAND

### Job Template
- **Nome:** `toq-media-video-preset-staging`
- **ARN:** `arn:aws:mediaconvert:us-east-1:058264253741:jobTemplates/toq-media-video-preset-staging`
- **Category:** toq-media-processing
- **Outputs:**
  - 1080p: H.264 (5 Mbps max), AAC 128 kbps, MP4
  - 720p: H.264 (3 Mbps max), AAC 96 kbps, MP4

---

## 8. Amazon SNS - Notificações

### Topic
- **Nome:** `toq-media-alerts-staging`
- **ARN:** `arn:aws:sns:us-east-1:058264253741:toq-media-alerts-staging`
- **Subscription:** giulio.alfieri@gmail.com (confirmar email)

---

## 9. Amazon EventBridge - Eventos

### Rule MediaConvert
- **Nome:** `toq-media-mediaconvert-status-staging`
- **ARN:** `arn:aws:events:us-east-1:058264253741:rule/toq-media-mediaconvert-status-staging`
- **Padrão:** Eventos de MediaConvert (COMPLETE/ERROR) da queue staging
- **Target:** SNS topic `toq-media-alerts-staging`

---

## 10. Amazon CloudWatch - Alarmes

### Alarme 1: DLQ Messages
- **Nome:** `listing-media-dlq-has-messages-staging`
- **Métrica:** ApproximateNumberOfMessagesVisible (SQS DLQ)
- **Threshold:** > 0
- **Período:** 5 minutos
- **Ação:** Notifica SNS

### Alarme 2: Step Functions Failures
- **Nome:** `stepfunctions-execution-failed-staging`
- **Métrica:** ExecutionsFailed (Step Functions)
- **Threshold:** > 0
- **Período:** 5 minutos
- **Ação:** Notifica SNS

### Alarme 3: Lambda Errors
- **Nome:** `lambda-errors-media-processing-staging`
- **Métrica:** Errors (Lambda - todas funções)
- **Threshold:** > 5
- **Período:** 5 minutos
- **Ação:** Notifica SNS

---

## 11. Próximos Passos - Implementação

### Para o Time de Backend (Go):
1. Implementar código real nas Lambdas (atualmente placeholders)
2. Adicionar Sharp layer na Lambda `listing-media-thumbnails-staging`
3. Configurar variáveis de ambiente no `env.yaml`:
   ```yaml
   aws:
     region: us-east-1
     kms_key_id: 2fd6d9e5-dc43-4275-ba77-79c99cad509c
     media_bucket: toq-listing-medias
     sqs_queue_url: https://sqs.us-east-1.amazonaws.com/058264253741/listing-media-processing-staging
     step_functions_arn: arn:aws:states:us-east-1:058264253741:stateMachine:listing-media-processing-sm-staging
   ```
4. Implementar endpoints:
   - `POST /internal/media-processing/callback` (recebe callbacks)
   - `POST /api/v1/listings/{id}/media/batch` (inicia upload)
5. Aplicar migrations SQL (criar tabelas `listing_media_*`)

### Para o Time de Cloud/DevOps:
1. Revisar e ajustar policies IAM se necessário
2. Configurar backups do bucket S3 se aplicável
3. Monitorar custos (MediaConvert é pay-per-use)
4. Configurar dashboards Grafana/CloudWatch personalizados
5. Implementar housekeeping automatizado (fase futura)

### Para o Time de QA:
1. Testar workflow completo end-to-end
2. Validar alarmes e notificações
3. Testar cenários de falha (DLQ, retries)
4. Verificar tracing X-Ray funcionando

---

## 12. Custos Estimados (Staging)

- **S3:** ~$0.023/GB/mês + requests
- **SQS:** Grátis até 1M requests/mês
- **Lambda:** Grátis até 1M requests + 400k GB-s/mês
- **Step Functions:** $0.025 por 1000 transições de estado
- **MediaConvert:** ~$0.015/minuto (1080p SD) + ~$0.0075/min (720p SD)
- **KMS:** $1/mês por chave + $0.03 per 10k requests
- **CloudWatch Logs:** $0.50/GB ingestão, $0.03/GB armazenamento

**Estimativa mensal (uso moderado):** $20-50 USD

---

## 13. Segurança e Compliance

✅ Criptografia em repouso (KMS) em todos os recursos  
✅ Criptografia em trânsito (TLS)  
✅ Block Public Access no S3  
✅ Least privilege IAM policies  
✅ Logging e auditoria habilitados  
✅ X-Ray tracing para troubleshooting  
✅ Alarmes para detecção de falhas  

---

## 14. Credenciais Locais

**Arquivo:** `/codigos/go_code/toq_server/configs/aws_credentials`  
**Symlink:** `~/.aws/credentials → configs/aws_credentials`  
**Perfil:** `admin`  

⚠️ **Não commitar** este arquivo no Git!

---

**Documento gerado automaticamente**  
Para dúvidas: contatar time de Cloud/Backend
