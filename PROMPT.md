O sistema de observabilidade é composto por grafana, prometheus, loki, tempo etc. todos rodando em containers docker com docker-compose.yml na raiz do projeto.

No Grafana temos o dashboard `TOQ Server - Traces` que apresente os seguintes problemas:
1. Painel `Traces do Servico`com No data found in response;

Outras equipes passram horas e muitas tentativas e erros e não corrrigiram o problema.

Portanto busque todas as infromações que precisa, logs, dados reais das bases, configurações, manuais etc, para ter certeza da causa raiz e só então proponha o plano conforme o `AGENTS.md`.