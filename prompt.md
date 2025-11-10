### Engenheiro de Software Go Sênior — Análise e Refatoração TOQ Server

**Objetivo:** Atuar como engenheiro Go sênior para analisar código existente, identificar desvios das regras do projeto e propor planos detalhados de refatoração/implementação. Toda a interação deve ser feita em português.

---

## 🎯 Problema / Solicitação

após a refatoração do sistema de gestão de usuários e permissionamento, o lint reporta:
internal/core/service/user_service/assign_role_to_user.go:129:42: Error return value of `us.permissionService.InvalidateUserCache` is not checked (errcheck)
        us.permissionService.InvalidateUserCache(ctx, userID)
                                                ^
internal/core/service/user_service/remove_role_from_user.go:91:42: Error return value of `us.permissionService.InvalidateUserCache` is not checked (errcheck)
        us.permissionService.InvalidateUserCache(ctx, userID)
                                                ^
internal/core/service/user_service/switch_active_role.go:79:42: Error return value of `us.permissionService.InvalidateUserCache` is not checked (errcheck)
        us.permissionService.InvalidateUserCache(ctx, userID) // TODO incluir mensagem, "switch_active_role_with_tx")
                                                ^
make: *** [Makefile:39: ci-lint] Error 1


Assim:
1. Analise os codigos de user_model, user_service, user_repository, permission_model, permission_service, permission_repository mapeando a causa raiz do problema.
2. Proponha um plano detalhado para corrigir o problema.
3. Existe um TODO que é necessário incluir a mensagem da causa da invalidação do cache. Isto se deve a refatoração que moveu a responsabilidade de gestão de user_roles para user_service. Proponha como incluir essa mensagem em cada chamada de invalidação de cache.
    3.1. Revise as assinaturas das funções de invalidação de cache em permission_service e permission_repository, garantindo que aceitem um parâmetro adicional para a mensagem.
    3.2. Revise todas as chamadas para essas funções em user_service, garantindo que a mensagem apropriada seja passada com base na operação realizada (ex: "assign_role", "remove_role", "switch_active_role").


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