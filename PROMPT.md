O endpoint `GET /listings` tem como resposta:
```json
{
  "data": [
    {
      "activeVersionId": 0,
      "complex": "string",
      "complexId": "string",
      "description": "string",
      "draftVersionId": 0,
      "favoritesCount": 0,
      "id": 0,
      "isFavorite": true,
      "listingIdentityId": 0,
      "listingUuid": "string",
      "number": "string",
      "price": 0,
      "propertyType": {
        "code": 0,
        "label": "string",
        "propertyBit": 0
      },
      "status": "string",
      "title": "string",
      "userId": 0,
      "version": 0,
      "zipCode": "string"
    }
  ],
  "pagination": {
    "limit": 0,
    "page": 0,
    "total": 0,
    "totalPages": 0
  }
}
```

É necessário que a respsota contenha a estrutura abaixo além dos campos já existentes, sem duplicações:
```json
      "city": "São Paulo",
      "complement": "apto 82",
      "description": "Apartamento amplo com três suítes e vista livre.",
      "listingIdentityId": 123,
      "neighborhood": "Moema",
      "number": "1234",
      "propertyType": {
        "code": 0,
        "label": "string",
        "propertyBit": 0
      },
      "state": "SP",
      "street": "Av. Ibirapuera",
      "title": "Cobertura incrível em Moema",
      "zipCode": "04534011"
```

Analise o código proponha o plano conforme o `AGENTS.md`.