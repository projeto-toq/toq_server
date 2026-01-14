O endpoint `GET /proposals/realtor` tem a respsota:
```json
{
  "items": [
    {
      "acceptedAt": "string",
      "cancelledAt": "string",
      "createdAt": "string",
      "documents": [
        {
          "base64Payload": "string",
          "fileName": "string",
          "fileSizeBytes": 0,
          "id": 0,
          "mimeType": "string"
        }
      ],
      "documentsCount": 0,
      "id": 0,
      "listing": {
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
      },
      "listingIdentityId": 0,
      "proposalText": "string",
      "realtor": {
        "acceptedProposals": 0,
        "accountAgeMonths": 0,
        "name": "string",
        "nickname": "string",
        "photoUrl": "string",
        "proposalsCreated": 0
      },
      "receivedAt": "string",
      "rejectedAt": "string",
      "rejectionReason": "string",
      "respondedAt": "string",
      "status": "string"
    }
  ],
  "total": 0
}
```
entretanto, os dados do proprietário estâo ausentes:
```json
  "owner": {
    "fullName": "string",
    "id": 0,
    "memberSinceMonths": 0,
    "photoUrl": "string",
    "proposalAverageSeconds": 0,
    "visitAverageSeconds": 0
  },
```
Já o endpoint `POST /proposals/detail` que tem como resposta:
```json
{
  "documents": [
    {
      "base64Payload": "string",
      "fileName": "string",
      "fileSizeBytes": 0,
      "id": 0,
      "mimeType": "string"
    }
  ],
  "proposal": {
    "acceptedAt": "string",
    "cancelledAt": "string",
    "createdAt": "string",
    "documents": [
      {
        "base64Payload": "string",
        "fileName": "string",
        "fileSizeBytes": 0,
        "id": 0,
        "mimeType": "string"
      }
    ],
    "documentsCount": 0,
    "id": 0,
    "listing": {
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
    },
    "listingIdentityId": 0,
    "proposalText": "string",
    "realtor": {
      "acceptedProposals": 0,
      "accountAgeMonths": 0,
      "name": "string",
      "nickname": "string",
      "photoUrl": "string",
      "proposalsCreated": 0
    },
    "receivedAt": "string",
    "rejectedAt": "string",
    "rejectionReason": "string",
    "respondedAt": "string",
    "status": "string"
  },
  "realtor": {
    "acceptedProposals": 0,
    "accountAgeMonths": 0,
    "name": "string",
    "nickname": "string",
    "photoUrl": "string",
    "proposalsCreated": 0
  }
}```

tem duplicado os dados do corretor, quando deveria ter os dados do proprietário ao invés do último conjunto "realtor", como no exemplo abaixo:
```json
  "owner": {
    "fullName": "string",
    "id": 0,
    "memberSinceMonths": 0,
    "photoUrl": "string",
    "proposalAverageSeconds": 0,
    "visitAverageSeconds": 0
  },
```
Portanto busque todas as infromações que precisa, para ter certeza da causa raiz, e só então proponha o plano conforme o `AGENTS.md`.