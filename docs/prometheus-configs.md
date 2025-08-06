# Configurações Alternativas do Prometheus

## Para Docker Compose (Padrão)
```yaml
scrape_configs:
  - job_name: 'prometheus'
    static_configs:
      - targets: ['prometheus:9090']  # Nome do serviço Docker

  - job_name: 'otel-collector'
    static_configs:
      - targets: ['otel-collector:8889']

  - job_name: 'toq_server_metrics'
    static_configs:
      - targets: ['172.17.0.1:4318']  # IP da bridge docker0
```

## Para Docker Desktop (Windows/Mac)
```yaml
scrape_configs:
  - job_name: 'prometheus'
    static_configs:
      - targets: ['prometheus:9090']

  - job_name: 'otel-collector'
    static_configs:
      - targets: ['otel-collector:8889']

  - job_name: 'toq_server_metrics'
    static_configs:
      - targets: ['host.docker.internal:4318']  # Docker Desktop
```

## Para Desenvolvimento Local (sem Docker)
```yaml
scrape_configs:
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']

  - job_name: 'otel-collector'
    static_configs:
      - targets: ['localhost:8889']

  - job_name: 'toq_server_metrics'
    static_configs:
      - targets: ['localhost:4318']
```

## Notas sobre Conectividade

### IP da Bridge Docker (172.17.0.1)
- Funciona na maioria dos ambientes Linux com Docker
- É o gateway padrão da rede bridge do Docker
- Permite que containers acessem serviços no host

### host.docker.internal
- Funciona no Docker Desktop (Windows/Mac)
- Pode não estar disponível em instalações Linux padrão
- DNS especial que resolve para o host interno

### Como descobrir o IP correto
```bash
# Ver a bridge docker0
ip addr show docker0

# Ver redes Docker
docker network ls
docker network inspect bridge

# Testar conectividade de dentro do container
docker exec -it <container_name> curl http://172.17.0.1:4318/metrics
```
