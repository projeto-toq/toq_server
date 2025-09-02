Eu preciso que você atue como um engenheiro de software Go sênior, especializado em arquitetura hexagonal e boas práticas de código. Siga as instruções abaixo de forma **ESTRITA** para analisar um problema e propor uma solução.

---
🛠️ Problema
O usuário está recebendo permission denied ao fazer login, o que não tem sentido.
{"time":"2025-09-02T13:08:05.731778183Z","level":"WARN","msg":"Permission denied","userID":4,"method":"POST","path":"/api/v1/user/signout"}
{"time":"2025-09-02T13:08:05.731914834Z","level":"WARN","source":{"function":"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/middlewares.StructuredLoggingMiddleware.func1","file":"/codigos/go_code/toq_server/internal/adapter/left/http/middlewares/structured_logging_middleware.go","line":126},"msg":"HTTP Error","request_id":"ec9e4982-aa06-4c89-9399-586b317a272a","method":"POST","path":"/api/v1/user/signout","status":403,"duration":3357187,"size":49,"client_ip":"179.110.194.42","user_agent":"PostmanRuntime/7.45.0","user_id":4,"user_role_id":4,"role_status":"pending_both"}
verifique as permissões carregadas nos CSVs de /data e infrome o que é necessário para incluir a permissão de signout a todos os usuários

---
**REGRAS OBRIGATÓRIAS DE DESENVOLVIMENTO EM GO**
1.  **Arquitetura e Fluxo de Código**
    * **Arquitetura:** Implemente estritamente a Arquitetura Hexagonal.
    * **Fluxo de Chamadas:** Mantenha a hierarquia de dependências: `Handlers` → `Services` → `Repositories`.
    * **Injeção de Dependência:** Use o padrão de factories (`/config/*`, `/factory/*`) para injetar dependências. Inicialize `adapters` e `services` **uma única vez** no início da aplicação.
    * **Localização de Repositórios:** Os repositórios devem residir em `/internal/adapter/right/mysql/`.
    * **Transações SQL:** Use exclusivamente `global_services/transactions` para todas as transações de banco de dados.

2.  **Tratamento de Erros**
    * **Padrão:** Erros devem ser tratados com o pacote `http/http_errors` (para `adapter errors`) ou `utils/http_errors` (para `DomainError`).
    * **Propagação:** Logue e transforme o erro **apenas no ponto de origem**. Funções intermediárias devem apenas repassar o erro sem logar ou recriar.
    * **Verificação:** Sempre verifique o retorno de erro de qualquer função.

3.  **Boas Práticas Gerais**
    * **Estilo de Código:** Siga o Go Best Practices e o Google Go Style Guide. Mantenha o código simples, eficiente e consistente.
    * **Separação:** Mantenha a clara separação entre arquivos de `domínio`, `interfaces` e suas implementações.
    * **Processo:** Não use mocks ou código temporário. O código legado deve ser completamente removido. Não gere scripts de migração de DB; alterações devem ser manuais via MySQL Workbench.

---
**INSTRUÇÕES FINAIS**
* **Ação:** Não implemente nenhum código.
* **Análise:** Analise cuidadosamente o problema (`log.md`) e os requisitos. Se necessário, solicite informações adicionais. Analise sempre o código existente.
* **Plano:** Apresente um plano detalhado para a refatoração. O plano deve incluir:
    * Descrição da arquitetura proposta e seu alinhamento com a arquitetura hexagonal.
    * Interfaces a serem criadas (com métodos e assinaturas).
    * Estrutura de diretórios e arquivos sugerida.
    * Ordem das etapas de refatoração para garantir uma transição suave.
* **Qualidade do Plano:** O plano deve ser completo, sem mocks ou soluções temporárias. Se for muito grande, divida-o em etapas que possam ser implementadas separadamente.