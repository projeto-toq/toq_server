### Engenheiro de Software Go Sênior — Análise e Refatoração TOQ Server

**Objetivo:** Atuar como engenheiro Go sênior para analisar código existente, entender claramente o que a regra de negócio exige e propor planos detalhados de refatoração/implementação da forma mais eficiente. Toda a interação deve ser feita em português.

---

## 🎯 Solicitação

A criação de listings no TOQ Server deve ser alterada para que caso o usuário deseje efetua alguma alteração no listing, seja criado uma nova versão do listing, ao invés de atualizar a versão existente.
Sobre esta nova versão do listing que é criada como draft, deve ser feita a validação através do endpoint de validação de listing, e caso esteja tudo correto, o usuário poderá tornar esta nova versão do listing, como a versão ativa do listing. Isto garante a preservação do histórico e do ciclo de vida do listing. Por exemplo, se o listing na V1 está no estado de 	StatusPendingPhotoScheduling, e o usuário deseja alterar alguma informação do listing, uma nova versão V2 será criada como draft. O usuário poderá então validar a V2, e caso esteja tudo correto, poderá promover a V2 para ser a versão ativa do listing que deverá manter o estado de 	StatusPendingPhotoScheduling. Assim, o histórico do listing permanece intacto, e o ciclo de vida é preservado.
Este processo precisa preservar as foreignkeys e relacionamentos existentes, como guarantias, features, exchange_places etc. entre versoes do mesmo listing.
Uma abordagem possível seria alterar o modelo de listing para, além do campo version que já existe, ter um campo uuid que identifique o grupo de versões do listing, e um campo active_version que identifica a versão activa dentro deste grupo. Assim, todas as versões do mesmo listing teriam o mesmo uuid, mas version_number diferentes (1, 2, 3, ...) mas só uma avtive_version. Isto permite inclusive retroceder a uma versão anterior. O endpoint de criação de listing então criaria um novo registro com o mesmo uuid e version_number incrementado gerenciando active_version.
As tabelas satelites que possuem foreign keys para listing precisariam referenciar o uuid e version_number para manter a integridade referencial e nÃo mais ter FK direta para o id do listing.


Assim:
1. Analise o código atual model, service, handler, repository, dto, converter relacionado ao listing e identifique a melhor forma de implementar a mudança.
   1.1) atenção especial as tabelas satelites de listing que possuem foreign keys para listing.
2. Proponha um plano detalhado de implementação, incluindo:
   - Diagnóstico: arquivos envolvidos, justificativa da abordagem, impacto e melhorias possíveis.
   - Code Skeletons: esqueletos para cada arquivo novo/alterado (handlers, services, repositories, DTOs, entities, converters) conforme templates da Seção 8 do guia.
   - Estrutura de Diretórios: organização final seguindo a Regra de Espelhamento (Seção 2.1 do guia).
   - Ordem de Execução: etapas numeradas com dependências.
   - Checklist de Conformidade: validação contra seções específicas do guia.
3. Siga todas as regras e padrões do projeto conforme documentado no guia do TOQ
4. Não se preocupe em garantir backend compatibilidade com versões anteriores, pois esta é uma alteração disruptiva e todos os listings serão apagados.
5. Não implemente alterações no script de DB, esta tarefa será feita manualmente pela equipe de DBA.
   5.1. o modelo de dados atual pode ser consultado em scripts/db_creation.sql;
   5.2. apresente as alteraçoes necessárias no modelo de dados para que a equipe de DBA possa implementar.

---

## 📘 Fonte da Verdade

**TODAS as regras de arquitetura, padrões de código, observabilidade e documentação estão em:**
- **`docs/toq_server_go_guide.md`** — Guia completo do projeto (seções 1-17)
- **`README.md`** — Configurações de ambiente e observabilidade

**⚠️ Consulte SEMPRE esses documentos antes de propor qualquer solução.**

---

## 🎯 Processo de Trabalho

1. **Leia o código** envolvido (adapters, services, handlers, entities, converters)
2. **Identifique a melhor forma de implementar** apresente evidencias no código
3. **Proponha plano detalhado** com code skeletons
4. **Não implemente código** — apenas análise e planejamento

---

## 📋 Formato do Plano

### 1. Diagnóstico
- Lista de arquivos analisados
- Porque esta é a melhor alternativa (apresente evidencias no código)
- Impacto da implementação
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