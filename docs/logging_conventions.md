# HTTP logging and error conventions

This guide defines how we log and propagate errors across handlers → services → repositories.

## Error propagation

- Services should create business errors with `utils.NewHTTPErrorWithSource(...)` or wrap with `utils.WrapDomainErrorWithSource(...)` to capture the origin (function, file, line, short stack).
- Repositories must log failures at the source with `slog.Error` and return pure Go errors. Services map them to domain errors as needed.
- Handlers must respond via `internal/adapter/left/http/http_errors.SendHTTPErrorObj`, which also marks the active span on errors.

## Structured logging policy

- Middlewares: `RequestIDMiddleware` → `StructuredLoggingMiddleware` → `CORSMiddleware` → `TelemetryMiddleware` → `ErrorRecoveryMiddleware` → `DeviceContextMiddleware`.
- Severity mapping for HTTP requests in StructuredLoggingMiddleware:
  - 5xx: ERROR (stderr). Includes error details and short stack when available.
  - 429/423: WARN (stderr). No stack; expected throttling/lock conditions.
  - Other 4xx: INFO (stdout). Business/client issues without stack.
  - 2xx/3xx: INFO (stdout).
- Logged fields (snake_case): `request_id`, `trace_id`, `span_id`, `method`, `path`, `status`, `duration`, `size`, `client_ip`, `user_agent`. Optionally enrich with `user_id`, `user_role_id`.

## Event naming examples

- permission.role.created | permission.role.assigned | permission.permission.granted | permission.http.check.denied | permission.user.blocked
- user.auth.signin | user.auth.signout | user.auth.refresh.ok | user.auth.refresh.reuse_detected
- session.created | session.rotated | session.revoked
- listing.created | listing.updated | listing.deleted | listing.fetched
- complex.created | complex.updated | complex.deleted | complex.fetched

## Correlation and tracing

- The logging middleware includes `trace_id` and `span_id` from OpenTelemetry when available to correlate logs and traces.

## Handlers requirement

- Always use `http_errors.SendHTTPErrorObj(c, err)` for error responses. Never build error payloads manually.

## Tratamento de Erros e Observabilidade (Guia para Devs)

2. Tratamento de Erros e Observabilidade

- Tracing
  - Inicie o tracing por operação com `utils.GenerateTracer(ctx)` no início de cada método público de Services e em Workers/Go routines.
  - Em Handlers HTTP, o tracing já é iniciado pelo `TelemetryMiddleware`. Não crie spans duplicados via `GenerateTracer` no handler.
  - Sempre chame a função de finalização retornada por `GenerateTracer` (ex.: `defer spanEnd()`). Erros devem marcar o span via `utils.SetSpanError` — nos handlers isso já é feito por `SendHTTPErrorObj` e no caso de panics pelo `ErrorRecoveryMiddleware`.

- Logging
  - Logs de domínio e segurança: use apenas `slog`.
    - `slog.Info`: eventos esperados do domínio (ex.: user status mudou de pending para active).
    - `slog.Warn`: condições anômalas, indícios de fraude/reuso, limites atingidos, falhas não fatais (ex.: 429/423 por throttling/lock).
    - `slog.Error`: exclusivamente para falhas internas de infraestrutura (DB, transação, providers externos). Devem ser registrados no ponto de ocorrência.
  - Repositórios (adapters): evite logs excessivos. Em falhas críticas de infraestrutura, logue com `slog.Error` incluindo somente contexto mínimo e útil (ex.: `user_id`, `key_query`). Sucessos devem ser no máximo `DEBUG` quando realmente necessário.
  - Handlers não devem gerar logs de acesso; o `StructuredLoggingMiddleware` já o faz centralmente com severidade baseada no status HTTP (5xx→ERROR, 429/423→WARN, demais 4xx→INFO, 2xx/3xx→INFO).

- Tratamento de Erros
  - Repositórios (Adapters): retornam erros "puros" (`error`). Nunca usar pacotes HTTP (`net/http` ou `http_errors`) nesta camada.
  - Serviços (Core): propagar erros de domínio usando `utils.WrapDomainErrorWithSource(derr)` para preservar a origem (função/arquivo/linha). Ao criar novos erros de domínio, usar `utils.NewHTTPErrorWithSource(...)`. Mapear erros de repositório para erros de domínio quando aplicável. Não serializar respostas HTTP aqui.
  - Handlers (HTTP): usar `http_errors.SendHTTPErrorObj(c, err)` para converter qualquer erro propagado em JSON `{code, message, details}`. O helper também executa `c.Error(err)` para que o middleware de log capte a origem/detalhes e marca o span no trace.
