# Gerenciamento de Timezones

Este documento descreve as convencoes de fuso horario adotadas pelo `toq_server`, os campos expostos via API e exemplos de chamadas que demonstram como clientes devem interagir com os novos contratos.

## Convencoes Gerais

- **Padrao de armazenamento**: toda data/hora persistida permanece em UTC (`time.Time.UTC()`).
- **Identificadores**: sempre utilize IDs IANA validos (ex.: `America/Sao_Paulo`).
- **Padrao de entrada**: quando o campo `timezone` e obrigatorio, a API rejeita valores invalidos (`400 INVALID_TIMEZONE`).
- **Conversao automatica**: respostas normalizam datas para o fuso solicitado; ausente o parametro, aplicamos o padrao `America/Sao_Paulo`.
- **Formatacao**: todos os exemplos usam RFC3339 completo, incluindo offset quando aplicavel.

## Calendarios de feriado (Admin)

### Criar calendario

`POST /admin/holidays/calendars`

```bash
curl -X POST "https://api.toq.dev/admin/holidays/calendars" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Feriados Sao Paulo",
    "scope": "STATE",
    "state": "SP",
    "cityIbge": "",
    "isActive": true,
    "timezone": "America/Sao_Paulo"
  }'
```

Resposta (`201 Created`):

```json
{
  "id": 42,
  "name": "Feriados Sao Paulo",
  "scope": "STATE",
  "state": "SP",
  "cityIbge": "",
  "isActive": true,
  "timezone": "America/Sao_Paulo"
}
```

### Atualizar calendario

`PUT /admin/holidays/calendars`

O payload e identico ao de criacao, incluindo `timezone`. Alterar o fuso atualiza a referencia usada para normalizar novas datas.

### Listar datas do calendario

`GET /admin/holidays/dates?calendarId=42&from=2025-12-01T00:00:00Z&to=2026-01-10T00:00:00Z&timezone=America/Sao_Paulo`

- `timezone` (opcional) controla o fuso das datas retornadas.

Exemplo de resposta:

```json
{
  "dates": [
    {
      "id": 101,
      "calendarId": 42,
      "holidayDate": "2025-12-25T00:00:00-03:00",
      "label": "Natal",
      "recurrent": true,
      "timezone": "America/Sao_Paulo"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 1
  }
}
```

## Agenda de listings (Owners)

### Criar bloqueio permanente ou temporario

`POST /schedules/listing/block`

```bash
curl -X POST "https://api.toq.dev/schedules/listing/block" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "listingId": 3241,
    "entryType": "BLOCK",
    "startsAt": "2025-06-15T09:00:00-03:00",
    "endsAt": "2025-06-15T11:00:00-03:00",
    "reason": "Janela de manutencao",
    "timezone": "America/Sao_Paulo"
  }'
```

Resposta (`200 OK`):

```json
{
  "entry": {
    "id": 5568,
    "entryType": "BLOCK",
    "startsAt": "2025-06-15T09:00:00-03:00",
    "endsAt": "2025-06-15T11:00:00-03:00",
    "blocking": true,
    "reason": "Janela de manutencao",
    "timezone": "America/Sao_Paulo"
  }
}
```

> O servico normaliza internamente para UTC e armazena a referencia de fuso. Requisicoes futuras para o mesmo listing respeitarao o timezone informado.

### Atualizar bloqueio existente

`PUT /schedules/listing/block`

Inclua `timezone` no payload para garantir que o intervalo seja interpretado corretamente.

### Listar agenda do listing

`GET /schedules/listing/agenda?listingId=3241&rangeFrom=2025-06-01T00:00:00Z&rangeTo=2025-06-30T23:59:59Z&timezone=America/Sao_Paulo`

A resposta devolvera cada entrada convertida para o fuso informado, com o campo `timezone` anexado.

### Listar disponibilidade

`GET /schedules/listing/availability?listingId=3241&rangeFrom=2025-06-15T00:00:00Z&rangeTo=2025-06-22T23:59:59Z&timezone=America/New_York`

Resposta (trecho):

```json
{
  "slots": [
    {
      "startsAt": "2025-06-16T08:00:00-04:00",
      "endsAt": "2025-06-16T09:00:00-04:00"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 12
  },
  "timezone": "America/New_York"
}
```

## Agenda do fotografo

### Criar time-off

`POST /photographer/agenda/time-off`

```bash
curl -X POST "https://api.toq.dev/photographer/agenda/time-off" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "startDate": "2025-07-05T09:00:00-03:00",
    "endDate": "2025-07-05T18:00:00-03:00",
    "reason": "Evento presencial",
    "timezone": "America/Sao_Paulo",
    "holidayCalendarId": 42,
    "horizonMonths": 3,
    "workdayStartHour": 8,
    "workdayEndHour": 19
  }'
```

Resposta (`201 Created`):

```json
{
  "message": "Time-off created successfully",
  "timeOffId": 901
}
```

### Listar agenda consolidada

`GET /photographer/agenda?startDate=2025-07-01T00:00:00Z&endDate=2025-07-31T23:59:59Z&timezone=Europe/Lisbon&holidayCalendarIds=42`

Trecho da resposta:

```json
{
  "slots": [
    {
      "slotId": 7788,
      "photographerId": 456,
      "start": "2025-07-10T13:00:00+01:00",
      "end": "2025-07-10T14:00:00+01:00",
      "status": "AVAILABLE",
      "groupId": "slot-2025-07-10-morning",
      "source": "SLOT",
      "isHoliday": false,
      "isTimeOff": false,
      "timezone": "Europe/Lisbon"
    },
    {
      "photographerId": 456,
      "start": "2025-07-15T00:00:00+01:00",
      "end": "2025-07-16T00:00:00+01:00",
      "status": "BLOCKED",
      "groupId": "holiday-2025-07-15",
      "source": "HOLIDAY",
      "isHoliday": true,
      "holidayLabels": ["Feriado Municipal"],
      "holidayCalendarIds": [42],
      "timezone": "Europe/Lisbon"
    }
  ],
  "total": 38,
  "page": 1,
  "size": 20,
  "timezone": "Europe/Lisbon"
}
```

### Remover time-off

`DELETE /photographer/agenda/time-off`

O payload deve repetir `timezone`, `horizonMonths`, `workdayStartHour` e `workdayEndHour` para que a projecao de slots seja recalculada no fuso correto apos a exclusao.

## Boas praticas para clientes

- Sempre converta a data local para RFC3339 com offset antes de enviar para a API.
- No consumo de respostas, utilize o campo `timezone` para ajustar exibicoes ou para normalizar novamente a UTC.
- Evite confiar em offsets fixos (ex.: `-03:00`); prefira o ID IANA pois contempla horario de verao.
- Ao migrar integracoes legadas, valide que qualquer cache ou camada de persistencia externa tambem armazene o timezone associado.

## Proximos passos internos

- Atualizar os exemplos no `swagger.yaml` para refletir os novos campos de timezone.
- Revisar scripts SQL (`scripts/atz_permissions.sql`, etc.) caso sejam necessarios dados de seed com fusos especificos.
- Manter testes cobrindo fluxos com parametros de timezone distintos (ex.: `America/New_York`, `Europe/Lisbon`).
