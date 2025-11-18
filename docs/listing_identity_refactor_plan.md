# Listing Identity & Media Workflow Consistency Plan

## 1. Objetivos
- Desambiguar o uso de `listingIdentityId` vs `listingVersionId` em todas as superfícies públicas.
- Uniformizar DTOs e handlers de mídia para expor explicitamente `listingIdentityId`.
- Introduzir tipos fortes opcionais para IDs de identidade e versão, reduzindo trocas acidentais.

## 2. Fluxos de Trabalho
1. **DTOs e Handlers HTTP** – renomear campos, validar IDs obrigatórios e atualizar conversores.
2. **Serviços de Media Processing** – ajustar inputs/outputs, transações e logs para usar `ListingIdentityID`.
3. **Tipos Fortes de IDs** – criar tipos no domínio `listing_model` e propagar gradualmente aos serviços.
4. **Observabilidade** – adicionar métricas de tracing/log para garantir visibilidade dos novos identificadores.

## 3. Ordem Recomendada de Execução
1. Ajustar DTOs (`internal/adapter/left/http/dto/listing_dto.go`).
2. Atualizar handlers HTTP de mídia (`create_upload_batch_handler.go`, `complete_upload_batch_handler.go`, `get_batch_status_handler.go`, `list_download_urls_handler.go`, `retry_media_batch_handler.go`).
3. Revisar conversores (`converters/media_processing_converters.go`).
4. Atualizar serviços (`internal/core/service/media_processing_service/*.go`).
5. Introduzir tipos fortes (`internal/core/model/listing_model/types.go`) e aplicar em inputs principais.
6. Revisar documentação e observabilidade.

## 4. Estrutura Final de Diretórios
- `internal/adapter/left/http/dto/listing_dto.go` (DTOs padronizados).
- `internal/adapter/left/http/handlers/listing_handlers/*.go` (handlers atualizados).
- `internal/adapter/left/http/handlers/listing_handlers/converters/media_processing_converters.go` (mapas DTO⇄serviço).
- `internal/core/service/media_processing_service/` (inputs/outputs e lógica de negócio).
- `internal/core/model/listing_model/types.go` (NOVO) contendo aliases `ListingIdentityID` e `ListingVersionID`.

## 5. Skeletons de Código por Arquivo

### `internal/adapter/left/http/dto/listing_dto.go`
```go
type CreateUploadBatchRequest struct {
    ListingIdentityID uint64 `json:"listingIdentityId" binding:"required,min=1"`
    BatchReference    string `json:"batchReference" binding:"required,max=120"`
    Files             []CreateUploadBatchFileRequest `json:"files" binding:"required,min=1,dive"`
}

type ListDownloadURLsRequest struct {
    ListingIdentityID uint64 `json:"listingIdentityId" binding:"required,min=1"`
    BatchID           uint64 `json:"batchId,omitempty"`
}
```

### `internal/adapter/left/http/handlers/listing_handlers/create_upload_batch_handler.go`
```go
func (lh *ListingHandler) CreateUploadBatch(c *gin.Context) {
    var request dto.CreateUploadBatchRequest
    // bind & validate
    input := converters.DTOToCreateUploadBatchInput(request)
    input.ListingIdentityID = request.ListingIdentityID
    output, err := lh.mediaProcessingService.CreateUploadBatch(ctx, input)
    // respond with dto
}
```

### `complete_upload_batch_handler.go`
```go
func (lh *ListingHandler) CompleteUploadBatch(c *gin.Context) {
    var request dto.CompleteUploadBatchRequest
    input := converters.DTOToCompleteUploadBatchInput(request)
    input.ListingIdentityID = request.ListingIdentityID
    // service call & response
}
```

### `get_batch_status_handler.go`
```go
func (lh *ListingHandler) GetBatchStatus(c *gin.Context) {
    var request dto.GetBatchStatusRequest
    input := converters.DTOToGetBatchStatusInput(request)
    input.ListingIdentityID = request.ListingIdentityID
    // call service
}
```

### `list_download_urls_handler.go`
```go
func (lh *ListingHandler) ListDownloadURLs(c *gin.Context) {
    var request dto.ListDownloadURLsRequest
    input := converters.DTOToListDownloadURLsInput(request)
    input.ListingIdentityID = request.ListingIdentityID
    // service call & response
}
```

### `retry_media_batch_handler.go`
```go
func (lh *ListingHandler) RetryMediaBatch(c *gin.Context) {
    var request dto.RetryMediaBatchRequest
    input := converters.DTOToRetryMediaBatchInput(request)
    input.ListingIdentityID = request.ListingIdentityID
    // call service
}
```

### `internal/adapter/left/http/handlers/listing_handlers/converters/media_processing_converters.go`
```go
func DTOToCreateUploadBatchInput(req dto.CreateUploadBatchRequest) mediaprocessingservice.CreateUploadBatchInput {
    return mediaprocessingservice.CreateUploadBatchInput{
        ListingIdentityID: int64(req.ListingIdentityID),
        BatchReference:    req.BatchReference,
        Files:             convertFiles(req.Files),
    }
}
```

### `internal/core/service/media_processing_service/create_upload_batch.go`
```go
type CreateUploadBatchInput struct {
    ListingIdentityID uint64
    RequestedBy       mediaprocessingmodel.Requester
    Files             []MediaUploadFile
}

func (s *mediaProcessingService) CreateUploadBatch(ctx context.Context, input CreateUploadBatchInput) (CreateUploadBatchOutput, error) {
    if input.ListingIdentityID == 0 {
        return CreateUploadBatchOutput{}, derrors.Validation(...)
    }
    listing, err := s.listingRepo.GetActiveListingVersion(ctx, tx, int64(input.ListingIdentityID))
    // remainder igual
}
```

### `complete_upload_batch.go`
```go
type CompleteUploadBatchInput struct {
    ListingIdentityID uint64
    BatchID           uint64
    Files             []CompletedUploadFile
}

func (s *mediaProcessingService) CompleteUploadBatch(...) (...) {
    // validar ListingIdentityID, carregar batch filtrando por identidade
}
```

### `get_batch_status.go`
```go
type GetBatchStatusInput struct {
    ListingIdentityID uint64
    BatchID           uint64
}

func (s *mediaProcessingService) GetBatchStatus(...) (...) {
    // consulta repo garantindo match de identidade
}
```

### `list_download_urls.go`
```go
type ListDownloadURLsInput struct {
    ListingIdentityID uint64
    BatchID           uint64
}

func (s *mediaProcessingService) ListDownloadURLs(...) (...) {
    // filtra por identidade e batch, retorna URLs
}
```

### `retry_media_batch.go`
```go
type RetryMediaBatchInput struct {
    ListingIdentityID uint64
    BatchID           uint64
    Reason            string
}

func (s *mediaProcessingService) RetryMediaBatch(...) (...) {
    // valida identidade e status terminal
}
```

### `internal/core/model/listing_model/types.go` (novo)
```go
package listingmodel

type ListingIdentityID int64
type ListingVersionID int64
```

## 6. Sugestões Adicionais
- Expor métricas (`listing.media.batch.identity_mismatch`) para capturar tentativas com IDs inválidos.
- Atualizar `docs/toq_server_go_guide.md` adicionando uma tabela de referência de IDs.
- Considerar middleware para extrair `listingIdentityId` de path/query automaticamente, evitando duplicação em handlers.
