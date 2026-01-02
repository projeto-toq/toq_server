# Plano de Refatoração — Fluxo de Visitas

## Contexto e Regras
- Seguir guia em docs/toq_server_go_guide.md (Seções 2, 5, 7, 8).
- Sem criação de campos de auditoria; usar apenas campos existentes (ex.: first_owner_action_at, métricas em listing_identities).
- Métricas de resposta do owner já em listing_identities: owner_avg_response_time_seconds, owner_total_visits_responded, owner_last_response_at.
- Permission service já controla papéis; handlers apenas consomem contexto autenticado.
- Não editar swagger.yaml/json; anotar via comentários no código.
- Proibido criar/alterar testes neste escopo.

## Fases e Sequência
1) **Config & contracts**
    - ✅ Config visits.min_hours_ahead e visits.max_days_ahead criada no env e carregada no Environment.
    - ✅ Propagação da config para visitService via factory/InitVisitService.
    - ✅ Ajustar DTOs/Swagger (handlers) para documentar regras de janela e headers padrão.

2) **Disponibilidade e conflitos**
    - ✅ Service: validar lead time/horizonte antes de persistir.
    - ✅ Usar scheduleRepo.GetAvailabilityData + rules/entries para recusar slots fora da disponibilidade.
    - ✅ Manter checagem de conflitos por entries (preservada no fluxo).

3) **Bloqueio bilateral de agenda**
    - ✅ Além de agenda do owner, criar/atualizar entry equivalente para o realtor (pending/confirmed/cancel/unblock).
    - ✅ Harmonizar EntryType para ambos (VISIT_PENDING, VISIT_CONFIRMED; unblock em CANCEL/NO_SHOW/COMPLETE quando aplicável).

4) **Notificações**
    - ✅ Integrar UnifiedNotificationService em Create/Approve/Reject/Cancel/Complete/NoShow (payloads específicos para owner/realtor).

5) **Métrica de resposta do owner**
     - ✅ Na primeira ação do owner (approve/reject; cancel se feito pelo owner):
         - ✅ Calcular delta = first_owner_action_at - visit.created_at (usar first_owner_action_at já persistido).
         - ✅ Atualizar listing_identities: avg = (avg*total + delta)/(total+1); total++; last_response_at = now.
     - ✅ Evitar recálculo se first_owner_action_at já setado.

6) **Listagens e filtros**
    - ✅ Garantir filtros por status/type/time aplicados em ListVisits; manter paginação segura.
    - ✅ Converters de resposta devem incluir first_owner_action_at quando existir.

7) **Documentação e revisão**
    - ✅ Revisar comentários Godoc/Swagger conforme Seção 8.
   - Conferir regra de “um método público por arquivo”.

## Arquivos-alvo (por fase)
- Config: configs/env.yaml, internal/core/config/* (bootstrap/factory), internal/core/service/visit_service/visit_service.go.
- DTO/Handlers: internal/adapter/left/http/dto/visit_request_dto.go, visit_response_dto.go; internal/adapter/left/http/handlers/visit_handlers/*.go.
- Services: internal/core/service/visit_service/{create_visit.go, approve_visit.go, reject_visit.go, cancel_visit.go, complete_visit.go, no_show_visit.go, list_visits.go, get_visit.go, helpers.go, validate_window.go (novo), notify_visit_events.go (novo), response_metrics.go (novo)}.
- Schedule: uso de schedule_repo interfaces já existentes; nenhuma mudança estrutural esperada.
- Listing repo: método para atualizar métricas de resposta em listing_identities (novo método público + adapter/port).
- Converters: internal/adapter/left/http/converters/visit_dto_converter.go (incluir first_owner_action_at).

## Skeletons (referência rápida)

### Handler CreateVisit (comentários Swagger conforme Seção 8.2)
```go
// RequestVisit ...
// @Summary Request a visit
// @Description Validates lead time/window, checks availability, and notifies owner.
// @Tags Visits
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param X-Device-Id header string false "Device ID"
// @Param request body dto.CreateVisitRequest true "Visit payload"
// @Success 201 {object} dto.VisitResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /visits [post]
func (h *VisitHandler) RequestVisit(c *gin.Context) {
    ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
    var req dto.CreateVisitRequest
    if err := c.ShouldBindJSON(&req); err != nil { httperrors.SendHTTPErrorObj(c, httputils.MapBindingError(err)); return }
    input, err := converters.CreateVisitDTOToInput(req)
    if err != nil { httperrors.SendHTTPErrorObj(c, err); return }
    visit, svcErr := h.visitService.CreateVisit(ctx, input)
    if svcErr != nil { httperrors.SendHTTPErrorObj(c, svcErr); return }
    c.JSON(http.StatusCreated, converters.VisitDomainToResponse(visit))
}
```

### Service CreateVisit (fluxo)
```go
func (s *visitService) CreateVisit(ctx context.Context, input CreateVisitInput) (listingmodel.VisitInterface, error) {
    ctx, end, err := utils.GenerateTracer(ctx); if err != nil { return nil, derrors.Infra("trace", err) }
    defer end()
    ctx = utils.ContextWithLogger(ctx)
    if err = s.validateWindow(ctx, input); err != nil { return nil, err }
    requesterID, err := s.globalService.GetUserIDFromContext(ctx); if err != nil { return nil, err }

    tx, err := s.globalService.StartTransaction(ctx); if err != nil { utils.SetSpanError(ctx, err); return nil, derrors.Infra("tx start", err) }
    defer s.globalService.SafeRollback(ctx, tx, &err)

    listingIdentity, activeVersion, agenda, err := s.loadListingAndAgenda(ctx, tx, input.ListingIdentityID)
    if err != nil { return nil, err }
    if err = s.ensureAvailability(ctx, tx, agenda, input); err != nil { return nil, err }

    visit := buildVisitDomain(input, requesterID, listingIdentity.UserID, activeVersion.Version())
    if _, err = s.visitRepo.InsertVisit(ctx, tx, visit); err != nil { utils.SetSpanError(ctx, err); return nil, derrors.Infra("insert visit", err) }

    if err = s.blockOwnerAgenda(ctx, tx, visit, schedulemodel.EntryTypeVisitPending); err != nil { return nil, err }
    if err = s.blockRealtorAgenda(ctx, tx, visit, schedulemodel.EntryTypeVisitPending); err != nil { return nil, err }

    if err = s.globalService.CommitTransaction(ctx, tx); err != nil { utils.SetSpanError(ctx, err); return nil, derrors.Infra("tx commit", err) }
    s.notifyOwnerAsync(ctx, visit)
    return visit, nil
}
```

### Service Approve/Reject (trecho para métrica owner)
```go
if _, ok := visit.FirstOwnerActionAt(); !ok {
    now := time.Now().UTC()
    visit.SetFirstOwnerActionAt(now)
    if err := s.updateOwnerResponseStats(ctx, tx, visit, now); err != nil { return nil, err }
}
```

### Repositório de listing (novo método)
```go
// UpdateOwnerResponseStats updates aggregated response metrics for a listing identity.
func (r *ListingRepository) UpdateOwnerResponseStats(ctx context.Context, tx *sql.Tx, identityID int64, deltaSeconds int64, respondedAt time.Time) error {
    // UPDATE listing_identities SET owner_avg_response_time_seconds = ?, owner_total_visits_responded = ?, owner_last_response_at = ? WHERE id = ?
}
```

### DTO CreateVisitRequest (campos)
```go
type CreateVisitRequest struct {
    ListingIdentityID int64  `json:"listingIdentityId" binding:"required"`
    ScheduledStart    string `json:"scheduledStart" binding:"required,datetime=2006-01-02T15:04:05Z07:00"`
    ScheduledEnd      string `json:"scheduledEnd" binding:"required,datetime=2006-01-02T15:04:05Z07:00"`
    Type              string `json:"type" binding:"required,oneof=WITH_CLIENT REALTOR_ONLY CONTENT_PRODUCTION"`
    RealtorNotes      string `json:"realtorNotes,omitempty" binding:"max=2000"`
    Source            string `json:"source,omitempty" binding:"omitempty,oneof=APP WEB ADMIN"`
}
```

### Converter VisitDomainToResponse
```go
func VisitDomainToResponse(v listingmodel.VisitInterface) dto.VisitResponse {
    // map básicos
    if ts, ok := v.FirstOwnerActionAt(); ok {
        formatted := ts.Format(time.RFC3339)
        resp.FirstOwnerActionAt = &formatted
    }
    return resp
}
```

## Checkpoints de entrega
- Após fase 1: config carregada e DTO/Swagger ajustados.
- Após fase 2/3: criação com disponibilidade + bloqueio bilateral funcional.
- Após fase 4: notificações enviadas em todas as transições.
- Após fase 5: métricas de resposta do owner atualizando listing_identities.
- Após fase 6: respostas retornando first_owner_action_at.
- Após fase 7: documentação e revisão de conformidade (Seções 7 e 8).
