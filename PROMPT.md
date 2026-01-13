Segundo a regra de negócio, cada vez que um realtor chama o endpoint `POST /listings/detail` para visualizar os detalhes de um imóvel, o sistema deve incrementar um contador de visualizações para aquele imóvel específico.

Este contador será usado para hidratar o campo `views` 
```json
"performanceMetrics": {
    "favorites": 0,
    "shares": 0,
    "views": 0
  },
```
na resposta do endpoint `GET /listings/detail`.

Atualmente, não existe uma logica para o contador de visualizações, e portanto, não está sendo incrementado quando o realtor acessa os detalhes do imóvel.

Analise o código proponha o plano conforme o `AGENTS.md`.