# Plano de migração para `auditService.RecordChange`

## Objetivo
Migrar todas as chamadas legadas de `globalService.CreateAudit` para `auditService.RecordChange`, preenchendo `Actor`, `Target` (Type/ID/Version), `Operation` e `Metadata` completos, preservando transacionalidade e telemetria.

## Arquivos que exigem refatoração
- **Infra/contratos**
  - [internal/core/service/global_service/global_service.go](internal/core/service/global_service/global_service.go)
  - [internal/core/service/global_service/create_audit.go](internal/core/service/global_service/create_audit.go)
  - [internal/core/service/audit_service/audit_service.go](internal/core/service/audit_service/audit_service.go)
  - [internal/core/service/audit_service/record_change.go](internal/core/service/audit_service/record_change.go)
  - [internal/core/service/audit_service/validate_input.go](internal/core/service/audit_service/validate_input.go)
  - (novo) `internal/core/service/audit_service/record_builder.go` (helper para montar `RecordInput`)
  - [internal/core/model/audit_model/*](internal/core/model/audit_model/) (caso precise novos `Operation`/`TargetType`)
  - [internal/core/config/inject_dependencies.go](internal/core/config/inject_dependencies.go)
- **DI nos services** (adição de `auditService` e remoção do uso da fachada)
  - [internal/core/service/listing_service/listing_service.go](internal/core/service/listing_service/listing_service.go)
  - [internal/core/service/proposal_service/proposal_service.go](internal/core/service/proposal_service/proposal_service.go)
  - [internal/core/service/schedule_service/schedule_service.go](internal/core/service/schedule_service/schedule_service.go)
  - [internal/core/service/user_service/user_service.go](internal/core/service/user_service/user_service.go)
- **Proposals (substituir CreateAudit por RecordChange)**
  - [internal/core/service/proposal_service/create_proposal.go](internal/core/service/proposal_service/create_proposal.go)
  - [internal/core/service/proposal_service/update_proposal.go](internal/core/service/proposal_service/update_proposal.go)
  - [internal/core/service/proposal_service/cancel_proposal.go](internal/core/service/proposal_service/cancel_proposal.go)
  - [internal/core/service/proposal_service/accept_proposal.go](internal/core/service/proposal_service/accept_proposal.go)
  - [internal/core/service/proposal_service/reject_proposal.go](internal/core/service/proposal_service/reject_proposal.go)
- **Listings**
  - [internal/core/service/listing_service/create_listing.go](internal/core/service/listing_service/create_listing.go)
  - [internal/core/service/listing_service/update_listing.go](internal/core/service/listing_service/update_listing.go)
  - [internal/core/service/listing_service/promote_listing_version.go](internal/core/service/listing_service/promote_listing_version.go)
  - [internal/core/service/listing_service/discard_draft_version.go](internal/core/service/listing_service/discard_draft_version.go)
  - [internal/core/service/listing_service/end_update_listing.go](internal/core/service/listing_service/end_update_listing.go)
  - [internal/core/service/listing_service/change_listing_status.go](internal/core/service/listing_service/change_listing_status.go)
- **Schedule**
  - [internal/core/service/schedule_service/finish_listing_agenda.go](internal/core/service/schedule_service/finish_listing_agenda.go)
- **Users**
  - Criação/atualização: [create_owner.go](internal/core/service/user_service/create_owner.go), [create_realtor.go](internal/core/service/user_service/create_realtor.go), [create_agency.go](internal/core/service/user_service/create_agency.go), [create_system_user.go](internal/core/service/user_service/create_system_user.go), [update_system_user.go](internal/core/service/user_service/update_system_user.go), [update_profile.go](internal/core/service/user_service/update_profile.go)
  - Segurança/contatos: [confirm_email_change.go](internal/core/service/user_service/confirm_email_change.go), [confirm_phone_change.go](internal/core/service/user_service/confirm_phone_change.go), [confirm_password_change.go](internal/core/service/user_service/confirm_password_change.go), [user_status_orchestrator.go](internal/core/service/user_service/user_status_orchestrator.go)
  - Opt-in/out: [push_optin.go](internal/core/service/user_service/push_optin.go), [push_optout.go](internal/core/service/user_service/push_optout.go), [update_opt_status.go](internal/core/service/user_service/update_opt_status.go)
  - Sessões: [signout.go](internal/core/service/user_service/signout.go)
  - Vínculos e convites: [invite_realtor.go](internal/core/service/user_service/invite_realtor.go), [accept_invitation.go](internal/core/service/user_service/accept_invitation.go), [reject_invitation.go](internal/core/service/user_service/reject_invitation.go), [delete_agency_of_realtor.go](internal/core/service/user_service/delete_agency_of_realtor.go), [delete_realtor_of_agency.go](internal/core/service/user_service/delete_realtor_of_agency.go), [add_alternative_role.go](internal/core/service/user_service/add_alternative_role.go), [approve_creci_manual.go](internal/core/service/user_service/approve_creci_manual.go)
  - Encerramento: [delete_account.go](internal/core/service/user_service/delete_account.go), [delete_system_user.go](internal/core/service/user_service/delete_system_user.go)

## Fases e interdependências
- **Fase 0 – Alinhamento (concluída)**
  - `AuditOperation`: já cobre os fluxos mapeados; não precisa novo enum neste momento.
  - `TargetType`: adicionar `users`, `user_roles`, `agency_invites`, `realtors_agency` para suportar fluxos de usuários/convites/vínculos.
  - Campos mínimos por domínio:
    - Proposals: `TargetType=proposals`, `Target.ID=proposal.ID`, metadata com `listing_identity_id`, `actor_role`, `reason`, transição de status.
    - Listings: `TargetType=listing_identities`, `Target.ID=identityID`, `Target.Version=versionNumber`, metadata com status `from/to`, `version_id`, contexto de agenda/timezone quando houver.
    - Schedule: `TargetType=listing_agendas` (agenda) e, quando aplicável, `listing_identities`; metadata com `agenda_id`, `listing_identity_id`, status `from/to`.
    - Users: `TargetType` conforme fluxo (`users`, `user_roles`, `agency_invites`, `realtors_agency`), metadata com campos alterados, `reason/origin`, IDs correlatos (invitation/agency), e uso de `ActorFromContext` para device/IP/UA.
- **Fase 1 – Infra/DI (concluída)**
  - `auditService` injetado em listing/proposal/schedule/user services e construtores/fábricas ajustados em [internal/core/config/inject_dependencies.go](internal/core/config/inject_dependencies.go).
  - `CreateAudit` mantido apenas para compatibilidade; sem novos usos.
- **Fase 2 – Helper e contratos (concluída)**
  - Helper `BuildRecordFromContext` criado em `internal/core/service/audit_service/record_builder.go` (Actor do contexto + fallback userID, correlation request/trace, version default=0, metadata normalizada).
  - `TargetType` ampliado com `users`, `user_roles`, `agency_invites`, `realtors_agency` em [internal/core/model/audit_model/audit_event.go](internal/core/model/audit_model/audit_event.go).
  - Nenhuma nova `AuditOperation` necessária nesta fase.
- **Fase 3 – Migração por domínio (concluída)**
  - **Proposals (concluído)**: create/update/cancel/accept/reject migrados para `RecordChange` com helper e metadata (`listing_identity_id`, `proposal_id`, `owner_id`, `realtor_id`, `actor_role`, `status_from/to`, `reason` quando aplicável).
  - **Listings (concluído)**: create/update/promote/discard/end-update/change-status migrados para `RecordChange` com metadata (`listing_identity_id`, `listing_version_id`, `version`, `status_from/to`, `actor_role`, `action/timezone/code/price_changed` quando aplicável) e `TargetType=listing_identities` com `Target.ID` e `Target.Version` preenchidos.
    - **Schedule (concluído)**: finish agenda migrado para `RecordChange` com `TargetType=listing_agendas`, metadata (`listing_identity_id`, `listing_version_id`, `agenda_id`, `status_from/to`, `actor_role`) e reuso de contexto para correlacionar request/trace.
  - **Users (concluído)**: todos os fluxos migrados para `RecordChange` com `TargetType` adequado (`users`, `user_roles`, `agency_invites`, `realtors_agency`), metadata detalhada (campos alterados, origem da ação, IDs de convite/relacionamento, status_from/to, device/IP/UA) e reuso de `ActorFromContext`.
- **Fase 4 – Limpeza e salvaguardas (concluída)**
  - Remoção concluída de usos remanescentes de `CreateAudit` nos serviços; mantida apenas a função legado no `globalService` para compatibilidade.
  - `make lint` executado com sucesso; spans/logs revisados nos fluxos migrados.
  - Pendência operacional: validar manualmente inserções em `audit_events` em ambiente de dev (target_id/operation/metadata/request_id/trace_id) conforme plano.

## Execução paralela sugerida
- Time A: Infra/DI + helper (F1, F2) – desbloqueia os demais.
- Time B: Proposals + Listings (F3 subset) – compartilham padrões de version/status.
- Time C: Users + Schedule (F3 subset) – foco em fluxos de segurança, opt-in/out e vínculos.

## Riscos e mitigação
- **Validação do `RecordChange`** exige `Target.ID > 0`; garantir preenchimento antes de migrar cada call site.
- **Enums faltantes**: adicionar operações específicas antes das migrações (F2) para evitar uso indevido de `update` genérico.
- **Metadata grande**: manter campos chave e serializáveis; evitar structs com ponteiros cíclicos.
- **Traços/logs**: reutilizar `ctx` atual para preservar `request_id/trace_id`; não criar spans extras nos handlers.
