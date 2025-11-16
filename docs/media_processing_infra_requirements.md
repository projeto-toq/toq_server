# Requisitos de Infraestrutura — Media Processing TOQ Server

Documento destinado aos times de Cloud Admin e DBA para provisionamento dos recursos necessários ao novo pipeline de processamento de mídias descrito em `docs/media_processing_guide.md`.

## 1. AWS — Recursos Obrigatórios

### 1.1 Amazon S3 — Bucket `toq-listing-medias`
- **Status**: já existente (confirmar). Caso contrário, criar bucket único por ambiente (`toq-listing-medias-${env}`) com as políticas abaixo.
- **Configurações**
  - Criptografia padrão: `SSE-KMS` (CMK específico do time TOQ).
  - Versionamento habilitado (mantém histórico para reprocessamento).
  - Bloqueio de acesso público: ativado.
  - Lifecycle Rule opcional para mover `raw/` para S3 Glacier após 180 dias.
- **Layout de chave**: `/{listingId}/{stage}/{mediaType}/{YYYY-MM-DD}/{filename}`
  - `stage`: `raw`, `processed`, `zip`, `thumb`.
  - `mediaType`: `photo/vertical`, `photo/horizontal`, `video/vertical`, `video/horizontal`, `project/doc`, `project/render`.
- **Fluxo owner**: quando `project/*`, o upload é disparado diretamente pelo owner autenticado (sem fotógrafo). Garante-se segregação criando política IAM específica que só permite `project/*` para requests assinadas pelo backend.
- **Políticas IAM**
  - Role `toq-media-processing-backend` com permissões `s3:PutObject`, `s3:GetObject`, `s3:HeadObject`, `s3:DeleteObject`, `s3:GetObjectAttributes` no bucket + prefixo.
  - Negar `s3:DeleteObject` via API pública; deletes acontecem apenas por operações backend (soft delete preferencial).
- **Logs/Auditoria**
  - Habilitar Server Access Logging para bucket dedicado `toq-logs`.

### 1.2 AWS KMS
- Criar/usar chave gerenciada `alias/toq-media-processing`.
- Permissões: backend app role + Lambdas do pipeline + MediaConvert service-linked-role.
- Habilitar automatic rotation (12 meses).

### 1.3 Amazon SQS
- **Fila principal**: `listing-media-processing`
  - Tipo: Standard.
  - Visibility timeout: 60 s (ajustável conforme Step Functions delay).
  - Encryption: KMS (mesmo alias da seção 1.2).
  - Atributos customizados obrigatórios nas mensagens: `Traceparent`, `ListingId`, `BatchId`.
- **DLQ**: `listing-media-processing-dlq`
  - Redrive policy: mover após 5 tentativas.
  - Configurar alerta CloudWatch para mensagens > 0 por 5 min.

### 1.4 AWS Step Functions
- **State machine**: `listing-media-processing-sm`
  - Tipo Standard (não Express) para rastreabilidade.
  - Estados sugeridos: `ValidateRawAssets` → `ParallelProcessing` (branches `ImageThumbnails`, `VideoTranscode`) → `ZipBundle` → `FinalizeMetadata` → `CallbackBackend`.
  - Timeout da execução: 1 hora.
  - Input esperado: `batchId`, `listingId`, `assets[]`, `traceparent`.
- **IAM Role**: permitir `lambda:InvokeFunction`, `mediaconvert:CreateJob`, `sqs:SendMessage`, `logs:CreateLogDelivery`.

### 1.5 AWS Lambda
| Função | Runtime | Descrição | Entradas/Saídas | Observações |
| --- | --- | --- | --- | --- |
| `listing-media-validate` | `nodejs20.x` (ou Go 1.x) | Confere existência no S3, checksum, bytes, status do batch e sinaliza Step Functions | Entrada: manifest do batch. Saída: manifest normalizado | Necessita `s3:GetObjectAttributes`, `s3:HeadObject` |
| `listing-media-thumbnails` | `nodejs20.x` + Sharp | Gera thumbnails, escreve em `processed/thumb/` e retorna metadados atualizados | Entrada: lista de fotos válidas | Acesso a camada compartilhada `libvips` |
| `listing-media-zip` | `python3.12` ou Go | Compacta assets processados em ZIPs completos/parciais | Entrada: lista de assets processados | Requer espaço temporário (EFS ou /tmp >= 1 GB) |
| `listing-media-callback-dispatch` | Go 1.x | Envia payload final para endpoint `/internal/media-processing/callback` | Entrada: status final do Step Functions | Assina requisições com token interno + headers Traceparent |

Configurações comuns:
- Variáveis de ambiente: `ENV`, `MEDIA_BUCKET`, `CALLBACK_URL`, `TRACE_HEADER_KEY`.
- Observabilidade: habilitar X-Ray e logs no CloudWatch (`/aws/lambda/<function>`).

### 1.6 AWS Elemental MediaConvert
- Criar **Template** `toq-media-video-preset` com saídas MP4 1080p e 720p + thumbnails opcionais.
- Criar **Queue** dedicada (`toq-media-queue`) com prioridade normal.
- Garantir role `MediaConvert_Default_Role` com acesso ao bucket e KMS.

### 1.7 EventBridge / API Gateway
- **Callback**: expor endpoint privado `POST /media-processing/callback` (API Gateway + VPC Link) chamando Lambda `listing-media-callback-dispatch`.
- **Eventos MediaConvert**: regra EventBridge `toq-media-mediaconvert-status` filtrando eventos `COMPLETED`/`ERROR` e redirecionando para Step Functions ou Lambda de reconciliação.

### 1.8 IAM Separação de Funções
- `toq-media-processing-backend`: usado pela aplicação Go para gerar URLs e consultar S3.
- `toq-media-processing-lambda`: assume pelas Lambdas; inclui permissões de SQS, Step Functions, S3, Logs.
- `toq-media-processing-admin`: role break-glass para Cloud Ops (acesso manual ao bucket e filas).

### 1.9 Observabilidade AWS
- **Logs**: centralizar no CloudWatch Log Group `/aws/step-functions/listing-media-processing-sm` e `/aws/lambda/listing-media-*`.
- **Métricas**: criar dashboards com `ExecutionsSucceeded`, `ExecutionsFailed`, `ApproximateNumberOfMessagesVisible` (fila principal e DLQ), `Errors` (por Lambda).
- **Alarmes sugeridos**
  - Execuções falhas > 0 em 5 min (Step Functions).
  - Mensagens na DLQ > 0 (5 min).
  - Lambda `listing-media-callback-dispatch` com `Errors` > 0 (threshold 1).
  - MediaConvert job falho (EventBridge envia para SNS `toq-media-alerts`).

## 2. Banco de Dados — Especificações para DBA

As alterações devem ser aplicadas no schema `toq_server`. Segue DDL sugerida (ajustar tipos conforme padrão do cluster). Conforme alinhado com o time de auditoria, campos `created_by`, `updated_by`, `deleted_by`, `created_at`, `updated_at` e `deleted_at` foram removidos de todas as tabelas abaixo sempre que eram utilizados apenas para rastreabilidade. Mantivemos somente `deleted_at` em `listing_media_batches` porque ele é necessário para o fluxo de soft delete/controlar housekeeping futuro.

### 2.1 Tabela `listing_media_batches`
```sql
CREATE TABLE listing_media_batches (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  listing_id BIGINT UNSIGNED NOT NULL,
  photographer_user_id BIGINT UNSIGNED NOT NULL,
  status ENUM('PENDING_UPLOAD','RECEIVED','PROCESSING','FAILED','READY') NOT NULL,
  upload_manifest_json JSON NOT NULL,
  processing_metadata_json JSON NULL,
  received_at DATETIME(6) NULL,
  processing_started_at DATETIME(6) NULL,
  processing_finished_at DATETIME(6) NULL,
  error_code VARCHAR(50) NULL,
  error_detail TEXT NULL,
  deleted_at DATETIME(6) NULL,
  PRIMARY KEY (id),
  KEY idx_listing_status (listing_id, status),
  KEY idx_listing_recent (listing_id, id DESC),
  CONSTRAINT fk_batches_listing FOREIGN KEY (listing_id) REFERENCES listing_identities(id) ON DELETE CASCADE,
  CONSTRAINT fk_batches_photographer FOREIGN KEY (photographer_user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 2.2 Tabela `listing_media_assets`
```sql
CREATE TABLE listing_media_assets (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  batch_id BIGINT UNSIGNED NOT NULL,
  asset_type ENUM('PHOTO_VERTICAL','PHOTO_HORIZONTAL','VIDEO_VERTICAL','VIDEO_HORIZONTAL','THUMBNAIL','ZIP','PROJECT_DOC','PROJECT_RENDER') NOT NULL,
  orientation ENUM('VERTICAL','HORIZONTAL') NULL,
  source_key VARCHAR(512) NOT NULL,
  processed_key VARCHAR(512) NULL,
  thumbnail_key VARCHAR(512) NULL,
  checksum_sha256 CHAR(64) NOT NULL,
  content_type VARCHAR(100) NOT NULL,
  bytes BIGINT UNSIGNED NOT NULL,
  resolution VARCHAR(20) NULL,
  duration_seconds INT UNSIGNED NULL,
  title VARCHAR(255) NULL,
  sequence INT UNSIGNED NULL,
  variant_metadata_json JSON NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uk_batch_sequence (batch_id, sequence),
  KEY idx_batch_type (batch_id, asset_type),
  CONSTRAINT fk_assets_batch FOREIGN KEY (batch_id) REFERENCES listing_media_batches(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 2.3 Tabela `listing_media_jobs`
```sql
CREATE TABLE listing_media_jobs (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  batch_id BIGINT UNSIGNED NOT NULL,
  external_job_id VARCHAR(255) NOT NULL,
  provider ENUM('STEP_FUNCTIONS','MEDIACONVERT') NOT NULL,
  status VARCHAR(50) NOT NULL,
  input_payload_json JSON NULL,
  output_payload_json JSON NULL,
  started_at DATETIME(6) NULL,
  finished_at DATETIME(6) NULL,
  PRIMARY KEY (id),
  KEY idx_batch_job (batch_id, created_at DESC),
  CONSTRAINT fk_jobs_batch FOREIGN KEY (batch_id) REFERENCES listing_media_batches(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 2.4 Considerações adicionais para DBA
- Garantir charset `utf8mb4` e collation `utf8mb4_unicode_ci` (padrão do projeto).
- Criar view ou procedure auxiliar para reportar batches ativos, se necessário.
- Ajustar backups/replicação para incluir novas tabelas.
- Conferir permissões de usuário de aplicação (`toq_app`) para `SELECT/INSERT/UPDATE/DELETE` nessas tabelas.
- Caso haja replicação assíncrona, validar impacto de JSON columns (InnoDB >= 5.7.9).

## 3. Dependências Entre Times
- **Backend** aguarda confirmação de provisionamento AWS + migração SQL antes de habilitar rotas.
- **Cloud Admin** deve fornecer outputs (ARNs de Step Functions, filas, roles) para parametrizar `env.yaml`.
- **DBA** deve comunicar janela de manutenção e confirmar criação em sandbox/homolog antes de produção.

## 4. Housekeeping & Limpeza de Mídias (Entrega Futuras)
- **Objetivo:** remover fisicamente objetos `raw/processed/zip` após o tempo de retenção e processar lotes marcados com `deleted_at`. Esta automação **não** faz parte do escopo imediato.
- **Infra necessária:** job dedicado (Lambda agendada, Step Functions ou ECS Fargate) que consuma lista de lotes elegíveis via endpoint interno. O job deve ter permissões restritas em `toq-listing-medias` e em `listing_media_*` para atualizar flags de limpeza.
- **Ações pendentes:** Cloud/Admin definir arquitetura e custos; Backend expor API de housekeeping na fase correspondente. Até lá, nenhum componente deve remover objetos automaticamente.

---
Para dúvidas ou alterações, registrar comentários diretamente neste arquivo e sincronizar com o time de backend responsável pelo TOQ Server.
