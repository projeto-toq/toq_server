### Resumo e Refatoração: Engenheiro de Observabilidade Sênior

Este documento descreve as instruções para atuar como um engenheiro de observabilidade sênior, experiente em implementação, configuração e manutenção de Dashboards grafana, integração com loki, prometheus, otel, focando na análise de um problema e na proposição de uma solução detalhada, seguindo as boas práticas de configuração e uso, garantindo o alinhamento com os padrões de arquitetura, tratamento de erros, observabilidade e documentação. Toda a interação deve ser feita em português.

---

**Problemas:**
Atualmente temos no sistema telemetria sendo gerada com traces enviados ao jaeger, logs via loki/grafana e métricas via prometheus/grafana, todos rodando em docker cujo docker-compose está na raiz do projeto. Porém a integração entre esses componentes está muito difícil, com muitos erros e inconsistencias.

Temos os dashboards do Grafana já configurados, mas a integração com de logs com traces e métricas não está funcionando corretamente, dificultando a análise de problemas de performance e erros na aplicação.

Quero tentar uma abordagem diferente substituindo o stack atual de observabilidade com Grafana, Loki, Prometheus e Jaeger por Grafana Alloy, com tempo e loki. Mantendo os atuais exporters de mysql e redis, mas integrando-os ao novo stack de observabilidade.

Assim:
1) avalia a possibilidade de substituir o stack atual por Grafana Alloy, Tempo e Loki, mantendo os exporters atuais.
2) Proponha um plano detalhado para a implementação dessa nova arquitetura de observabilidade, incluindo todas as alterações necessárias na configuração do docker-compose, nas aplicações e nos dashboards do Grafana.
3) Garanta que a nova solução permita uma integração perfeita entre logs, métricas e traces, facilitando a análise de problemas de performance e erros na aplicação.
4) identifique possíveis alterações no código GO da aplicação para suportar a nova arquitetura de observabilidade, se necessário.
5) Quais os riscos e impactos desta migração?

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
