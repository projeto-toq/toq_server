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
| Grafana | `3000` | Dashboards (`TOQ Server - Operations Overview`, `TOQ Server - Dependencies Observability`, `TOQ Server - Logs Analytics`, `TOQ Server - Observability Triage`) e consultas Explore |
| Jaeger | `16686` | Consulta de traces vinculados via `trace_id` |
| MySQL Exporter | `9104` | Exposição de métricas MySQL (threads, buffer pool, latency) para Prometheus |
| Redis Exporter | `9121` | Exposição de métricas Redis (clientes, memória, hits/miss) para Prometheus |

## Variáveis de Ambiente Importantes
- `LOKI_RETENTION_DAYS`: dias de retenção do Loki. Dev = `7`, sugerido prod = `3`.
- `OTEL_RESOURCE_SERVICE_NAME`: nome lógico exportado como `service.name` (default `toq_server`).
- `OTEL_RESOURCE_ENVIRONMENT`: ambiente (`dev`, `homo`, `prod`), replicado para métricas/logs/traces.
- `OTEL_RESOURCE_SERVICE_VERSION`: versão da aplicação exibida nos dashboards.
- `OTEL_EXPORTER_OTLP_TRACES_ENDPOINT`: endpoint protegido do backend de traces (Jaeger/Tempo).
- `OTEL_EXPORTER_OTLP_TRACES_AUTH_HEADER`: header `Authorization` usado no exporter OTLP gRPC (opcional).
- `LOKI_OTLP_ENDPOINT`: endpoint OTLP/HTTP do Loki ou gateway de logs.
- `LOKI_OTLP_AUTH_HEADER`: header `Authorization` para Loki (opcional).
- `LOKI_TENANT_ID`: identificador multi-tenant (ex.: `toq`).

> Altere valores via `.env` para manter os stacks sincronizados. PR-05 não muda defaults, apenas registra expectativas.

## Passos para Subir o Stack (dev)
```bash
export OTEL_RESOURCE_SERVICE_NAME=toq_server
export OTEL_RESOURCE_ENVIRONMENT=homo
export OTEL_RESOURCE_SERVICE_VERSION=2.0.0
export OTEL_EXPORTER_OTLP_TRACES_ENDPOINT=jaeger:4317
export OTEL_EXPORTER_OTLP_TRACES_INSECURE=true
export LOKI_OTLP_ENDPOINT=http://loki:3100/otlp
export LOKI_OTLP_INSECURE=true
export LOKI_TENANT_ID=toq
export LOKI_RETENTION_DAYS=7
docker compose up -d loki otel-collector grafana jaeger mysql-exporter redis-exporter
```
Inicie o servidor (`go run cmd/toq_server.go` ou binário). O bootstrap já conecta o `ctxlogger` aos handlers e workers.

## Dashboards disponíveis

- **TOQ Server - Operations Overview** (`grafana/dashboard-files/toq-server-operations-overview.json`)
	- Golden signals (latência p95/p99, tráfego por método, erros por classe, saturação de memória/goroutines).
	- Painéis adicionais: GC pause médio, `cache_operations_total` (hits/miss/expired em 5 min), `redis_connected_clients`, heap alloc e HTTP in-flight.
	- Variáveis: `Environment`, `Version` (via métricas Prometheus com labels constantes) para isolar instâncias.

- **TOQ Server - Dependencies Observability** (`grafana/dashboard-files/toq-server-dependencies-observability.json`)
	- Métricas de MySQL/Redis (throughput, latência p95, erros) + saturação de CPU/memória coletadas via hostmetrics.
	- Útil para validar regressões em queries, lock contention e impactos de cache.

- **TOQ Server - Logs Analytics** (`grafana/dashboard-files/toq-server-logs-analytics.json`)
	- Timeseries de volume por `severity`, top paths de erro, painel de logs contextualizados e bargauge por `user_agent`.
	- Variáveis: `Environment`, `Version`, `Severity`, `Request ID`, `Path` — todas alimentadas por labels Loki produzidos pelo collector (`deployment.environment`, `service.version`, etc.).

- **TOQ Server - Observability Triage** (`grafana/dashboard-files/toq-server-observability-triage.json`)
	- Vista única para correlação HTTP status vs. latência, logs filtrados, traces Jaeger (via tags `deployment.environment` e `request_id`) e tabela de requisições com maior incidência de erros.

Cada dashboard está pensado para investigação rápida: escolha o ambiente/versão, isolate o `request_id` (se conhecido) e acompanhe métricas, logs e traces no mesmo contexto.

### Derived Fields
No painel **Logs Analytics**, o campo `trace_id` permanece com link configurável para Jaeger (`http://localhost:16686/trace/${__value.raw}` por padrão). Ajuste conforme o domínio real do cluster.

## Métricas Redis Cache
- O adapter Redis (`internal/core/cache/redis_cache.go`) agora injeta instrumentação OpenTelemetry (`redisotel`) para traces/métricas nativas do cliente.
- Integração com o Prometheus Adapter registra operações em `cache_operations_total{operation="get|set|delete", result="hit|miss|expired|success|error"}`.
- Painéis “Redis Cache Operations (5m)” e “Redis Connected Clients” no dashboard de operações exibem taxa de hits/miss e saúde do cache.
- Qualquer erro de marshal/Redis é refletido como `result="error"`, permitindo alarmes baseados em proporção de falhas.

## Fluxo de Investigação Recomendado
1. **Execute o servidor com `ENVIRONMENT=homo`** para garantir que logs, métricas e traces sejam exportados para o collector.
2. Use os filtros `Environment`, `Version`, `Severity`, `Path` e `Request ID` no dashboard “TOQ Server - Logs Analytics” para isolar o cenário alvo.
3. Abra uma linha de log, utilize o link “Jaeger Trace” (campo `trace_id`) e valide o span correspondente no Jaeger.
4. Migre para o painel “TOQ Server - Observability Triage” para cruzar métricas e traces com o mesmo filtro.

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
