### Resumo e Refatoração: Engenheiro de Observabilidade Sênior

Este documento descreve as instruções para atuar como um engenheiro de observabilidade sênior, experiente em implementação, configuração e manutenção de Dashboards grafana, integração com loki, prometheus, otel, focando na análise de um problema e na proposição de uma solução detalhada, seguindo as boas práticas de configuração e uso, garantindo o alinhamento com os padrões de arquitetura, tratamento de erros, observabilidade e documentação. Toda a interação deve ser feita em português.

---

**Problemas:**
Atualmente o dashboard de observabilidade do Grafana TOQ Server - Logs e Traces não possui filtro por severidade de logs (debug, info, warn, error). Isso dificulta a análise de problemas específicos, especialmente quando se deseja focar em logs de erro ou avisos críticos.
Além disso o processo de correlação de logs e traces requer que após localizar o erro no log seja necessário copiar o trace_id e colar na busca de traces, o que é um processo manual e propenso a erros.

1) será que não seria elhor termos 2 dashboards separados, um para logs e outro para traces?
    1.1) no dashboard de logs teríamos mais espaço de visualização dos logs e busca por filtros dedicados, com filtros por severidade, data/hora além dos atuais.
    1.2) uma opção de clicar no trace_id ou em algum link para abrir o dashboard de traces já filtrado por esse trace_id.
2) no dashboard de traces, teríamos um filtro por trace_id, e também filtros por data/hora, operação, duração, status (erro/sucesso).
    2.1) ao clicar em um trace específico, o painel de detalhes do trace seria exibido, mostrando todas as informações relevantes de forma clara e organizada.
    2.2) uma opção de clicar em um link para abrir o dashboard de logs já filtrado por esse trace_id.


### **Instruções para a Proposição do Plano**

- **Ação:** Apenas a análise e a geração do plano são solicitadas. **Nenhum código deve ser implementado**.
- **Análise:** O problema e os requisitos devem ser analisados cuidadosamente. O código e arquivos de configuração existentes devem ser revisados para um plano preciso. Não faça suposições e confirme todos os detalhes necessários.
- **Plano:** Um plano detalhado deve ser apresentado, incluindo a descrição da arquitetura proposta, as interfaces, a estrutura de diretórios e a ordem de execução das etapas.
- **Qualidade do Plano:** O plano deve ser **extremamente prescritivo**, sem o uso de _mocks_ ou soluções temporárias. Para cada arquivo novo ou alterado, inclua um **esqueleto de código (code skeleton)** ou a **assinatura completa da função (full function signature)**.
- **Acompanhamento:** As etapas já planejadas e as próximas a serem analisadas devem ser sempre informadas para acompanhamento.
- **Ambiente:** O plano deve considerar que estamos em ambiente de desvolvimento, portanto não deve haver back compatibility, migração de dados, preocupação com janela de manutenção ou _downtime_.

- **Documentação:** Deve haver ajustes no documento docs/sre_guide.md para refletir as mudanças implementadas, garantindo que o guia de observabilidade esteja completo e atualizado.
---
