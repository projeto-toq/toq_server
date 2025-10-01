# Observabilidade de Logs com Loki + Grafana

## Visão Geral
- **Objetivo**: Persistir e explorar logs estruturados do TOQ Server com correlação de traces.
- **Pipeline**: `slog` → *bridge* OpenTelemetry → Collector (OTLP HTTP) → Exporter Loki → Grafana (Explore/Dashboards).
- **Compatibilidade**: Mantém formato JSON atual e split stdout/stderr, adicionando exportação OTLP sem alterar handlers existentes.

## Componentes
| Serviço | Porta | Responsabilidade |
| --- | --- | --- |
| Loki | `3100` | Armazenamento TSDB e API de consulta para logs |
| OpenTelemetry Collector | `4317/4318` | Recebe logs OTLP, enriquece labels e envia ao Loki |
| Grafana | `3000` | Visualização (dashboard `TOQ Server Logs`) e consultas Explore |
| Jaeger | `16686` | Consulta de traces para correlação via campo `trace_id` |

## Variáveis de Ambiente
| Variável | Default (dev) | Produção sugerida | Descrição |
| --- | --- | --- | --- |
| `LOKI_RETENTION_DAYS` | `7` | `3` | Define período de retenção em dias. Ajuste via arquivo `.env` ou export no shell antes do `docker compose up`. |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | `http://localhost:4318` | conforme ambiente | Já utilizada para traces/métricas; logs compartilham o endpoint |
| `OTEL_EXPORTER_OTLP_INSECURE` | `true` | `false` quando houver TLS | Controla uso de TLS com o collector |
| `OTEL_RESOURCE_ENVIRONMENT` | `development` | `production` | Rotula métricas, traces e logs com o ambiente lógico; consumido pelo collector |

> ℹ️ Para ambientes que desejam valores diferentes em produção, configure `LOKI_RETENTION_DAYS=3` no arquivo `.env` consumido pelo Docker Compose ou no pipeline de deploy.

## Subindo o Stack (dev)
1. Gere o binário do servidor (ou use `go run` com flags de log desejadas).
2. Configure variáveis conforme necessidade:
   ```bash
   export OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318
   export OTEL_EXPORTER_OTLP_INSECURE=true
   export LOKI_RETENTION_DAYS=7
   ```
3. Levante os serviços de observabilidade:
   ```bash
  docker compose up -d loki otel-collector grafana jaeger
   ```
4. Inicie o `toq_server` com `--log-output file` (ou outro formato); a ponte `slog → OTel` é configurada automaticamente durante o bootstrap.

## Exploração no Grafana
- Data source `Loki` é provisionado automaticamente (`uid = loki`).
- Dashboard `TOQ Server Logs` oferece:
  - Série temporal de eventos por severidade.
  - Painel de logs com filtro por ambiente (`$environment`).
  - Tabela de erros agregados (service/severity).
- Para correlação com Jaeger:
  - No painel de logs, abra um evento e clique na *derived field* **TraceID** → redireciona para a tela do Jaeger com o trace correspondente.
- Use Grafana Explore para consultas personalizadas, ex.: `count_over_time({service_name="toq_server", severity="ERROR"}[5m])`.

## Smoke Test Manual (checklist)
1. Com o stack ativo, gere requisições HTTP que produzam logs INFO e ERROR.
2. Verifique no Grafana Explore se os registros chegam com labels `service_name`, `environment`, `severity`, `trace_id` e `request_id`.
3. Clique no link **TraceID** para validar a navegação para o Jaeger.
4. Ajuste `LOKI_RETENTION_DAYS` para `3` temporariamente, recrie o serviço Loki (`docker compose up -d loki`) e confirme via API (`/loki/api/v1/status/buildinfo`) que o parâmetro foi aplicado.

## Logger e Compatibilidade
- O handler atual (`SplitLevelHandler`) continua responsável por stdout/stderr e arquivo NDJSON para retrocompatibilidade.
- O bootstrap anexa um segundo handler (`otelslog`) via *tee* interno, enviando cada `slog.Record` também ao Collector sem duplicar estrutura.
- Não há mudança de formato no arquivo em disco; ferramentas existentes seguem operando.

## Referências e Próximos Passos
- Configurações: `docker-compose.yml`, `loki-config.yaml`, `otel-collector-config.yaml`.
- Dashboard: `grafana/dashboard-files/toq-server-logs.json`.
- Próximos incrementos sugeridos:
  - Habilitar TLS no endpoint OTLP quando exposto fora da máquina local.
  - Criar regras de alerta Loki (ex.: `count_over_time` para picos de ERROR).
  - Automação IaC (Terraform) para replicar configuração em ambientes superiores.
