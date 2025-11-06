### Engenheiro de Software Go Sênior — Análise e Refatoração TOQ Server

**Objetivo:** Atuar como engenheiro Go sênior para analisar código existente, identificar desvios das regras do projeto e propor planos detalhados de refatoração/implementação. Toda a interação deve ser feita em português.

---

## 🎯 Problema / Solicitação

O sistema de gestão de usuários, implementado pelo serviço service/user_service, pelo repositorio repository/user_repositoy, e pela persistencia representados pelas tabelas users e user_roles. Cada usuário terá necessariamente ao menos 1 role e alguns podem ter mais que um role. Caso tenha mais de um role associado, um deles deve ser o role "ativo", que indica o papel atual do usuário no sistema.

O sistema de permissionamento, implementado pelo serviço de service/permission_service, pelo repositorio permission/repository, e pela persistencia representada pelas tabelas roles, roles_permission e permissions. Cada role possui um conjunto de permissions associadas, que definem as ações que o usuário com aquele role pode executar no sistema.

Assim, ao chamar algum endpoint protegido, o sistema, atraves do permission_middleware, verifica se o user_role daquele usuário possui as permissions necessárias para executar a ação, com base no seu role ativo e nas permissions associadas a esse role.

O sistema de permissionamento gerencia as tabelas de roles, permissions e roles_permissions, enquanto o sistema de gestão de usuários gerencia as tabelas de users e user_roles. A associação entre usuários e seus roles é feita na tabela user_roles, onde um usuário pode ter múltiplos roles, mas apenas um deles é marcado como ativo.

Ocorre que em algum momento da construção do código, foi delegado a permission_repository a gestão de user_roles, o que gera complexidade para user_service construir um usuário inteiro com suas roles, sendo obrigado a chamar permisson_repository para obter as roles do usuário.

Tarefas:
1. Analise os codigos de user_service, user_repository, permission_service e permission_repository. Mapeando se a situação descrita procede.
2. No caso de ser procedente, isto viola alguma regra do guia de arquitetura do projeto ou de boas práticas de código? Justifique citando as seções específicas do guia.
3. Proponha um plano detalhado para corrigir o problema, realocando a responsabilidade de gestão de user_roles para user_repository, incluindo code skeletons para os arquivos que precisariam ser criados ou alterados, seguindo o formato descrito abaixo.
4. Apresente a estrutura final de diretórios e arquivos após a implementação do plano, seguindo a Regra de Espelhamento Port ↔ Adapter do guia.
5. Forneça uma ordem de execução numerada para implementar o plano, considerando dependências entre etapas.
6. Inclua um checklist de conformidade para garantir que o plano atende todas as regras do guia de arquitetura e padrões de código relevantes.

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