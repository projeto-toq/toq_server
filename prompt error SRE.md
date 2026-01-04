### SRE Sênior — Análise e Refatoração TOQ Server

**Objetivo:** Atuar como SRE sênior para analisar código existente, entender claramente o erro apresentado e identificar a causa raiz do problema para propor planos detalhados de refatoração. Toda a interação deve ser feita em português.

---

## 🎯 Problema / Solicitação

Atualmente o Grafana é usado para apresentar Dashboards de observabilidade do TOQ Server.
Existem 2 Dashboards que não estão funcionando corretamente:
1. **Dashboard TOQ Server - Logs:** Apresenta os dados do Log estruturado, mas não possue uma forma de redirecionar diretamente ao Dashboard de Traces, baseado em request_id ou trace_id.
2. **Dashboard TOQ Server - Traces:** Apresenta os dados do traces, mas não permite a correlação direta com os logs, baseado em request_id ou trace_id.
Todos os componentes de observabilidade estão em docker `/codigos/go_code/toq_server/docker-compose.yml`.

Assim:
1. Analise o guia do projeto `docs/toq_server_go_guide.md`, o código atual, as configuraçoes, os dashboards atuais e identifique a causa raiz do problema
2. Proponha um plano detalhado de refatoração com code skeletons para corrigir o problema, seguindo estritamente as regras de arquitetura do manual `docs/toq_server_go_guide.md` (observabilidade, erros, transações, etc).
3. Garanta que existe uma funcionalide de nas versões em uso dos utilitários de observabilidade que permitam a solução proposta (ex: correlação logs/traces via request_id/trace_id)
4. Ao final do plano deve haver uma atualização de `/codigos/go_code/toq_server/docs/observability/sre_guide.md`, readme.md e guia do projeto para refletir as mudanças propostas.
---

**TODAS as regras de arquitetura, padrões de código, observabilidade e documentação estão em:**
- **`docs/toq_server_go_guide.md`** — Guia completo do projeto (seções 1-17)
- **`README.md`** — Configurações de ambiente e observabilidade

**⚠️ Consulte SEMPRE esses documentos antes de propor qualquer solução.**

---

## 🎯 Processo de Trabalho

1. **Leia o código e as configurações reais de containers** envolvido (adapters, services, handlers, entities, converters)
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