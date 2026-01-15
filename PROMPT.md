Filtros de Busca de Imóveis via endpoint `GET /listings` teve a inclusão destes novos filtros:

### Filtros Básicos

| Filtro | Descrição | Valores Possíveis |
|--------|-----------|-------------------|
| **Tipos de imóvel** | Categoria do imóvel | Casa, Apartamento, Terreno, etc. |
| **Tipo de transação** | Modalidade de negócio | Venda, Aluguel |
| **Uso do imóvel** | Finalidade de uso | `RESIDENCIAL`, `COMERCIAL` |
| **Aceita permuta** | Se aceita troca | `SIM`, `NÃO` |
| **Aceita financiamento** | Se aceita financiamento | `SIM`, `NÃO` |

### Filtros Especiais

| Filtro | Descrição |
|--------|-----------|
| **Apenas imóveis novos** | Filtra somente imóveis recém-cadastrados |
| **Apenas com alteração de preço** | Filtra imóveis que tiveram mudança de preço recente |
| **Apenas imóveis vendidos** | Filtra imóveis já vendidos |


entretano os filtros estão com funcionamento incorreto:
- TransactionType: pesquisei apenas ALUGUEL e retornou imóvel que está apenas para VENDA;
- onlyPriceChanged: pesquisei e retornou a lista mesmo sem o valor ter sido alterado;
- Tipo de imóvel: pesquisei pelo tipo de imóvel apartamento e foram retornados todos os que estavam publicados, inclusive do tipo Casa;

Busque todas as informações que precisa consultando o código e as configurações reais, não confiando na documentação, para ter certeza da situação, e só então proponha o plano de correção.