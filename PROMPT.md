# Filtros de Busca de Imóveis e Correções de API

## Correções Necessárias

## Endpoint `/listings/catalog`

**Problema:** O endpoint está fora do padrão REST. Atualmente usa o método `GET` com dados sendo passados via body da requisição.

- **Mobile:** Funciona sem problemas
- **Web:** Não aceita body em requisições GET (comportamento padrão de navegadores)

**Solução Recomendada:** O endpoint deve ser alterado para usar **POST** ao invés de GET, permitindo o envio de filtros no body da requisição.
---

## Filtros de Busca de Imóveis

O endpoint `GET /listings` deve incluir, além dos filtros já existentes, os seguintes filtros adicionais:

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

Busque todas as informações que precisa consultando o código e as configurações reais, não confiando na documentação, para ter certeza da situação, e só então proponha o plano de correção.