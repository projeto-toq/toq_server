# HTTP Logging and Error Conventions (Phase 0)

This document captures the standardized conventions for structured logging and error propagation used across handlers, services, and repositories.

## Error Propagation

- Services must return DomainError for business cases via `utils.NewHTTPErrorWithSource(...)` or wrap incoming domain errors with `utils.WrapDomainErrorWithSource(...)` to capture origin.
- Infrastructure failures (DB/Redis/Tx/External providers) must be logged with `slog.Error` and returned as `utils.InternalError("")` to capture origin.
- Handlers translate any error with `internal/adapter/left/http/http_errors.SendHTTPErrorObj` which attaches the error to Gin context for centralized logging.

## Structured Logging

- Middlewares: `RequestIDMiddleware` sets `request_id` and `TelemetryMiddleware` provides trace/span context. `StructuredLoggingMiddleware` logs one entry per request.
- Severity:
  - Info: successful/expected operations
  - Warn: business anomalies (denied, limits, reuse)
  - Error: infra failures
- Fields (snake_case): `request_id`, `trace_id`, `span_id`, `method`, `path`, `status`, `duration`, `size`, `client_ip`, `user_agent`, optionally `user_id`, `user_role_id`, `role_status`.
- Error enrichment: when a `DomainErrorWithSource` is present, logs include `function`, `file`, `line`, `stack`, `error_code`, `error_message`.

## Event Naming Examples

- permission.role.created | permission.role.assigned | permission.permission.granted | permission.http.check.denied | permission.user.blocked
- user.auth.signin | user.auth.signout | user.auth.refresh.ok | user.auth.refresh.reuse_detected
- session.created | session.rotated | session.revoked
- listing.created | listing.updated | listing.deleted | listing.fetched
- complex.created | complex.updated | complex.deleted | complex.fetched

## Correlation

- The logging middleware now includes `trace_id` and `span_id` from OpenTelemetry when available to correlate logs and traces.

## Handlers Requirement

- Always use `http_errors.SendHTTPErrorObj(c, err)` for error responses. Do not build error payloads manually.
