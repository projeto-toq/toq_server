# Plano de Implementação do Sistema de Visitas

## Status Atual
✅ Domínios de visita corrigidos (visit_status.go, visit_interface.go, visit_domain.go)
✅ Adapter MySQL de visita corrigido (entities, converters, queries)
✅ Helpers de agenda para visitas implementados (Create/Update/Delete/Check) com transações alinhadas ao GlobalService
✅ Schedule repo/adapter exposto para lookup por visit_id
✅ VisitService criado (create/approve/reject/cancel/complete/no-show/get/list) com transações e integração de agenda
✅ Lint (`make lint`) está passando
⏳ Alterações de banco pendentes (serão aplicadas pelo time de DBA)

## Etapas Pendentes

- Aplicar alterações de banco (DBA) e validar consistência das novas colunas/enum
- Atualizar DTOs, converters e handlers HTTP de visita (+ rotas) seguindo o espelhamento
- Revisar métricas/telemetria e notificações específicas de visita

---

## Etapa 1: Atualização do Schema do Banco de Dados

### 1.1 Instruções para DBA (não aplicar no repositório)

**Contexto:** O repositório não deve ser alterado; o time de DBA executará as mudanças diretamente no banco. Enviar a eles o script abaixo.

**Mudanças necessárias em listing_visits:**

```sql
CREATE TABLE IF NOT EXISTS `toq_db`.`listing_visits` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `listing_identity_id` INT UNSIGNED NOT NULL,
  `listing_version` INT UNSIGNED NOT NULL DEFAULT 1,
  `user_id` INT UNSIGNED NOT NULL,
  `scheduled_date` DATE NOT NULL,
  `scheduled_time_start` TIME NOT NULL,
  `scheduled_time_end` TIME NOT NULL,
  `status` ENUM('PENDING', 'APPROVED', 'REJECTED', 'CANCELLED', 'COMPLETED', 'NO_SHOW') NOT NULL DEFAULT 'PENDING',
  `source` ENUM('APP', 'WEB', 'ADMIN') NOT NULL DEFAULT 'APP',
    `notes` TEXT NULL,
    `rejection_reason` VARCHAR(255) NULL,
    `first_owner_action_at` DATETIME NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_visits_listing_identity_idx` (`listing_identity_id` ASC) VISIBLE,
  INDEX `fk_visits_user_idx` (`user_id` ASC) VISIBLE,
  INDEX `idx_scheduled_date` (`scheduled_date` ASC) VISIBLE,
  INDEX `idx_status` (`status` ASC) VISIBLE,
  CONSTRAINT `fk_visits_listing_identity`
    FOREIGN KEY (`listing_identity_id`)
    REFERENCES `toq_db`.`listing_identities` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_visits_user`
    FOREIGN KEY (`user_id`)
    REFERENCES `toq_db`.`users` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB;
```

**Mudanças necessárias em listing_identities (métricas de resposta do owner):**

```sql
-- Campos para tracking de tempo de resposta do owner
ALTER TABLE `toq_db`.`listing_identities`
    ADD COLUMN `owner_avg_response_time_seconds` INT UNSIGNED NULL 
        COMMENT 'Tempo médio de resposta do owner em segundos (EMA)' 
        AFTER `active_version_id`,
    ADD COLUMN `owner_total_visits_responded` INT UNSIGNED NOT NULL DEFAULT 0 
        COMMENT 'Total de visitas respondidas pelo owner' 
        AFTER `owner_avg_response_time_seconds`,
    ADD COLUMN `owner_last_response_at` DATETIME NULL 
        COMMENT 'Timestamp da última resposta do owner' 
        AFTER `owner_total_visits_responded`;
```

**Arquivo de migração (referência para o DBA executar):** `scripts/migrations/add_visit_fields.sql`

```sql
-- Renomear listing_id para listing_identity_id (consistência arquitetural)
ALTER TABLE `toq_db`.`listing_visits`
  CHANGE COLUMN `listing_id` `listing_identity_id` INT UNSIGNED NOT NULL;

-- Remover campos desnecessários
ALTER TABLE `toq_db`.`listing_visits`
  DROP COLUMN IF EXISTS `owner_id`,
  DROP COLUMN IF EXISTS `realtor_id`;

-- Adicionar novos campos à tabela listing_visits
ALTER TABLE `toq_db`.`listing_visits`
  ADD COLUMN `listing_version` INT UNSIGNED NOT NULL DEFAULT 1 AFTER `listing_identity_id`,
  ADD COLUMN `source` ENUM('APP', 'WEB', 'ADMIN') NOT NULL DEFAULT 'APP' AFTER `status`,
    ADD COLUMN `first_owner_action_at` DATETIME NULL AFTER `rejection_reason`;

-- Atualizar ENUM de status
ALTER TABLE `toq_db`.`listing_visits`
  MODIFY COLUMN `status` ENUM('PENDING', 'APPROVED', 'REJECTED', 'CANCELLED', 'COMPLETED', 'NO_SHOW') NOT NULL DEFAULT 'PENDING';

-- Atualizar registros existentes
# Plano de Implementação do Sistema de Visitas (Nova Versão)

## 1. Diagnóstico
- Arquivos analisados: [docs/toq_server_go_guide.md](docs/toq_server_go_guide.md), [docs/visits_system_specification.md](docs/visits_system_specification.md), [internal/core/model/listing_model/visit_domain.go](internal/core/model/listing_model/visit_domain.go), [internal/core/model/listing_model/visit_filters.go](internal/core/model/listing_model/visit_filters.go), [internal/core/model/listing_model/visit_status.go](internal/core/model/listing_model/visit_status.go), [internal/core/model/schedule_model/agenda_domain.go](internal/core/model/schedule_model/agenda_domain.go), [internal/core/model/schedule_model/enums.go](internal/core/model/schedule_model/enums.go), [internal/adapter/right/mysql/visit](internal/adapter/right/mysql/visit), [internal/core/port/right/repository/visit_repository/visit_repository_interface.go](internal/core/port/right/repository/visit_repository/visit_repository_interface.go), [internal/core/service/listing_service/approve_visit.go](internal/core/service/listing_service/approve_visit.go), [internal/core/service/listing_service/reject_visit.go](internal/core/service/listing_service/reject_visit.go), [internal/core/service/listing_service/cancel_visit.go](internal/core/service/listing_service/cancel_visit.go), [internal/core/service/listing_service/get_all_visits_by_user.go](internal/core/service/listing_service/get_all_visits_by_user.go), [scripts/db_creation.sql](scripts/db_creation.sql).
- Problemas encontrados:
    - Domain `visit` é rascunho: não contém owner ID, tipo de visita, duração, notas, response time ou campos para métricas; não há validação de fuso/agenda ([internal/core/model/listing_model/visit_domain.go](internal/core/model/listing_model/visit_domain.go)).
    - Ports/Services para visita estão vazios/stubados: métodos em `listing_service` não possuem lógica ([internal/core/service/listing_service/approve_visit.go](internal/core/service/listing_service/approve_visit.go)). Não há service dedicado de visitas.
    - Filtros de listagem divergem do adapter: `VisitListFilter` usa `ListingID/OwnerID/RealtorID` enquanto o adapter usa `listing_identity_id/user_id/status` ([internal/core/model/listing_model/visit_filters.go](internal/core/model/listing_model/visit_filters.go) vs [internal/adapter/right/mysql/visit/list_visits.go](internal/adapter/right/mysql/visit/list_visits.go)).
    - Schema atual está obsoleto: tabela `listing_visits` usa colunas `listing_id/owner_id/realtor_id` e enum `PENDING_OWNER/CONFIRMED/DONE`, sem `listing_version`, sem tipos ou métricas, e `listing_agenda_entries` usa `starts_at/ends_at` com `visit_id` sem constraints de unicidade ([scripts/db_creation.sql](scripts/db_creation.sql#L610-L660)).
    - Agenda de disponibilidade existe, mas não há integração real com visitas (nenhum helper em schedule para criar/atualizar entries de visita; service atual não chama schedule).
    - Requisito de negócio de notificar owner/realtor não está mapeado em nenhum handler/service.
- Impacto: precisamos criar fluxo completo (domínio, DTOs, handlers, service, repository) alinhado ao guia (Seção 2.1 espelhamento, Seção 7/8 documentação, uma função pública por arquivo) e ajustar schema para novos campos/enum. Alteração disruptiva, sem compatibilidade retroativa.
- Melhorias adicionais: normalizar filtros (owner/realtor/listing), alinhar enums (`PENDING/APPROVED/REJECTED/CANCELLED/COMPLETED/NO_SHOW` + `WITH_CLIENT/REALTOR_ONLY/CONTENT_PRODUCTION`), registrar métricas de resposta do owner no próprio `listing_identity`, e usar agenda para bloqueio de horários.

## 2. Code Skeletons (prontos para implementar)
> Código em inglês, um método público por arquivo, obedecendo templates da Seção 8 do guia. Apenas assinaturas e estrutura; sem regra de negócio implementada.

### 2.1 Domínio e Dados
- Atualizar domínio de visita ([internal/core/model/listing_model/visit_domain.go](internal/core/model/listing_model/visit_domain.go)) e filtros ([internal/core/model/listing_model/visit_filters.go](internal/core/model/listing_model/visit_filters.go)):
```go
package listingmodel

type VisitType string
const (
        VisitTypeWithClient      VisitType = "WITH_CLIENT"
        VisitTypeRealtorOnly     VisitType = "REALTOR_ONLY"
        VisitTypeContentProduction VisitType = "CONTENT_PRODUCTION"
)

        nil,

func NewVisit() Visit { return &visit{} }

type VisitListFilter struct {
        ListingIdentityID *int64
        OwnerUserID       *int64
        RequesterUserID   *int64
        Statuses          []VisitStatus
        Types             []VisitType
        From              *time.Time
        To                *time.Time
        Page              int
        Limit             int
}

        time.Now(),
        Visits []Visit
        Total  int64
}
```
- Novo enum status já existe; manter `IsBlocking()` para `PENDING/APPROVED` ([internal/core/model/listing_model/visit_status.go](internal/core/model/listing_model/visit_status.go)).

### 2.2 DTOs (HTTP Left)
- Request/response em [internal/adapter/left/http/dto/visit](internal/adapter/left/http/dto/visit):
```go
package dto

// CreateVisitRequest carries visit creation data
// Fields validated via binding tags; times in RFC3339 with timezone.
        time.Now(),
        ListingIdentityID int64  `json:"listingIdentityId" binding:"required" example:"123"`
        ScheduledStart    string `json:"scheduledStart" binding:"required,datetime=2006-01-02T15:04:05Z07:00" example:"2025-01-10T14:00:00Z"`
        ScheduledEnd      string `json:"scheduledEnd" binding:"required,datetime=2006-01-02T15:04:05Z07:00" example:"2025-01-10T14:30:00Z"`
        Type              string `json:"type" binding:"required,oneof=WITH_CLIENT REALTOR_ONLY CONTENT_PRODUCTION" example:"WITH_CLIENT"`
        RealtorNotes      string `json:"realtorNotes" binding:"max=2000" example:"Client has time only in the afternoon"`
}

// UpdateVisitStatusRequest centralizes approve/reject/cancel/complete/no-show
// Action determines which optional fields are mandatory.
    )
        VisitID         int64  `json:"visitId" binding:"required" example:"456"`
        Action          string `json:"action" binding:"required,oneof=APPROVE REJECT CANCEL COMPLETE NO_SHOW" example:"APPROVE"`
        OwnerNotes      string `json:"ownerNotes" binding:"max=2000" example:"Ring the bell"`
        RejectionReason string `json:"rejectionReason" binding:"max=2000" example:"Slot unavailable"`
        CancelReason    string `json:"cancelReason" binding:"max=2000" example:"Client emergency"`
}


        ID                int64   `json:"id" example:"456"`
        ListingIdentityID int64   `json:"listingIdentityId" example:"123"`
        ListingVersion    uint8   `json:"listingVersion" example:"3"`
        RequesterUserID   int64   `json:"requesterUserId" example:"5"`
        OwnerUserID       int64   `json:"ownerUserId" example:"10"`
        ScheduledStart    string  `json:"scheduledStart" example:"2025-01-10T14:00:00Z"`
        ScheduledEnd      string  `json:"scheduledEnd" example:"2025-01-10T14:30:00Z"`
        DurationMinutes   int64   `json:"durationMinutes" example:"30"`
        Status            string  `json:"status" example:"PENDING"`
        Type              string  `json:"type" example:"WITH_CLIENT"`
        RealtorNotes      string  `json:"realtorNotes,omitempty"`
        OwnerNotes        string  `json:"ownerNotes,omitempty"`
        RejectionReason   string  `json:"rejectionReason,omitempty"`
        CancelReason      string  `json:"cancelReason,omitempty"`
        FirstOwnerActionAt *string `json:"firstOwnerActionAt,omitempty" example:"2025-01-10T14:05:00Z"`
        CreatedAt         string  `json:"createdAt" example:"2025-01-09T12:00:00Z"`
        UpdatedAt         string  `json:"updatedAt" example:"2025-01-09T12:00:00Z"`
}

type VisitListResponse struct {
        Items      []VisitResponse `json:"items"`
        Total      int64           `json:"total" example:"120"`
        Page       int             `json:"page" example:"1"`
        PageSize   int             `json:"pageSize" example:"20"`
}
```

### 2.3 Converters (HTTP DTO ↔ Domínio)
- [internal/adapter/left/http/dto/converters/visit_dto_converter.go](internal/adapter/left/http/dto/converters/visit_dto_converter.go):
```go
package converters

func CreateVisitDTOToDomain(req *dto.CreateVisitRequest, requesterID int64, activeVersion uint8, ownerID int64) (listingmodel.Visit, error) {
        // parse times, set duration, set created/updated by, status PENDING, type
}

func VisitDomainToResponse(v listingmodel.Visit) dto.VisitResponse {
        // map all fields, format time RFC3339, optional pointers
}

func VisitDomainsToResponse(items []listingmodel.Visit, total int64, page, pageSize int) dto.VisitListResponse {
        // iterate and build list response
}
```

### 2.4 Repository Port e Adapter (MySQL)
- Port: atualizar [internal/core/port/right/repository/visit_repository/visit_repository_interface.go](internal/core/port/right/repository/visit_repository/visit_repository_interface.go):
```go
    return visit, nil
        InsertVisit(ctx context.Context, tx *sql.Tx, visit listingmodel.Visit) (int64, error)
        UpdateVisit(ctx context.Context, tx *sql.Tx, visit listingmodel.Visit) error
        GetVisitByID(ctx context.Context, tx *sql.Tx, id int64) (listingmodel.Visit, error)
        ListVisits(ctx context.Context, tx *sql.Tx, filter listingmodel.VisitListFilter) (listingmodel.VisitListResult, error)
}
```
- Adapter (espelhamento) em [internal/adapter/right/mysql/visit](internal/adapter/right/mysql/visit) (um método por arquivo, queries sem SELECT *):
```go
// insert_visit.go
func (a *VisitAdapter) InsertVisit(ctx context.Context, tx *sql.Tx, visit listingmodel.Visit) (int64, error) {
        // INSERT listing_identity_id, listing_version, requester_user_id, owner_user_id, scheduled_start, scheduled_end, duration_minutes, status, type, notes, rejection_reason, cancel_reason, first_owner_action_at
}

// update_visit.go
func (a *VisitAdapter) UpdateVisit(ctx context.Context, tx *sql.Tx, visit listingmodel.Visit) error {
        // UPDATE same set; rows affected check
}

// get_visit_by_id.go
func (a *VisitAdapter) GetVisitByID(ctx context.Context, tx *sql.Tx, id int64) (listingmodel.Visit, error) {
        // SELECT explicit columns; scan via mapper; convert entity→domain
}

// list_visits.go
func (a *VisitAdapter) ListVisits(ctx context.Context, tx *sql.Tx, filter listingmodel.VisitListFilter) (listingmodel.VisitListResult, error) {
        // build WHERE for listing_identity_id, owner_user_id, requester_user_id, statuses, types, from/to; COUNT + SELECT with LIMIT/OFFSET
}

// entities/visit_entity.go
// converters/visit_entity_to_domain.go
// converters/visit_domain_to_entity.go
```
- Manter `defaultPagination` e `visit_row_mapper` com nova coluna order.

### 2.5 Service de Visitas (Core)
- Novo service dedicado em [internal/core/service/visit_service](internal/core/service/visit_service) obedecendo template da Seção 8 (um método público por arquivo, tracing, transação, logs de infra):
```go
// visit_service.go
package visitservice

type VisitService interface {
        CreateVisit(ctx context.Context, input CreateVisitInput) (listingmodel.Visit, error)
        ApproveVisit(ctx context.Context, visitID, ownerID int64, ownerNotes string) (listingmodel.Visit, error)
        RejectVisit(ctx context.Context, visitID, ownerID int64, reason string) (listingmodel.Visit, error)
        CancelVisit(ctx context.Context, visitID, requesterID int64, reason string) (listingmodel.Visit, error)
        CompleteVisit(ctx context.Context, visitID, ownerID int64, notes string) (listingmodel.Visit, error)
        MarkNoShow(ctx context.Context, visitID, ownerID int64, notes string) (listingmodel.Visit, error)
        ListVisits(ctx context.Context, filter listingmodel.VisitListFilter) (listingmodel.VisitListResult, error)
        GetVisit(ctx context.Context, visitID int64, requesterID int64) (listingmodel.Visit, error)
}

type visitService struct {
        globalService globalservice.GlobalServiceInterface
        visitRepo     visitrepository.VisitRepositoryInterface
        listingRepo   listingrepository.ListingRepoPortInterface
        scheduleRepo  schedulerepository.ScheduleRepositoryInterface
        notification  notificationservice.NotificationServiceInterface
}

func NewVisitService(...) VisitService { /* wiring only */ }
```
- Métodos públicos (arquivos separados: `create_visit.go`, `approve_visit.go`, `reject_visit.go`, `cancel_visit.go`, `complete_visit.go`, `mark_no_show.go`, `get_visit.go`, `list_visits.go`):
```go
// create_visit.go
}
        ListingIdentityID int64
        RequesterUserID   int64
        ScheduledStart    time.Time
        ScheduledEnd      time.Time
        Type              listingmodel.VisitType
        RealtorNotes      string
}

func (s *visitService) CreateVisit(ctx context.Context, input CreateVisitInput) (listingmodel.Visit, error) {
        // tracer; validate listing identity + owner; compute active version; validate availability via schedule; start tx; insert visit; create agenda entry VISIT_PENDING blocking; commit; send notification owner post-commit
}
```
(Seguem os demais métodos com Godoc descrevendo fluxo, validações de status, bloqueio/desbloqueio de agenda, métricas de resposta, e envio de notificações; sempre marcar infra com `utils.SetSpanError`).

### 2.6 Integração com Agenda (Schedule Service)
- Adicionar helpers no schedule service (arquivos separados) para reaproveitar entry model [internal/core/model/schedule_model](internal/core/model/schedule_model):
```go
// create_visit_entry.go
func (s *scheduleService) CreateVisitEntry(ctx context.Context, agendaID uint64, visitID uint64, start, end time.Time, pending bool) error

// update_visit_entry_type.go
func (s *scheduleService) UpdateVisitEntryType(ctx context.Context, visitID uint64, newType schedulemodel.EntryType) error

// delete_visit_entry.go
func (s *scheduleService) DeleteVisitEntry(ctx context.Context, visitID uint64) error

// check_visit_conflict.go
func (s *scheduleService) CheckVisitConflict(ctx context.Context, agendaID uint64, start, end time.Time, excludeVisitID *uint64) (bool, error)
```
- Repositório de schedule precisa de método `GetAgendaEntryByVisitID` (adapter MySQL) respeitando espelhamento.

### 2.7 Handlers HTTP e Port Left
- Criar novo port [internal/core/port/left/http/visithandler/visit_handler_port.go](internal/core/port/left/http/visithandler/visit_handler_port.go) com métodos `CreateVisit`, `ListVisitsRealtor`, `ListVisitsOwner`, `GetVisit`, `UpdateVisitStatus`.
- Handlers (um por arquivo) em [internal/adapter/left/http/handlers/visit_handler](internal/adapter/left/http/handlers/visit_handler) com Swagger completo; nenhum POST com id em path (IDs sempre no body):
```go
// create_visit_handler.go
// @Summary Request a visit
// @Description Creates a visit request using listing availability. IDs in body.
// @Tags Visits
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param request body dto.CreateVisitRequest true "Visit payload"
// @Success 201 {object} dto.VisitResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Router /api/v2/visits [post]
func (h *VisitHandler) CreateVisit(c *gin.Context) { /* bind → convert → service → response */ }

// update_status_handler.go (approve/reject/cancel/complete/no-show via body.action)
// @Router /api/v2/visits/status [post]

// list_realtor_handler.go
// @Router /api/v2/visits/realtor [get]

// list_owner_handler.go
// @Router /api/v2/visits/owner [get]

// get_visit_handler.go
// @Router /api/v2/visits/{visitId} [get]
```
- Rotas centralizadas em [internal/adapter/left/http/routes/routes.go](internal/adapter/left/http/routes/routes.go): criar helpers ou blocos explícitos no grupo `/api/v2/visits` (seguindo padrão do arquivo), registrando `CreateVisit`, `UpdateVisitStatus`, `ListVisitsRealtor`, `ListVisitsOwner`, `GetVisit` e reaproveitando middlewares `AuthMiddleware` e `PermissionMiddleware` já aplicados no group.

### 2.8 Converters (DB ↔ Domínio)
- [internal/adapter/right/mysql/visit/converters/visit_domain_to_entity.go](internal/adapter/right/mysql/visit/converters/visit_domain_to_entity.go) mapeando novos campos `owner_user_id`, `requester_user_id`, `duration_minutes`, `type`, `owner_notes`, `realtor_notes`, `rejection_reason`, `cancel_reason`, `first_owner_action_at`.
- [internal/adapter/right/mysql/visit/converters/visit_entity_to_domain.go](internal/adapter/right/mysql/visit/converters/visit_entity_to_domain.go) mantendo `sql.Null*` para opcionais.

### 2.9 Notificações
- Criar template constants em [internal/core/service/notification_templates/visit_notifications.go](internal/core/service/notification_templates/visit_notifications.go) e helpers no visit service (`sendVisitRequested`, `sendVisitApproved`, `sendVisitRejected`, `sendVisitCancelled`, `sendVisitCompleted`, `sendVisitNoShow`).

### 2.10 Factory / Wiring
- AdapterFactory: já injeta visit adapter; manter.
- ServiceFactory: criar `visitService` com dependências (global, visit repo, listing repo para owner/versão ativa, schedule repo, notification service).
- HandlerFactory: criar `visitHandler` e registrar em rotas `api/v2/visits`.

## 3. Estrutura de Diretórios (espelhamento)
```
internal/core/model/listing_model/visit_domain.go
internal/core/model/listing_model/visit_filters.go
internal/core/model/listing_model/visit_status.go
internal/core/port/left/http/visithandler/visit_handler_port.go
internal/core/port/right/repository/visit_repository/visit_repository_interface.go
internal/core/service/visit_service/
    visit_service.go
    create_visit.go
    approve_visit.go
    reject_visit.go
    cancel_visit.go
    complete_visit.go
    mark_no_show.go
    get_visit.go
    list_visits.go
internal/adapter/left/http/dto/visit/
    visit_request_dto.go
    visit_response_dto.go
    converters/visit_dto_converter.go
internal/adapter/left/http/handlers/visit_handler/
    visit_handler.go (struct + ctor)
    create_visit_handler.go
    update_status_handler.go
    list_realtor_handler.go
    list_owner_handler.go
    get_visit_handler.go
internal/adapter/left/http/routes/routes.go (atualizar grupo de visitas existente)
internal/adapter/right/mysql/visit/
    visit_adapter.go
    insert_visit.go
    update_visit.go
    get_visit_by_id.go
    list_visits.go
    pagination.go
    visit_row_mapper.go
    entities/visit_entity.go
    converters/visit_entity_to_domain.go
    converters/visit_domain_to_entity.go
internal/core/service/schedule_service/
    create_visit_entry.go
    update_visit_entry_type.go
    delete_visit_entry.go
    check_visit_conflict.go
internal/core/service/notification_templates/visit_notifications.go
```

## 4. Ordem de Execução
1. **Modelo & Schema**: atualizar domínios (visit/type/filter) e especificar DDL novo (sem aplicar migração). Depende do guia (Seção 2.1, 8).
2. **Port/Adapter Repo**: alinhar interface e adapter MySQL ao novo schema (queries, converters, pagination). Depende do passo 1.
3. **Schedule helpers**: adicionar métodos para inserir/atualizar/remover entries de visita. Depende do schema de agenda existente.
4. **Service de Visitas**: implementar métodos públicos com transações, métricas de resposta, integração com agenda e notificações. Depende dos passos 2 e 3.
5. **DTOs e Converters**: criar DTOs e conversores HTTP ↔ domínio. Depende do domínio (passo 1) e regras do serviço (passo 4).
6. **Handlers & Rotas**: criar handlers com Swagger e registrar rotas sem IDs em POST. Depende de DTOs e service (passos 4 e 5).
7. **Factory/Wiring & Notificações**: ligar dependências e adicionar templates de notificação. Depende dos passos 4 e 6.
8. **Documentação Swagger**: gerar após código (`make swagger`).

## 5. Ajustes de Dados Necessários (sem migração neste repositório)
- `listing_visits` (substituir schema atual):
    - Colunas: `id PK`, `listing_identity_id` (FK listing_identities), `listing_version TINYINT`, `requester_user_id` (FK users), `owner_user_id` (FK users), `scheduled_start DATETIME`, `scheduled_end DATETIME`, `duration_minutes SMALLINT`, `status ENUM('PENDING','APPROVED','REJECTED','CANCELLED','COMPLETED','NO_SHOW')`, `type ENUM('WITH_CLIENT','REALTOR_ONLY','CONTENT_PRODUCTION')`, `realtor_notes TEXT NULL`, `owner_notes TEXT NULL`, `rejection_reason VARCHAR(255) NULL`, `cancel_reason VARCHAR(255) NULL`, `first_owner_action_at DATETIME NULL`.
    - Índices: `idx_listing_identity_status`, `idx_requester`, `idx_owner`, `idx_scheduled_start`, `idx_status_type`.
- `listing_agenda_entries`: manter colunas existentes, garantir índice em `visit_id` e unicidade 1:1 visita↔entry (unique `visit_id`), `entry_type` deve aceitar `VISIT_PENDING`/`VISIT_CONFIRMED` já existentes.
- `listing_identities`: adicionar métricas do owner (sem aplicar migração): `owner_avg_response_time_seconds INT UNSIGNED NULL`, `owner_total_visits_responded INT UNSIGNED NOT NULL DEFAULT 0`, `owner_last_response_at DATETIME NULL`, `owner_within_sla_2h INT UNSIGNED NOT NULL DEFAULT 0` (para percentual de respostas ≤ 2h).

## 6. Observações finais
- Nenhum POST com ID em path; IDs via body conforme handlers propostos.
- Obedecer regra "uma função pública por arquivo" em services/handlers/adapters.
- Não criar/alterar testes ou swagger.json manualmente; gerar com `make swagger` após implementação.
- Notificações devem usar serviço existente; enviar apenas após commit da transação para evitar efeitos colaterais em caso de rollback.
```

### 2.4 Criar Conversor Domain → Response DTO

**Arquivo:** `internal/adapter/left/http/dto/converters/visit_domain_to_response.go`

```go
package converters

import (
    "toq_server/internal/adapter/left/http/dto"
    "toq_server/internal/core/model/listing_model"
)

// VisitDomainToResponse converte Visit domain para VisitResponse
func VisitDomainToResponse(visit *listing_model.Visit) *dto.VisitResponse {
    response := &dto.VisitResponse{
        ID:                 visit.GetID(),
        ListingIdentityID:  visit.GetListingIdentityID(),
        ListingVersion:     visit.GetListingVersion(),
        UserID:             visit.GetUserID(),
        ScheduledDate:      visit.GetScheduledDate().Format("2006-01-02"),
        ScheduledTimeStart: visit.GetScheduledTimeStart().Format("15:04"),
        ScheduledTimeEnd:   visit.GetScheduledTimeEnd().Format("15:04"),
        Status:             visit.GetStatus().String(),
        Source:             visit.GetSource(),
        Notes:              visit.GetNotes(),
        RejectionReason:    visit.GetRejectionReason(),
        CreatedAt:          visit.GetCreatedAt().Format(time.RFC3339),
        UpdatedAt:          visit.GetUpdatedAt().Format(time.RFC3339),
    }

    if firstAction := visit.GetFirstOwnerActionAt(); firstAction != nil {
        formatted := firstAction.Format(time.RFC3339)
        response.FirstOwnerActionAt = &formatted
    }

    return response
}

// VisitDomainListToResponse converte lista de Visit domain para VisitListResponse
func VisitDomainListToResponse(visits []*listing_model.Visit, totalCount, page, pageSize int) *dto.VisitListResponse {
    visitResponses := make([]dto.VisitResponse, 0, len(visits))
    for _, visit := range visits {
        visitResponses = append(visitResponses, *VisitDomainToResponse(visit))
    }

    totalPages := (totalCount + pageSize - 1) / pageSize

    return &dto.VisitListResponse{
        Visits:     visitResponses,
        TotalCount: totalCount,
        Page:       page,
        PageSize:   pageSize,
        TotalPages: totalPages,
    }
}
```

---

## Etapa 3: Service Layer

### 3.1 Criar Visit Service Interface

**Arquivo:** `internal/core/port/service/visit_service.go`

```go
package service

import (
    "context"
    "time"
    "toq_server/internal/core/model/listing_model"
)

type VisitService interface {
    // CreateVisit cria uma nova solicitação de visita
    CreateVisit(ctx context.Context, visit *listing_model.Visit, requesterID uint) (*listing_model.Visit, error)
    
    // GetVisitByID retorna uma visita por ID
    GetVisitByID(ctx context.Context, visitID uint, requesterID uint) (*listing_model.Visit, error)
    
    // ListVisits lista visitas com filtros
    ListVisits(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*listing_model.Visit, int, error)
    
    // ApproveVisit aprova uma solicitação de visita (owner apenas)
    ApproveVisit(ctx context.Context, visitID uint, ownerID uint, notes string) error
    
    // RejectVisit rejeita uma solicitação de visita (owner apenas)
    RejectVisit(ctx context.Context, visitID uint, ownerID uint, reason string) error
    
    // CancelVisit cancela uma visita (requester ou owner)
    CancelVisit(ctx context.Context, visitID uint, userID uint, reason string) error
    
    // CompleteVisit marca visita como completada (owner apenas)
    CompleteVisit(ctx context.Context, visitID uint, ownerID uint, notes string) error
    
    // MarkNoShow marca visita como não comparecimento (owner apenas)
    MarkNoShow(ctx context.Context, visitID uint, ownerID uint, notes string) error
    
    // CheckAvailability verifica disponibilidade de horário para visita
    CheckAvailability(ctx context.Context, listingIdentityID uint, date time.Time, timeStart, timeEnd time.Time) (bool, error)
    
    // UpdateOwnerResponseMetrics atualiza métricas de tempo de resposta do owner
    UpdateOwnerResponseMetrics(ctx context.Context, listingIdentityID uint, responseTimeSeconds int) error
}
```

### 3.2 Implementar Visit Service

**Arquivo:** `internal/core/service/visit_service_impl.go`

**Estrutura principal:**

```go
package service

import (
    "context"
    "errors"
    "fmt"
    "time"
    "toq_server/internal/core/model/listing_model"
    "toq_server/internal/core/model/schedule_model"
    "toq_server/internal/core/port/repository"
    "toq_server/internal/core/port/service"
)

type visitServiceImpl struct {
    globalService     service.GlobalService
    visitRepo         repository.VisitRepository
    listingRepo       repository.ListingRepository
    scheduleRepo      repository.ScheduleRepository
    notificationService service.NotificationService
}

func NewVisitService(
    globalService service.GlobalService,
    visitRepo repository.VisitRepository,
    listingRepo repository.ListingRepository,
    scheduleRepo repository.ScheduleRepository,
    notificationService service.NotificationService,
) service.VisitService {
    return &visitServiceImpl{
        globalService:       globalService,
        visitRepo:           visitRepo,
        listingRepo:         listingRepo,
        scheduleRepo:        scheduleRepo,
        notificationService: notificationService,
    }
}
```

**Métodos principais a implementar:**

1. **CreateVisit:**
   - Validar que listing_identity existe e tem versão ativa
   - Obter active_version_id do listing_identity
   - Validar disponibilidade do horário (CheckAvailability)
   - Iniciar transação
   - Inserir visit com status PENDING (listing_identity_id + listing_version)
   - Criar agenda_entry tipo VISIT_PENDING, blocking=true
   - Enviar notificação para owner
   - Commit transação

2. **ApproveVisit:**
   - Validar permissão (user é owner do listing_identity)
   - Validar status atual (deve ser PENDING)
   - Iniciar transação
   - Atualizar visit: status=APPROVED, firstOwnerActionAt=now
   - Atualizar agenda_entry: tipo=VISIT_CONFIRMED
    - Calcular tempo de resposta usando o timestamp da criação em memória (now - request_time) em segundos
   - Atualizar métricas em listing_identity (avg_response_time, total_visits_responded, last_response_at)
   - Enviar notificação para requester
   - Commit transação

3. **RejectVisit:**
   - Validar permissão (user é owner do listing_identity)
   - Validar status (deve ser PENDING)
   - Iniciar transação
   - Atualizar visit: status=REJECTED, rejectionReason, firstOwnerActionAt=now
   - Remover agenda_entry
    - Calcular tempo de resposta usando o timestamp da criação em memória (now - request_time) em segundos
   - Atualizar métricas em listing_identity (avg_response_time, total_visits_responded, last_response_at)
   - Enviar notificação para requester
   - Commit transação

4. **CancelVisit:**
   - Validar permissão (requester ou owner)
   - Validar status (não pode estar COMPLETED/NO_SHOW/CANCELLED)
   - Iniciar transação
   - Atualizar visit: status=CANCELLED
   - Remover agenda_entry
   - Enviar notificação
   - Commit transação

5. **CompleteVisit / MarkNoShow:**
   - Validar permissão (owner)
   - Validar status (deve ser APPROVED)
   - Validar data (deve ser após a data agendada)
   - Atualizar visit status
   - Remover agenda_entry
   - Enviar notificação

6. **CheckAvailability:**
   - Buscar agenda do listing
   - Listar entries no período (date + time range)
   - Verificar conflitos com entries blocking=true
   - Retornar disponibilidade

### 3.3 Helpers para Schedule Integration

**Adicionar métodos no schedule_service ou criar helpers:**

```go
// CreateVisitPendingEntry cria entry de visita pendente
func (s *scheduleService) CreateVisitPendingEntry(ctx context.Context, agendaID, visitID uint, date time.Time, start, end time.Time) error

// UpdateVisitEntryToConfirmed atualiza entry para visita confirmada
func (s *scheduleService) UpdateVisitEntryToConfirmed(ctx context.Context, visitID uint) error

// RemoveVisitEntry remove entry de visita
func (s *scheduleService) RemoveVisitEntry(ctx context.Context, visitID uint) error

// CheckVisitConflict verifica conflitos de horário
func (s *scheduleService) CheckVisitConflict(ctx context.Context, agendaID uint, date time.Time, start, end time.Time, excludeVisitID *uint) (bool, error)
```

---

## Etapa 4: Handlers e Rotas

### 4.1 Criar Visit Handler

**Arquivo:** `internal/adapter/left/http/handlers/visit_handler.go`

**Endpoints a implementar:**

```go
// POST /v1/listing-identities/{listing_identity_id}/visits - Criar solicitação de visita
// @Summary Solicitar visita
// @Tags Visits
// @Accept json
// @Produce json
// @Param listing_identity_id path int true "Listing Identity ID"
// @Param request body dto.CreateVisitRequest true "Visit request"
// @Success 201 {object} dto.VisitResponse
// @Router /v1/listing-identities/{listing_identity_id}/visits [post]
func (h *visitHandler) CreateVisit(c *gin.Context)

// GET /v1/visits/{id} - Obter visita por ID
// @Summary Obter detalhes da visita
// @Tags Visits
// @Produce json
// @Param id path int true "Visit ID"
// @Success 200 {object} dto.VisitResponse
// @Router /v1/visits/{id} [get]
func (h *visitHandler) GetVisit(c *gin.Context)

// GET /v1/visits - Listar visitas (com filtros)
// @Summary Listar visitas
// @Tags Visits
// @Produce json
// @Param listing_id query int false "Filter by listing"
// @Param user_id query int false "Filter by user"
// @Param status query string false "Filter by status"
// @Param date_from query string false "Filter from date (YYYY-MM-DD)"
// @Param date_to query string false "Filter to date (YYYY-MM-DD)"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} dto.VisitListResponse
// @Router /v1/visits [get]
func (h *visitHandler) ListVisits(c *gin.Context)

// PUT /v1/visits/{id}/approve - Aprovar visita (owner)
// @Summary Aprovar solicitação de visita
// @Tags Visits
// @Accept json
// @Produce json
// @Param id path int true "Visit ID"
// @Param request body dto.UpdateVisitStatusRequest true "Approval data"
// @Success 200 {object} dto.VisitResponse
// @Router /v1/visits/{id}/approve [put]
func (h *visitHandler) ApproveVisit(c *gin.Context)

// PUT /v1/visits/{id}/reject - Rejeitar visita (owner)
// @Summary Rejeitar solicitação de visita
// @Tags Visits
// @Param id path int true "Visit ID"
// @Param request body dto.UpdateVisitStatusRequest true "Rejection data"
// @Success 200 {object} dto.VisitResponse
// @Router /v1/visits/{id}/reject [put]
func (h *visitHandler) RejectVisit(c *gin.Context)

// PUT /v1/visits/{id}/cancel - Cancelar visita
// @Summary Cancelar visita
// @Tags Visits
// @Param id path int true "Visit ID"
// @Param request body dto.UpdateVisitStatusRequest true "Cancellation data"
// @Success 200 {object} dto.VisitResponse
// @Router /v1/visits/{id}/cancel [put]
func (h *visitHandler) CancelVisit(c *gin.Context)

// PUT /v1/visits/{id}/complete - Marcar visita como completada (owner)
// @Summary Completar visita
// @Tags Visits
// @Param id path int true "Visit ID"
// @Param request body dto.UpdateVisitStatusRequest true "Completion data"
// @Success 200 {object} dto.VisitResponse
// @Router /v1/visits/{id}/complete [put]
func (h *visitHandler) CompleteVisit(c *gin.Context)

// PUT /v1/visits/{id}/no-show - Marcar não comparecimento (owner)
// @Summary Marcar não comparecimento
// @Tags Visits
// @Param id path int true "Visit ID"
// @Param request body dto.UpdateVisitStatusRequest true "No-show data"
// @Success 200 {object} dto.VisitResponse
// @Router /v1/visits/{id}/no-show [put]
func (h *visitHandler) MarkNoShow(c *gin.Context)
```

**Padrão de implementação:**

```go
func (h *visitHandler) CreateVisit(c *gin.Context) {
    ctx, span := utils.GenerateTracer(c, "VisitHandler.CreateVisit")
    defer span.End()

    // 1. Extrair userID do contexto (middleware de auth)
    userID := c.GetUint("user_id")
    
    // 2. Parse listingIdentityID da URL
    listingIdentityID, err := strconv.ParseUint(c.Param("listing_identity_id"), 10, 32)
    if err != nil {
        c.JSON(400, gin.H{"error": "invalid listing_identity_id"})
        return
    }

    // 3. Bind request body
    var req dto.CreateVisitRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    req.ListingIdentityID = uint(listingIdentityID)

    // 4. Converter DTO para Domain
    visit, err := converters.CreateVisitRequestToDomain(&req, userID)
    if err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // 5. Chamar service
    createdVisit, err := h.visitService.CreateVisit(ctx, visit, userID)
    if err != nil {
        // Tratar erros específicos
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    // 6. Converter para response DTO
    response := converters.VisitDomainToResponse(createdVisit)

    c.JSON(201, response)
}
```

### 4.2 Registrar Rotas

**Arquivo:** `internal/adapter/left/http/routes/visit_routes.go`

```go
package routes

import (
    "github.com/gin-gonic/gin"
    "toq_server/internal/adapter/left/http/handlers"
    "toq_server/internal/adapter/left/http/middleware"
)

func RegisterVisitRoutes(router *gin.RouterGroup, visitHandler *handlers.VisitHandler, authMiddleware *middleware.AuthMiddleware) {
    visits := router.Group("/visits")
    visits.Use(authMiddleware.Authenticate())
    {
        visits.POST("", visitHandler.CreateVisit)  // POST /v1/visits (ou via listings)
        visits.GET("", visitHandler.ListVisits)
        visits.GET("/:id", visitHandler.GetVisit)
        visits.PUT("/:id/approve", visitHandler.ApproveVisit)
        visits.PUT("/:id/reject", visitHandler.RejectVisit)
        visits.PUT("/:id/cancel", visitHandler.CancelVisit)
        visits.PUT("/:id/complete", visitHandler.CompleteVisit)
        visits.PUT("/:id/no-show", visitHandler.MarkNoShow)
    }

    // Rota alternativa via listing-identities
    listingIdentities := router.Group("/listing-identities")
    listingIdentities.Use(authMiddleware.Authenticate())
    {
        listingIdentities.POST("/:listing_identity_id/visits", visitHandler.CreateVisit)
    }
}
```

**Atualizar:** `internal/adapter/left/http/routes/routes.go`

```go
// Adicionar no InitializeRoutes:
RegisterVisitRoutes(v1, visitHandler, authMiddleware)
```

### 4.3 Wiring no Factory

**Arquivo:** `internal/adapter/left/concrete_adapter_factory.go`

**Adicionar:**

```go
// No método setupServices:
visitService := service.NewVisitService(
    globalService,
    f.adapterRepository.Visit,
    f.adapterRepository.Listing,
    f.adapterRepository.Schedule,
    notificationService,
)

// No método setupHandlers:
visitHandler := handlers.NewVisitHandler(visitService)

// Passar visitHandler para InitializeRoutes
```

---

## Etapa 5: Notificações

### 5.1 Adicionar Templates de Notificação

**Arquivo:** `internal/core/service/notification_templates.go` (ou criar novo)

**Templates a adicionar:**

```go
const (
    // Visit request criado (para owner)
    NotificationTypeVisitRequested = "VISIT_REQUESTED"
    NotificationTitleVisitRequested = "Nova Solicitação de Visita"
    NotificationBodyVisitRequested = "{{.RequesterName}} solicitou uma visita para {{.ListingTitle}} em {{.ScheduledDate}} às {{.ScheduledTime}}"

    // Visit aprovado (para requester)
    NotificationTypeVisitApproved = "VISIT_APPROVED"
    NotificationTitleVisitApproved = "Visita Aprovada"
    NotificationBodyVisitApproved = "Sua visita para {{.ListingTitle}} em {{.ScheduledDate}} às {{.ScheduledTime}} foi aprovada"

    // Visit rejeitado (para requester)
    NotificationTypeVisitRejected = "VISIT_REJECTED"
    NotificationTitleVisitRejected = "Visita Rejeitada"
    NotificationBodyVisitRejected = "Sua visita para {{.ListingTitle}} foi rejeitada: {{.Reason}}"

    // Visit cancelado (para ambos)
    NotificationTypeVisitCancelled = "VISIT_CANCELLED"
    NotificationTitleVisitCancelled = "Visita Cancelada"
    NotificationBodyVisitCancelled = "A visita para {{.ListingTitle}} em {{.ScheduledDate}} foi cancelada"

    // Lembrete de visita (para ambos, 24h antes)
    NotificationTypeVisitReminder = "VISIT_REMINDER"
    NotificationTitleVisitReminder = "Lembrete de Visita"
    NotificationBodyVisitReminder = "Você tem uma visita agendada para {{.ListingTitle}} amanhã às {{.ScheduledTime}}"
)
```

### 5.2 Implementar Envio de Notificações no Service

**No visit_service_impl.go, adicionar chamadas:**

```go
// Após criar visit (CreateVisit):
h.sendVisitRequestedNotification(ctx, visit, listing, requester)

// Após aprovar (ApproveVisit):
h.sendVisitApprovedNotification(ctx, visit, listing, requester)

// Após rejeitar (RejectVisit):
h.sendVisitRejectedNotification(ctx, visit, listing, requester)

// Após cancelar (CancelVisit):
h.sendVisitCancelledNotification(ctx, visit, listing, requester, owner)
```

**Métodos helper:**

```go
func (s *visitServiceImpl) sendVisitRequestedNotification(ctx context.Context, visit *listing_model.Visit, listing *listing_model.Listing, requester *user_model.User) {
    data := map[string]interface{}{
        "RequesterName":  requester.GetName(),
        "ListingTitle":   listing.GetTitle(),
        "ScheduledDate":  visit.GetScheduledDate().Format("02/01/2006"),
        "ScheduledTime":  visit.GetScheduledTimeStart().Format("15:04"),
    }

    s.notificationService.SendNotification(ctx, &notification_model.Notification{
        UserID:   listing.GetOwnerID(), // Para o owner
        Type:     NotificationTypeVisitRequested,
        Title:    NotificationTitleVisitRequested,
        Body:     s.renderTemplate(NotificationBodyVisitRequested, data),
        Data:     data,
    })
}

// Implementar métodos similares para outros eventos
```

---

## Etapa 6: Schedule Helpers e Validações

### 6.1 Helpers no Schedule Service

**Arquivo:** `internal/core/service/schedule_service_impl.go`

**Adicionar métodos:**

```go
// CreateVisitPendingEntry cria entry na agenda para visita pendente
func (s *scheduleServiceImpl) CreateVisitPendingEntry(ctx context.Context, listingIdentityID, visitID uint, date time.Time, start, end time.Time) error {
    // 1. Obter ou criar agenda do listing_identity
    agenda, err := s.GetOrCreateAgenda(ctx, listingIdentityID)
    if err != nil {
        return err
    }

    // 2. Criar entry tipo VISIT_PENDING, blocking=true
    entry := schedule_model.NewAgendaEntry(
        0,
        agenda.GetID(),
        date,
        start,
        end,
        schedule_model.EntryTypeVisitPending,
        true, // blocking
        "",   // reason
        &visitID,
        nil, // photo_booking_id
    )

    // 3. Inserir entry
    return s.scheduleRepo.InsertAgendaEntry(ctx, entry)
}

// UpdateVisitEntryToConfirmed atualiza tipo da entry para confirmada
func (s *scheduleServiceImpl) UpdateVisitEntryToConfirmed(ctx context.Context, visitID uint) error {
    // 1. Buscar entry pelo visit_id
    entry, err := s.scheduleRepo.GetAgendaEntryByVisitID(ctx, visitID)
    if err != nil {
        return err
    }

    // 2. Atualizar tipo para VISIT_CONFIRMED
    entry.SetType(schedule_model.EntryTypeVisitConfirmed)
    
    return s.scheduleRepo.UpdateAgendaEntry(ctx, entry)
}

// RemoveVisitEntry remove entry da agenda
func (s *scheduleServiceImpl) RemoveVisitEntry(ctx context.Context, visitID uint) error {
    entry, err := s.scheduleRepo.GetAgendaEntryByVisitID(ctx, visitID)
    if err != nil {
        return err
    }

    return s.scheduleRepo.DeleteAgendaEntry(ctx, entry.GetID())
}

// CheckVisitConflict verifica se há conflito de horário
func (s *scheduleServiceImpl) CheckVisitConflict(ctx context.Context, listingIdentityID uint, date time.Time, start, end time.Time, excludeVisitID *uint) (bool, error) {
    // 1. Obter agenda
    agenda, err := s.GetOrCreateAgenda(ctx, listingIdentityID)
    if err != nil {
        return false, err
    }

    // 2. Listar entries no dia
    entries, err := s.scheduleRepo.ListAgendaEntriesByDate(ctx, agenda.GetID(), date)
    if err != nil {
        return false, err
    }

    // 3. Verificar sobreposição com entries blocking=true
    for _, entry := range entries {
        // Ignorar entry da própria visita se fornecido
        if excludeVisitID != nil && entry.GetVisitID() != nil && *entry.GetVisitID() == *excludeVisitID {
            continue
        }

        if !entry.IsBlocking() {
            continue
        }

        // Verificar overlap
        if timesOverlap(start, end, entry.GetStartTime(), entry.GetEndTime()) {
            return true, nil // Há conflito
        }
    }

    return false, nil // Sem conflito
}

func timesOverlap(start1, end1, start2, end2 time.Time) bool {
    return start1.Before(end2) && end1.After(start2)
}

// UpdateOwnerResponseMetrics atualiza métricas de tempo de resposta
func (s *visitServiceImpl) UpdateOwnerResponseMetrics(ctx context.Context, listingIdentityID uint, responseTimeSeconds int) error {
    // 1. Buscar listing_identity
    identity, err := s.listingIdentityRepo.GetByID(ctx, listingIdentityID)
    if err != nil {
        return err
    }

    // 2. Calcular nova média usando EMA (Exponential Moving Average)
    // Peso de 20% para nova medição, 80% para histórico
    currentAvg := identity.GetOwnerAvgResponseTime()
    var newAvg int
    if currentAvg == 0 {
        // Primeira medição
        newAvg = responseTimeSeconds
    } else {
        // Média móvel exponencial
        newAvg = int(float64(currentAvg)*0.8 + float64(responseTimeSeconds)*0.2)
    }

    // 3. Atualizar campos
    identity.SetOwnerAvgResponseTime(newAvg)
    identity.SetOwnerTotalVisitsResponded(identity.GetOwnerTotalVisitsResponded() + 1)
    identity.SetOwnerLastResponseAt(time.Now())

    // 4. Persistir
    return s.listingIdentityRepo.Update(ctx, identity)
}
```

### 6.2 Adicionar Query no Schedule Repository

**Arquivo:** `internal/adapter/right/mysql/schedule/get_agenda_entry_by_visit_id.go`

```go
package schedule

import (
    "context"
    "toq_server/internal/core/model/schedule_model"
)

func (r *scheduleRepository) GetAgendaEntryByVisitID(ctx context.Context, visitID uint) (*schedule_model.AgendaEntry, error) {
    query := `
        SELECT id, agenda_id, entry_date, start_time, end_time, type, blocking, reason, visit_id, photo_booking_id
        FROM listing_agenda_entries
        WHERE visit_id = ?
        LIMIT 1
    `

    row := r.db.QueryRowContext(ctx, query, visitID)

    entry, err := mapRowToAgendaEntry(row)
    if err != nil {
        return nil, err
    }

    return entry, nil
}
```

---

## Etapa 7: Validações e Regras de Negócio

### 7.1 Validações a implementar no Service

**No visit_service_impl.go:**

```go
// validateVisitPermissions valida se user pode executar ação
func (s *visitServiceImpl) validateVisitPermissions(ctx context.Context, visitID, userID uint, requireOwner bool) (*listing_model.Visit, *listing_model.Listing, error) {
    // 1. Buscar visit
    visit, err := s.visitRepo.GetVisitByID(ctx, visitID)
    if err != nil {
        return nil, nil, err
    }

    // 2. Buscar listing_identity para obter owner
    listingIdentity, err := s.listingIdentityRepo.GetByID(ctx, visit.GetListingIdentityID())
    if err != nil {
        return nil, nil, err
    }

    // 3. Verificar permissão
    isOwner := listingIdentity.GetUserID() == userID
    isRequester := visit.GetUserID() == userID

    if requireOwner && !isOwner {
        return nil, nil, errors.New("only listing owner can perform this action")
    }

    if !requireOwner && !isOwner && !isRequester {
        return nil, nil, errors.New("user not authorized for this visit")
    }

    return visit, listing, nil
}

// validateVisitStatus valida se status permite transição
func validateVisitStatus(current listing_model.VisitStatus, action string) error {
    switch action {
    case "approve":
        if current != listing_model.VisitStatusPending {
            return errors.New("only pending visits can be approved")
        }
    case "reject":
        if current != listing_model.VisitStatusPending {
            return errors.New("only pending visits can be rejected")
        }
    case "cancel":
        if current == listing_model.VisitStatusCompleted || 
           current == listing_model.VisitStatusNoShow ||
           current == listing_model.VisitStatusCancelled {
            return errors.New("cannot cancel visit in current status")
        }
    case "complete":
        if current != listing_model.VisitStatusApproved {
            return errors.New("only approved visits can be completed")
        }
    case "no_show":
        if current != listing_model.VisitStatusApproved {
            return errors.New("only approved visits can be marked as no-show")
        }
    }
    return nil
}

// validateVisitDateTime valida data/hora da visita
func (s *visitServiceImpl) validateVisitDateTime(ctx context.Context, date, start, end time.Time) error {
    now := time.Now()

    // 1. Data não pode ser no passado
    if date.Before(now.Truncate(24 * time.Hour)) {
        return errors.New("visit date cannot be in the past")
    }

    // 2. Data não pode ser muito distante (ex: max 90 dias)
    maxFutureDate := now.AddDate(0, 0, 90)
    if date.After(maxFutureDate) {
        return errors.New("visit date cannot be more than 90 days in the future")
    }

    // 3. Horário de início não pode ser no passado (se data é hoje)
    if date.Truncate(24*time.Hour).Equal(now.Truncate(24*time.Hour)) {
        visitDateTime := time.Date(date.Year(), date.Month(), date.Day(), start.Hour(), start.Minute(), 0, 0, date.Location())
        if visitDateTime.Before(now) {
            return errors.New("visit time cannot be in the past")
        }
    }

    // 4. Duração mínima (ex: 30 min)
    duration := end.Sub(start)
    if duration < 30*time.Minute {
        return errors.New("visit duration must be at least 30 minutes")
    }

    // 5. Duração máxima (ex: 4 horas)
    if duration > 4*time.Hour {
        return errors.New("visit duration cannot exceed 4 hours")
    }

    return nil
}
```

---

## Etapa 8: Testes

### 8.1 Testes Unitários do Service

**Arquivo:** `internal/core/service/visit_service_test.go`

**Testes a implementar:**

```go
func TestCreateVisit_Success(t *testing.T)
func TestCreateVisit_ListingNotFound(t *testing.T)
func TestCreateVisit_TimeConflict(t *testing.T)
func TestCreateVisit_InvalidDateTime(t *testing.T)

func TestApproveVisit_Success(t *testing.T)
func TestApproveVisit_NotOwner(t *testing.T)
func TestApproveVisit_InvalidStatus(t *testing.T)

func TestRejectVisit_Success(t *testing.T)
func TestRejectVisit_MissingReason(t *testing.T)

func TestCancelVisit_ByRequester(t *testing.T)
func TestCancelVisit_ByOwner(t *testing.T)
func TestCancelVisit_InvalidStatus(t *testing.T)

func TestCompleteVisit_Success(t *testing.T)
func TestCompleteVisit_BeforeScheduledDate(t *testing.T)

func TestCheckAvailability_NoConflict(t *testing.T)
func TestCheckAvailability_WithConflict(t *testing.T)
```

### 8.2 Testes de Integração

**Arquivo:** `internal/adapter/left/http/handlers/visit_handler_test.go`

**Testes a implementar:**

```go
func TestCreateVisitEndpoint(t *testing.T)
func TestGetVisitEndpoint(t *testing.T)
func TestListVisitsEndpoint(t *testing.T)
func TestApproveVisitEndpoint(t *testing.T)
func TestRejectVisitEndpoint(t *testing.T)
// etc...
```

---

## Etapa 9: Documentação e Deploy

### 9.1 Atualizar Swagger

```bash
# Gerar nova documentação swagger
swag init -g cmd/toq_server.go
```

### 9.2 Executar Migração do Banco

```bash
# Aplicar migração
mysql -u root -p toq_db < scripts/migrations/add_visit_fields.sql
```

### 9.3 Testar Endpoints

**Postman/Insomnia collection com exemplos:**

```json
POST /v1/listings/123/visits
{
  "scheduled_date": "2025-01-15",
  "scheduled_time_start": "14:00",
  "scheduled_time_end": "15:00",
  "notes": "Gostaria de visitar o imóvel",
  "source": "APP"
}

PUT /v1/visits/456/approve
{
  "status": "APPROVED",
  "notes": "Aprovado"
}

GET /v1/visits?listing_id=123&status=PENDING&page=1&page_size=10
```

---

## Checklist Final

- [ ] Etapa 1: Schema atualizado (listing_visits + listing_identities métricas) e migração aplicada
- [ ] Etapa 2: DTOs e conversores criados
- [ ] Etapa 3: Service implementado com todas as operações
- [ ] Etapa 4: Handlers e rotas registradas
- [ ] Etapa 5: Templates de notificação adicionados
- [ ] Etapa 6: Helpers de schedule implementados
- [ ] Etapa 7: Validações de negócio implementadas
- [ ] Etapa 8: Testes unitários e integração criados
- [ ] Etapa 9: Swagger atualizado e migração aplicada
- [ ] Testes manuais realizados
- [ ] Code review e deploy

---

## Ordem Recomendada de Implementação

1. **Etapa 1** (Schema) - Executar migração primeiro
2. **Etapa 2** (DTOs) - Criar DTOs e conversores
3. **Etapa 6** (Schedule helpers) - Preparar helpers antes do service
4. **Etapa 3** (Service) - Implementar lógica de negócio
5. **Etapa 5** (Notificações) - Adicionar templates
6. **Etapa 7** (Validações) - Adicionar regras de negócio
7. **Etapa 4** (Handlers) - Criar endpoints HTTP
8. **Etapa 8** (Testes) - Criar testes
9. **Etapa 9** (Deploy) - Atualizar docs e deploy

---

## Notas Importantes

- **Arquitetura de versionamento:** Visitas são vinculadas a `listing_identity_id` (identidade do anúncio), não a versões específicas. O campo `listing_version` captura qual versão estava ativa no momento da solicitação.
- **Métricas de resposta:** Calculadas incrementalmente usando EMA (média móvel exponencial) para performance. Atualizar apenas em ApproveVisit/RejectVisit.
- **SLA firstOwnerActionAt:** Setar apenas na primeira ação do owner (approve/reject)
- **Transaction management:** Usar globalService.BeginTransaction() para operações multi-step
- **Agenda blocking:** VISIT_PENDING e VISIT_CONFIRMED devem sempre ter blocking=true
- **Cascade delete:** FK em agenda_entries garante remoção automática ao deletar visit
- **Unique constraint:** visit_id em agenda_entries garante 1:1 entre visit e entry
- **Time validation:** Validar que horários não sejam no passado e respeitem limites de duração
- **Permission checks:** Owner pode aprovar/rejeitar/completar; Requester pode cancelar próprias visitas
- **Notification timing:** Enviar notificações após commit bem-sucedido da transação
