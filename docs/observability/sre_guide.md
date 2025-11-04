# Observabilidade – Stack TOQ Server

## Visão Geral
- **Objetivo:** monitorar métricas, logs e traces da API com correlação imediata entre requisições, dependências e infraestrutura.
- **Pipeline:** aplicação (`slog` + Otel SDK) → OpenTelemetry Collector →
	- **Traces:** exportados para Jaeger via OTLP gRPC.
	- **Métricas:** expostas pelo collector em `/metrics` e consumidas pelo Prometheus.
	- **Logs:** enviados para o Loki exporter nativo com labels (`trace_id`, `request_id`, `method`, `path`, `status`, `service.name`).
- **Metadados de serviço:** `service.name`, `service.namespace`, `service.version` e `deployment.environment` são preservados ponta a ponta. Isso permite filtrar por serviço em qualquer dashboard ou consulta no Grafana/Jaeger.

## Componentes
| Serviço | Porta | Responsabilidade |
| --- | --- | --- |
| OpenTelemetry Collector | `4317` / `4318` | Recebe sinais OTLP, agrega atributos e expõe métricas consolidando hostmetrics. |
| Loki | `3100` | Banco de logs para consultas estruturadas e alertas. |
| Jaeger | `16686` | UI de traces distribuídos (busca por `trace_id`). |
| Prometheus | `9091` (loopback) | Scrape de métricas do collector, aplicação, MySQL e Redis. |
| Grafana | `3000` (loopback) | Dashboards provisionados e modo Explore. |
| MySQL Exporter | `9104` | Métricas detalhadas do MySQL (threads, InnoDB, locks). |
| Redis Exporter | `9121` | Métricas de cache (hits/misses, memória, clientes). |

## Inicialização (Ambiente de Desenvolvimento)
```bash
docker compose up -d prometheus loki jaeger otel-collector grafana mysql-exporter redis-exporter

# Opcional: expor nome/namespace do serviço ao iniciar o binário localmente
export OTEL_SERVICE_NAME=toq_server
export OTEL_RESOURCE_NAMESPACE=projeto-toq
export OTEL_RESOURCE_ENVIRONMENT=homo
```
Em seguida execute `go run cmd/toq_server.go` (ou utilize o binário em `bin/toq_server`). O bootstrap ativa tracing, métricas e logs automaticamente quando `telemetry.enabled=true` (perfil `homo`).

## Dashboards Provisionados
- **TOQ Server - Aplicação (`toq-app-overview`)** – Golden signals (QPS, p95, erros, in-flight), saúde do runtime Go (CPU, heap, goroutines) e tabela de status por rota. Variáveis principais: `service`, `path`, `method`, `status`.
- **TOQ Server - Infraestrutura (`toq-infra-overview`)** – Host metrics coletadas pelo collector, MySQL (threads, taxa de queries, locks InnoDB) e Redis (clientes, hit/miss, memória). Variáveis: `service`, `mysql_instance`, `redis_instance`.
- **TOQ Server - Logs e Traces (`toq-logs-traces`)** – Painel de correlação com logs Loki, trace viewer Jaeger e séries de erro por rota. Variáveis: `service`, `request_id`, `trace_id`, `method`, `path`, `status`.

### Filtros e Correlação
- **`$service`** restringe ambiente e versão automaticamente (os labels coexistem nas métricas e logs).
- **`$request_id`** agrupa todos os logs de uma requisição HTTP; o painel de traces usa o mesmo valor.
- **`$trace_id`** abre diretamente o trace correspondente no painel Jaeger, permitindo navegar span a span.
- **Derived field:** no datasource Loki, o `trace_id` já está configurado com link para `http://localhost:16686/trace/${__value.raw}`.

### Fluxo de Investigação Recomendado
1. Abra **TOQ Server - Aplicação** e valide os “Golden Signals” (latência, taxa, erro, saturação). Ajuste `$service/$path/$method` conforme necessário.
2. Identifique anomalias de dependência em **TOQ Server - Infraestrutura** (ex.: aumento de `mysql_global_status_threads_connected`).
3. No dashboard **Logs e Traces**, filtre por `request_id` ou `trace_id` e:
	 - inspecione o log estruturado com labels;
	 - clique no link do `trace_id` para abrir o trace completo no Jaeger;
	 - analise a tabela “Erros por rota” para confirmar o impacto.
4. Caso precise investigar payloads específicos, utilize Grafana Explore → Loki com a mesma query gerada pelo painel.

## Checklist Pós-Subida
- [ ] Prometheus exibe métricas `http_requests_total` com labels `service`, `method`, `path`, `status`.
- [ ] Variáveis dos dashboards retornam valores (teste `service`, `request_id`, `mysql_instance`, `redis_instance`).
- [ ] Logs em Loki aparecem com labels `trace_id` e `request_id` e permitem abrir o trace no Jaeger.
- [ ] Traces exibem `service.name=toq_server` na UI do Jaeger.
- [ ] Host metrics (`system_cpu_time`, `system_memory_usage`) estão sendo coletadas pelo collector.

## Referências
- Configurações: `docker-compose.yml`, `otel-collector-config.yaml`, `loki-config.yaml`.
- Dashboards provisionados: `grafana/dashboard-files/*.json`.
- Telemetria no código: `internal/core/config/telemetry.go`, middlewares em `internal/adapter/left/http/middlewares/*`.
