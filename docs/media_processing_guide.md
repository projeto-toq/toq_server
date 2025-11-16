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
5. **Pipeline AWS**: Step Functions + Lambdas/MediaConvert validam entradas, geram thumbnails, transcoding de vídeo, compactações (ZIP), atualizam estado do job na SQS/DynamoRDS.
6. **Atualização do domínio**: worker backend recebe callback da Step Function (webhook ou fila de saída), atualiza DB (`listing_media_batches`), cria registros de mídias processadas e atualiza listing para `listingmodel.StatusPendingOwnerApproval`.
7. **Distribuição**: APIs de download (expostas via `POST`) retornam URLs GET pré-assinadas para ZIP completo, mídias individuais e thumbnails (com TTL configurável). Frontend pode exibir progresso consultando endpoint de status.

## 3. Estrutura de Domínio e Persistência
- **Novos modelos core (`internal/core/model/media_processing_model`)**
  - `MediaBatch`: representa um lote de uploads (campos sugeridos: `ID`, `ListingID`, `PhotographerUserID`, `Status` = `PENDING_UPLOAD|RECEIVED|PROCESSING|FAILED|READY`, `UploadManifest`, `ProcessingMetadata`).
  - `MediaAsset`: item individual associado a um lote (`MediaBatchID`, `Type` = `PHOTO_VERTICAL|PHOTO_HORIZONTAL|VIDEO_VERTICAL|VIDEO_HORIZONTAL|THUMBNAIL|ZIP|PROJECT_DOC|PROJECT_RENDER`, `SourceKey`, `ProcessedKey`, `Checksum`, `ContentType`, `Resolution`, `DurationSeconds`, `Title`, `Sequence`). Os campos `Title` e `Sequence` são persistidos e vinculados ao objeto S3 para ordenação do carrossel/prateleira de projetos e exibição ao cliente.
  - `MediaProcessingJob`: metadados de execução assíncrona (`BatchID`, `ExternalJobID` do Step Functions, timestamps, mensagens de erro).
- **Modelagem física proposta**
  - `listing_media_batches`
    - Colunas chave: `id`, `listing_id`, `photographer_user_id`, `status`, `upload_manifest_json`, `processing_metadata_json`, `received_at`, `processing_started_at`, `processing_finished_at`, `error_code`, `error_detail`, `deleted_at`.
    - Índices: (`listing_id`, `id` DESC) para recuperar o lote mais recente; `status` para dashboards/observabilidade.
    - Serviço permite múltiplos lotes ativos por listing; soft delete (`deleted_at`) remove lotes obsoletos da UI sem perder histórico para auditoria ou reprocessamento.
  - `listing_media_assets`
    - Colunas: `id`, `batch_id`, `asset_type`, `orientation`, `source_key`, `processed_key`, `thumbnail_key`, `checksum_sha256`, `content_type`, `bytes`, `resolution`, `duration_seconds`, `title`, `sequence`, `variant_metadata_json` (metadados adicionais como "tipo de planta" ou "render diurno/noturno").
    - Constraint de unicidade (`batch_id`, `sequence`) garante ordenação determinística do carrossel; `title` fica armazenado para exibição.
    - Foreign key `batch_id` → `listing_media_batches.id` com `ON DELETE CASCADE` para limpeza automática.
  - `listing_media_jobs`
    - Colunas: `id`, `batch_id`, `external_job_id`, `provider`, `status`, `input_payload_json`, `output_payload_json`, `started_at`, `finished_at`.
    - Mantém histórico de tentativas Step Functions/MediaConvert e permite reconciliar callbacks com batches.
  - Opcional futuro: `listing_media_asset_events` para trilhar auditoria fina (estado do asset, geração de variações). Só criar se houver necessidade explícita de rastreabilidade além de logs/métricas.
- **Relacionamentos e ciclo de vida**
  - `listing_media_batches` tem relação 1:N com `listing_media_assets` e 1:N com `listing_media_jobs`. O `listing_id` liga o lote ao anúncio e permite histórico de reprocessamentos.
  - Após o upload (status `RECEIVED`), o batch continua sendo a âncora do pipeline: armazena status de processamento, links para assets finais e erros estruturados. Mesmo em `READY`, o registro permanece para audit trail, reprocessamento e geração de downloads adicionais.
  - `MediaAsset` acompanha tanto o objeto bruto (`source_key`) quanto os derivados (`processed_key`, `thumbnail_key`). Assim, o backend pode servir downloads individuais e gerar ZIPs sem nova consulta ao S3.
  - `MediaProcessingJob` documenta cada execução assíncrona e facilita retentativas controladas. Em caso de falha, o batch fica em `FAILED`, preservando contexto para diagnosticar e permitir retry manual/automático. Cada solicitação de reprocessamento gera um novo registro em `listing_media_jobs`, mantendo o relacionamento com o batch e permitindo comparar tentativas.
- **Acompanhamento pós-upload**
  - Services atualizam `listing_media_batches.status` (→ `PROCESSING`, `READY`, `FAILED`) de acordo com callbacks SQS/Step Functions.
  - Endpoints de status e download consomem `listing_media_batches` + `listing_media_assets` para montar respostas (sem varrer S3).
  - Métricas/logs referenciam sempre `batch_id` e `listing_id`, garantindo rastreabilidade ponta a ponta.
  - Um batch `READY` continua disponível para: gerar novas URLs de download, iniciar reprocessamentos (criando novos registros em `listing_media_jobs`) e servir como snapshot histórico do conjunto de mídias aprovadas.
- **Repositórios/Adapters**
  - Criar port `internal/core/port/right/repository/mediaprocessingrepository` com métodos `CreateBatch`, `UpsertAssets`, `UpdateBatchStatus`, `ListBatchesByListing`, `GetBatchByID`, `RegisterProcessingJob`, `FinalizeProcessingJob`.
  - Adapter MySQL espelhando regra 2.1 (um arquivo por método, `converters`, `entities`).
- **Relacionamento com Listing**
  - `listing_service` consome novo `mediaprocessingservice` para transições de status (`StatusPendingPhotoProcessing` → `StatusPendingOwnerApproval` ao término do job). Nenhuma lógica HTTP dentro do service.

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
- **Regras adicionais**
  - `sequence` deve ser numeração inteira positiva e única dentro do lote para definir a ordem do carrossel ou da prateleira de projetos.
  - `title` será persistido e exibido na galeria; deve refletir o ambiente, render ou documento.
  - `mediaType` aceita `project/doc` (PDF, imagens de plantas) e `project/render` (renders) quando o listing está marcado como empreendimento em construção e o solicitante é o owner associado.
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
        "s3Key": "123/raw/photo/vertical/2025-11-11/photo-vertical-1.jpg",
        "title": "Sala Social",
        "sequence": 1
      }
    ]
  }
  ```
- **Códigos HTTP**
  - `201 Created`: lote criado e URLs emitidas.
  - `400 Bad Request`: manifesto inválido (tipos/limites).
  - `401 Unauthorized` / `403 Forbidden`: usuário sem permissão no listing.
  - `404 Not Found`: listing inexistente ou status diferente de `PENDING_PHOTO_PROCESSING`.
  - `409 Conflict`: lote anterior ainda em processamento pendente (regra: um lote aberto por vez).

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
- **Regras**
  - Backend executa `HEAD`/`GetObjectAttributes` para validar existência, tamanho, checksum.
  - Atualiza `MediaBatch.Status` → `RECEIVED`, persiste metadados por arquivo e publica job na fila.
  - Campos `title` e `sequence` são consolidados com o manifesto inicial e gravados no `MediaAsset` correspondente para uso posterior no carrossel e nos downloads. Para `project/doc` e `project/render`, validar `contentType` (`application/pdf` ou `image/*`) e aplicar sanitização de nome de arquivo.
  - `batchId` obrigatório: confirma apenas o lote especificado (idempotência garantida).
- **Responses**
  - `202 Accepted`: job enviado a processamento (retorna `processingJobId` e tempo estimado).
  - `409 Conflict`: arquivos faltantes ou divergência de checksum.
  - `410 Gone`: URLs expiradas (frontend deve reiniciar fluxo 4.1).

### 4.3 `POST /api/v2/listings/media/uploads/status` — Status do processamento
- **Justificativa**: seguindo a diretriz interna, buscas/listagens sem paginação ou filtros devem usar `POST` com parâmetros no corpo.
- **Request Body**
  ```json
  {
    "listingId": 123,
    "batchId": "UUID"
  }
  ```
- **Observações**
  - `listingId` é obrigatório e referencia o anúncio associado ao lote.
  - `batchId` é obrigatório e identifica o lote monitorado.
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
- **HTTP**: `200 OK`, `404 Not Found`, `403 Forbidden`.

### 4.4 `POST /api/v2/listings/media/downloads/query` — URLs de download
- **Justificativa**: consulta sem paginação ou filtros extras; parâmetros enviados no corpo (regra interna).
- **Request Body**
  ```json
  {
    "listingId": 123,
    "batchId": "UUID"
  }
  ```
- **Observações**
  - `listingId` é obrigatório.
  - `batchId` é opcional; ausência implica retornar o último lote com status `READY`.
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
            "preview": "https://cdn/...thumb.jpg",
            "url": "https://s3/...",
            "resolution": "4000x6000",
            "title": "Sala Social",
            "sequence": 1
          }
        ]
      }
    ]
  }
  ```
- **Códigos**: `200 OK`, `204 No Content` (nenhum lote pronto), `404` (listing/batch inexistente).

### 4.5 `POST /api/v2/listings/media/uploads/retry` — Reprocessar lote existente
- **Objetivo**: permitir que o fotógrafo reenvie um lote para processamento quando o status estiver `FAILED` ou quando alguma derivação precise ser regenerada.
- **Request Body**
  ```json
  {
    "listingId": 123,
    "batchId": "UUID",
    "reason": "manual-check"
  }
  ```
- **Regras**
  - Apenas lotes em `FAILED` ou `READY` podem ser reprocessados. Em `READY`, o reprocessamento gera novas versões dos assets preservando o histórico anterior.
  - Sempre cria um novo registro em `listing_media_jobs` vinculado ao batch, com `provider` = `STEP_FUNCTIONS` e `status` inicial `PENDING`.
  - O reprocessamento não gera novos uploads; reutiliza os objetos S3 existentes na pasta `raw/`.
  - Soft delete opcional pode ocultar lotes obsoletos após o novo processamento.
- **Responses**
  - `202 Accepted`: job reenfileirado, retorna `jobId` e status atual.
  - `403 Forbidden`: usuário não autorizado.
  - `404 Not Found`: batch inexistente para o listing fornecido.
  - `409 Conflict`: lote em `PENDING_UPLOAD`/`RECEIVED` ainda não finalizou upload.

## 5. Serviços Core e Dependências (Hexagonal)
- **Service**: `internal/core/service/media_processing_service`
  - Métodos públicos (cada um em arquivo próprio): `CreateUploadBatch`, `CompleteUploadBatch`, `GetBatchStatus`, `ListDownloadUrls`, `HandleProcessingCallback`.
  - Depende de: `mediaprocessingrepository`, `storageport.CloudStoragePortInterface`, `queueport.MediaProcessingQueue`, `listingrepository` (para atualizar status), `globalservice` (transações, notificações push).
- **Ports adicionais**
  - `internal/core/port/right/queue/mediaprocessingqueue` (Adapter SQS). Métodos: `EnqueueJob`, `ParseJobMessage`, `AcknowledgeJob`.
  - `internal/core/port/right/functions/mediaprocessingcallback` caso usemos WebHook (API Gateway + Lambda) para retornos do Step Functions.
- **Handlers HTTP** usam DTOs em `internal/adapter/left/http/dto/listing_dto`. Validação com `validation_service`.

## 6. AWS — Serviços Recomendados e Justificativas
- **Amazon S3**: armazenamento de mídias brutas e processadas. Pré-assinados PUT/GET já suportados pelo adapter atual; utilizar o bucket `toq-listing-medias` com prefixos `/{listingId}/{stage}/...` (`stage`: `raw`, `processed`, `zip`, `thumb`).
- **AWS SQS (Standard Queue)**: desacopla API HTTP do pipeline de processamento; permite retry e DLQ (`listing-media-processing-dlq`).
- **AWS Step Functions**: orquestra pipeline (validação → fan-out → fan-in). Facilita paralelizar fotos/vídeos, integrar com MediaConvert e garantir idempotência.
- **AWS Lambda**: funções menores para: validar manifest vs S3, gerar thumbnails (Sharp/libvips), criar ZIP incremental, gerar payload de callback.
- **AWS Elemental MediaConvert**: transcodificação de vídeo para múltiplas resoluções (480p/720p/1080p) e codecs otimizados (H.264 + HLS/MP4). Alternativa: Elastic Transcoder (legado) ou AWS Batch se preferir docker custom.
- **Amazon DynamoDB ou S3 Metadata** (opcional): rastrear progresso do Step Functions; entretanto recomendável persistir estado final no MySQL via callback para manter single source of truth.
- **Amazon CloudFront** (futuro): CDN para downloads públicos temporários; alternativa via S3 + presigned URLs com restrição IP.
- **AWS KMS**: criptografia de objetos sensíveis (`SSE-KMS`) com políticas mínimas para fotógrafos.

### 6.1 Sequência detalhada do pipeline de processamento
1. **Ingestão e validação (Lambda #1)**
  - Disparada ao mudar `MediaBatch` para `RECEIVED`.
  - Tarefas: verificar existência dos objetos brutos (`GetObjectAttributes` S3), checar checksum/bytes, normalizar metadados (orientação, content-type) e enriquecer `upload_manifest_json` com dados consolidados.
  - Em caso de falha: marca batch como `FAILED`, envia mensagem para DLQ e notifica observabilidade.
2. **Geração de thumbnails (Lambda #2)**
  - Usa `Sharp`/`libvips` ou alternativa para processar imagens JPEG/HEIC.
  - Outputs: escreve versões `thumb` no bucket `toq-listing-medias` (`processed/{variant}`), atualiza `listing_media_assets.thumbnail_key`.
3. **Transcodificação de vídeos (AWS MediaConvert)**
  - Step Functions cria job MediaConvert com presets otimizados (ex.: 1080p/720p MP4 + HLS opcional).
  - MediaConvert lê do prefixo `raw/` e grava no prefixo `processed/video/{quality}/`.
  - Status do job é monitorado via CloudWatch Events → Step Functions task; em caso de erro, lote vai para `FAILED` e o job é documentado em `listing_media_jobs`.
4. **Geração de ZIPs (Lambda #3)**
  - Após fan-in das etapas de imagem e vídeo, Lambda monta pacote ZIP com todas as mídias processadas.
  - Escreve `processed/zip/full.zip` e (opcional) `processed/zip/thumbnails.zip`.
  - Atualiza `listing_media_assets` com entradas específicas de tipo `ZIP` e respectivo `processed_key`.
5. **Atualização de metadados finais (Lambda #4)**
  - Consolida dados (resoluções, duração, URLs de preview) e persiste em `listing_media_assets`.
  - Marca `listing_media_batches.status` como `READY`, registra `processing_finished_at` e gera evento (e-mail/push) através do `global_service`.
6. **Callback para backend**
  - Step Functions envia payload para endpoint interno (`media_processing_service.HandleProcessingCallback`) contendo `batchId`, `status`, métricas e caminhos gerados.
  - Service valida idempotência, atualiza banco, dispara notificações e, se `READY`, chama `listing_service` para transição do anúncio.

### 6.2 Componentes auxiliares
- **Amazon EventBridge/CloudWatch Events**: monitoram jobs MediaConvert e acionar Step Functions/Lambdas de follow-up.
- **Amazon SQS (DLQ)**: captura falhas irrecuperáveis; operadores podem consultar e reiniciar lote.
- **AWS Step Functions Map/Parallel**: orquestra fases distintas (imagens vs vídeos) em paralelo respeitando dependências.
- **AWS Batch (opcional)**: se for necessário processamento pesado custom (ex.: IA para seleção de fotos), pode substituir Lamdbas específicas mantendo integração com Step Functions.

## 7. Convenções de Segurança e Compliance
- **Controle de acesso**: validar `PhotographerUserID` associado ao booking ativo; somente o fotógrafo do booking (ou um operador com privilégio explícito `media:manage` em ferramentas internas) pode criar lotes, confirmar uploads ou solicitar reprocessamento de fotos/vídeos. **Exceção**: listings marcados como empreendimento em construção aceitam que o owner autenticado (ou procurador cadastrado) crie lotes `project/doc` e `project/render`, desde que possua relação comprovada com o listing.
- **Prefixos S3**: `/{listingId}/raw/{mediaType}/` e `/{listingId}/processed/{variant}/` dentro do bucket `toq-listing-medias`. Nunca reutilizar chaves previsíveis sem UUID (`batchId` + `clientId`). `mediaType` inclui `photo/vertical`, `photo/horizontal`, `video/vertical`, `video/horizontal`, `project/doc`, `project/render`.
- **TTL de URLs**: upload 15 min, download 60 min (config via `env.S3.signedUrlTTL`).
- **Checksums**: exigir `x-amz-checksum-sha256` para uploads; backend valida no `Complete`.
- **Límites**: máximo 60 arquivos por lote, 1 GiB total; rejeitar `contentType` fora da whitelist (`image/jpeg`, `image/heic`, `video/mp4`, `video/quicktime`, `application/pdf`).
- **Auditoria**: registrar logs estruturados com `listing_id`, `batch_id`, `photographer_user_id`, `job_id` em cada etapa. Utilizar `utils.LoggerFromContext`.
- **Idempotência**: este fluxo **não** suporta `Idempotency-Key`. A prevenção de duplicidade é garantida por regras de status (`PENDING_UPLOAD` → `RECEIVED` → `PROCESSING` → `READY|FAILED`) e por validações de manifesto (`batchId` obrigatório em todas as operações).

## 8. Orquestração Assíncrona & Estados
- **Estados do lote** *(confirmado em 16/11/2025 com plataforma + frontend; não haverá novos estados no MVP)*
  - `PENDING_UPLOAD`: URLs geradas, aguardando confirmação.
  - `RECEIVED`: todos os arquivos confirmados.
  - `PROCESSING`: job Step Functions em execução.
  - `READY`: processamento finalizado com sucesso; listing pode avançar no funil.
  - `FAILED`: Step Functions retornou erro; permitir retry (gera novo job) mantendo mídia bruta.
- **Eventos**
  - Ao receber `READY`, service atualiza listing → `StatusPendingOwnerApproval` e dispara notificação push/e-mail.
  - Em `FAILED`, mantém listing no status atual e gera alerta Slack/Prometheus (ver seção observabilidade).
- **Retry**
  - Lambda com retry automático; Step Functions configurado com `catch` e envio à DLQ (SQS). Backend expõe endpoint `POST /api/v2/listings/media/uploads/retry` que sempre cria um novo `listing_media_jobs` e reaproveita os objetos brutos existentes.

## 9. Observabilidade
- **Tracing**: usar `utils.GenerateTracer` nos métodos públicos do `media_processing_service`; propagar `trace_id` para mensagens SQS (atributo `Traceparent`).
- **Métricas** (todas obrigatórias)
  - Counter `media_batches_created_total{status="PENDING_UPLOAD|RECEIVED|PROCESSING|READY|FAILED"}`.
  - Counter `media_batches_failed_total` (incrementa a cada batch que termina com `FAILED`).
  - Counter `media_reprocess_requests_total` (incrementa a cada chamada aceita em `/uploads/retry`).
  - Histogram `media_processing_duration_seconds` (tempo entre `RECEIVED` e `READY`).
  - Gauge `media_processing_inflight_jobs` (batches com status `PROCESSING`).
  - Counter `media_processing_dlq_messages_total` (mensagens redirecionadas para DLQ).
  - Counter `media_processing_callback_errors_total` (falhas ao aplicar callback no backend).
- **Logs**: chaves padrão (`listing_id`, `batch_id`, `job_id`, `stage`). Nível `Error` somente infraestrutura (S3/SQS/MediaConvert). Regra de span error no catch.
- **Alertas**: se `media_processing_duration_seconds` p95 > 10 min ou `FAILED` > 0 nos últimos 15 min, disparar alerta.

## 10. Considerações para o Frontend
- Fazer upload direto via `fetch`/XHR com `PUT` para URL retornada.
- Respeitar headers obrigatórios (especialmente `Content-Type` e checksum).
- Armazenar `clientId` por arquivo para reconciliar na confirmação.
- Incluir `title` e `sequence` em cada item do manifesto, garantindo unicidade da sequência dentro do lote para ordenar o carrossel exibido aos clientes finais.
- Exibir progresso através do endpoint de status; fallback com polling a cada 10s até `READY` ou `FAILED`.
- Em caso de `FAILED`, permitir retry automático solicitando novo lote (chamar novamente 4.1).
- Para download final, usar URLs GET pré-assinadas; frontend deve iniciar download imediatamente para evitar expiração.

## 11. Próximos Passos (macro backlog)
1. Modelagem persitência (`listing_media_batches`, `listing_media_assets`, `listing_media_jobs`).
2. Implementar `media_processing_service` + ports/repositorios + DTOs + handlers.
3. Criar adapter SQS (`internal/adapter/right/aws_sqs/mediaprocessing`) com instrumentação.
4. Provisionar pipeline AWS (Infra as Code) e integração Step Functions ↔ callback HTTP.
5. Atualizar `docs/swagger.yaml` com novos endpoints e esquemas.
6. Ajustar frontend conforme contrato definido neste guia.

## 12. Housekeeping & Limpeza de Mídias (futuro)
- **Escopo confirmado:** a exclusão física/expiração de mídias antigas não fará parte desta entrega. Os serviços atuais devem apenas registrar `deleted_at` (soft delete) e manter os objetos brutos/processados disponíveis para reprocessamentos.
- **Job dedicado:** será criado futuramente um worker/batch dedicado para limpeza (permanently delete) com base em políticas de retenção e sinalizações de soft delete. O job consumirá uma API interna que listará lotes elegíveis.
- **Dependências:** Cloud/Admin devem provisionar recursos necessários (ex.: Lambda Scheduler ou ECS Fargate) quando o backlog autorizar. Até lá, nenhum componente deve remover objetos automaticamente.

---
Este documento é ponto de partida oficial para desenvolvimento paralelo frontend/backend. Qualquer alteração de contrato deve ser negociada e refletida aqui antes da implementação.
