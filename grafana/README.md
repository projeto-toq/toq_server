# Grafana Configuration - TOQ Server

Esta pasta contém as configurações para o Grafana no projeto TOQ Server.

## Estrutura

```
grafana/
├── datasources/           # Configuração automática de datasources
│   └── datasource.yml     # Datasource do Prometheus
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

### Acesso Local (no servidor Linux)
- **URL**: http://localhost:3000
- **Usuário**: admin
- **Senha**: admin

### Acesso Remoto (de outras máquinas na rede)
- **URL**: http://192.168.10.137:3000
- **Usuário**: admin
- **Senha**: admin

> **Nota**: Se você configurou redirecionamento no arquivo hosts (`C:\Windows\System32\drivers\etc\hosts` no Windows), 
> pode usar `localhost:3000` mesmo de máquinas remotas.

## Troubleshooting

### Dashboard não aparece
- Verifique os logs: `docker-compose logs grafana`
- Confirme que o arquivo JSON está válido
- Verifique permissões da pasta `dashboard-files`

### Métricas não aparecem
- Verifique se o Prometheus está coletando as métricas
- Confirme que o datasource está configurado corretamente
- Teste queries diretamente no Prometheus primeiro

### Acesso remoto não funciona
- Verifique se o firewall do Linux permite conexões na porta 3000
- Confirme que o Docker está fazendo bind em `0.0.0.0:3000` (não apenas `127.0.0.1`)
- Teste conectividade: `telnet 192.168.10.137 3000`
- Se usar arquivo hosts no Windows, verifique: `C:\Windows\System32\drivers\etc\hosts`

### Confusão entre múltiplos serviços na porta 3000
Se `localhost:3000` no Windows mostra um Grafana diferente do esperado:
```powershell
# Verifique se há serviços locais na porta 3000
netstat -ano | findstr :3000

# Identifique o processo servidor (LISTENING)
tasklist /FI "PID eq NUMERO_DO_PID_LISTENING"
```

**Causas comuns:**
- **VS Code**: Live Server, Preview, ou extensões de desenvolvimento
- **Docker Desktop**: Port forwarding automático
- **Node.js/React**: Servidor de desenvolvimento local
- **WSL2**: Redirecionamento de portas do Linux para Windows
- **Grafana local**: Instalação separada do Grafana

**Identificação:**
```powershell
# Exemplo de saída
Code.exe                      1544    # VS Code servidor
chrome.exe                   13240    # Browser cliente
```

**Solução:** Sempre use o IP direto `http://192.168.10.137:3000` para acessar o Grafana do projeto.

### Grafana não inicia
- Verifique se a porta 3000 não está em uso
- Confirme que os volumes estão montados corretamente
- Verifique logs de inicialização
