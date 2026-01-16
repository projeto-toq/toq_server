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
- ✅ Fase 1 concluída: status `StatusPendingPlanLoading` já presente em `listing_model/constants.go`; asset types `PROJECT_DOC`/`PROJECT_RENDER` e `application/pdf` já suportados; feature flags de upload/admin review disponíveis em `env`/`NewConfigFromEnvironment`. Nenhum ajuste de código necessário nesta fase.

### Fase 2 — DTOs e contratos HTTP (pode rodar em paralelo com Fase 3)
- Criar `internal/adapter/left/http/dto/listing_dto_project_media.go` com:
  - `RequestProjectUploadURLsRequest`: `listingIdentityId` obrigatório (min=1) e `files` (min=1). Cada `ProjectUploadFileRequest` com `assetType` (enums `PROJECT_DOC`/`PROJECT_RENDER`, uppercase), `sequence` (gt 0), `filename`, `contentType`, `bytes` (gt 0), `checksum`, `title`, `metadata` (map opcional). Bindings e examples em todos os campos.
  - `RequestProjectUploadURLsResponse`: reuso de `RequestUploadInstructionResponse` contendo `uploadUrl`, `method`, `headers`, `objectKey`, `assetType`, `sequence`, `title`, `listingIdentityId`, `uploadUrlTTLSeconds`.
  - `CompleteProjectMediaRequest`: apenas `listingIdentityId` obrigatório (min=1).
- Ajustar DTOs existentes em `listing_dto.go` para cobrir o fluxo de projeto:
  - `ListMediaRequest.assetType` documentar enums incluindo `PROJECT_DOC`, `PROJECT_RENDER`, `ZIP` e permitir `sequence` opcional.
  - `GenerateDownloadURLsRequest.requests.resolution`: permitir `zip` quando `assetType=ZIP` e `original` para `PROJECT_DOC`/`PROJECT_RENDER`; descrever que `zip` ignora sequence.
- Garantir documentação Swagger completa (descrição em inglês, exemplos, enums, binding tags) conforme guia.
- ✅ Fase 2 concluída: criado `listing_dto_project_media.go` com DTOs dedicados e documentação completa; `listing_dto.go` atualizado para aceitar `PROJECT_DOC`/`PROJECT_RENDER` e resolução `zip`, adicionando validação de enums em listagem e download; gofmt aplicado.

### Fase 3 — Converters HTTP (paralelo à Fase 2)
- Criar `internal/adapter/left/http/handlers/listing_handlers/converters/project_media_converters.go` para mapear DTO ↔ inputs de domínio, reutilizando structs de media_processing.
- Normalizar `assetType` para uppercase, rejeitar fora da whitelist `PROJECT_DOC`/`PROJECT_RENDER`.
- Reusar conversores já existentes para respostas (upload instructions, listagem, download) sem alterar shape atual.
- ✅ Fase 3 concluída: criado `project_media_converters.go` com normalização/whitelist de `PROJECT_DOC`/`PROJECT_RENDER`, mapeando para `RequestUploadURLsInput` e `CompleteMediaInput`. Mantidos converters de resposta existentes. gofmt e make lint executados com sucesso.

### Fase 4 — Handlers HTTP (depende das Fases 2 e 3)
- Novas rotas públicas autenticadas em `internal/adapter/left/http/handlers/media_processing_handlers/`:
  - `request_project_upload_urls.go` (`POST /listings/project-media/uploads`): bind JSON, converte via novo converter, chama `RequestProjectUploadURLs`, responde 200 com instruções.
  - `complete_project_media.go` (`POST /listings/project-media/complete`): bind JSON simples, chama `CompleteProjectMedia`, responde 200/204 conforme retorno.
  - `delete_project_media.go` (`DELETE /listings/project-media`): bind JSON com `listingIdentityId`, `assetType` (whitelist projeto), `sequence`, chama `DeleteMedia` (reuso), responde 204.
- Swagger: Tag `Listings Media`, exemplos claros de payload incluindo `PROJECT_DOC`/`PROJECT_RENDER`, mencionar status gate `StatusPendingPlanLoading` e restrição a `OffPlanHouse`.
- Uso obrigatório de `SendHTTPErrorObj`; não criar spans nos handlers.
- ✅ Fase 4 concluída: criados handlers `request_project_upload_urls.go`, `complete_project_media.go` e `delete_project_media.go` com Swagger, validação de whitelist e uso de converters; delete responde 204. Handlers reutilizam service atual (`RequestUploadURLs`, `CompleteMedia`, `DeleteMedia`) até a especialização na Fase 5. gofmt e make lint executados.

### Fase 5 — Serviço de mídia (depende das Fases 2 e 3)
- Em `internal/core/service/media_processing_service/` adicionar métodos (um arquivo por método):
  - `RequestProjectUploadURLs(ctx, input)`: checar flag `allow_owner_project_uploads`; validar listing `OffPlanHouse` em `StatusPendingPlanLoading`; whitelist `PROJECT_DOC`/`PROJECT_RENDER`; reusar validação de limites/tipos; gerar URLs via `storage.GenerateRawUploadURL`; `UpsertAsset` mantendo status `PENDING_UPLOAD`.
  - `CompleteProjectMedia(ctx, input)`: checar flag; revalidar listing em `StatusPendingPlanLoading`; listar assets de projeto; copiar `raw` → `processed` (sem transformação) via storage; marcar assets como `PROCESSED`; registrar job em `media_processing_jobs` e gerar ZIP via StepFunctions Finalization (mantém compatibilidade com `latestZipBundle`); atualizar status para `StatusPendingAdminReview` se `listing_approval_admin_review` true, senão `StatusReady`; registrar logs de infra + `SetSpanError` em falhas.
- Helpers privados: validação de property type/status/flag; operação de cópia `raw`→`processed` reutilizável.
- Ajustes em métodos existentes para compatibilidade:
  - `ListMedia` deve aceitar/filter `PROJECT_DOC`/`PROJECT_RENDER` e expor ZIP quando existir job.
  - `GenerateDownloadURLs` deve aceitar `PROJECT_DOC`/`PROJECT_RENDER` e `ZIP` (resolution `zip`), permitindo `original` para projeto mesmo que não processado; manter regras de status atual.
- Nenhuma mudança no fluxo legado de fotos/vídeos além da aceitação de novos enums.

### Riscos e mitigações (Fase 5)
- Path de processed: espelhar convenção do adapter S3 (`/{listingId}/processed/{mediaType}/{resolution}/{filename}`) implementando função local no service para calcular a key (mediaType segment + filename extraído de metadata/keys); evita depender de método não exportado e garante compatibilidade com download/listagem.
- ZIP/finalização: reutilizar StepFunctions Finalization existente, registrando job mesmo sem etapa de processamento; payload via `buildFinalizationInput` com assets já `PROCESSED`, preservando leitura de ZIP por `latestZipBundle`.

### Fase 6 — Roteamento e wiring (depende da Fase 4)
- ✅ Handlers de projeto agora chamam os métodos específicos `RequestProjectUploadURLs` e `CompleteProjectMedia`; rotas já registradas em `/listings/project-media/*` permanecem válidas.
- Garantir injeção do mesmo service de mídia já usado pelos handlers existentes (já provido pelo wiring atual).

### Fase 7 — Observabilidade e documentação (depende das Fases 4 e 5)
- ✅ `make swagger` executado (warnings de go list em raiz, geração concluída). Sem alterações manuais em swagger.json/yaml.
- Confirmar logs mínimos (infra errors) e uso de `SetSpanError` no serviço (já aplicado no fluxo de projeto).

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
