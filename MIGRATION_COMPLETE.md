# Migração para Grafana Alloy + Tempo + Loki - CONCLUÍDA

## ✅ IMPLEMENTAÇÃO COMPLETA

Data: 4 de novembro de 2025

### Alterações Realizadas

#### 1. Nova Estrutura de Diretórios
- ✅ Criado `alloy/` para configurações do Grafana Alloy
- ✅ Criado `tempo/` para configurações do Grafana Tempo

#### 2. Configurações Criadas/Atualizadas
- ✅ **alloy/config.alloy**: Configuração completa do Alloy (HTTP only, sem gRPC)
  - Receiver OTLP HTTP (:4318)
  - Processadores de enriquecimento de atributos
  - Exporters para Tempo, Loki e Prometheus
  - Scraping de MySQL e Redis exporters
  
- ✅ **tempo/tempo-config.yaml**: Configuração do Tempo (HTTP only)
  - Receiver OTLP HTTP
  - Metrics generator com RED metrics automáticas
  - Armazenamento local com retenção de 7 dias
  
- ✅ **loki-config.yaml**: Atualizado com structured metadata e retention habilitado

- ✅ **prometheus.yml**: Simplificado (apenas self-monitoring, scraping via Alloy)

- ✅ **configs/env.yaml**: Endpoint OTLP atualizado para `alloy:4318`

#### 3. Datasources do Grafana
- ✅ **grafana/datasources/tempo.yml**: Criado com correlação automática
  - tracesToLogs configurado
  - tracesToMetrics configurado
  - nodeGraph habilitado
  
- ✅ **grafana/datasources/loki.yml**: Atualizado
  - Derived field aponta para Tempo (antes era Jaeger)
  - maxLines aumentado para 5000
  
- ❌ **grafana/datasources/jaeger.yml**: REMOVIDO

#### 4. Docker Compose
- ✅ **docker-compose.yml**: Totalmente reescrito
  - Serviço `alloy` adicionado
  - Serviço `tempo` adicionado
  - Serviço `otel-collector` REMOVIDO
  - Serviço `jaeger` REMOVIDO
  - Prometheus com flag `--web.enable-remote-write-receiver`
  - Grafana com feature toggles para TraceQL

#### 5. Dashboards
- ✅ Script `scripts/migrate_dashboards.sh` criado e executado
- ✅ Dashboard `toq-server-logs-traces.json` migrado (Jaeger → Tempo)
- ✅ Backup criado em `grafana/dashboard-files_backup_20251104_174901/`

#### 6. Arquivos Obsoletos Removidos
- ❌ `otel-collector-config.yaml`: DELETADO
- ❌ `grafana/datasources/jaeger.yml`: DELETADO

#### 7. Documentação
- ✅ **docs/observability/sre_guide.md**: Completamente reescrito
  - Arquitetura atualizada para Alloy + Tempo
  - Comandos de troubleshooting específicos
  - Fluxo de correlação automática documentado
  - Checklists de validação

---

## 🚀 PRÓXIMOS PASSOS (VALIDAÇÃO)

### Fase 1: Iniciar Nova Stack

```bash
cd /codigos/go_code/toq_server

# Parar serviços obsoletos (se ainda rodando)
docker compose down otel-collector jaeger

# Iniciar infraestrutura base
docker compose up -d mysql redis mysql-exporter redis-exporter

# Iniciar backends de observabilidade
docker compose up -d prometheus loki tempo

# Iniciar Alloy (aguardar backends estarem prontos)
sleep 5
docker compose up -d alloy

# Iniciar Grafana
docker compose up -d grafana
```

### Fase 2: Validar Health dos Serviços

```bash
# Alloy
curl http://localhost:12345/ready
# Esperado: HTTP 200

# Tempo
curl http://localhost:3200/ready
# Esperado: ready

# Loki
curl http://localhost:3100/ready
# Esperado: ready

# Prometheus
curl http://localhost:9091/-/ready
# Esperado: HTTP 200

# Grafana
curl http://localhost:3000/api/health
# Esperado: {"database":"ok","version":"..."}
```

### Fase 3: Verificar Logs dos Serviços

```bash
# Ver logs do Alloy (procurar por erros de parsing)
docker logs alloy | tail -50

# Ver logs do Tempo
docker logs tempo | tail -50

# Ver logs do Loki
docker logs loki | tail -50
```

### Fase 4: Reiniciar Aplicação Go

```bash
# A aplicação precisa ser reiniciada para conectar ao novo endpoint (alloy:4318)
# Se rodando em Docker:
docker compose restart toq-server

# Se rodando no host:
# Parar processo atual e executar:
go run cmd/toq_server.go
```

### Fase 5: Gerar Tráfego e Validar Telemetria

```bash
# Fazer algumas requisições HTTP à API
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/...

# Verificar métricas no Prometheus
curl -s "http://localhost:9091/api/v1/query?query=http_requests_total" | jq

# Verificar logs no Loki
curl -G -s "http://localhost:3100/loki/api/v1/query" \
  --data-urlencode 'query={service_name="toq_server"}' \
  --data-urlencode 'limit=5' | jq
```

### Fase 6: Validar Grafana UI

1. Acesse http://localhost:3000 (admin / Senh@123)

2. **Validar Datasources:**
   - Configuration → Data Sources
   - Verificar `Tempo` (verde)
   - Verificar `Loki` (verde)
   - Verificar `Prometheus` (verde)
   - Jaeger NÃO deve aparecer na lista

3. **Validar Dashboards:**
   - Dashboards → Browse
   - Abrir "TOQ Server - Aplicação": deve mostrar métricas
   - Abrir "TOQ Server - Logs e Traces": deve ter logs com trace_id
   - Clicar em um `trace_id` → deve abrir Tempo (não Jaeger)

4. **Testar Explore:**
   - Explore → Selecionar `Tempo`
   - Buscar por trace recente
   - Clicar em span → Botão "Logs for this span" deve aparecer
   - Clicar → Deve abrir Loki com logs correlacionados

5. **Testar Correlação Logs → Traces:**
   - Explore → Selecionar `Loki`
   - Query: `{service_name="toq_server"}`
   - Clicar em linha de log
   - Link "View Trace in Tempo" deve aparecer e funcionar

### Fase 7: Validar Métricas RED do Tempo

```bash
# Tempo gera métricas RED automaticamente
# Verificar se estão sendo enviadas ao Prometheus:
curl -s "http://localhost:9091/api/v1/query?query=traces_spanmetrics_calls_total" | jq
```

---

## 🔍 CHECKLIST FINAL

- [ ] Alloy UI acessível e componentes saudáveis (http://localhost:12345)
- [ ] Tempo retorna traces via API
- [ ] Loki retorna logs com labels `trace_id`
- [ ] Prometheus recebe métricas via remote write do Alloy
- [ ] Grafana mostra 3 datasources (Tempo, Loki, Prometheus)
- [ ] Dashboards carregam sem erros
- [ ] Correlação Logs → Traces funciona (clicar em trace_id abre Tempo)
- [ ] Correlação Traces → Logs funciona (botão "Logs for this span")
- [ ] Métricas RED aparecem nos traces do Tempo
- [ ] Exporters MySQL e Redis sendo scrapedos pelo Alloy

---

## 📋 TROUBLESHOOTING COMUM

### Alloy não inicia
```bash
docker logs alloy
# Procurar por parsing errors na config River
```

### Tempo não recebe traces
```bash
# Verificar se aplicação Go está conectando ao Alloy
docker logs alloy | grep otlp

# Verificar se Alloy está exportando para Tempo
docker logs alloy | grep tempo
```

### Correlação não funciona
```bash
# Verificar derived fields no datasource Loki
# Grafana UI → Configuration → Data Sources → Loki → Derived Fields
# Deve ter entry para trace_id apontando para datasource 'tempo'
```

### Prometheus não recebe métricas
```bash
# Verificar flag remote-write-receiver
docker inspect prometheus | grep enable-remote

# Verificar logs do Alloy
docker logs alloy | grep prometheus
```

---

## 🎯 BENEFÍCIOS DA NOVA STACK

✅ **Redução de Complexidade**: 3 componentes backend (vs. 4 anteriormente)  
✅ **Correlação Nativa**: Sem regex frágil, Tempo injeta links automaticamente  
✅ **Configuração Unificada**: 1 arquivo River (Alloy) vs. 3 YAMLs separados  
✅ **Métricas RED Automáticas**: Geradas a partir de spans, sem instrumentação extra  
✅ **100% HTTP**: Eliminação total de gRPC (mais simples para REST API)  
✅ **Pronto para Produção**: Tempo com armazenamento persistente e retention  
✅ **Service Discovery**: Alloy preparado para ambientes dinâmicos (K8s/Swarm)  

---

## 📚 REFERÊNCIAS

- Documentação completa: `docs/observability/sre_guide.md`
- Guia do projeto: `docs/toq_server_go_guide.md`
- Configuração Alloy: `alloy/config.alloy`
- Configuração Tempo: `tempo/tempo-config.yaml`
- Script de migração de dashboards: `scripts/migrate_dashboards.sh`

---

**Status Final**: ✅ PRONTO PARA VALIDAÇÃO E USO
