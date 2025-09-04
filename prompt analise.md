Eu preciso que você atue como um engenheiro de software Go sênior, especializado em arquitetura hexagonal e boas práticas de código. Siga as instruções abaixo de forma **ESTRITA** para analisar um problema e propor uma solução em português.

### 🛠️ Análise e Solução

**Problema:** O processo de logging da aplicação está confuso e inconsistente. Preciso de uma anlaise sua como especialista em go e metodologias ágeis para propor um plano de ação detalhado para padronizar o logging e tratamento de erros em toda a aplicação, alinhado com as melhores práticas do mercado.

Considere que:
- é um Rest API e portanto o handler sempre tem que retornar HHTP status code e payload json.
- a arquitetura é hexagonal, com handlers, services e repositories.
- o logging deve ser estruturado e consistente, com níveis de severidade claros (info, warn, error).
- o tratamento de erros deve ser padronizado, com erros de domínio e erros de infraestrutura claramente diferenciados.
- o plano deve incluir exemplos de código para cada camada (handler, service, repository).
- o plano deve prever a implementação de middlewares para logging e tratamento de erros.
- o plano deve garantir que o código siga as melhores práticas de Go, incluindo estilo, organização e documentação.
- o plano deve prever a documentação das mudanças, incluindo atualizações nos handlers/DTO permitindo gerar a doc swagger.
- o log criado tem que indicar claramente o local do erro e não o wrapper/util do log
- o plano deve prever a correlação entre logs e traces, utilizando trace_id e span_id quando disponíveis.

---

### REGRAS OBRIGATÓRIAS DE ANÁLISE E PLANEJAMENTO

1.  **Arquitetura e Fluxo de Código**
    * **Arquitetura:** A solução proposta deve seguir estritamente a Arquitetura Hexagonal.
    * **Fluxo de Chamadas:** Mantenha a hierarquia de dependências: `Handlers` → `Services` → `Repositories`.
    * **Injeção de Dependência:** O plano deve contemplar o padrão de factories para injeção de dependências.
    * **Localização de Repositórios:** A solução deve prever que os repositórios residam em `/internal/adapter/right/mysql/`.
    * **Transações SQL:** Todas as transações de banco de dados devem utilizar `global_services/transactions`.

2.  Tratamento de Erros e Observabilidade

- Tracing
  - Inicie o tracing por operação com `utils.GenerateTracer(ctx)` no início de cada método público de Services e em Workers/Go routines.
  - Em Handlers HTTP, o tracing já é iniciado pelo `TelemetryMiddleware`. Não crie spans duplicados via `GenerateTracer` no handler.
  - Sempre chame a função de finalização retornada por `GenerateTracer` (ex.: `defer spanEnd()`). Erros devem marcar o span via `utils.SetSpanError` — nos handlers isso já é feito por `SendHTTPErrorObj` e no caso de panics pelo `ErrorRecoveryMiddleware`.

- Logging
  - Logs de domínio e segurança: use apenas `slog`.
    - `slog.Info`: eventos esperados do domínio (ex.: user status mudou de pending para active).
    - `slog.Warn`: condições anômalas, indícios de fraude/reuso, limites atingidos, falhas não fatais (ex.: 429/423 por throttling/lock).
    - `slog.Error`: exclusivamente para falhas internas de infraestrutura (DB, transação, providers externos). Devem ser registrados no ponto de ocorrência.
  - Repositórios (adapters): evite logs excessivos. Em falhas críticas de infraestrutura, logue com `slog.Error` incluindo somente contexto mínimo e útil (ex.: `user_id`, `key_query`). Sucessos devem ser no máximo `DEBUG` quando realmente necessário.
  - Handlers não devem gerar logs de acesso; o `StructuredLoggingMiddleware` já o faz centralmente com severidade baseada no status HTTP (5xx→ERROR, 429/423→WARN, demais 4xx→INFO, 2xx/3xx→INFO).

- Tratamento de Erros
  - Repositórios (Adapters): retornam erros "puros" (`error`). Nunca usar pacotes HTTP (`net/http` ou `http_errors`) nesta camada.
  - Serviços (Core): propagar erros de domínio usando `utils.WrapDomainErrorWithSource(derr)` para preservar a origem (função/arquivo/linha). Ao criar novos erros de domínio, usar `utils.NewHTTPErrorWithSource(...)`. Mapear erros de repositório para erros de domínio quando aplicável. Não serializar respostas HTTP aqui.
  - Handlers (HTTP): usar `http_errors.SendHTTPErrorObj(c, err)` para converter qualquer erro propagado em JSON `{code, message, details}`. O helper também executa `c.Error(err)` para que o middleware de log capte a origem/detalhes e marca o span no trace.

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