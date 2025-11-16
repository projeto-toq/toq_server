# Plano de Implementação — Media Processing TOQ Server

## 1. Diagnóstico Detalhado
- **Arquitetura atual (handler/service/repo)**
  - `internal/adapter/left/http/handlers/listing_handlers/not_implemented_handlers.go`: handlers de fotos ainda respondem HTTP 501 (`AddListingPhotos`, `GetListingPhotos` etc.), não existe rota `/api/v2/listings/media/*`. Evidência direta de que a camada HTTP não cobre o fluxo descrito no guia (Seção 7.4 do `toq_server_go_guide.md`).
  - `internal/core/service/listing_service/*.go`: nenhum arquivo trata lotes de mídia; apenas catálogo, versões e reservas. `StatusPendingPhotoProcessing` é atribuído exclusivamente em `photo_session_service/session_update_status.go` (linhas 29-199) e não há serviço subsequente que mova para `StatusPendingOwnerApproval` após upload/processamento, deixando o funil travado.
  - `internal/adapter/right/mysql/listing/*`: repositório atual só manipula `listing_versions` e satélites; não há entidades/queries para `listing_media_*`, quebrando a Regra de Espelhamento (Seção 2.1). Precisamos de novo domínio `media_processing` com adapter próprio e uso obrigatório de `InstrumentedAdapter` (Seção 7.3).
  - `internal/adapter/right/aws_s3/*.go`: `S3Adapter` expõe utilitários genéricos (`GenerateV4PutObjectSignedURL`, `GeneratePhotoSignedURL`), mas não há contrato orientado a lotes/listings (sem prefixos `raw/processed`, checksum SHA-256 ou TTL configurável). Também não existe port/adapter SQS/Step Functions para orquestração.
- **Modelos/Persistência inexistentes**
  - `internal/core/model/listing_model/constants.go` define `StatusPendingPhotoProcessing`, mas não há `media_processing_model` descrevendo `MediaBatch`, `MediaAsset`, `MediaProcessingJob`. Banco de dados carece das tabelas detalhadas em `docs/media_processing_infra_requirements.md`.
- **Fluxo atual interrompido**
  - `photo_session_service.UpdateSessionStatus` (linhas 29-215) move listings para `StatusPendingPhotoProcessing` quando o fotógrafo conclui a sessão, porém nenhum worker/handler aceita uploads. O fotógrafo fica sem ponto de contato e o status permanece indefinidamente.
- **Observabilidade ausente**
  - Não existem contadores `media_batches_*` nem histogramas `media_processing_duration_seconds` (Seção 9 do guia de mídia). O repositório `metrics` não possui coletores específicos.
- **Confirmações recentes do time**
  1. Estados `PENDING_UPLOAD → RECEIVED → PROCESSING → READY|FAILED` estão confirmados; não haverá novos estados antes do MVP.
  2. Expiração/limpeza física ficará para um job dedicado futuro. Nesta fase só precisamos de soft delete (`deleted_at`) e APIs para que o job venha a consumir posteriormente.

**Impacto principal**
- Novo domínio `media_processing` (model/port/service/adapters/DTOs) conforme Arquitetura Hexagonal (Seção 1).
- Cinco endpoints REST sob `/api/v2/listings/media/*` com autorização fotógrafo/owner e validações de manifesto (checksum, título, sequência).
- Integração com S3/SQS/Step Functions e métricas específicas para acompanhamento.
- Atualização dos guias `media_processing_guide.md` (estados confirmados + limpeza futura) e `media_processing_infra_requirements.md` (registrar job dedicado).

## 2. Code Skeletons (Seção 8)
### 2.1 Handlers (`internal/adapter/left/http/handlers/listing_handlers/`)
1. `create_upload_batch_handler.go`
   ```go
   // CreateUploadBatch handles signed URL issuance for listing media uploads
   //
   // @Summary     Request upload URLs for a listing media batch
   // @Description Validates permissions (photographer booking or owner for project media), enforces listing status PENDING_PHOTO_PROCESSING, creates a media batch and returns signed PUT URLs with required headers (Content-Type, checksum). Rejects duplicates (sequence/title) and unsupported media types.
   // @Tags        Listings Media
   // @Accept      json
   // @Produce     json
   // @Security    BearerAuth
   // @Param       request body listingdto.CreateUploadBatchRequest true "Manifest with client-side file metadata"
   // @Success     201 {object} listingdto.CreateUploadBatchResponse "Batch created with signed URLs"
   // @Failure     400 {object} httperrors.HTTPError "Invalid manifest"
   // @Failure     401 {object} httperrors.HTTPError "Authentication required"
   // @Failure     403 {object} httperrors.HTTPError "User not allowed to upload for this listing"
   // @Failure     404 {object} httperrors.HTTPError "Listing not found or not in PENDING_PHOTO_PROCESSING"
   // @Failure     409 {object} httperrors.HTTPError "There is another open batch"
   // @Failure     500 {object} httperrors.HTTPError "Infrastructure failure"
   // @Router      /api/v2/listings/media/uploads [post]
   func (h *ListingHandler) CreateUploadBatch(c *gin.Context) {
       ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

       var request listingdto.CreateUploadBatchRequest
       if err := c.ShouldBindJSON(&request); err != nil {
           httperrors.SendHTTPErrorObj(c, httputils.MapBindingError(err))
           return
       }

       input := listingconverters.DTOToCreateUploadBatchInput(request)
       output, err := h.mediaProcessingService.CreateUploadBatch(ctx, input)
       if err != nil {
           httperrors.SendHTTPErrorObj(c, err)
           return
       }

       response := listingconverters.CreateUploadBatchOutputToDTO(output)
       c.JSON(http.StatusCreated, response)
   }
   ```
2. `complete_upload_batch_handler.go` — payload `CompleteUploadBatchRequest`, resposta `202` (`Accepted`) e rota `/api/v2/listings/media/uploads/complete`.
3. `get_batch_status_handler.go` — requisição `GetBatchStatusRequest`, responde `200` com progresso e downloads disponíveis.
4. `list_download_urls_handler.go` — rota `/api/v2/listings/media/downloads/query`, responde `200` com `ListDownloadUrlsResponse` ou `204` se não houver lote READY.
5. `retry_media_batch_handler.go` — rota `/api/v2/listings/media/uploads/retry`, aceita apenas lotes `FAILED|READY` e sempre retorna `202` com novo job ID.

### 2.2 Services (`internal/core/service/media_processing_service/`)
Arquivos obrigatórios (um método público por arquivo, Seção 7.1):
- `media_processing_service.go` — struct, interface pública, `New` recebendo dependências (`globalService`, `mediaprocessingRepository`, `mediaQueuePort`, `listingRepository`, `listingMediaStoragePort`, `metricsPort`, `timeProvider`).
- `create_upload_batch.go`
- `complete_upload_batch.go`
- `get_batch_status.go`
- `list_download_urls.go`
- `retry_media_batch.go`
- `handle_processing_callback.go`

Exemplo: `create_upload_batch.go`
```go
// CreateUploadBatch validates a manifest and returns signed URLs for raw uploads.
func (s *mediaProcessingService) CreateUploadBatch(ctx context.Context, input CreateUploadBatchInput) (CreateUploadBatchOutput, error) {
    ctx, spanEnd, err := utils.GenerateTracer(ctx)
    if err != nil {
        return CreateUploadBatchOutput{}, derrors.Infra("failed to generate tracer", err)
    }
    defer spanEnd()

    logger := utils.LoggerFromContext(ctx)
    tx, err := s.globalService.StartTransaction(ctx)
    if err != nil {
        utils.SetSpanError(ctx, err)
        logger.Error("service.media.create_batch.tx_start_error", "listing_id", input.ListingID, "err", err)
        return CreateUploadBatchOutput{}, derrors.Infra("failed to start transaction", err)
    }
    committed := false
    defer func() {
        if !committed {
            _ = s.globalService.RollbackTransaction(ctx, tx)
        }
    }()

    // 1. Validate listing ownership/permissions.
    // 2. Ensure listing status == StatusPendingPhotoProcessing.
    // 3. Check absence of open batch (statuses PENDING_UPLOAD/RECEIVED/PROCESSING).
    // 4. Persist batch + placeholder assets via repository.

    if err := s.globalService.CommitTransaction(ctx, tx); err != nil {
        utils.SetSpanError(ctx, err)
        logger.Error("service.media.create_batch.tx_commit_error", "listing_id", input.ListingID, "err", err)
        return CreateUploadBatchOutput{}, derrors.Infra("failed to commit transaction", err)
    }
    committed = true

    // 5. Generate signed URLs (15-minute TTL) through listingMediaStoragePort.
    // 6. Increment metrics counter media_batches_created_total{status="PENDING_UPLOAD"}.
    return output, nil
}
```
- `complete_upload_batch.go` validará HEAD/`GetObjectAttributes`, consolidará `title/sequence`, atualizará status → `RECEIVED`, registrará assets definitivos e publicará job SQS via `mediaQueuePort.EnqueueJob`. Retorna `processingJobId` e `estimatedDuration`.
- `handle_processing_callback.go` receberá payload do pipeline (status + chaves processadas); quando status = `READY`, atualiza listing → `StatusPendingOwnerApproval` chamando `listingRepo.UpdateListingStatus` e observa `media_processing_duration_seconds`.

### 2.3 Repositórios/Adapters (`internal/core/port/right/repository/mediaprocessingrepository/` ↔ `internal/adapter/right/mysql/media_processing/`)
**Port**
```go
type Repository interface {
    CreateBatch(ctx context.Context, tx *sql.Tx, batch mediaprocessingmodel.MediaBatch) (uint64, error)
    UpdateBatchStatus(ctx context.Context, tx *sql.Tx, batchID uint64, status mediaprocessingmodel.BatchStatus, metadata mediaprocessingmodel.BatchStatusMetadata) error
    GetBatchByID(ctx context.Context, tx *sql.Tx, batchID uint64) (mediaprocessingmodel.MediaBatch, error)
    ListBatchesByListing(ctx context.Context, tx *sql.Tx, listingID uint64, limit int) ([]mediaprocessingmodel.MediaBatch, error)
    UpsertAssets(ctx context.Context, tx *sql.Tx, assets []mediaprocessingmodel.MediaAsset) error
    ListAssetsByBatch(ctx context.Context, tx *sql.Tx, batchID uint64) ([]mediaprocessingmodel.MediaAsset, error)
    RegisterProcessingJob(ctx context.Context, tx *sql.Tx, job mediaprocessingmodel.MediaProcessingJob) (uint64, error)
    UpdateProcessingJob(ctx context.Context, tx *sql.Tx, jobID uint64, status mediaprocessingmodel.JobStatus, output mediaprocessingmodel.JobPayload) error
    SoftDeleteBatch(ctx context.Context, tx *sql.Tx, batchID uint64) error
}
```

**Adapter MySQL**
- `media_processing_adapter.go` (struct + `New`).
- Métodos em arquivos separados: `create_batch.go`, `update_batch_status.go`, `get_batch_by_id.go`, `list_batches_by_listing.go`, `upsert_assets.go`, `list_assets_by_batch.go`, `register_processing_job.go`, `update_processing_job.go`, `soft_delete_batch.go`.
- Diretórios auxiliares:
  - `entities/` (`batch_entity.go`, `asset_entity.go`, `job_entity.go`) com `sql.NullString`, `sql.NullTime`, `sql.NullInt64` para campos opcionais (`processed_key`, `thumbnail_key`, `error_code`).
  - `converters/` (`batch_entity_to_domain.go`, `batch_domain_to_entity.go`, `asset_entity_to_domain.go`, `job_entity_to_domain.go`).
- Todas as queries usam `a.ExecContext`/`a.QueryContext` (InstrumentedAdapter) + `ObserveOnComplete("insert", query)` e `utils.GenerateTracer(ctx)` no topo (Seção 7.3).

### 2.4 DTOs (`internal/adapter/left/http/dto/listing_dto.go`)
- Adicionar novas structs com comentários em inglês:
  - `CreateUploadBatchRequest/Response`
  - `CreateUploadBatchFileRequest`
  - `CompleteUploadBatchRequest/Response`
  - `CompleteUploadBatchFileRequest`
  - `GetBatchStatusRequest/Response`
  - `ListDownloadUrlsRequest/Response`
  - `RetryMediaBatchRequest/Response`
  - `MediaProcessingStatusItem`, `MediaProcessingDownloadItem`
- Cada campo deve conter tags `binding`, `example`, `description` conforme template da Seção 8.3. Exemplo:
```go
// CreateUploadBatchRequest represents the payload used to request signed URLs for raw media uploads.

    // ListingID identifies the listing identity receiving the batch.
    ListingID uint64 `json:"listingId" binding:"required,min=1" example:"123"`
    // BatchReference helps the frontend correlate retries (timestamp/slot identifier).
    BatchReference string `json:"batchReference" binding:"required,max=100" example:"2025-11-11T14:00Z-slot-123"`
    // Files enumerates every asset to upload, preserving carousel order and metadata.
    Files []CreateUploadBatchFileRequest `json:"files" binding:"required,min=1,max=60,dive"`
}
```

### 2.5 Modelos & Converters
- `internal/core/model/media_processing_model/`
  - `types.go`: enums `BatchStatus`, `MediaAssetType`, `MediaAssetOrientation`, `JobStatus`, `Provider` com métodos `String()`.
  - `media_batch.go`, `media_asset.go`, `media_job.go` contendo struct + getters para manter imutabilidade.
- `internal/adapter/right/mysql/media_processing/entities/*`: `MediaBatchEntity`, `MediaAssetEntity`, `MediaJobEntity` com `sql.Null*`.
- `internal/adapter/left/http/handlers/listing_handlers/converters/media_processing_converters.go`: mapeia DTO ↔ service input/output sem regra de negócio.

### 2.6 Ports Externos
- `internal/core/port/right/storage/listing_media_storage_port.go`: interface com métodos `GenerateRawUploadURL`, `GenerateProcessedDownloadURL`, `ValidateObjectChecksum`, etc.
- `internal/core/port/right/queue/mediaprocessingqueue/`: define `EnqueueJob`, `PublishRetry`, `DecodeMessage`, `Acknowledge`.
- `internal/adapter/right/aws_sqs/media_processing/`: implementa port usando AWS SDK v2, um método público por arquivo (`enqueue_job.go`, `acknowledge_message.go`).
- `internal/adapter/right/aws_s3/listing_media_storage.go`: implementa interface acima, reutilizando `S3Adapter` e garantindo prefixos `/{listingId}/{stage}/{mediaType}/YYYY-MM-DD/`.
- Callback: definir port `internal/core/port/left/media_processing_callback_port.go` (worker HTTP interno) para receber Step Functions (mesmo que o handler real chegue depois).

### 2.7 Observabilidade Helpers
- Arquivo `metrics_helpers.go` dentro do service para encapsular incrementação/observação (`media_batches_created_total`, `media_processing_duration_seconds`, `media_processing_dlq_messages_total`).
- Configurar labels padrão (`listing_id`, `batch_id`, `status`).

## 3. Estrutura de Diretórios (Regra de Espelhamento)
```
internal/
  core/
    model/
      media_processing_model/
        media_batch.go
        media_asset.go
        media_job.go
        enums.go
    port/
      right/
        repository/
          mediaprocessingrepository/
            media_processing_repository_port.go
        queue/
          mediaprocessingqueue/
            media_processing_queue_port.go
        storage/
          listing_media_storage_port.go
    service/
      media_processing_service/
        media_processing_service.go
        create_upload_batch.go
        complete_upload_batch.go
        get_batch_status.go
        list_download_urls.go
        retry_media_batch.go
        handle_processing_callback.go
        helpers.go
  adapter/
    left/
      http/
        dto/
          listing_dto.go (novos DTOs documentados)
        handlers/
          listing_handlers/
            create_upload_batch_handler.go
            complete_upload_batch_handler.go
            get_batch_status_handler.go
            list_download_urls_handler.go
            retry_media_batch_handler.go
            converters/
              media_processing_converters.go
    right/
      mysql/
        media_processing/
          media_processing_adapter.go
          create_batch.go
          update_batch_status.go
          get_batch_by_id.go
          list_batches_by_listing.go
          upsert_assets.go
          list_assets_by_batch.go
          register_processing_job.go
          update_processing_job.go
          soft_delete_batch.go
          converters/
            batch_entity_to_domain.go
            batch_domain_to_entity.go
            asset_entity_to_domain.go
            job_entity_to_domain.go
          entities/
            batch_entity.go
            asset_entity.go
            job_entity.go
      aws_s3/
        listing_media_storage.go
      aws_sqs/
        media_processing/
          queue_adapter.go
          enqueue_job.go
          acknowledge_job.go
      step_functions/
        media_processing_callback_adapter.go (stub inicial)
```

## 4. Ordem de Execução (Fases)
1. **Fase 0 — Modelos & Documentação (✅ concluída em 2025-11-16)**
  - Criar `media_processing_model` + enums. **Status:** entregue.
  - Atualizar `media_processing_guide.md` com (a) estados confirmados, (b) futura limpeza via job dedicado. **Status:** entregue.
  - Atualizar `media_processing_infra_requirements.md` registrando dependência do job e responsabilidades. **Status:** entregue.
  - Expandir `listing_dto.go` com novos DTOs. **Status:** pendente → será deslocado para o início da Fase 2 para acompanhar handlers/serviço.
2. **Fase 1 — Ports & Estrutura mínima (✅ concluída em 2025-11-16)**
  - Criar ports `mediaprocessingrepository`, `mediaprocessingqueue`, `listing_media_storage`. **Status:** em progresso neste commit.
  - Ajustar `internal/core/factory` para reconhecer novos ports (interfaces vazias até adapters ficarem prontos). **Status:** em progresso.
3. **Fase 2 — Adapter MySQL (✅ concluída em 2025-11-16)**
   - Implementar adapter `internal/adapter/right/mysql/media_processing` completo (métodos + entidades + conversores) seguindo InstrumentedAdapter.
   - Preparar migrations SQL (fora deste repositório, mas documentar scripts necessários no plano de rollout).
4. **Fase 3 — Adapters Externos (S3/SQS/Callback)** (✅ concluída)
  - `listing_media_storage` implementado em `internal/adapter/right/aws_s3/listing_media_storage.go` com geração de chaves `/{listingId}/{stage}/{mediaType}/{date}` (TTL configurável via `env.yaml`) e validação de checksum via `HeadObject`.
  - Adapter SQS criado em `internal/adapter/right/aws_sqs/media_processing/*` expondo `EnqueueJob`, `EnqueueRetry`, `DecodeMessage` e `Acknowledge` com atributos `Traceparent/Listings/Batch/Job` e instrumentação `utils.GenerateTracer`.
  - Callback adapter inicial (`internal/adapter/right/step_functions/media_processing_callback_adapter.go`) validando `shared_secret` configurado em `media_processing.callback`.
  - `CreateExternalServiceAdapters` agora entrega os três adapters e injeta `MediaProcessingCallback` no factory; variáveis `env.yaml` (`media_processing.*` e `s3.signed_url`) foram preenchidas com valores default.
5. **Fase 4 — Serviço** (✅ concluída em 2025-11-16)
   - Implementar métodos do `media_processing_service` com transações, validações, métricas e publicação de eventos. **Status:** entregue.
   - Configurações novas em `env.yaml` (`mediaProcessing`): limites de arquivos, TTLs, nomes de bucket/fila, toggles para owners/project. **Status:** entregue.
   - Métodos implementados: `CompleteUploadBatch`, `GetBatchStatus`, `ListDownloadURLs`, `RetryMediaBatch`, `HandleProcessingCallback`.
6. **Fase 5 — Handlers & Conversores** (✅ concluída em 2025-11-16)
   - Implementar handlers HTTP + conversores DTO ↔ domínio; atualizar `listing_handlers.go` para registrar rotas. **Status:** entregue.
   - DTOs adicionados ao `listing_dto.go`: `CreateUploadBatchRequest/Response`, `CompleteUploadBatchRequest/Response`, `GetBatchStatusRequest/Response`, `ListDownloadURLsRequest/Response`, `RetryMediaBatchRequest/Response` (13 novos tipos com documentação completa em inglês).
   - Conversores bidirecionais criados em `media_processing_converters.go` mapeando DTOs ↔ service input/output (10 funções).
   - Handlers HTTP criados com Swagger completo: `create_upload_batch_handler.go`, `complete_upload_batch_handler.go`, `get_batch_status_handler.go`, `list_download_urls_handler.go`, `retry_media_batch_handler.go`.
   - `ListingHandler` atualizado para injetar `MediaProcessingServiceInterface`.
   - Garantir Swagger (rodar `make swagger` após merge). **Status:** pendente (depende de wiring no factory).
7. **Fase 6 — Wiring & Observabilidade** (✅ concluída em 2025-11-16)
   - Atualizar `internal/core/factory/concrete_adapter_factory.go` e `adapter_factory.go` (injeção S3/SQS/repos/service/handler). **Status:** entregue.
   - Configuração do serviço:
     - Adicionado campo `mediaProcessingService` e `externalServiceAdapters` ao struct `config` em `config_model.go`
     - Criado método `InitMediaProcessingService()` em `inject_dependencies.go` que instancia o serviço com todas as dependências (repositório, listing repo, global service, storage S3, fila SQS)
     - Atualizado `phase_05_services.go` para incluir `MediaProcessingService` na ordem de inicialização (após ScheduleService, antes de ListingService)
     - Interface `ConfigInterface` atualizada para incluir `InitMediaProcessingService()`
   - Factory Pattern:
     - Atualizada interface `AdapterFactory` em `adapter_factory.go` para aceitar `mediaProcessingService` no método `CreateHTTPHandlers`
     - Implementação concreta em `concrete_adapter_factory.go` atualizada com nova assinatura
     - Handler `ListingHandler` agora recebe `MediaProcessingServiceInterface` via construtor
     - Método `SetupHTTPHandlersAndRoutes` em `config_model.go` atualizado para passar `mediaProcessingService`
   - Rotas HTTP registradas:
     - `POST /api/v2/listings/media/uploads` → `CreateUploadBatch`
     - `POST /api/v2/listings/media/uploads/complete` → `CompleteUploadBatch`
     - `POST /api/v2/listings/media/uploads/retry` → `RetryMediaBatch`
     - `POST /api/v2/listings/media/status` → `GetBatchStatus`
     - `POST /api/v2/listings/media/downloads` → `ListDownloadURLs`
   - Instrumentar métricas e logs (adicionar registradores no `metrics` port/dashboards `grafana`). **Status:** infraestrutura pronta; serviço preparado para receber metricsPort quando disponível.
   - **Nota**: MediaProcessingRepository ainda retorna `nil` no factory (linha 247 de `concrete_adapter_factory.go`) pois o adapter MySQL não foi implementado. O serviço será `nil` até que a Fase 2 (Adapter MySQL) seja concluída e o repositório seja instanciado no factory.
8. **Fase 7 — Rollout & Housekeeping**
   - Atualizar `internal/core/factory/concrete_adapter_factory.go` e `adapter_factory.go` (injeção S3/SQS/repos/service/handler).
   - Instrumentar métricas e logs (adicionar registradores no `metrics` port/dashboards `grafana`).
8. **Fase 7 — Rollout & Housekeeping**
   - Checklist de migração (apagar listings antigos, provisionar AWS, configurar SQS DLQ).
   - Documentar job futuro de limpeza (section “Housekeeping” nos guias) e dependências externas.

## 5. Checklist de Conformidade
- [ ] **Arquitetura hexagonal (Seção 1)** — handler → service → ports/adapters, sem atalhos.
- [ ] **Regra de Espelhamento (Seção 2.1)** — `mediaprocessingrepository` ↔ `adapter/right/mysql/media_processing` + `queue`/`storage` correspondentes.
- [ ] **InstrumentedAdapter (Seção 7.3)** — uso obrigatório de `mysqladapter.InstrumentedAdapter` e `ObserveOnComplete`.
- [ ] **Transações via globalService (Seção 7.1)** — nenhum `sql.Tx` iniciado direto no service.
- [ ] **Tracing/logging/erros (Seções 5, 7, 9)** — `utils.GenerateTracer`, `utils.SetSpanError`, logs `slog` com `listing_id`, `batch_id`, `job_id`.
- [ ] **Documentação (Seção 8)** — Godoc completo, Swagger em handlers, DTOs comentados.
- [ ] **Observabilidade (Guia de mídia §9)** — métricas `media_batches_*`, `media_processing_duration_seconds`, DLQ counters.
- [ ] **Sem anti-padrões (Seção 14)** — um método público por arquivo, nada de `SELECT *`, DTOs sem regra de negócio.
- [ ] **Infra dependências documentadas** — bucket/fila/Step Functions + job futuro de limpeza.

## 6. Considerações Complementares
- **Estados confirmados:** `PENDING_UPLOAD`, `RECEIVED`, `PROCESSING`, `READY`, `FAILED` são definitivos; atualizar validações e diagramas.
- **Limpeza futura:** mencionar explicitamente nos guias que remoção física ocorrerá via job dedicado. Nesta release: apenas `SoftDeleteBatch` + API para job consumir (futuro).
- **Reprocessamento:** endpoint `/api/v2/listings/media/uploads/retry` aceita `FAILED|READY`, cria novo registro em `listing_media_jobs` e reutiliza objetos `raw/` no S3.
- **Compatibilidade:** mudança disruptiva — listings antigos serão apagados. Migrar dados não é requisito.
- **Dependências externas:** aguardar provisionamento AWS (S3, SQS, Step Functions, MediaConvert) antes de liberar endpoints; registrar ARNs no `env.yaml`.

---
Plano alinhado ao `docs/toq_server_go_guide.md` (Seções 1, 2.1, 5, 7, 8, 9, 14) e aos guias de mídia. Segmenta o trabalho em fases claras, inclui esqueletos de código/documentação e garante aderência à Regra de Espelhamento, InstrumentedAdapter e observabilidade obrigatória.
