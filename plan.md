# Plano de Auditoria — TOQ Server

## Contexto
- Auditoria atual grava apenas `table_name`, `action`, `executed_by`, `executed_at`, sem `table_id` válido, sem metadados e com nomes de tabela divergentes do schema real. Isso impede reconstruir o ciclo de vida de anúncios e outras operações críticas.
- Objetivo: criar trilha completa, consultável e correlacionada com tracing/logs, cobrindo ciclo de vida de anúncios, propostas, visitas, mídia e mudanças de estado.

## Objetivos
- Padronizar modelo de auditoria (evento estruturado) com ator, alvo, operação, timestamps e correlação (`request_id`/`trace_id`).
- Garantir rastreabilidade do ciclo de vida de anúncio (criação, edição, promoção, publicação/despublicação, descarte, status changes, propostas, visitas, mídia, agenda).
- Possibilitar consultas por alvo (listing_identity/listing_version) e linha do tempo.

## Escopo
- Criar novo modelo/porta/serviço/adapter de auditoria (infraestrutura). 
- Instrumentar serviços críticos com chamadas `auditService.RecordChange` contendo metadados mínimos para reconstrução.
- Não incluir migrações no repositório; alinhamento com DBA é obrigatório para criação/alteração de tabela(s) de auditoria.

## Fora de Escopo
- Remover colunas `created_by/updated_by` existentes (apenas avaliar depois da migração completa).
- Ajustar dashboards/relatórios legados neste momento (apenas registrar necessidades).

## Faseamento e Dependências
1) **Infraestrutura de Auditoria** (pré-requisito para demais fases)
   - Modelagem: `AuditEvent`, `AuditActor`, `AuditTarget`, `AuditOperation`, `AuditCorrelation`, `RecordInput`, filtros/paginação.
   - Porta `audit_repository` + adapter MySQL (`audit_adapter.go`, `create_event.go`, entities/converters) usando InstrumentedAdapter.
   - Serviço `audit_service` com `RecordChange` e helper `WithRequestContext` para enriquecer com `request_id`/`trace_id`/device/ip/UA.
   - Ajuste de DI (factory) para expor `AuditService`; manter compatibilidade temporária com `GlobalService.CreateAudit` apenas como fachada (opcional) até refatorar call sites.
   - **Dependência**: aprovação de schema com DBA (nova tabela `audit_events` ou evolução da atual). Sem merge de migração aqui.

2) **Cobertura de Fluxos Prioritários** (depende da fase 1 concluída)
   - Ciclo de anúncio: create/update/promote/status_change/publish/unpublish/discard/delete + snapshots de preço/status/versão.
   - Propostas: create/accept/reject/cancel + flags no listing.
   - Visitas/agendas: schedule/approve/reject/cancel/complete/no-show + creation/finish agenda.
   - Mídia: approvals/rejections de batch e assets.
   - Autenticação/sessões: signin/signout/refresh/reset password (mínimo registro de ator e resultado).
   - Permissões/admin: criação/alteração de usuários do sistema, roles e vínculos (owner/realtor/agency).
   - Cada chamada deve preencher: `Target` (tipo + IDs reais do schema), `Operation` enum, `Actor` (user_id, role_slug, device_id, ip, user_agent), `OccurredAt` UTC, `Correlation` (request_id, trace_id) e `Metadata` com deltas relevantes (ex.: status_from/to, version_from/to, price_from/to, proposal_id, visit_id, agenda_id).

3) **Rollout e Observabilidade** (depende da fase 2 parcial para testes)
   - Consultas internas: endpoints/queries de suporte para timeline por `target_type + target_id (+version)` paginada por `occurred_at`.
   - Dashboards/logging: validar correlação com Tempo/Loki via `trace_id`/`request_id` armazenados.
   - Playbook de uso: como registrar operações, campos obrigatórios e exemplos.

## Detalhes de Implementação (guias por componente)
- **Tabela sugerida (alinhamento com DBA)**: `audit_events( id PK, occurred_at, actor_id, actor_role, actor_device_id, actor_ip, actor_user_agent, target_type, target_id, target_version NULL, operation, metadata JSON, request_id, trace_id )` com índices por `(target_type, target_id, occurred_at)` e `(request_id)`. Sem migrations no repo.
- **Enums/nomes de recurso**: usar nomes do schema real (`listing_identities`, `listing_versions`, `listing_visits`, `proposals`, `media_assets`, `listing_agendas`).
- **Metadata mínima por fluxo**:
  - Listing: `{version_from, version_to, status_from, status_to, price_from, price_to}` quando aplicável.
  - Proposta: `{proposal_id, status_from, status_to}`.
  - Visita/agenda: `{visit_id|agenda_id, status_from, status_to}`.
  - Mídia: `{asset_id or batch_id, status_from, status_to}`.
  - Auth: `{session_id?, result}` (sem dados sensíveis).
- **Segurança**: evitar PII no metadata; se necessário, mascarar. Guardar somente IDs e estados.
- **Tracing/logs**: não criar spans em handler; `RecordChange` inicia tracer leve e propaga logger; marcar `SetSpanError` em falhas de infra.

## Critérios de Aceite por Fase
- Fase 1: serviço/porta/adapter compilam; tracer/logger/correlação presentes; schema validado com DBA; testes manuais de inserção funcionam (table_id/target_id preenchido).
- Fase 2: todos os pontos de negócio priorizados chamam `RecordChange` com IDs e metadata corretos; nenhuma chamada deixa `actor_id` ou `target_id` vazio; enums correspondem ao schema real; preço/status/version mudando gera delta no metadata.
- Fase 3: é possível consultar timeline completa de um listing (identity e versões) e ver proposta/visita/mídia associadas, com request_id/trace_id; playbook/documentação publicada.

## Riscos e Mitigações
- Divergência de schema (nomes/índices): mitigar com validação antecipada com DBA antes de codar adapters.
- Metadata sensível: definir whitelist de campos e mascarar PII.
- Falta de cobertura em fluxos legados: listar endpoints/serviços sem auditoria e tratá-los como débito até completar.
- Performance: usar índices propostos; metadata enxuta; spans leves.

## Itens Ações (sequência sugerida)
1. Validar com DBA o schema final de `audit_events` (índices, tipos, JSON vs TEXT) e políticas de retenção.
2. Implementar domínio/porta/adapter/serviço de auditoria (Fase 1) + DI/factory.
3. Mapear todos os fluxos alvo e produzir checklist de cobertura por domínio (listing, proposal, visit, media, auth, admin/roles).
4. Instrumentar ciclo de anúncio completo e propostas (Fase 2 priorização alta).
5. Instrumentar visitas/agendas e mídia (Fase 2 médio).
6. Instrumentar auth/sessões e admin/roles (Fase 2 baixo).
7. Criar endpoints/queries internas de timeline e validar correlação com tracing/logs (Fase 3).
8. Documentar playbook e atualizar guia interno com exemplos de uso e campos obrigatórios.

## Blueprint Detalhado para Codificação (sem aplicar ainda)
> Objetivo: descrever todo o código a ser criado/alterado para antecipar riscos. Não executar migração nem merge; usar como guia único.

### 0. Schema (para alinhamento com DBA; não comitar)
```sql
CREATE TABLE audit_events (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  occurred_at DATETIME(6) NOT NULL,
  actor_id BIGINT UNSIGNED NOT NULL,
  actor_role VARCHAR(64) NULL,
  actor_device_id VARCHAR(100) NULL,
  actor_ip VARCHAR(64) NULL,
  actor_user_agent VARCHAR(255) NULL,
  target_type VARCHAR(64) NOT NULL,
  target_id BIGINT UNSIGNED NOT NULL,
  target_version INT UNSIGNED NULL,
  operation VARCHAR(64) NOT NULL,
  metadata JSON NULL,
  request_id VARCHAR(64) NULL,
  trace_id VARCHAR(64) NULL,
  INDEX idx_target_time (target_type, target_id, occurred_at),
  INDEX idx_request (request_id),
  INDEX idx_trace (trace_id)
);
```

### 1. Domínio (novo) — `internal/core/model/audit_model`
- `audit_event.go`
```go
type AuditOperation string
// consts: OperationCreate, OperationUpdate, OperationPromote, OperationStatusChange,
// OperationPublish, OperationUnpublish, OperationDiscard, OperationDelete,
// OperationProposalCreate/Accept/Reject/Cancel, OperationVisitRequest/Approve/Reject/Cancel/Complete/NoShow,
// OperationMediaApprove/Reject, OperationAgendaCreate/Finish, OperationAuthSignin/Signout/PasswordReset.

type TargetType string
// consts: TargetListingIdentity, TargetListingVersion, TargetListingVisit,
// TargetProposal, TargetMediaAsset, TargetListingAgenda, TargetSession.

type AuditActor struct {
   ID int64
   RoleSlug string
   DeviceID string
   IP string
   UserAgent string
}

type AuditTarget struct {
   Type TargetType
   ID int64
   Version *int64 // optional
}

type AuditCorrelation struct {
   RequestID string
   TraceID string
}

type AuditEvent interface {
   ID() int64; SetID(int64)
   OccurredAt() time.Time; SetOccurredAt(time.Time)
   Actor() AuditActor; SetActor(AuditActor)
   Target() AuditTarget; SetTarget(AuditTarget)
   Operation() AuditOperation; SetOperation(AuditOperation)
   Metadata() map[string]any; SetMetadata(map[string]any)
   Correlation() AuditCorrelation; SetCorrelation(AuditCorrelation)
}

type RecordInput struct {
   Actor AuditActor
   Target AuditTarget
   Operation AuditOperation
   Metadata map[string]any
   OccurredAt time.Time // optional; default now UTC
   Correlation AuditCorrelation // optional
}
```

- `audit_event_domain.go`: struct privado implementando a interface com getters/setters.
- `errors.go` (opcional): validações de operação/target.

### 2. Porta (novo) — `internal/core/port/right/repository/audit_repository/audit_repository_interface.go`
```go
type AuditRepoPort interface {
   CreateEvent(ctx context.Context, tx *sql.Tx, event auditmodel.AuditEvent) error
   ListEventsByTarget(ctx context.Context, tx *sql.Tx, filter auditmodel.TargetFilter, page auditmodel.Page) ([]auditmodel.AuditEvent, error)
}
```
`TargetFilter` inclui `TargetType`, `TargetID`, `Version *int64`, intervalo de datas opcional.

### 3. Serviço (novo) — `internal/core/service/audit_service`
- `audit_service.go`: struct + interface + `NewAuditService(repo AuditRepoPort, globalSvc GlobalServiceInterface)`.
- `record_change.go` (método público único):
```go
func (s *auditService) RecordChange(ctx context.Context, tx *sql.Tx, input auditmodel.RecordInput) (err error) {
   ctx, end, tracerErr := utils.GenerateTracer(ctx); if tracerErr != nil { return derrors.Infra("audit tracer", tracerErr) }
   defer end(); ctx = utils.ContextWithLogger(ctx); logger := utils.LoggerFromContext(ctx)

   if err = s.validateInput(input); err != nil { return derrors.BadRequest("invalid audit input", map[string]string{"reason": err.Error()}) }

   event := s.builder.FromRecordInput(ctx, input) // preenche occurred_at UTC, correlation se ausente

   if err = s.repo.CreateEvent(ctx, tx, event); err != nil {
      utils.SetSpanError(ctx, err); logger.Error("audit.record.persist_error", "err", err, "target", input.Target.Type, "operation", input.Operation)
      return derrors.Infra("persist audit event", err)
   }
   return nil
}
```
- `validate_input.go`: garante `Actor.ID != 0`, `Target.ID != 0`, `Target.Type`/`Operation` válidos.
- `builder.go`: aplica defaults (`OccurredAt = time.Now().UTC()`), correlaciona `request_id`/`trace_id` se faltarem usando `globalmodel.RequestIDKey` e tracing (`utils.GetTraceID` se houver helper) ou logger fields.

### 4. Adapter MySQL (novo) — `internal/adapter/right/mysql/audit`
- `audit_adapter.go`: struct + `NewAuditAdapter(db, metrics)` com InstrumentedAdapter.
- `create_event.go` (método público):
```go
const insertAudit = `INSERT INTO audit_events
 (occurred_at, actor_id, actor_role, actor_device_id, actor_ip, actor_user_agent,
  target_type, target_id, target_version, operation, metadata, request_id, trace_id)
 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

func (a *AuditAdapter) CreateEvent(ctx context.Context, tx *sql.Tx, event auditmodel.AuditEvent) error {
   ctx, end, _ := utils.GenerateTracer(ctx); defer end(); ctx = utils.ContextWithLogger(ctx); logger := utils.LoggerFromContext(ctx)
   entity := auditconverters.EventDomainToEntity(event)
   res, err := a.ExecContext(ctx, tx, "insert", insertAudit,
      entity.OccurredAt, entity.ActorID, entity.ActorRole, entity.ActorDeviceID, entity.ActorIP, entity.ActorUserAgent,
      entity.TargetType, entity.TargetID, entity.TargetVersion, entity.Operation, entity.Metadata, entity.RequestID, entity.TraceID,
   )
   if err != nil { utils.SetSpanError(ctx, err); logger.Error("mysql.audit.create.exec_error", "err", err); return fmt.Errorf("insert audit_event: %w", err) }
   id, err := res.LastInsertId(); if err != nil { utils.SetSpanError(ctx, err); logger.Error("mysql.audit.create.last_insert_id_error", "err", err); return fmt.Errorf("audit_event last insert id: %w", err) }
   event.SetID(id); return nil
}
```
- `list_events_by_target.go` (opcional Fase 3): SELECT com filtros e paginação.
- `entities/audit_event_entity.go`: struct com comentários do schema, campos `sql.Null*` onde aplicável.
- `converters/event_domain_to_entity.go`: serializa `Metadata` via `json.Marshal` para `[]byte`/`sql.NullString`.
- `converters/event_entity_to_domain.go`: parse de JSON; mapeia `[]byte` -> string.

### 5. Factory/DI
- Atualizar fábrica para criar `AuditAdapter` e injetar em `audit_service`.
- Expor `AuditService` onde `GlobalService` é montado. Opcional: manter `GlobalService.CreateAudit` chamando `auditService.RecordChange` como compat temporária.

### 6. Refino de Constantes
- `global_model/global_constants.go`: alinhar `TableName` ou criar novo enum para `TargetType` com valores do schema real (`listing_identities`, `listing_versions`, `listing_visits`, `proposals`, `media_assets`, `listing_agendas`, `sessions`).

### 7. Instrumentação (fase 2, por domínio)
Para cada serviço, substituir/adiar `globalService.CreateAudit` por `auditService.RecordChange` com metadata adequada:
- Listing create/update/promote/status/publish/unpublish/discard/delete: preencher `TargetType` = listing_identity ou listing_version conforme operação; metadata com status_from/to, version_from/to, price deltas.
- Proposals: create/accept/reject/cancel -> target proposal id, metadata status_from/to.
- Visits/agenda: schedule/approve/reject/cancel/complete/no-show, finish agenda -> target visit/agenda id, status_from/to.
- Mídia: approvals/rejections -> target asset/batch id, status_from/to.
- Auth/sessões: signin/signout/refresh/password reset -> target session/token id (se aplicável), metadata result (success/fail) sem PII.
- Admin/roles: criação/alteração/deleção de usuários do sistema e vínculos (owner/realtor/agency) -> target user/user_role id.

### 8. Consultas/Timeline (fase 3)
- Endpoint interno ou função de suporte para timeline por `target_type + target_id (+version opcional)`, ordenado por `occurred_at DESC`, paginado.
- Validar correlação: request_id/trace_id em cada entrada deve bater com Tempo/Loki.

### 9. Checklist de Qualidade antes do código
- Confirmar com DBA nomes de tabela/colunas/índices.
- Confirmar se `metadata` será JSON ou TEXT; definir política de tamanhos.
- Definir lista branca de campos em metadata por operação para evitar PII.
- Garantir que `Actor.ID` nunca seja zero; se não houver usuário (tarefas automáticas), definir `system` actor id acordado.
