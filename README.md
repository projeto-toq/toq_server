# toq_server
TOQ Server is an HTTP API server for the TOQ App, built with Go, Gin, and a hexagonal architecture. It exposes REST endpoints under `/api/v2`, with centralized error handling, tracing, metrics, and clean DI via factories.

Developer docs:
- See `docs/toq_server_go_guide.md` — Global Go developer guide (architecture, DI, repos/DTOs/handlers, logging, tracing, errors, transactions, checklists).

## API path conventions
- Base path: `/api/v2`
- Email change: `/user/email/{request|confirm|resend}`
- Phone change: `/user/phone/{request|confirm|resend}`
- Password change: `/auth/password/{request|confirm}`

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
