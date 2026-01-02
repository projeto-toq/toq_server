### Engenheiro de Software Go Sênior — Análise e Refatoração TOQ Server

**Objetivo:** Atuar como engenheiro Go sênior para analisar código existente, entender claramente o erro apresentado e identificar a causa raiz do problema para propor planos detalhados de refatoração. Toda a interação deve ser feita em português.

---

## 🎯 Problema / Solicitação

Ao aprovar uma visita ao imovel através do endpoint `POST /visits/status` com o body:

```json
{
  "action": "APPROVE",
  "visitId": 1,
  "notes": "Owner approved without constraints"
}
```

recebo como resposta:
```json
{
    "code": 500,
    "details": null,
    "message": "Internal server error"
}
```

e estas entraddas no log:

```json
{"time":"2026-01-02T14:59:12.921933019Z","level":"ERROR","msg":"mysql.listing.update_owner_response_stats.exec_error","request_id":"46accc34-3f13-401c-86a3-14b2fd4ccf22","error":"Error 1264 (22003): Out of range value for column 'owner_avg_response_time_seconds' at row 1","listing_identity_id":2}
{"time":"2026-01-02T14:59:12.922107853Z","level":"ERROR","msg":"visit.approve.owner_response_metrics_error","request_id":"46accc34-3f13-401c-86a3-14b2fd4ccf22","visit_id":1,"err":"update owner response stats: Error 1264 (22003): Out of range value for column 'owner_avg_response_time_seconds' at row 1"}
{"time":"2026-01-02T14:59:12.922968288Z","level":"ERROR","msg":"HTTP Error","request_id":"46accc34-3f13-401c-86a3-14b2fd4ccf22","request_id":"46accc34-3f13-401c-86a3-14b2fd4ccf22","method":"POST","path":"/api/v2/visits/status","status":500,"duration":9862177,"size":61,"client_ip":"78.152.121.74","user_agent":"PostmanRuntime/7.51.0","user_id":2,"user_role_id":6,"errors":["update owner response stats: Error 1264 (22003): Out of range value for column 'owner_avg_response_time_seconds' at row 1"]}
```


Assim:
1. Analise o guia do projeto `docs/toq_server_go_guide.md`, o código atual e identifique a causa raiz do problema
2. Proponha um plano detalhado de refatoração com code skeletons para corrigir o problema, seguindo estritamente as regras de arquitetura do manual `docs/toq_server_go_guide.md` (observabilidade, erros, transações, etc).

---

**TODAS as regras de arquitetura, padrões de código, observabilidade e documentação estão em:**
- **`docs/toq_server_go_guide.md`** — Guia completo do projeto (seções 1-17)
- **`README.md`** — Configurações de ambiente e observabilidade

**⚠️ Consulte SEMPRE esses documentos antes de propor qualquer solução.**

---

## 🎯 Processo de Trabalho

1. **Leia o código** envolvido (adapters, services, handlers, entities, converters)
2. **Identifique a causa raiz** apresente evidencias no código
3. **Proponha plano detalhado** com code skeletons
4. **Não implemente código** — apenas análise e planejamento

---

## 📋 Formato do Plano

### 1. Diagnóstico
- Lista de arquivos analisados
- Causa raiz identificada (apresente evidencias no código)
- Impacto de cada desvio/problema
- Melhorias possíveis

### 2. Code Skeletons
Para cada arquivo novo/alterado, forneça **esqueletos** conforme templates da **Seção 8 do guia**:
- **Handlers:** Assinatura + Swagger completo (sem implementação)
- **Services:** Assinatura + Godoc + estrutura tracing/transação
- **Repositories:** Assinatura + Godoc + query + InstrumentedAdapter
- **DTOs:** Struct completa com tags e comentários
- **Entities:** Struct completa com sql.Null* quando aplicável
- **Converters:** Lógica completa de conversão

### 3. Estrutura de Diretórios
Mostre organização final seguindo **Regra de Espelhamento (Seção 2.1 do guia)**

### 4. Ordem de Execução
Etapas numeradas com dependências

### 5. Checklist de Conformidade
Valide contra **seções específicas do guia**:
- [ ] Arquitetura hexagonal (Seção 1)
- [ ] Regra de Espelhamento Port ↔ Adapter (Seção 2.1)
- [ ] InstrumentedAdapter em repos (Seção 7.3)
- [ ] Transações via globalService (Seção 7.1)
- [ ] Tracing/Logging/Erros (Seções 5, 7, 9)
- [ ] Documentação (Seção 8)
- [ ] Sem anti-padrões (Seção 14)

---

## 🚫 Restrições

### Permitido (ambiente dev)
- Alterações disruptivas, quebrar compatibilidade, alterar assinaturas

### Proibido
- ❌ Criar/alterar testes unitários
- ❌ Scripts de migração de dados
- ❌ Editar swagger.json/yaml manualmente
- ❌ Executar git/go test
- ❌ Mocks ou soluções temporárias

---

## 📝 Documentação

- **Código:** Inglês (seguir Seção 8 do guia)
- **Plano:** Português (citar seções do guia ao justificar)
- **Swagger:** `make swagger` (anotações no código)