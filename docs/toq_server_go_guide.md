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
  "RoleSlug": "<slug>" // novo campo aditivo (root|owner|realtor|agency), pode faltar em tokens antigos
}
```
Refresh token (claim `infos`):
```
{
  "ID": <int>,
  "UserRoleID": 0,
  "RoleStatus": <int>,
  "RoleSlug": "<slug|opcional>"
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

## 4. Injeção de dependências (Factory Pattern)

- Orquestração em `config.InjectDependencies(lm)` com `AdapterFactory`.
- Abstract Factory (`internal/core/factory`):
  - CreateStorageAdapters(ctx, env, db) → Database, Cache, CloseFunc.
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
- Tracing:
  - Use `utils.GenerateTracer(ctx)` no início de métodos públicos de Services, Repositories e em Workers/Go routines.
  - Em Handlers HTTP, NÃO crie spans: o `TelemetryMiddleware` cuida disso.
  - Sempre `defer spanEnd()`.
  - Em erros de infraestrutura, chame `utils.SetSpanError(ctx, err)`.
- Logging (slog):
  - `slog.Info`: eventos esperados de domínio (mudança de status, criação de recurso).
  - `slog.Warn`: anomalias/limites, condições não fatais (ex.: 429/423).
  - `slog.Error`: exclusivamente falhas de infraestrutura (DB, cache, providers externos, transações).
  - Em Repositórios, evite logs verbosos; sucesso no máximo `DEBUG` quando realmente necessário.
- Erros e HTTP:
  - Repositórios retornam erros “puros” (ex.: `sql.ErrNoRows`, `fmt.Errorf`), sem pacotes HTTP.
  - Services mapeiam para erros de domínio (sentinelas/Kind) quando aplicável; erros de infra geram logs + span error.
  - Handlers SEMPRE usam `http_errors.SendHTTPErrorObj(c, err)`; não serialize manualmente.

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

- Inicie tracer no topo; trate domínio vs infraestrutura; marque spans em erros de infra; não serialize HTTP.

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

- Inicie tracer no início do método (adapters têm tracing mínimo):

```go
ctx, end, _ := utils.GenerateTracer(ctx)
 defer end()
```

- Em falhas, `slog.Error` com contexto enxuto (ids, chave da operação) e retorne erro puro (`error`).
- Nunca use pacotes HTTP nesta camada; não mapeie para domínio.
- Sucesso: no máximo `slog.Debug` quando necessário.

Padrão de repositórios com Converters:
- Conversões entre linhas/DTOs DB ↔ entidades de domínio não devem ser feitas inline na função do repositório.
- Centralize conversões em pacotes utilitários ou helpers por domínio (ex.: `internal/core/utils/converters` ou `internal/adapter/right/mysql/<entity>/converters.go`).
- Repositório foca em: construir query, executar, checar `RowsAffected`, lidar com `sql.ErrNoRows`, e retornar entidades/domínio já convertidas.
- Em `RowsAffected == 0`, retorne `sql.ErrNoRows`. O mapeamento para “não encontrado/sem pendência” é feito no Service.

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
- Swagger: documente via anotações no código do handler/DTO.

### 7.5 Workers/Go routines

- Inicie tracer com `utils.GenerateTracer(ctx)` e `defer spanEnd()`.
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

- Services (novo padrão recomendado):
  - Prefira erros de domínio com `internal/core/derrors` (Kinds/sentinelas). Para erros de infraestrutura, use `derrors.Infra(...)` e registre `slog.Error` + `utils.SetSpanError` no ponto de falha.
- Camada legada (compatível):
  - Criar erros de negócio com `utils.NewHTTPErrorWithSource(...)` ou envolver com `utils.WrapDomainErrorWithSource(...)` para capturar origem (função, arquivo, linha, short stack).
- Repositories:
  - Logam falhas com `slog.Error` (contexto mínimo) e retornam erros puros (`error`, por ex. `sql.ErrNoRows`). Não usar pacotes HTTP.
- Handlers:
  - Sempre responder via `internal/adapter/left/http/http_errors.SendHTTPErrorObj`, que serializa `{code,message,details}` e marca o span ativo em erros.

## 11. Organização de interfaces e arquivos

- Ports (interfaces) ficam em `internal/core/port/...` e são separadas por contexto (left/right) e por módulo (authhandler, userhandler, repositories, etc.).
- Interface em arquivo distinto dos domínios: não misture a definição de interface com os modelos de domínio; mantenha os modelos em `internal/core/model` e interfaces em `internal/core/port`.
- Cada função exposta relevante pode ter seu próprio arquivo no Service para granularidade e histórico limpo (ex.: `confirm_email_change.go`, `confirm_phone_change.go`). Isso facilita teste, revisão e rastreabilidade.
- Em adapters/handlers, prefira também quebrar por caso de uso quando crescer (ex.: `update_user_role_status_tx.go`).

## 12. Padrão para análise/refatoração

- Siga os boilerplates: `prompt analise.md`, `prompt bug.md`, `prompt quick.md`.
- Ao refatorar:
  - Não mude contratos públicos sem atualizar Swagger/Docs.
  - Preserve o isolamento das camadas: adapters sem semântica de domínio nem HTTP; services sem HTTP; handlers sem spans.
  - Padronize conversões em Converters; elimine conversões “inline” dentro de queries.
  - Tracing/logs conforme seções acima; evite duplicação de spans.
  - Sempre rodar build/lint/tests rápidos e relatar PASS/FAIL.

Checklist rápido de refatoração:
- [ ] Mantém arquitetura hexagonal e fluxo Handlers → Services → Repos.
- [ ] Repositórios usam Converters; retornam `sql.ErrNoRows` quando aplicável.
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

Service (público) — erro de infra vs domínio:

```go
if err := repo.UpdateUser(...); err != nil {
    utils.SetSpanError(ctx, err)
    slog.Error("user.update.stage_error", "stage", "update_user", "err", err)
    return derrors.Infra("Failed to update user", err)
}
// Domínio: apenas retorne a sentinela/Kind
return derrors.ErrPhoneChangeNotPending
```

Repository — erro puro e log enxuto:

```go
res, err := tx.ExecContext(ctx, q, args...)
if err != nil {
    slog.Error("mysql.user.update: exec_failed", "err", err)
    return err
}
```

Handler — serialização padronizada:

```go
if err != nil {
    httperrors.SendHTTPErrorObj(c, err)
    return
}
```

## 16. Referências

- `internal/adapter/left/http/http_errors` — serialização de erros para HTTP.
- `internal/core/utils` — tracing (`GenerateTracer`, `SetSpanError`).
- `internal/core/derrors` — erros de domínio (Kind/sentinelas).
