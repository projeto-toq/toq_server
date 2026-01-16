# Plano para habilitar uploads de mídia de projetos (casa em construção)

## Objetivo
Permitir que proprietários de imóveis em construção façam upload de plantas (PDF) e renders (JPEG) durante a criação do anúncio, sem depender do status PENDING_PHOTO_PROCESSING, reutilizando a infraestrutura de mídia existente com segregação de fluxos e governança por feature flag.

## Contexto e premissas
- Pipeline atual de mídia aceita `PROJECT_DOC` e `PROJECT_RENDER`, mas só libera uploads em `StatusPendingPhotoProcessing` ou `StatusRejectedByOwner`.
- Listing type "Casa na Planta" (`OffPlanHouse`) já existe no catálogo; content-type `application/pdf` está permitido em `media_processing.limits.allowed_content_types`.
- Feature flag `media_processing.features.allow_owner_project_uploads` está true (pode ser usada para gating do novo fluxo).
- Nenhum script de migração deve ser criado; sem alteração de schema.
- Códigos e Swagger devem seguir o guia `docs/toq_server_go_guide.md` (uma função pública por arquivo, documentação em inglês, comentários mínimos e úteis).

## Fases (executar sequencialmente; paralelizar apenas onde indicado)

### Fase 1 — Alinhamento de regras de negócio (owner vs admin)
- Status alvo definido: para `OffPlanHouse`, o ciclo é `StatusDraft` → `StatusPendingAvailability` → `StatusPendingPlanLoading` → `StatusReady`, salvo se `media_processing.features.listing_approval_admin_review` estiver `true`, caso em que há passo intermediário `StatusPendingAdminReview` antes de `StatusReady` (sem sessão de fotos, sem aprovação de owner).
- Escopo definido: exclusivo para `OffPlanHouse`.
- Armazenamento: sem processamento/derivação; copiar `raw` → `processed` para manter padrão e permitir listagem/download homogêneos.
- Job de finalização: sim, registrar em `media_processing_jobs` e gerar ZIP mesmo sem transformação.
- Saída: decisões consolidadas acima (sem pendências).

### Fase 2 — DTOs e contratos HTTP (pode rodar em paralelo com Fase 3)
- Criar `internal/adapter/left/http/dto/listing_dto_project_media.go` com:
  - `RequestProjectUploadURLsRequest` (lista de arquivos; asset types restritos a `PROJECT_DOC`/`PROJECT_RENDER`).
  - `RequestProjectUploadURLsResponse` (reuso de instruções de upload).
  - `CompleteProjectMediaRequest` (somente `listingIdentityId`).
- Garantir documentação Swagger completa nos DTOs (descrição, exemplos, enums, binding tags).

### Fase 3 — Converters HTTP (paralelo à Fase 2)
- Criar `internal/adapter/left/http/handlers/listing_handlers/converters/project_media_converters.go` para mapear DTO ↔ inputs de domínio, reutilizando structs de media_processing.
- Validar/normalizar asset types para uppercase e filtrar somente `PROJECT_DOC`/`PROJECT_RENDER`.

### Fase 4 — Handlers HTTP (depende das Fases 2 e 3)
- Novas rotas públicas autenticadas em `internal/adapter/left/http/handlers/media_processing_handlers/`:
  - `request_project_upload_urls.go` com `POST /listings/project-media/uploads`.
  - `complete_project_media.go` com `POST /listings/project-media/complete`.
  - `delete_project_media.go` com `DELETE /listings/project-media` (remover asset de projeto por listingIdentityId/assetType/sequence, respeitando status permitido).
- Uso de `SendHTTPErrorObj`, nenhum span manual, documentação Swagger completa.

### Fase 5 — Serviço de mídia (depende das Fases 2 e 3)
- Em `internal/core/service/media_processing_service/` adicionar métodos:
  - `RequestProjectUploadURLs(ctx, input)`:
    - Checar feature flag `allow_owner_project_uploads`.
    - Validar listing: status permitido apenas `StatusPendingPlanLoading` e property type `OffPlanHouse`.
    - Reutilizar validação de manifest, porém com whitelist de asset types de projeto.
    - Gerar URLs via `storage.GenerateRawUploadURL`, `UpsertAsset` (mesmo fluxo existente).
  - `CompleteProjectMedia(ctx, input)`:
    - Checar feature flag.
    - Validar listing em `StatusPendingPlanLoading` (único estado permitido).
    - Copiar assets `raw` → `processed` (sem transformação) para preservar padrão de armazenamento e permitir download pelos endpoints atuais.
    - Registrar job de finalização em `media_processing_jobs` e gerar ZIP (reaproveitando workflow ou implementando etapa síncrona de zip) apenas para empacotar, não para converter.
    - Atualizar status final conforme flag: se `listing_approval_admin_review` true, avançar para `StatusPendingAdminReview`; caso contrário, `StatusReady` (sem owner approval para OffPlanHouse).
    - Garantir compatibilidade: `GET /listings/media` e `POST /listings/media/download` devem listar/servir PDF e JPEG desde que os assets estejam em `processed` (cópia feita no passo acima) e registrados no repositório com os tipos `PROJECT_DOC`/`PROJECT_RENDER`. Atualizar a documentação Swagger desses endpoints para explicitar enums/variantes de assetType incluindo os tipos de projeto, permitindo upload/listagem/download claros.
- Adicionar helpers privados para checagem de property type/flag/status e para operação de cópia `raw` → `processed` quando não houver workflow.

### Fase 6 — Roteamento e wiring (depende da Fase 4)
- Registrar novas rotas no router v2 mantendo prefixo `/api/v2` e documentá-las na tag Swagger `Listings Media`.
- Garantir injeção do mesmo service de mídia já usado pelos handlers existentes.

### Fase 7 — Observabilidade e documentação (depende das Fases 4 e 5)
- Confirmar logs mínimos (infra errors) e uso de `SetSpanError` no serviço.
- Rodar `make swagger` para atualizar JSON/YAML (gerado, não editado manualmente).

### Fase 8 — Revisão funcional e checklist (final)
- Verificar aderência ao guia: uma função pública por arquivo, Godoc em inglês, comentários úteis apenas.
- Validar que não há migrações nem alterações de schema.
- Confirmar que feature flags controlam o novo fluxo sem afetar o pipeline atual de fotos/vídeos.

## Interdependências entre times
- Produto/Negócio: decisão da Fase 1 (status alvo, escopo de property types) é bloqueio para Fase 5.
- Backend (DTO/Service/Handlers): pode paralelizar Fases 2 e 3; Handlers (Fase 4) aguardam DTOs e converters.
- Plataforma/DevOps: nenhuma mudança de infra requerida, apenas observabilidade; deve validar se `make swagger` é permitido no pipeline.

## Critérios de aceite
- Proprietário de anúncio em construção consegue solicitar URLs e concluir lote de projeto sem transitar por `PENDING_PHOTO_PROCESSING`.
- Uploads aceitam somente `PROJECT_DOC` (PDF) e `PROJECT_RENDER` (imagens) com validações de tamanho e checksum existentes.
- Status do listing avança conforme regra decidida na Fase 1, respeitando `listing_approval_admin_review`.
- Fluxo legado de fotos/vídeos permanece inalterado (gated por status original).
