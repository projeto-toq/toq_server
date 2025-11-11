# Guia Global para Desenvolvedores Go — TOQ Server

Este é o guia único do projeto. Aqui você encontra: arquitetura (hexagonal), estrutura de pastas, bootstrap, injeção de dependências, padrões por camada (handlers/DTOs, services, repositórios com converters, workers), erros/observabilidade (logging/tracing), transações e checklists de refatoração/PR.

Sumário
- 1. Arquitetura do projeto
- 2. Estrutura de pastas
- 3. Bootstrapping e ciclo de vida
- 4. Injeção de dependências (Factory Pattern)
- 5. Princípios essenciais
- 6. Middleware HTTP (referência)
- 7. Padrões por camada
  - 7.1 Services (métodos públicos)
  - 7.2 Services (métodos privados)
  - 7.3 Repositórios (adapters)
  - 7.4 Handlers HTTP
  - 7.5 Workers/Go routines
- 8. Padrões de Documentação
  - 8.1 Princípios Gerais
  - 8.2 Handlers (HTTP Left Adapters)
  - 8.3 DTOs (Data Transfer Objects)
  - 8.4 Services (Core Business Logic)
  - 8.5 Repositories (Right Adapters)
  - 8.6 Entities e Converters
  - 8.7 Models (Domain)
  - 8.8 Helpers e Utils
  - 8.9 Factories
  - 8.10 Checklist de Documentação
  - 8.11 Ferramentas e Automação
- 9. Marcação de erro no span
- 10. Nomenclatura e campos dos logs
- 11. Propagação de erros (detalhado)
- 12. Organização de interfaces e arquivos
- 13. Padrão para análise/refatoração
- 14. Anti‑padrões a evitar
- 15. Checklist de PR (observabilidade)
- 16. Exemplos rápidos
- 17. Referências

## 1. Arquitetura do projeto

- Estilo: Arquitetura Hexagonal (Ports & Adapters).
- Fluxo: Handlers (Left) → Services (Core) → Repositórios/Providers (Right Adapters).
- Responsabilidades:
  - Handlers: borda HTTP. Convertem DTOs ⇄ domínio; chamam Services; serializam resultado/erros.
  - Services: regras de negócio e orquestração; transações; mapeiam erros de infra → domínio; sem HTTP.
  - Right Adapters: implementam Ports (repositórios, e-mail, SMS, S3 etc.); fazem I/O; retornam erros “puros”.
  - Ports: interfaces do domínio (left/right). Mantêm o core desacoplado de frameworks/vendors.

## 2. Estrutura de pastas

- cmd/ — Entrypoints (ex.: `toq_server.go`).
- internal/
  - adapter/
    - left/http/ — Handlers, middlewares, serialização de erro.
      - handlers/admin_handlers/ — um arquivo por endpoint admin, mantendo handlers curtos e focados.
    - right/mysql/ — Repositórios MySQL (implementações de Ports).
    - right/* — Providers externos (aws_s3, sms, email, fcm, ...).
  - core/
    - config/ — Bootstrap (fases 01–08), lifecycle, telemetry, HTTP server.
    - factory/ — Abstract Factory para criar adapters/handlers.
    - port/ — Ports (interfaces) left/right (HTTP e repositórios/providers).
    - service/ — Services; um arquivo por caso de uso/método público (OBRIGATÓRIO).
    - model/ — Modelos de domínio (interfaces e structs).
    - utils/ — Tracing, errors, converters, validators, etc.
    - events/, go_routines/ — Barramento e rotinas de sistema.
  - templates/ — Auxiliares (se houver).
- docs/ — Este guia, Swagger, outras docs.
- configs/ — `env.yaml`, credenciais de dev.
- Observabilidade — `grafana/`, `otel-collector-config.yaml`, `prometheus.yml`.
- scripts/ — SQL e utilitários de setup em dev.

Observação: modelos de domínio ficam em `internal/core/model/*`. Não importar pacotes HTTP em modelos.

### 2.1 Regra de Espelhamento Port ↔ Adapter (Arquitetura Hexagonal)

**Princípio Fundamental:** A estrutura de diretórios dos adapters DEVE espelhar a estrutura dos ports para garantir navegabilidade, clareza arquitetural e aderência à Arquitetura Hexagonal.

**Regra Obrigatória:**

Para cada Port definido em:
```
internal/core/port/right/repository/{DOMAIN_NAME}/
```

DEVE existir um Adapter correspondente em:
```
internal/adapter/right/mysql/{DOMAIN_NAME}/
```

**Estrutura Padrão de um Adapter:**

```
internal/adapter/right/mysql/{DOMAIN_NAME}/
├── {domain}_adapter.go          # Struct do adapter + NewFunc APENAS
├── create_{entity}.go            # Um método público por arquivo
├── get_{entity}_by_id.go
├── update_{entity}.go
├── delete_{entity}.go
├── list_{entities}.go
├── converters/                   # Conversões DB ↔ Domain
│   ├── {entity}_entity_to_domain.go
│   └── {entity}_domain_to_entity.go
└── entities/                     # Structs que representam schema DB
    └── {entity}_entity.go
```

**Exemplo Completo:**

```
Port (Interface):
  internal/core/port/right/repository/device_token_repository/
    └── device_token_repo_port.go    # DeviceTokenRepoPortInterface

Adapter (Implementação MySQL):
  internal/adapter/right/mysql/device_token/
    ├── device_token_adapter.go      # DeviceTokenAdapter struct + NewFunc
    ├── add_token.go                 # func (a *DeviceTokenAdapter) AddToken(...)
    ├── remove_token.go
    ├── list_by_user_id.go
    ├── converters/
    │   └── device_token_entity_to_domain.go
    └── entities/
        └── device_token_entity.go
```

**Benefícios desta Regra:**

1. **Navegabilidade 1:1:** Localização de implementações é intuitiva e previsível.
2. **Separação de Responsabilidades (SRP):** Cada adapter gerencia apenas seu domínio.
3. **Testabilidade:** Testes isolados por domínio, facilitando mocks e coverage.
4. **Escalabilidade:** Novos domínios seguem padrão consistente, evitando "god adapters".
5. **Desacoplamento Real:** Facilita substituição de tecnologias (MySQL → PostgreSQL, Redis, etc).
6. **Code Reviews Eficientes:** Revisores sabem exatamente onde procurar implementações.
7. **Onboarding de Novos Desenvolvedores:** Estrutura física reflete estrutura lógica da arquitetura.

**Anti-Padrões a Evitar:**

❌ **Implementar múltiplos domínios em um único adapter:**
```
# ERRADO: device_token implementado dentro de user/
internal/adapter/right/mysql/user/
├── user_adapter.go
├── device_token_repository.go   # ❌ Violação de SRP
└── get_user_by_id.go
```

✅ **Cada domínio tem seu próprio adapter:**
```
# CORRETO: separação clara de responsabilidades
internal/adapter/right/mysql/user/
└── user_adapter.go

internal/adapter/right/mysql/device_token/
└── device_token_adapter.go
```

❌ **Métodos de negócio no arquivo do adapter:**
```go
// user_adapter.go - ERRADO
type UserAdapter struct { ... }
func NewUserAdapter(...) { ... }
func (ua *UserAdapter) CreateUser(...) { ... }  // ❌ Método aqui
```

✅ **Apenas struct e NewFunc no arquivo principal:**
```go
// user_adapter.go - CORRETO
type UserAdapter struct { ... }
func NewUserAdapter(...) { ... }
// create_user.go (arquivo separado)
func (ua *UserAdapter) CreateUser(...) { ... }  // ✅
```

**Regra OBRIGATÓRIA: Uma Função Pública Por Arquivo**

Cada método público de um adapter/service/handler DEVE estar em seu próprio arquivo dedicado. Esta regra é **obrigatória e sem exceções** para garantir:

1. **Manutenibilidade:** Mudanças isoladas com histórico Git limpo (1 arquivo = 1 responsabilidade)
2. **Legibilidade:** Desenvolvedores localizam funções rapidamente por nome de arquivo
3. **Code Review:** PRs focados e revisões granulares (não revisar múltiplas funções de uma vez)
4. **Testabilidade:** Testes unitários espelham estrutura (1 arquivo de teste por função)
5. **Refatoração Segura:** Mudanças em uma função não afetam acidentalmente outras no mesmo arquivo

**Nomenclatura de Arquivos:**
- Adapters: `{action}_{entity}.go` → `create_user.go`, `get_user_by_id.go`, `update_user_status.go`
- Services: `{action}_{entity}.go` → `confirm_email_change.go`, `revoke_user_session.go`
- Handlers: `{action}_{entity}_handler.go` → `create_user_handler.go`, `update_user_handler.go`

**Exemplo Correto:**
```
internal/adapter/right/mysql/user/
├── user_adapter.go                    # Struct + NewFunc APENAS
├── block_user_temporarily.go          # func (ua *UserAdapter) BlockUserTemporarily(...)
├── unblock_user.go                    # func (ua *UserAdapter) UnblockUser(...)
├── get_expired_temp_blocked_users.go  # func (ua *UserAdapter) GetExpiredTempBlockedUsers(...)
├── create_user.go
├── update_user_by_id.go
└── ...
```

**Exemplo Incorreto (ANTI-PADRÃO):**
```go
// user_blocking.go - ❌ ERRADO: múltiplas funções no mesmo arquivo
func (ua *UserAdapter) BlockUserTemporarily(...) error      // ❌ Deve estar em block_user_temporarily.go
func (ua *UserAdapter) UnblockUser(...) error               // ❌ Deve estar em unblock_user.go
func (ua *UserAdapter) GetExpiredTempBlockedUsers(...) (...)// ❌ Deve estar em get_expired_temp_blocked_users.go
```

**Arquivos Legacy Não Conformes (Requerem Refatoração):**
- `internal/adapter/right/mysql/user/user_blocking.go` - Deve ser dividido em 3 arquivos separados

**Checklist de Conformidade Arquitetural:**

- [ ] Cada Port em `/port/right/repository/{DOMAIN}/` tem Adapter em `/adapter/right/mysql/{DOMAIN}/`
- [ ] Arquivo principal do adapter contém APENAS struct + NewFunc
- [ ] Cada método público está em arquivo próprio (OBRIGATÓRIO, sem exceções)
- [ ] Converters separados em subdiretório `/converters/`
- [ ] Entities separadas em subdiretório `/entities/`
- [ ] Adapter usa `InstrumentedAdapter` para queries (tracing + métricas)
- [ ] Nenhum adapter implementa Ports de outros domínios

## 3. Bootstrapping e ciclo de vida

Orquestrado por `internal/core/config` (struct `Bootstrap`) com Lifecycle Manager e health tracking.

Fases (ordem):
- 01 Context: contexto base, sinais, usuário sistema, ajuste do diretório, pprof opcional.
- 02 Config: carrega env + YAML; valida DB/HTTP/telemetry/security; aplica JWT/TTLs.

### JWT Claims (Access / Refresh)

Access token (claim `infos`):
```
{
  "ID": <int>,
  "UserRoleID": <int>,
  "RoleStatus": <int>,
  "RoleSlug": "<slug>" (root|owner|realtor|agency),
}
```
Refresh token (claim `infos`):
```
{
  "ID": <int>,
  "UserRoleID": 0,
  "RoleStatus": <int>,
  "RoleSlug": "<slug>"
}
```
Backward-compatible: middleware trata ausência de `RoleSlug`.
- 03 Infra: MySQL, Redis, OpenTelemetry (tracing/metrics), adapter de métricas.
- 04 DI: Factory Pattern criando Storage, Repositories, Validation, External Services.
- 05 Services: ordem crítica — Global → Permission → User → Complex → Listing.
- 06 HTTP: servidor, middlewares, rotas/handlers, health checks.
- 07 Workers: workers do sistema, activity tracker, verificação de schema.
- 08 Startup: inicia HTTP, readiness, health monitor, shutdown gracioso.

Shutdown: `Bootstrap.Shutdown()` cancela contexto, aguarda workers e executa cleanup do Lifecycle.

### Perfis de ambiente (`ENVIRONMENT`)
- `ENVIRONMENT=homo` (default) mantém a porta `:8080`, workers ativos e habilita todos os pipelines de telemetria (traces, métricas e logs OTLP → collector).
- `ENVIRONMENT=dev` troca a porta para `127.0.0.1:18080`, desativa os workers em background **e interrompe os pipelines de telemetria** (nenhum sinal é enviado ao collector/Loki; logs permanecem locais).
- Para sobrescrever apenas a porta HTTP sem mudar o perfil, exporte `TOQ_HTTP_PORT` antes de iniciar o binário.

## 4. Injeção de dependências (Factory Pattern)

- Orquestração em `config.InjectDependencies(lm)` com `AdapterFactory`.
- Abstract Factory (`internal/core/factory`):
  - CreateStorageAdapters(ctx, env, db, metrics) → Database, Cache, CloseFunc.
  - CreateRepositoryAdapters(database) → repositórios MySQL.
  - CreateValidationAdapters(env) → CEP, CPF, CNPJ, ...
  - CreateExternalServiceAdapters(ctx, env) → FCM, Email, SMS, Storage.
  - CreateHTTPHandlers(...) → handlers HTTP.
- Lifecycle: adapters podem registrar `CloseFunc` no LifecycleManager para cleanup.

Transações: use o provedor padronizado (global_services/transactions). Services iniciam/commit/rollback. Adapters recebem `ctx` e `*sql.Tx` quando aplicável.

## 5. Princípios essenciais

- Arquitetura por camadas:
  - Handlers → Services → Repositories (Adapters).
  - Injeção de dependências via factories já existente.
- Funções pequenas e focadas: cada função/método deve ter responsabilidade única e clara.
  - Funções que implementam os repositórios/services/handlers/model devem ser separadas dos arquivos de interface/domain.
- Tracing:
  - Use `utils.GenerateTracer(ctx)` no início de métodos públicos de Services, Repositories e em Workers/Go routines.
  - Em Handlers HTTP, NÃO crie spans: o `TelemetryMiddleware` cuida disso.
  - Sempre `defer spanEnd()`.
  - Em erros de infraestrutura, chame `utils.SetSpanError(ctx, err)`.
- Logging (slog):
  - `slog.Info`: eventos esperados de domínio (mudança de status, criação de recurso).
  - `slog.Warn`: anomalias/limites, condições não fatais (ex.: 429/423).
  - `slog.Error`: exclusivamente falhas de infraestrutura (DB, cache, providers externos, transações).
  - Sempre derive o logger via `utils.LoggerFromContext(ctx)` para garantir `request_id`/`trace_id` automáticos; caso precise reaproveitar o logger, utilize `utils.ContextWithLogger(ctx)` antes de repassar para outras funções.
  - Em Repositórios, evite logs verbosos; sucesso no máximo `DEBUG` quando realmente necessário.
- Erros e HTTP:
  - Repositórios retornam erros “puros” (ex.: `sql.ErrNoRows`, `fmt.Errorf`), sem pacotes HTTP.
  - Services mapeiam para erros de domínio (sentinelas/Kind) quando aplicável; erros de infra geram logs + span error.
  - Handlers SEMPRE usam `http_errors.SendHTTPErrorObj(c, err)`; não serialize manualmente.
  - Mensagens enviadas aos clients (400/422) devem ser curtas em inglês simples; o frontend faz a localização para o usuário final.

Observação: o projeto suporta erros de domínio via `derrors` (novo) e tipos `utils.HTTPError` (legado). Em ambos os casos, o handler converte corretamente para HTTP.

## 6. Middleware HTTP (referência)

- Ordem dos middlewares HTTP:
  - `RequestIDMiddleware` → `StructuredLoggingMiddleware` → `CORSMiddleware` → `TelemetryMiddleware` → `ErrorRecoveryMiddleware` → `DeviceContextMiddleware`.
- `StructuredLoggingMiddleware` escreve o log de acesso com os campos (snake_case):
  - `request_id`, `trace_id`, `span_id`, `method`, `path`, `status`, `duration`, `size`, `client_ip`, `user_agent` e, quando disponível, `user_id`, `user_role_id`.
- Severidade automática do log de acesso:
  - 5xx → ERROR; 429/423 → WARN; demais 4xx → INFO; 2xx/3xx → INFO.
- `TelemetryMiddleware` cria o span raiz do request; não crie spans nos handlers.

## 7. Padrões por camada

### 7.1 Services (métodos públicos)

- os serviços terão seu proprio diretório em internal/core/service/<service_name>/**
- o arquivo de interface ficará em internal/core/service/<service_name>/<service_name>_service.go
- o arquivo de interface terá apenas a struct do service, a interface e a func New
- cada método público do service terá seu próprio arquivo em internal/core/service/<service_name>/<method_name>.go
- Inicie transação via serviço global (`s.globalService.StartTransaction
- Inicie tracer no topo; trate domínio vs infraestrutura; marque spans em erros de infra; não serialize HTTP.
- Se houver necessidade de um helper com funções simples e pontuais, comun a várias funções públicas, crie um método privado helper.go. Se for função mais complexa crie um metodo privado.

```go
func (s *someService) DoSomething(ctx context.Context, input Input) (Output, error) {
    ctx, end, err := utils.GenerateTracer(ctx)
    if err != nil { return Output{}, derrors.Infra("Failed to generate tracer", err) }
    defer end()

    tx, err := s.globalService.StartTransaction(ctx)
    if err != nil {
        slog.Error("service.do_something.tx_start_error", "err", err)
        utils.SetSpanError(ctx, err)
        return Output{}, derrors.Infra("Failed to start transaction", err)
    }
    defer func() {
        if err != nil { _ = s.globalService.RollbackTransaction(ctx, tx) }
    }()

    // ... chamar privados/repos; propague erros de domínio; logue infra + SetSpanError.

    if err := s.globalService.CommitTransaction(ctx, tx); err != nil {
        slog.Error("service.do_something.tx_commit_error", "err", err)
        utils.SetSpanError(ctx, err)
        return Output{}, derrors.Infra("Failed to commit transaction", err)
    }
    return out, nil
}
```

Notas:
- Em erros de domínio (regras de negócio), apenas retorne o erro; não logue como `Error`.
- Em erros de infra, faça `slog.Error(...)` no ponto de falha e `utils.SetSpanError`.

### 7.2 Services (métodos privados)

- Não iniciar novo tracer; reutilize `ctx`.
- Faça logs `slog.Error` e `SetSpanError` somente em falhas de infraestrutura.
- Evite logs redundantes; comente os “stages” quando ajuda no diagnóstico.

### 7.3 Repositórios (adapters MySQL, Redis, etc.)

- MySQL utiliza obrigatoriamente o executor compartilhado `SQLExecutor` exposto por `InstrumentedAdapter`. Sempre chame `a.ExecContext`, `a.QueryContext`, `a.QueryRowContext` ou `a.PrepareContext`; nunca invoque `db.ExecContext`/`tx.ExecContext` diretamente.
- O executor gera métricas (via `metricsport.MetricsPortInterface`), logging contextual (`utils.LoggerFromContext`) e tracing automático. Métodos auxiliares `basic_*` foram removidos; não recrie utilitários semelhantes.
- Inicie tracer no início do método (adapters têm tracing mínimo):

```go
ctx, end, _ := utils.GenerateTracer(ctx)
 defer end()
```

- Utilize o executor para executar as queries e delegue o tratamento de span/log dentro do helper:

```go
logger := utils.LoggerFromContext(ctx)
result, err := a.ExecContext(ctx, tx, "insert", insertQuery, args...)
if err != nil {
  utils.SetSpanError(ctx, err)
  logger.Error("mysql.entity.insert.exec_error", "err", err)
  return 0, fmt.Errorf("insert entity: %w", err)
}
```

**REGRA OBRIGATÓRIA: Nunca use SELECT ***

- **TODAS as queries devem listar colunas explicitamente** para evitar:
  - **Schema evolution breaking Scan()**: Adicionar/reordenar colunas no schema quebra `rows.Scan()` silenciosamente
  - **Type assertion errors**: Mudanças de tipo de coluna causam panics em runtime
  - **Performance issues**: `SELECT *` trafega colunas não utilizadas (ex: BLOBs pesados)
  - **Manutenibilidade**: Code reviewers não sabem quais campos são realmente usados

✅ **CORRETO:**
```go
query := `SELECT id, name, email, created_at FROM users WHERE id = ?`
```

❌ **INCORRETO:**
```go
query := `SELECT * FROM users WHERE id = ?`  // NUNCA use SELECT *
```

**Exceção permitida:** `SELECT COUNT(*)` é aceitável pois retorna tipo conhecido (int64).

**Justificativa técnica:**
- Colunas explícitas garantem que `rows.Scan()` sempre recebe valores na ordem esperada
- Previne bugs silenciosos quando schema é alterado (nova coluna adicionada/removida)
- Facilita debug: erro de scan aponta exatamente qual campo/tipo está incorreto
- Query optimizer trabalha melhor com colunas específicas

- Em falhas, `slog.Error` com contexto enxuto (ids, chave da operação) e retorne erro puro (`error`).
- Nunca use pacotes HTTP nesta camada; não mapeie para domínio.
- Sucesso: no máximo `slog.Debug` quando necessário.

Padrão de repositórios:
- Cada função deve ter seu próprio arquivo.
- Conversões entre linhas/DTOs DB ↔ entidades de domínio não devem ser feitas inline na função do repositório.
  - Crie entidades de DB que representam ROWs em `internal/adapter/right/mysql/<repo_name>/entities/*.go`.
  - Centralize conversões em pacotes utilitários por domínio em `internal/adapter/right/mysql/<repo_name>/converters/*.go`.
- Repositório foca em: construir query, executar, checar `RowsAffected`, lidar com `sql.ErrNoRows`, e retornar entidades/domínio já convertidas.
- Em `RowsAffected == 0`, retorne `sql.ErrNoRows`. O mapeamento para "não encontrado/sem pendência" é feito no Service.
- Se houver necessidade de um helper com funções simples e pontuais, comun a várias funções públicas, crie um método privado helper.go. Se for função mais complexa crie um metodo privado.

**Documentação Interna de Repositórios (OBRIGATÓRIO):**

All public repository functions MUST have robust internal documentation in English, explaining the logic in detail to facilitate future maintenance.

For complete templates and detailed examples, see **Section 8.5 (Documentation Standards - Repositories)**.

**Mandatory Documentation Points:**
- ✅ Complete Godoc with description, parameters, and returns
- ✅ Comment explaining tracing initialization
- ✅ Comment on logger context propagation
- ✅ Explanation of business rules in query (e.g., `deleted = 0`)
- ✅ Comment on InstrumentedAdapter usage
- ✅ Explanation of error handling (span marking, logging)
- ✅ Justification for returning `sql.ErrNoRows`
- ✅ Explanation of edge cases (multiple results, empty data)
- ✅ Comment on domain/entity conversion

### 7.4 Handlers HTTP

- Não crie spans nem logs de acesso.
- Ao lidar com erro, use:

```go
httperrors.SendHTTPErrorObj(c, err)
```

- Esse helper serializa `{code,message,details}`, marca o span em erros e permite que o middleware de log inclua os detalhes corretos.

Padrão de DTOs e Handlers:
- DTOs (request/response) residem no lado HTTP (left adapter), desacoplados do domínio. Converta DTO ⇄ domínio na borda do handler/assembler.
- Validação de DTO no handler (ou binder), sem regras de negócio; regras ficam no Service.
- Erros: sempre `http_errors.SendHTTPErrorObj(c, err)`.
- Swagger: documente via anotações no código do handler/DTO e não diretamente nos arquivos json/yaml.
- Convenção obrigatória de verbos HTTP:
  - Endpoints de listagem devem ser `GET` com filtros e paginação via query string.
  - Endpoints de criação, atualização e deleção devem usar `POST`, `PUT` e `DELETE`, respectivamente.
  - Endpoints de detalhe de item devem ser `POST`, recebendo identificadores no corpo.
  - Somente endpoints GET podem receber filtros via query; demais verbos devem transportar identificadores e dados no corpo JSON (sem path params).
- Regras para os endpoints admin:
  - Cada rota deve possuir seu handler em arquivo dedicado dentro de `handlers/admin_handlers/`.
- Todo handler público precisa possuir comentários em inglês com exemplos de utilização (incluindo parâmetros e payloads), garantindo a geração automática da documentação Swagger.

### 7.5 Workers/Go routines

- Inicie tracer com `utils.GenerateTracer(ctx)` e `defer spanEnd()`.
- Para rotinas periódicas de manutenção que rodam sem interação do usuário, derive o contexto com `utils.WithSkipTracing(ctx)` antes de chamar o service para evitar spans duplicados.
- Propague `ctx` pelas chamadas; marque spans em erros de infra.

## 8. Padrões de Documentação

### 8.1 Princípios Gerais

**Language:**
- **All documentation (Godoc, inline comments, Swagger annotations)**: ALWAYS in English

**Documentation Levels:**
1. **External Documentation (Godoc)**: For those who WILL USE/CALL the function
   - Describes WHAT the function does
   - Expected parameters
   - Return values
   - Special cases and edge cases
   
2. **Internal Documentation (inline comments)**: For those who WILL MAINTAIN/UNDERSTAND the implementation
   - Explains HOW the function works
   - Justifies technical decisions
   - Describes complex flow steps
   - Alerts about edge cases and pitfalls

**Mandatory Rules:**
- ✅ Every public function/method MUST have complete Godoc
- ✅ Complex functions MUST have inline comments explaining logic
- ✅ Handlers and DTOs MUST have practical examples for Swagger
- ✅ Comments must be kept updated with code
- ❌ Do not comment the obvious (e.g., `// Set name` for `SetName()`)
- ❌ Do not leave outdated or obsolete comments

---

### 8.2 Handlers (HTTP Left Adapters)

**Location**: `internal/adapter/left/http/handlers/`

**Goal**: Documentation must generate complete Swagger and be clear for API consumers.

**Mandatory Template:**

```go
// CreateUser handles user creation with full validation
//
// @Summary     Create a new user account
// @Description Create a new user account with validation of CPF, email uniqueness, and address via CEP lookup.
//              The endpoint validates all required fields and returns detailed error messages for validation failures.
//              Address fields (street, city, state) are populated from CEP lookup, except number and complement which honor the request payload.
// @Tags        Users
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       Authorization  header  string                  true   "Bearer token for authentication" example(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...)
// @Param       X-Device-Id    header  string                  false  "Device ID (UUIDv4) for tracking user sessions" example(550e8400-e29b-41d4-a716-446655440000)
// @Param       request        body    dto.CreateUserRequest   true   "User creation data with all required fields"
// @Success     201            {object} dto.CreateUserResponse  "User successfully created"
// @Failure     400            {object} dto.ErrorResponse       "Invalid request data (malformed JSON, missing required fields)"
// @Failure     409            {object} dto.ErrorResponse       "User already exists (duplicate CPF or email)"
// @Failure     422            {object} dto.ErrorResponse       "Validation failed (invalid CPF format, CEP not found, etc.)" example({"code":422,"message":"Validation failed","details":{"field":"nationalID","error":"Invalid CPF format"}})
// @Failure     500            {object} dto.ErrorResponse       "Internal server error"
// @Router      /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
    // Note: tracing already provided by TelemetryMiddleware
    ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
    logger := coreutils.LoggerFromContext(ctx)

    // Parse and validate request body
    var request dto.CreateUserRequest
    if err := c.ShouldBindJSON(&request); err != nil {
        httperrors.SendHTTPErrorObj(c, httputils.MapBindingError(err))
        return
    }

    // Convert DTO to domain model
    user, err := h.dtoToDomain(request)
    if err != nil {
        httperrors.SendHTTPErrorObj(c, err)
        return
    }

    // Call service layer
    createdUser, err := h.userService.CreateUser(ctx, user)
    if err != nil {
        httperrors.SendHTTPErrorObj(c, err)
        return
    }

    // Convert domain to response DTO
    response := h.domainToDTO(createdUser)
    c.JSON(http.StatusCreated, response)
}
```

**Mandatory Elements for Handlers:**
1. ✅ `@Summary`: One-line action description
2. ✅ `@Description`: Detailed description including:
   - What the endpoint does
   - Validations applied
   - Data transformations (e.g., CEP lookup)
   - Special behaviors
3. ✅ `@Tags`: Swagger category (Users, Authentication, Listings, etc.)
4. ✅ `@Accept` / `@Produce`: Input/output formats
5. ✅ `@Security`: Authentication type when applicable
6. ✅ `@Param`: ALL parameters (headers, path, query, body) with:
   - Name
   - Location (header/path/query/body)
   - Type
   - Required (true/false)
   - Clear and detailed description
   - Practical example when applicable
7. ✅ `@Success`: Success response with code and type
8. ✅ `@Failure`: ALL possible error responses (400, 401, 403, 404, 409, 422, 500, etc.) with:
   - HTTP code
   - Error type
   - Scenario description that generates the error
   - Example when applicable
9. ✅ `@Router`: Path and HTTP method

**List Handlers (GET with filters):**

```go
// ListUsers retrieves a paginated list of users with optional filters
//
// @Summary     List users with pagination and filters
// @Description Retrieve a paginated list of users with optional filtering by role, status, and search term.
//              Supports sorting by multiple fields and custom page sizes.
//              All timestamp filters accept RFC3339 format and are interpreted in UTC unless timezone is specified.
// @Tags        Users
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       page           query   int     false  "Page number (1-indexed)" minimum(1) default(1) example(1)
// @Param       pageSize       query   int     false  "Number of items per page" minimum(1) maximum(100) default(20) example(20)
// @Param       sortBy         query   string  false  "Field to sort by" Enums(name, email, createdAt) default(createdAt) example(name)
// @Param       sortOrder      query   string  false  "Sort direction" Enums(asc, desc) default(desc) example(asc)
// @Param       role           query   string  false  "Filter by role slug" example(owner)
// @Param       status         query   int     false  "Filter by status (1=active, 2=inactive, 3=blocked)" Enums(1, 2, 3) example(1)
// @Param       search         query   string  false  "Search term for name or email (case-insensitive partial match)" example(john)
// @Param       createdAfter   query   string  false  "Filter users created after this timestamp (RFC3339)" example(2024-01-01T00:00:00Z)
// @Param       createdBefore  query   string  false  "Filter users created before this timestamp (RFC3339)" example(2024-12-31T23:59:59Z)
// @Success     200            {object} dto.ListUsersResponse    "Paginated list of users with metadata"
// @Failure     400            {object} dto.ErrorResponse         "Invalid filter parameters"
// @Failure     401            {object} dto.ErrorResponse         "Unauthorized (missing or invalid token)"
// @Failure     500            {object} dto.ErrorResponse         "Internal server error"
// @Router      /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
    // Implementation...
}
```

**Internal Comments in Handlers:**
```go
func (h *UserHandler) CreateOwner(c *gin.Context) {
    // Note: request tracing already provided by TelemetryMiddleware; avoid duplicate spans
    ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

    var request dto.CreateOwnerRequest
    if err := c.ShouldBindJSON(&request); err != nil {
        httperrors.SendHTTPErrorObj(c, httputils.MapBindingError(err))
        return
    }

    // Validate and parse date fields with precise error attribution
    // (bornAt and creciValidity may have invalid formats)
    bornAt, creciValidity, derr := httputils.ValidateUserDates(request.Owner, "owner")
    if derr != nil {
        httperrors.SendHTTPErrorObj(c, derr)
        return
    }

    // Extract request context for security logging and session metadata
    reqContext := coreutils.ExtractRequestContext(c)

    // Validate Device ID (required for session tracking and push notifications)
    rawDeviceID := c.GetHeader("X-Device-Id")
    trimmedDeviceID := strings.TrimSpace(rawDeviceID)
    if trimmedDeviceID == "" {
        httperrors.SendHTTPErrorObj(c, coreutils.NewHTTPErrorWithSource(
            http.StatusBadRequest, "X-Device-Id header is required"))
        return
    }

    // Call service: creates account and authenticates via standard SignIn flow
    tokens, err := h.userService.CreateOwner(ctx, user, request.Owner.Password, 
        trimmedDeviceToken, reqContext.IPAddress, reqContext.UserAgent)
    if err != nil {
        // Convert domain errors to HTTP response
        httperrors.SendHTTPErrorObj(c, err)
        return
    }

    c.JSON(http.StatusCreated, dto.CreateOwnerResponse{Tokens: tokens})
}
```

---

### 8.3 DTOs (Data Transfer Objects)

**Location**: `internal/adapter/left/http/dto/`

**Goal**: Clear documentation for Swagger generation and API contract validation.

**Template for Request DTOs:**

```go
// CreateUserRequest represents the payload for user creation
//
// This DTO is used for user registration and contains all required fields for account creation.
// Address fields (street, city, state) are populated from CEP lookup service.
// Validation is performed both at struct level (tags) and business logic level (service layer).
type CreateUserRequest struct {
    // FullName is the user's complete legal name
    // Required for legal identification and contracts
    // Example: "João Silva Santos"
    FullName string `json:"fullName" binding:"required,min=2,max=100" example:"João Silva Santos"`
    
    // NickName is the user's display name in the application
    // Used in UI elements and notifications
    // Example: "João"
    NickName string `json:"nickName" binding:"required,min=2,max=50" example:"João"`
    
    // NationalID is the user's CPF (Cadastro de Pessoas Físicas) for individuals
    // or CNPJ (Cadastro Nacional da Pessoa Jurídica) for companies
    // Format: digits only (no punctuation). Validation includes checksum verification.
    // Example CPF: "12345678901"
    // Example CNPJ: "12345678000195"
    NationalID string `json:"nationalID" binding:"required,min=11,max=14" example:"12345678901"`
    
    // CreciNumber is the CRECI (Conselho Regional de Corretores de Imóveis) registration number
    // Required ONLY for realtor role onboarding flows
    // Format: numeric followed by "-F" (e.g., "12345-F")
    // Optional for other roles
    CreciNumber string `json:"creciNumber,omitempty" example:"12345-F"`
    
    // CreciState is the Brazilian state (UF) where CRECI is registered
    // Required when CreciNumber is provided
    // Must be a valid 2-letter Brazilian state code
    // Example: "SP", "RJ", "MG"
    CreciState string `json:"creciState,omitempty" binding:"omitempty,len=2" example:"SP"`
    
    // CreciValidity is the expiration date of CRECI registration
    // Format: "YYYY-MM-DD" (ISO 8601 date)
    // Must be a future date for active realtors
    // Example: "2025-12-31"
    CreciValidity string `json:"creciValidity,omitempty" example:"2025-12-31"`
    
    // BornAt is the user's date of birth
    // Format: "YYYY-MM-DD" (ISO 8601 date)
    // User must be at least 18 years old
    // Example: "1990-05-15"
    BornAt string `json:"bornAt" binding:"required" example:"1990-05-15"`
    
    // PhoneNumber is the user's mobile phone in E.164 international format
    // Used for SMS notifications and two-factor authentication
    // Must include country code with + prefix
    // Example: "+5511999999999" (Brazil mobile)
    PhoneNumber string `json:"phoneNumber" binding:"required" example:"+5511999999999"`
    
    // Email is the user's email address
    // Used for account recovery and email notifications
    // Must be unique in the system
    // Example: "joao.silva@example.com"
    Email string `json:"email" binding:"required,email" example:"joao.silva@example.com"`
    
    // ZipCode is the Brazilian postal code (CEP)
    // Format: 8 digits without separators (no hyphen)
    // Used to auto-populate address fields via CEP lookup service
    // Example: "01310100" (Avenida Paulista, São Paulo)
    ZipCode string `json:"zipCode" binding:"required,len=8" example:"01310100"`
    
    // Street is populated automatically from CEP lookup
    // Can be overridden by user if CEP service data is incorrect
    // Example: "Avenida Paulista"
    Street string `json:"street,omitempty" example:"Avenida Paulista"`
    
    // Number is the building/property number in the street
    // Required as it's not provided by CEP lookup
    // Use "S/N" for addresses without number
    // Example: "1578"
    Number string `json:"number" binding:"required" example:"1578"`
    
    // Complement provides additional address information
    // Optional field for apartment number, building name, etc.
    // Example: "Apto 501", "Bloco B"
    Complement string `json:"complement,omitempty" example:"Apto 501"`
    
    // Neighborhood is populated automatically from CEP lookup
    // Can be overridden if CEP data is incorrect
    // Example: "Bela Vista"
    Neighborhood string `json:"neighborhood,omitempty" example:"Bela Vista"`
    
    // City is populated automatically from CEP lookup
    // Example: "São Paulo"
    City string `json:"city,omitempty" example:"São Paulo"`
    
    // State is the Brazilian state (UF) code
    // Populated automatically from CEP lookup
    // Must be a valid 2-letter state code
    // Example: "SP"
    State string `json:"state,omitempty" binding:"omitempty,len=2" example:"SP"`
    
    // Password is the user's chosen password for authentication
    // Minimum 8 characters (enforced by binding tag)
    // Should contain mix of uppercase, lowercase, numbers, and symbols (enforced by service layer)
    // Never logged or exposed in responses
    // Example: "SecureP@ssw0rd"
    Password string `json:"password" binding:"required,min=8" example:"SecureP@ssw0rd"`
}
```

**Mandatory Elements for DTOs:**
1. ✅ Struct comment explaining purpose and usage context
2. ✅ Comment for EVERY field including:
   - Clear and concise description
   - Expected format (dates, phones, etc.)
   - Validation rules
   - Allowed values (enums)
   - Default values when applicable
   - Practical example using `example` tag
3. ✅ Binding tags (`binding`) with Gin validations
4. ✅ JSON tags (`json`) with camelCase names
5. ✅ `example` tags for ALL fields
6. ✅ `default` tags when there's a default value
7. ✅ `enums` tags when there's a limited set of values

---

### 8.4 Services (Core Business Logic)

**Location**: `internal/core/service/`

**Goal**: Documentation explaining business rules, orchestration, and complex behaviors.

**Template for Public Service Methods:**

```go
// CreateUser creates a new user account with full validation and setup
//
// This method orchestrates the complete user creation flow:
//  1. Validates business rules (age requirements, role-specific validations)
//  2. Checks for duplicate CPF/email
//  3. Creates user record in database
//  4. Assigns initial role and permissions
//  5. Creates validation records for email/phone
//  6. Initializes user storage folder in cloud
//  7. Logs audit trail
//
// The operation is transactional: if any step fails, all changes are rolled back.
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging. Must contain request metadata.
//   - user: UserInterface with all required fields populated. Password will be hashed before storage.
//   - roleSlug: Target role for the user ("owner", "realtor", "agency"). Determines validation rules.
//
// Returns:
//   - created: UserInterface with ID and active role populated
//   - err: Domain error with appropriate HTTP status code:
//     * 409 (Conflict) if CPF or email already exists
//     * 422 (Unprocessable) if validation fails (invalid CPF, underage user, invalid CRECI for realtors)
//     * 500 (Internal) for infrastructure failures (DB, cloud storage, external services)
//
// Business Rules:
//   - Users must be at least 18 years old
//   - CPF must be unique and pass checksum validation
//   - Email must be unique across all users
//   - Phone must be in E.164 format
//   - Realtors MUST provide valid CRECI number and state
//   - Password must meet complexity requirements (handled by validation service)
//
// Side Effects:
//   - Creates user record in users table
//   - Creates user_roles record
//   - Creates user_validations record (email and phone unverified)
//   - Creates cloud storage folder at /users/{user_id}/
//   - Logs audit entry
//
// Example:
//   user := usermodel.NewUser()
//   user.SetFullName("João Silva")
//   user.SetNationalID("12345678901")
//   // ... set other fields
//   created, err := svc.CreateUser(ctx, user, "owner")
//   if err != nil {
//       // Handle error (already logged by service)
//       return err
//   }
//   // User created with ID: created.GetID()
func (us *userService) CreateUser(ctx context.Context, user usermodel.UserInterface, roleSlug string) (created usermodel.UserInterface, err error) {
    // Initialize tracing for distributed observability
    ctx, spanEnd, err := utils.GenerateTracer(ctx)
    if err != nil {
        return nil, derrors.Infra("Failed to generate tracer", err)
    }
    defer spanEnd()

    // Ensure logger propagation with request_id and trace_id
    ctx = utils.ContextWithLogger(ctx)
    logger := utils.LoggerFromContext(ctx)

    // Start transaction to ensure atomicity of all operations
    tx, err := us.globalService.StartTransaction(ctx)
    if err != nil {
        logger.Error("user.create_user.tx_start_error", "err", err)
        utils.SetSpanError(ctx, err)
        return nil, derrors.Infra("Failed to start transaction", err)
    }
    defer func() {
        // Rollback on error (explicit commit at the end if success)
        if err != nil {
            if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
                logger.Error("user.create_user.tx_rollback_error", "err", rbErr)
                utils.SetSpanError(ctx, rbErr)
            }
        }
    }()

    // Validate user data (CPF, email, phone, minimum age)
    if err = us.ValidateUserData(ctx, tx, user, roleSlug); err != nil {
        // Domain/validation error: just return (already mapped by ValidateUserData)
        return nil, err
    }

    // Fetch role for assignment (ensure it exists before creating user)
    role, err := us.permissionService.GetRoleBySlugWithTx(ctx, tx, roleSlug)
    if err != nil {
        // Infrastructure error: log and mark span
        logger.Error("user.create_user.get_role_error", "role_slug", roleSlug, "err", err)
        utils.SetSpanError(ctx, err)
        return nil, derrors.Infra("Failed to fetch role", err)
    }

    // Create user record in database
    if err = us.repo.CreateUser(ctx, tx, user); err != nil {
        // Repository already logged the error; just propagate
        utils.SetSpanError(ctx, err)
        return nil, derrors.Infra("Failed to create user", err)
    }

    // Assign role to newly created user
    userRole, err := us.permissionService.AssignRoleToUserWithTx(ctx, tx, user.GetID(), role.GetID(), nil, nil)
    if err != nil {
        logger.Error("user.create_user.assign_role_error", "user_id", user.GetID(), "role_id", role.GetID(), "err", err)
        utils.SetSpanError(ctx, err)
        return nil, derrors.Infra("Failed to assign role", err)
    }
    user.SetActiveRole(userRole)

    // Create validation records for email and phone (both unverified initially)
    if err = us.CreateUserValidations(ctx, tx, user); err != nil {
        logger.Error("user.create_user.create_validations_error", "user_id", user.GetID(), "err", err)
        utils.SetSpanError(ctx, err)
        return nil, derrors.Infra("Failed to create user validations", err)
    }

    // Create user folder in cloud storage (operation outside DB transaction)
    // Note: if it fails, DB transaction will be reverted; orphaned folder cleaned by cleanup job
    if err = us.cloudStorageService.CreateUserFolder(ctx, user.GetID()); err != nil {
        logger.Error("user.create_user.create_folder_error", "user_id", user.GetID(), "err", err)
        utils.SetSpanError(ctx, err)
        return nil, derrors.Infra("Failed to create user folder", err)
    }

    // Register audit trail for creation
    if err = us.globalService.CreateAudit(ctx, tx, globalmodel.TableUsers, "User account created", user.GetID()); err != nil {
        logger.Error("user.create_user.create_audit_error", "user_id", user.GetID(), "err", err)
        utils.SetSpanError(ctx, err)
        return nil, derrors.Infra("Failed to create audit record", err)
    }

    // Commit transaction (all DB operations successful)
    if err = us.globalService.CommitTransaction(ctx, tx); err != nil {
        logger.Error("user.create_user.tx_commit_error", "err", err)
        utils.SetSpanError(ctx, err)
        return nil, derrors.Infra("Failed to commit transaction", err)
    }

    // Success log (Info level for domain event)
    logger.Info("user.created", "user_id", user.GetID(), "role", roleSlug)

    return user, nil
}
```

**Template for Private Service Methods:**

```go
// validateUserAge verifies if user meets minimum age requirement
//
// Parameters:
//   - ctx: context for logging
//   - bornAt: user's date of birth
//
// Returns domain error (422) if user is under 18 years old
func (us *userService) validateUserAge(ctx context.Context, bornAt time.Time) error {
    // Do not start new tracer (private methods reuse context from public)
    logger := utils.LoggerFromContext(ctx)

    // Calculate age based on current date (America/Sao_Paulo timezone)
    now := time.Now().In(utils.GetDefaultLocation())
    age := now.Year() - bornAt.Year()

    // Adjust if birthday hasn't occurred this year yet
    if now.Month() < bornAt.Month() || (now.Month() == bornAt.Month() && now.Day() < bornAt.Day()) {
        age--
    }

    // Validate minimum age (business rule: 18 years)
    if age < 18 {
        logger.Warn("user.validation.underage", "age", age, "born_at", bornAt)
        return derrors.Validation("User must be at least 18 years old", map[string]string{
            "field": "bornAt",
            "age":   fmt.Sprintf("%d", age),
        })
    }

    return nil
}

// sanitizeDeviceContext validates and normalizes device identifiers
//
// Validates UUID format of device ID and trims spaces from tokens.
// Returns 400 error if device ID is invalid.
//
// Parameters:
//   - ctx: current context
//   - deviceToken: FCM token for push notifications (may contain spaces)
//   - deviceID: unique device identifier (must be valid UUIDv4)
//   - operation: calling operation name for contextual logging
//
// Returns:
//   - updated context with device ID
//   - sanitized deviceToken (trimmed)
//   - sanitized deviceID (trimmed)
//   - domain error if validation fails
func (us *userService) sanitizeDeviceContext(ctx context.Context, deviceToken string, deviceID string, operation string) (context.Context, string, string, error) {
    logger := utils.LoggerFromContext(ctx)

    // Trim whitespace from both fields
    trimmedToken := strings.TrimSpace(deviceToken)
    trimmedID := strings.TrimSpace(deviceID)

    // Validate device ID presence (required for session tracking)
    if trimmedID == "" {
        logger.Warn(operation+".missing_device_id")
        return ctx, "", "", derrors.BadRequest("X-Device-Id header is required", nil)
    }

    // Validate device ID UUID format
    if _, err := uuid.Parse(trimmedID); err != nil {
        logger.Warn(operation+".invalid_device_id", "device_id", trimmedID, "err", err)
        return ctx, "", "", derrors.BadRequest("X-Device-Id must be a valid UUIDv4", nil)
    }

    // Add device ID to context for downstream use
    ctx = context.WithValue(ctx, globalmodel.DeviceIDKey, trimmedID)

    return ctx, trimmedToken, trimmedID, nil
}
```

**Mandatory Elements for Services:**
1. ✅ Complete Godoc describing:
   - What the function does (high level)
   - Flow of operations (numbered list)
   - Parameters with types and purposes
   - Returns with types and error scenarios
   - Business rules applied
   - Side effects (DB, external services)
   - Usage example when complex
2. ✅ Inline comments explaining:
   - Why each operation is necessary
   - Business rules being applied
   - Special handling of edge cases
   - Non-obvious technical decisions
3. ✅ Logging for infrastructure errors
4. ✅ `utils.SetSpanError` on all infrastructure errors
5. ✅ Proper propagation of domain errors (no ERROR log)

---

### 8.5 Repositories (Right Adapters)

**Location**: `internal/adapter/right/mysql/`

**Goal**: Technical documentation explaining queries, edge cases, and data conversions.

**Complete Template for Repositories:**

```go
// GetUserByID retrieves a user by their unique ID from the users table.
// Returns sql.ErrNoRows if no user is found with the given ID or if the user is marked as deleted.
// This function ensures only active (non-deleted) users are returned.
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging
//   - tx: Database transaction (can be nil for standalone queries)
//   - id: User's unique identifier
//
// Returns:
//   - user: UserInterface containing all user data
//   - error: sql.ErrNoRows if not found, or other database errors
func (ua *UserAdapter) GetUserByID(ctx context.Context, tx *sql.Tx, id int64) (user usermodel.UserInterface, err error) {
    // Initialize tracing for observability (metrics + distributed tracing)
    ctx, spanEnd, err := utils.GenerateTracer(ctx)
    if err != nil {
        return nil, err
    }
    defer spanEnd()

    // Attach logger to context to ensure request_id/trace_id propagation
    ctx = utils.ContextWithLogger(ctx)
    logger := utils.LoggerFromContext(ctx)

    // Query only active users (deleted = 0) to maintain data integrity
    query := `SELECT id, full_name, nick_name, national_id, creci_number, creci_state, 
              creci_validity, born_at, phone_number, email, zip_code, street, number, 
              complement, neighborhood, city, state, password, opt_status, 
              last_activity_at, deleted
              FROM users 
              WHERE id = ? AND deleted = 0`

    // Execute query using instrumented adapter (auto-generates metrics + tracing)
    rows, queryErr := ua.QueryContext(ctx, tx, "select", query, id)
    if queryErr != nil {
        // Mark span as error for distributed tracing analysis
        utils.SetSpanError(ctx, queryErr)
        logger.Error("mysql.user.get_user_by_id.query_error", "error", queryErr)
        return nil, queryErr
    }
    defer rows.Close()

    // Convert database rows to strongly-typed entities
    entities, err := rowsToUserEntities(rows)
    if err != nil {
        utils.SetSpanError(ctx, err)
        logger.Error("mysql.user.get_user_by_id.scan_error", "error", err)
        return nil, fmt.Errorf("scan user rows: %w", err)
    }

    // Handle no results: return standard sql.ErrNoRows for service layer handling
    if len(entities) == 0 {
        return nil, sql.ErrNoRows
    }

    // Safety check: unique constraint should prevent multiple rows
    if len(entities) > 1 {
        errMultiple := fmt.Errorf("multiple users found with the same ID: %d", id)
        utils.SetSpanError(ctx, errMultiple)
        logger.Error("mysql.user.get_user_by_id.multiple_users_error", "user_id", id, "error", errMultiple)
        return nil, errMultiple
    }

    // Convert database entity to domain model (separation of concerns)
    user = userconverters.UserEntityToDomain(entities[0])

    return user, nil
}
```

**Template for Write Operations:**

```go
// UpdateUserPassword updates the password hash for a specific user
//
// This function performs a targeted update of only the password field.
// Returns sql.ErrNoRows if user not found or already deleted.
//
// Parameters:
//   - ctx: Context for tracing and logging
//   - tx: Database transaction (must not be nil for consistency)
//   - userID: ID of the user whose password will be updated
//   - hashedPassword: Bcrypt hash of the new password (must be pre-hashed)
//
// Returns:
//   - error: sql.ErrNoRows if user not found, or other database errors
//
// Security Considerations:
//   - Password must be hashed BEFORE calling this function
//   - Old password is not validated here (must be done in service layer)
//   - Operation must run within a transaction
func (ua *UserAdapter) UpdateUserPassword(ctx context.Context, tx *sql.Tx, userID int64, hashedPassword string) error {
    // Initialize tracing
    ctx, spanEnd, err := utils.GenerateTracer(ctx)
    if err != nil {
        return err
    }
    defer spanEnd()

    ctx = utils.ContextWithLogger(ctx)
    logger := utils.LoggerFromContext(ctx)

    // Query updates only password; WHERE ensures user exists and is not deleted
    query := `UPDATE users 
              SET password = ?, updated_at = NOW() 
              WHERE id = ? AND deleted = 0`

    // Execute update via instrumented adapter
    result, execErr := ua.ExecContext(ctx, tx, "update", query, hashedPassword, userID)
    if execErr != nil {
        utils.SetSpanError(ctx, execErr)
        logger.Error("mysql.user.update_password.exec_error", "user_id", userID, "error", execErr)
        return fmt.Errorf("update user password: %w", execErr)
    }

    // Check if any rows were affected (does user exist?)
    rowsAffected, raErr := result.RowsAffected()
    if raErr != nil {
        utils.SetSpanError(ctx, raErr)
        logger.Error("mysql.user.update_password.rows_affected_error", "user_id", userID, "error", raErr)
        return fmt.Errorf("get rows affected: %w", raErr)
    }

    // Return sql.ErrNoRows if user not found (service layer maps to 404)
    if rowsAffected == 0 {
        return sql.ErrNoRows
    }

    return nil
}
```

**Mandatory Elements for Repositories:**
1. ✅ Complete Godoc in English with:
   - Description of what the function does
   - sql.ErrNoRows return scenarios
   - Documented parameters and returns
   - Business rules in query (e.g., `deleted = 0`)
2. ✅ Inline comments explaining:
   - Why certain query conditions
   - Edge case handling
   - InstrumentedAdapter usage
   - Entity/domain conversions
3. ✅ `utils.GenerateTracer` at beginning
4. ✅ `utils.SetSpanError` on all errors
5. ✅ Mandatory use of `InstrumentedAdapter` (ExecContext, QueryContext, QueryRowContext)
6. ✅ Return `sql.ErrNoRows` when appropriate

---

### 8.6 Entities e Converters

**Location**: 
- Entities: `internal/adapter/right/mysql/{domain}/entities/`
- Converters: `internal/adapter/right/mysql/{domain}/converters/`

**Goal**: Clear documentation about data layer mapping.

**Template for Entities:**

```go
package userentity

import (
    "database/sql"
    "time"
)

// UserEntity represents a row from the users table in the database
//
// This struct maps directly to the database schema and uses sql.Null* types
// for nullable columns. It should ONLY be used within the MySQL adapter layer.
//
// Schema Mapping:
//   - Database: users table (InnoDB, utf8mb4_unicode_ci)
//   - Primary Key: id (BIGINT AUTO_INCREMENT)
//   - Unique Constraints: national_id, email, phone_number
//   - Indexes: idx_users_deleted, idx_users_national_id, idx_users_email
//
// NULL Handling:
//   - sql.NullString: Used for VARCHAR columns that allow NULL
//   - sql.NullTime: Used for DATETIME columns that allow NULL
//   - Direct types: Used for NOT NULL columns
//
// Conversion:
//   - To Domain: Use userconverters.UserEntityToDomain()
//   - From Domain: Use userconverters.UserDomainToEntity()
//
// Important:
//   - DO NOT use this struct outside the adapter layer
//   - DO NOT add business logic methods to this struct
//   - DO NOT import core/model packages here
type UserEntity struct {
    // ID is the user's unique identifier (PRIMARY KEY, AUTO_INCREMENT)
    ID int64
    
    // FullName is the user's complete legal name (NOT NULL, VARCHAR(100))
    FullName string
    
    // NickName is the user's display name (NULL, VARCHAR(50))
    NickName sql.NullString
    
    // NationalID is the user's CPF or CNPJ (NOT NULL, VARCHAR(14), UNIQUE)
    NationalID string
    
    // CreciNumber is the CRECI registration number (NULL, VARCHAR(20))
    CreciNumber sql.NullString
    
    // CreciState is the Brazilian state where CRECI is registered (NULL, CHAR(2))
    CreciState sql.NullString
    
    // CreciValidity is the CRECI expiration date (NULL, DATE)
    CreciValidity sql.NullTime
    
    // BornAT is the user's date of birth (NOT NULL, DATE)
    BornAT time.Time
    
    // PhoneNumber is the user's mobile in E.164 format (NOT NULL, VARCHAR(20), UNIQUE)
    PhoneNumber string
    
    // Email is the user's email address (NOT NULL, VARCHAR(100), UNIQUE)
    Email string
    
    // ZipCode is the Brazilian CEP (NOT NULL, CHAR(8))
    ZipCode string
    
    // Street is the street name (NOT NULL, VARCHAR(100))
    Street string
    
    // Number is the building number (NOT NULL, VARCHAR(10))
    Number string
    
    // Complement provides additional address info (NULL, VARCHAR(50))
    Complement sql.NullString
    
    // Neighborhood is the district/neighborhood name (NOT NULL, VARCHAR(50))
    Neighborhood string
    
    // City is the city name (NOT NULL, VARCHAR(50))
    City string
    
    // State is the Brazilian state code (NOT NULL, CHAR(2))
    State string
    
    // Photo is the user's profile photo URL (NULL, VARCHAR(255))
    Photo sql.NullString
    
    // Password is the bcrypt hash of user's password (NOT NULL, VARCHAR(255))
    Password string
    
    // OptStatus indicates if user opted in for marketing (NOT NULL, TINYINT(1))
    OptStatus bool
    
    // LastActivityAT is the timestamp of user's last action (NOT NULL, DATETIME)
    LastActivityAT time.Time
    
    // Deleted indicates soft delete status (NOT NULL, TINYINT(1), DEFAULT 0)
    Deleted bool
    
    // LastSignInAttempt is the timestamp of last sign-in attempt (NULL, DATETIME)
    LastSignInAttempt sql.NullTime
}
```

**Template for Converters:**

```go
package userconverters

import (
    "database/sql"
    
    userentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/entities"
    usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
)

// UserEntityToDomain converts a database entity to a domain model
//
// This converter handles the translation from database-specific types (sql.Null*)
// to clean domain types, ensuring the core layer remains decoupled from database concerns.
//
// Conversion Rules:
//   - sql.NullString → string (empty string if NULL)
//   - sql.NullTime → time.Time (zero time if NULL)
//   - TINYINT(1) → bool (1 = true, 0 = false)
//
// Parameters:
//   - entity: UserEntity from database query
//
// Returns:
//   - user: UserInterface with all fields populated from entity
func UserEntityToDomain(entity userentity.UserEntity) usermodel.UserInterface {
    user := usermodel.NewUser()
    
    // Map mandatory fields (NOT NULL in schema)
    user.SetID(entity.ID)
    user.SetFullName(entity.FullName)
    user.SetNationalID(entity.NationalID)
    user.SetBornAt(entity.BornAT)
    user.SetPhoneNumber(entity.PhoneNumber)
    user.SetEmail(entity.Email)
    user.SetZipCode(entity.ZipCode)
    user.SetStreet(entity.Street)
    user.SetNumber(entity.Number)
    user.SetNeighborhood(entity.Neighborhood)
    user.SetCity(entity.City)
    user.SetState(entity.State)
    user.SetPassword(entity.Password)
    user.SetOptStatus(entity.OptStatus)
    user.SetLastActivityAt(entity.LastActivityAT)
    user.SetDeleted(entity.Deleted)
    
    // Map optional fields (NULL in schema) - check Valid before accessing
    if entity.NickName.Valid {
        user.SetNickName(entity.NickName.String)
    }
    
    if entity.CreciNumber.Valid {
        user.SetCreciNumber(entity.CreciNumber.String)
    }
    
    if entity.CreciState.Valid {
        user.SetCreciState(entity.CreciState.String)
    }
    
    if entity.CreciValidity.Valid {
        user.SetCreciValidity(entity.CreciValidity.Time)
    }
    
    if entity.Complement.Valid {
        user.SetComplement(entity.Complement.String)
    }
    
    if entity.LastSignInAttempt.Valid {
        user.SetLastSignInAttempt(entity.LastSignInAttempt.Time)
    }
    
    return user
}

// UserDomainToEntity converts a domain model to a database entity
//
// This converter handles the translation from clean domain types to database-specific
// types (sql.Null*), preparing data for database insertion/update.
//
// Conversion Rules:
//   - string → sql.NullString (Valid=true if non-empty)
//   - time.Time → sql.NullTime (Valid=true if not zero time)
//   - bool → TINYINT(1) (true = 1, false = 0)
//
// Parameters:
//   - domain: UserInterface from core layer
//
// Returns:
//   - entity: UserEntity ready for database operations
//
// Important:
//   - ID may be 0 for new records (populated by AUTO_INCREMENT)
//   - Empty strings are converted to NULL for optional fields
//   - Zero times (IsZero()) are converted to NULL for optional date fields
func UserDomainToEntity(domain usermodel.UserInterface) userentity.UserEntity {
    entity := userentity.UserEntity{}
    
    // Map mandatory fields
    entity.ID = domain.GetID()
    entity.FullName = domain.GetFullName()
    entity.NationalID = domain.GetNationalID()
    entity.BornAT = domain.GetBornAt()
    entity.PhoneNumber = domain.GetPhoneNumber()
    entity.Email = domain.GetEmail()
    entity.ZipCode = domain.GetZipCode()
    entity.Street = domain.GetStreet()
    entity.Number = domain.GetNumber()
    entity.Neighborhood = domain.GetNeighborhood()
    entity.City = domain.GetCity()
    entity.State = domain.GetState()
    entity.Password = domain.GetPassword()
    entity.OptStatus = domain.IsOptStatus()
    entity.LastActivityAT = domain.GetLastActivityAt()
    entity.Deleted = domain.IsDeleted()
    
    // Map optional fields - convert to sql.Null* with Valid based on value presence
    nickName := domain.GetNickName()
    entity.NickName = sql.NullString{
        String: nickName,
        Valid:  nickName != "",
    }
    
    creciNumber := domain.GetCreciNumber()
    entity.CreciNumber = sql.NullString{
        String: creciNumber,
        Valid:  creciNumber != "",
    }
    
    creciState := domain.GetCreciState()
    entity.CreciState = sql.NullString{
        String: creciState,
        Valid:  creciState != "",
    }
    
    creciValidity := domain.GetCreciValidity()
    entity.CreciValidity = sql.NullTime{
        Time:  creciValidity,
        Valid: !creciValidity.IsZero(),
    }
    
    complement := domain.GetComplement()
    entity.Complement = sql.NullString{
        String: complement,
        Valid:  complement != "",
    }
    
    lastSignIn := domain.GetLastSignInAttempt()
    entity.LastSignInAttempt = sql.NullTime{
        Time:  lastSignIn,
        Valid: !lastSignIn.IsZero(),
    }
    
    return entity
}
```

**Mandatory Elements for Entities:**
1. ✅ Package comment explaining purpose
2. ✅ Struct comment describing:
   - Which table it represents
   - Schema details (engine, charset, constraints)
   - Usage rules (adapter layer only)
   - References to converters
3. ✅ Comment for EVERY field including:
   - Field description
   - SQL type (VARCHAR, DATETIME, etc.)
   - Constraints (NOT NULL, UNIQUE, etc.)
   - Expected format
   - Business rules when relevant

**Mandatory Elements for Converters:**
1. ✅ Godoc explaining:
   - Conversion direction (entity→domain or domain→entity)
   - Conversion rules applied
   - NULL handling
   - Parameters and returns
2. ✅ Inline comments explaining:
   - Why certain fields are mapped specially
   - Valid checks for sql.Null*
   - Edge cases in mapping

---

### 8.7 Models (Domain)

**Location**: `internal/core/model/`

**Goal**: Clear documentation of domain interfaces and structures.

**Template for Domain Interfaces:**

```go
package usermodel

import (
    "time"
    
    permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
)

// UserInterface defines the contract for user domain entities
//
// This interface represents a user in the TOQ platform's core domain.
// It provides a clean abstraction over user data, decoupling the business logic
// from storage and presentation concerns.
//
// Design Principles:
//   - Getter/Setter pattern for all fields (no direct field access)
//   - No database or HTTP dependencies
//   - Suitable for mocking in tests
//   - Implementation hidden from consumers
//
// Implementations:
//   - user (private struct in user_domain.go)
//
// Usage:
//   - Services: Orchestrate business logic using this interface
//   - Repositories: Return/accept this interface (convert from/to entities)
//   - Handlers: Convert to/from DTOs at the boundary
//
// Example:
//   user := usermodel.NewUser()
//   user.SetFullName("João Silva")
//   user.SetNationalID("12345678901")
//   // ... set other required fields
//   createdUser, err := userService.CreateUser(ctx, user)
type UserInterface interface {
    // GetID returns the user's unique identifier
    // Returns 0 for new users not yet persisted
    GetID() int64
    
    // SetID sets the user's unique identifier
    // Typically called by repository after INSERT
    SetID(int64)
    
    // GetActiveRole returns the user's currently active role
    // Returns nil if no role assigned yet
    GetActiveRole() permissionmodel.UserRoleInterface
    
    // SetActiveRole sets the user's active role
    SetActiveRole(active permissionmodel.UserRoleInterface)
    
    // GetFullName returns the user's complete legal name
    GetFullName() string
    
    // SetFullName sets the user's complete legal name
    // Must be at least 2 characters (enforced by service layer)
    SetFullName(string)
    
    // GetNickName returns the user's display name
    GetNickName() string
    
    // SetNickName sets the user's display name
    SetNickName(string)
    
    // GetNationalID returns the user's CPF or CNPJ
    // Format: digits only (no punctuation)
    GetNationalID() string
    
    // SetNationalID sets the user's CPF or CNPJ
    SetNationalID(string)
    
    // Additional methods... (keeping it concise for space)
}

// NewUser creates a new user instance with default values
//
// Returns a UserInterface implementation with:
//   - ID = 0 (not persisted yet)
//   - Deleted = false
//   - OptStatus = false (marketing opted out by default)
//   - All timestamps = zero time
//   - All strings = empty
//
// Caller must set all required fields before persistence.
func NewUser() UserInterface {
    return &user{}
}
```

**Mandatory Elements for Models:**
1. ✅ Interface comment explaining:
   - Purpose and responsibilities
   - Design principles
   - How and where to use
   - Basic usage example
2. ✅ Comment for EACH interface method including:
   - What the method returns/receives
   - Data format when applicable
   - Special rules or edge cases
   - Deprecation status when applicable
3. ✅ Constructor (NewXXX) comment explaining:
   - Initial/default values
   - What caller must do after creation
4. ✅ Private implementation struct: brief comment only

---

### 8.8 Helpers e Utils

**Location**: `internal/core/utils/`, `internal/adapter/left/http/utils/`

**Goal**: Clear documentation of reusable utility functions.

**Template for Util Functions:**

```go
package utils

import (
    "strings"
    "time"
)

// ParseRFC3339Relaxed parses RFC3339 timestamps tolerating spaces in the offset portion
//
// This function handles a common issue with URL query decoding where the '+' character
// in timezone offsets (e.g., "+03:00") is converted to a space (" 03:00") by HTTP libraries.
//
// The function normalizes the input by detecting space-separated offset suffixes and
// restoring the '+' sign before parsing with Go's standard RFC3339 parser.
//
// Parameters:
//   - field: Name of the field being parsed (used in error messages for clarity)
//   - raw: The RFC3339 string to parse (may contain space instead of '+' in offset)
//
// Returns:
//   - time.Time: Parsed timestamp with timezone information preserved
//   - error: ValidationError if input is empty or not valid RFC3339 format
//
// Supported Formats:
//   - Standard: "2024-11-06T14:30:00+03:00" ✅
//   - With space: "2024-11-06T14:30:00 03:00" ✅ (space converted to +)
//   - UTC: "2024-11-06T14:30:00Z" ✅
//   - Negative offset: "2024-11-06T14:30:00-03:00" ✅
//
// Example:
//   // URL query: ?createdAfter=2024-01-01T00:00:00 03:00
//   // ('+' decoded as space by HTTP library)
//   timestamp, err := utils.ParseRFC3339Relaxed("createdAfter", "2024-01-01T00:00:00 03:00")
//   // timestamp = 2024-01-01 00:00:00 +0300 +0300
func ParseRFC3339Relaxed(field, raw string) (time.Time, error) {
    // Validate value presence before processing
    trimmed := strings.TrimSpace(raw)
    if trimmed == "" {
        return time.Time{}, ValidationError(field, field+" is required and must follow RFC3339")
    }

    // Normalize offset (restore '+' if necessary)
    normalized := normalizeRFC3339Offset(trimmed)
    
    // Parse using standard RFC3339 format
    value, err := time.Parse(time.RFC3339, normalized)
    if err != nil {
        return time.Time{}, ValidationError(field, field+" must be a valid RFC3339 timestamp")
    }
    
    return value, nil
}

// normalizeRFC3339Offset detects and fixes space-separated timezone offsets
//
// Detects offset suffix after last space and restores '+' when appropriate.
// Returns string unchanged if offset pattern not detected.
func normalizeRFC3339Offset(value string) string {
    // Find last space (possible offset separator)
    lastSpace := strings.LastIndex(value, " ")
    if lastSpace == -1 {
        return value // No spaces, return unchanged
    }

    // Extract suffix after space
    suffix := value[lastSpace+1:]
    if !looksLikeOffsetSuffix(suffix) {
        return value // Suffix doesn't look like offset, return unchanged
    }

    prefix := value[:lastSpace]
    if len(suffix) == 0 {
        return value
    }

    // If suffix already has sign (+/-), return assembled
    if suffix[0] == '+' || suffix[0] == '-' {
        return prefix + suffix
    }

    // Suffix is numeric without sign, assume '+'
    return prefix + "+" + suffix
}

// looksLikeOffsetSuffix validates if a string matches timezone offset pattern
//
// Checks if string has "HH:MM" format with or without initial sign.
// Returns true if it looks like a timezone offset.
func looksLikeOffsetSuffix(s string) bool {
    // Offset must be exactly 5 characters: XX:XX
    if len(s) != 5 {
        return false
    }
    
    // Third character must be ':'
    if s[2] != ':' {
        return false
    }
    
    // First character can be sign or digit
    first := s[0]
    if first == '+' || first == '-' {
        return true
    }
    
    // If digit, validate (0-9)
    return first >= '0' && first <= '9'
}
```

**Mandatory Elements for Utils:**
1. ✅ Complete Godoc in English explaining:
   - Function purpose
   - Why it exists (problem it solves)
   - Validation/transformation rules
   - Parameters and returns
   - Usage examples with success and failure cases
2. ✅ Inline comments explaining:
   - Implementation logic
   - Edge cases handled
   - Reasons for specific technical choices

---

### 8.9 Factories

**Location**: `internal/core/factory/`

**Goal**: Clear documentation of creation patterns and dependency injection.

**Template for Factory Interface:**

```go
package factory

// AdapterFactory defines the main interface for dependency creation and injection
//
// This interface implements the Abstract Factory pattern to centralize and organize
// the creation of all adapters, services, and handlers in the application.
//
// Design Goals:
//   - Centralize dependency creation logic
//   - Enable easy mocking for tests
//   - Enforce consistent initialization patterns
//   - Support lifecycle management through LifecycleManager
//
// Initialization Order:
//   1. CreateMetricsAdapter (standalone, no dependencies)
//   2. CreateStorageAdapters (database, cache)
//   3. CreateRepositoryAdapters (depends on storage)
//   4. CreateValidationAdapters (external APIs: CEP, CPF, CNPJ)
//   5. CreateExternalServiceAdapters (FCM, Email, SMS, Cloud Storage)
//   6. CreateHTTPHandlers (depends on all services)
//
// Lifecycle Management:
//   - Adapters can register cleanup functions via LifecycleManager
//   - Cleanup functions are called in reverse order during shutdown
//   - Example: Database connections closed before metrics shutdown
//
// Usage:
//   factory := factory.NewAdapterFactory(lifecycleManager)
//   storage, err := factory.CreateStorageAdapters(ctx, env, db, metrics)
//   repos, err := factory.CreateRepositoryAdapters(storage.Database, metrics)
type AdapterFactory interface {
    // CreateValidationAdapters creates all validation adapters for external services
    //
    // These adapters validate business data using external APIs:
    //   - CEP: Brazilian postal code validation and address lookup
    //   - CPF: Brazilian individual taxpayer ID validation
    //   - CNPJ: Brazilian company taxpayer ID validation
    //
    // Parameters:
    //   - env: Environment configuration with API credentials and endpoints
    //
    // Returns:
    //   - ValidationAdapters: Struct containing all validation adapter instances
    //   - error: Configuration or initialization errors
    CreateValidationAdapters(env *globalmodel.Environment) (ValidationAdapters, error)

    // CreateExternalServiceAdapters creates all external service adapters
    //
    // These adapters integrate with third-party services:
    //   - FCM: Firebase Cloud Messaging for push notifications
    //   - Email: Email delivery service
    //   - SMS: SMS delivery for OTPs and notifications
    //   - GCS: Google Cloud Storage for file uploads
    //
    // Parameters:
    //   - ctx: Context for initialization and connection setup
    //   - env: Environment configuration with service credentials
    //
    // Returns:
    //   - ExternalServiceAdapters: Struct containing all service adapter instances
    //   - error: Authentication, connectivity, or configuration errors
    CreateExternalServiceAdapters(ctx context.Context, env *globalmodel.Environment) (ExternalServiceAdapters, error)

    // Additional methods...
}

// NewAdapterFactory creates a new concrete adapter factory instance
//
// Parameters:
//   - lm: LifecycleManager for registering cleanup functions
//
// Returns:
//   - AdapterFactory: Factory instance ready for component creation
func NewAdapterFactory(lm LifecycleManager) AdapterFactory {
    return &ConcreteAdapterFactory{
        lm: lm,
    }
}
```

**Mandatory Elements for Factories:**
1. ✅ Interface comment explaining:
   - Pattern implemented (Abstract Factory)
   - Design goals
   - Initialization order
   - Lifecycle management
2. ✅ Comment for EACH factory method including:
   - What is created
   - Required dependencies
   - Parameters and returns
   - Possible errors
3. ✅ Constructor (NewXXX) comment

---

### 8.10 Checklist de Documentação

Use this checklist when creating or reviewing code:

**Handlers:**
- [ ] Godoc with `@Summary`, `@Description`, `@Tags`
- [ ] All `@Param` documented with examples
- [ ] All `@Failure` codes documented with scenarios
- [ ] `@Router` with correct path and method
- [ ] Inline comments explaining non-obvious logic

**DTOs:**
- [ ] Struct comment explaining purpose
- [ ] ALL fields documented with description, format, examples
- [ ] `example` tags on all fields
- [ ] `enums` tags when applicable
- [ ] `default` tags when applicable

**Services:**
- [ ] Godoc explaining: what it does, flow, business rules, side effects
- [ ] Documented parameters and returns
- [ ] Usage example when complex
- [ ] Inline comments explaining logic
- [ ] `utils.GenerateTracer` at beginning
- [ ] `utils.SetSpanError` on infrastructure errors

**Repositories:**
- [ ] Godoc explaining query and edge cases
- [ ] Documented parameters and returns
- [ ] Comments explaining: tracing, logging, query logic, conversions
- [ ] Use of `InstrumentedAdapter` (ExecContext/QueryContext)
- [ ] Return `sql.ErrNoRows` when appropriate

**Entities:**
- [ ] Struct comment with schema details
- [ ] ALL fields documented with SQL type, constraints
- [ ] References to converters

**Converters:**
- [ ] Godoc explaining conversion direction and rules
- [ ] Inline comments about NULL handling

**Models:**
- [ ] Interface comment explaining purpose and usage
- [ ] ALL methods documented
- [ ] Constructors (NewXXX) documented

**Utils:**
- [ ] Complete Godoc with purpose, rules, parameters, returns, examples
- [ ] Inline comments explaining logic

**Factories:**
- [ ] Interface comment with pattern and goals
- [ ] All factory methods documented
- [ ] Initialization order documented

---

### 8.11 Ferramentas e Automação

**Swagger Generation:**
```bash
# Generate Swagger documentation from code comments
make swagger

# Never edit swagger.json or swagger.yaml manually
# All documentation must come from code annotations
```

**Documentation Validation:**
```bash
# Check if all public exports have documentation
go vet ./...

# Specific lint for documentation (if configured)
golangci-lint run --enable=godot,godox
```

---

## 9. Marcação de erro no span

- Sempre que uma falha de infraestrutura ocorrer (DB, transação, provider, cache, IO), chame:

```go
utils.SetSpanError(ctx, err)
```

- Em handlers, não é necessário: `SendHTTPErrorObj` já marca o span automaticamente.

## 10. Nomenclatura e campos dos logs

- Nomes de eventos (exemplos):
  - permission.role.created | permission.role.assigned | permission.permission.granted | permission.http.check.denied | permission.user.blocked
  - user.auth.signin | user.auth.signout | user.auth.refresh.ok | user.auth.refresh.reuse_detected
  - session.created | session.rotated | session.revoked
  - listing.created | listing.updated | listing.deleted | listing.fetched
  - complex.created | complex.updated | complex.deleted | complex.fetched
  - user.confirm_email_change.stage_error | user.confirm_phone_change.tx_commit_error
- Campos em snake_case e objetivos: `user_id`, `role`, `stage`, `err`.

## 11. Propagação de erros (detalhado)

- Services:
  - Prefira erros de domínio com `internal/core/derrors` (Kinds/sentinelas). Para erros de infraestrutura, use `derrors.Infra(...)` e registre `slog.Error` + `utils.SetSpanError` no ponto de falha.
- Repositories:
  - Logam falhas com `slog.Error` (contexto mínimo) e retornam erros puros (`error`, por ex. `sql.ErrNoRows`). Não usar pacotes HTTP.
- Handlers:
  - Sempre responder via `internal/adapter/left/http/http_errors.SendHTTPErrorObj`, que serializa `{code,message,details}` e marca o span ativo em erros.

## 12. Organização de interfaces e arquivos

- Ports (interfaces) ficam em `internal/core/port/...` e são separadas por contexto (left/right) e por módulo (authhandler, userhandler, repositories, etc.).
- Interface em arquivo distinto dos domínios: não misture a definição de interface com os modelos de domínio; mantenha os modelos em `internal/core/model` e interfaces em `internal/core/port`.
- Services estarão no diretório service e possuem um diretório por módulo (ex.: user_service, permission_service, listing_service).
- Cada service terá seu arquivo de interface com o nome de seu módulo (ex.: user_service.go) apenas com struct, interface e func New e **cada método público estará em um arquivo separado** (ex.: create_user.go, update_user.go) - OBRIGATÓRIO.
- Cada função exposta relevante deve ter seu próprio arquivo no Service para granularidade e histórico limpo (ex.: `confirm_email_change.go`, `confirm_phone_change.go`) - OBRIGATÓRIO.
- Handlers estarão no diretório http/handlers e possuirão um diretório por módulo (ex.: user_handlers, admin_handlers, auth_handlers).
- Cada handler estará em um arquivo separado (ex.: create_user_handler.go, update_user_handler.go) apenas com struct, interface e func New e **cada método público estará em um arquivo separado** - OBRIGATÓRIO.
- Em adapters, **cada método público deve estar em arquivo separado** (ex.: `create_user.go`, `update_user_role_status.go`, `block_user_temporarily.go`) - OBRIGATÓRIO.

## 13. Padrão para análise/refatoração

- Ao refatorar:
  - Não mude contratos públicos sem atualizar Swagger/Docs.
  - Preserve o isolamento das camadas: adapters sem semântica de domínio nem HTTP; services sem HTTP; handlers sem spans.
  - Padronize conversões em Converters; elimine conversões “inline” dentro de queries.
  - Tracing/logs conforme seções acima; evite duplicação de spans.
  - Sempre rodar build/lint

Checklist rápido de refatoração:
- [ ] Mantém arquitetura hexagonal e fluxo Handlers → Services → Repos.
- [ ] Repositórios usam Converters; retornam `sql.ErrNoRows` quando aplicável.
- [ ] Adapters MySQL utilizam `InstrumentedAdapter` (`ExecContext`/`QueryContext`/`QueryRowContext`/`PrepareContext`) sem recriar helpers `basic_*`.
- [ ] Services mapeiam erros de infra para domínio; marcam span e logam infra.
- [ ] Handlers usam DTOs e `SendHTTPErrorObj`; Swagger atualizado.
- [ ] Transações via serviço/infra padrão; commit/rollback corretos.
- [ ] Logs/Spans seguindo este guia.

## 14. Anti‑padrões a evitar

- Criar spans em handlers HTTP (duplicação com o middleware).
- Logar erros de domínio como `slog.Error` (use domínio como retorno, sem ERROR).
- Serializar respostas de erro manualmente no handler (sempre use `SendHTTPErrorObj`).
- Mapear erros de repositório para HTTP dentro do adapter.
- Fazer logs verbosos nos adapters sem necessidade.

## 15. Checklist de PR (observabilidade)

- [ ] Services públicos iniciam tracer e finalizam com `defer`.
- [ ] Handlers não criam spans e usam `SendHTTPErrorObj`.
- [ ] Adapters retornam erros “puros” e evitam verbosidade.
- [ ] Erros de infra possuem `slog.Error` no ponto da falha e `SetSpanError`.
- [ ] Erros de domínio são propagados sem `slog.Error`.
- [ ] Logs usam campos em snake_case e mensagens curtas.

## 16. Exemplos rápidos

### 16.1 Service (público) — erro de infra vs domínio:

```go
if err := repo.UpdateUser(...); err != nil {
    utils.SetSpanError(ctx, err)
    slog.Error("user.update.stage_error", "stage", "update_user", "err", err)
    return derrors.Infra("Failed to update user", err)
}
// Domínio: apenas retorne a sentinela/Kind
return derrors.ErrPhoneChangeNotPending
```

### 16.2 Repository — erro puro e log enxuto:

```go
res, err := tx.ExecContext(ctx, q, args...)
if err != nil {
    slog.Error("mysql.user.update: exec_failed", "err", err)
    return err
}
```

### 16.3 Handler — serialização padronizada:

```go
if err != nil {
    httperrors.SendHTTPErrorObj(c, err)
    return
}
```

## 17. Referências

- `internal/adapter/left/http/http_errors` — serialização de erros para HTTP.
- `internal/core/utils` — tracing (`GenerateTracer`, `SetSpanError`).
- `internal/core/derrors` — erros de domínio (Kind/sentinelas).
- `docs/observability/logs.md` — guia de Loki/Grafana e correlação de logs.
