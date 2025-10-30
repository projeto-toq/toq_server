### Resumo e Refatoração: Engenheiro de Software Go Sênior

Este documento descreve as instruções para atuar como um engenheiro de software Go sênior, focando na análise de um problema e na proposição de uma solução detalhada, seguindo a arquitetura hexagonal e boas práticas de código, garantindo o alinhamento com os padrões de arquitetura, tratamento de erros, observabilidade e documentação. Toda a interação deve ser feita em português.

---

**Problemas:**
A agenda dos fotografos, gerenciadas pelos serviço service/photo_session_service/*.go tem como logica criar entradas de bloqueio para horários fora dos horários de trabalho conforme definidi em env.yaml:
```env.yaml
photo_session:
  slot_duration_minutes: 240
  slots_per_period: 1
  business_start_hour: 8
  business_end_hour: 18
  morning_start_hour: 8
  afternoon_start_hour: 14
  photographer_horizon_months: 3
  photographer_timezone: "America/Sao_Paulo"
  photographer_agenda_refresh_interval_hours: 24
```
Entretanto esta logica é ineficiente, pois gera muitas entradas de bloqueio desnecessárias, e exige serviços de criação para manter horizontes de bloqueios.
O correto é que não haja entradas de bloqueios para horários fora do horário de trabalho, e que a consulta da agenda do fotografo leve em consideração os horários de trabalho definidos no env.yaml.
Adicionalmente, os feriados, gerenciados pelo serviço service/holiday_service/*.go, não são considerados nas consultas as agendas do fotografo via GET /photographer/agenda, mesmo a reposta do endpoint estar preparada para os feriados.
```
{
            "entryId": 8,
            "photographerId": 5,
            "entryType": "BLOCK",
            "source": "ONBOARDING",
            "start": "2025-11-02T00:00:00-03:00",
            "end": "2025-11-03T00:00:00-03:00",
            "status": "BLOCKED",
            "groupId": "onboarding-2025-11-02",
            "isHoliday": false,
            "isTimeOff": false,
            "timezone": "America/Sao_Paulo"
        }
```
A tabela toq_db.photographer_holiday_calendars que deveria ter o link estre o usuário fotografo e os feriados, através do endereço city/state do usuário e city/state da tabela de feriados não está sendo populada. Por consequencia o GET /photographer/agenda tem como filtro o campo holidayCalendarIds que não deveria existir. Esta logica de tabela de link de feriados com ids previamente polupados é ineficiente.por cidade/estado é desnecessária, O correto é que o ao consultar a agenda do fotografo exista uma consulta trazendo os feriados nacionais, do estado e da cidade do fotógrafo, segundo seu endereço. isso em tempo de consulta e não persistido em tabela específica.


**Solicitação:** Analise o problema, **leia o código** envolvido, **ache a causa raiz** e proponha um plano detalhado para a implementação da solução, após ler o o manual do projeto em docs/toq_server_go_guide.md.

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