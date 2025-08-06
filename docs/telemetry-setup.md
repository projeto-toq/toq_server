# Configuração de Telemetria - TOQ Server

Este documento descreve como a telemetria está configurada no projeto TOQ Server usando OpenTelemetry, Prometheus e Docker Compose.

## Arquitetura de Telemetria

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   TOQ Server    │───▶│ OTEL Collector   │───▶│   Prometheus    │───▶│    Grafana      │
│   (Go App)      │    │                  │    │                 │    │                 │
│   Port: 4318    │    │ Ports: 4317,8889 │    │   Port: 9091    │    │   Port: 3000    │
└─────────────────┘    └──────────────────┘    └─────────────────┘    └─────────────────┘
        │                        │                        │                        │
        │                        │                        │                        │
        ▼                        ▼                        ▼                        ▼
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   App Metrics   │    │  Traces & Logs   │    │  Time Series    │    │   Dashboards    │
│   /metrics      │    │     Debug        │    │     Storage     │    │  Visualization  │
└─────────────────┘    └──────────────────┘    └─────────────────┘    └─────────────────┘
```

## Componentes

### 1. TOQ Server (Aplicação Go)
- **Traces**: Enviados via OTLP gRPC para o collector (porta 4317)
- **Metrics**: Expostos diretamente via Prometheus handler (porta 4318)
- **Resource**: `service.name=toq_server`, `service.version=v2.1-grpc`

### 2. OpenTelemetry Collector
- **Receivers**: OTLP gRPC (4317) e HTTP (4318)
- **Processors**: batch, memory_limiter
- **Exporters**: 
  - debug (logs detalhados)
  - prometheus (métricas na porta 8889)
  - logging (info level)

### 3. Prometheus
- **Scrape Intervals**: 15s (global), 10s (específicos)
- **Targets**:
  - `otel-collector:8889` - Métricas do collector
  - `172.17.0.1:4318` - Métricas diretas da aplicação (bridge docker0)
  - `prometheus:9090` - Auto-monitoramento

### 4. Grafana
- **Datasource**: Prometheus (configurado automaticamente)
- **Dashboards**: 
  - TOQ Server - OpenTelemetry Monitoring
  - TOQ Server - Application Metrics
- **Credentials**: admin/admin (padrão)
- **Provisioning**: Datasources e dashboards configurados automaticamente

## Portas Utilizadas

| Serviço | Porta | Protocolo | Descrição |
|---------|--------|-----------|-----------|
| TOQ Server | 4318 | HTTP | Prometheus metrics endpoint |
| OTEL Collector | 4317 | gRPC | OTLP receiver |
| OTEL Collector | 8888 | HTTP | Collector internal metrics |
| OTEL Collector | 8889 | HTTP | Prometheus exporter |
| Prometheus | 9091 | HTTP | Prometheus UI |
| Grafana | 3000 | HTTP | Grafana UI |

## Configuração do Código Go

```go
// Em initialize_telemetry.go
otlpEndpoint := "otel-collector:4317"  // Nome do serviço Docker
resource.WithAttributes(
    semconv.ServiceNameKey.String("toq_server"),
    semconv.ServiceVersionKey.String("v2.1-grpc"),
)
```

## Como Testar

### 1. Iniciar os serviços:
```bash
docker-compose up -d prometheus otel-collector grafana
```

### 2. Executar o script de teste:
```bash
./scripts/test-telemetry.sh
```

### 3. Verificar endpoints:
- Prometheus UI: http://localhost:9091
- Grafana UI: http://localhost:3000 (admin/admin)
- OTEL Collector metrics: http://localhost:8889/metrics
- Application metrics: http://localhost:4318/metrics

### 4. Verificar métricas no Prometheus:
- Acesse http://localhost:9091
- Procure por métricas com prefixo `toq_server_` ou `otel_`
- Verifique se os targets estão "UP" na aba Status > Targets

### 5. Verificar dashboards no Grafana:
- Acesse http://localhost:3000 (admin/admin)
- Os dashboards estarão disponíveis na pasta "TOQ Server"
- Dashboards incluídos:
  - **TOQ Server - OpenTelemetry Monitoring**: Métricas do collector
  - **TOQ Server - Application Metrics**: Métricas da aplicação gRPC

## Métricas Disponíveis

### Do OpenTelemetry Collector:
- `otelcol_receiver_*` - Métricas do receiver
- `otelcol_processor_*` - Métricas dos processadores
- `otelcol_exporter_*` - Métricas dos exporters

### Da Aplicação (quando implementadas):
- `toq_server_requests_total` - Total de requests
- `toq_server_request_duration` - Duração dos requests
- `toq_server_grpc_*` - Métricas específicas do gRPC

## Resolução de Problemas

### 1. Collector não recebe dados:
- Verificar se a aplicação está enviando para `otel-collector:4317`
- Verificar logs do collector: `docker-compose logs otel-collector`

### 2. Prometheus não coleta métricas:
- Verificar configuração em `prometheus.yml`
- Verificar status dos targets na UI do Prometheus

### 3. Aplicação não expõe métricas:
- Verificar se a porta 4318 está sendo exposta
- Verificar se o handler `/metrics` está configurado

### 4. Grafana não mostra dados:
- Verificar se o datasource Prometheus está configurado
- Verificar se as métricas existem no Prometheus primeiro
- Verificar se os dashboards foram provisionados corretamente

## Variáveis de Ambiente

```bash
# Para a aplicação
OTLP_ENDPOINT=otel-collector:4317

# Para desenvolvimento local (fora do Docker)
OTLP_ENDPOINT=localhost:4317
```

## Estrutura de Arquivos

```
├── docker-compose.yml              # Configuração dos serviços
├── prometheus.yml                  # ⚠️  CONFIGURAÇÃO DO PROMETHEUS (serviço)
├── otel-collector-config.yaml      # Configuração do OpenTelemetry Collector
├── grafana/
│   ├── datasources/
│   │   └── prometheus.yml          # ⚠️  CONFIGURAÇÃO DO GRAFANA (datasource)
│   ├── dashboards/
│   │   └── dashboards.yml          # Configuração de provisioning dos dashboards
│   └── dashboard-files/
│       ├── toq-server-otel.json    # Dashboard de monitoramento do OTEL Collector
│       └── toq-server-app.json     # Dashboard de métricas da aplicação
├── docs/
│   └── prometheus-configs.md       # Configurações alternativas do Prometheus
└── scripts/
    └── test-telemetry.sh           # Script de teste dos serviços
```

### ⚠️ Importante: Dois arquivos prometheus.yml diferentes

1. **`/prometheus.yml`** - Configuração do **serviço Prometheus**
   - Define quais targets coletar métricas
   - Intervalos de scraping
   - Jobs de monitoramento

2. **`/grafana/datasources/prometheus.yml`** - Configuração do **datasource no Grafana**
   - Define como o Grafana se conecta ao Prometheus
   - Configurações de proxy e timeouts
   - Provisioning automático

## Dashboards Disponíveis

### 1. TOQ Server - OpenTelemetry Monitoring
- **Métricas do Collector**: Taxa de métricas recebidas
- **Status dos Serviços**: Verificação de saúde
- **Uso de Memória**: Consumo de recursos do collector

### 2. TOQ Server - Application Metrics
- **gRPC Request Rate**: Requisições por segundo
- **Success Rate**: Percentual de sucessos
- **Response Time**: Percentis P50, P95, P99
- **Error Rate**: Taxa de erros
- **Requests by Method**: Distribuição por método gRPC

## Configuração Automática

O Grafana está configurado com **provisioning automático**:
- **Datasources**: Prometheus configurado automaticamente
- **Dashboards**: Carregados automaticamente da pasta `dashboard-files`
- **Credenciais**: admin/admin (configurável via environment variables)

## Próximos Passos

1. ✅ **Adicionado**: Grafana para visualização avançada
2. Configurar alertas no Prometheus e Grafana
3. Implementar métricas customizadas na aplicação
4. Configurar Jaeger para traces distribuídos (se necessário)
5. Adicionar dashboards específicos para business metrics
6. Configurar notificações de alertas (Slack, email, etc.)
