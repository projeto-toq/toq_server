### Resumo e Refatoração: Engenheiro de Software Go Sênior

Este documento descreve as instruções para atuar como um engenheiro de software Go sênior, focando na análise de um problema e na proposição de uma solução detalhada, seguindo a arquitetura hexagonal e boas práticas de código, garantindo o alinhamento com os padrões de arquitetura, tratamento de erros, observabilidade e documentação. Toda a interação deve ser feita em português.

---

**Problemas:**
Diversos endpoint que aceitam datas no formato RFC3339 estão falhando em receber datas como o endpoint /schedules/listing/availability?listingId=3&rangeFrom=2025-10-31T08:00:00+03:00&rangeTo=2025-11-02T19:00:00+03:00&slotDurationMinute=60&page=1&limit=50&timezone=America/Sao_Paulo

A linha 56-60 do get_listing_availability.go está com erro de parsing de data.

Além disso, como está sendo passado o timezone na query string, está ficando confuso qual timezone deve ser considerado para o rangeFrom e rangeTo.

Assim:
- faça uma lista de todos os endpoints que aceitam datas no formato RFC3339 e apresente-a.
- faça uma lista de todos os endpoints que aceitam apenas datas no formato YYYY-MM-DD e apresente-a.
- proponha um plano para que:
  - altere todos os os endpoints que aceitam datas no formato RFC3339 para remover timezone na query string ou no body, e que o timezone seja extraído diretamente da data enviada no formato RFC3339.
    - esta alteração deve ser propagada, de forma que o service receba o timezone no propio location do time enviado na data.
    - internamente o time é tratado como UTC, então garanta que essa conversão seja feita corretamente.
  - altere todos o endpoint GET /listings/photo-session/slots, que aceitam datas no formato YYYY-MM-DD, para que necessariamente passe timezone na query string ou no body, para evitar ambiguidades.
    - esta alteração deve ser propagada, de forma que o service receba o timezone no propio location do time enviado na data.
    - internamente o time é tratado como UTC, então garanta que essa conversão seja feita corretamente.
    - os demais endpoints que aceitam datas no formato YYYY-MM-DD devem continuar como estão, sem necessidade de alteração, visto que as datas informadas são datas de nascimento, que não possuem timezone. confirme isso no código.
  - garanta que o parsing de datas em RFC3339 seja corrigido em todos os endpoints que aceitam esse formato, pois hoje existe um erro de parsing.


**Solicitação:** Analise o problema, **leia o código** envolvido, **ache a causa raiz** e proponha um plano detalhado para a implementação/refatoração da solução, após ler o o manual do projeto em docs/toq_server_go_guide.md.

### **Instruções para a Proposição do Plano**

- **Ação:** Apenas a análise e a geração do plano são solicitadas. **Nenhum código deve ser implementado**.
- **Análise:** O problema e os requisitos devem ser analisados cuidadosamente. O código e arquivos de configuração existentes devem ser revisados para um plano preciso. Não faça suposições e confirme todos os detalhes necessários.
- **Plano:** Um plano detalhado deve ser apresentado, incluindo a descrição da arquitetura proposta, as interfaces, a estrutura de diretórios e a ordem de execução das etapas.
- **Qualidade do Plano:** O plano deve ser **extremamente prescritivo**, sem o uso de _mocks_ ou soluções temporárias. Para cada arquivo novo ou alterado, inclua um **esqueleto de código (code skeleton)** ou a **assinatura completa da função (full function signature)**, mostrando:
    * Structs, Interfaces (Ports) e DTOs com todos os campos.
    * Assinaturas completas de métodos públicos e privados.
    * Uso explícito de `utils.GenerateTracer`, `defer spanEnd()`, e `utils.SetSpanError` nos pontos aplicáveis (Services e Repositories).
- **Acompanhamento:** As etapas já planejadas e as próximas a serem analisadas devem ser sempre informadas para acompanhamento.
- **Ambiente:** O plano deve considerar que estamos em ambiente de desvolvimento, portanto não deve haver back compatibility, migração de dados, preocupação com janela de manutenção ou _downtime_.
- **Testes:** O plano **NÃO** deve incluir a criação/alteração de testes unitários e de integração para garantir a qualidade do código.
- **Documentação:** A documentação Swagger/docs deve ser criada por comentários em DTO/Handler e execução de make swagger. Sem alterações manuais no swagger.yaml/json.
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
- Não execute git status, git diff nem go test.

---

### **Regras de Documentação e Comentários**

- A documentação da solução deve ser clara e concisa.
- A documentação das funções deve ser em **inglês**.
- Os comentários internos devem ser em **português**.
- A API deve ser documentada com **Swagger**, usando anotações diretamente no código, em inglês e não alterando swagger.yaml/json manualmente.