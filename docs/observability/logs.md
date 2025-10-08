# Observabilidade de Logs com Loki + Grafana

> **Foco desta fase (PR-05)**: assegurar que toda a cadeia de observabilidade reflita o novo logger contextual (com `request_id`, `trace_id` e metadados de requisição) e que a experiência no Grafana destaque esses campos.

## Visão Geral
- **Objetivo**: Persistir e explorar logs estruturados do TOQ Server com correlação de traces e recursos de filtro por `request_id`.
- **Pipeline**: `ContextWithLogger` (handlers/middlewares/workers) → `slog` → ponte OpenTelemetry → Collector (OTLP HTTP) → Loki → Grafana.
- **Compatibilidade**: Mantém JSON/stdout existentes, adicionando labels/atributos enriquecidos (request/trace/client) para consultas.

## Componentes
| Serviço | Porta | Responsabilidade |
| --- | --- | --- |
| Loki | `3100` | Time-series DB e API de consulta para logs |
| OpenTelemetry Collector | `4317/4318` | Recebe spans/logs OTLP, agrega atributos e envia ao Loki |
| Grafana | `3000` | Dashboards (`TOQ Server Logs`) e consultas Explore |
| Jaeger | `16686` | Consulta de traces vinculados via `trace_id` |

## Variáveis de Ambiente Importantes
- `LOKI_RETENTION_DAYS`: dias de retenção do Loki. Dev = `7`, sugerido prod = `3`.
- `OTEL_EXPORTER_OTLP_ENDPOINT`: endpoint do collector (compartilhado por métricas/traces/logs).
- `OTEL_RESOURCE_ENVIRONMENT`: rótulo do ambiente (`development`, `staging`, `production`).
- `OTEL_EXPORTER_OTLP_INSECURE`: defina `false` quando houver TLS.

> Altere valores via `.env` para manter os stacks sincronizados. PR-05 não muda defaults, apenas registra expectativas.

## Passos para Subir o Stack (dev)
```bash
export OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318
export OTEL_EXPORTER_OTLP_INSECURE=true
export LOKI_RETENTION_DAYS=7
docker compose up -d loki otel-collector grafana jaeger
```
Inicie o servidor (`go run cmd/toq_server.go` ou binário). O bootstrap já conecta o `ctxlogger` aos handlers e workers.

## Painel "TOQ Server Logs"
Arquivo: `grafana/dashboard-files/toq-server-logs.json`

Principais atualizações (PR-05):
- Variável adicional **Request ID** (`$requestId`) para filtrar consultas.
- Campo derivado `request_id` exibido na tabela/logs.
- Link rápido para Jaeger via `trace_id`.

### Painéis
1. **Eventos por severidade** — Timeseries usando `count_over_time`, segmentado por `severity`.
2. **Fluxo de logs** — Painel de logs com labels comuns visíveis (`environment`, `request_id`, `trace_id`).
3. **Erros por serviço** — Tabela agregada filtrável por `request_id` via variável.

### Variáveis de Template
```text
environment: label_values({service_name="toq_server"}, environment)
requestId: label_values({service_name="toq_server", environment=~"$environment"}, request_id)
```
`requestId` permite selecionar um ID específico ou usar `.*` (valor padrão) para mostrar todos.

### Derived Fields
No painel de logs, o campo `trace_id` está configurado como link para o Jaeger (`http://localhost:16686/trace/${__value.raw}` por padrão). Ajuste a URL conforme ambiente.

## Fluxo de Investigação Recomendado
1. Use o filtro `Request ID` para isolar requisições reportadas.
2. Abra a entrada no painel de logs e copie `trace_id`.
3. Clique no link `TraceID` para abrir o Jaeger no trace correspondente.
4. Utilize o `request_id` nos handlers/serviços para correlacionar logs com métricas e eventos de auditoria.

## Checklist de Smoke
- Gera logs INFO/ERROR após iniciar o stack.
- No Grafana Explore, confirma labels `service_name`, `environment`, `severity`, `request_id`, `trace_id`, `client_ip` e `user_agent`.
- Valida `request_id` derivado dos middlewares (RequestID + ContextWithLogger).
- Garante navegação Jaeger pela derived field.
- Ajusta `LOKI_RETENTION_DAYS` (ex.: `3`) e reinicia Loki para confirmar retenção.

## Compatibilidade do Logger
- Handler `SplitLevelHandler` mantém escrita stdout/stderr/arquivo NDJSON.
- `otelslog` (tee interno) transfere registros com atributos enriquecidos ao collector.
- Não há alteração no formato em disco; ferramentas antigas seguem suportadas.

## Referências
- Configuração: `docker-compose.yml`, `loki-config.yaml`, `otel-collector-config.yaml`.
- Dashboards: `grafana/dashboard-files/*`.
- Guia rápido de Go (`docs/toq_server_go_guide.md`) atualizado para uso do `coreutils.ContextWithLogger`.
- Para produção, considere: TLS no OTLP, retenção menor, alertas Loki (`count_over_time`), IaC para ambiente gerenciado.
