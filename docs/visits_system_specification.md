# Especificação Técnica - Sistema de Visitas

## 📋 Visão Geral

Sistema completo de agendamento e gerenciamento de visitas entre **Corretores** e **Proprietários**.

---

## 🎯 Requisitos Funcionais

### Para CORRETOR:
1. ✅ Agendar visita a um imóvel
2. ✅ Visualizar histórico de visitas solicitadas
3. ✅ Visualizar status de cada visita
4. ✅ Cancelar visita agendada (antes da aprovação)

### Para PROPRIETÁRIO:
1. ✅ Visualizar solicitações de visita recebidas
2. ✅ Aprovar ou recusar visita
3. ✅ Visualizar histórico de visitas (aprovadas/recusadas)
4. ✅ Cancelar visita aprovada (com antecedência)

---

## 📊 Modelos de Dados

### 1. Visit (Modelo Principal)

```json
{
  "id": 456,
  "listingIdentityID": 123,
  "realtorId": 5,
  "ownerId": 10,
  
  "propertyTitle": "Apartamento 3 quartos",
  "propertyAddress": "Rua Exemplo, 123",
  "propertyImageUrl": "https://...",
  
  "realtorName": "João Silva",
  "realtorPhone": "+5511999999999",
  "realtorEmail": "joao@example.com",
  
  "scheduledAt": "2025-12-30T14:00:00Z",
  "durationMinutes": 30,
  
  "status": "pending",
  "type": "withClient",
  
  "realtorNotes": "Cliente interessado em imóveis na região",
  "ownerNotes": null,
  "rejectionReason": null,
  
  "createdAt": "2025-12-29T10:00:00Z",
  "approvedAt": null,
  "rejectedAt": null,
  "cancelledAt": null,
  "updatedAt": "2025-12-29T10:00:00Z"
}
```

### 2. VisitStatus (Enum)

- `pending` - Aguardando aprovação do proprietário
- `approved` - Aprovada pelo proprietário
- `rejected` - Recusada pelo proprietário
- `cancelled` - Cancelada (por corretor ou proprietário)
- `completed` - Visita realizada
- `noShow` - Corretor não compareceu

### 3. VisitType (Enum)

- `withClient` - Visita com cliente
- `realtorOnly` - Apenas corretor (conhecer imóvel)
- `contentProduction` - Produção de conteúdo/fotos

### 4. VisitCreateRequest (Request Body)

```json
{
  "listingIdentityID": 123,
  "scheduledAt": "2025-12-30T14:00:00Z",
  "type": "withClient",
  "durationMinutes": 30,
  "realtorNotes": "Cliente interessado em imóveis na região"
}
```

### 5. VisitUpdateStatusRequest (Request Body)

```json
{
  "visitId": 456,
  "newStatus": "approved",
  "notes": "Visita aprovada. Por favor, tocar a campainha."
}
```

### 6. VisitListResponse (Response)

```json
{
  "items": [
    {
      "id": 456,
      "listingIdentityID": 123,
      "propertyTitle": "Apartamento 3 quartos",
      "propertyAddress": "Rua Exemplo, 123",
      "propertyImageUrl": "https://...",
      "scheduledAt": "2025-12-30T14:00:00Z",
      "status": "pending",
      "type": "withClient",
      "createdAt": "2025-12-29T10:00:00Z"
    }
  ],
  "pagination": {
    "currentPage": 1,
    "totalPages": 3,
    "totalItems": 45,
    "itemsPerPage": 20
  },
  "appliedFilters": {
    "status": "pending",
    "startDate": null,
    "endDate": null,
    "listingIdentityID": null
  }
}
```

### 7. VisitAvailability (Objeto)

```json
{
  "listingIdentityID": 123,
  "date": "2025-12-30",
  "availableSlots": [
    {
      "startTime": "2025-12-30T09:00:00Z",
      "endTime": "2025-12-30T09:30:00Z",
      "isAvailable": true,
      "blockedReason": null
    },
    {
      "startTime": "2025-12-30T10:00:00Z",
      "endTime": "2025-12-30T10:30:00Z",
      "isAvailable": false,
      "blockedReason": "Já existe visita agendada"
    }
  ]
}
```

---

## 🔌 Endpoints da API

### Base URL: `/api/v2/visits`

---

### 1. **Criar Visita** (CORRETOR)
```
POST /visits
```

**Request Body:**
```json
{
  "listingIdentityID": 123,
  "scheduledAt": "2025-12-30T14:00:00Z",
  "type": "withClient",
  "durationMinutes": 30,
  "realtorNotes": "Cliente interessado em imóveis na região"
}
```

**Response:** `201 Created`
```json
{
  "id": 456,
  "listingIdentityID": 123,
  "realtorId": 5,
  "ownerId": 10,
  "propertyTitle": "Apartamento 3 quartos",
  "propertyAddress": "Rua Exemplo, 123",
  "scheduledAt": "2025-12-30T14:00:00Z",
  "status": "pending",
  "type": "withClient",
  "createdAt": "2025-12-29T10:00:00Z"
}
```

**Erros:**
- `400` - Dados inválidos
- `404` - Imóvel não encontrado
- `409` - Horário já ocupado

---

### 2. **Listar Visitas do Corretor** (CORRETOR)
```
GET /visits/realtor
```

**Query Parameters:**
- `status` (opcional): `pending`, `approved`, `rejected`, `cancelled`, `completed`
- `startDate` (opcional): Data inicial (ISO 8601)
- `endDate` (opcional): Data final (ISO 8601)
- `page` (opcional): Número da página (padrão: 1)
- `limit` (opcional): Itens por página (padrão: 20)

**Response:** `200 OK`
```json
{
  "items": [
    {
      "id": 456,
      "listingIdentityID": 123,
      "propertyTitle": "Apartamento 3 quartos",
      "propertyAddress": "Rua Exemplo, 123",
      "propertyImageUrl": "https://...",
      "scheduledAt": "2025-12-30T14:00:00Z",
      "status": "pending",
      "type": "withClient",
      "createdAt": "2025-12-29T10:00:00Z"
    }
  ],
  "pagination": {
    "currentPage": 1,
    "totalPages": 3,
    "totalItems": 45,
    "itemsPerPage": 20
  }
}
```

---

### 3. **Listar Solicitações de Visita** (PROPRIETÁRIO)
```
GET /visits/owner 
```

**Query Parameters:**
- `status` (opcional): `pending`, `approved`, `rejected`, `cancelled`
- `listingIdentityID` (opcional): Filtrar por imóvel específico
- `page` (opcional): Número da página
- `limit` (opcional): Itens por página

**Response:** `200 OK`
```json
{
  "items": [
    {
      "id": 456,
      "listingIdentityID": 123,
      "propertyTitle": "Apartamento 3 quartos",
      "realtorName": "João Silva",
      "realtorPhone": "+5511999999999",
      "scheduledAt": "2025-12-30T14:00:00Z",
      "status": "pending",
      "type": "withClient",
      "realtorNotes": "Cliente interessado",
      "createdAt": "2025-12-29T10:00:00Z"
    }
  ],
  "pagination": {
    "currentPage": 1,
    "totalPages": 2,
    "totalItems": 15,
    "itemsPerPage": 20
  }
}
```

---

### 4. **Detalhes da Visita**
```
GET /visits/{visitId}
```

**Response:** `200 OK`
```json
{
  "id": 456,
  "listingIdentityID": 123,
  "realtorId": 5,
  "ownerId": 10,
  "propertyTitle": "Apartamento 3 quartos",
  "propertyAddress": "Rua Exemplo, 123",
  "propertyImageUrl": "https://...",
  "realtorName": "João Silva",
  "realtorPhone": "+5511999999999",
  "realtorEmail": "joao@example.com",
  "scheduledAt": "2025-12-30T14:00:00Z",
  "durationMinutes": 30,
  "status": "pending",
  "type": "withClient",
  "realtorNotes": "Cliente interessado",
  "ownerNotes": null,
  "rejectionReason": null,
  "createdAt": "2025-12-29T10:00:00Z",
  "updatedAt": "2025-12-29T10:00:00Z"
}
```

**Erros:**
- `404` - Visita não encontrada
- `403` - Sem permissão para visualizar

---

### 5. **Aprovar Visita** (PROPRIETÁRIO)
```
POST /visits/{visitId}/approve
```

**Request Body:**
```json
{
  "ownerNotes": "Visita aprovada. Por favor, tocar a campainha."
}
```

**Response:** `200 OK`
```json
{
  "id": 456,
  "status": "approved",
  "approvedAt": "2025-12-29T11:00:00Z",
  "ownerNotes": "Visita aprovada. Por favor, tocar a campainha."
}
```

**Erros:**
- `404` - Visita não encontrada
- `403` - Apenas proprietário pode aprovar
- `409` - Visita não está pendente

---

### 6. **Recusar Visita** (PROPRIETÁRIO)
```
POST /visits/{visitId}/reject
```

**Request Body:**
```json
{
  "rejectionReason": "Horário não disponível. Por favor, reagendar."
}
```

**Response:** `200 OK`
```json
{
  "id": 456,
  "status": "rejected",
  "rejectedAt": "2025-12-29T11:00:00Z",
  "rejectionReason": "Horário não disponível. Por favor, reagendar."
}
```

**Erros:**
- `404` - Visita não encontrada
- `403` - Apenas proprietário pode recusar
- `409` - Visita não está pendente

---

### 7. **Cancelar Visita**
```
POST /visits/{visitId}/cancel
```

**Request Body:**
```json
{
  "reason": "Imprevisto do cliente"
}
```

**Response:** `200 OK`
```json
{
  "id": 456,
  "status": "cancelled",
  "cancelledAt": "2025-12-29T12:00:00Z"
}
```

**Regras:**
- Corretor pode cancelar visitas `pending` ou `approved`
- Proprietário pode cancelar apenas visitas `approved`
- Cancelamento com menos de 2h de antecedência pode gerar penalidade

**Erros:**
- `404` - Visita não encontrada
- `403` - Sem permissão para cancelar
- `409` - Visita já foi realizada ou cancelada

---

### 8. **Verificar Disponibilidade** (CORRETOR)
```
GET /visits/availability/{listingIdentityID}
```

**Query Parameters:**
- `date`: Data desejada (formato: YYYY-MM-DD)

**Response:** `200 OK`
```json
{
  "listingIdentityID": 123,
  "date": "2025-12-30",
  "availableSlots": [
    {
      "startTime": "2025-12-30T09:00:00Z",
      "endTime": "2025-12-30T09:30:00Z",
      "isAvailable": true
    },
    {
      "startTime": "2025-12-30T10:00:00Z",
      "endTime": "2025-12-30T10:30:00Z",
      "isAvailable": false,
      "blockedReason": "Já existe visita agendada"
    }
  ]
}
```

---

### 9. **Marcar Visita como Realizada**
```
POST /visits/{visitId}/complete
```

**Request Body:**
```json
{
  "notes": "Visita realizada com sucesso. Cliente demonstrou interesse."
}
```

**Response:** `200 OK`
```json
{
  "id": 456,
  "status": "completed",
  "completedAt": "2025-12-30T14:30:00Z"
}
```

---

### 10. **Estatísticas de Visitas** (CORRETOR)
```
GET /visits/realtor/stats
```

**Response:** `200 OK`
```json
{
  "total": 45,
  "pending": 5,
  "approved": 20,
  "rejected": 8,
  "completed": 10,
  "cancelled": 2,
  "thisMonth": 12,
  "thisWeek": 3
}
```

---

### 11. **Estatísticas de Visitas** (PROPRIETÁRIO)
```
GET /visits/owner/stats
```

**Query Parameters:**
- `listingIdentityID` (opcional): Estatísticas de um imóvel específico

**Response:** `200 OK`
```json
{
  "total": 28,
  "pending": 3,
  "approved": 15,
  "rejected": 5,
  "completed": 5,
  "byListing": [
    {
      "listingIdentityID": 123,
      "propertyTitle": "Apartamento 3 quartos",
      "totalVisits": 15,
      "pending": 2,
      "approved": 8
    }
  ]
}
```

---

## 🔐 Regras de Negócio

### Agendamento:
1. ✅ Visita só pode ser agendada em horários disponíveis do proprietário
2. ✅ Mínimo de 2h de antecedência para agendar
3. ✅ Máximo de 30 dias de antecedência
4. ✅ Duração padrão: 30 minutos
5. ✅ Intervalo mínimo entre visitas: 15 minutos

### Aprovação/Recusa:
1. ✅ Apenas proprietário pode aprovar/recusar
2. ✅ Apenas visitas `pending` podem ser aprovadas/recusadas
3. ✅ Motivo de recusa é obrigatório
4. ✅ Notificação enviada ao corretor após decisão

### Cancelamento:
1. ✅ Corretor pode cancelar visitas `pending` ou `approved`
2. ✅ Proprietário pode cancelar apenas visitas `approved`
3. ✅ Cancelamento com menos de 2h gera penalidade
4. ✅ Notificação enviada à outra parte

### Histórico:
1. ✅ Corretor vê todas as suas visitas solicitadas
2. ✅ Proprietário vê visitas de seus imóveis
3. ✅ Filtros por status, data, imóvel

---

## 📱 Fluxos de Tela

### Fluxo do Corretor:

```
1. Lista de Imóveis
   ↓
2. Detalhes do Imóvel
   ↓ [Agendar Visita]
3. Verificar Disponibilidade
   ↓ [Selecionar Data/Hora]
4. Confirmar Agendamento
   ↓
5. Status da Visita (pending)
   ↓ [Proprietário aprova]
6. Visita Confirmada (approved)
   ↓ [Após visita]
7. Marcar como Realizada (completed)
```

### Fluxo do Proprietário:

```
1. Notificação de Nova Solicitação
   ↓
2. Lista de Solicitações Pendentes
   ↓ [Selecionar visita]
3. Detalhes da Solicitação
   ↓ [Aprovar ou Recusar]
4. Visita Aprovada/Recusada
   ↓
5. Histórico de Visitas
```

---

## 🔔 Notificações

### Para Corretor:
- ✅ Visita aprovada pelo proprietário
- ✅ Visita recusada pelo proprietário
- ✅ Visita cancelada pelo proprietário
- ✅ Lembrete 1h antes da visita

### Para Proprietário:
- ✅ Nova solicitação de visita
- ✅ Visita cancelada pelo corretor
- ✅ Lembrete 1h antes da visita aprovada

---

## 🗄️ Estrutura de Banco de Dados (Sugestão)

```sql
CREATE TABLE visits (
  id SERIAL PRIMARY KEY,
  listing_id INTEGER NOT NULL REFERENCES listings(id),
  realtor_id INTEGER NOT NULL REFERENCES users(id),
  owner_id INTEGER NOT NULL REFERENCES users(id),
  
  scheduled_at TIMESTAMP WITH TIME ZONE NOT NULL,
  duration_minutes INTEGER DEFAULT 30,
  
  status VARCHAR(20) NOT NULL DEFAULT 'pending',
  type VARCHAR(30) NOT NULL,
  
  realtor_notes TEXT,
  owner_notes TEXT,
  rejection_reason TEXT,
  
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  approved_at TIMESTAMP WITH TIME ZONE,
  rejected_at TIMESTAMP WITH TIME ZONE,
  cancelled_at TIMESTAMP WITH TIME ZONE,
  completed_at TIMESTAMP WITH TIME ZONE,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  
  CONSTRAINT valid_status CHECK (status IN ('pending', 'approved', 'rejected', 'cancelled', 'completed', 'noShow')),
  CONSTRAINT valid_type CHECK (type IN ('withClient', 'realtorOnly', 'contentProduction'))
);

CREATE INDEX idx_visits_realtor ON visits(realtor_id, status);
CREATE INDEX idx_visits_owner ON visits(owner_id, status);
CREATE INDEX idx_visits_listing ON visits(listing_id, scheduled_at);
CREATE INDEX idx_visits_scheduled ON visits(scheduled_at) WHERE status IN ('pending', 'approved');
```

---

## ✅ Checklist de Implementação

### Backend:
- [ ] Criar tabela `visits`
- [ ] Implementar endpoints CRUD
- [ ] Validações de horário e disponibilidade
- [ ] Sistema de notificações
- [ ] Testes unitários e integração

### Frontend (Flutter):
- [ ] Atualizar `Visit` model com novos campos
- [ ] Criar DTOs de request/response
- [ ] Implementar repository de visitas
- [ ] Criar notifiers/controllers
- [ ] Telas de listagem (corretor e proprietário)
- [ ] Tela de agendamento
- [ ] Tela de detalhes da visita
- [ ] Sistema de notificações push

### Adicionar em `api_paths.dart`:
```dart
// Visits
static const String visitsCreate = '/visits';
static const String visitsRealtor = '/visits/realtor';
static const String visitsOwner = '/visits/owner';
static const String visitsDetail = '/visits/detail';
static const String visitsApprove = '/visits/approve';
static const String visitsReject = '/visits/reject';
static const String visitsCancel = '/visits/cancel';
static const String visitsComplete = '/visits/complete';
static const String visitsAvailability = '/visits/availability';
static const String visitsRealtorStats = '/visits/realtor/stats';
static const String visitsOwnerStats = '/visits/owner/stats';
```

---

## 📝 Notas Finais

- Todos os timestamps devem usar **ISO 8601** com timezone
- Paginação padrão: 20 itens por página
- Autenticação via JWT obrigatória em todos os endpoints
- Rate limiting: 100 requisições por minuto por usuário
- Logs de auditoria para todas as ações (criar, aprovar, recusar, cancelar)
