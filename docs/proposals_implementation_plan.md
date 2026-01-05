# Plano de Implementação – Sistema de Propostas (TOQ Server)

## 1. Objetivo
Registrar um plano detalhado, faseado e executável por times distintos para implementar o fluxo completo de propostas entre corretores e proprietários, conforme regras de negócio e padrões definidos no guia Go do TOQ Server.

## 2. Escopo
- Criar domínio de propostas (model, ports, adapters, services, handlers) seguindo Arquitetura Hexagonal e Regra de Espelhamento.
- Suportar envio (texto livre) e/ou anexo PDF (até 1MB), edição enquanto `pending`, cancelamento, aceite, recusa com motivo, histórico e listagens.
- Adicionar sinalização em listing sobre proposta pendente/aceita.
- Enviar notificações push (sistema unificado) para owner e realtor em mudanças de status.
- Fora de escopo: scripts de migração, testes automatizados, alteração manual de swagger.json/yaml.

## 3. Referências Obrigatórias
- Guia: docs/toq_server_go_guide.md (seções 2, 7, 8, Regra de Espelhamento, tracing/logs, 1 função pública por arquivo).
- Observabilidade/execução: README.md.
- Esquema atual de BD: scripts/db_creation.sql.

## 4. Requisitos de Negócio (resumo)
- Realtor: criar proposta (texto ou PDF), listar suas propostas, ver status, editar se `pending`, cancelar antes da aceitação.
- Owner: listar recebidas, aceitar, recusar com motivo, ver histórico.
- Notificações push: envio ao owner na criação/cancelamento; envio ao realtor em aceite/recusa.
- Status suportados: `pending`, `accepted`, `rejected`, `cancelled` (expiração futura opcional via worker).
- Listing expõe flag de proposta pendente/aceita.
- Campos financeiros e dados de cliente **não** fazem parte da proposta (removidos de payload, modelo e base). Proposta contém apenas texto livre e anexos PDF.

## 5. Modelo de Dados Proposto (sem migração aqui)
- Tabela `proposals` (InnoDB, utf8mb4):
  - Chaves: `id` PK; `listing_identity_id` FK; `realtor_id` FK users; `owner_id` FK users.
  - Campos: `status` ENUM(`pending`,`accepted`,`rejected`,`cancelled`), `proposal_text` TEXT NULL, `rejection_reason` VARCHAR(500) NULL, `accepted_at` DATETIME NULL, `rejected_at` DATETIME NULL, `cancelled_at` DATETIME NULL, `created_at` DATETIME(6) NOT NULL, `updated_at` DATETIME(6) NOT NULL, `deleted` TINYINT NOT NULL DEFAULT 0.
  - Índices: (`listing_identity_id`,`status`), (`realtor_id`,`status`), (`owner_id`,`status`).
- Tabela `proposal_documents`: `id` PK; `proposal_id` FK; `file_name`, `file_type`, `file_url`, `file_size_bytes` BIGINT, `uploaded_at` DATETIME(6).
- Colunas em `listing_identities`: `has_pending_proposal` TINYINT DEFAULT 0; `has_accepted_proposal` TINYINT DEFAULT 0; `accepted_proposal_id` INT UNSIGNED NULL FK proposals(id).

## 6. Arquitetura e Pastas (Regra de Espelhamento)
- Domínio: internal/core/model/proposal_model/ (interfaces, enums simples).
- Ports (right): internal/core/port/right/repository/proposal_repository/.
- Adapter MySQL: internal/adapter/right/mysql/proposal/ com proposal_adapter.go (struct+New) e um arquivo por método público; subpastas entities/ e converters/.
- Service: internal/core/service/proposal_service/ (interface em proposal_service.go, um método público por arquivo).
- Ports (left): internal/core/port/left/http/proposalhandler/.
- Handlers: internal/adapter/left/http/handlers/proposal_handlers/ (um arquivo por endpoint + Swagger); DTOs em internal/adapter/left/http/dto/.
- Rotas: ajustes em internal/adapter/left/http/routes/routes.go para /api/v2/proposals.
- Factory wiring: adicionar proposal adapter/service/handler em factories.

## 7. Fases e Entregáveis (podem rodar em paralelo quando indicado)
### Fase 0 — Alinhamento e contratos
- Revisar requisitos e DTOs alvo com produto.
- Validar impactos de permissão (roles realtor/owner) com time de Auth/Permissão.

### Fase 1 — Domínio e Ports (Time Core)
- Criar enum `Status` (`pending`,`accepted`,`rejected`,`cancelled`).
- Criar interfaces `ProposalInterface`, `ProposalDocumentInterface`, filtros/pagination structs.
- Publicar Port `ProposalRepoPortInterface` com assinaturas: create, update (texto), update_status (with expected current), get_by_id, list_realtor, list_owner, add_document, list_documents, mark_listing_flags.

### Fase 2 — Adapter MySQL (Time Dados)
- Entities mapeando colunas com sql.Null* (texto, rejection_reason, timestamps de status).
- Converters entity↔domain isolados.
- Métodos (um por arquivo) usando InstrumentedAdapter e colunas explícitas: create_proposal, update_proposal, update_proposal_status (check status atual), get_proposal_by_id, list_realtor_proposals (filtros + paginação), list_owner_proposals, add_document, list_documents, mark_listing_proposal_flags (atualiza flags em listing_identities).
- Tracing em cada método (utils.GenerateTracer), utils.SetSpanError em erros de infra.

### Fase 3 — Services (Time Core)
- proposalService com injeções de repos (proposal, listing, user), gsi (tx/audit), notifier, storage.
- Métodos públicos (arquivos dedicados):
  - CreateProposal: valida ownership do listing, bloqueia múltiplas pending/accepted, persiste texto, seta flags, notifica owner.
  - UpdateProposal: somente status pending, atualiza texto, registra audit.
  - CancelProposal: apenas realtor autor e status pending; seta status, flags, notifica owner.
  - AcceptProposal / RejectProposal: apenas owner do listing; valida status pending; registra motivo na recusa; atualiza flags; notifica realtor.
  - ListRealtorProposals / ListOwnerProposals: paginação, filtros por status/listing/data; retorno com contagem total.
  - GetProposalDetail: checa permissão (realtor autor ou owner do listing); inclui documentos.
  - AddProposalDocument / RequestUploadURL: gera URL pré-assinada, valida tamanho ≤ 1MB, salva metadados do PDF.
- Todas com tracer, tx, rollback/commit e mapeamento de erros de domínio vs infra.

### Fase 4 — DTOs e Handlers (Time HTTP)
- DTOs request/response com comentários em inglês, tags json/binding/example; sem campos financeiros ou dados de cliente.
- Handlers em proposal_handlers com Swagger completo (seção 8 do guia), usando http_errors.SendHTTPErrorObj.
- Endpoints:
  - POST /proposals (create; texto opcional, sem campos financeiros).
  - PUT /proposals/update (editar texto; apenas pending).
  - POST /proposals/cancel.
  - POST /proposals/accept.
  - POST /proposals/reject (motivo obrigatório).
  - GET /proposals/realtor (filtros: status, listingId, data; sem faixa de valor).
  - GET /proposals/owner (filtros: status, listingId, data).
  - POST /proposals/detail (id no body).
  - POST /proposals/documents/upload-url (PDF até 1MB).

### Fase 5 — Notificações e UX (Time Core + Infra)
- Mensagens push via UnifiedNotificationService (FCM) usando template existente como base (title, body, orientation_msg, data).
- Variantes:
  - Criação → destino owner: status pending.
  - Cancelamento → destino owner: status cancelled.
  - Aceite → destino realtor: status accepted.
  - Recusa → destino realtor: status rejected com rejection_reason.
- Payload mínimo: title/body/orientation_msg + data (proposal_id, listing_identity_id, status, role).

### Fase 6 — Observabilidade e Auditoria (Time Observabilidade)
- Logs slog: info em eventos de domínio (criado, status alterado), error só infra.
- Tracing nos adapters/services; garantir propagação de logger no ctx.
- Audit: gsi.CreateAudit para mudanças de status.

### Fase 7 — Hardening e Rollout
- Revisão de segurança: ownership/role em cada endpoint.
- Performance: checar índices; paginação com limit/offset.
- Worker futuro opcional para expirar propostas antigas (status passaria a cancelled/expired).
- Checklist contra guia (seções 2 e 7): 1 função pública por arquivo; espelhamento Port/Adapter; SELECT colunas explícitas.

## 8. Divisão por Times
- Core Domain/Service: model, ports, services, regras de negócio, validação de status.
- Dados/DB: entities, converters, queries, índices, flags de listing.
- HTTP/API: DTOs, handlers, rotas, binding/validação e Swagger.
- Infra/Obs: notificações, audit, tracing/logs verificados, storage upload URL.
- QA/Produto: validação funcional, cenários de status, notificações, permissões.

## 9. Riscos e Mitigações
- Concorrência em status: usar update com expected current status; retornar conflito em mismatch.
- Flags de listing: manter consistência via transação (atualizar proposals e listing_identities no mesmo tx).
- Upload: validar tamanho/tipo (PDF, ≤1MB) antes de gerar URL; registrar metadados.
- Ausência de campos financeiros: garantir alinhamento de produto/UX e mensagens de notificação coerentes.

## 10. Entregáveis por Fase
- F1: Model + Ports aprovados em review.
- F2: Adapter MySQL com queries e conversores revisados.
- F3: Services com transações, notificações integradas e auditoria.
- F4: DTOs/Handlers/Rotas com Swagger gerado (make swagger).
- F5: Mensagens push testadas em ambiente homo.
- F6: Logs/Traces/Audit validados em Grafana/Tempo.
- F7: Go-live checklist e handover para QA.
