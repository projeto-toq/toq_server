# Plano de Refatoração: Sistema de Processamento de Mídia (TOQ Server)

Este documento detalha o plano de execução para a refatoração do sistema de processamento de mídia, visando atender aos requisitos de upload incremental, gestão granular de assets e separação de responsabilidades entre processamento e finalização.

**Objetivo:** Migrar de uma arquitetura baseada em Lotes (Batch) para uma arquitetura baseada em Assets (Mídias Individuais), permitindo uploads parciais, atualizações de metadados e feedback visual imediato.

---

## 📋 Visão Geral das Fases

1.  **Fase 1: Fundação (Banco de Dados e Modelos)** [CONCLUÍDO]
2.  **Fase 2: Camada de Persistência (Repositórios)** [CONCLUÍDO]
3.  **Fase 3: Lógica Core - Upload e Processamento** [CONCLUÍDO]
4.  **Fase 4: Lógica Core - Gestão (CRUD)** [CONCLUÍDO]
5.  **Fase 5: Lógica Core - Finalização** [CONCLUÍDO]
6.  **Fase 6: Camada HTTP (Handlers)** [CONCLUÍDO]
7.  **Fase 7: Infraestrutura AWS** [CONCLUÍDO]
8.  **Fase 8: Documentação** [CONCLUÍDO]

---

## 🚀 Detalhamento das Fases

### Fase 1: Fundação (Banco de Dados e Modelos)

**Objetivo:** Estabelecer a estrutura de dados que suportará a gestão granular de mídias.

#### 1.1. Modelagem de Dados (SQL)
*Ação:* Solicitar ao DBA a execução do script abaixo.
*Arquivo:* `scripts/refactor_media_assets.sql` (Sugestão)

```sql
-- Tabela para gestão individual de assets
-- Chave única composta garante a regra: Listing + Tipo + Sequência
CREATE TABLE media_assets (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    listing_id BIGINT NOT NULL,
    asset_type VARCHAR(50) NOT NULL, -- Ex: PHOTO_HORIZONTAL, VIDEO_VERTICAL
    sequence INT NOT NULL,
    status VARCHAR(50) NOT NULL,     -- PENDING_UPLOAD, PROCESSING, PROCESSED, FAILED
    s3_key_raw VARCHAR(255),
    s3_key_processed VARCHAR(255),
    title VARCHAR(255),
    metadata JSON,                   -- Metadados flexíveis (ex: client_id, checksum)
    
    -- Garante que não existam duas "Foto Horizontal 1" para o mesmo listing
    UNIQUE KEY uk_listing_asset_seq (listing_id, asset_type, sequence),
    INDEX idx_listing_status (listing_id, status)
);
```

#### 1.2. Entidade de Domínio
*Arquivo:* `internal/core/model/media_processing_model/media_asset.go`
*Descrição:* Criar a struct `MediaAsset` sem campos de auditoria (`created_at`, `updated_at`), focada na regra de negócio.

```go
type MediaAsset struct {
    id             uint64
    listingID      uint64
    assetType      MediaAssetType
    sequence       uint8
    status         MediaAssetStatus
    s3KeyRaw       sql.NullString
    s3KeyProcessed sql.NullString
    title          sql.NullString
    metadata       sql.NullString // JSON
}
// Métodos: NewMediaAsset, Getters, Setters, IsProcessed(), etc.
```

#### 1.3. DTOs (Data Transfer Objects)
*Arquivo:* `internal/core/domain/dto/media_dto.go`
*Descrição:* Definir as estruturas de entrada e saída para as novas operações.
*   `RequestUploadURLsInput`: Lista de arquivos para upload (sem batchId).
*   `ProcessMediaInput`: Gatilho para processamento.
*   `ListDownloadURLsInput`: Filtros para listagem.
*   `UpdateMediaInput`: Edição de título/sequência.
*   `DeleteMediaInput`: Remoção de asset.

---

### Fase 2: Camada de Persistência (Repositórios)

**Objetivo:** Implementar o acesso a dados para a nova tabela `media_assets`.

#### 2.1. Definição do Port
*Arquivo:* `internal/core/port/right/repository/mediaprocessingrepository/media_repo_port.go`
*Interface:* `MediaRepositoryPort`
*   `UpsertAsset`: Cria ou atualiza (on duplicate key update).
*   `GetAsset`: Busca por chave composta (ListingID, Type, Sequence).
*   `ListAssets`: Busca lista com filtros (Status, Type).
*   `DeleteAsset`: Remove registro.

#### 2.2. Implementação do Adapter MySQL
*Diretório:* `internal/adapter/right/mysql/media_processing/`
*Arquivos:*
*   `media_processing_adapter.go`: Struct e construtor.
*   `upsert_asset.go`: Implementação do Upsert.
*   `get_asset.go`: Implementação do Get.
*   `list_assets.go`: Implementação do List.
*   `delete_asset.go`: Implementação do Delete.
*   `converters/`: Mapeamento `MediaAsset` (Domain) <-> `MediaAssetEntity` (DB).

---

### Fase 3: Lógica Core - Upload e Processamento

**Objetivo:** Permitir o fluxo de upload incremental e processamento parcial.

#### 3.1. Refatorar Solicitação de Upload
*Arquivo:* `internal/core/service/media_processing_service/request_upload_urls.go` (Renomear de `create_upload_batch.go`)
*Lógica:*
1.  Receber lista de arquivos.
2.  Validar regras de negócio (tipos permitidos, tamanhos).
3.  **Mudança Crítica:** Validar unicidade baseada em `(AssetType, Sequence)` e não globalmente.
4.  Gerar URLs pré-assinadas (PUT) para o S3.
5.  Persistir assets com status `PENDING_UPLOAD` usando `UpsertAsset`.
6.  Retornar URLs para o frontend.

#### 3.2. Implementar Gatilho de Processamento
*Arquivo:* `internal/core/service/media_processing_service/process_media.go` (Novo)
*Lógica:*
1.  Listar assets do listing com status `PENDING_UPLOAD`.
2.  Se houver assets:
    *   Montar payload para Step Function.
    *   Invocar Step Function de **Processamento** (Validação + Thumbnails).
    *   Atualizar status dos assets para `PROCESSING`.
3.  Retornar sucesso (202 Accepted).

#### 3.3. Refatorar Listagem de Downloads
*Arquivo:* `internal/core/service/media_processing_service/list_download_urls.go`
*Lógica:*
1.  Listar **todos** os assets do listing (independente de estarem prontos ou não).
2.  Para cada asset:
    *   Se `PROCESSED`: Gerar URL GET assinada para a chave processada (otimizada).
    *   Se `PENDING/PROCESSING`: Gerar URL GET assinada para a chave RAW (se permitido) ou retornar apenas metadados indicando status.
3.  Permitir que o frontend mostre o progresso real ("Processando...", "Pronto").

---

### Fase 4: Lógica Core - Gestão (CRUD)

**Objetivo:** Permitir correções e organização das mídias antes da finalização.

#### 4.1. Atualização de Mídia
*Arquivo:* `internal/core/service/media_processing_service/update_media.go` (Novo)
*Lógica:*
1.  Buscar asset por `(ListingID, Type, Sequence)`.
2.  Atualizar campos permitidos (`Title`, `Metadata`, `Sequence` - cuidado com colisão de sequência).
3.  Persistir alterações.
4.  **Opcional:** Disparar reprocessamento se necessário.

#### 4.2. Remoção de Mídia
*Arquivo:* `internal/core/service/media_processing_service/delete_media.go` (Novo)
*Lógica:*
1.  Buscar asset.
2.  Remover arquivos do S3 (Raw e Processed).
3.  Remover registro do banco.

#### 4.3. Listagem Geral (Backoffice/Frontend)
*Arquivo:* `internal/core/service/media_processing_service/list_media.go` (Novo)
*Lógica:*
1.  Expor funcionalidade de `ListAssets` do repositório com filtros ricos (paginação, tipos específicos).

---

### Fase 5: Lógica Core - Finalização

**Objetivo:** Gerar o pacote final (ZIP) e avançar o status do Listing.

#### 5.1. Refatorar Finalização
*Arquivo:* `internal/core/service/media_processing_service/complete_media.go` (Renomear de `complete_upload_batch.go`)
*Lógica:*
1.  Verificar se existem assets em `PENDING_UPLOAD` ou `PROCESSING`. Se sim, bloquear ou aguardar.
2.  Invocar Step Function de **Finalização** (Gerar ZIP).
3.  Atualizar status do Listing para `StatusPendingOwnerApproval`.
4.  Não reprocessar imagens (isso já foi feito na fase 3).

---

### Fase 6: Camada HTTP (Handlers)

**Objetivo:** Expor as novas funcionalidades via API REST.

*Diretório:* `internal/adapter/left/http/handlers/listing_handlers/`

1.  **`request_upload_urls_handler.go`**: Endpoint `POST /listings/media/uploads`.
2.  **`process_media_handler.go`**: Endpoint `POST /listings/media/uploads/process`.
3.  **`list_download_urls_handler.go`**: Endpoint `POST /listings/media/downloads` (Ajustar contrato).
4.  **`update_media_handler.go`**: Endpoint `POST /listings/media/update`.
5.  **`delete_media_handler.go`**: Endpoint `DELETE /listings/media`.
6.  **`list_media_handler.go`**: Endpoint `GET /listings/` (com filtros de media).
7.  **`complete_media_handler.go`**: Endpoint `POST /listings/media/uploads/complete`.

*Ação:* Atualizar anotações Swagger em todos os handlers.

---

### Phase 7: AWS Infrastructure (Serverless)
- [x] Update Step Function definition (`media_processing_pipeline.json`) to remove `BatchID` and use `ListingID`.
- [x] Update `validate` Lambda to parse SQS events correctly and trigger Step Function with `ListingID`.
- [x] Update `thumbnails` Lambda to use `ListingID` for logging/metrics.
- [x] Update `zip` Lambda to group by `ListingID` instead of `BatchID`.
- [x] Update `consolidate` Lambda to aggregate results by `ListingID`.

---

### Fase 8: Documentação
- [x] Atualizar `docs/media_processing_guide.md` com novo fluxo de assets.
- [x] Atualizar `docs/aws_media_processing_useful_commands.md` com novos payloads.
- [x] Atualizar `docs/aws_media_processing_implementation_summary.md` com novas estruturas.
- [x] Atualizar `aws/README.md` com novos caminhos S3.

---

## ⚠️ Pontos de Atenção

*   **Compatibilidade:** Como é um ambiente de dev sem back-compatibility, podemos apagar os dados antigos das tabelas de batch se necessário.
*   **Concorrência:** O uso de `Upsert` e chaves únicas no banco deve prevenir condições de corrida em uploads simultâneos.
*   **Observabilidade:** Manter logs estruturados e tracing em todas as etapas (Service -> Repo -> AWS).
