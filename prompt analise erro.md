### Resumo e Refatoração: Engenheiro de Software Go Sênior

Este documento descreve as instruções para atuar como um engenheiro de software Go sênior, focando na análise de um problema e na proposição de uma solução detalhada, seguindo a arquitetura hexagonal e boas práticas de código, garantindo o alinhamento com os padrões de arquitetura, tratamento de erros, observabilidade e documentação. Toda a interação deve ser feita em português.

---

**Problema:** Durante a chamada para o endpoint `DELETE /api/v2/user/account`, o sistema retorna um erro 409 com a mensagem "Active role missing for user", mesmo quando o usuário possui uma função ativa, pois é impossível o usuário não ter um role ativo. O log relevante é o seguinte:


{"time":"2025-09-22T07:40:04.285628577Z","level":"INFO","msg":"permission.check.allowed","user_id":2,"resource":"http","action":"GET:/api/v2/user/profile","permission_id":34}
{"time":"2025-09-22T07:40:04.288763323Z","level":"INFO","msg":"HTTP Request","request_id":"61efff10-d7ee-4d6c-a716-53bf2ac83886","method":"GET","path":"/api/v2/user/profile","status":200,"duration":6350962,"size":481,"client_ip":"93.36.219.10","user_agent":"Dart/3.9 (dart:io)","trace_id":"3dbd6a6b22efe91d1336389280307215","span_id":"6c063c8a940241ea","user_id":2,"user_role_id":2}
{"time":"2025-09-22T07:40:22.190753461Z","level":"INFO","msg":"permission.check.allowed","user_id":2,"resource":"http","action":"DELETE:/api/v2/user/account","permission_id":47}
{"time":"2025-09-22T07:40:22.20879252Z","level":"INFO","msg":"starting efficient user folder deletion in S3","userID":2,"bucket":"toq-app-media","prefix":"2/"}
{"time":"2025-09-22T07:40:22.389030856Z","level":"INFO","msg":"collected all objects for deletion","userID":2,"totalCount":6}
{"time":"2025-09-22T07:40:22.507343774Z","level":"INFO","msg":"user folder completely deleted from S3","userID":2,"bucket":"toq-app-media"}
{"time":"2025-09-22T07:40:22.509740594Z","level":"WARN","msg":"cannot issue access token without active role","user_id":2}
{"time":"2025-09-22T07:40:22.509850932Z","level":"ERROR","msg":"user.delete_account.delete_account_error","error":"HTTP 409: Active role missing for user","user_id":2}
{"time":"2025-09-22T07:40:22.519949833Z","level":"INFO","msg":"HTTP Response","request_id":"747bd88e-e016-4bf5-91e5-55bf1cc83090","method":"DELETE","path":"/api/v2/user/account","status":409,"duration":332388091,"size":68,"client_ip":"93.36.219.10","user_agent":"Dart/3.9 (dart:io)","trace_id":"01167dce55e510824992bbd4aadf4938","span_id":"410fbba027afcf0c","user_id":2,"user_role_id":2}


**Solicitação:** Analise o problema, **leia o código** envolvido, **ache a causa raiz** e proponha um plano detalhado para a implementação da solução. 

### **Instruções para a Proposição do Plano**

- **Ação:** Apenas a análise e a geração do plano são solicitadas. **Nenhum código deve ser implementado**.
- **Análise:** O problema e os requisitos devem ser analisados cuidadosamente. O código e arquivos de configuração existentes devem ser revisados para um plano preciso. Não faça suposições e confirme todos os detalhes necessários.
- **Plano:** Um plano detalhado deve ser apresentado, incluindo a descrição da arquitetura proposta, as interfaces, a estrutura de diretórios e a ordem de execução das etapas.
- **Qualidade do Plano:** O plano deve ser completo, sem o uso de _mocks_ ou soluções temporárias. Caso seja extenso, deve ser dividido em etapas implementáveis.
- **Acompanhamento:** As etapas já planejadas e as próximas a serem analisadas devem ser sempre informadas para acompanhamento.
- **Ambiente:** O plano deve considerar que estamos em ambiente de desvolvimento, portanto não deve haver back compatibility, migração de dados, preocupação com janela de manutenção ou _downtime_.

---

### **Regras Obrigatórias de Análise e Planejamento**

#### 1. Arquitetura e Fluxo de Código
- **Arquitetura:** A solução deve seguir estritamente a **Arquitetura Hexagonal**.
- **Fluxo de Chamadas:** As chamadas de função devem seguir a hierarquia `Handlers` → `Services` → `Repositories`.
- **Injeção de Dependência:** O padrão de _factories_ deve ser usado para a injeção de dependências.
- **Localização de Repositórios:** Os repositórios devem ser localizados em `/internal/adapter/right/mysql/` e deve fazer uso dos convertess para mapear entidades de banco de dados para entidades e vice versa.
- **Transações SQL:** Todas as transações de banco de dados devem utilizar `global_services/transactions`.


#### 2. Tratamento de Erros e Observabilidade

- **Tracing:**
  - Iniciar _tracing_ com `utils.GenerateTracer(ctx)` em métodos públicos de **Services**, **Repositories** e em **Workers/Go routines**.
  - Evitar _spans_ duplicados em **Handlers HTTP**, pois o `TelemetryMiddleware` já inicia o _tracing_.
  - Chamar a função de finalização (`defer spanEnd()`) e usar `utils.SetSpanError` para marcar erros.

- **Logging:**
  - Usar `slog` para _logs_ de domínio e segurança.
    - `slog.Info`: Eventos esperados do domínio.
    - `slog.Warn`: Condições anômalas ou falhas não fatais.
    - `slog.Error`: Falhas internas de infraestrutura.
  - Evitar _logs_ excessivos em **Repositórios (adapters)**.
  - **Handlers** não devem gerar _logs_ de acesso, pois o `StructuredLoggingMiddleware` já faz isso.

- **Tratamento de Erros:**
  - **Repositórios (Adapters):** Retornam erros "puros" (`error`).
  - **Serviços (Core):** Propagam erros de domínio usando `utils.WrapDomainErrorWithSource(derr)` e criam novos erros com `utils.NewHTTPErrorWithSource(...)`.
  - **Handlers (HTTP):** Usam `http_errors.SendHTTPErrorObj(c, err)` para converter erros em JSON.

#### 3. Boas Práticas Gerais
- **Estilo de Código:** A proposta deve seguir as **Go Best Practices** e o **Google Go Style Guide**.
- **Separação:** Manter a clara separação entre arquivos de **domínio**, **interfaces** e suas implementações.
- **Processo:** O plano não deve incluir a geração de _scripts_ de migração ou soluções temporárias.

---

### **Regras de Documentação e Comentários**

- A documentação da solução deve ser clara e concisa.
- A documentação das funções deve ser em **inglês**.
- Os comentários internos devem ser em **português**.
- A API deve ser documentada com **Swagger**, usando anotações diretamente no código.