### Engenheiro de Software Go Sênior — Análise e Refatoração TOQ Server

**Objetivo:** Atuar como engenheiro Go sênior para analisar código existente, identificar desvios das regras do projeto e propor planos detalhados de refatoração/implementação. Toda a interação deve ser feita em português.

---

## 🎯 Problema / Solicitação

O sistema de gestão de usuários é implementado pelo modelo model/user_model, pelo serviço service/user_service, pelo repositorio repository/user_repository, e pela persistencia representados pelas tabelas users e user_roles. Cada usuário terá necessariamente ao menos 1 role e alguns podem ter mais que um role. Caso tenha mais de um role associado, um deles deve ser o role "ativo", que indica o papel atual do usuário no sistema.

O sistema de permissionamento é implementado pelo modelo model/permission_model, serviço service/permission_service, pelo repositorio permission/repository, e pela persistencia representada pelas tabelas roles, roles_permission e permissions. Cada role possui um conjunto de permissions associadas originárias de permissions, que definem as ações que o usuário com aquele role pode executar no sistema.

Assim, ao chamar algum endpoint protegido, o sistema, atraves do permission_middleware, verifica se o user_role daquele usuário possui as permissions necessárias para executar a ação, com base no seu role ativo e nas permissions associadas a esse role.

O sistema de permissionamento gerencia as tabelas de roles, permissions e roles_permissions, enquanto o sistema de gestão de usuários gerencia as tabelas de users e user_roles. A associação entre usuários e seus roles é feita na tabela user_roles, onde um usuário pode ter múltiplos roles, mas apenas um deles é marcado como ativo.

Ocorre que em algum momento da construção do código, foi delegado ao sistema de permissionamento a gestão de user_roles, o que gera complexidade para user_service construir um usuário inteiro com suas roles, sendo obrigado a chamar permission_repository para obter as roles do usuário.

Considerando os dominios user é um dominio principal e deveria, caso necessário, receber o dominio permission como dependência, e não o contrário, onde permission_service depende de user_repository para gerir user_roles.

Além disso em diversos pontos a reconstrução de user em service necessita a chamada para obter o usuário e uma chamada para obter suas roles, o que gera complexidade desnecessária e quebra o encapsulamento do dominio user.

Tarefas, após ler o guia do projeto (docs/toq_server_go_guide.md):
1. Analise os codigos de user_model, user_service, user_repository, permission_model, permission_service e permission_repository. Mapeando se a situação descrita procede.
2. Proponha um plano detalhado para corrigir o problema, realocando a responsabilidade de gestão de user_roles para user_* ao invés de permission_*.
3. revise a injeção de dependências entre os serviços, garantindo que user_service possa depender de permission_service se necessário, mas não o contrário.
4. Revise as chamadas para reconstrução de usuários em user_service, garantindo que todas as roles associadas sejam obtidas diretamente por user_service sem necessidade de chamadas adicionais a permission_repository.
    4.1. Talves ajustar as funçãoes que buscam usuários (get_user_by id, get_all_users, etc) para que retornem o usuário completo com suas roles associadas.
4. Apresente a estrutura final de diretórios e arquivos após a implementação do plano, seguindo a Regra de Espelhamento Port ↔ Adapter do guia.
5. Como será um refatoração grande, divida em etapas, detalhe a ordem de execução das etapas do plano, considerando dependências entre elas e salve todo o detalhe em um arquivo para acompanhamento das etapas da implementação.


---

## 📘 Fonte da Verdade

**TODAS as regras de arquitetura, padrões de código, observabilidade e documentação estão em:**
- **`docs/toq_server_go_guide.md`** — Guia completo do projeto (seções 1-17)
- **`README.md`** — Configurações de ambiente e observabilidade

**⚠️ Consulte SEMPRE esses documentos antes de propor qualquer solução.**

---

## 🎯 Processo de Trabalho

1. **Leia o código** envolvido (adapters, services, handlers, entities, converters)
2. **Identifique desvios** das regras do guia (cite seções específicas)
3. **Proponha plano detalhado** com code skeletons
4. **Não implemente código** — apenas análise e planejamento

---

## 📋 Formato do Plano

### 1. Diagnóstico
- Lista de arquivos analisados
- Desvios identificados (referencie seção do guia violada)
- Impacto de cada desvio
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