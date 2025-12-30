### Engenheiro de Software Go Sênior — Análise e Implementação de funçoes TOQ Server

**Objetivo:** Atuar como engenheiro Go sênior para analisar código existente, entender claramente o que a regra de negócio exige e propor planos detalhados de refatoração/implementação da forma mais eficiente. Toda a interação deve ser feita em português.

---

## 🎯 Solicitação

Quando um realtor navega pelos listings publicados ele precisa ter a possibilidade de enviar um pedido de visita ao owner do imóvel. Atualmente essa funcionalidade não existe no TOQ Server e precisa ser implementada.
Em `/codigos/go_code/toq_server/docs/visits_system_specification.md`existe um rascunho de especificação do sistema de visitas que deve ser usado como referencia, e não fonte da verdade, para implementar essa funcionalidade. Os endpoints/payloads/respostas/enum devem ser seguidos sempre que for recomendado e interessante. Nenhum POST deve ter o id no path, sempre deve ser passado via body.
A regra de negócio preve:
1. O realtor envia um pedido de visita para o owner do imóvel.
   1.1. O modelo da visita `/codigos/go_code/toq_server/internal/core/model/listing_model/visit_domain.go` é um rascunho do que deve ser utilizado para representar o pedido de visita. Deve ser adequado conforme a necessidade.
   1.2. O pedido de visita deve ser baseado na agenda de disponibilidade que o owner criou durante a criação do listing representada em `/codigos/go_code/toq_server/internal/core/model/schedule_model/agenda_domain.go`.
   1.3. O pedido de visita dever ser enviado ao owner do imóvel via push notification (utilize o sistema de notificações já existente no TOQ Server).
2. O owner pode aceitar ou recusar o pedido de visita.
   2.1. Ao aceitar o pedido de visita, o sistema deve bloquear o horário na agenda do owner e na agenda do realtor para que não haja conflitos.
   2.2. Ao recusar o pedido de visita, o sistema deve enviar uma notificação ao realtor informando a recusa.
3. O realtor pode cancelar o pedido de visita a qualquer momento.
   3.1. Ao cancelar o pedido de visita, o sistema deve enviar uma notificação ao owner informando o cancelamento e retirar da agenda do owner e do realtor o bloqueio do horário.
4. Após a visita o realtor deve informar o status da visita (realizada, não realizada, reagendada).
   4.1. O owner deve ser notificado sobre o status da visita.
5. Deve haver um contador de tempo desde o envio do pedido de visitas até aceite/recusa do proprietário.
   5.1. Esta informação deve ser contabilizada pelo proprietário cobrindo todos os seus imoveis
   5.2. Esta informação deve ser armazenada para futuras análises de performance do owner e será mostrada em seus anuncios. EX: "Respondeu 90% dos pedidos de visita em até 2 horas".


Assim:
1. Analise o código atual model, service, handler, repository, dto, converter do projeto, leia o `toq_server_go_guide.md` e identifique a melhor forma de implementar a nova funcionalidade.
2. Proponha um plano detalhado de implementação incluindo:
   - Diagnóstico: arquivos envolvidos, justificativa da abordagem, impacto e melhorias possíveis.
   - O Codigo completo a ser implementado (handlers, services, repositories, DTOs, entities, converters), fazendo com a implementação seja simples e sem mais análises.
   - Estrutura de Diretórios: organização final seguindo a Regra de Espelhamento (Seção 2.1 do guia).
   - Ordem de Execução: etapas numeradas com dependências.
3. Siga todas as regras e padrões do projeto conforme documentado no guia do TOQ
4. Não se preocupe em garantir backend compatibilidade com versões anteriores, pois esta é uma alteração disruptiva.
5. Em `scripts/db_creation.sql` existe o modelo de dados atual do banco. Proponha as alterações necessárias para suportar a nova funcionalidade (sem scripts de migração).

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