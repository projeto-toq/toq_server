# Documentação da cadeia de middlewares

Ordem e responsabilidades dos middlewares globais e específicos.

## Ordem de execução (global)

Request → 1. RequestID → 2. StructuredLogging → 3. CORS → 4. Telemetry → 5. ErrorRecovery → 6. DeviceContext → [Route Specific] → Handler

1) RequestIDMiddleware
- Gera `request_id` para correlação.

2) StructuredLoggingMiddleware
- Log JSON por requisição com separação stdout/stderr conforme severidade.

3) CORSMiddleware
- Configuração de headers CORS.

4) TelemetryMiddleware
- Tracing OpenTelemetry + métricas.

5) ErrorRecoveryMiddleware
- Converte panics em HTTP 500 e marca o span de erro.

6) DeviceContextMiddleware
- Injeta `device_id` do header no contexto.

## Middlewares específicos (rotas protegidas)

[Globais] → AuthMiddleware → PermissionMiddleware → Handler

- AuthMiddleware: valida JWT, injeta usuário e faz activity tracking.
- PermissionMiddleware: verifica permissões e roles.

## Notas

- O `TelemetryMiddleware` deve anteceder o `ErrorRecoveryMiddleware` para que panics marquem o span.
- Swagger e handlers devem sempre usar `http_errors.SendHTTPErrorObj` para respostas de erro.

Última atualização: 04 de Setembro de 2025
