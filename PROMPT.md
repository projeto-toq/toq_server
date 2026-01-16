O endpoiint `GET /listings` tem como parametro de busca :
zipCode
city
neighborhood
street
number
complement
state
que compoem o endereço do imóvel.

Devemos alterar a forma de busca para que o enpoint aceite um único campo chamado address que deve conter o endereço completo do imóvel em uma única string.

O serviço deve fazer uma Wildcard Search, ou se possível Fuzzy Search, no banco de dados utilizando o campo address, que deve ser uma concatenação dos campos zipCode, city, neighborhood, street, number, complement e state. 

A resposta da endpoint deve continuar a mesma, retornando a lista de imóveis que batem com o critério de busca.


Busque todas as informações que precisa consultando o código e as configurações reais, não confiando na documentação, para ter certeza da situação, e só então proponha o plano de correção.