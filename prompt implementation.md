### Engenheiro de Software Go Sênior — Análise e Implementação de funçoes TOQ Server

**Objetivo:** Atuar como engenheiro Go sênior para analisar código existente, entender claramente o que a regra de negócio exige e propor planos detalhados de refatoração/implementação da forma mais eficiente. Toda a interação deve ser feita em português.

---

## 🎯 Solicitação

Quando um realtor navega pelos listings publicados ele precisa ter a possibilidade de enviar uma proposta ao owner do imóvel. Atualmente essa funcionalidade não existe no TOQ Server e precisa ser implementada.

A regra de negócio preve:
### Para CORRETOR:
1. Enviar proposta para um imóvel
2. Visualizar histórico de propostas enviadas
3. Visualizar status de cada proposta
4. Editar proposta (apenas se status = `pending`)
5. Cancelar proposta (antes da aceitação)
   
### Para PROPRIETÁRIO:
1. Visualizar propostas recebidas
2. Aceitar proposta
3. Recusar proposta (com motivo)
4. Visualizar histórico de propostas recebidas e seus status

1. O realtor envia uma proposta para o owner do imóvel.
   1.1. A proposta pode ser enviada por um texto livre ou por um pdf(com tamanho máximo de 1MB). ambos devem ser armazenados na base de dados.
   1.2. Deve ser enviado um push notification (utilize o sistema de notificações já existente no TOQ Server) quando uma proposta for enviada ao owner do imóvel.
2. O owner pode aceitar ou recusar a proposta.
   2.1. Ao aceitar ou recusar a proposta um push notificatioin deve ser enviado ao realtor informando o status da proposta.
   2.2. Ao recusar a proposta o owner deve informar um motivo (texto livre).
3. O realtor pode cancelar a proposta a qualquer momento antes do aceite pelo owner.
   3.1. Ao cancelar a proposta, o sistema deve enviar uma notificação ao owner informando o cancelamento.
4. Ambos realtor e owner podem visualizar o histórico de propostas enviadas/recebidas com seus respectivos status (pending, accepted, refused, cancelled).
5. o listing deve ter um campo que indique se existe propsota aceita ou pendente.

O plano em `/codigos/go_code/toq_server/docs/proposals_implementation_plan.md` foi criado para implementar este funcionalidade, mas não foi finalizado e não atende a totalidade dos requisitos. Sua tarefa é analisar o plano existente, o código do TOQ Server e propor ajustes para tornar este plano, um plano completo de implementação seguindo todas as regras e padrões do projeto.

Assim:
1. Analise o `toq_server_go_guide.md` e identifique a melhor forma de implementar a nova funcionalidade.
2. Proponha um plano detalhado de implementação incluindo:
   - Diagnóstico: arquivos envolvidos, justificativa da abordagem, impacto e melhorias possíveis.
   - O Codigo completo a ser implementado (handlers, services, repositories, DTOs, entities, converters), fazendo com a implementação seja simples e sem mais análises.
   - Estrutura de Diretórios: organização final seguindo a Regra de Espelhamento (Seção 2.1 do guia).
   - Ordem de Execução: etapas numeradas com dependências.
3. Siga todas as regras e padrões do projeto conforme documentado no guia do TOQ
4. Não se preocupe em garantir backend compatibilidade com versões anteriores, pois esta é uma alteração disruptiva.
5. Em `scripts/db_creation.sql` existe o modelo de dados atual do banco. Proponha as alterações necessárias para suportar a nova funcionalidade (sem scripts de migração), que será implemtnentada posteriormente por outro time.

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