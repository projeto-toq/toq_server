Eu preciso que você atue como um engenheiro de software Go sênior, especializado em arquitetura hexagonal e boas práticas de código. Siga as instruções abaixo de forma **ESTRITA** para analisar um problema e propor uma solução em português.

---
🛠️ Problema
verifique os erros abaixo e analise a causa raiz do problema: o solicitar o reenvio do codigo de validação de troca de e-mail, após sign in bem sucedido.
{"time":"2025-09-03T18:43:05.132621868Z","level":"INFO","msg":"Security event logged","eventType":"signin_success","result":"success","timestamp":"2025-09-03T18:43:05.132609288Z","userID":2,"nationalID":"04679654805","ipAddress":"179.110.194.42","userAgent":"PostmanRuntime/7.45.0"}
{"time":"2025-09-03T18:43:05.132708949Z","level":"INFO","msg":"User signed in successfully","userID":2}
{"time":"2025-09-03T18:43:05.140500699Z","level":"INFO","msg":"HTTP Request","request_id":"fdbb7f08-9586-4415-968c-47b16f5e61d5","method":"POST","path":"/api/v2/auth/signin","status":200,"duration":153802609,"size":599,"client_ip":"179.110.194.42","user_agent":"PostmanRuntime/7.45.0"}
{"time":"2025-09-03T18:43:43.212901377Z","level":"WARN","source":{"function":"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/middlewares.StructuredLoggingMiddleware.func1","file":"/codigos/go_code/toq_server/internal/adapter/left/http/middlewares/structured_logging_middleware.go","line":126},"msg":"HTTP Error","request_id":"318167be-f9cb-4eba-84a8-5db39a3d993b","method":"POST","path":"/api/v2/user/email/resend","status":401,"duration":17870991,"size":47,"client_ip":"179.110.194.42","user_agent":"PostmanRuntime/7.45.0","user_id":2,"user_role_id":2,"role_status":"pending_both"}
---
### REGRAS OBRIGATÓRIAS DE ANÁLISE E PLANEJAMENTO

1.  **Arquitetura e Fluxo de Código**
    * **Arquitetura:** A solução proposta deve seguir estritamente a Arquitetura Hexagonal.
    * **Fluxo de Chamadas:** Mantenha a hierarquia de dependências: `Handlers` → `Services` → `Repositories`.
    * **Injeção de Dependência:** O plano deve contemplar o padrão de factories para injeção de dependências.
    * **Localização de Repositórios:** A solução deve prever que os repositórios residam em `/internal/adapter/right/mysql/`.
    * **Transações SQL:** Todas as transações de banco de dados devem utilizar `global_services/transactions`.

2.  **Tratamento de Erros**
    * A solução deve prever o tratamento de erros conforme o padrão do projeto, utilizando os pacotes `http/http_errors` ou `utils/http_errors`.
    * O plano deve garantir a correta propagação e log de erros.

3.  **Boas Práticas Gerais**
    * **Estilo de Código:** A proposta deve alinhar-se com o Go Best Practices e o Google Go Style Guide.
    * **Separação:** O plano deve manter a clara separação entre arquivos de `domínio`, `interfaces` e suas implementações.
    * **Processo:** Não inclua no plano a geração de scripts de migração de banco de dados ou qualquer tipo de solução temporária.

---

### REGRAS DE DOCUMENTAÇÃO E COMENTÁRIOS
* A documentação da solução deve ser clara e concisa.
* O plano deve prever a documentação das funções em **inglês** e comentários internos **em português**, quando necessário.
* Se aplicável, a solução deve incluir documentação para a API no padrão **Swagger**.

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