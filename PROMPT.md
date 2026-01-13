O sistema de observabilidade é composto por grafana, prometheus, loki, tempo etc. todos rodando em containers docker com docker-compose.yml na raiz do projeto.

No Grafana temos o dashboard `TOQ Server - Traces` que apresente os seguintes problemas:
1. Painel `Traces do Servico` temos
    1.1 A coluna `Service` que não faz sentido, sempre será toq_server. deve ser removida.
2. O filtro `Operação/Rota` não está filtrando os traces corretamente, por exemplo ao selecionar a rota `/api/v2;listing` ele retorna traces que não são nem rotas http

Outras equipes passram horas e muitas tentativas e erros e não corrrigiram o problema.

Portanto busque todas as infromações que precisa, logs, dados reais das bases, configurações, manuais etc, para ter certeza da causa raiz e só então proponha o plano conforme o `AGENTS.md`.