# Grafana Configuration - TOQ Server

Esta pasta contém as configurações para o Grafana no projeto TOQ Server.

## Estrutura

```
grafana/
├── datasources/           # Configuração automática de datasources
│   └── prometheus.yml     # Datasource do Prometheus
├── dashboards/           # Configuração de provisioning
│   └── dashboards.yml    # Define onde encontrar os dashboards
└── dashboard-files/      # Arquivos JSON dos dashboards
    ├── toq-server-otel.json    # Dashboard do OpenTelemetry Collector
    └── toq-server-app.json     # Dashboard da aplicação
```

## Datasources

### Prometheus
- **URL**: `http://prometheus:9090`
- **Access**: Proxy (através do Grafana)
- **Default**: Sim
- **Editable**: Sim

## Dashboards

### 1. TOQ Server - OpenTelemetry Monitoring (`toq-server-otel`)
Monitora o OpenTelemetry Collector:
- Taxa de métricas recebidas pelo collector
- Status dos serviços (UP/DOWN)
- Uso de memória do collector

### 2. TOQ Server - Application Metrics (`toq-server-app`)
Monitora a aplicação gRPC:
- **gRPC Request Rate**: Requisições por segundo
- **Success Rate**: Percentual de requisições bem-sucedidas
- **P95 Response Time**: Tempo de resposta do percentil 95
- **Error Rate**: Taxa de erros
- **Requests by Method**: Gráfico de requisições por método gRPC
- **Response Time Percentiles**: P50, P95, P99

## Métricas Utilizadas

### OpenTelemetry Collector
```promql
# Taxa de métricas recebidas
rate(otelcol_receiver_accepted_metric_points_total[5m])

# Status dos serviços
up

# Uso de memória
otelcol_process_memory_rss
```

### Aplicação gRPC
```promql
# Taxa de requisições
rate(grpc_server_handled_total{service_name="toq_server"}[5m])

# Taxa de sucesso
rate(grpc_server_handled_total{service_name="toq_server",grpc_code="OK"}[5m]) / 
rate(grpc_server_handled_total{service_name="toq_server"}[5m]) * 100

# Percentis de tempo de resposta
histogram_quantile(0.95, rate(grpc_server_handling_seconds_bucket{service_name="toq_server"}[5m]))

# Taxa de erros
rate(grpc_server_handled_total{service_name="toq_server",grpc_code!="OK"}[5m])
```

## Personalização

### Adicionando Novos Dashboards
1. Exporte o dashboard do Grafana como JSON
2. Coloque o arquivo em `dashboard-files/`
3. Reinicie o Grafana ou aguarde o reload automático (10s)

### Modificando Datasources
1. Edite `datasources/prometheus.yml`
2. Reinicie o container do Grafana

### Alterando Credenciais
No `docker-compose.yml`:
```yaml
environment:
  - GF_SECURITY_ADMIN_USER=seu_usuario
  - GF_SECURITY_ADMIN_PASSWORD=sua_senha
```

## Acesso

- **URL**: http://localhost:3000
- **Usuário**: admin
- **Senha**: admin

## Troubleshooting

### Dashboard não aparece
- Verifique os logs: `docker-compose logs grafana`
- Confirme que o arquivo JSON está válido
- Verifique permissões da pasta `dashboard-files`

### Métricas não aparecem
- Verifique se o Prometheus está coletando as métricas
- Confirme que o datasource está configurado corretamente
- Teste queries diretamente no Prometheus primeiro

### Grafana não inicia
- Verifique se a porta 3000 não está em uso
- Confirme que os volumes estão montados corretamente
- Verifique logs de inicialização
