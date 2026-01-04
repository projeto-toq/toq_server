# Observabilidade – Stack TOQ Server

## Visão Geral

- **Objetivo:** monitorar métricas, logs e traces da API REST com correlação automática nativa entre sinais.
- **Pipeline:** aplicação (`slog` + Otel SDK) → Grafana Alloy →
	- **Traces:** Grafana Tempo (armazenamento + métricas RED derivadas)
	- **Logs:** Loki (com labels automáticos para correlação)
	- **Métricas:** Prometheus (via remote write do Alloy)
- **Correlação:** automática via Tempo (trace-to-logs, trace-to-metrics) sem configuração manual.
- **Metadados de serviço:** `service.name`, `service.namespace`, `service.version` e `deployment.environment` são preservados ponta a ponta.
- **Protocolo:** 100% HTTP (OTLP HTTP). Sem gRPC.

## Componentes

| Serviço | Porta | Responsabilidade |
| --- | --- | --- |
| Grafana Alloy | `12345` (UI), `4318` (OTLP HTTP) | Coleta unificada de telemetria, processamento, scraping de métricas |
| Grafana Tempo | `3200` | Backend de traces, métricas RED derivadas, correlação automática |
| Loki | `3100` | Banco de logs estruturados com labels automáticos |
| Prometheus | `9091` | Armazenamento de métricas (recebe remote write do Alloy) |
| Grafana | `3000` | Dashboards provisionados, Explore com correlação nativa |
| MySQL Exporter | `9104` | Métricas do MySQL (scrapedas pelo Alloy) |
| Redis Exporter | `9121` | Métricas do Redis (scrapedas pelo Alloy) |

## Inicialização (Ambiente de Desenvolvimento)

```bash
# Iniciar stack completa
docker compose up -d mysql redis mysql-exporter redis-exporter
docker compose up -d prometheus loki tempo
docker compose up -d alloy
docker compose up -d grafana

# Opcional: metadados de serviço (já configurados via env.yaml)
export ALLOY_SERVICE_NAME=toq_server
export ALLOY_NAMESPACE=projeto-toq
export ALLOY_ENVIRONMENT=homo
```

Execute a aplicação Go: `go run cmd/toq_server.go`. O endpoint OTLP agora aponta para `alloy:4318` (configurado em `configs/env.yaml`).

## Dashboards Provisionados

| Dashboard | UID | Propósito | Variáveis Principais |
|-----------|-----|-----------|---------------------|
| TOQ Server - Aplicação | `toq-app-overview` | Golden signals, runtime Go, métricas RED | `$service`, `$path`, `$method`, `$status` |
| TOQ Server - Infraestrutura | `toq-infra-overview` | Host metrics, MySQL, Redis | `$service`, `$mysql_instance`, `$redis_instance` |
| TOQ Server - Logs | `toq-logs` | **NOVO** - Análise de logs com filtro por severidade | `$service`, `$level`, `$request_id`, `$trace_id`, `$method`, `$path`, `$status`, `$search` |
| TOQ Server - Traces | `toq-traces` | **NOVO** - Análise de traces com correlação bidirecional | `$service`, `$trace_id`, `$operation`, `$status_filter`, `$min_duration`, `$max_duration` |

### Dashboards Separados (Logs e Traces)

A partir de **novembro de 2025**, os dashboards de logs e traces foram separados para melhorar a usabilidade e permitir análises mais profundas.

#### TOQ Server - Logs (`toq-logs`)

**Propósito**: Análise centralizada de logs estruturados com filtros avançados.

**Variáveis principais**:
- `$service`: filtro por serviço (ex.: `toq_server`)
- `$level`: **NOVO** - filtro por severidade (All, ERROR, WARN, INFO, DEBUG)
- `$request_id`: isola logs de uma requisição HTTP específica
- `$trace_id`: correlaciona logs com traces
- `$method`, `$path`, `$status`: filtros HTTP
- `$search`: busca full-text no conteúdo dos logs

**Painéis**:
1. **Logs Estruturados**: painel principal com 20 unidades de altura, exibindo logs com datalink clicável no `trace_id`
2. **Taxa de Logs por Severidade**: time series com cores customizadas (ERROR=vermelho, WARN=laranja)
3. **Distribuição de Logs por Status HTTP**: pie chart
4. **Top 10 Rotas com Erros/Avisos**: tabela de rotas ordenadas por quantidade de ERROR/WARN
5. **Usuários com Mais Erros**: tabela de `user_id` com mais logs de erro

**Navegação para Traces**:
- Clicar no `trace_id` de qualquer log → abre dashboard de traces ou Explore do Tempo
- Link no topo do dashboard → "Ver Dashboard de Traces"

#### TOQ Server - Traces (`toq-traces`)

**Propósito**: Análise de traces distribuídos com correlação para logs e métricas.

**Variáveis principais**:
- `$service`: filtro por serviço
- `$trace_id`: **NOVO** - busca direta por ID de trace
- `$operation`: filtro por rota/operação (ex.: `/api/v1/users`)
- `$status_filter`: **NOVO** - filtro por status HTTP (Todos/Sucesso/Erro/Cliente)
- `$min_duration` / `$max_duration`: **NOVO** - filtros de duração

**Painéis**:
1. **Traces do Serviço**: tabela com datalinks duplos:
   - "Ver Detalhes do Trace" → abre Explore do Tempo
   - "Ver Logs Relacionados" → abre dashboard de logs filtrado por trace_id
2. **Distribuição de Duração de Traces**: histograma
3. **Taxa de Requisições por Status**: time series com dados do Prometheus
4. **Top 10 Rotas Mais Lentas (p95)**: tabela de rotas ordenadas por latência
5. **Mapa de Serviços**: node graph mostrando dependências
6. **Taxa de Erro (%)**: time series calculando percentual de erros 5xx
7. **Latência de Requisições (Percentis)**: time series com p50, p95 e p99

**Navegação para Logs**:
- Clicar no segundo datalink do `Trace ID` na tabela → abre dashboard de logs
- Link no topo do dashboard → "Ver Dashboard de Logs"

### Filtros e Correlação

- **`$service`** restringe ambiente e versão automaticamente (os labels coexistem nas métricas e logs).
- **`$request_id`** agrupa todos os logs de uma requisição HTTP.
- **`$trace_id`** abre diretamente o trace correspondente no Tempo, permitindo navegar span a span.
- **`$level`** filtra logs por severidade (DEBUG, INFO, WARN, ERROR).
- **Derived fields:** no datasource Loki, o `trace_id` está configurado com links para dashboard de traces e Explore.
- **Navegação sem busca manual:**
  - Logs → Traces: clique em `trace_id` **ou** `request_id` para abrir o dashboard de traces já filtrado e com o mesmo intervalo de tempo.
  - Traces → Logs: use “Ver Logs Relacionados” (coluna Trace ID) para abrir o dashboard de logs já filtrado por `trace_id`/`request_id`, preservando `from/to`.

### Fluxo de Investigação Recomendado (Atualizado)

1. Abra **TOQ Server - Aplicação** e valide os "Golden Signals" (latência, taxa, erro).
2. Identifique anomalias de dependência em **TOQ Server - Infraestrutura**.
3. Abra **TOQ Server - Logs**:
   - Filtre por `$level=ERROR` ou `$level=WARN`
   - Use `$path` e `$method` para isolar rota problemática
   - Identifique `request_id` ou `trace_id` relevante
   - Clique no datalink do `trace_id` para abrir traces
4. No **TOQ Server - Traces**:
   - Examine o trace completo no painel principal
   - Clique em "Ver Detalhes do Trace" para Explore detalhado
   - Analise "Top 10 Rotas Mais Lentas" para confirmar padrão
   - Clique em "Ver Logs Relacionados" se precisar voltar ao contexto de logs
5. Confirme resolução voltando ao dashboard de aplicação.

### Correlação Automática

- **Logs → Traces**: Derived fields em Loki extraem `trace_id` e criam links diretos para dashboard de traces e Explore do Tempo.
- **Traces → Logs**: Datalinks na tabela de traces permitem navegar para logs filtrados por `trace_id`.
- **Traces → Métricas**: Tempo exibe métricas RED derivadas automaticamente dos spans.

## Checklist Pós-Subida

- [ ] Alloy UI acessível em `http://localhost:12345` e mostra componentes saudáveis.
- [ ] Tempo retorna traces: `curl http://localhost:3200/api/search`.
- [ ] Loki retorna logs com labels `service_name`, `trace_id`, `level`: `curl http://localhost:3100/ready`.
- [ ] Prometheus exibe métricas `http_requests_total{service="toq_server"}`.
- [ ] Grafana Explore → Tempo: buscar trace por `trace_id` funciona.
- [ ] Correlação: clicar em `trace_id` em log abre trace no Tempo.
- [ ] Métricas RED derivadas aparecem no painel de traces do Tempo.
- [ ] Variáveis dos dashboards retornam valores (teste `service`, `request_id`, `trace_id`).

### Checklist Pós-Refatoração (Logs e Traces Separados)

Após aplicar as mudanças de separação de dashboards:

- [ ] Alloy reiniciado e logs exibindo label `level` no Loki
  ```bash
  # Verificar label_values de level
  curl -G "http://localhost:3100/loki/api/v1/label/level/values" | jq
  # Deve retornar: ["DEBUG", "INFO", "WARN", "ERROR"]
  ```
- [ ] Dashboard "TOQ Server - Logs" carrega sem erros
- [ ] Variável `$level` retorna opções (All, ERROR, WARN, INFO, DEBUG)
- [ ] Filtrar `$level=ERROR` exibe apenas logs de erro
- [ ] Clicar em `trace_id` no painel de logs abre dashboard de traces
- [ ] Dashboard "TOQ Server - Traces" carrega sem erros
- [ ] Variável `$trace_id` (textbox) aceita entrada manual
- [ ] Variável `$status_filter` funciona (filtrar por "Erro (5xx)" funciona)
- [ ] Clicar em "Ver Logs Relacionados" na tabela de traces abre dashboard de logs filtrado
- [ ] Painel "Mapa de Serviços" exibe dependências (se houver mais de um serviço)
- [ ] Links bidirecionais entre dashboards preservam time range (`from`/`to`)
- [ ] Navegação não requer busca manual: dos logs para traces via `trace_id`/`request_id`; dos traces para logs via "Ver Logs Relacionados"

## Configurações

- **Alloy**: `alloy/config.alloy` (River syntax, HTTP only)
  - Label `level` adicionado em `loki.attribute.labels`
- **Tempo**: `tempo/tempo-config.yaml` (HTTP receiver only)
- **Loki**: `loki-config.yaml` (atualizado com retention e structured metadata)
- **Prometheus**: `prometheus.yml` (simplificado, apenas self-monitoring)
- **Grafana Datasources**: `grafana/datasources/{tempo,loki,prometheus}.yml`
  - Loki configurado com derived fields bidirecionais
- **App config**: `configs/env.yaml` (endpoint OTLP → `alloy:4318`)

## Troubleshooting

- **Traces não aparecem**: verificar Alloy logs (`docker logs alloy`) e health do Tempo (`curl http://localhost:3200/ready`).
- **Correlação quebrada**: validar derived fields no datasource Loki (Grafana UI → Configuration → Data Sources → Loki).
- **Métricas RED não aparecem**: confirmar `metrics_generator` ativo no Tempo (checar `tempo/tempo-config.yaml`).
- **Logs sem label `level`**: verificar `loki.attribute.labels` no processador de atributos do Alloy (`alloy/config.alloy`).
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

# Testar query por severidade
curl -G -s "http://localhost:3100/loki/api/v1/query" \
  --data-urlencode 'query={service_name="toq_server", level="ERROR"}' \
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

## Separação de Dashboards de Logs e Traces

Implementada em **6 de novembro de 2025**:
- ✅ Dashboard único `toq-server-logs-traces` separado em dois dashboards dedicados
- ✅ Novo dashboard `toq-logs` com filtro por severidade (`level`)
- ✅ Novo dashboard `toq-traces` com filtros avançados e métricas
- ✅ Label `level` promovido para indexação no Loki
- ✅ Derived fields bidirecionais configurados (logs ↔ traces)
- ✅ Documentação atualizada com novos fluxos de investigação

## Referências

- Grafana Alloy: https://grafana.com/docs/alloy/
- Grafana Tempo: https://grafana.com/docs/tempo/
- TraceQL: https://grafana.com/docs/tempo/latest/traceql/
- Telemetria no código: `internal/core/config/telemetry.go`, middlewares em `internal/adapter/left/http/middlewares/*`
- Guia do projeto: `docs/toq_server_go_guide.md`
