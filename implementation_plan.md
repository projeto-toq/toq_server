# Plano de Implementação — Extensões Listings

## 1. Diagnóstico
- **Documentação**: Regras arquiteturais confirmadas em `docs/toq_server_go_guide.md`; esquema `users` avaliado em `scripts/db_creation.sql` (campo `created_at` disponível para cálculo de tempo de cadastro). Requisitos descritos em `PROMPT.md`.
- **GET /listings**: `ListListingsRequest`/handler/service/repositório apenas filtram/ordenam por `zip_code`, `city`, `neighborhood`. Faltam `street`, `number`, `complement`, `complex`, `state` e ordenação por campos de endereço.
- **POST /listings/detail**: Handler [`get_listing.go`](internal/adapter/left/http/handlers/listing_handlers/get_listing.go), serviço [`get_listing_detail.go`](internal/core/service/listing_service/get_listing_detail.go) e conversor [`listing_detail_to_dto.go`](internal/adapter/left/http/converters/listing_detail_to_dto.go) não retornam dados do proprietário nem placeholders de métricas do imóvel.
- **Dependências**: `listingService` já injeta `userRepository`, mas não consome `owner_metrics`. Handler de listings não recebe `userService`, inviabilizando `GetPhotoDownloadURL`. Domínio `UserInterface` não expõe `CreatedAt`, logo não há cálculo de “tempo de cadastro”.

## 2. Code Skeletons
```go
// internal/adapter/left/http/handlers/listing_handlers/list_listings_handler.go
sortBy, err := parseSortBy(req.SortBy) // aceitar novos campos
input := listingservices.ListListingsInput{
    Street:  strings.TrimSpace(req.Street),
    Number:  strings.TrimSpace(req.Number),
    // ...
}
```

```go
// internal/core/port/right/repository/listing_repository/listing_repository_interface.go
type ListListingsFilter struct {
    Street, Number, Complement, Complex, State string
    // existentes…
}
```

```go
// internal/adapter/right/mysql/listing/list_listings.go
if filter.Street != "" {
    conditions = append(conditions, "lv.street LIKE ?")
    args = append(args, filter.Street)
}
// idem para number/complement/complex/state
columnMap := map[string]string{
    "street": "lv.street",
    "number": "lv.number",
    // ...
}
```

```go
// internal/core/model/user_model/user_domain.go
type user struct {
    // ...
    createdAt time.Time
}
func (u *user) GetCreatedAt() time.Time { return u.createdAt }
func (u *user) SetCreatedAt(t time.Time) { u.createdAt = t }
```

```go
// internal/adapter/right/mysql/user/get_user_by_id.go
query := `SELECT u.id, ..., u.permanently_blocked, u.created_at, ...`
```

```go
// internal/core/service/listing_service/listing_service.go
type listingService struct {
    listingRepository listingrepository.ListingRepoPortInterface
    ownerMetricsRepo  ownermetricsrepository.Repository
    // ...
}
```

```go
// internal/core/service/listing_service/get_listing_detail.go
owner, err := ls.userRepository.GetUserByID(ctx, tx, identity.UserID)
metrics, _ := ls.ownerMetricsRepo.GetByOwnerID(ctx, tx, owner.GetID())
memberMonths := monthsSince(owner.GetCreatedAt(), time.Now().UTC())
output.OwnerDetail = &OwnerDetail{
    FullName: owner.GetFullName(),
    MemberSinceMonths: memberMonths,
    Metrics: metrics,
}
```

```go
// internal/adapter/left/http/handlers/listing_handlers/listing_handlers.go
type ListingHandler struct {
    listingService listingservice.ListingServiceInterface
    globalService  globalservice.GlobalServiceInterface
    userService    userservices.UserServiceInterface
}
```

```go
// internal/adapter/left/http/handlers/listing_handlers/get_listing.go
photoURL, err := lh.userService.GetPhotoDownloadURL(ctx, "medium")
if err == nil {
    detail.OwnerDetail.PhotoURL = photoURL
}
```

```go
// internal/adapter/left/http/dto/listing_dto.go
type ListingOwnerInfoResponse struct {
    FullName          string `json:"fullName"`
    PhotoURL          string `json:"photoUrl,omitempty"`
    MemberSinceMonths int    `json:"memberSinceMonths"`
    VisitAverageSeconds    *int64 `json:"visitAverageSeconds,omitempty"`
    ProposalAverageSeconds *int64 `json:"proposalAverageSeconds,omitempty"`
}

type ListingPerformanceMetricsResponse struct {
    Shares    int64 `json:"shares"`
    Views     int64 `json:"views"`
    Favorites int64 `json:"favorites"`
}
```

```go
// internal/adapter/left/http/converters/listing_detail_to_dto.go
if detail.OwnerDetail != nil {
    resp.OwnerInfo = &dto.ListingOwnerInfoResponse{ /* mapear campos */ }
}
resp.PerformanceMetrics = dto.ListingPerformanceMetricsResponse{
    Shares: detail.Performance.Shares,
    Views: detail.Performance.Views,
    Favorites: detail.Performance.Favorites,
    // TODO: service populates later
}
```

```go
// internal/core/service/listing_service/helpers.go
func monthsSince(start, now time.Time) int {
    // cálculo seguro evitando negativos
}
```

## 3. Estrutura de Diretórios Afetada
```
internal/adapter/left/http/dto/listing_dto.go
internal/adapter/left/http/handlers/listing_handlers/
internal/adapter/left/http/converters/listing_detail_to_dto.go
internal/adapter/right/mysql/listing/list_listings.go
internal/adapter/right/mysql/user/{get_*,scan_helpers.go,entities,converters}
internal/core/model/user_model/{user_domain.go,user_interface.go}
internal/core/service/listing_service/{listing_service.go,list_listings.go,get_listing_detail.go,helpers.go}
internal/core/port/right/repository/listing_repository/listing_repository_interface.go
internal/core/config/inject_dependencies.go
internal/core/factory/concrete_adapter_factory.go
```

## 4. Ordem de Execução
1. **Modelo de Usuário** — adicionar `CreatedAt` ao domínio/interface/entidades/conversores e incluir a coluna em todas as queries/scanners `GetUserBy*`.
2. **Injeção de Dependências** — atualizar `listingService` (novo campo `ownerMetricsRepo`), `NewListingService`, `InitListingHandler`, `InitListingHandlerAdapter` e factory para repassar `userService` e repositório.
3. **GET /listings** — estender DTO/form bindings, parsing e validações para campos de endereço + sortBy; propagar para `ListListingsInput`, port e adapter SQL (filtros + ORDER BY).
4. **POST /listings/detail (service)** — carregar dados do proprietário e métricas (`owner_metrics`), calcular `MemberSince`, criar helper `monthsSince`, preparar estrutura `OwnerDetail` e placeholders de métricas do imóvel.
5. **POST /listings/detail (handler/DTO)** — obter `photoUrl` via `userService`, atualizar converter/DTO para expor `ownerInfo` e `performanceMetrics` com `TODO` indicando preenchimento futuro pelo serviço.
6. **Validação/Docs** — revisar comentários Swagger/Godoc nos pontos alterados e executar `make test`/`make lint` (se aplicável) garantindo aderência ao guia.

## 5. Status de Execução
- [x] **Fase 1 — Modelo de Usuário**: `CreatedAt` propagado para domínio, interface, entidades SQL, conversores e queries `GetUserBy*`/scanner.
- [x] **Fase 2 — Injeção de Dependências**: `listingService` agora injeta `ownerMetricsRepo`; config valida repositório e factory/handler recebem `userService`.
- [x] **Fase 3 — GET /listings**: DTO, handler, serviço e MySQL adapter aceitam filtros/sort por street/number/complement/complex/state e zip/city/neighborhood ordenáveis.
- [x] **Fase 4 — Serviço POST /listings/detail**: Service enriquece response com `OwnerDetail` (nome, tempo de cadastro calculado via `monthsSince`, métricas de SLA) e inicializa placeholders de `PerformanceMetrics`.
- [x] **Fase 5 — Handler/DTO POST /listings/detail**: Handler consome `userService` para gerar `ownerPhoto` impersonando o proprietário, DTO/conversor expõem `ownerInfo` (incluindo métricas SLA em ponteiro) e `performanceMetrics` placeholders.
- [x] **Fase 6 — Validação/Docs**: Ajustados comentários Godoc/Swagger dos novos helpers e DTOs; `make lint` executado com sucesso garantindo conformidade das fases anteriores.
``