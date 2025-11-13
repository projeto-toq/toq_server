### Engenheiro de Software Go Sênior — Análise e Refatoração TOQ Server

**Objetivo:** Atuar como engenheiro Go sênior para analisar código existente, identificar desvios das regras do projeto e propor planos detalhados de refatoração/implementação. Toda a interação deve ser feita em português.

---

## 🎯 Problema / Solicitação

Após a refatoraçÃo onde repository responde sql.ErrNoRows, foi identificado que o endpoint de confirmaçÃo de troca de email estÃ falhando com erro 500 quando o usuÃrio nÃo Ã encontrado. Logs de erro:
{"time":"2025-11-12T16:15:57.239971148Z","level":"ERROR","msg":"user.confirm_email_change.stage_error","request_id":"b7a2173a-eb7e-441a-9346-ad84ff927e0e","stage":"update_user","err":"sql: no rows in result set"}
{"time":"2025-11-12T16:15:57.244005132Z","level":"ERROR","msg":"HTTP Error","request_id":"b7a2173a-eb7e-441a-9346-ad84ff927e0e","request_id":"b7a2173a-eb7e-441a-9346-ad84ff927e0e","method":"POST","path":"/api/v2/user/email/confirm","status":500,"duration":12936655,"size":61,"client_ip":"79.55.106.238","user_agent":"PostmanRuntime/7.50.0","user_id":6,"user_role_id":6,"function":"github.com/projeto-toq/toq_server/internal/core/utils.InternalError","file":"/codigos/go_code/toq_server/internal/core/utils/http_errors.go","line":248,"stack":["github.com/projeto-toq/toq_server/internal/core/utils.InternalError (http_errors.go:248)"],"error_code":500,"error_message":"Failed to update user","errors":["HTTP 500: Failed to update user"]}

sql.ErrnoRows deveria ser tratado pelo service para determinar se é erro ou não.

Assim:
1. Analise o código onde o erro ocorreu e identifique a causa raiz do problema.
2. Verifique outros pontos service onde sql.ErrNoRows pode estar sendo mal interpretado.

---

## 📘 Fonte da Verdade

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