# Sistema de Permissionamento

Este documento descreve como o permissionamento funciona no projeto, como mapear endpoints HTTP para permissões e como criar, editar e remover permissões de forma consistente com a arquitetura hexagonal.

## Visão geral das entidades

- Roles (`roles`)
  - Perfis de usuário (ex.: root, owner, realtor, agency).
  - Campos principais: `id`, `name`, `slug`, `is_system_role`, `is_active`.
- Permissions (`permissions`)
  - Ações de negócio identificadas por `resource` e `action`.
  - Campos principais: `id`, `name`, `resource`, `action`, `description` (opcional), `conditions` (JSON opcional), `is_active`.
  - Observação: não existe coluna física `slug` em `permissions`. Quando necessário, o código deriva `slug` como `CONCAT(resource, ':', action)`.
- RolePermissions (`role_permissions`)
  - Vincula roles a permissions com um flag `granted` e `conditions` (JSON) opcionais.
  - Campos: `id`, `role_id`, `permission_id`, `granted`, `conditions`.
- UserRoles (`user_roles`)
  - Vincula usuários a roles com estado (`is_active`, `status`, `expires_at`, etc.).
  - Status `deleted` (13) é aplicado quando o usuário solicita exclusão de conta; o status legado `invite_pending` foi descontinuado e convites agora permanecem em `pending_manual` até aceite.

## Como a checagem de permissão funciona (HTTP)

1. O serviço resolve o role ativo do usuário (apenas `user_roles.is_active = 1` e não expirado).
2. Para endpoints HTTP, o mapeamento é:
   - `resource = "http"`
  - `action = "<METHOD>:<PATH>"` (ex.: `POST:/api/v2/user/signout`).
3. O serviço consulta as permissões efetivas do usuário via `role_permissions` + `permissions` (somente `granted = 1` e `permissions.is_active = 1`).
4. Se existir uma permissão com `resource=http` e `action` exatamente igual ao endpoint, o acesso é permitido (condições podem restringir; ver abaixo).
5. `conditions` (JSON) em `permissions` ou `role_permissions` podem aplicar regras (ex.: `{ "owner_only": true }`), avaliadas contra o contexto do usuário.

Observação importante: o caminho deve casar exatamente (método e path) com o roteamento real da API. Diferenças como `/auth/signout` vs `/user/signout` causam negação.

## "Slug" de permissão

- Não existe `permissions.slug` no banco.
- Para conveniência em listagens, o repositório seleciona `CONCAT(resource, ':', action) AS slug`.

## Operações comuns

### Criar uma permissão

1) Inserir na tabela `permissions`:
- `name`: nome amigável (ex.: "HTTP SignOut").
- `resource`: para HTTP, sempre `http`.
- `action`: exatamente `METHOD:/path` (ex.: `POST:/api/v2/user/signout`).
- `description`: descrição (opcional).
- `conditions`: JSON válido ou `NULL` (opcional).
- `is_active`: `1` para ativa.

2) Conceder a roles na `role_permissions`:
- Uma linha por role com `(role_id, permission_id, granted=1, conditions=NULL ou JSON)`.

3) Nenhuma mudança de código é necessária; o sistema é data-driven.

Exemplo SQL (ilustrativo):

- Criar a permissão:
  - `INSERT INTO permissions (name, resource, action, description, conditions, is_active)
     VALUES ('HTTP SignOut', 'http', 'POST:/api/v2/user/signout', 'Permite encerrar a sessão', NULL, 1);`
- Conceder a todos os roles (IDs 1..7, exemplo):
  - `INSERT INTO role_permissions (role_id, permission_id, granted, conditions)
     VALUES (1, <perm_id>, 1, NULL), (2, <perm_id>, 1, NULL), ... ;`

### Editar uma permissão

- Atualizar campos em `permissions` (ex.: `description`, `is_active`).
- Se alterar `resource`/`action`, a identidade lógica da permissão muda; mantenha o mesmo `id` para preservar concessões e garanta que o novo `action` corresponda ao endpoint real.
- Ajuste `conditions` em `role_permissions` por role quando necessário.

Exemplos SQL:
- `UPDATE permissions SET description='Nova descrição' WHERE id=<perm_id>;`
- `UPDATE permissions SET is_active=0 WHERE id=<perm_id>;` (desativação global)
- `UPDATE role_permissions SET granted=0 WHERE role_id=2 AND permission_id=<perm_id>;` (revogar para um role específico)

### Remover uma permissão

- Revogação por role: `UPDATE role_permissions SET granted=0 ...` ou `DELETE FROM role_permissions ...`.
- Desativação global: `UPDATE permissions SET is_active=0 WHERE id=<perm_id>;`.
- Exclusão: `DELETE FROM permissions WHERE id=<perm_id>;` (FKs de `role_permissions` tratam a cascata conforme configurado).

## Uso dos CSVs (seeds)

- `data/base_permissions.csv`: define permissões (`id;name;resource;action;description;conditions;is_active`).
- `data/base_permission_roles.csv`: define roles (`id;name;slug;description;is_system_role;is_active`).
- `data/base_role_permissions.csv`: mapeia role → permission (`id;role_id;permission_id;granted;conditions`).

Exemplo (adicionar signout para todos os usuários via CSV):

- Em `data/base_permissions.csv` (próximo `id` livre, ex.: `33`):
  - `33;HTTP SignOut;http;POST:/api/v2/user/signout;Permite encerrar a sessão;NULL;1`
- Em `data/base_role_permissions.csv` (continuando a sequência de `id`):
  - `75;1;33;1;NULL`
  - `76;2;33;1;NULL`
  - `77;3;33;1;NULL`
  - `78;4;33;1;NULL`
  - `79;5;33;1;NULL`
  - `80;6;33;1;NULL`
  - `81;7;33;1;NULL`

Se o endpoint legado `/api/v1/auth/signout` também existir, crie outra permissão (ex.: `34`) com `action=POST:/api/v1/auth/signout` e repita as concessões.

### Permissão para download administrativo dos documentos CRECI

- `resource = http`
- `action = POST:/api/v2/admin/users/creci/download-url`
- Concedida inicialmente ao role **Administrador** (`role_id = 1`) em `data/base_role_permissions.csv`.
- Descrição: possibilita ao time admin gerar URLs assinadas para selfie/front/back do CRECI do usuário alvo.

## Cache de permissões e invalidação

- As permissões agregadas de cada usuário são materializadas no Redis (`toq_cache:user_permissions:<id>`) por até 15 minutos.
- Operações de domínio que alteram roles/permissões agora invalidam automaticamente o cache afetado:
  - `AssignRoleToUser`, `RemoveRoleFromUser`, `SwitchActiveRole`, `ActivateUserRole`, `DeactivateAllUserRoles`.
  - `GrantPermissionToRole` e `RevokePermissionFromRole` invalidam todos os usuários que possuem o role impactado.
- Em cenários de atualização direta via SQL (ex.: scripts de manutenção), invoque `permissionService.RefreshUserPermissions(ctx, userID)` após concluir as alterações.
- O serviço `HasPermission` faz um "refresh-on-miss": se o cache existir mas não contiver a permissão requerida, ele força a recarga a partir do banco e reavalia o acesso.
- Falhas na invalidação são registradas com `permission.cache.invalidate_safe_failed` e marcadas no trace para facilitar investigação.
- Observabilidade: cada operação de cache (lookup, store, invalidate, refresh) emite contadores no Prometheus (`cache_operations_total`), diferenciando `operation` e `result` (`hit`, `miss`, `success`, `error`, `disabled`).

## Boas práticas e armadilhas

- Case exato do `action` HTTP (método e path) com o roteamento real.
- `conditions` deve ser JSON válido; use `NULL` quando não houver regra.
- Permissões só valem com `permissions.is_active = 1` e `role_permissions.granted = 1`.
- O usuário precisa ter um role ativo e não expirado em `user_roles`.
- Evite concessões redundantes; mantenha consistência por role.

## Notas de arquitetura (Hexagonal)

- O permissionamento é aplicado via Services que usam Repositories; Handlers não acessam DB diretamente.
- Repositórios ficam em `internal/adapter/right/mysql/...` e derivam `slug` quando necessário com `CONCAT(resource, ':', action)`.
- Transações são gerenciadas pelos serviços globais (global service) e injetadas via factories na inicialização.
- Alterações de permissão são 100% data-driven, dispensando mudanças de código.
