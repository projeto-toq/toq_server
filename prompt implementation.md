### Engenheiro de Software Go Sênior — Análise e Refatoração TOQ Server

**Objetivo:** Atuar como engenheiro Go sênior para analisar código existente, entender claramente o que a regra de negócio exige e propor planos detalhados de refatoração/implementação da forma mais eficiente. Toda a interação deve ser feita em português.

---

## 🎯 Solicitação

o processo de criação e atualização de versões de listings deve ser modificado para:
1) o endpoint de criação de listing POST /listings deve ser usado APENAS para criar a versão inicial (versão 1) em status DRAFT
   1.1. deve haver validação se existe uma versão ativa (não-expirada/não-fechada) para o listingIdentityId; ou se existe um listing par ao endereço selecionado, se existir, retornar erro 409 
2) sobre esta versÃo inciail, o endpoint PUT /listings faz todas as atualizações, sempre verificando que a versão está em status DRAFT
3) ao terminar as atualizações o endpoint POST /listings/versions/promote deve ser chamado para promover a versão DRAFT para:
   3.1 - Se for a primeira versão (v1), muda o status para `StatusPendingAvailability` e cria a agenda básica do imóvel
	3.2 - Se for uma versão posterior, mantém o status da versão ativa anterior (preserva o ciclo de vida do listing)
4) para criar uma nova versão DRAFT a partir de uma versão ativa existente, deve ser usado o novo endpoint POST /listings/versions/draft
   4.1 - este endpoint deve validar se a versão ativa está em um dos status permitidos para cópia (ver regras abaixo)
   4.2 - se já existir uma versão DRAFT não-promovida, retornar o versionId desta versão
   4.3 - caso contrário, criar uma nova versão DRAFT, copiando todos os dados da versão ativa (incluindo entidades satélite: features, exchange_places, financing_blockers, guarantees, etc)
   4.4 - retornar o versionId e status da nova versão DRAFT criada

### Regras de Cópia de Versão Ativa para DRAFT
- permitir cópia APENAS de: StatusSuspended, StatusRejectedByOwner, StatusPendingPhotoProcessing, StatusPhotosScheduled, StatusPendingPhotoConfirmation, StatusPendingPhotoScheduling, StatusPendingAvailability;
- bloquear StatusPublished com mensagem "Listing is published. Suspend it via status update before creating a draft version";
- bloquear StatusUnderNegotiation/StatusPendingAdminReview/StatusPendingOwnerApproval com "Listing is locked in workflow and cannot be copied";
- bloquear StatusExpired/StatusArchived/StatusClosed com "Listing is permanently closed and cannot be edited"


Assim:
1. Analise o código atual model, service, handler, repository, dto, converter relacionado ao listing e identifique a melhor forma de implementar a mudança.
2. Proponha um plano detalhado de implementação, incluindo:
   - Diagnóstico: arquivos envolvidos, justificativa da abordagem, impacto e melhorias possíveis.
   - Code Skeletons: esqueletos para cada arquivo novo/alterado (handlers, services, repositories, DTOs, entities, converters) conforme templates da Seção 8 do guia.
   - Estrutura de Diretórios: organização final seguindo a Regra de Espelhamento (Seção 2.1 do guia).
   - Ordem de Execução: etapas numeradas com dependências.
   - Checklist de Conformidade: validação contra seções específicas do guia.
3. Siga todas as regras e padrões do projeto conforme documentado no guia do TOQ
4. Não se preocupe em garantir backend compatibilidade com versões anteriores, pois esta é uma alteração disruptiva e todos os listings serão apagados.
5. Verifique se os endpoints podem ter uma nomenclatura melhor, mas mantenha os verbos HTTP conforme descrito.

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