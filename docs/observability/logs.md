# Observabilidade de Logs com Loki + Grafana

> **Foco desta fase (PR-05)**: assegurar que toda a cadeia de observabilidade reflita o novo logger contextual (com `request_id`, `trace_id` e metadados de requisição) e que a experiência no Grafana destaque esses campos.

## Visão Geral
- **Objetivo**: Persistir e explorar logs estruturados do TOQ Server com correlação de traces e recursos de filtro por `request_id`.
- **Pipeline**: `ContextWithLogger` (handlers/middlewares/workers) → `slog` → ponte OpenTelemetry → Collector (`transform/log_labels`) → Loki (OTLP) → Grafana.
- **Compatibilidade**: Mantém JSON/stdout existentes e promove `request_id`, `severity`, `path`, `method`, `user_agent`, `trace_id` e `span_id` a labels no Loki via atributo especial `loki.attribute.labels`.

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

Atualizações recentes:
- Variáveis de filtro **Request ID**, **Severity** e **HTTP Path** aplicadas diretamente sobre labels Loki (`request_id`, `severity`, `path`).
- Campo derivado `trace_id` segue apontando para Jaeger por meio do datasource configurado no Loki.
- Consultas (`count_over_time` e painéis de logs) usam os filtros selecionados para reduzir ruído e acelerar a investigação.

### Painéis
1. **Eventos por severidade** — Timeseries usando `count_over_time`, segmentado por `severity` com filtros de `request_id` e `path`.
2. **Fluxo de logs** — Painel de logs com labels comuns (`request_id`, `severity`, `trace_id`, `path`).
3. **Erros por serviço** — Tabela agregada filtrável por `request_id` e `path`, destacando requisições com `severity="ERROR"`.

### Variáveis de Template
```text
requestId: label_values({service_name="toq_server"}, request_id)
severity:  label_values({service_name="toq_server", request_id=~"$requestId"}, severity)
path:      label_values({service_name="toq_server", request_id=~"$requestId"}, path)
```
Selecione valores específicos quando precisar investigar um fluxo; o valor `.*` aplica o filtro sobre todos os registros.

### Derived Fields
No painel de logs, o campo `trace_id` está configurado como link para o Jaeger (`http://localhost:16686/trace/${__value.raw}` por padrão). Ajuste a URL conforme ambiente.

## Fluxo de Investigação Recomendado
1. **Execute o servidor com `ENVIRONMENT=homo`** para garantir que logs, métricas e traces sejam exportados para o collector.
2. Use os filtros `Request ID`, `Severity` e `Path` no dashboard “TOQ Server Logs” para isolar o cenário alvo.
3. Abra uma linha de log, utilize o link “Jaeger Trace” (campo `trace_id`) e valide o span correspondente no Jaeger.
4. Opcionalmente, migre para o painel “TOQ Server - Observability Correlation” e reaproveite os filtros para cruzar métricas e traces.

## Checklist de Smoke
- Verifique no Grafana Explore se as labels `service_name`, `severity`, `request_id`, `trace_id`, `path`, `method`, `client_ip` e `user_agent` estão disponíveis.
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
