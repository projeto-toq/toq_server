Atue como um desenvolvedor GO Senior e faça toda a interação em português.

---

## 1) Objetivo do Pedido
- Tipo: Somente análise e apresentaçao do plano para aprovação (sem implementação).

## 2) Requisição

Com a evolução do sistema Toq Server, é essencial garantir uma estratégia robusta de monitoramento e correlação entre métricas, logs e traces para facilitar a detecção e resolução de problemas.

Abaixo está uma sugestão detalhada para implementar essa estratégia utilizando Grafana, Prometheus, Loki e Jaeger.

É uma sugestão de quem tem uma visão superficial da realidade no nosso projeto. Assim, analise cuidadosamente e prepare um plano para implementação, considerando ajustes conforme necessário para o contexto específico do Toq Server, considerando o código atual, infraestrutura, docker-compose onde definimos todos os serviços exceto o proprio Toq Server, e definições do manual do projeto `docs/toq_server_go_guide.md`.

Os atuais dashboards do Grafana devem ser apagados e recriados com base nesta nova estratégia.

Caso sejam necessário incluir novos pontos de instrumentação no código Go (métricas, logs, traces), inclua isso no plano de implementação.

---

## Estratégia de Monitoramento e Correlação (Grafana + Prometheus + Loki + **Jaeger**)

A estratégia principal se mantém: usar as **Golden Signals** (Prometheus) para *alertar* sobre problemas e os **Traces** (Jaeger) para *diagnosticar* a causa raiz.

### 1. Golden Signals (Indicadores Essenciais de Saúde)

**Origem da Métrica:** Prometheus (coletado do sistema Go via OpenTelemetry/Prometheus client).

| Categoria | Indicador (Métrica) | Por que é relevante / Ação |
| :--- | :--- | :--- |
| **Latência** | P95 e P99 de Latência por Endpoint. | Mede a experiência dos piores usuários e identifica gargalos. Use para alertas de lentidão. |
| **Tráfego** | QPS (Queries Per Second) total e por Endpoint. | Mede volume e demanda, essencial para planejamento de capacidade. |
| **Erros** | Taxa de Erro 5xx (Server Errors) e Distribuição de Códigos HTTP. | **Indicador de alerta crucial.** Sobe -> ALERTA. |
| **Saturação** | Utilização da CPU e Uso de Memória/Goroutines do Go Runtime. | Ajuda a identificar problemas de recurso antes que a latência dispare. |

---

### 2. Correlação Métrica-Trace-Log com Jaeger

O coração da observabilidade é a capacidade de "saltar" entre as ferramentas. Como seu sistema Go está instrumentado com OpenTelemetry, ele deve injetar o **Trace ID** e o **Span ID** em tudo:

1.  **Métricas (Prometheus/Grafana):** O Prometheus monitora as Golden Signals.
2.  **Logs (Loki):** Cada log gerado no sistema Go precisa conter os campos `trace_id` e `span_id`.
3.  **Traces (Jaeger/OTel):** O Jaeger armazena o caminho completo da requisição.

#### Novo Indicador Crucial para Correlação:

| Indicador | Categoria | Como usar para Correlação |
| :--- | :--- | :--- |
| **Latência por Componente/Span** | Tracing (Jaeger) | Qual porcentagem da latência total é gasta no Handler HTTP, na chamada de Banco de Dados, ou em um serviço externo. **Isto é o que você investiga no Jaeger após um alerta.** |
| **Trace ID em Logs** | Logs (Loki) | Use a `Trace ID` de um *span* no Jaeger para buscar logs no Loki. |

### Fluxo de Análise e Debug com Jaeger 

A adaptação é puramente de *data source* no Grafana. A lógica de investigação se mantém, garantindo a correlação entre as 3 pilares:

| Passo | Ação | Ferramenta | Objetivo |
| :--- | :--- | :--- | :--- |
| **1. Alerta** | A **Taxa de Erro 5xx** (Prometheus) ultrapassa o limite. | Grafana (Alerting) | Detectar o problema de saúde geral. |
| **2. Pivô para Trace (Métrica -> Trace)** | No painel de Erros (Métrica no Grafana), você deve ter um painel mostrando uma amostra dos **Trace IDs** das requisições com erro (isto requer que você exponha o Trace ID como um label nas métricas, ou mais facilmente, no Log). | Grafana + Jaeger/Loki | Obter um *exemplo* de `Trace ID` de uma requisição que falhou. |
| **3. Investigação Detalhada do Trace** | Copie o `Trace ID` e cole-o na interface do Jaeger ou use o Data Source do Jaeger no Grafana. | **Jaeger** | Ver o caminho completo do serviço. Identificar o **Span** que gerou o erro (ex: `DB Query Error`, `HTTP 500 from ServiceX`). |
| **4. Pivô para Log (Trace -> Log)** | Dentro do Span que falhou no Jaeger, o Log geralmente contém o `Trace ID` e `Span ID`. Use o `Trace ID` para buscar todos os logs relacionados no Loki. | **Loki** | Ver o *stack trace* exato ou a mensagem de erro que a aplicação Go gerou no momento da falha. |

**Como configurar a correlação no Grafana (Jaeger):**

1.  **DataSource:** No seu Grafana, o Jaeger deve ser configurado como uma fonte de dados de **Tracing**.
2.  **Links de Campos (Field Links):** No Data Source do Loki, configure um *Field Link* para que o campo `trace_id` (que está nos seus logs) aponte diretamente para o Data Source do **Jaeger**, permitindo o "salto" do Log para o Trace com um clique.
3.  **Annotations (Opcional):** Você pode usar o Loki ou Jaeger para criar *annotations* (marcas) nos gráficos do Prometheus no Grafana, correlacionando eventos (ex: *deploy*, picos de erro) nos gráficos de métricas.


- Documentação de referência: `docs/toq_server_go_guide.md`


## 3) Requisitos
  - Observabilidade alinhada (logs/traces/metrics existentes)
  - Performance/concorrência (se aplicável)
  - Sem back compatibility/downtime (ambiente de desenvolvimento)

## 4) Artefatos a atualizar (se aplicável)
- Código (handlers/services/adapters)
- Swagger (anotações no código)
- Documentação (README/docs/*)
- Observabilidade (métricas/dashboards) — apenas quando estritamente pertinente
- Não crie/altere arquivos de testes

## 5) Arquitetura e Fluxo (resumo)
- Siga o guia: `docs/toq_server_go_guide.md` (Seções 1–4 e 7–11).
- Em alto nível: `Handlers` → `Services` → `Repositories`; DI por factories; converters nos repositórios; transações via serviço padrão.

## 6) Erros e Observabilidade (resumo)
- Siga o guia: `docs/toq_server_go_guide.md` (Seções 5, 6, 8, 9 e 10).
- Pontos-chave: spans só fora de handlers HTTP; `SetSpanError` em falhas de infra; handlers usam `http_errors.SendHTTPErrorObj`; adapters retornam erros puros; logs com severidade correta.

## 7) Dados e Transações
- Tabelas/entidades afetadas: <listar>
- Necessidade de transação: <sim/não; justificar>
- Regras de concorrência e idempotência: <se aplicável>

## 8) Interfaces/Contratos
- Endpoints HTTP envolvidos: <listar caminhos e métodos>
- Payloads de request/response: <resumo>
- Erros de domínio esperados e mapeamento HTTP: <ex.: ErrX → 409>

## 9) Entregáveis Esperados do Agente
- Análise detalhada com checklist dos requisitos e plano por etapas (com ordem de execução).
- Quality gates rápidos no final: Build, Lint/Typecheck e mapeamento Requisito → Status.
- Não execute git status, git diff nem go test.

## 10) Restrições e Assunções
- Ambiente: desenvolvimento (sem back compatibility, sem janela de manutenção, sem migração).
- Sem uso de mocks nem soluções temporárias em entregas finais.

## 11) Anexos e Referências
- Arquivos relevantes: `docs/toq_server_go_guide.md`.

---

### Regras Operacionais do Pedido (a serem cumpridas pelo agente)

1) Antes de propor/alterar código:
- Revisar arquivos relevantes (código e configs); evitar suposições sem checagem.
- Produzir um checklist de requisitos explícitos e implícitos, mantendo-o visível.

2) Durante a análise/implementação:
- Seguir o guia `docs/toq_server_go_guide.md` (arquitetura, erros/observabilidade, transações).
- Manter adapters com erros “puros”; sem HTTP/semântica de domínio nessa camada.
- A documentação Swagger/docs deve ser criada por comentários, em inglês em DTO/Handler e execução de make swagger. Sem alterações manuais no swagger.yaml/json.

3) Após mudanças:
- Relatar “requirements coverage” (Requisito → Done/Deferred + motivo).

4) Em casos de dúvidas: consultar o requerente.

5) Caso a tarefa seja grande/demorada, dividir em fases menores e entregáveis curtos.

6) Caso haja remoção de arquivos limpe o conteúdo e informa a lista para remoção manual.

7) Comentários internos em português; docstrings de funções em inglês; Swagger por anotações no código.