### Resumo e Refatoração: Engenheiro de Software Go Sênior

Este documento descreve as instruções para atuar como um engenheiro de software Go sênior, focando na análise de um problema e na proposição de uma solução detalhada, seguindo a arquitetura hexagonal e boas práticas de código, garantindo o alinhamento com os padrões de arquitetura, tratamento de erros, observabilidade e documentação. Toda a interação deve ser feita em português.

---

**Problemas:**
Após confirmar a disponibilidade de slots para o fotografo através de GET /listings/photo-session/slots?from=2025-11-01&to=2025-11-04&period=MORNING&page=1&size=20&sort=start_asc&listingId=3&timezone=America/Sao_Paulo e receber como reposta:
{
    "data": [
        {
            "slotId": 23236831280,
            "photographerUserId": 5,
            "slotStart": "2025-11-01T11:00:00Z",
            "slotEnd": "2025-11-01T15:00:00Z",
            "status": "AVAILABLE"
        },
        {
            "slotId": 23236917680,
            "photographerUserId": 5,
            "slotStart": "2025-11-02T11:00:00Z",
            "slotEnd": "2025-11-02T15:00:00Z",
            "status": "AVAILABLE"
        },
        {
            "slotId": 23237004080,
            "photographerUserId": 5,
            "slotStart": "2025-11-03T11:00:00Z",
            "slotEnd": "2025-11-03T15:00:00Z",
            "status": "AVAILABLE"
        }
    ],
    "pagination": {
        "page": 1,
        "limit": 20,
        "total": 3,
        "totalPages": 1
    }
}

chamei o endpoint POST /listings/photo-session/reserve com o seguinte payload:
{
  "listingId": 3,
  "slotId": 23236917680
}
para reservar o slot.

Recebi como resposta:
{
    "code": 500,
    "details": null,
    "message": "failed to create booking"
}
e a mensagem de erro no log está em logs/log_example.json.

1) analise a causa raiz do problema e verifique se os outros endpoints relacionados ao CRUD de reservas de sessões fotográficas estão funcionando corretamente, em especiao repositórios e serviços.
  1.1)o script usado para a criação das tabelas do banco de dados está em scripts/db_creation.sql

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