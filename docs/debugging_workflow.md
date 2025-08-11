# TOQ Server - Debugging Workflow Guide

## Overview
Esta documentação descreve o workflow completo de debugging e performance profiling do TOQ Server, incluindo configurações do VS Code e uso do pprof.

## VS Code Debugging Configurations

### 1. Development Mode with pprof
```json
"TOQ Server (Development with pprof)"
```
- **Propósito**: Desenvolvimento com profiling habilitado
- **Variáveis de Ambiente**:
  - `ENABLE_PPROF=true` - Habilita servidor pprof
  - `PPROF_PORT=6060` - Porta do servidor pprof
  - `LOG_LEVEL=debug` - Logs detalhados
  - `ENV=development` - Modo desenvolvimento
- **Terminal**: Integrado para monitoramento
- **Uso**: Debugging diário e análise de performance

### 2. Production Mode
```json
"TOQ Server (Production Mode)"
```
- **Propósito**: Simulação de ambiente de produção
- **Variáveis de Ambiente**:
  - `ENABLE_PPROF=false` - pprof desabilitado por segurança
  - `LOG_LEVEL=info` - Logs de produção
  - `ENV=production` - Modo produção
- **Uso**: Testes de deployment e validação final

### 3. Remote Debugging
```json
"Attach to running TOQ Server"
```
- **Propósito**: Attach a processo rodando externamente
- **Configuração**: Porta 2345 (padrão Delve)
- **Uso**: Debugging de containers ou processos remotos

## Performance Profiling com pprof

### Acessando o Servidor pprof
Quando `ENABLE_PPROF=true`, o servidor estará disponível em:
```
http://localhost:6060/debug/pprof/
```

### Perfis Disponíveis

#### 1. CPU Profiling
```bash
# Capturar perfil CPU por 30 segundos
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# Análise interativa
(pprof) top10
(pprof) web
(pprof) list function_name
```

#### 2. Memory Profiling
```bash
# Heap memory
go tool pprof http://localhost:6060/debug/pprof/heap

# Memory allocations
go tool pprof http://localhost:6060/debug/pprof/allocs
```

#### 3. Goroutine Analysis
```bash
# Visualizar goroutines
go tool pprof http://localhost:6060/debug/pprof/goroutine

# Stack traces de goroutines
curl http://localhost:6060/debug/pprof/goroutine?debug=2
```

#### 4. Blocking Analysis
```bash
# Análise de bloqueios
go tool pprof http://localhost:6060/debug/pprof/block

# Mutex contention
go tool pprof http://localhost:6060/debug/pprof/mutex
```

### Comandos pprof Úteis
```bash
# Top functions por CPU
(pprof) top10

# Gráfico web (requer Graphviz)
(pprof) web

# Lista código de função específica
(pprof) list function_name

# Disassembly de função
(pprof) disasm function_name

# Comparar dois perfis
go tool pprof -base profile1.pb.gz profile2.pb.gz
```

## Workflow de Debugging

### 1. Desenvolvimento Diário
1. Use "TOQ Server (Development with pprof)" no VS Code
2. Inicie debugging com F5
3. Monitore logs no terminal integrado
4. Acesse pprof em http://localhost:6060/debug/pprof/

### 2. Análise de Performance
```bash
# 1. Capture perfil durante carga
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=60

# 2. Análise interativa
(pprof) top10
(pprof) web

# 3. Salvar perfil para comparação
go tool pprof -proto http://localhost:6060/debug/pprof/profile > profile_$(date +%Y%m%d_%H%M%S).pb.gz
```

### 3. Debugging de Produção
1. Use "TOQ Server (Production Mode)" para simular produção
2. Valide comportamento sem pprof habilitado
3. Teste graceful shutdown com Ctrl+C

### 4. Debugging Remoto
1. Inicie servidor com debugger:
```bash
dlv debug cmd/toq_server.go --headless --listen=:2345 --api-version=2
```
2. Use "Attach to running TOQ Server" no VS Code

## Monitoramento em Tempo Real

### Endpoints de Health Check
```bash
# Status do servidor
curl http://localhost:8080/health

# Métricas Prometheus
curl http://localhost:9090/metrics

# pprof endpoints
curl http://localhost:6060/debug/pprof/
```

### OpenTelemetry Integration
- **Jaeger UI**: http://localhost:16686
- **Traces**: Visualização completa de requests
- **Spans**: Análise detalhada de latência

## Troubleshooting

### Problemas Comuns

#### pprof não acessível
1. Verifique se `ENABLE_PPROF=true`
2. Confirme porta `PPROF_PORT=6060`
3. Verifique logs para erros de inicialização

#### Debugging não conecta
1. Verifique se Delve está instalado: `go install github.com/go-delve/delve/cmd/dlv@latest`
2. Confirme porta 2345 disponível
3. Verifique firewall/network

#### Performance Issues
1. Use CPU profiling para identificar hotspots
2. Analise goroutines para vazamentos
3. Monitore heap para memory leaks

## Ferramentas Recomendadas

### VS Code Extensions
- Go (Google) - Suporte completo Go
- Go Nightly - Recursos experimentais
- Delve Debugger - Debugging avançado

### Ferramentas Externas
- Graphviz - Para gráficos pprof
- wrk/hey - Load testing
- Jaeger - Distributed tracing

## Configurações Avançadas

### Environment Variables Reference
```bash
# Core Configuration
ENV=development|production
LOG_LEVEL=debug|info|warn|error

# pprof Configuration
ENABLE_PPROF=true|false
PPROF_PORT=6060

# Database Configuration
MYSQL_HOST=localhost
MYSQL_PORT=3306
REDIS_HOST=localhost
REDIS_PORT=6379

# gRPC Configuration
GRPC_PORT=50051
HTTP_PORT=8080

# Observability
JAEGER_ENDPOINT=http://localhost:14268/api/traces
ENABLE_METRICS=true
```

### Performance Tuning
```bash
# Go runtime configuration
GOMAXPROCS=4
GOGC=100
GODEBUG=gctrace=1
```

## Best Practices

1. **Desenvolvimento**:
   - Sempre use development mode com pprof
   - Monitore goroutines regularmente
   - Faça profiling durante desenvolvimento

2. **Testing**:
   - Teste em production mode antes de deploy
   - Valide graceful shutdown
   - Execute load testing com profiling

3. **Produção**:
   - Mantenha pprof desabilitado
   - Configure logging adequado
   - Monitor métricas continuamente

4. **Debugging**:
   - Use breakpoints estrategicamente
   - Analise contexto completo
   - Documente findings

## Resources

- [Go pprof Documentation](https://golang.org/pkg/net/http/pprof/)
- [Delve Debugger Guide](https://github.com/go-delve/delve)
- [VS Code Go Extension](https://github.com/golang/vscode-go)
- [OpenTelemetry Go](https://opentelemetry.io/docs/instrumentation/go/)
