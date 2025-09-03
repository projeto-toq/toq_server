# Política de Status do Usuário (FSM via YAML)

Este guia explica como usar e manter o novo sistema de transição de status do usuário. Ele substitui o update_user_status legado por uma política baseada em regras (tabela/YAML), obedecendo a Arquitetura Hexagonal, transações e auditoria.

## Visão geral

- Motor de decisão: porta `UserStatusPolicy` avalia transições.
  - Entrada: (role, fromStatus, action)
  - Saída: (toStatus, notification, changed)
- Fonte das regras: arquivo YAML configurado via `auth.status_rules_path` (env).
- Persistência: atualização via repositório dentro de transação.
- Orquestração: serviços consultam a política, persistem a mudança, criam auditoria e, após commit, podem notificar.

Arquitetura (arquivos principais):
- Porta: `internal/core/port/policy/user_status_policy_port.go`
- Implementação (carrega do YAML): `internal/adapter/right/config/user_status_rules_loader.go`
- Loader YAML: `internal/adapter/right/config/user_status_rule_source.go`
- Orquestradores: `internal/core/service/user_service/user_status_orchestrator.go`
- Repositório (MySQL): `internal/adapter/right/mysql/user/update_user_role_status_tx.go`
- DI/wiring: `internal/core/config/inject_dependencies.go`
- Regras (YAML): `configs/status_rules.yaml`

## Habilitar e configurar

1) Defina o caminho do YAML em `configs/env.yaml`:

```yaml
auth:
  status_rules_path: "configs/status_rules.yaml"
```

2) Reinicie o serviço após editar o YAML (ou implemente endpoint admin chamando `UserStatusPolicy.Reload`).

Logs úteis:
- "User status policy rules loaded" (contagem de regras)
- "No matching user status transition rule" quando não há regra correspondente

## Como os serviços usam

Orquestradores prontos e integrados:
- `ApplyUserStatusTransitionAfterEmailConfirmed(ctx)`
- `ApplyUserStatusTransitionAfterPhoneConfirmed(ctx)`

Fluxo:
1) Serviço determina a ação (ex.: verificação concluída vs. outro fator pendente).
2) Política `Evaluate(role, from, action)` decide a transição.
3) Se `changed`, serviço chama repo `UpdateUserRoleStatus` dentro da mesma transação e cria auditoria.
4) Após commit, se `notification != 0`, encaminha notificação via Unified Notification Service.

## Criar/editar o arquivo `configs/status_rules.yaml`

Cada item mapeia (role, from, action) -> (to, notification), com `priority` opcional (maior vence em caso de sobreposição).

Campos por regra:
- role: string (slug da role: `owner`, `realtor`, `agency`, `root`, ...)
- from: int (permission_model.UserRoleStatus)
- action: int (user_model.ActionFinished)
- to: int (permission_model.UserRoleStatus)
- notification: int (global_model.NotificationType) — 0 para nenhuma
- priority: int (opcional; maior valor vence)

Exemplo mínimo:

```yaml
- role: owner
  from: 4    # StatusPendingEmail
  action: 16 # ActionProfileEmailVerifiedPhonePending
  to: 5      # StatusPendingPhone
  notification: 0
  priority: 10

- role: owner
  from: 4    # StatusPendingEmail
  action: 18 # ActionProfileVerificationCompleted
  to: 0      # StatusActive
  notification: 0
  priority: 10
```

Onde encontrar valores das enums:
- Status (UserRoleStatus): `internal/core/model/permission_model/constants.go`
- Ações (ActionFinished): `internal/core/model/user_model/constants.go`
- Notificações (NotificationType): `internal/core/model/global_model`

Boas práticas ao editar:
- Casamento é exato: não há curingas. Adicione regras explícitas para cada caso esperado.
- Use `priority` somente quando duas regras poderiam casar; a de maior prioridade vence.
- Comente valores no YAML (como no exemplo) para facilitar revisão.

Passo a passo para criar/editar:
1) Abra `configs/status_rules.yaml` (crie se não existir) e adicione regras conforme o esquema.
2) Confirme que os slugs de role existem em `base_roles.slug`.
3) Valide números contra as enums nos arquivos acima.
4) Reinicie o serviço (ou acione `Reload`).
5) Verifique os logs após o boot.

## Expandir para outros fluxos

- Fluxos de convite (criado/aceito/rejeitado) podem ser migrados:
  - Adicione regras para `ActionFinishedInviteCreated`, `ActionFinishedInviteAccepted`, `ActionFinishedInviteRejected` por role/from.
  - Crie um orquestrador genérico que aceite `userID` alvo + `action` (para atualizar outro usuário na mesma transação) e substitua chamadas legadas.

## Solução de problemas

- Status não mudou: confira logs “No matching user status transition rule” e verifique (role, from, action). Revise o YAML e o path em `env.yaml`.
- Erro ao ler YAML: corrija indentação/campos. O loader espera exatamente os nomes acima.
- Status inesperado: verifique sobreposição de regras e ajuste `priority`.

## Nota operacional

- A política não tem efeitos colaterais; apenas decide. Persistência/auditoria ficam no serviço.
- Notificações são enviadas após commit para evitar falsos positivos.
- Mantenha regras curtas e explícitas.
