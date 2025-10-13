# toq_server
TOQ Server is an HTTP API server for the TOQ App, built with Go, Gin, and a hexagonal architecture. It exposes REST endpoints under `/api/v2`, with centralized error handling, tracing, metrics, and clean DI via factories.

Developer docs:
- See `docs/toq_server_go_guide.md` — Global Go developer guide (architecture, DI, repos/DTOs/handlers, logging, tracing, errors, transactions, checklists).
- S3 credentials now depend on the AWS credential chain (instance profile, environment variables). Do not commit static keys; use `TOQ_S3_*` environment variables only for local overrides.

## Execução em duas instâncias (nohup + F5)
- **Instância principal (nohup)**: execute `nohup ./bin/toq_server &` sem variáveis extras. O servidor sobe com `ENVIRONMENT=homo`, porta `:8080`, workers em execução e telemetria apontando para o collector de homologação.
- **Instância de debug (VS Code / F5)**: configure o launch `TOQ Server (Development)` para incluir `ENVIRONMENT=dev` (já presente em `.vscode/launch.json`). Essa instância usa a porta `127.0.0.1:18080`, mantém os workers desativados e envia telemetria marcada como ambiente dev.
- **Sobrescrita manual de porta**: defina `TOQ_HTTP_PORT` antes de iniciar o binário para forçar uma porta específica sem alterar o YAML.
- **Observabilidade**: ambas as instâncias continuam exportando métricas via `/metrics`; utilize o label `deployment.environment` nas consultas para diferenciar `homo` de `dev`.

## API path conventions
- Base path: `/api/v2`
- Email change: `/user/email/{request|confirm|resend}`
- Phone change: `/user/phone/{request|confirm|resend}`
- Password change: `/auth/password/{request|confirm}`
- Auth validation: `/auth/validate/{cpf|cnpj|cep}` (signed requests)

## Auth validation with shared HMAC

- Configuration lives under `security.hmac` in `configs/env.yaml`:

  ```yaml
  security:
    hmac:
      secret: "changeme-frontend-shared-secret" # shared with trusted clients only
      algorithm: "SHA256"                       # currently only SHA256 is accepted
      encoding: "hex"                           # hex or base64
      skew_seconds: 300                          # max clock drift allowed (5 minutes)
  ```

- Every request to `/api/v2/auth/validate/{cpf|cnpj|cep}` must include:
  - Body JSON with the domain fields (`nationalID` + optional `bornAt` or `postalCode`), `timestamp` (Unix seconds) and `hmac`.
  - Canonical string for the signature: `METHOD|PATH|timestamp|payload`.
    - `METHOD`: uppercase HTTP verb (e.g. `POST`).
    - `PATH`: the exact route, including the `/api/v2` prefix.
    - `timestamp`: same integer sent in the body.
    - `payload`: JSON minified, **without the `hmac` field**, keys sorted alphabetically.
  - Example canonical message for CPF validation:

    ```text
    POST|/api/v2/auth/validate/cpf|1712345678|{"bornAt":"1990-01-01","nationalID":"12345678901","timestamp":1712345678}
    ```

- Compute the digest with `HMAC(secret, canonical_message)` and encode using the configured `encoding` (`hex` default). Place the resulting value in the request body (`hmac`).
- The server enforces the `timestamp` drift (`skew_seconds`) and returns HTTP 401 if the signature is missing, malformed, expired or mismatched.
- Successful responses:
  - CPF/CNPJ: `{ "valid": true }` only.
  - CEP: `{ "valid": true, "postalCode": "...", "street": "...", ... }` with full address payload.
- All errors keep the standardized error schema documented below (HTTP 4xx/5xx only, never 2xx in failure scenarios).

Note: paths intentionally do not include a `/change` segment. Keep Swagger annotations and clients aligned to these routes to avoid 404s.

## Content Security Policy (CSP)
- A atualização das diretivas não acontece mais via endpoints administrativos. Use o arquivo `configs/security/csp_policy.json` como fonte da verdade.
- O time de frontend deve gerar o JSON seguindo o modelo descrito em `docs/security/csp-policy-model.md` e submetê-lo para revisão do time de plataforma.
- Após merge, o pipeline executa `scripts/render_csp_snippets.sh` para converter o JSON em snippets Nginx e aplicar a nova política.

## Error schema (standardized)
All error responses follow a flat, consistent schema returned by the centralized HTTP serializer:

```
{
  "code": number,          // HTTP status code
  "message": string,       // Human-readable message
  "details": object?       // Optional structured details (free-form JSON)
}
```

This applies to all 4xx/5xx responses, including validation and authorization failures.

## 📧 Email Configuration

The server now supports robust email delivery with the following features:

### Configuration (env.yaml)
```yaml
email:
  smtp_server: "smtp.gmail.com"     # SMTP server address
  smtp_port: 587                    # SMTP port (587 for TLS, 465 for SSL)
  smtp_user: "your-email@gmail.com" # SMTP username
  smtp_password: "your-app-password" # SMTP password or app password
  use_tls: true                     # Enable TLS (recommended)
  use_ssl: false                    # Enable SSL (alternative to TLS)
  skip_verify: false                # Skip TLS certificate verification (development only)
  from_email: "noreply@yourapp.com" # From email address
  from_name: "Your App Name"        # From display name
  max_retries: 3                    # Maximum retry attempts
  timeout_seconds: 30               # Connection timeout
```

### Features
- **Retry Logic**: Automatic retry with exponential backoff
- **TLS/SSL Support**: Secure email transmission
- **Configurable Headers**: Dynamic From address and name
- **Robust Error Handling**: Detailed logging and error reporting
- **Performance Optimized**: Connection reuse and timeout management

### Gmail Configuration for AWS
1. Enable 2-factor authentication in your Gmail account
2. Generate an App Password: Google Account → Security → App Passwords
3. Use the generated password in `smtp_password` field
4. Ensure `use_tls: true` and `smtp_port: 587`
