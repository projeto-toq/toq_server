Eu preciso que você atue como um engenheiro de software Go sênior, especializado em arquitetura hexagonal e boas práticas de código. Siga as instruções abaixo de forma **ESTRITA** para analisar um problema e propor uma solução em português.

### 🛠️ Análise e Solução

**Problema:** Analise o extrato do log abaixo e veja que está verboso e inutil, pois source fala de 

- /codigos/go_code/toq_server/internal/adapter/left/http/middlewares/structured_logging_middleware.go;
- function":"github.com/giulioalfieri/toq_server/internal/core/utils.WrapDomainErrorWithSource","file":"/codigos/go_code/toq_server/internal/core/utils/http_errors.go","line":292,"stack":["github.com/giulio-alfieri/toq_server/internal/core/utils.WrapDomainErrorWithSource (http_errors.go:292)"

que são as funções de logging e não funções relacionadas ao erro em si.

{"time":"2025-09-04T17:24:47.056887306Z","level":"WARN","source":{"function":"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/middlewares.StructuredLoggingMiddleware.func1","file":"/codigos/go_code/toq_server/internal/adapter/left/http/middlewares/structured_logging_middleware.go","line":163},"msg":"HTTP Error","request_id":"e7fc8c37-80f0-40e3-8087-b4f8a539a475","method":"POST","path":"/api/v2/user/phone/confirm","status":401,"duration":427709,"size":45,"client_ip":"179.110.194.42","user_agent":"PostmanRuntime/7.45.0","trace_id":"81d165a5783a27fb6fa2982129524ebb","span_id":"297d9dee22eee1c2","function":"github.com/giulio-alfieri/toq_server/internal/core/utils.WrapDomainErrorWithSource","file":"/codigos/go_code/toq_server/internal/core/utils/http_errors.go","line":292,"stack":["github.com/giulio-alfieri/toq_server/internal/core/utils.WrapDomainErrorWithSource (http_errors.go:292)"],"error_code":401,"error_message":"Invalid access token","errors":["HTTP 401: Invalid access token"]}

Analise o código e o log em detalhes e apresenta acausa raiz.

Apresente um plano de refatoração para correção.

---

### REGRAS OBRIGATÓRIAS DE ANÁLISE E PLANEJAMENTO

1.  **Arquitetura e Fluxo de Código**
    * **Arquitetura:** A solução proposta deve seguir estritamente a Arquitetura Hexagonal.
    * **Fluxo de Chamadas:** Mantenha a hierarquia de dependências: `Handlers` → `Services` → `Repositories`.
    * **Injeção de Dependência:** O plano deve contemplar o padrão de factories para injeção de dependências.
    * **Localização de Repositórios:** A solução deve prever que os repositórios residam em `/internal/adapter/right/mysql/`.
    * **Transações SQL:** Todas as transações de banco de dados devem utilizar `global_services/transactions`.

2.  **Tratamento de Erros e Observabilidade**
    * **Tracing:** A solução deve iniciar o tracing para cada operação com `utils.GenerateTracer(ctx)`.
    * **Logging:**
        * **Logs de Domínio e Segurança:** Utilize o pacote `slog`.
            * `slog.Info`: Para eventos de domínio esperados (ex: status do usuário mudou de pendente para ativo).
            * `slog.Warn`: Para condições anômalas, como indícios de fraude/reuso ou falhas não fatais.
            * `slog.Error`: Exclusivamente para falhas internas de infraestrutura, como problemas de transação com o banco de dados.
        * **Logs em Repositórios:** Evite logs excessivos. Em caso de falha crítica de infraestrutura (ex: erro de conexão com DB), use `slog.Error` com contexto mínimo (ex: `user_id` ou `key_query`).
    * **Tratamento de Erros:**
        * **Repositórios (Adapters):** Retorne erros "puros" (`error`) ou erros de domínio. **Nunca** use pacotes HTTP (`http` ou `http_errors`) nesta camada.
        * **Serviços (Core):** Propague erros de domínio utilizando `utils.WrapDomainErrorWithSource(derr)` para preservar a origem (função/arquivo/linha). Se for um erro novo, use `utils.NewHTTPErrorWithSource(...)` para criá-lo. Não serializar respostas HTTP diretamente aqui.
        * **Handlers (HTTP):**
            * Use `http_errors.SendHTTPErrorObj(c, err)` para converter qualquer erro propagado em uma resposta JSON com o formato `{code, message, details}`. Este helper também anexará o erro no contexto (`c.Error`) para que o middleware de log possa capturar a origem e os detalhes.
            * Evite construir payloads de erro manualmente.

3.  **Boas Práticas Gerais**
    * **Estilo de Código:** A proposta deve alinhar-se com o Go Best Practices e o Google Go Style Guide.
    * **Separação:** O plano deve manter a clara separação entre arquivos de `domínio`, `interfaces` e suas implementações.
    * **Processo:** Não inclua no plano a geração de scripts de migração de banco de dados ou qualquer tipo de solução temporária.

---

### REGRAS DE DOCUMENTAÇÃO E COMENTÁRIOS
* A documentação da solução deve ser clara e concisa.
* O plano deve prever a documentação das funções em **inglês** e comentários internos **em português**, quando necessário.
* Se aplicável, a solução deve incluir documentação para a API no padrão **Swagger**, feitas no código e não no swagger.yaml/json diretamente.

---

### INSTRUÇÕES FINAIS PARA O PLANO
* **Ação:** Não implemente nenhum código. Apenas analise e gere o plano.
* **Análise:** Analise cuidadosamente o problema e os requisitos. Se necessário, solicite informações adicionais. Analise sempre o código e os arquivos de configuração existentes.
* **Plano:** Apresente um plano detalhado para a implementação. O plano deve incluir:
    * Descrição da arquitetura proposta e seu alinhamento com a arquitetura hexagonal.
    * Interfaces a serem criadas (com métodos e assinaturas).
    * Estrutura de diretórios e arquivos sugerida.
    * Ordem das etapas de refatoração para garantir uma transição suave.
* **Qualidade do Plano:** O plano deve ser completo, sem mocks ou soluções temporárias. Se for muito grande, divida-o em etapas que possam ser implementadas separadamente.
* **Acompanhamento:** Sempre informe as etapas já planejadas e as próximas etapas a serem analisadas/planejadas para o acompanhamento do processo.