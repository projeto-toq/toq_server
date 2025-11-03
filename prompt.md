### Resumo e Refatoração: Engenheiro de Software Go Sênior

Este documento descreve as instruções para atuar como um engenheiro de software Go sênior, focando na análise de um problema e na proposição de uma solução detalhada, seguindo a arquitetura hexagonal e boas práticas de código, garantindo o alinhamento com os padrões de arquitetura, tratamento de erros, observabilidade e documentação. Toda a interação deve ser feita em português.

---

**Problemas:**
Atualmente temos no sistema telemetria sendo gerada com traces enviados ao jaeger, logs via loki/grafana e métricas via prometheus/grafana, todos rodando em docker cujo docker-compose está na raiz do projeto. Porém a integração entre esses componentes não está otimizada para facilitar a investigação de problemas. A proposta é melhorar essa integração, especialmente focando na correlação entre logs e traces, além de aprimorar os filtros disponíveis nos dashboards do Grafana para facilitar a análise dos dados.

Precismaos ter o conjunto de informações básicas para monitorar o sistema que é um rest-api, portanto precisamos garantir ao menos:
1) dashboard de sinais vitais GO e host (CPU, memória, concorrencia, GC, etc)
2) monitoramento de database mysql
  2.1) talvez seja importante ter um coletor/exporter específico para mysql
  2.2) utilzar o que já existe no projeto 
  2.3) existe cache com redis, necessário algo específico para ele?
3) dashboard de logs estruturados com filtros por request_id, trace_id, path, method
4) dashboard de traces com correlação com logs via trace_id
5) monitoramento de HTTP Request Rate, HTTP Requests In Flight, HTTP Request Duration e request in flight,  latência, erros 4xx/5xx, traffic qps, erros por segundo etc.
6) algum outro dashboard que lhe pareça relevante para monitorar uma aplicação REST API em Go.

Ja existe alguns dashboards criados no grafana, porém estão mal construídos e devem ser apagados. Não tem uso, portanto pode apagar e substitui-los pelos novos.

Assim:
1) verifique as metricas que já estão sendo coletadas e quais estão faltando ou precisam melhorias.
2) verifique os campos do log que estão sendo gerados e quais estão faltando ou precisam melhorias
3) verifique os traces que estão sendo gerados e quais estão faltando ou precisam melhorias.
4) verifique se falta algo para coletar/exportar metricas do mysql
5) verifique se falta algo para coletar/exportar metricas do redis
6) crie um plano detalhado para implementar os dashboards necessários, incluindo as queries do grafana, os painéis, as variáveis de filtro e as derived fields para correlação entre logs e traces.
7) sugira qualquer melhoria necessária no middleware de telemetria, campos de logs e traces para garantir que todos os dados necessários estão sendo coletados e exportados corretamente.
8) todos os dashboards e painéis devem ter explicação clara do objetivo e instruções de uso.
9) em docs/observability/ crie um documento com a explicação dos dashboards criados, como usá-los e como interpretar os dados apresentados. Apague o atual documento logs.md. só deve existir um guia completo e atualizado.
10) verifique porque existem 2 diretorios dashboards em grafana/ e mantenha somente um.

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