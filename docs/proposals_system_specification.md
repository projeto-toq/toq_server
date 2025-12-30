# Especificação Técnica - Sistema de Propostas

## 📋 Visão Geral

Sistema completo de envio e gerenciamento de propostas de compra/locação entre **Corretores** (em nome de clientes) e **Proprietários**.

---

## 🎯 Requisitos Funcionais

### Para CORRETOR:
1. Enviar proposta para um imóvel
2. Visualizar histórico de propostas enviadas
3. Visualizar status de cada proposta
4. Editar proposta (apenas se status = `pending`)
5. Cancelar proposta (antes da aceitação)

### Para PROPRIETÁRIO:
1. Visualizar propostas recebidas
2. Aceitar proposta
3. Recusar proposta (com motivo)
4. Visualizar histórico de propostas

---

## Modelos de Dados

### 1. Proposal (Modelo Principal)

```json
{
  "id": 789,
  "listingId": 123,
  "realtorId": 5,
  "ownerId": 10,
  
  "propertyTitle": "Apartamento 3 quartos",
  "propertyAddress": "Rua Exemplo, 123",
  "propertyImageUrl": "https://...",
  "transactionType": "sale",
  
  "realtorName": "João Silva",
  "realtorPhone": "+5511999999999",
  "realtorEmail": "joao@example.com",
  "clientName": "Maria Silva",
  "clientPhone": "+5511988887777",
  
  "proposedValue": 450000.00,
  "originalValue": 500000.00,
  "paymentMethod": "cashAndFinancing",
  "financingDetails": {
    "bankName": "Banco do Brasil",
    "approvedAmount": 350000.00,
    "preApproved": true,
    "approvalDocument": "https://..."
  },
  
  "downPayment": 100000.00,
  "installments": null,
  "acceptsExchange": false,
  "exchangeDetails": null,
  
  "rentalMonths": null,
  "guaranteeType": null,
  "securityDeposit": null,
  
  "status": "pending",
  
  "proposalNotes": "Cliente tem interesse imediato",
  "ownerNotes": null,
  "rejectionReason": null,
  
  "documents": [
    {
      "id": 101,
      "fileName": "pre_aprovacao.pdf",
      "fileType": "application/pdf",
      "fileUrl": "https://...",
      "fileSizeBytes": 245678,
      "uploadedAt": "2025-12-29T14:30:00Z"
    }
  ],
  
  "createdAt": "2025-12-29T15:00:00Z",
  "acceptedAt": null,
  "rejectedAt": null,
  "cancelledAt": null,
  "expiresAt": "2026-01-15T23:59:59Z",
  "updatedAt": "2025-12-29T15:00:00Z"
}
```

### 2. ProposalStatus (Enum)

- `pending` - Aguardando resposta do proprietário
- `accepted` - Aceita pelo proprietário
- `rejected` - Recusada pelo proprietário
- `cancelled` - Cancelada pelo corretor
- `expired` - Proposta expirou


### 3. TransactionType (Enum)

- `sale` - Venda
- `rent` - Locação

### 4. PaymentMethod (Enum)

- `cash` - À vista
- `financing` - Financiamento bancário
- `installments` - Parcelado direto com proprietário
- `cashAndFinancing` - Entrada + financiamento
- `exchange` - Permuta
- `cashAndExchange` - À vista + permuta

### 5. GuaranteeType (Enum)

- `securityDeposit` - Caução
- `suretyBond` - Fiança
- `rentalInsurance` - Seguro fiança
- `guarantor` - Fiador

### 6. FinancingDetails (Objeto)

```json
{
  "bankName": "Banco do Brasil",
  "approvedAmount": 350000.00,
  "preApproved": true,
  "approvalDocument": "https://..."
}
```

### 7. ExchangeDetails (Objeto)

```json
{
  "propertyDescription": "Casa 2 quartos",
  "propertyAddress": "Rua das Flores, 456",
  "estimatedValue": 300000.00,
  "propertyType": "house",
  "propertyImages": [
    "https://...",
    "https://..."
  ]
}
```

### 8. ProposalDocument (Objeto)

```json
{
  "id": 101,
  "fileName": "pre_aprovacao.pdf",
  "fileType": "application/pdf",
  "fileUrl": "https://...",
  "fileSizeBytes": 245678,
  "uploadedAt": "2025-12-29T14:30:00Z"
}
```

### 9. ProposalCreateRequest (Request Body)

```json
{
  "listingId": 123,
  "transactionType": "sale",
  "proposedValue": 450000.00,
  "paymentMethod": "cashAndFinancing",
  "financingDetails": {
    "bankName": "Banco do Brasil",
    "approvedAmount": 350000.00,
    "preApproved": true
  },
  "downPayment": 100000.00,
  "installments": null,
  "acceptsExchange": false,
  "exchangeDetails": null,
  "rentalMonths": null,
  "guaranteeType": null,
  "securityDeposit": null,
  "clientName": "Maria Silva",
  "clientPhone": "+5511988887777",
  "proposalNotes": "Cliente tem interesse imediato. Financiamento pré-aprovado.",
  "expiresAt": "2026-01-15T23:59:59Z",
  "documentIds": [101, 102]
}
```

### 10. ProposalUpdateStatusRequest (Request Body)

```json
{
  "proposalId": 789,
  "newStatus": "rejected",
  "notes": "Valor abaixo do esperado",
  "rejectionReason": "Valor mínimo aceitável: R$ 480.000"
}
```

### 11. ProposalListResponse (Response)

```json
{
  "items": [
    {
      "id": 789,
      "listingId": 123,
      "propertyTitle": "Apartamento 3 quartos",
      "proposedValue": 450000.00,
      "status": "pending",
      "createdAt": "2025-12-29T15:00:00Z"
    }
  ],
  "pagination": {
    "currentPage": 1,
    "totalPages": 2,
    "totalItems": 25,
    "itemsPerPage": 20
  },
  "appliedFilters": {
    "status": "pending",
    "transactionType": null,
    "startDate": null,
    "endDate": null,
    "listingId": null,
    "minValue": null,
    "maxValue": null
  },
  "summary": {
    "totalProposals": 25,
    "pendingCount": 8,
    "acceptedCount": 10,
    "rejectedCount": 5,
    "highestProposal": 800000.00,
    "lowestProposal": 420000.00,
    "averageProposal": 520000.00
  }
}
```

---

## 🔌 Endpoints da API

### Base URL: `/api/v2/proposals`

---

### 1. **Criar Proposta** (CORRETOR)
```
POST /proposals
```

**Request Body:**
```json
{
  "listingId": 123,
  "transactionType": "sale",
  "proposedValue": 450000.00,
  "paymentMethod": "cashAndFinancing",
  "financingDetails": {
    "bankName": "Banco do Brasil",
    "approvedAmount": 350000.00,
    "preApproved": true
  },
  "downPayment": 100000.00,
  "clientName": "Maria Silva",
  "clientPhone": "+5511988887777",
  "proposalNotes": "Cliente tem interesse imediato. Financiamento pré-aprovado.",
  "expiresAt": "2026-01-15T23:59:59Z",
  "documentIds": [101, 102]
}
```

**Response:** `201 Created`
```json
{
  "id": 789,
  "listingId": 123,
  "realtorId": 5,
  "ownerId": 10,
  "propertyTitle": "Apartamento 3 quartos",
  "transactionType": "sale",
  "proposedValue": 450000.00,
  "originalValue": 500000.00,
  "paymentMethod": "cashAndFinancing",
  "status": "pending",
  "origin": "original",
  "createdAt": "2025-12-29T15:00:00Z",
  "expiresAt": "2026-01-15T23:59:59Z"
}
```

**Erros:**
- `400` - Dados inválidos
- `404` - Imóvel não encontrado
- `409` - Já existe proposta ativa para este imóvel

---

### 2. **Listar Propostas do Corretor** (CORRETOR)
```
GET /proposals/realtor
```

**Query Parameters:**
- `status` (opcional): `pending`, `accepted`, `rejected`, `cancelled`, `expired`
- `transactionType` (opcional): `sale`, `rent`
- `listingId` (opcional): Filtrar por imóvel
- `startDate` (opcional): Data inicial
- `endDate` (opcional): Data final
- `page` (opcional): Número da página (padrão: 1)
- `limit` (opcional): Itens por página (padrão: 20)

**Response:** `200 OK`
```json
{
  "items": [
    {
      "id": 789,
      "listingId": 123,
      "propertyTitle": "Apartamento 3 quartos",
      "propertyAddress": "Rua Exemplo, 123",
      "propertyImageUrl": "https://...",
      "transactionType": "sale",
      "proposedValue": 450000.00,
      "originalValue": 500000.00,
      "status": "pending",
      "clientName": "Maria Silva",
      "createdAt": "2025-12-29T15:00:00Z",
      "expiresAt": "2026-01-15T23:59:59Z"
    }
  ],
  "pagination": {
    "currentPage": 1,
    "totalPages": 2,
    "totalItems": 25,
    "itemsPerPage": 20
  },
  "summary": {
    "totalProposals": 25,
    "pendingCount": 8,
    "acceptedCount": 10,
    "rejectedCount": 5,
    "highestProposal": 800000.00,
    "averageProposal": 520000.00
  }
}
```

---

### 3. **Listar Propostas Recebidas** (PROPRIETÁRIO)
```
GET /proposals/owner
```

**Query Parameters:**
- `status` (opcional): Filtrar por status
- `listingId` (opcional): Filtrar por imóvel específico
- `sortBy` (opcional): `value_desc`, `value_asc`, `date_desc`, `date_asc`
- `page` (opcional): Número da página
- `limit` (opcional): Itens por página

**Response:** `200 OK`
```json
{
  "items": [
    {
      "id": 789,
      "listingId": 123,
      "propertyTitle": "Apartamento 3 quartos",
      "realtorName": "João Silva",
      "realtorPhone": "+5511999999999",
      "clientName": "Maria Silva",
      "transactionType": "sale",
      "proposedValue": 450000.00,
      "originalValue": 500000.00,
      "paymentMethod": "cashAndFinancing",
      "status": "pending",
      "proposalNotes": "Cliente tem interesse imediato",
      "createdAt": "2025-12-29T15:00:00Z",
      "expiresAt": "2026-01-15T23:59:59Z"
    }
  ],
  "pagination": {
    "currentPage": 1,
    "totalPages": 1,
    "totalItems": 12,
    "itemsPerPage": 20
  },
  "summary": {
    "totalProposals": 12,
    "pendingCount": 4,
    "acceptedCount": 5,
    "rejectedCount": 3,
    "highestProposal": 480000.00,
    "lowestProposal": 420000.00,
    "averageProposal": 445000.00
  }
}
```

---

### 4. **Detalhes da Proposta**
```
GET /proposals/{proposalId}
```

**Response:** `200 OK`
```json
{
  "id": 789,
  "listingId": 123,
  "realtorId": 5,
  "ownerId": 10,
  "propertyTitle": "Apartamento 3 quartos",
  "propertyAddress": "Rua Exemplo, 123",
  "propertyImageUrl": "https://...",
  "transactionType": "sale",
  "realtorName": "João Silva",
  "realtorPhone": "+5511999999999",
  "clientName": "Maria Silva",
  "clientPhone": "+5511988887777",
  "proposedValue": 450000.00,
  "originalValue": 500000.00,
  "paymentMethod": "cashAndFinancing",
  "financingDetails": {
    "bankName": "Banco do Brasil",
    "approvedAmount": 350000.00,
    "preApproved": true
  },
  "downPayment": 100000.00,
  "status": "pending",
  "proposalNotes": "Cliente tem interesse imediato",
  "documents": [
    {
      "id": 101,
      "fileName": "pre_aprovacao.pdf",
      "fileType": "application/pdf",
      "fileUrl": "https://...",
      "fileSizeBytes": 245678,
      "uploadedAt": "2025-12-29T14:30:00Z"
    }
  ],
  "createdAt": "2025-12-29T15:00:00Z",
  "expiresAt": "2026-01-15T23:59:59Z",
  "updatedAt": "2025-12-29T15:00:00Z"
}
```

**Erros:**
- `404` - Proposta não encontrada
- `403` - Sem permissão para visualizar

---

### 5. **Aceitar Proposta** (PROPRIETÁRIO)
```
POST /proposals/{proposalId}/accept
```

**Request Body:**
```json
{
  "ownerNotes": "Proposta aceita. Aguardando próximos passos."
}
```

**Response:** `200 OK`
```json
{
  "id": 789,
  "status": "accepted",
  "acceptedAt": "2025-12-29T16:00:00Z",
  "ownerNotes": "Proposta aceita. Aguardando próximos passos."
}
```

**Erros:**
- `404` - Proposta não encontrada
- `403` - Apenas proprietário pode aceitar
- `409` - Proposta não está pendente ou expirou

---

### 6. **Recusar Proposta** (PROPRIETÁRIO)
```
POST /proposals/{proposalId}/reject
```

**Request Body:**
```json
{
  "rejectionReason": "Valor abaixo do esperado. Valor mínimo aceitável: R$ 480.000"
}
```

**Response:** `200 OK`
```json
{
  "id": 789,
  "status": "rejected",
  "rejectedAt": "2025-12-29T16:00:00Z",
  "rejectionReason": "Valor abaixo do esperado. Valor mínimo aceitável: R$ 480.000"
}
```

**Erros:**
- `404` - Proposta não encontrada
- `403` - Apenas proprietário pode recusar
- `409` - Proposta não está pendente

---

### 7. **Cancelar Proposta** (CORRETOR)
```
POST /proposals/{proposalId}/cancel
```

**Request Body:**
```json
{
  "reason": "Cliente desistiu da compra"
}
```

**Response:** `200 OK`
```json
{
  "id": 789,
  "status": "cancelled",
  "cancelledAt": "2025-12-29T17:00:00Z"
}
```

**Regras:**
- Corretor pode cancelar propostas `pending`
- Não pode cancelar propostas `accepted`

**Erros:**
- `404` - Proposta não encontrada
- `403` - Sem permissão para cancelar
- `409` - Proposta já foi aceita

---

### 8. **Editar Proposta** (CORRETOR)
```
PUT /proposals/{proposalId}
```

**Request Body:**
```json
{
  "proposedValue": 460000.00,
  "downPayment": 110000.00,
  "proposalNotes": "Cliente aumentou a oferta",
  "expiresAt": "2026-01-20T23:59:59Z"
}
```

**Response:** `200 OK`
```json
{
  "id": 789,
  "proposedValue": 460000.00,
  "downPayment": 110000.00,
  "proposalNotes": "Cliente aumentou a oferta",
  "updatedAt": "2025-12-29T17:30:00Z"
}
```

**Regras:**
- Apenas propostas `pending` podem ser editadas
- Não pode alterar `listingId`, `transactionType`

---

### 9. **Upload de Documento** (CORRETOR)
```
POST /proposals/documents/upload-url
```

**Request Body:**
```json
{
  "fileName": "pre_aprovacao.pdf",
  "fileType": "application/pdf",
  "fileSizeBytes": 245678
}
```

**Response:** `200 OK`
```json
{
  "documentId": 101,
  "uploadUrl": "https://s3.amazonaws.com/...",
  "headers": {
    "Content-Type": "application/pdf"
  }
}
```

---

### 10. **Estatísticas de Propostas** (CORRETOR)
```
GET /proposals/realtor/stats
```

**Query Parameters:**
- `startDate` (opcional): Data inicial
- `endDate` (opcional): Data final

**Response:** `200 OK`
```json
{
  "total": 45,
  "pending": 8,
  "accepted": 20,
  "rejected": 12,
  "cancelled": 5,
  "thisMonth": 15,
  "thisWeek": 4,
  "acceptanceRate": 44.4,
  "averageResponseTime": "2.5 days",
  "byTransactionType": {
    "sale": 30,
    "rent": 15
  }
}
```

---

### 11. **Estatísticas de Propostas** (PROPRIETÁRIO)
```
GET /proposals/owner/stats
```

**Query Parameters:**
- `listingId` (opcional): Estatísticas de um imóvel específico

**Response:** `200 OK`
```json
{
  "total": 28,
  "pending": 4,
  "accepted": 15,
  "rejected": 9,
  "byListing": [
    {
      "listingId": 123,
      "propertyTitle": "Apartamento 3 quartos",
      "totalProposals": 12,
      "pending": 2,
      "accepted": 7,
      "highestProposal": 480000.00,
      "averageProposal": 445000.00
    }
  ],
  "averageResponseTime": "1.8 days"
}
```

---

---

## 🔐 Regras de Negócio

### Criação de Proposta:
1. ✅ Corretor só pode criar proposta para imóveis publicados
2. ✅ Valor proposto deve ser > 0
3. ✅ Data de expiração máxima: 90 dias
4. ✅ Documentos obrigatórios para financiamento
5. ✅ Validar dados de permuta se `acceptsExchange = true`

### Aceitação/Recusa:
1. ✅ Apenas proprietário pode aceitar/recusar
2. ✅ Apenas propostas `pending` podem ser aceitas/recusadas
3. ✅ Motivo de recusa é obrigatório
4. ✅ Ao aceitar, todas as outras propostas do mesmo imóvel são automaticamente recusadas
5. ✅ Notificação enviada ao corretor

### Cancelamento:
1. ✅ Corretor pode cancelar propostas `pending`
2. ✅ Não pode cancelar propostas `accepted`
3. ✅ Notificação enviada ao proprietário

### Edição:
1. ✅ Apenas propostas `pending` podem ser editadas
2. ✅ Não pode alterar tipo de transação ou imóvel
3. ✅ Histórico de edições mantido

### Expiração:
1. ✅ Propostas expiradas automaticamente mudam para status `expired`
2. ✅ Job automático verifica expiração a cada hora
3. ✅ Notificação enviada 24h antes da expiração

---

## 📱 Fluxos de Tela

### Fluxo do Corretor:

```
1. Lista de Imóveis / Detalhes do Imóvel
   ↓ [Fazer Proposta]
2. Formulário de Proposta
   ↓ [Preencher dados]
3. Upload de Documentos (opcional)
   ↓ [Confirmar]
4. Proposta Enviada (status: pending)
   ↓ [Aguardar resposta]
5a. Proposta Aceita (status: accepted) → Próximos passos
5b. Proposta Recusada (status: rejected) → Nova proposta
```

### Fluxo do Proprietário:

```
1. Notificação de Nova Proposta
   ↓
2. Lista de Propostas Recebidas
   ↓ [Selecionar proposta]
3. Detalhes da Proposta
   ↓ [Decidir]
4a. Aceitar Proposta → Próximos passos
4b. Recusar Proposta → Fim
```

---

## 🔔 Notificações

### Para Corretor:
- ✅ Proposta aceita pelo proprietário
- ✅ Proposta recusada pelo proprietário
- ✅ Proposta expirando em 24h
- ✅ Proposta expirou

### Para Proprietário:
- ✅ Nova proposta recebida
- ✅ Proposta editada pelo corretor
- ✅ Proposta cancelada pelo corretor

---

## 🗄️ Estrutura de Banco de Dados (Sugestão)

```sql
CREATE TABLE proposals (
  id SERIAL PRIMARY KEY,
  listing_id INTEGER NOT NULL REFERENCES listings(id),
  realtor_id INTEGER NOT NULL REFERENCES users(id),
  owner_id INTEGER NOT NULL REFERENCES users(id),
  
  transaction_type VARCHAR(10) NOT NULL,
  
  proposed_value DECIMAL(15,2) NOT NULL,
  original_value DECIMAL(15,2),
  payment_method VARCHAR(30) NOT NULL,
  
  -- Financiamento
  financing_bank_name VARCHAR(100),
  financing_approved_amount DECIMAL(15,2),
  financing_pre_approved BOOLEAN DEFAULT false,
  
  -- Venda
  down_payment DECIMAL(15,2),
  installments INTEGER,
  accepts_exchange BOOLEAN DEFAULT false,
  exchange_property_description TEXT,
  exchange_property_address TEXT,
  exchange_estimated_value DECIMAL(15,2),
  
  -- Locação
  rental_months INTEGER,
  guarantee_type VARCHAR(30),
  security_deposit DECIMAL(15,2),
  
  -- Cliente
  client_name VARCHAR(200),
  client_phone VARCHAR(20),
  
  -- Status
  status VARCHAR(20) NOT NULL DEFAULT 'pending',
  
  -- Observações
  proposal_notes TEXT,
  owner_notes TEXT,
  rejection_reason TEXT,
  
  -- Timestamps
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  accepted_at TIMESTAMP WITH TIME ZONE,
  rejected_at TIMESTAMP WITH TIME ZONE,
  cancelled_at TIMESTAMP WITH TIME ZONE,
  expires_at TIMESTAMP WITH TIME ZONE,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  
  CONSTRAINT valid_status CHECK (status IN ('pending', 'accepted', 'rejected', 'cancelled', 'expired')),
  CONSTRAINT valid_transaction_type CHECK (transaction_type IN ('sale', 'rent')),
  CONSTRAINT positive_value CHECK (proposed_value > 0)
);

CREATE TABLE proposal_documents (
  id SERIAL PRIMARY KEY,
  proposal_id INTEGER NOT NULL REFERENCES proposals(id) ON DELETE CASCADE,
  file_name VARCHAR(255) NOT NULL,
  file_type VARCHAR(100) NOT NULL,
  file_url TEXT NOT NULL,
  file_size_bytes INTEGER NOT NULL,
  uploaded_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_proposals_realtor ON proposals(realtor_id, status);
CREATE INDEX idx_proposals_owner ON proposals(owner_id, status);
CREATE INDEX idx_proposals_listing ON proposals(listing_id, status);
CREATE INDEX idx_proposals_expires ON proposals(expires_at) WHERE status = 'pending';
CREATE INDEX idx_proposal_documents ON proposal_documents(proposal_id);
```

---

## ✅ Checklist de Implementação

### Backend:
- [ ] Criar tabelas `proposals` e `proposal_documents`
- [ ] Implementar endpoints CRUD
- [ ] Validações de valores e condições
- [ ] Job de expiração automática
- [ ] Sistema de notificações
- [ ] Upload de documentos (S3)
- [ ] Testes unitários e integração

### Frontend (Flutter):
- [ ] Criar models `Proposal`, `ProposalDocument`, etc.
- [ ] Criar DTOs de request/response
- [ ] Implementar repository de propostas
- [ ] Criar notifiers/controllers
- [ ] Tela de listagem (corretor e proprietário)
- [ ] Tela de criação de proposta
- [ ] Tela de detalhes da proposta
- [ ] Upload de documentos
- [ ] Sistema de notificações push

### Adicionar em `api_paths.dart`:
```dart
// Proposals
static const String proposalsCreate = '/proposals';
static const String proposalsRealtor = '/proposals/realtor';
static const String proposalsOwner = '/proposals/owner';
static const String proposalsDetail = '/proposals/detail';
static const String proposalsAccept = '/proposals/accept';
static const String proposalsReject = '/proposals/reject';
static const String proposalsCancel = '/proposals/cancel';
static const String proposalsUpdate = '/proposals/update';
static const String proposalsRealtorStats = '/proposals/realtor/stats';
static const String proposalsOwnerStats = '/proposals/owner/stats';
static const String proposalsDocumentsUploadUrl = '/proposals/documents/upload-url';
```

---

## 📝 Notas Finais

- Todos os valores monetários em **DECIMAL(15,2)**
- Timestamps em **ISO 8601** com timezone
- Paginação padrão: 20 itens por página
- Autenticação via JWT obrigatória
- Rate limiting: 100 requisições por minuto
- Logs de auditoria para todas as ações
- Documentos armazenados no S3
- Limite de 10MB por documento
- Formatos aceitos: PDF, JPG, PNG, DOCX
- Histórico de edições mantido
