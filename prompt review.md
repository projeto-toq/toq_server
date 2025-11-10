### Engenheiro de Software Go Sênior — Análise e Refatoração TOQ Server

**Objetivo:** Atuar como engenheiro Go sênior para analisar código existente, identificar desvios das regras do projeto e propor planos detalhados de refatoração/implementação. Toda a interação deve ser feita em português.

---

## 🎯 Problema / Solicitação

O sistema de gestão de usuários é implementado pelo modelo model/user_model, pelo serviço service/user_service, pelo repositorio repository/user_repository, e pela persistencia representados pelas tabelas users e user_roles que podem ser consultadas pelo scripts/db_creation.sql.

Após inúmeras refatorações e adições de funcionalidades, fica a dúvida se as regras definidas no guia do projeto (docs/toq_server_go_guide.md) estão sendo seguidas corretamente.

Considerando as extesão da verificação, vamos focar em lotes de arquivos, iniciando pelos arquivos 
├── exists_email_for_another_user.go
├── exists_phone_for_another_user.go
├── get_active_user_role_by_user_id.go
├── get_agency_of_realtor.go
├── get_invite_by_phone_number.go
├── get_realtors_by_agency.go
├── get_user_by_id.go
├── get_user_by_nationalid.go
├── get_user_by_phone_number.go
├── get_user_role_by_user_id_and_role_id.go
├── get_user_roles_by_user_id.go
├── get_user_validations.go
├── get_users_by_role_and_status.go
├── get_wrong_signin_by_userid.go
├── has_user_duplicate.go

Tarefas, após ler o guia do projeto (docs/toq_server_go_guide.md):
1. Analise o código de cada um dos arquivos em busca de desvios das regras do guia.
2. Para cada desvio identificado, explique qual regra foi violada e o impacto disso no sistema.
3. Proponha um plano detalhado para corrigir os desvios, incluindo code skeletons para cada arquivo que precisa ser alterado ou criado.
    3.1. Caso a alteração seja apenas sobre a documentação, não é necessário apresentar o code skeleton.
4. Organize o plano em uma estrutura clara, incluindo a ordem de execução das tarefas e a estrutura de diretórios final.
5. Caso haja alguma sugestão de melhoria além da correção dos desvios, inclua no plano.


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