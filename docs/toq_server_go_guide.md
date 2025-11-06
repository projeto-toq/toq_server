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
- 8. Marcação de erro no span
- 9. Nomenclatura e campos dos logs
- 10. Propagação de erros (detalhado)
- 11. Organização de interfaces e arquivos
- 12. Padrão para análise/refatoração
- 13. Anti‑padrões a evitar
- 14. Checklist de PR (observabilidade)
- 15. Exemplos rápidos
- 16. Referências

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
    - service/ — Services; um arquivo por caso de uso/método público.
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

**Checklist de Conformidade Arquitetural:**

- [ ] Cada Port em `/port/right/repository/{DOMAIN}/` tem Adapter em `/adapter/right/mysql/{DOMAIN}/`
- [ ] Arquivo principal do adapter contém APENAS struct + NewFunc
- [ ] Cada método público está em arquivo próprio
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

Todas as funções públicas de repositórios DEVEM possuir documentação interna robusta em inglês, explicando detalhadamente a lógica para facilitar futuras manutenções. A documentação deve incluir:

1. **Godoc Comment (Público):** Descrição concisa do propósito da função.
2. **Comentários Inline (Detalhados):** Explicação de cada etapa crítica da implementação.

**Pontos Obrigatórios da Documentação:**
- ✅ Godoc completo com descrição, parâmetros e retornos
- ✅ Comentário explicando inicialização de tracing
- ✅ Comentário sobre logger context propagation
- ✅ Explicação de regras de negócio na query (ex: `deleted = 0`)
- ✅ Comentário sobre uso de InstrumentedAdapter
- ✅ Explicação de tratamento de erros (span marking, logging)
- ✅ Justificativa para retorno de `sql.ErrNoRows`
- ✅ Explicação de edge cases (múltiplos resultados, dados vazios)
- ✅ Comentário sobre conversão domain/entity

Ver **Seção 15 (Exemplos Rápidos)** para template completo de documentação.

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

## 8. Marcação de erro no span

- Sempre que uma falha de infraestrutura ocorrer (DB, transação, provider, cache, IO), chame:

```go
utils.SetSpanError(ctx, err)
```

- Em handlers, não é necessário: `SendHTTPErrorObj` já marca o span automaticamente.

## 9. Nomenclatura e campos dos logs

- Nomes de eventos (exemplos):
  - permission.role.created | permission.role.assigned | permission.permission.granted | permission.http.check.denied | permission.user.blocked
  - user.auth.signin | user.auth.signout | user.auth.refresh.ok | user.auth.refresh.reuse_detected
  - session.created | session.rotated | session.revoked
  - listing.created | listing.updated | listing.deleted | listing.fetched
  - complex.created | complex.updated | complex.deleted | complex.fetched
  - user.confirm_email_change.stage_error | user.confirm_phone_change.tx_commit_error
- Campos em snake_case e objetivos: `user_id`, `role`, `stage`, `err`.

## 10. Propagação de erros (detalhado)

- Services:
  - Prefira erros de domínio com `internal/core/derrors` (Kinds/sentinelas). Para erros de infraestrutura, use `derrors.Infra(...)` e registre `slog.Error` + `utils.SetSpanError` no ponto de falha.
- Repositories:
  - Logam falhas com `slog.Error` (contexto mínimo) e retornam erros puros (`error`, por ex. `sql.ErrNoRows`). Não usar pacotes HTTP.
- Handlers:
  - Sempre responder via `internal/adapter/left/http/http_errors.SendHTTPErrorObj`, que serializa `{code,message,details}` e marca o span ativo em erros.

## 11. Organização de interfaces e arquivos

- Ports (interfaces) ficam em `internal/core/port/...` e são separadas por contexto (left/right) e por módulo (authhandler, userhandler, repositories, etc.).
- Interface em arquivo distinto dos domínios: não misture a definição de interface com os modelos de domínio; mantenha os modelos em `internal/core/model` e interfaces em `internal/core/port`.
- Services estarão no diretório service e possuem um diretório por módulo (ex.: user_service, permission_service, listing_service).
- Cada service terá seu arquivo de interface com o nome de seu módulo (ex.: user_service.go) apenas com struct, interface e func New e cada método público estará em um arquivo separado (ex.: create_user.go, update_user.go).
- Cada função exposta relevante deve ter seu próprio arquivo no Service para granularidade e histórico limpo (ex.: `confirm_email_change.go`, `confirm_phone_change.go`).
- Handlers estarão no diretório http/handlers e possuirão um diretório por módulo (ex.: user_handlers, admin_handlers, auth_handlers).
- Cada handler estará em um arquivo separado (ex.: create_user_handler.go, update_user_handler.go) apenas com struct, interface e func New e cada método público estará em um arquivo separado.
- Em adapters/handlers, mantenha cada endpoint em arquivo separado (ex.: `update_user_role_status.go`).

## 12. Padrão para análise/refatoração

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

## 13. Anti‑padrões a evitar

- Criar spans em handlers HTTP (duplicação com o middleware).
- Logar erros de domínio como `slog.Error` (use domínio como retorno, sem ERROR).
- Serializar respostas de erro manualmente no handler (sempre use `SendHTTPErrorObj`).
- Mapear erros de repositório para HTTP dentro do adapter.
- Fazer logs verbosos nos adapters sem necessidade.

## 14. Checklist de PR (observabilidade)

- [ ] Services públicos iniciam tracer e finalizam com `defer`.
- [ ] Handlers não criam spans e usam `SendHTTPErrorObj`.
- [ ] Adapters retornam erros “puros” e evitam verbosidade.
- [ ] Erros de infra possuem `slog.Error` no ponto da falha e `SetSpanError`.
- [ ] Erros de domínio são propagados sem `slog.Error`.
- [ ] Logs usam campos em snake_case e mensagens curtas.

## 15. Exemplos rápidos

### 15.1 Service (público) — erro de infra vs domínio:

```go
if err := repo.UpdateUser(...); err != nil {
    utils.SetSpanError(ctx, err)
    slog.Error("user.update.stage_error", "stage", "update_user", "err", err)
    return derrors.Infra("Failed to update user", err)
}
// Domínio: apenas retorne a sentinela/Kind
return derrors.ErrPhoneChangeNotPending
```

### 15.2 Repository — erro puro e log enxuto:

```go
res, err := tx.ExecContext(ctx, q, args...)
if err != nil {
    slog.Error("mysql.user.update: exec_failed", "err", err)
    return err
}
```

### 15.3 Handler — serialização padronizada:

```go
if err != nil {
    httperrors.SendHTTPErrorObj(c, err)
    return
}
```

### 15.4 Template Completo de Documentação para Repositórios:

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
	          last_activity_at, deleted, last_signin_attempt 
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

**Benefícios da Documentação Robusta:**
- ✅ Novos desenvolvedores entendem a lógica rapidamente
- ✅ Facilita code reviews e auditoria de código
- ✅ Reduz bugs causados por mal-entendidos de lógica
- ✅ Serve como documentação viva que evolui com o código
- ✅ Ajuda em troubleshooting e debugging
- ✅ Melhora onboarding de novos membros da equipe

## 16. Referências

- `internal/adapter/left/http/http_errors` — serialização de erros para HTTP.
- `internal/core/utils` — tracing (`GenerateTracer`, `SetSpanError`).
- `internal/core/derrors` — erros de domínio (Kind/sentinelas).
- `docs/observability/logs.md` — guia de Loki/Grafana e correlação de logs.
