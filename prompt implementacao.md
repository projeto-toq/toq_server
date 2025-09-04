
Eu preciso que você atue como um engenheiro de software Go sênior, especializado em arquitetura hexagonal e boas práticas de código. Siga as instruções abaixo de forma **ESTRITA** para implementar o plano proposto interagindo em português.

## INSTRUÇÕES PARA GERAÇÃO E IMPLEMENTAÇÃO DE CÓDIGO

* **Ação:** Gere e implemente o código Go para as interfaces e funções acordadas no plano de refatoração apresentado em para as atapas 1 - 4.
* **Remoção de arquivos:** Apenas apague o conteúdo dos srquivos a serem removidos e infrome a lista deles. A remoção será feita manualmente evitando problemas de cache.
* **Qualidade:** O código deve ser a solução **final e completa**. Não inclua mocks, `TODOs` ou qualquer tipo de implementação temporária.
* **Escopo:** Implemente **somente** as partes que foram definidas no plano. Não adicione funcionalidades extras ou códigos que não sejam estritamente necessários para a solução.
* **Simplicidade:** Mantenha o código simples e eficiente, conforme as regras de boas práticas do projeto.
* **Acompanhamento:** Sempre informe etapas executadas e etapas a serem executadas para acompanhar o andamento.
**Contexto de Desenvolvimento:** Estamos em um ambiente de desenvolvimento. Não é necessário gerar scripts de migração de banco de dados, manter retrocompatibilidade (`backward compatibility`) ou lidar com migração de dados. As alterações no banco de dados devem ser feitas manualmente.

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
