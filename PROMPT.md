Estamos recebendo o seguinte erro no log:
```json
{"time":"2026-01-13T10:19:46.350466709Z","level":"ERROR","msg":"HTTP Error","request_id":"d3caf5c87334136980afc8b334582512","request_id":"d3caf5c87334136980afc8b334582512","method":"POST","path":"/api/v2/listings/detail","status":500,"duration":7223324,"size":46,"client_ip":"45.188.201.28","user_agent":"insomnia/11.6.1","trace_id":"d2b8136191387aaf33d821a5debe60a1","span_id":"7696cb49b7d6418f","user_id":7,"user_role_id":11,"function":"github.com/projeto-toq/toq_server/internal/core/utils.InternalError","file":"/codigos/go_code/toq_server/internal/core/utils/http_errors.go","line":248,"stack":["github.com/projeto-toq/toq_server/internal/core/utils.InternalError (http_errors.go:248)"],"error_code":500,"error_message":"Internal server error","errors":["HTTP 500: Internal server error"]}
{"time":"2026-01-13T10:19:46.350563511Z","level":"ERROR","msg":"HTTP Error","request_id":"d810532abce56df678886462699fe9e3","request_id":"d810532abce56df678886462699fe9e3","method":"POST","path":"/api/v2/listings/detail","status":500,"duration":7377887,"size":46,"client_ip":"45.188.201.28","user_agent":"insomnia/11.6.1","trace_id":"c0bb6e9ae36e7a65b71fa26641edff7a","span_id":"c73bde02dada02ed","user_id":7,"user_role_id":11,"function":"github.com/projeto-toq/toq_server/internal/core/utils.InternalError","file":"/codigos/go_code/toq_server/internal/core/utils/http_errors.go","line":248,"stack":["github.com/projeto-toq/toq_server/internal/core/utils.InternalError (http_errors.go:248)"],"error_code":500,"error_message":"Internal server error","errors":["HTTP 500: Internal server error"]}
```

Aqui temos 2 atividades que precisam ser realizadas:
1. analisar o erro e identificar a causa raiz do problema.
2. a mesangem de erro está muito genérica, precisamos melhorar a mensagem de erro para que possamos ter mais detalhes sobre o que está causando o erro 500.

Analise o código proponha o plano conforme o `AGENTS.md`.

