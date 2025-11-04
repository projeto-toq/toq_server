# Observabilidade – Stack TOQ Server (Grafana Alloy + Tempo + Loki)# Observabilidade – Stack TOQ Server



## Visão Geral## Visão Geral

- **Objetivo:** monitorar métricas, logs e traces da API REST com correlação automática nativa entre sinais.- **Objetivo:** monitorar métricas, logs e traces da API com correlação imediata entre requisições, dependências e infraestrutura.

- **Pipeline:** aplicação (`slog` + Otel SDK) → Grafana Alloy →- **Pipeline:** aplicação (`slog` + Otel SDK) → OpenTelemetry Collector →

	- **Traces:** Grafana Tempo (armazenamento + métricas RED derivadas)	- **Traces:** exportados para Jaeger via OTLP gRPC.

	- **Logs:** Loki (com labels automáticos para correlação)	- **Métricas:** expostas pelo collector em `/metrics` e consumidas pelo Prometheus.

	- **Métricas:** Prometheus (via remote write do Alloy)	- **Logs:** enviados para o Loki exporter nativo com labels (`trace_id`, `request_id`, `method`, `path`, `status`, `service.name`).

- **Correlação:** automática via Tempo (trace-to-logs, trace-to-metrics) sem configuração manual.- **Metadados de serviço:** `service.name`, `service.namespace`, `service.version` e `deployment.environment` são preservados ponta a ponta. Isso permite filtrar por serviço em qualquer dashboard ou consulta no Grafana/Jaeger.

- **Metadados de serviço:** `service.name`, `service.namespace`, `service.version` e `deployment.environment` são preservados ponta a ponta.

- **Protocolo:** 100% HTTP (OTLP HTTP). Sem gRPC.## Componentes

| Serviço | Porta | Responsabilidade |

## Componentes| --- | --- | --- |

| Serviço | Porta | Responsabilidade || OpenTelemetry Collector | `4317` / `4318` | Recebe sinais OTLP, agrega atributos e expõe métricas consolidando hostmetrics. |

| --- | --- | --- || Loki | `3100` | Banco de logs para consultas estruturadas e alertas. |

| Grafana Alloy | `12345` (UI), `4318` (OTLP HTTP) | Coleta unificada de telemetria, processamento, scraping de métricas || Jaeger | `16686` | UI de traces distribuídos (busca por `trace_id`). |

| Grafana Tempo | `3200` | Backend de traces, métricas RED derivadas, correlação automática || Prometheus | `9091` (loopback) | Scrape de métricas do collector, aplicação, MySQL e Redis. |

| Loki | `3100` | Banco de logs estruturados com labels automáticos || Grafana | `3000` (loopback) | Dashboards provisionados e modo Explore. |

| Prometheus | `9091` | Armazenamento de métricas (recebe remote write do Alloy) || MySQL Exporter | `9104` | Métricas detalhadas do MySQL (threads, InnoDB, locks). |

| Grafana | `3000` | Dashboards provisionados, Explore com correlação nativa || Redis Exporter | `9121` | Métricas de cache (hits/misses, memória, clientes). |

| MySQL Exporter | `9104` | Métricas do MySQL (scrapedas pelo Alloy) |

| Redis Exporter | `9121` | Métricas do Redis (scrapedas pelo Alloy) |## Inicialização (Ambiente de Desenvolvimento)

```bash

## Arquitetura de Fluxodocker compose up -d prometheus loki jaeger otel-collector grafana mysql-exporter redis-exporter

```

Aplicação Go (REST API) → Alloy (OTLP HTTP :4318) → {Tempo, Loki, Prometheus}# Opcional: expor nome/namespace do serviço ao iniciar o binário localmente

Exporters (MySQL/Redis) → Alloy (scrape) → Prometheusexport OTEL_SERVICE_NAME=toq_server

```export OTEL_RESOURCE_NAMESPACE=projeto-toq

export OTEL_RESOURCE_ENVIRONMENT=homo

## Inicialização (Ambiente de Desenvolvimento)```

```bashEm seguida execute `go run cmd/toq_server.go` (ou utilize o binário em `bin/toq_server`). O bootstrap ativa tracing, métricas e logs automaticamente quando `telemetry.enabled=true` (perfil `homo`).

# Iniciar stack completa

docker compose up -d mysql redis mysql-exporter redis-exporter## Dashboards Provisionados

docker compose up -d prometheus loki tempo- **TOQ Server - Aplicação (`toq-app-overview`)** – Golden signals (QPS, p95, erros, in-flight), saúde do runtime Go (CPU, heap, goroutines) e tabela de status por rota. Variáveis principais: `service`, `path`, `method`, `status`.

docker compose up -d alloy- **TOQ Server - Infraestrutura (`toq-infra-overview`)** – Host metrics coletadas pelo collector, MySQL (threads, taxa de queries, locks InnoDB) e Redis (clientes, hit/miss, memória). Variáveis: `service`, `mysql_instance`, `redis_instance`.

docker compose up -d grafana- **TOQ Server - Logs e Traces (`toq-logs-traces`)** – Painel de correlação com logs Loki, trace viewer Jaeger e séries de erro por rota. Variáveis: `service`, `request_id`, `trace_id`, `method`, `path`, `status`.



# Opcional: metadados de serviço (já configurados via env.yaml)### Filtros e Correlação

export ALLOY_SERVICE_NAME=toq_server- **`$service`** restringe ambiente e versão automaticamente (os labels coexistem nas métricas e logs).

export ALLOY_NAMESPACE=projeto-toq- **`$request_id`** agrupa todos os logs de uma requisição HTTP; o painel de traces usa o mesmo valor.

export ALLOY_ENVIRONMENT=homo- **`$trace_id`** abre diretamente o trace correspondente no painel Jaeger, permitindo navegar span a span.

```- **Derived field:** no datasource Loki, o `trace_id` já está configurado com link para `http://localhost:16686/trace/${__value.raw}`.



Execute a aplicação Go: `go run cmd/toq_server.go`. O endpoint OTLP agora aponta para `alloy:4318` (configurado em `configs/env.yaml`).### Fluxo de Investigação Recomendado

1. Abra **TOQ Server - Aplicação** e valide os “Golden Signals” (latência, taxa, erro, saturação). Ajuste `$service/$path/$method` conforme necessário.

## Dashboards Provisionados2. Identifique anomalias de dependência em **TOQ Server - Infraestrutura** (ex.: aumento de `mysql_global_status_threads_connected`).

- **TOQ Server - Aplicação** (`toq-server-app-overview`): Golden signals (QPS, p95, erros), runtime Go, métricas RED derivadas de traces. Variáveis: `service`, `path`, `method`, `status`.3. No dashboard **Logs e Traces**, filtre por `request_id` ou `trace_id` e:

- **TOQ Server - Infraestrutura** (`toq-server-infra-overview`): Host metrics, MySQL (threads, queries, InnoDB), Redis (cache hit/miss, memória). Variáveis: `service`, `mysql_instance`, `redis_instance`.	 - inspecione o log estruturado com labels;

- **TOQ Server - Logs e Traces** (`toq-server-logs-traces`): Correlação automática; clique em `trace_id` nos logs para abrir trace no Tempo; painel de erros por rota. Variáveis: `service`, `request_id`, `trace_id`, `method`, `path`, `status`.	 - clique no link do `trace_id` para abrir o trace completo no Jaeger;

	 - analise a tabela “Erros por rota” para confirmar o impacto.

### Correlação Automática4. Caso precise investigar payloads específicos, utilize Grafana Explore → Loki com a mesma query gerada pelo painel.

- **Logs → Traces**: Derived field em Loki extrai `trace_id` e cria link direto para Tempo Explore.

- **Traces → Logs**: Tempo injeta botão "Logs for this span" com filtro automático por `trace_id` + timerange.## Checklist Pós-Subida

- **Traces → Métricas**: Tempo exibe RED metrics (Rate/Error/Duration) derivadas automaticamente dos spans.- [ ] Prometheus exibe métricas `http_requests_total` com labels `service`, `method`, `path`, `status`.

- [ ] Variáveis dos dashboards retornam valores (teste `service`, `request_id`, `mysql_instance`, `redis_instance`).

### Fluxo de Investigação Recomendado- [ ] Logs em Loki aparecem com labels `trace_id` e `request_id` e permitem abrir o trace no Jaeger.

1. Abra dashboard **TOQ Server - Aplicação** e valide os "Golden Signals" (latência p95, taxa de requisições, erro rate).- [ ] Traces exibem `service.name=toq_server` na UI do Jaeger.

2. Identifique anomalias de dependência em **TOQ Server - Infraestrutura** (ex.: aumento de `mysql_global_status_threads_connected`).- [ ] Host metrics (`system_cpu_time`, `system_memory_usage`) estão sendo coletadas pelo collector.

3. No dashboard **Logs e Traces**, filtre por `request_id` ou `trace_id` e:

	 - Inspecione o log estruturado com labels.## Referências

	 - Clique no link do `trace_id` para abrir o trace completo no Tempo.- Configurações: `docker-compose.yml`, `otel-collector-config.yaml`, `loki-config.yaml`.

	 - Analise a tabela "Erros por rota" para confirmar o impacto.- Dashboards provisionados: `grafana/dashboard-files/*.json`.

4. No Tempo, clique em um span específico → "Logs for this span" → Loki com contexto exato do timerange.- Telemetria no código: `internal/core/config/telemetry.go`, middlewares em `internal/adapter/left/http/middlewares/*`.

5. Analise métricas RED derivadas no Tempo para confirmar impacto sistêmico.

## Checklist Pós-Subida
- [ ] Alloy UI acessível em `http://localhost:12345` e mostra componentes saudáveis.
- [ ] Tempo retorna traces: `curl http://localhost:3200/api/search`.
- [ ] Loki retorna logs com labels `service_name`, `trace_id`: `curl http://localhost:3100/ready`.
- [ ] Prometheus exibe métricas `http_requests_total{service="toq_server"}`.
- [ ] Grafana Explore → Tempo: buscar trace por `trace_id` funciona.
- [ ] Correlação: clicar em `trace_id` em log abre trace no Tempo.
- [ ] Métricas RED derivadas aparecem no painel de traces do Tempo.
- [ ] Variáveis dos dashboards retornam valores (teste `service`, `request_id`, `trace_id`).

## Configurações
- **Alloy**: `alloy/config.alloy` (River syntax, HTTP only)
- **Tempo**: `tempo/tempo-config.yaml` (HTTP receiver only)
- **Loki**: `loki-config.yaml` (atualizado com retention e structured metadata)
- **Prometheus**: `prometheus.yml` (simplificado, apenas self-monitoring)
- **Grafana Datasources**: `grafana/datasources/{tempo,loki,prometheus}.yml`
- **App config**: `configs/env.yaml` (endpoint OTLP → `alloy:4318`)

## Troubleshooting
- **Traces não aparecem**: verificar Alloy logs (`docker logs alloy`) e health do Tempo (`curl http://localhost:3200/ready`).
- **Correlação quebrada**: validar derived fields no datasource Loki (Grafana UI → Configuration → Data Sources → Loki).
- **Métricas RED não aparecem**: confirmar `metrics_generator` ativo no Tempo (checar `tempo/tempo-config.yaml`).
- **Logs sem labels**: verificar `loki.attribute.labels` no processador de atributos do Alloy (`alloy/config.alloy`).
- **Alloy não inicia**: validar sintaxe River com `docker logs alloy` e procurar por parsing errors.
- **Prometheus não recebe métricas**: confirmar flag `--web.enable-remote-write-receiver` no Prometheus (docker-compose.yml).

## Comandos Úteis

### Validar Health dos Serviços
```bash
# Alloy
curl http://localhost:12345/ready

# Tempo
curl http://localhost:3200/ready

# Loki
curl http://localhost:3100/ready

# Prometheus
curl http://localhost:9091/-/ready

# Grafana
curl http://localhost:3000/api/health
```

### Reiniciar Stack Completa
```bash
docker compose restart alloy tempo loki prometheus grafana
```

### Ver Logs de um Serviço Específico
```bash
docker logs -f alloy
docker logs -f tempo
docker logs -f loki
```

### Testar Query no Loki
```bash
curl -G -s "http://localhost:3100/loki/api/v1/query" \
  --data-urlencode 'query={service_name="toq_server"}' \
  --data-urlencode 'limit=10' | jq
```

### Testar Query no Tempo (TraceQL)
```bash
# Buscar traces com erro
curl -G "http://localhost:3200/api/search" \
  -d 'q={status=error}' \
  -d 'limit=10' | jq
```

## Migração do Stack Anterior (Jaeger + OTel Collector)
A migração foi executada em **4 de novembro de 2025**:
- ✅ Jaeger substituído por Grafana Tempo
- ✅ OpenTelemetry Collector substituído por Grafana Alloy
- ✅ Dashboards migrados automaticamente (script: `scripts/migrate_dashboards.sh`)
- ✅ Datasource Jaeger removido, Tempo provisionado
- ✅ Configurações de correlação automática ativadas

## Referências
- Grafana Alloy: https://grafana.com/docs/alloy/
- Grafana Tempo: https://grafana.com/docs/tempo/
- TraceQL: https://grafana.com/docs/tempo/latest/traceql/
- Telemetria no código: `internal/core/config/telemetry.go`, middlewares em `internal/adapter/left/http/middlewares/*`
- Guia do projeto: `docs/toq_server_go_guide.md`
