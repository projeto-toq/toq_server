### Resumo e Refatoração: Engenheiro de Observabilidade Sênior

Este documento descreve as instruções para atuar como um engenheiro de observabilidade sênior, experiente em implementação, configuração e manutenção de Dashboards grafana, integração com loki, prometheus, otel, focando na análise de um problema e na proposição de uma solução detalhada, seguindo as boas práticas de configuração e uso, garantindo o alinhamento com os padrões de arquitetura, tratamento de erros, observabilidade e documentação. Toda a interação deve ser feita em português.

---

**Problemas:**
Atualmente temos no sistema telemetria sendo gerada com traces enviados ao jaeger, logs via loki/grafana e métricas via prometheus/grafana, todos rodando em docker cujo docker-compose está na raiz do projeto. Porém a integração entre esses componentes não está otimizada para facilitar a investigação de problemas. A proposta é melhorar essa integração, especialmente focando na correlação entre logs e traces, além de aprimorar os filtros disponíveis nos dashboards do Grafana para facilitar a análise dos dados.

Precismaos ter o conjunto de informações básicas para monitorar o sistema que é um rest-api, portanto precisamos garantir ao menos:
1) dashboard de sinais vitais GO e host (CPU, memória, concorrencia, GC, etc)
2) monitoramento de database mysql e do redis
3) dashboard de logs estruturados com filtros por request_id, trace_id, path, method, etc
4) dashboard de traces com correlação com logs via trace_id e/ou request_id
5) monitoramento de HTTP Request Rate, HTTP Requests In Flight, HTTP Request Duration e request in flight,  latência, erros 4xx/5xx, traffic qps, erros por segundo etc.
6) algum outro dashboard que lhe pareça relevante para monitorar uma aplicação REST API em Go.

Ja existe alguns dashboards criados no grafana mas se necessário podem ser completamente alterados ou removidos.

Assim:
1) verifique as metricas que já estão sendo coletadas e quais estão faltando ou precisam melhorias.
  1.1) Estas devem ser verificadas no prometheus, para garantir que estão sendo coletadas corretamente.
2) verifique os campos do log que estão sendo gerados e quais estão faltando ou precisam melhorias
  2.1) Estes devem ser verificados no loki/grafana, para garantir que estão sendo coletados corretamente.
3) verifique os traces que estão sendo gerados e quais estão faltando ou precisam melhorias.
  3.1) Estes devem ser verificados no jaeger/grafana, para garantir que estão sendo coletados corretamente.
4) verifique as metricas do mysql e do redis que estão sendo coletadas e quais estão faltando ou precisam melhorias
  4.1) Estas devem ser verificadas no prometheus, para garantir que estão sendo coletadas corretamente.
5) Os dashboards do grafana possuem filtros que não funcionam ou estão incompletos, corrija estes filtros para garantir que todos os dados possam ser filtrados corretamente.
6) garanta que todos os dashboards estejam utilizando as melhores práticas de visualização de dados, como uso adequado de gráficos, tabelas, alertas e cores para facilitar a interpretação dos dados.
7) os valores environment, service.name e service.version estão sendo exportados, mas não deveriam, pois isso é uma configuração interna da aplicação. Corrija isso.


**Solicitação:** Analise o problema, **leia as configurações e docker compose** envolvido, **ache a causa raiz** e proponha um plano detalhado para a implementação/refatoração da solução, após ler o o manual do projeto em docs/toq_server_go_guide.md. O plano deve prever já toda a alteração de configuração necessária, para que quando aprovado, possa ser implementado diretamente.

### **Instruções para a Proposição do Plano**

- **Ação:** Apenas a análise e a geração do plano são solicitadas. **Nenhum código deve ser implementado**.
- **Análise:** O problema e os requisitos devem ser analisados cuidadosamente. O código e arquivos de configuração existentes devem ser revisados para um plano preciso. Não faça suposições e confirme todos os detalhes necessários.
- **Plano:** Um plano detalhado deve ser apresentado, incluindo a descrição da arquitetura proposta, as interfaces, a estrutura de diretórios e a ordem de execução das etapas.
- **Qualidade do Plano:** O plano deve ser **extremamente prescritivo**, sem o uso de _mocks_ ou soluções temporárias. Para cada arquivo novo ou alterado, inclua um **esqueleto de código (code skeleton)** ou a **assinatura completa da função (full function signature)**.
- **Acompanhamento:** As etapas já planejadas e as próximas a serem analisadas devem ser sempre informadas para acompanhamento.
- **Ambiente:** O plano deve considerar que estamos em ambiente de desvolvimento, portanto não deve haver back compatibility, migração de dados, preocupação com janela de manutenção ou _downtime_.

- **Documentação:** Deve haver ajustes no documento docs/sre_guide.md para refletir as mudanças implementadas, garantindo que o guia de observabilidade esteja completo e atualizado.
---

### **Regras Obrigatórias de Análise e Planejamento**

#### 1. Arquitetura e Fluxo de Código
- **Arquitetura:** A solução deve seguir estritamente a **Arquitetura Hexagonal**.
- **Fluxo de Chamadas:** As chamadas de função devem seguir a hierarquia `Handlers` → `Services` → `Repositories`.
- **Injeção de Dependência:** O padrão de _factories_ deve ser usado para a injeção de dependências.
- **Localização de Repositórios:** Os repositórios devem ser localizados em `/internal/adapter/right/mysql/` e deve fazer uso dos convertess para mapear entidades de banco de dados para entidades e vice versa.
- **Transações SQL:** Todas as transações de banco de dados devem utilizar `global_services/transactions`.
