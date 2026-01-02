### Engenheiro de Software Go Sênior — Análise e Refatoração TOQ Server

**Objetivo:** Atuar como engenheiro Go sênior, para analisar código existente, identificar desvios das regras do projeto, implementações mal feitas ou mal arquitetadas, códigos errôneos e propor planos detalhados de refatoração/implementação. Toda a interação deve ser feita em português.

---

## 🎯 Problema / Solicitação

O sistema de gestão de pedidos de visitas do TOQ Server foi implementado de forma incompleta e com diversos desvios das regras e padrões do projeto.

A regra de negócio preve:
1. O realtor envia um pedido de visita para o owner do imóvel.
   1.1. O modelo da visita `/codigos/go_code/toq_server/internal/core/model/listing_model/visit_domain.go`.
   1.2. O pedido de visita deve ser baseado na agenda de disponibilidade do imovle que o owner criou durante a criação do listing representada em `/codigos/go_code/toq_server/internal/core/model/schedule_model/agenda_domain.go`, portanto visitas fora da disponibilidade não podem ser solicitadas.
   1.3. Um alerta do pedido de visita dever ser enviado ao owner do imóvel via push notification
   1.4. Utilize o sistema de notificações já existente no TOQ Server em `/codigos/go_code/toq_server/internal/core/service/global_service/notification_service.go`
2. O owner pode aceitar ou recusar o pedido de visita.
   2.1. Ao aceitar o pedido de visita, o sistema deve bloquear o horário na agenda do imovel e na agenda do realtor para que não haja conflitos.
   2.2. Ao aceitar o pedido de visita, o sistema deve enviar uma notificação ao realtor informando o aceite.
   2.3. Ao recusar o pedido de visita, o sistema deve enviar uma notificação ao realtor informando a recusa.
3. O realtor pode cancelar o pedido de visita a qualquer momento.
   3.1. Ao cancelar o pedido de visita, o sistema deve enviar uma notificação ao owner informando o cancelamento e retirar da agenda do owner e do realtor o bloqueio do horário.
4. Após a visita o realtor deve informar o status da visita (realizada, não realizada, reagendada).
   4.1. O owner deve ser notificado sobre o status da visita.
5. Deve haver um contador de tempo desde o envio do pedido de visitas até aceite/recusa do proprietário.
   5.1. Esta informação deve ser contabilizada pelo proprietário cobrindo todos os seus imoveis
   5.2. Esta informação deve ser armazenada para futuras análises de performance do owner e será mostrada em seus anuncios. EX: "Respondeu 90% dos pedidos de visita em até 2 horas".
6. Visitas podem ser solicitadas X horas a partir do pedido e no máximo Y dias no futuro.
   6.1. Estes valores X e Y devem ser configuráveis no env.yaml
   6.2. Caso o realtor tente solicitar uma visita fora destes limites, o sistema deve rejeitar a solicitação com a mensagem apropriada.

Portanto, o objetivo aqui é uma análise profunda e completa para identificara desvios/erros e propor um plano de refatoração detalhado.

Tarefas, após ler o guia do projeto em `docs/toq_server_go_guide.md`:
1. Analise o código dos handler, services, adapters, entities, converters e DTOs envolvidos no processamento das vistas.
2. Identifique todos os desvios e ausencias das regras de negócio, padrões e boas práticas descritas no guia do projeto (cite seções específicas) e na regra de negócio acima.
3. Proponha um plano detalhado para atender ao descritos nos manuais incluindo code skeletons para cada arquivo que precisa ser alterado ou criado.
    3.1. A refatoração pode ser disruptiva, pois este é um ambiente de dev e não temos back compatibility.
    3.2. se for necessário alterar o modelo da base de dados, apresente no novo modelo de dados que o DBA fará manualmente.
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