### Resumo e Refatoração: Engenheiro de Software Go Sênior

Este documento descreve as instruções para atuar como um engenheiro de software Go sênior, focando na análise de um problema e na proposição de uma solução detalhada, seguindo a arquitetura hexagonal e boas práticas de código, garantindo o alinhamento com os padrões de arquitetura, tratamento de erros, observabilidade e documentação. Toda a interação deve ser feita em português.

---

**Problemas:**
A agenda dos fotografos foi implementada com uma lógica muito complicada.
Existe uma tabela, com horizonte de 3 meses com todas as disponibilidades de horários dos fotografos. photographer_default_availability. Processo contrário ao processo de uma agenda normal de compromissos, todas os slots estão livres exceto quanto existe um bloqueio ou uma sessão agendada.
O correto seria que:
- cada fotografo tivesse uma agenda com todos os compromissos agendados e bloqueios, e a disponibilidade fosse calculada a partir disso. Sem bloqueios ou sessões agendadas, o fotografo estaria disponível.
- Como o proprietário necessita buscar quais dias e horários teremos um fotografo disponível, via GET /listings/photo-session/slots, e o slot para cada sessão de fotos é de 4 horas, definido por env.yaml, o serviço deveria buscar os fotógrafos que possuem horários livres de 4 horas consecutivas entre seus compromissos agendados e bloqueios.
- a criação da agenda base do fotografo, efetuada durante a criação do usuário do sistema  (create_system_user.go) com roleSlug photographer, deveria colocar como bloqueio padrão todos os horários fora do expediente, 08:00 - 18:00 como default mas definido por variáveis no env.yaml, e deixar livres os horários dentro do expediente.
- Dias feriados, nacionais, estaduais e municipais, também deveriam ser bloqueados na criação da agenda padrão.
- para definir dias de feriados que afetam o fotografo, é necessário que o sistema de calendarios holiday_service na tabela holiday_calendars, ao invés de tratar o campo city_ibge como codigo ibge trate como city, assim como é feito na tabela users. Assim é possével associar o calendário de feriados ao usuário fotografo baseado na cidade do usuário (city + state), ou no estado do usuário (state). Nacional afeta a todos.
- eventuais alterações no banco de dados serão feitas manualmente, apenas indique as alteraçoes necessárias.
- eventuais códigos morto ou não utilizados devem ser removidos.

**Solicitação:** Analise o problema, **leia o código** envolvido, **ache a causa raiz** e proponha um plano detalhado para a implementação da solução, após ler o o manual do projeto em docs/toq_server_go_guide.md.

### **Instruções para a Proposição do Plano**

- **Ação:** Apenas a análise e a geração do plano são solicitadas. **Nenhum código deve ser implementado**.
- **Análise:** O problema e os requisitos devem ser analisados cuidadosamente. O código e arquivos de configuração existentes devem ser revisados para um plano preciso. Não faça suposições e confirme todos os detalhes necessários.
- **Plano:** Um plano detalhado deve ser apresentado, incluindo a descrição da arquitetura proposta, as interfaces, a estrutura de diretórios e a ordem de execução das etapas.
- **Qualidade do Plano 1:** O plano deve ser completo, sem o uso de _mocks_ ou soluções temporárias. Caso seja extenso, deve ser dividido em etapas implementáveis.
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