# Plano de Refatoração e Implementação: Processamento de Mídia (Media Processing)

**Data:** 01/12/2025
**Status:** Planejado
**Objetivo:** Alinhar o módulo de processamento de mídia à Arquitetura Hexagonal estrita do projeto, segregar responsabilidades e implementar funcionalidades faltantes (`CompleteMedia`).

---

## 1. Diagnóstico e Motivação

A análise do código atual revelou os seguintes problemas críticos:
1.  **Violação de Arquitetura:** Handlers de mídia misturados em `listing_handlers`, violando a regra de espelhamento (Seção 2.1 do Guia).
2.  **Interface Ausente:** `ListingHandlerPort` não define os métodos de mídia, mas a implementação os possui.
3.  **Funcionalidade Incompleta:** O endpoint `CompleteMedia` é apenas um placeholder retornando 501 Not Implemented.
4.  **Acoplamento em Lambdas:** A lambda `validate` contém lógica de negócio no `main.go` e acoplamento indevido com eventos SQS.

---

## 2. Estrutura de Diretórios Alvo

```text
internal/
  core/
    port/
      left/
        http/
          mediaprocessinghandler/       # [NOVO] Port dedicado
            media_processing_handler_port.go
  adapter/
    left/
      http/
        handlers/
          media_processing_handlers/    # [NOVO] Adapter dedicado
            media_processing_handler.go
            request_upload_urls.go
            process_media.go
            complete_media.go
            ...
aws/
  lambdas/
    go_src/
      internal/
        core/
          service/
            validation/                 # [NOVO] Service para Lambda
```

---

## 3. Plano de Execução Faseado

### Fase 1: Fundação (Core Ports & Models)
**Objetivo:** Estabelecer os contratos corretos antes de mover a lógica.

1.  **Criar Port do Handler:**
    *   Arquivo: `internal/core/port/left/http/mediaprocessinghandler/media_processing_handler_port.go`
    *   Definir interface com todos os métodos (`RequestUploadURLs`, `ProcessMedia`, `CompleteMedia`, etc.).
2.  **Revisar DTOs:**
    *   Garantir que `internal/core/domain/dto` tenha `CompleteMediaInput` e `CompleteMediaOutput`.

### Fase 2: Lógica de Negócio (Service Layer)
**Objetivo:** Implementar a lógica faltante no Core.

1.  **Implementar `CompleteMedia`:**
    *   Arquivo: `internal/core/service/media_processing_service/complete_media.go`
    *   Lógica: Validar status do listing, disparar job de ZIP (se houver), atualizar status para `PENDING_OWNER_APPROVAL`.
2.  **Refatorar Service Interface:**
    *   Atualizar `MediaProcessingServiceInterface` para incluir `CompleteMedia`.

### Fase 3: Camada HTTP (Adapters & Handlers)
**Objetivo:** Segregar os handlers e conectar ao novo Port.

1.  **Criar Adapter Package:**
    *   Pasta: `internal/adapter/left/http/handlers/media_processing_handlers/`
    *   Factory: `media_processing_handler.go` (implementando `NewMediaProcessingHandler`).
2.  **Migrar Handlers Existentes:**
    *   Mover lógica de `listing_handlers/request_upload_urls_handler.go` para `media_processing_handlers/request_upload_urls.go`.
    *   Repetir para `process_media`, `list_download_urls`, `update_media`, `delete_media`, `handle_callback`.
    *   *Atenção:* Ajustar imports e nomes de pacotes.
3.  **Implementar Handler `CompleteMedia`:**
    *   Arquivo: `internal/adapter/left/http/handlers/media_processing_handlers/complete_media.go`
    *   Conectar ao service criado na Fase 2.
4.  **Wiring (Injeção de Dependência):**
    *   Atualizar `internal/core/factory/factory.go` (ou `main.go`) para usar o novo handler.
    *   Atualizar rotas no Gin (`internal/core/config/server/routes.go` ou similar).

### Fase 4: Infraestrutura AWS (Lambdas)
**Objetivo:** Corrigir arquitetura da Lambda `validate`.

1.  **Criar Service de Validação:**
    *   Arquivo: `aws/lambdas/go_src/internal/core/service/validation/service.go`
    *   Mover lógica de validação de assets para cá.
2.  **Refatorar Handler da Lambda:**
    *   Arquivo: `aws/lambdas/go_src/internal/adapter/left/lambda/validate/handler.go`
    *   Remover dependência de `events.SQSEvent`. Aceitar `StepFunctionPayload`.
3.  **Limpar `main.go`:**
    *   Deixar apenas a inicialização e chamada do `lambda.Start`.

### Fase 5: Limpeza e Finalização
**Objetivo:** Remover código morto e atualizar documentação.

1.  **Remover Código Antigo:**
    *   Apagar os arquivos de mídia de `internal/adapter/left/http/handlers/listing_handlers/`.
    *   Remover métodos de mídia de `ListingHandlerPort` e `ListingHandler`.
2.  **Swagger:**
    *   Rodar `make swagger` para refletir a nova estrutura (tags devem permanecer as mesmas para não quebrar clientes, mas a organização do código muda).

---

## 4. Code Skeletons (Templates)

### 4.1 Port (Fase 1)
```go
package mediaprocessinghandler

import "github.com/gin-gonic/gin"

type MediaProcessingHandlerPort interface {
    RequestUploadURLs(c *gin.Context)
    ProcessMedia(c *gin.Context)
    ListDownloadURLs(c *gin.Context)
    UpdateMedia(c *gin.Context)
    DeleteMedia(c *gin.Context)
    CompleteMedia(c *gin.Context)
    HandleProcessingCallback(c *gin.Context)
}
```

### 4.2 Service Implementation (Fase 2)
```go
// internal/core/service/media_processing_service/complete_media.go
func (s *mediaProcessingService) CompleteMedia(ctx context.Context, input dto.CompleteMediaInput) error {
    // 1. Tracer & Logger
    // 2. Transaction
    // 3. Get Listing & Validate Status (PENDING_PHOTO_PROCESSING)
    // 4. Trigger ZIP Job (Async)
    // 5. Update Listing Status -> PENDING_OWNER_APPROVAL
    // 6. Commit
    return nil
}
```

### 4.3 Handler Adapter (Fase 3)
```go
// internal/adapter/left/http/handlers/media_processing_handlers/media_processing_handler.go
type MediaProcessingHandler struct {
    service mediaprocessingservice.MediaProcessingServiceInterface
    logger  *slog.Logger
}

func NewMediaProcessingHandler(s mediaprocessingservice.MediaProcessingServiceInterface, l *slog.Logger) mediaprocessinghandlerport.MediaProcessingHandlerPort {
    return &MediaProcessingHandler{service: s, logger: l}
}
```
