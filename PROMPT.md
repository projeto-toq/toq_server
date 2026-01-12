Segundo a regra de negócios o realtor deve poder favoritar listings.

Assim deveria haver um endpoint para:
- adicionar um listing aos favoritos - `POST /listings/favorites/add` com o body contendo o `listingIdentityId`
- remover um listing dos favoritos - `POST /listings/favorites/remove` com o body contendo o `listingIdentityId`
- listar os listings favoritos - `GET /listings/favorites` deve retornar a lista de listings favoritos do realtor autenticado. A respsota deve ser semelhante a do endpoint `GET /listings`, mas contendo apenas os listings favoritados pelo realtor.

O modelo de dados não está preparado para isso, então será necessário criar uma nova tabela com os campos adequados.

O endpoint `POST /listings/detail` tem como parte de sua resposta os dados de performande do listinfg
 "performanceMetrics": {
    "favorites": 0,
    "shares": 0,
    "views": 0
  },
Entretanto favorites não está sendo hidratado na resposta. É necessário corrigir isso para que o campo `favorites` retorne a quantidade correta de vezes que o listing foi favoritado por qualquer realtor.

Os endpoints `POST /listings/detail` e `GET /listings` deve ter um campo adicional na resposta chamado `favorite` booleano que indica se o realtor autenticado favoritou ou não o listing em questão.

Analise o código atual e proponha o plano conforme o `AGENTS.md`.

