### SRE Sênior — Análise e Refatoração TOQ Server

**Objetivo:** Atuar como SRE sênior para analisar código existente, entender claramente o erro apresentado e identificar a causa raiz do problema para propor planos detalhados de refatoração. Toda a interação deve ser feita em português.

---

## 🎯 Problema / Solicitação

Existe um manual de observabilidade em `/codigos/go_code/toq_server/docs/observability/sre_guide.md` que foi criado na primeira implementaçÃo, mas nÃo deve ser tomado como fonte da verdade. As configurações efetivamente implementadas são a fonte da verdade.
A telemetria do TOQ Server está sendo feita por:
`/codigos/go_code/toq_server/internal/core/config/telemetry.go`.
Grafana concentra a análise dos dados coletados de prometheus, tempo, loki. Todos rodando em Docker segundo `/codigos/go_code/toq_server/docker-compose.yml`.
Ocorre que os logs sÃo hidratados e geram o seguinte registro no dashboard TOQ Server - Logs do grafana:
```json
{"body":"notification.async_send_error","severity":"ERROR","attributes":{"code.file.path":"/codigos/go_code/toq_server/internal/core/service/global_service/notification_service.go","code.function.name":"github.com/projeto-toq/toq_server/internal/core/service/global_service.(*unifiedNotificationService).SendNotification.func1","code.line.number":95,"deployment.environment":"homo","err":"HTTP 500: Internal server error","service.name":"toq_server","service.namespace":"projeto-toq","service.version":"2.0.0","to":"","token":"euho9KY_5EPUm-EnTMPAe6:APA91bGJ6alJhbEutQ7Nz3DyVt2JE6Yw5KHc0TUlF6QZmwmSnSSMM2b1fzSmdq92zB0fPkgf4yB_VyVmLtaKVyp8wTrGgrVqGJCDhJkWcdpKAapns5HMMb0","type":"unhandled: (globalservice.NotificationType) fcm"},"resources":{"deployment.environment":"homo","host.name":"bbf1a8bbc4e9","os.type":"linux","service.instance.id":"ip-172-31-81-196-2231546","service.name":"toq_server","service.namespace":"projeto-toq","service.version":"2.0.0","telemetry.sdk.language":"go","telemetry.sdk.name":"beyla","telemetry.sdk.version":"1.38.0"},"instrumentation_scope":{"name":"toq_server","version":"2.0.0"}}
```
Existe muita informaçÃo irrelevante neste log que dificulta a análise do problema, a mensgem de erro em si é `HTTP 500: Internal server error` o que nÃo ajuda a identificar a causa raiz do problema,
Adicionamelmente não existe trace correspondente a esta entrada no log o que impossibilita a corelação do erro com o fluxo de execução do código.
Por ser um servidor REST-API todo o fluxo de execução deveria ser rastreável via traces e logs correlacionados. Cada chamada http deve gerar um request-id único que deve ser propagado por todo o fluxo de execução do código, permitindo a correlação entre logs e traces. Entretanto existe um trace-id e um request-id diferente para cada log, o que indica que o trace-id e request-id não estão sendo propagados corretamente.

Assim:
1. Analise o guia do projeto `docs/toq_server_go_guide.md`, o código atual e identifique a causa raiz do problema
2. Proponha um plano detalhado de refatoração com code skeletons para corrigir o problema, seguindo estritamente as regras de arquitetura do manual `docs/toq_server_go_guide.md` (observabilidade, erros, transações, etc).
3. Leia atentamente as cofiguraçòes atuais dos containers para evitar quebraas. As última refatorações foram traumaticas por `assumir` configurações que não existiam e quebrar o ambiente.
4. Ao final do plano deve haver uma atualização de `/codigos/go_code/toq_server/docs/observability/sre_guide.md`, readme.md e guia do projeto para refletir as mudanças propostas.
---

**TODAS as regras de arquitetura, padrões de código, observabilidade e documentação estão em:**
- **`docs/toq_server_go_guide.md`** — Guia completo do projeto (seções 1-17)
- **`README.md`** — Configurações de ambiente e observabilidade

**⚠️ Consulte SEMPRE esses documentos antes de propor qualquer solução.**

---

## 🎯 Processo de Trabalho

1. **Leia o código e as configurações reais de containers** envolvido (adapters, services, handlers, entities, converters)
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