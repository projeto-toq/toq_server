# Plano de Implementação — Media Processing TOQ Server

## 1. Diagnóstico e Escopo
- **Domínios afetados:** `internal/core/model/listing_model`, `internal/core/service/listing_service`, `internal/adapter/left/http/handlers/listing_handlers`, `internal/adapter/right/mysql/listing`, `internal/core/port/right/repository/listingrepository`, além da criação do novo domínio `media_processing` (model/service/repository/handlers/DTOs). Evidências: handlers de fotos em `listing_handlers` retornam HTTP 501; não há repositórios/serviços para lotes ou assets; S3 adapter atual só cobre uploads genéricos.
- **Motivação:** implementar o fluxo descrito em `docs/media_processing_guide.md` assegurando que listings em `StatusPendingPhotoProcessing` avancem apenas após pipeline assíncrono concluir. Ausência atual impede fotógrafos de subir mídias e bloqueia aprovação de anúncios.
- **Impacto:** novas tabelas (`listing_media_batches`, `listing_media_assets`, `listing_media_jobs`), novos ports/adapters, novos serviços core, 5 endpoints REST (`/api/v2/listings/media/*`, incluindo `/uploads/retry`), suporte a novos tipos de mídia (`PROJECT_DOC`, `PROJECT_RENDER`) com uploads feitos pelo owner em listings de empreendimento, integração S3/SQS/Step Functions, métricas e logs adicionais. Alteração disruptiva — todos os listings serão resetados.
- **Melhorias incorporadas:** autorização restrita ao fotógrafo do booking ativo ou ao owner (para `project/*`), reprocessamento com novo job, suporte a múltiplos batches simultâneos (com soft delete lógico), métricas específicas, checklist completo de observabilidade.

## 2. Code Skeletons (seguir Seção 8 do guia)
### 2.1 Handlers (`internal/adapter/left/http/handlers/listing_handlers/`)
- `create_upload_batch_handler.go`
- `complete_upload_batch_handler.go`
- `get_batch_status_handler.go`
- `list_download_urls_handler.go`
- `retry_media_batch_handler.go` (reprocessamento)
  
Cada handler conterá:
```go
// CreateUploadBatch godoc
// @Summary Request upload URLs for listing media batch
// @Tags Listings Media
// @Accept json
// @Produce json
// @Param request body listingdto.CreateUploadBatchRequest true "Upload manifest"
// @Success 201 {object} listingdto.CreateUploadBatchResponse
// @Failure 4xx {object} httperrors.HTTPError
// @Router /api/v2/listings/media/uploads [post]
func (h *ListingHandler) CreateUploadBatch(c *gin.Context) { /* … */ }
```

### 2.2 Services (`internal/core/service/media_processing_service/`)
Arquivos:
- `media_processing_service.go` (interface + struct + `New`)
- `create_upload_batch.go`
- `complete_upload_batch.go`
- `get_batch_status.go`
- `list_download_urls.go`
- `retry_media_batch.go`
- `handle_processing_callback.go`

Template:
```go
// CreateUploadBatch issues signed URLs for photographers to upload raw media.
func (s *mediaProcessingService) CreateUploadBatch(ctx context.Context, input CreateUploadBatchInput) (CreateUploadBatchOutput, error) {
    ctx, spanEnd, err := utils.GenerateTracer(ctx)
    if err != nil {
        return CreateUploadBatchOutput{}, derrors.Infra("Failed to generate tracer", err)
    }
    defer spanEnd()

    logger := utils.LoggerFromContext(ctx)
    tx, err := s.globalService.StartTransaction(ctx)
    if err != nil {
        logger.Error("service.media.create_batch.tx_start_error", "err", err)
        utils.SetSpanError(ctx, err)
        return CreateUploadBatchOutput{}, derrors.Infra("Failed to start transaction", err)
    }
    defer func() {
        if err != nil { _ = s.globalService.RollbackTransaction(ctx, tx) }
    }()

    // … regras de autorização (fotógrafo ou owner para projetos), validações, persistência via repo …

    if err := s.globalService.CommitTransaction(ctx, tx); err != nil {
        logger.Error("service.media.create_batch.tx_commit_error", "err", err)
        utils.SetSpanError(ctx, err)
        return CreateUploadBatchOutput{}, derrors.Infra("Failed to commit transaction", err)
    }

    return output, nil
}
```

### 2.3 Repositórios (`internal/core/port/right/repository/mediaprocessingrepository/` + `internal/adapter/right/mysql/media_processing/`)
Métodos mínimos no Port:
- `CreateBatch`
- `GetBatchByID`
- `ListBatchesByListing`
- `UpdateBatchStatus`
- `SoftDeleteBatch`
- `CreateAsset`
- `UpsertAssets`
- `ListAssetsByBatch`
- `CreateJob`
- `UpdateJob`
- `ListJobsByBatch`

Adapter MySQL segue InstrumentedAdapter e um arquivo por método (ex.: `create_batch.go`). Inserir comentários destacando queries, colunas explícitas e conversão via `converters/`.

### 2.4 DTOs (`internal/adapter/left/http/dto/listing_dto/`)
- `CreateUploadBatchRequest/Response` (inclui campos `title`, `sequence`)
- `CompleteUploadBatchRequest/Response`
- `GetBatchStatusRequest/Response`
- `ListDownloadUrlsRequest/Response`
- `RetryMediaBatchRequest/Response`
- `MediaProcessingCallbackRequest`

Exemplo:
```go
// CreateUploadBatchRequest defines the payload used to request signed URLs.
type CreateUploadBatchRequest struct {
  ListingID      uint64                         `json:"listingId" binding:"required"`
  BatchReference string                         `json:"batchReference" binding:"required,max=100"`
  Files          []CreateUploadBatchFileRequest `json:"files" binding:"required,dive"` // aceita mediaType = PHOTO_*|VIDEO_*|PROJECT_DOC|PROJECT_RENDER
}
```

### 2.5 Entities/Converters
- `internal/adapter/right/mysql/media_processing/entities/batch_entity.go` com `sql.Null*`
- `asset_entity.go`, `job_entity.go` (asset inclui campos para `variant_metadata_json` e novos tipos `PROJECT_DOC|PROJECT_RENDER`)
- Conversores `entity_to_domain` e `domain_to_entity`
- Converters HTTP `listing_converters/media_batch_converters.go` para mapear DTO ⇄ Service

## 3. Estrutura de Diretórios (Regra de Espelhamento)
```
internal/
  core/
    model/
      media_processing_model/
        media_batch_interface.go
    port/
      right/
        repository/
          mediaprocessingrepository/
            media_processing_repo_port.go
    service/
      media_processing_service/
        media_processing_service.go
        create_upload_batch.go
        ...
  adapter/
    left/http/handlers/listing_handlers/
      create_upload_batch_handler.go
      ...
    left/http/dto/listing_dto/
      media_processing_dto.go
    right/mysql/media_processing/
      media_processing_adapter.go
      create_batch.go
      ...
    right/aws_s3/
      media_processing_s3.go (novos helpers)
    right/aws_sqs/media_processing/
      media_processing_queue_adapter.go
```

## 4. Ordem de Execução
1. **Modelagem & Schema:** definir structs domínio, entidades MySQL e script SQL.
2. **Port & Adapter MySQL:** criar interface + adapter com todos os métodos.
3. **Service:** implementar métodos públicos e regras de autorização/validação.
4. **Adapters externos:** estender S3 e criar SQS/Callback port.
5. **Handlers/DTOs/Converters:** expor endpoints `/api/v2/listings/media/*`.
6. **Factory/DI:** registrar novos adapters/serviços nas fábricas.
7. **Observabilidade & Docs:** métricas, logging, Swagger, checklist final.

## 5. Checklist de Conformidade
- [ ] Arquitetura hexagonal (Seção 1)
- [ ] Regra de Espelhamento Port ↔ Adapter (Seção 2.1)
- [ ] InstrumentedAdapter em todos os repos (Seção 7.3)
- [ ] Transações via `globalService` (Seção 7.1)
- [ ] Tracing/logging/erros (`utils.GenerateTracer`, `utils.SetSpanError`) (Seções 5, 7, 9)
- [ ] Documentação completa: Godoc + Swagger + DTO comments (Seção 8)
- [ ] Sem anti-padrões (Seção 14) — um método por arquivo, nada de SELECT *, sem lógica HTTP no service
- [ ] Observabilidade reforçada (métricas media_*, logs com `listing_id`, `batch_id`, `job_id`)

## 6. Regras Complementares Confirmadas
1. **Autorização:** fotógrafo vinculado ao booking ativo controla fotos/vídeos; owner autenticado pode criar/reprocessar lotes de `project/doc` e `project/render` em listings de empreendimento.
2. **Idempotência:** não suportar `Idempotency-Key`; controles via status/batch explícitos.
3. **Reprocessamento:** endpoint dedicado `POST /api/v2/listings/media/uploads/retry` gera novo registro em `listing_media_jobs` e reenvia para pipeline.
4. **Múltiplos batches:** permitir histórico completo; aplicar soft delete (`deleted_at` ou flag) para batches substituídos manualmente.
5. **Métricas obrigatórias:** `media_batches_created_total`, `media_processing_duration_seconds`, `media_processing_inflight_jobs`, `media_processing_dlq_messages_total`, `media_batches_failed_total`.

## 7. Recursos AWS Necessários (Cloud Admin)
- **S3 Bucket `toq-listing-medias`**: habilitar SSE-KMS, versionamento opcional; pastas `/{listingId}/{stage}/{type}/`; políticas para fotógrafos via backend (presigned) e roles internas. Habilitar `mediaType` adicional `project/doc` e `project/render`, disponível apenas para owners de listings de empreendimento.
- **SQS Standard Queue `listing-media-processing`** + **DLQ `listing-media-processing-dlq`** com redrive policy; atributos extras: `Traceparent`.
- **AWS Step Functions State Machine** `listing-media-processing-sm` com etapas: validate → thumbnails → transcode → zip → finalize.
- **Lambda Functions** (todas com IAM mínimo):
  1. `listing-media-validate`
  2. `listing-media-thumbnails`
  3. `listing-media-zip`
  4. `listing-media-callback-dispatch`
- **AWS Elemental MediaConvert Presets** para vídeo (1080p/720p MP4 + HLS).
- **EventBridge Rule** ou **API Gateway** para receber callbacks de MediaConvert/Step Functions e acionar Lambda #4.
- **IAM Roles/Policies** específicas para Step Functions, Lambdas, MediaConvert e acesso restrito ao bucket/pastas.
- **CloudWatch Logs/Metrics + Alarms** (duração de processamento p95 >10min, mensagens DLQ >0, falhas consecutivas em Lambdas).

## 8. Sugestões de Melhoria para `media_processing_guide.md`
- Explicitar regra de autorização (somente fotógrafo do booking ativo ou administradores) na seção 7.
- Registrar decisão de não suportar `Idempotency-Key` e explicar o controle por status do lote.
- Adicionar subseção de reprocessamento incluindo endpoint, fluxo Step Functions e novo registro em `listing_media_jobs`.
- Detalhar política de múltiplos batches + soft delete para manter histórico porém ocultar batches obsoletos.
- Incluir tabela de métricas/alertas obrigatórios (nomes exatos, descrição e limiares).

---
Este arquivo servirá como referência viva para implementação incremental e validação futura. Ajustes adicionais podem ser inseridos aqui antes do desenvolvimento efetivo.
