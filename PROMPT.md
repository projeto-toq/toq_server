O endpoint `GET /visits/owner` tem a respsota:
```json
{
  "items": [
    {
      "firstOwnerActionAt": "2025-01-10T14:05:00Z",
      "id": 456,
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
      "listingIdentityId": 123,
      "listingVersion": 1,
      "liveStatus": "AO_VIVO",
      "notes": "string",
      "owner": {
        "avgResponseHours": 4.5,
        "fullName": "Maria Souza",
        "memberSince": "2021-05-10T12:00:00Z",
        "memberSinceDays": 980,
        "photoUrl": "https://signed.cdn/photos/owner.jpg",
        "userId": 10
      },
      "ownerUserId": 10,
      "realtor": {
        "fullName": "João Corretor",
        "memberSince": "2022-02-01T09:30:00Z",
        "memberSinceDays": 600,
        "photoUrl": "https://signed.cdn/photos/realtor.jpg",
        "userId": 5,
        "visitsPerformed": 37
      },
      "rejectionReason": "string",
      "requesterUserId": 5,
      "scheduledEnd": "2025-01-10T14:30:00Z",
      "scheduledStart": "2025-01-10T14:00:00Z",
      "source": "APP",
      "status": "PENDING",
      "timeline": {
        "createdAt": "2025-01-05T12:00:00Z",
        "receivedAt": "2025-01-05T12:05:00Z",
        "respondedAt": "2025-01-05T13:15:00Z"
      }
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
Porem os dados de
```json
        "propertyType": {
          "code": 0,
          "label": "string",
          "propertyBit": 0
        },
```
estão audentes.

Portanto busque todas as infromações que precisa, para ter certeza da causa raiz, e só então proponha o plano conforme o `AGENTS.md`.

Em /proposal/realtor

- dados do proprietário ausente, tendo apenas dados do corretor.
- proprietário "recebeu proposta" está ausente.