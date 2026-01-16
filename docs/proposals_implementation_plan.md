# Plano de Implementação – Sistema de Propostas (TOQ Server)

## 1. Diagnóstico
- **Arquivos analisados:** [docs/toq_server_go_guide.md](docs/toq_server_go_guide.md), [internal/adapter/right/mysql/proposal/create_proposal.go](internal/adapter/right/mysql/proposal/create_proposal.go#L31-L105), [internal/adapter/right/mysql/proposal/list_proposals.go](internal/adapter/right/mysql/proposal/list_proposals.go#L33-L260), [internal/core/model/proposal_model/proposal.go](internal/core/model/proposal_model/proposal.go#L8-L175), [internal/core/port/right/repository/proposal_repository/proposal_repository_interface.go](internal/core/port/right/repository/proposal_repository/proposal_repository_interface.go), [internal/core/service/global_service/notification_service.go](internal/core/service/global_service/notification_service.go#L1-L190), [internal/core/service/global_service/transactions.go](internal/core/service/global_service/transactions.go#L1-L78), diretórios [internal/core/service](internal/core/service) e [internal/adapter/left/http/handlers](internal/adapter/left/http/handlers), além do schema em [scripts/db_creation.sql](scripts/db_creation.sql).
- **Justificativa da abordagem:**
  - O adapter atual grava campos financeiros, clientes e favoritos inexistentes na nova regra ([internal/adapter/right/mysql/proposal/create_proposal.go#L31-L89](internal/adapter/right/mysql/proposal/create_proposal.go#L31-L89)) e ignora o `proposal_text`. Esses dados não serão utilizados e conflitam com o pedido de armazenar texto/PDF.
  - As consultas ainda retornam estatísticas monetárias e filtros por valores ([internal/adapter/right/mysql/proposal/list_proposals.go#L33-L195](internal/adapter/right/mysql/proposal/list_proposals.go#L33-L195)), o que mantém complexidade desnecessária e impede a exposição simples do histórico solicitado.
  - O domínio carece de campos para texto/documentos e mantém enums (`TransactionType`, `PaymentMethod`, `GuaranteeType`) que serão removidos ([internal/core/model/proposal_model/proposal.go#L8-L175](internal/core/model/proposal_model/proposal.go#L8-L175)). Também não há forma de representar o binário do PDF.
  - Não existe service/handler dedicado (ausência em [internal/core/service](internal/core/service) e [internal/adapter/left/http/handlers](internal/adapter/left/http/handlers)), impossibilitando controles de permissão, transação (Seção 7.1) e geração de Swagger (Seção 8.2).
  - O schema MySQL não contém `proposals`, `proposal_documents` e colunas nas `listing_identities` ([scripts/db_creation.sql](scripts/db_creation.sql)), portanto flags `has_pending_proposal`/`has_accepted_proposal` ainda não podem ser atualizadas.
  - O serviço global já expõe transações e UnifiedNotificationService ([internal/core/service/global_service/transactions.go#L1-L73](internal/core/service/global_service/transactions.go#L1-L73), [internal/core/service/global_service/notification_service.go#L1-L190](internal/core/service/global_service/notification_service.go#L1-L190)); reutilizar esse stack é o caminho aderente ao guia (seções 2 e 7) para push notifications e auditoria.
  - Sobre `created_at`/`updated_at`: o modelo de dados padronizado não inclui esses campos na tabela `proposals`, portanto a API/adapter não deve depender deles. O controle de status permanece em `accepted_at/rejected_at/cancelled_at` e nas flags de listing, garantindo rastreabilidade via status e auditoria.
- **Impacto esperado:**
  - Quebra intencional de compatibilidade na camada de repositório/serviço (novas assinaturas com `*sql.Tx`, filtros simplificados e remoção de favoritos/estatísticas).
  - Novos endpoints HTTP `/api/v2/proposals/*`, DTOs e documentação Swagger regenerada.
  - Ajustes de schema em `proposals`, `proposal_documents` e `listing_identities`, além de novos índices por status/atores.
  - Integração com UnifiedNotificationService para todos os eventos (envio/cancelamento/aceite/recusa) e auditoria via `auditService.RecordChange` (actor/target/correlation).
- **Melhorias possíveis:**
  - Simplificação de payloads (apenas texto + PDF opcional) reduzindo superfície de dados sensíveis.
  - Atualização única das flags de listing dentro da mesma transação para manter consistência eventual.
  - Remoção dos métodos `SetFavorite`, `GetStats` e filtros monetários, reduzindo o adapter a um subconjunto mínimo e aderente à Regra de Espelhamento.

## 2. Code Skeletons
### 2.1 DTOs
- **Arquivo:** [internal/adapter/left/http/dto/proposal_dto.go](internal/adapter/left/http/dto/proposal_dto.go)
```go
package dto

import "time"

// CreateProposalRequest represents the realtor payload to submit a new proposal.
type CreateProposalRequest struct {
  ListingIdentityID int64                    `json:"listingIdentityId" binding:"required,min=1" example:"981"`
  ProposalText      string                  `json:"proposalText" binding:"required,min=1,max=5000" example:"Gostaria de propor pagamento em 30 dias"`
  Document          *ProposalDocumentUpload `json:"document,omitempty"`
}

// UpdateProposalRequest allows editing a pending proposal text/document.
type UpdateProposalRequest struct {
  ProposalID   int64                    `json:"proposalId" binding:"required,min=1" example:"120"`
  ProposalText string                  `json:"proposalText" binding:"required,min=1,max=5000"`
  Document     *ProposalDocumentUpload `json:"document,omitempty"`
}

// ProposalDocumentUpload carries the base64 PDF metadata limited to 1MB.
type ProposalDocumentUpload struct {
  FileName      string `json:"fileName" binding:"required,min=1,max=120" example:"proposta.pdf"`
  MimeType      string `json:"mimeType" binding:"required,oneof=application/pdf" example:"application/pdf"`
  Base64Payload string `json:"base64Payload" binding:"required"`
}

// CancelProposalRequest is used by realtors before owner acceptance.
type CancelProposalRequest struct {
  ProposalID int64 `json:"proposalId" binding:"required,min=1"`
}

// AcceptProposalRequest is triggered by the owner.
type AcceptProposalRequest struct {
  ProposalID int64 `json:"proposalId" binding:"required,min=1"`
}

// RejectProposalRequest stores the owner reason for refusal.
type RejectProposalRequest struct {
  ProposalID int64  `json:"proposalId" binding:"required,min=1"`
  Reason     string `json:"reason" binding:"required,min=1,max=500"`
}

// ListProposalsQuery is shared by realtor/owner GET endpoints.
type ListProposalsQuery struct {
  Statuses          []string `form:"status" binding:"omitempty,dive,oneof=pending accepted refused cancelled"`
  ListingIdentityID int64    `form:"listingIdentityId" binding:"omitempty,min=1"`
  Page              int      `form:"page" binding:"omitempty,min=1" default:"1"`
  PageSize          int      `form:"pageSize" binding:"omitempty,min=1,max=100" default:"20"`
}

// GetProposalDetailRequest returns the full payload including documents.
type GetProposalDetailRequest struct {
  ProposalID int64 `json:"proposalId" binding:"required,min=1"`
}

// ProposalResponse summarizes proposal information for list views.
type ProposalResponse struct {
  ID                int64     `json:"id"`
  ListingIdentityID int64     `json:"listingIdentityId"`
  Status            string    `json:"status"`
  ProposalText      string    `json:"proposalText"`
  AcceptedAt        *time.Time `json:"acceptedAt,omitempty"`
  RejectedAt        *time.Time `json:"rejectedAt,omitempty"`
  CancelledAt       *time.Time `json:"cancelledAt,omitempty"`
  DocumentsCount    int       `json:"documentsCount"`
}

// ProposalDocumentResponse exposes metadata and optional base64 payload.
type ProposalDocumentResponse struct {
  ID            int64  `json:"id"`
  FileName      string `json:"fileName"`
  MimeType      string `json:"mimeType"`
  FileSizeBytes int64  `json:"fileSizeBytes"`
  Base64Payload string `json:"base64Payload,omitempty"`
}

// ProposalDetailResponse aggregates summary + documents and owner metadata.
type ProposalDetailResponse struct {
  Proposal  ProposalResponse            `json:"proposal"`
  Documents []ProposalDocumentResponse `json:"documents"`
}

// ListProposalsResponse is returned by both realtor/owner endpoints.
type ListProposalsResponse struct {
  Items []ProposalResponse `json:"items"`
  Total int64              `json:"total"`
}
```

### 2.2 Ports e Handlers
- **Arquivo:** [internal/core/port/left/http/proposalhandler/proposal_handler_port.go](internal/core/port/left/http/proposalhandler/proposal_handler_port.go)
```go
package proposalhandler

import "github.com/gin-gonic/gin"

// Handler defines all HTTP entrypoints for the proposal domain.
type Handler interface {
  CreateProposal(c *gin.Context)
  UpdateProposal(c *gin.Context)
  CancelProposal(c *gin.Context)
  AcceptProposal(c *gin.Context)
  RejectProposal(c *gin.Context)
  ListRealtorProposals(c *gin.Context)
  ListOwnerProposals(c *gin.Context)
  GetProposalDetail(c *gin.Context)
}
```

- **Arquivo:** [internal/adapter/left/http/handlers/proposal_handlers/proposal_handler.go](internal/adapter/left/http/handlers/proposal_handlers/proposal_handler.go)
```go
package proposalhandlers

import (
  proposalservice "github.com/projeto-toq/toq_server/internal/core/service/proposal_service"
)

// ProposalHandler wires DTO conversion, authentication context and the service port.
type ProposalHandler struct {
  proposalService proposalservice.Service
}

// NewProposalHandler builds a handler with its dependencies injected by the factory.
func NewProposalHandler(service proposalservice.Service) *ProposalHandler {
  return &ProposalHandler{proposalService: service}
}
```

- **Arquivo:** [internal/adapter/left/http/handlers/proposal_handlers/create_proposal_handler.go](internal/adapter/left/http/handlers/proposal_handlers/create_proposal_handler.go)
```go
package proposalhandlers

import (
  net/http

  "github.com/gin-gonic/gin"

  dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
  httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
  proposalservice "github.com/projeto-toq/toq_server/internal/core/service/proposal_service"
  coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateProposal handles realtor submissions of new proposals.
//
// @Summary     Submit a proposal for a published listing
// @Description Allows authenticated realtors to send a free-text proposal and optional PDF attachment (≤1MB) for a listing identity they do not own.
// @Tags        Proposals
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       Authorization header string true "Bearer <token>"
// @Param       request body dto.CreateProposalRequest true "Proposal payload with optional PDF encoded in base64"
// @Success     201 {object} dto.ProposalResponse
// @Failure     400 {object} dto.ErrorResponse "Invalid payload or PDF too large"
// @Failure     401 {object} dto.ErrorResponse "Authentication required"
// @Failure     409 {object} dto.ErrorResponse "Listing already has pending/accepted proposal"
// @Failure     422 {object} dto.ErrorResponse "Business validation failure"
// @Failure     500 {object} dto.ErrorResponse "Infrastructure failure"
// @Router      /proposals [post]
func (h *ProposalHandler) CreateProposal(c *gin.Context) {
  ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
  var request dto.CreateProposalRequest
  if err := c.ShouldBindJSON(&request); err != nil {
    httperrors.SendHTTPErrorObj(c, httperrors.MapBindingError(err))
    return
  }
  actor := proposalservice.ExtractActorFromContext(c)
  input, err := proposalservice.ConvertCreateRequest(request, actor)
  if err != nil {
    httperrors.SendHTTPErrorObj(c, err)
    return
  }
  result, err := h.proposalService.CreateProposal(ctx, input)
  if err != nil {
    httperrors.SendHTTPErrorObj(c, err)
    return
  }
  response := dto.FromProposalDomain(result)
  c.JSON(http.StatusCreated, response)
}
```

- **Arquivo:** [internal/adapter/left/http/handlers/proposal_handlers/update_proposal_handler.go](internal/adapter/left/http/handlers/proposal_handlers/update_proposal_handler.go)
```go
package proposalhandlers

import (
  net/http

  "github.com/gin-gonic/gin"

  dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
  httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
  proposalservice "github.com/projeto-toq/toq_server/internal/core/service/proposal_service"
  coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateProposal allows the author to edit a pending proposal.
//
// @Summary     Edit a pending proposal
// @Description Updates text and/or PDF while the proposal is still pending, ensuring idempotency and actor ownership.
// @Tags        Proposals
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       Authorization header string true "Bearer <token>"
// @Param       request body dto.UpdateProposalRequest true "Editable fields"
// @Success     200 {object} dto.ProposalResponse
// @Failure     400,401,403,409,422,500 {object} dto.ErrorResponse
// @Router      /proposals [put]
func (h *ProposalHandler) UpdateProposal(c *gin.Context) {
  ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
  var request dto.UpdateProposalRequest
  if err := c.ShouldBindJSON(&request); err != nil {
    httperrors.SendHTTPErrorObj(c, httperrors.MapBindingError(err))
    return
  }
  actor := proposalservice.ExtractActorFromContext(c)
  input, err := proposalservice.ConvertUpdateRequest(request, actor)
  if err != nil {
    httperrors.SendHTTPErrorObj(c, err)
    return
  }
  result, err := h.proposalService.UpdateProposal(ctx, input)
  if err != nil {
    httperrors.SendHTTPErrorObj(c, err)
    return
  }
  c.JSON(http.StatusOK, dto.FromProposalDomain(result))
}
```

- **Arquivo:** [internal/adapter/left/http/handlers/proposal_handlers/cancel_proposal_handler.go](internal/adapter/left/http/handlers/proposal_handlers/cancel_proposal_handler.go)
```go
package proposalhandlers

import (
  net/http

  "github.com/gin-gonic/gin"

  dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
  httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
  proposalservice "github.com/projeto-toq/toq_server/internal/core/service/proposal_service"
  coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// CancelProposal lets the realtor withdraw a pending proposal.
//
// @Summary     Cancel a proposal before owner decision
// @Description Realtor-owned proposals can be cancelled anytime before acceptance, generating notifications to the owner.
// @Tags        Proposals
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       request body dto.CancelProposalRequest true "Proposal identifier"
// @Success     204 {string} string "No Content"
// @Failure     400,401,403,409,500 {object} dto.ErrorResponse
// @Router      /proposals/cancel [post]
func (h *ProposalHandler) CancelProposal(c *gin.Context) {
  ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
  var request dto.CancelProposalRequest
  if err := c.ShouldBindJSON(&request); err != nil {
    httperrors.SendHTTPErrorObj(c, httperrors.MapBindingError(err))
    return
  }
  actor := proposalservice.ExtractActorFromContext(c)
  if err := h.proposalService.CancelProposal(ctx, proposalservice.StatusChangeInput{ProposalID: request.ProposalID, Actor: actor}); err != nil {
    httperrors.SendHTTPErrorObj(c, err)
    return
  }
  c.Status(http.StatusNoContent)
}
```

- **Arquivo:** [internal/adapter/left/http/handlers/proposal_handlers/accept_proposal_handler.go](internal/adapter/left/http/handlers/proposal_handlers/accept_proposal_handler.go)
```go
package proposalhandlers

import (
  net/http

  "github.com/gin-gonic/gin"

  dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
  httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
  proposalservice "github.com/projeto-toq/toq_server/internal/core/service/proposal_service"
  coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// AcceptProposal confirms the owner's approval and updates listing flags.
// @Summary     Accept a proposal
// @Description Owners can accept a single pending proposal per listing, triggering notifications and flag updates.
// @Tags        Proposals
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       request body dto.AcceptProposalRequest true "Proposal identifier"
// @Success     200 {object} dto.ProposalResponse
// @Failure     400,401,403,409,500 {object} dto.ErrorResponse
// @Router      /proposals/accept [post]
func (h *ProposalHandler) AcceptProposal(c *gin.Context) {
  ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
  var request dto.AcceptProposalRequest
  if err := c.ShouldBindJSON(&request); err != nil {
    httperrors.SendHTTPErrorObj(c, httperrors.MapBindingError(err))
    return
  }
  actor := proposalservice.ExtractActorFromContext(c)
  result, err := h.proposalService.AcceptProposal(ctx, proposalservice.StatusChangeInput{ProposalID: request.ProposalID, Actor: actor})
  if err != nil {
    httperrors.SendHTTPErrorObj(c, err)
    return
  }
  c.JSON(http.StatusOK, dto.FromProposalDomain(result))
}
```

- **Arquivo:** [internal/adapter/left/http/handlers/proposal_handlers/reject_proposal_handler.go](internal/adapter/left/http/handlers/proposal_handlers/reject_proposal_handler.go)
```go
package proposalhandlers

import (
  net/http

  "github.com/gin-gonic/gin"

  dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
  httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
  proposalservice "github.com/projeto-toq/toq_server/internal/core/service/proposal_service"
  coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// RejectProposal stores the owner's reason and informs the realtor.
// @Summary     Reject a proposal with reason
// @Description Owners must provide a free text reason; realtor receives a push notification with the status update.
// @Tags        Proposals
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       request body dto.RejectProposalRequest true "Proposal identifier and reason"
// @Success     200 {object} dto.ProposalResponse
// @Failure     400,401,403,409,422,500 {object} dto.ErrorResponse
// @Router      /proposals/reject [post]
func (h *ProposalHandler) RejectProposal(c *gin.Context) {
  ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
  var request dto.RejectProposalRequest
  if err := c.ShouldBindJSON(&request); err != nil {
    httperrors.SendHTTPErrorObj(c, httperrors.MapBindingError(err))
    return
  }
  actor := proposalservice.ExtractActorFromContext(c)
  result, err := h.proposalService.RejectProposal(ctx, proposalservice.StatusChangeInput{ProposalID: request.ProposalID, Actor: actor, Reason: request.Reason})
  if err != nil {
    httperrors.SendHTTPErrorObj(c, err)
    return
  }
  c.JSON(http.StatusOK, dto.FromProposalDomain(result))
}
```

- **Arquivo:** [internal/adapter/left/http/handlers/proposal_handlers/list_realtor_proposals_handler.go](internal/adapter/left/http/handlers/proposal_handlers/list_realtor_proposals_handler.go)
```go
package proposalhandlers

import (
  net/http

  "github.com/gin-gonic/gin"

  dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
  httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
  proposalservice "github.com/projeto-toq/toq_server/internal/core/service/proposal_service"
  coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListRealtorProposals returns paginated history filtered by realtor context.
// @Summary     List realtor proposals
// @Description Returns paginated proposals created by the authenticated realtor, supporting filters by status, listing identity and creation dates.
// @Tags        Proposals
// @Security    BearerAuth
// @Produce     json
// @Param       status query []string false "Status filter" collectionFormat(multi)
// @Param       listingIdentityId query int false "Listing identity"
// @Param       page query int false "Page number" default(1)
// @Param       pageSize query int false "Page size" default(20)
// @Success     200 {object} dto.ListProposalsResponse
// @Failure     400,401,500 {object} dto.ErrorResponse
// @Router      /proposals/realtor [get]
func (h *ProposalHandler) ListRealtorProposals(c *gin.Context) {
  ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
  var query dto.ListProposalsQuery
  if err := c.ShouldBindQuery(&query); err != nil {
    httperrors.SendHTTPErrorObj(c, httperrors.MapBindingError(err))
    return
  }
  actor := proposalservice.ExtractActorFromContext(c)
  result, err := h.proposalService.ListRealtorProposals(ctx, proposalservice.ListFilterFromQuery(query, actor))
  if err != nil {
    httperrors.SendHTTPErrorObj(c, err)
    return
  }
  c.JSON(http.StatusOK, dto.FromProposalList(result))
}
```

- **Arquivo:** [internal/adapter/left/http/handlers/proposal_handlers/list_owner_proposals_handler.go](internal/adapter/left/http/handlers/proposal_handlers/list_owner_proposals_handler.go)
```go
package proposalhandlers

import (
  net/http

  "github.com/gin-gonic/gin"

  dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
  httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
  proposalservice "github.com/projeto-toq/toq_server/internal/core/service/proposal_service"
  coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListOwnerProposals lists proposals received by the owner.
// @Summary     List owner proposals
// @Description Owners can see all proposals received across listings with the same filters as realtors, but scoped to identities they own.
// @Tags        Proposals
// @Security    BearerAuth
// @Produce     json
// @Param       status query []string false "Status filter"
// @Param       listingIdentityId query int false "Listing identity"
// @Param       page query int false "Page number"
// @Param       pageSize query int false "Page size"
// @Success     200 {object} dto.ListProposalsResponse
// @Failure     400,401,500 {object} dto.ErrorResponse
// @Router      /proposals/owner [get]
func (h *ProposalHandler) ListOwnerProposals(c *gin.Context) {
  ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
  var query dto.ListProposalsQuery
  if err := c.ShouldBindQuery(&query); err != nil {
    httperrors.SendHTTPErrorObj(c, httperrors.MapBindingError(err))
    return
  }
  actor := proposalservice.ExtractActorFromContext(c)
  result, err := h.proposalService.ListOwnerProposals(ctx, proposalservice.ListFilterFromQuery(query, actor))
  if err != nil {
    httperrors.SendHTTPErrorObj(c, err)
    return
  }
  c.JSON(http.StatusOK, dto.FromProposalList(result))
}
```

- **Arquivo:** [internal/adapter/left/http/handlers/proposal_handlers/get_proposal_detail_handler.go](internal/adapter/left/http/handlers/proposal_handlers/get_proposal_detail_handler.go)
```go
package proposalhandlers

import (
  net/http

  "github.com/gin-gonic/gin"

  dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
  httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
  proposalservice "github.com/projeto-toq/toq_server/internal/core/service/proposal_service"
  coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetProposalDetail returns proposal metadata plus documents (base64 encoded) to the owner or the realtor.
// @Summary     Retrieve proposal detail
// @Description Provides proposal metadata and the binary (base64) of any PDF attachments for authorized actors.
// @Tags        Proposals
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       request body dto.GetProposalDetailRequest true "Proposal identifier"
// @Success     200 {object} dto.ProposalDetailResponse
// @Failure     400,401,403,404,500 {object} dto.ErrorResponse
// @Router      /proposals/detail [post]
func (h *ProposalHandler) GetProposalDetail(c *gin.Context) {
  ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
  var request dto.GetProposalDetailRequest
  if err := c.ShouldBindJSON(&request); err != nil {
    httperrors.SendHTTPErrorObj(c, httperrors.MapBindingError(err))
    return
  }
  actor := proposalservice.ExtractActorFromContext(c)
  detail, err := h.proposalService.GetProposalDetail(ctx, proposalservice.DetailInput{ProposalID: request.ProposalID, Actor: actor})
  if err != nil {
    httperrors.SendHTTPErrorObj(c, err)
    return
  }
  c.JSON(http.StatusOK, dto.FromProposalDetail(detail))
}
```

### 2.3 Services
- **Arquivo:** [internal/core/service/proposal_service/proposal_service.go](internal/core/service/proposal_service/proposal_service.go)
```go
package proposalservice

import (
  proposalrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/proposal_repository"
  listingrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/listing_repository"
  globalservice "github.com/projeto-toq/toq_server/internal/core/service/global_service"
)

// Service exposes the orchestration required by handlers.
type Service interface {
  CreateProposal(ctx context.Context, input CreateProposalInput) (proposalmodel.ProposalInterface, error)
  UpdateProposal(ctx context.Context, input UpdateProposalInput) (proposalmodel.ProposalInterface, error)
  CancelProposal(ctx context.Context, input StatusChangeInput) error
  AcceptProposal(ctx context.Context, input StatusChangeInput) (proposalmodel.ProposalInterface, error)
  RejectProposal(ctx context.Context, input StatusChangeInput) (proposalmodel.ProposalInterface, error)
  ListRealtorProposals(ctx context.Context, filter ListFilter) (ListResult, error)
  ListOwnerProposals(ctx context.Context, filter ListFilter) (ListResult, error)
  GetProposalDetail(ctx context.Context, input DetailInput) (DetailResult, error)
}

type proposalService struct {
  proposalRepo proposalrepository.Repository
  listingRepo  listingrepository.ListingRepoPortInterface
  globalSvc    globalservice.GlobalServiceInterface
  notifier     globalservice.UnifiedNotificationService
  maxDocBytes  int64
}

// New builds a Service respecting the factory order (Seção 4 do guia).
func New(
  proposalRepo proposalrepository.Repository,
  listingRepo listingrepository.ListingRepoPortInterface,
  globalSvc globalservice.GlobalServiceInterface,
) Service {
  return &proposalService{
    proposalRepo: proposalRepo,
    listingRepo:  listingRepo,
    globalSvc:    globalSvc,
    notifier:     globalSvc.GetUnifiedNotificationService(),
    maxDocBytes:  1_000_000,
  }
}
```

- **Arquivo:** [internal/core/service/proposal_service/types.go](internal/core/service/proposal_service/types.go)
```go
package proposalservice

import (
  "context"

  proposalmodel "github.com/projeto-toq/toq_server/internal/core/model/proposal_model"
)

// Actor stores the authenticated user metadata extracted from middlewares.
type Actor struct {
  UserID  int64
  RoleSlug string
}

// CreateProposalInput aggregates validated data from handlers.
type CreateProposalInput struct {
  ListingIdentityID int64
  RealtorID         int64
  ProposalText      string
  Document          *DocumentPayload
}

// UpdateProposalInput extends CreateProposalInput with an ID.
type UpdateProposalInput struct {
  ProposalID  int64
  EditorID    int64
  ProposalText string
  Document     *DocumentPayload
}

// DocumentPayload carries decoded bytes and metadata.
type DocumentPayload struct {
  FileName      string
  MimeType      string
  Bytes         []byte
  SizeBytes     int64
}

// StatusChangeInput is reused by cancel/accept/reject flows.
type StatusChangeInput struct {
  ProposalID int64
  Actor      Actor
  Reason     string
}

// ListFilter stores normalized filters for repository queries.
type ListFilter struct {
  Actor      Actor
  Statuses   []proposalmodel.Status
  ListingID  *int64
  CreatedGTE *time.Time
  CreatedLTE *time.Time
  Page       int
  PageSize   int
}

// ListResult is returned to handlers before DTO serialization.
type ListResult struct {
  Items []proposalmodel.ProposalInterface
  Total int64
}

// DetailInput ensures only owners or authors can inspect documents.
type DetailInput struct {
  ProposalID int64
  Actor      Actor
}

// DetailResult stores proposal and documents.
type DetailResult struct {
  Proposal  proposalmodel.ProposalInterface
  Documents []proposalmodel.ProposalDocumentInterface
}

// ExtractActorFromContext reads user metadata already set by DeviceContextMiddleware.
func ExtractActorFromContext(c *gin.Context) Actor {
  // TODO: replace with the actual helper used across handlers (e.g. httpctx.GetAuthenticatedUser).
  return Actor{}
}

// ConvertCreateRequest validates DTO + actor and builds input for services.
func ConvertCreateRequest(request dto.CreateProposalRequest, actor Actor) (CreateProposalInput, error) {
  // TODO: decode base64, enforce max size and ensure actor role is realtor.
  return CreateProposalInput{}, nil
}
```

- **Arquivo:** [internal/core/service/proposal_service/create_proposal.go](internal/core/service/proposal_service/create_proposal.go)
```go
package proposalservice

import (
  "context"
  "database/sql"
  "fmt"
  "log/slog"

  derrors "github.com/projeto-toq/toq_server/internal/core/derrors"
  globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
  proposalmodel "github.com/projeto-toq/toq_server/internal/core/model/proposal_model"
  "github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateProposal orchestrates validation, persistence, flag updates and notifications.
func (s *proposalService) CreateProposal(ctx context.Context, input CreateProposalInput) (proposalmodel.ProposalInterface, error) {
  ctx, spanEnd, err := utils.GenerateTracer(ctx)
  if err != nil {
    return nil, derrors.Infra("failed to start tracer", err)
  }
  defer spanEnd()
  ctx = utils.ContextWithLogger(ctx)
  logger := utils.LoggerFromContext(ctx)

  tx, err := s.globalSvc.StartTransaction(ctx)
  if err != nil {
    logger.Error("proposal.create.tx_start_error", "listing_identity_id", input.ListingIdentityID, "err", err)
    return nil, err
  }
  defer s.rollbackOnError(ctx, tx, &err)

  identity, err := s.listingRepo.GetListingIdentityByID(ctx, tx, input.ListingIdentityID)
  if err != nil {
    return nil, s.mapListingError(err)
  }
  if identity.UserID == input.RealtorID {
    return nil, derrors.Forbidden("owners cannot send proposals to themselves", nil)
  }
  proposal := proposalmodel.NewProposal()
  proposal.SetListingIdentityID(input.ListingIdentityID)
  proposal.SetRealtorID(input.RealtorID)
  proposal.SetOwnerID(identity.UserID)
  proposal.SetProposalText(input.ProposalText)
  proposal.SetStatus(proposalmodel.StatusPending)

  if err = s.proposalRepo.CreateProposal(ctx, tx, proposal); err != nil {
    return nil, derrors.Infra("failed to persist proposal", err)
  }
  if input.Document != nil {
    if err = s.createProposalDocument(ctx, tx, proposal.ID(), input.Document); err != nil {
      return nil, err
    }
  }
  if err = s.listingRepo.UpdateProposalFlags(ctx, tx, listingrepository.ProposalFlagsUpdate{
    ListingIdentityID: proposal.ListingIdentityID(),
    HasPending:        true,
    AcceptedProposalID: sql.NullInt64{},
  }); err != nil {
    return nil, derrors.Infra("failed to set listing flags", err)
  }
  auditRecord := auditservice.BuildRecordFromContext(
    ctx,
    proposal.RealtorID(),
    auditmodel.AuditTarget{Type: auditmodel.TargetProposal, ID: proposal.ID()},
    auditmodel.OperationCreate,
    map[string]any{"listingIdentityId": proposal.ListingIdentityID()},
  )
  if err = s.auditService.RecordChange(ctx, tx, auditRecord); err != nil {
    return nil, err
  }
  if err = s.globalSvc.CommitTransaction(ctx, tx); err != nil {
    return nil, err
  }
  go s.notifyOwnerNewProposal(context.Background(), proposal)
  return proposal, nil
}
```

- **Arquivos de métodos adicionais:** seguir o mesmo padrão com tracer, transação e comentários explicando cada etapa conforme Seção 7.1 do guia:
  - [internal/core/service/proposal_service/update_proposal.go](internal/core/service/proposal_service/update_proposal.go)
  - [internal/core/service/proposal_service/cancel_proposal.go](internal/core/service/proposal_service/cancel_proposal.go)
  - [internal/core/service/proposal_service/accept_proposal.go](internal/core/service/proposal_service/accept_proposal.go)
  - [internal/core/service/proposal_service/reject_proposal.go](internal/core/service/proposal_service/reject_proposal.go)
  - [internal/core/service/proposal_service/list_realtor_proposals.go](internal/core/service/proposal_service/list_realtor_proposals.go)
  - [internal/core/service/proposal_service/list_owner_proposals.go](internal/core/service/proposal_service/list_owner_proposals.go)
  - [internal/core/service/proposal_service/get_proposal_detail.go](internal/core/service/proposal_service/get_proposal_detail.go)
  - [internal/core/service/proposal_service/helpers.go](internal/core/service/proposal_service/helpers.go) para `rollbackOnError`, validação de PDF, conversões e envio de notificações.

Cada arquivo deve conter Godoc descrevendo fluxo, regras de negócio, chamadas auxiliares e o uso de `utils.SetSpanError`/`slog` em falhas de infraestrutura.

### 2.4 Repositórios MySQL
- **Arquivo:** [internal/core/port/right/repository/proposal_repository/proposal_repository_interface.go](internal/core/port/right/repository/proposal_repository/proposal_repository_interface.go)
```go
package proposalrepository

import (
  "context"
  "database/sql"

  proposalmodel "github.com/projeto-toq/toq_server/internal/core/model/proposal_model"
)

// Repository exposes CRUD operations for proposals and documents (Seção 7.3 do guia).
type Repository interface {
  CreateProposal(ctx context.Context, tx *sql.Tx, proposal proposalmodel.ProposalInterface) error
  UpdateProposalText(ctx context.Context, tx *sql.Tx, proposal proposalmodel.ProposalInterface) error
  UpdateProposalStatus(ctx context.Context, tx *sql.Tx, proposal proposalmodel.ProposalInterface, expected proposalmodel.Status) error
  GetProposalByID(ctx context.Context, tx *sql.Tx, proposalID int64) (proposalmodel.ProposalInterface, error)
  GetProposalByIDForUpdate(ctx context.Context, tx *sql.Tx, proposalID int64) (proposalmodel.ProposalInterface, error)
  ListProposals(ctx context.Context, tx *sql.Tx, filter proposalmodel.ListFilter) (proposalmodel.ListResult, error)
  CreateDocument(ctx context.Context, tx *sql.Tx, document proposalmodel.ProposalDocumentInterface) error
  ListDocuments(ctx context.Context, tx *sql.Tx, proposalID int64, includeBlob bool) ([]proposalmodel.ProposalDocumentInterface, error)
}
```

- **Arquivo:** [internal/adapter/right/mysql/proposal/create_proposal.go](internal/adapter/right/mysql/proposal/create_proposal.go)
```go
package mysqlproposaladapter

import (
  "context"
  "fmt"
  "time"

  proposalmodel "github.com/projeto-toq/toq_server/internal/core/model/proposal_model"
  "github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateProposal inserts a new proposal row using the provided transaction.
func (a *ProposalAdapter) CreateProposal(ctx context.Context, tx *sql.Tx, proposal proposalmodel.ProposalInterface) error {
  ctx, spanEnd, err := utils.GenerateTracer(ctx)
  if err != nil {
    return err
  }
  defer spanEnd()
  ctx = utils.ContextWithLogger(ctx)
  logger := utils.LoggerFromContext(ctx)

  entity := converters.ToProposalEntity(proposal)

  query := `INSERT INTO proposals (
    listing_identity_id,
    realtor_id,
    owner_id,
    status,
    proposal_text,
    rejection_reason,
    accepted_at,
    rejected_at,
    cancelled_at,
    deleted
  ) VALUES (?,?,?,?,?,?,?,?,?,0)`

  result, execErr := a.ExecContext(ctx, tx, "insert_proposal", query,
    entity.ListingIdentityID,
    entity.RealtorID,
    entity.OwnerID,
    entity.Status,
    entity.ProposalText,
    entity.RejectionReason,
    entity.AcceptedAt,
    entity.RejectedAt,
    entity.CancelledAt,
  )
  if execErr != nil {
    utils.SetSpanError(ctx, execErr)
    logger.Error("mysql.proposal.create.exec_error", "listing_identity_id", entity.ListingIdentityID, "err", execErr)
    return fmt.Errorf("create proposal: %w", execErr)
  }
  id, idErr := result.LastInsertId()
  if idErr != nil {
    utils.SetSpanError(ctx, idErr)
    return fmt.Errorf("proposal last insert id: %w", idErr)
  }
  proposal.SetID(id)
  proposal.SetCreatedAt(entity.CreatedAt)
  proposal.SetUpdatedAt(entity.UpdatedAt)
  return nil
}
```

- **Demais arquivos (mesma estrutura, cada um em arquivo próprio conforme Regra de 1 função pública por arquivo):**
  - `update_proposal_text.go`: `UPDATE proposals SET proposal_text = ? WHERE id = ? AND status = 'pending' AND deleted = 0`.
  - `update_proposal_status.go`: inclui cláusula `AND status = ?` para optimistic locking (`expected`).
  - `get_proposal_by_id.go`: `SELECT` das colunas explícitas com `deleted = 0`.
  - `get_proposal_by_id_for_update.go`: mesma query com `FOR UPDATE` para garantir consistência no service de status.
  - `list_proposals.go`: aplica filtros (`status`, `listing_identity_id`, `realtor_id`, `owner_id`, range de datas) e retorna `proposalmodel.ListResult` sem estatísticas financeiras.
  - `create_document.go`: grava `file_blob` (LONGBLOB), `file_name`, `mime_type`, `file_size_bytes`, `uploaded_at`.
  - `list_documents.go`: parâmetro `includeBlob` decide se o SELECT traz `file_blob` (detail) ou apenas metadados (list).

Cada função deve seguir o template da Seção 7.3: iniciar tracer, usar `InstrumentedAdapter`, logar erros infra em `slog`, retornar `sql.ErrNoRows` quando apropriado e documentar a query.

- **Arquivo:** [internal/core/port/right/repository/listing_repository/listing_repository_interface.go](internal/core/port/right/repository/listing_repository/listing_repository_interface.go)
```go
// UpdateProposalFlags updates listing_identity proposal indicators inside the same transaction used by proposal flows.
UpdateProposalFlags(ctx context.Context, tx *sql.Tx, input ProposalFlagsUpdate) error

type ProposalFlagsUpdate struct {
  ListingIdentityID int64
  HasPending        bool
  HasAccepted       bool
  AcceptedProposalID sql.NullInt64
}
```

- **Arquivo:** [internal/adapter/right/mysql/listing/update_proposal_flags.go](internal/adapter/right/mysql/listing/update_proposal_flags.go) deve executar `UPDATE listing_identities SET has_pending_proposal = ?, has_accepted_proposal = ?, accepted_proposal_id = ? WHERE id = ?` usando o mesmo padrão de tracing/logging.

### 2.5 Entities e Converters
- **Arquivo:** [internal/adapter/right/mysql/proposal/entities/proposal_entity.go](internal/adapter/right/mysql/proposal/entities/proposal_entity.go)
```go
package entities

import (
  "database/sql"
  "time"
)

// ProposalEntity mirrors proposals table after the schema change.
type ProposalEntity struct {
  ID                int64
  ListingIdentityID int64
  RealtorID         int64
  OwnerID           int64
  Status            string
  ProposalText      sql.NullString
  RejectionReason   sql.NullString
  AcceptedAt        sql.NullTime
  RejectedAt        sql.NullTime
  CancelledAt       sql.NullTime
  CreatedAt         time.Time
  UpdatedAt         time.Time
  Deleted           bool
}
```

- **Arquivo:** [internal/adapter/right/mysql/proposal/entities/proposal_document_entity.go](internal/adapter/right/mysql/proposal/entities/proposal_document_entity.go)
```go
package entities

import "time"

// ProposalDocumentEntity mirrors proposal_documents table including BLOB content.
type ProposalDocumentEntity struct {
  ID            int64
  ProposalID    int64
  FileName      string
  MimeType      string
  FileSizeBytes int64
  FileBlob      []byte
  UploadedAt    time.Time
}
```

- **Arquivo:** [internal/adapter/right/mysql/proposal/converters/proposal_entity_to_domain.go](internal/adapter/right/mysql/proposal/converters/proposal_entity_to_domain.go)
```go
package converters

import (
  "database/sql"

  "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/proposal/entities"
  proposalmodel "github.com/projeto-toq/toq_server/internal/core/model/proposal_model"
)

// ToProposalModel maps DB entity to domain struct, converting sql.Null* types.
func ToProposalModel(entity entities.ProposalEntity) proposalmodel.ProposalInterface {
  model := proposalmodel.NewProposal()
  model.SetID(entity.ID)
  model.SetListingIdentityID(entity.ListingIdentityID)
  model.SetRealtorID(entity.RealtorID)
  model.SetOwnerID(entity.OwnerID)
  model.SetStatus(proposalmodel.Status(entity.Status))
  model.SetProposalText(entity.ProposalText.String)
  model.SetRejectionReason(entity.RejectionReason)
  model.SetAcceptedAt(entity.AcceptedAt)
  model.SetRejectedAt(entity.RejectedAt)
  model.SetCancelledAt(entity.CancelledAt)
  model.SetCreatedAt(entity.CreatedAt)
  model.SetUpdatedAt(entity.UpdatedAt)
  return model
}

// ToProposalDocumentModel maps proposal_documents row to domain.
func ToProposalDocumentModel(entity entities.ProposalDocumentEntity, includeBlob bool) proposalmodel.ProposalDocumentInterface {
  doc := proposalmodel.NewProposalDocument()
  doc.SetID(entity.ID)
  doc.SetProposalID(entity.ProposalID)
  doc.SetFileName(entity.FileName)
  doc.SetMimeType(entity.MimeType)
  doc.SetFileSizeBytes(entity.FileSizeBytes)
  if includeBlob {
    doc.SetFileData(entity.FileBlob)
  }
  doc.SetUploadedAt(entity.UploadedAt)
  return doc
}
```

- **Arquivo:** [internal/adapter/right/mysql/proposal/converters/proposal_domain_to_entity.go](internal/adapter/right/mysql/proposal/converters/proposal_domain_to_entity.go)
```go
// ToProposalEntity converts domain to SQL entity, encapsulating sql.Null handling.
func ToProposalEntity(model proposalmodel.ProposalInterface) entities.ProposalEntity {
  return entities.ProposalEntity{
    ID:                model.ID(),
    ListingIdentityID: model.ListingIdentityID(),
    RealtorID:         model.RealtorID(),
    OwnerID:           model.OwnerID(),
    Status:            string(model.Status()),
    ProposalText:      sql.NullString{String: model.ProposalText(), Valid: model.ProposalText() != ""},
    RejectionReason:   model.RejectionReason(),
    AcceptedAt:        model.AcceptedAt(),
    RejectedAt:        model.RejectedAt(),
    CancelledAt:       model.CancelledAt(),
    CreatedAt:         model.CreatedAt(),
    UpdatedAt:         model.UpdatedAt(),
  }
}

// ToProposalDocumentEntity converts domain doc to DB entity.
func ToProposalDocumentEntity(doc proposalmodel.ProposalDocumentInterface) entities.ProposalDocumentEntity {
  return entities.ProposalDocumentEntity{
    ID:            doc.ID(),
    ProposalID:    doc.ProposalID(),
    FileName:      doc.FileName(),
    MimeType:      doc.MimeType(),
    FileSizeBytes: doc.FileSizeBytes(),
    FileBlob:      doc.FileData(),
    UploadedAt:    doc.UploadedAt(),
  }
}
```

### 2.6 Domain Models
- **Arquivo:** [internal/core/model/proposal_model/proposal.go](internal/core/model/proposal_model/proposal.go)
```go
package proposalmodel

import (
  "database/sql"
  "time"
)

// ProposalInterface now represents only the attributes required by the new feature.
type ProposalInterface interface {
  ID() int64
  SetID(int64)
  ListingIdentityID() int64
  SetListingIdentityID(int64)
  RealtorID() int64
  SetRealtorID(int64)
  OwnerID() int64
  SetOwnerID(int64)
  ProposalText() string
  SetProposalText(string)
  RejectionReason() sql.NullString
  SetRejectionReason(sql.NullString)
  Status() Status
  SetStatus(Status)
  AcceptedAt() sql.NullTime
  SetAcceptedAt(sql.NullTime)
  RejectedAt() sql.NullTime
  SetRejectedAt(sql.NullTime)
  CancelledAt() sql.NullTime
  SetCancelledAt(sql.NullTime)
  CreatedAt() time.Time
  SetCreatedAt(time.Time)
  UpdatedAt() time.Time
  SetUpdatedAt(time.Time)
}

type proposal struct {
  id                int64
  listingIdentityID int64
  realtorID         int64
  ownerID           int64
  proposalText      string
  rejectionReason   sql.NullString
  status            Status
  acceptedAt        sql.NullTime
  rejectedAt        sql.NullTime
  cancelledAt       sql.NullTime
  createdAt         time.Time
  updatedAt         time.Time
}

func NewProposal() ProposalInterface { return &proposal{} }

// Getters and setters omitted for brevity (same pattern as guia, focusing on new fields).
```

- **Arquivo:** [internal/core/model/proposal_model/document.go](internal/core/model/proposal_model/document.go)
```go
// ProposalDocumentInterface now includes the binary data (FileData).
type ProposalDocumentInterface interface {
  ID() int64
  SetID(int64)
  ProposalID() int64
  SetProposalID(int64)
  FileName() string
  SetFileName(string)
  MimeType() string
  SetMimeType(string)
  FileSizeBytes() int64
  SetFileSizeBytes(int64)
  FileData() []byte
  SetFileData([]byte)
  UploadedAt() time.Time
  SetUploadedAt(time.Time)
}
```

- **Arquivo:** [internal/core/model/proposal_model/filter.go](internal/core/model/proposal_model/filter.go)
```go
// ListFilter replaces legacy financial filters by actor-scoped parameters.
type ListFilter struct {
  ActorScope   ActorScope // enum {realtor, owner}
  ActorID      int64
  ListingID    *int64
  Statuses     []Status
  Page         int
  Limit        int
}

type ListResult struct {
  Items []ProposalInterface
  Total int64
}
```

- **Arquivo:** [internal/core/model/proposal_model/status.go](internal/core/model/proposal_model/status.go)
```go
const (
  StatusPending  Status = "pending"
  StatusAccepted Status = "accepted"
  StatusRefused  Status = "refused"
  StatusCancelled Status = "cancelled"
)
```

### 2.7 Banco de Dados (scripts/db_creation.sql)
```sql
CREATE TABLE `proposals` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `listing_identity_id` INT UNSIGNED NOT NULL,
  `realtor_id` INT UNSIGNED NOT NULL,
  `owner_id` INT UNSIGNED NOT NULL,
  `status` ENUM('pending','accepted','refused','cancelled') NOT NULL DEFAULT 'pending',
  `proposal_text` TEXT NULL,
  `rejection_reason` VARCHAR(500) NULL,
  `accepted_at` DATETIME(6) NULL,
  `rejected_at` DATETIME(6) NULL,
  `cancelled_at` DATETIME(6) NULL,
  `deleted` TINYINT(1) NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  KEY `idx_proposals_realtor` (`realtor_id`,`status`),
  KEY `idx_proposals_owner` (`owner_id`,`status`),
  KEY `idx_proposals_listing_status` (`listing_identity_id`,`status`),
  CONSTRAINT `fk_proposals_listing` FOREIGN KEY (`listing_identity_id`) REFERENCES `listing_identities`(`id`),
  CONSTRAINT `fk_proposals_realtor` FOREIGN KEY (`realtor_id`) REFERENCES `users`(`id`),
  CONSTRAINT `fk_proposals_owner` FOREIGN KEY (`owner_id`) REFERENCES `users`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `proposal_documents` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `proposal_id` BIGINT UNSIGNED NOT NULL,
  `file_name` VARCHAR(255) NOT NULL,
  `mime_type` VARCHAR(60) NOT NULL DEFAULT 'application/pdf',
  `file_size_bytes` BIGINT NOT NULL,
  `file_blob` LONGBLOB NOT NULL,
  `uploaded_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  PRIMARY KEY (`id`),
  KEY `idx_proposal_documents_proposal` (`proposal_id`),
  CONSTRAINT `fk_proposal_documents_proposal` FOREIGN KEY (`proposal_id`) REFERENCES `proposals`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

ALTER TABLE `listing_identities`
  ADD COLUMN `has_pending_proposal` TINYINT(1) NOT NULL DEFAULT 0 AFTER `owner_last_response_at`,
  ADD COLUMN `has_accepted_proposal` TINYINT(1) NOT NULL DEFAULT 0 AFTER `has_pending_proposal`,
  ADD COLUMN `accepted_proposal_id` BIGINT UNSIGNED NULL AFTER `has_accepted_proposal`,
  ADD CONSTRAINT `fk_listing_identity_accepted_proposal` FOREIGN KEY (`accepted_proposal_id`) REFERENCES `proposals`(`id`) ON DELETE SET NULL;
```

O controle temporal fica restrito às colunas de status (`accepted_at`, `rejected_at`, `cancelled_at`) e à auditoria do serviço global, já que o schema não inclui `created_at`/`updated_at` para `proposals`.

## 3. Estrutura de Diretórios
```
internal/
  adapter/
    left/http/dto/proposal_dto.go
    left/http/handlers/proposal_handlers/
      proposal_handler.go
      create_proposal_handler.go
      update_proposal_handler.go
      cancel_proposal_handler.go
      accept_proposal_handler.go
      reject_proposal_handler.go
      list_realtor_proposals_handler.go
      list_owner_proposals_handler.go
      get_proposal_detail_handler.go
    right/mysql/proposal/
      proposal_adapter.go
      create_proposal.go
      update_proposal_text.go
      update_proposal_status.go
      get_proposal_by_id.go
      get_proposal_by_id_for_update.go
      fetch_proposal.go
      list_proposals.go
      create_document.go
      list_documents.go
      converters/
        proposal_entity_to_domain.go
        proposal_domain_to_entity.go
      entities/
        proposal_entity.go
        proposal_document_entity.go
    right/mysql/listing/update_proposal_flags.go
  core/
    model/proposal_model/
      proposal.go
      document.go
      filter.go
      status.go
    port/
      left/http/proposalhandler/proposal_handler_port.go
      right/repository/proposal_repository/proposal_repository_interface.go
    service/proposal_service/
      proposal_service.go
      types.go
      create_proposal.go
      update_proposal.go
      cancel_proposal.go
      accept_proposal.go
      reject_proposal.go
      list_realtor_proposals.go
      list_owner_proposals.go
      get_proposal_detail.go
      helpers.go
docs/
  proposals_implementation_plan.md
scripts/
  db_creation.sql
```
Essa organização mantém a Regra de Espelhamento (Seção 2.1 do guia) e a regra "uma função pública por arquivo".

## 4. Ordem de Execução
1. **Persistência** – O DBA deverá atualizar `scripts/db_creation.sql` com as novas tabelas/colunas.
2. **Domínio** – Refatorar `internal/core/model/proposal_model/*` para remover campos financeiros, adicionar texto/PDF e atualizar `status`. Ajustar o port `proposal_repository` e criar as novas estruturas de filtro.
3. **Adapter MySQL** – Reescrever `internal/adapter/right/mysql/proposal/*` seguindo o novo contrato, adicionando suporte a transações (`*sql.Tx`), BLOBs e filtros mínimos. Criar `update_proposal_flags.go` no adapter de listings.
4. **Serviço** – Implementar `internal/core/service/proposal_service/*`, compondo com `global_service` para transações, auditoria e notificações. Garantir observabilidade (tracing/log) conforme Seção 7 do guia.
5. **DTOs/Handlers** – Criar DTOs e handlers em `proposal_handlers`, adicionando Swagger annotations e validações (Seção 8.2). Reutilizar `httperrors.SendHTTPErrorObj` e extrair contexto do request conforme guia.
6. **Ports/Factories** – Expor `proposalhandler.Handler` no port esquerdo, ajustar factories (AdapterFactory e HTTP Factory) para injetar repository/service/handler, atualizar rotas `/api/v2/proposals/*`.
7. **Observabilidade/Notificações** – Conectar `globalService.GetUnifiedNotificationService()` em transições de status, criar auditorias e logs `slog.Info` para eventos de domínio.
8. **Documentação e validação** – Rodar `make swagger`, revisar diffs das novas rotas e garantir que a documentação (este arquivo + README se necessário) descreva os novos fluxos. Somente após isso liberar para QA.

## 5. Status de Implementação (2026-01-06)
- [x] Domínio atualizado (modelos Proposal, ProposalDocument, filtros e status) em `internal/core/model/proposal_model`.
- [x] Repositório MySQL migrado para o novo contrato (`CreateProposal`, `UpdateProposalText`, `UpdateProposalStatus`, `GetProposalByID`, `ListProposals`, `CreateDocument`, `ListDocuments`).
- [x] Serviços de aplicação implementados em `internal/core/service/proposal_service` com DTOs e handlers já integrados conforme Seções 2.2/2.3.
- [x] Integração com notificações/auditoria e sincronização das flags no adapter de listings através das novas operações do serviço.
- [x] Execução de `make lint` e validações finais concluída em 2026-01-06.
