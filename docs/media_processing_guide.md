# Guia de Processamento de Mídias — Listings em Status `PENDING_PHOTO_PROCESSING`

## 1. Contexto e Objetivo
- Fluxo vigente: `listingmodel.StatusPendingPhotoProcessing` indica que a sessão de fotos foi concluída e o fotógrafo precisa subir mídias brutas (fotos/vídeos) para posterior tratamento.
- Objetivo deste guia: alinhar backend, frontend e plataforma sobre o processo completo de upload, processamento assíncrono e disponibilização de mídias finais (ZIP completo, mídias individuais e thumbnails), respeitando arquitetura hexagonal, segurança e escalabilidade do TOQ Server.
- Escopo: novas APIs REST (`/api/v2`), serviços core, adapters S3/SQS/Lambda/MediaConvert, modelos de domínio, monitoramento e contratos com o frontend.

## 2. Visão Geral do Fluxo (alto nível)
1. **Manifesto de upload**: frontend envia lista de arquivos que serão enviados (tipo, orientação, extensão, tamanho estimado, checksum). Backend valida permissão e status do listing.
2. **Emissão de URLs pré-assinadas**: backend retorna URLs PUT S3 com metadados obrigatórios. Fotógrafo envia diretamentre para o bucket dedicado `toq-listing-medias`, sempre dentro do diretório raiz `/{listingId}/`. Listings classificados como empreendimento em construção permitem que o **owner** autenticado (ou representante designado) gere URLs para projetos e renders sem depender do fotógrafo.
3. **Confirmação de upload**: após concluir todos os uploads com sucesso HTTP 200 na S3, frontend chama endpoint de confirmação entregando o manifesto final com ETags retornadas pelo S3.
4. **Orquestração assíncrona**: backend persiste lote, muda status interno do lote para `RECEIVED` e publica job em SQS para pipeline de processamento.
5. **Pipeline AWS**: Step Functions + Lambdas (Go) validam entradas, geram thumbnails, transcoding de vídeo, compactações (ZIP), atualizam estado do job na SQS/DynamoRDS.
6. **Atualização do domínio**: worker backend recebe callback da Step Function (webhook ou fila de saída), atualiza DB (`listing_media_batches`), cria registros de mídias processadas e atualiza listing para `listingmodel.StatusPendingOwnerApproval`.
7. **Distribuição**: APIs de download (expostas via `POST`) retornam URLs GET pré-assinadas para ZIP completo, mídias individuais e thumbnails (com TTL configurável). Frontend pode exibir progresso consultando endpoint de status.

## 3. Estrutura de Domínio e Persistência
- **Novos modelos core (`internal/core/model/media_processing_model`)**
  - `MediaBatch`: representa um lote de uploads (campos sugeridos: `ID`, `ListingID`, `PhotographerUserID`, `Status` = `PENDING_UPLOAD|RECEIVED|PROCESSING|FAILED|READY`, `UploadManifest`, `ProcessingMetadata`).
  - `MediaAsset`: item individual associado a um lote (`MediaBatchID`, `Type` = `PHOTO_VERTICAL|PHOTO_HORIZONTAL|VIDEO_VERTICAL|VIDEO_HORIZONTAL|THUMBNAIL|ZIP|PROJECT_DOC|PROJECT_RENDER`, `SourceKey`, `ProcessedKey`, `ThumbnailKey`, `Checksum`, `ContentType`, `Resolution`, `DurationSeconds`, `Title`, `Sequence`). Os campos `Title` e `Sequence` são persistidos e vinculados ao objeto S3 para ordenação do carrossel/prateleira de projetos e exibição ao cliente.
  - `MediaProcessingJob`: metadados de execução assíncrona (`BatchID`, `ExternalJobID` do Step Functions, timestamps, mensagens de erro).

## 4. Endpoints HTTP (Left Adapter)
Todos sob `/api/v2/listings/media/*`, com `listingId` presente apenas no corpo das requisições (padrão adotado pelo projeto). Handlers em `internal/adapter/left/http/handlers/listing_handlers/` (um arquivo por handler). Autorização: somente fotógrafo designado ou roles administrativas conforme política de permissionamento.

### 4.1 `POST /api/v2/listings/media/uploads` — Solicitar URLs de upload
- **Objetivo**: cliente envia manifesto contendo arquivos que pretende subir.
- **Request Body**
  ```json
  {
    "listingId": 123,
    "batchReference": "2025-11-11T14:00Z-slot-123",
    "files": [
      {
        "clientId": "photo-vertical-1",
        "mediaType": "PHOTO_VERTICAL",
        "orientation": "VERTICAL",
        "filename": "livingroom_vert_01.jpg",
        "contentType": "image/jpeg",
        "bytes": 4231987,
        "checksum": "sha256:...",
        "title": "Sala Social",
        "sequence": 1
      }
    ]
  }
  ```
- **Response 201**
  ```json
  {
    "batchId": "UUID",
    "uploadUrlTtlSeconds": 900,
    "files": [
      {
        "clientId": "photo-vertical-1",
        "uploadUrl": "https://s3/...",
        "method": "PUT",
        "headers": {
          "Content-Type": "image/jpeg",
          "x-amz-checksum-sha256": "..."
        },
        "s3Key": "123/raw/photo/vertical/photo-vertical-1.jpg",
        "title": "Sala Social",
        "sequence": 1
      }
    ]
  }
  ```
  *Nota: O caminho S3 não inclui mais segmentos de data (YYYY-MM-DD).*

### 4.2 `POST /api/v2/listings/media/uploads/complete` — Confirmação de upload
- **Request Body**
  ```json
  {
    "listingId": 123,
    "batchId": "UUID",
    "files": [
      {
        "clientId": "photo-vertical-1",
        "s3Key": "123/raw/photo/vertical/...",
        "etag": "\"9b2cf535f27731c974343645a3985328\"",
        "bytes": 4231987
      }
    ]
  }
  ```

### 4.3 `POST /api/v2/listings/media/uploads/status` — Status do processamento
- **Request Body**
  ```json
  {
    "listingId": 123,
    "batchId": "UUID"
  }
  ```
- **Response 200**
  ```json
  {
    "batchId": "UUID",
    "status": "PROCESSING",
    "submittedAt": "2025-11-11T17:03:12Z",
    "processing": {
      "jobId": "arn:aws:states:...",
      "steps": [
        { "name": "validate-raw-assets", "status": "SUCCEEDED" },
        { "name": "image-thumbnails", "status": "RUNNING" },
        { "name": "video-transcode", "status": "PENDING" }
      ]
    },
    "availableDownloads": [
      {
        "type": "ZIP_FULL",
        "ready": false
      },
      {
        "type": "THUMBNAILS",
        "ready": true
      }
    ]
  }
  ```

### 4.4 `POST /api/v2/listings/media/downloads/query` — URLs de download
- **Request Body**
  ```json
  {
    "listingId": 123,
    "batchId": "UUID"
  }
  ```
- **Response 200**
  ```json
  {
    "batchId": "UUID",
    "generatedAt": "2025-11-11T18:22:00Z",
    "ttlSeconds": 3600,
    "downloads": [
      {
        "type": "ZIP_FULL",
        "url": "https://s3/...",
        "expiresAt": "2025-11-11T19:22:00Z"
      },
      {
        "type": "PHOTO_VERTICAL",
        "items": [
          {
            "label": "Sala - Vertical",
            "preview": "https://cdn/.../thumbnail/photo.jpg",
            "url": "https://s3/.../large/photo.jpg",
            "resolution": "4000x6000",
            "title": "Sala Social",
            "sequence": 1
          }
        ]
      }
    ]
  }
  ```

## 5. Serviços Core e Dependências (Hexagonal)
- **Service**: `internal/core/service/media_processing_service`
- **Ports**: `mediaprocessingrepository`, `storageport`, `queueport`.

## 6. AWS — Serviços Recomendados e Justificativas
- **Amazon S3**: armazenamento de mídias brutas e processadas.
  - Prefixos: `/{listingId}/raw/{mediaType}/` e `/{listingId}/processed/{mediaType}/{size}/`.
  - Tamanhos de imagem: `thumbnail` (200px), `small` (400px), `medium` (800px), `large` (1200px).
- **AWS SQS**: desacopla API HTTP do pipeline.
- **AWS Step Functions**: orquestra pipeline.
- **AWS Lambda (Go)**:
  - `validate`: Validação de manifesto.
  - `thumbnails`: Redimensionamento e correção de rotação (EXIF).
  - `zip`: Geração de pacotes para download.
  - `consolidate`: Agregação de resultados.
  - `callback`: Notificação ao backend.

### 6.1 Sequência detalhada do pipeline de processamento
1. **Ingestão e validação (Lambda #1)**
  - Disparada ao mudar `MediaBatch` para `RECEIVED`.
2. **Geração de thumbnails (Lambda #2)**
  - Usa `disintegration/imaging` (Go) para processar imagens.
  - Outputs: escreve versões `thumbnail`, `small`, `medium`, `large` no bucket.
3. **Transcodificação de vídeos (AWS MediaConvert)**
  - Step Functions cria job MediaConvert.
4. **Geração de ZIPs (Lambda #3)**
  - Lambda monta pacote ZIP com todas as mídias processadas.
  - Estrutura interna do ZIP é limpa (sem pastas de sistema).
5. **Atualização de metadados finais (Lambda #4)**
  - Consolida dados e persiste.
6. **Callback para backend**
  - Step Functions envia payload para endpoint interno.

## 7. Convenções de Segurança e Compliance
- **Controle de acesso**: validar `PhotographerUserID`.
- **Prefixos S3**: `/{listingId}/raw/{mediaType}/` e `/{listingId}/processed/{mediaType}/{size}/`.
- **TTL de URLs**: upload 15 min, download 60 min.
- **Checksums**: exigir `x-amz-checksum-sha256`.

## 8. Orquestração Assíncrona & Estados
- **Estados do lote**: `PENDING_UPLOAD`, `RECEIVED`, `PROCESSING`, `READY`, `FAILED`.

## 9. Observabilidade
- **Tracing**: `utils.GenerateTracer`.
- **Métricas**: Prometheus counters e histograms.

## 10. Considerações para o Frontend
- Fazer upload direto via `fetch`/XHR com `PUT`.
- Respeitar headers obrigatórios.
- Exibir progresso através do endpoint de status.
- Para download final, usar URLs GET pré-assinadas.

---
Este documento é ponto de partida oficial para desenvolvimento paralelo frontend/backend. Qualquer alteração de contrato deve ser negociada e refletida aqui antes da implementação.
