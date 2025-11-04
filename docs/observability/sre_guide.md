# Observabilidade de Logs com Loki + Grafana

## Visão Geral
- **Objetivo**: Persistir e explorar logs estruturados do TOQ Server com correlação de traces e recursos de filtro por identificadores de requisição.
- **Pipeline**: `slog` na aplicação → ponte OpenTelemetry → Collector → Loki (OTLP) → Grafana.
- **Correlação**: A correlação entre logs e traces é feita através do `trace_id`. O `request_id` permite agrupar todos os logs de uma mesma requisição.

## Componentes
| Serviço | Porta | Responsabilidade |
| --- | --- | --- |
| Loki | `3100` | Time-series DB e API de consulta para logs |
| OpenTelemetry Collector | `4317/4318` | Recebe spans/logs OTLP, agrega atributos e envia ao Loki |
| Grafana | `3000` | Dashboards e consultas Explore |
| Jaeger | `16686` | Consulta de traces vinculados via `trace_id` |
| MySQL Exporter | `9104` | Exposição de métricas MySQL para Prometheus |
| Redis Exporter | `9121` | Exposição de métricas Redis para Prometheus |

## Passos para Subir o Stack (dev)
```bash
# Não é mais necessário exportar variáveis de ambiente para serviço, versão ou ambiente.
docker compose up -d
```
Inicie o servidor (`go run cmd/toq_server.go` ou binário). O bootstrap já conecta o logger ao pipeline de telemetria.

## Dashboards
- **TOQ Server - Observability Triage**: Este é o dashboard principal e o ponto de partida para qualquer investigação. Ele consolida os "Golden Signals" (Latência, Tráfego, Erros, Saturação), sinais vitais da aplicação (CPU, Memória, Goroutines) e painéis de logs e traces.
- **TOQ Server - Dependencies Observability**: Fornece uma visão detalhada da saúde e performance do MySQL e do Redis.

### Filtros e Correlação
Os dashboards foram simplificados e não utilizam mais filtros de `serviço`, `versão` ou `ambiente`. A investigação deve ser focada no uso dos seguintes filtros:
- **`$path`**: Filtra por rota HTTP.
- **`$method`**: Filtra por método HTTP.
- **`$request_id`**: Isola todos os logs de uma única requisição.
- **`$trace_id`**: Isola uma transação distribuída completa (logs e spans).

### Derived Fields
No painel de Logs, o campo `trace_id` possui um link configurado para a UI do Jaeger (`http://localhost:16686/trace/${__value.raw}`), permitindo a navegação direta do log para o trace correspondente.

## Fluxo de Investigação Recomendado
1. **Acesse o dashboard "TOQ Server - Observability Triage"** no Grafana.
2. **Analise os painéis de "Golden Signals"** para identificar anomalias (picos de latência, aumento na taxa de erros, etc.).
3. **Utilize os filtros** (`$path`, `$method`) para isolar o escopo do problema.
4. **Examine o painel de Logs** para encontrar mensagens de erro ou logs relevantes.
5. **Copie um `trace_id`** de um log de interesse e **cole no filtro `$trace_id`** para ver o trace completo no painel do Jaeger e todos os logs associados.
6. Alternativamente, **clique no link do `trace_id`** no log para abrir a UI do Jaeger em uma nova aba.
7. Se a suspeita recair sobre o banco de dados ou cache, **navegue para o dashboard "TOQ Server - Dependencies Observability"** para analisar métricas específicas de MySQL e Redis.

## Referências
- Configuração: `docker-compose.yml`, `loki-config.yaml`, `otel-collector-config.yaml`.
- Dashboards: `grafana/dashboard-files/*`.
