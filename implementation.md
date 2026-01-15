# Implementation Plan — Auth Hardening (Access Token Blocklist, Refresh/Device Binding, Admin Ops)

## Scope
- Introduce access-token blocklist (Redis) to invalidate JWT access tokens immediately after logout/revocation/rotation.
- Enforce access token type/JTI in auth middleware; reject refresh tokens used as access.
- Preserve device binding across refresh flows (DeviceID propagated).
- Add admin endpoints to list/revoke access-token JTIs (uses existing permission service for access control).
- No session count limits (multiple concurrent sessions allowed).
- JWT secret migration to env is **out of scope** (future security review).

## Files to touch/create (no further discovery needed)
- New: `internal/core/port/right/cache/token_blocklist_port.go`
- New: `internal/adapter/right/redis/token_blocklist/`:
  - `token_blocklist_adapter.go`
  - `add.go`, `exists.go`, `delete.go`, `list.go`, `count.go`
- Update: `internal/core/config/env_loader.go` (wire blocklist adapter) — adjust actual DI point used today.
- Update: `internal/core/factory/adapter_factory.go` (or equivalent) to expose blocklist.
- Update: `internal/adapter/left/http/middlewares/auth_middleware.go` (validate typ/jti, check blocklist).
- Update: `internal/adapter/left/http/handlers/auth_handlers/refresh_token.go` (require/validate X-Device-Id, inject into ctx).
- Update: `internal/core/service/user_service/refresh_token.go` (propagate deviceID from old session; optional blocklist previous access jti if available).
- Update: `internal/core/service/user_service/signout.go` (capture access jti from ctx/header and add to blocklist).
- Optional: `internal/core/service/user_service/device_context.go` (if helper needed to extract jti/device).
- New admin endpoints under `internal/adapter/left/http/handlers/admin_handlers/`:
  - `list_blocklisted_jtis_handler.go`
  - `add_blocklisted_jti_handler.go`
  - `delete_blocklisted_jti_handler.go`
- DTOs: `internal/adapter/left/http/dto/admin_blocklist_dto.go`
- Routes: `internal/adapter/left/http/routes/routes.go` (register admin routes)
- Permissions: ensure permission slug (e.g., `auth.blocklist.manage`) is enforced via permission_service in the admin handlers.
- Optional (recommended): `internal/core/port/right/repository/session_repository/session_repo_port.go` + adapter file `internal/adapter/right/mysql/session/get_active_session_by_token_jti.go` for session lookup by JTI.

## Ready-to-paste code skeletons

### Port: token blocklist (Redis)
**File:** `internal/core/port/right/cache/token_blocklist_port.go`
```go
package cacheport

import "context"

// TokenBlocklistPort defines operations to manage blocked access-token JTIs.
type TokenBlocklistPort interface {
   Add(ctx context.Context, jti string, ttlSeconds int64) error
   Exists(ctx context.Context, jti string) (bool, error)
   Delete(ctx context.Context, jti string) error
   List(ctx context.Context, page, pageSize int64) (items []BlocklistItem, err error)
   Count(ctx context.Context) (int64, error)
}

// BlocklistItem represents a blocked JTI entry.
type BlocklistItem struct {
   JTI       string
   ExpiresAt int64 // unix seconds
}
```

### Adapter: Redis blocklist
**Dir:** `internal/adapter/right/redis/token_blocklist/`

`token_blocklist_adapter.go`
```go
package tokenblocklist

import (
   "context"
   "fmt"
   "time"

   "github.com/redis/go-redis/v9"
   cacheport "github.com/projeto-toq/toq_server/internal/core/port/right/cache"
   "github.com/projeto-toq/toq_server/internal/core/utils"
)

const keyPrefix = "toq:blocklist:jti:"

type Adapter struct {
   client *redis.Client
}

func NewAdapter(client *redis.Client) cacheport.TokenBlocklistPort {
   return &Adapter{client: client}
}

func key(jti string) string { return keyPrefix + jti }
```

`add.go`
```go
package tokenblocklist

import (
   "context"
   "fmt"
   "time"

   "github.com/projeto-toq/toq_server/internal/core/utils"
)

func (a *Adapter) Add(ctx context.Context, jti string, ttlSeconds int64) error {
   ctx, end, _ := utils.GenerateTracer(ctx)
   defer end()
   if ttlSeconds <= 0 {
      return fmt.Errorf("ttl must be positive")
   }
   return a.client.Set(ctx, key(jti), "1", time.Duration(ttlSeconds)*time.Second).Err()
}
```

`exists.go`
```go
package tokenblocklist

import "context"

func (a *Adapter) Exists(ctx context.Context, jti string) (bool, error) {
   res, err := a.client.Exists(ctx, key(jti)).Result()
   if err != nil {
      return false, err
   }
   return res == 1, nil
}
```

`delete.go`
```go
package tokenblocklist

import "context"

func (a *Adapter) Delete(ctx context.Context, jti string) error {
   _, err := a.client.Del(ctx, key(jti)).Result()
   return err
}
```

`list.go`
```go
package tokenblocklist

import (
   "context"
   "strings"
   "time"

   cacheport "github.com/projeto-toq/toq_server/internal/core/port/right/cache"
)

// Uses SCAN with pagination window; page/pageSize are coarse (not strict ordering).
func (a *Adapter) List(ctx context.Context, page, pageSize int64) ([]cacheport.BlocklistItem, error) {
   if page < 1 { page = 1 }
   if pageSize <= 0 { pageSize = 100 }
   start := (page - 1) * pageSize
   cursor := uint64(0)
   items := make([]cacheport.BlocklistItem, 0, pageSize)
   skipped := int64(0)
   for {
      keys, next, err := a.client.Scan(ctx, cursor, keyPrefix+"*", pageSize*2).Result()
      if err != nil { return nil, err }
      for _, k := range keys {
         if skipped < start { skipped++; continue }
         ttl, _ := a.client.TTL(ctx, k).Result()
         parts := strings.Split(k, ":")
         jti := parts[len(parts)-1]
         exp := time.Now().Add(ttl).Unix()
         items = append(items, cacheport.BlocklistItem{JTI: jti, ExpiresAt: exp})
         if int64(len(items)) >= pageSize {
            return items, nil
         }
      }
      if next == 0 { break }
      cursor = next
   }
   return items, nil
}
```

`count.go`
```go
package tokenblocklist

import "context"

func (a *Adapter) Count(ctx context.Context) (int64, error) {
   var cursor uint64
   var total int64
   for {
      keys, next, err := a.client.Scan(ctx, cursor, keyPrefix+"*", 1000).Result()
      if err != nil { return 0, err }
      total += int64(len(keys))
      if next == 0 { break }
      cursor = next
   }
   return total, nil
}
```

### Middleware auth (enforce typ/jti + blocklist)
**File:** `internal/adapter/left/http/middlewares/auth_middleware.go`
- Parse JWT with `jwt.ParseWithClaims` ensuring `SigningMethodHMAC`.
- Validate claim `typ == "access"`; if missing/wrong → 401.
- Require `jti` string; extract.
- Call blocklist `Exists`; if true → 401 "Token revoked".
- Keep existing `infos` extraction logic.

Snippet inside `validateAccessToken`:
```go
if typ, ok := (*claims)["typ"].(string); !ok || typ != "access" {
   return usermodel.UserInfos{}, utils.AuthenticationError("Invalid access token type")
}
jtiRaw, ok := (*claims)["jti"].(string)
if !ok || jtiRaw == "" {
   return usermodel.UserInfos{}, utils.AuthenticationError("Invalid access token")
}
if blocklist != nil {
   blocked, err := blocklist.Exists(context.Background(), jtiRaw)
   if err == nil && blocked {
      return usermodel.UserInfos{}, utils.AuthenticationError("Token revoked")
   }
}
```
(Inject blocklist adapter in middleware constructor; wire via routes/bootstrap.)

### Refresh handler
**File:** `internal/adapter/left/http/handlers/auth_handlers/refresh_token.go`
- Read header `X-Device-Id`; if present, validate UUIDv4; if absent, recommended: return 400 to enforce device binding.
- Inject into ctx: `ctx = context.WithValue(ctx, globalmodel.DeviceIDKey, trimmedDeviceID)`
- Call `userService.RefreshTokens(ctx, request.RefreshToken)`.

### Refresh service
**File:** `internal/core/service/user_service/refresh_token.go`
- After loading session, before `CreateTokens`, set:
```go
if did := session.GetDeviceID(); did != "" {
   ctx = context.WithValue(ctx, globalmodel.DeviceIDKey, did)
}
```
- Optional: if access `jti` is available from context/header, add to blocklist on rotation (reuse helper from logout).

### Logout service
**File:** `internal/core/service/user_service/signout.go`
- Extract current access `jti` from ctx/header (parse Authorization if needed).
- Add `jti` to blocklist with TTL = `globalmodel.GetAccessTTL().Seconds()`.
- Keep existing revoke sessions + device token removal.

Helper (suggested): `extractAccessJTIFromContext(ctx)` placed in `user_service` private helper file to reuse in logout/refresh.

### Admin DTOs
**File:** `internal/adapter/left/http/dto/admin_blocklist_dto.go`
```go
package dto

// BlocklistItemResponse represents one blocked JTI
type BlocklistItemResponse struct {
   JTI       string `json:"jti"`
   ExpiresAt int64  `json:"expiresAt"` // unix seconds
}

type ListBlocklistResponse struct {
   Items []BlocklistItemResponse `json:"items"`
   Total int64                   `json:"total"`
}

type AddBlocklistRequest struct {
   JTI string `json:"jti" binding:"required"`
   TTL int64  `json:"ttlSeconds" binding:"omitempty,gt=0"`
}
```

### Admin handlers (one public method per file)
Dir: `internal/adapter/left/http/handlers/admin_handlers/`

`list_blocklisted_jtis_handler.go`
```go
// @Summary List blocked JTIs
// @Tags AdminAuth
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page (1-based)"
// @Param pageSize query int false "Page size"
// @Success 200 {object} dto.ListBlocklistResponse
// @Router /admin/auth/blocklist [get]
func (h *AdminAuthHandler) ListBlocklistedJTIs(c *gin.Context) { /* parse pagination, call service/adapter, return JSON */ }
```

`add_blocklisted_jti_handler.go`
```go
// @Summary Add JTI to blocklist
// @Tags AdminAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.AddBlocklistRequest true "JTI and optional TTL"
// @Success 204
// @Router /admin/auth/blocklist [post]
func (h *AdminAuthHandler) AddBlocklistedJTI(c *gin.Context) { /* bind, default ttl = accessTTL, call blocklist.Add */ }
```

`delete_blocklisted_jti_handler.go`
```go
// @Summary Remove JTI from blocklist
// @Tags AdminAuth
// @Security BearerAuth
// @Produce json
// @Param jti path string true "JTI"
// @Success 204
// @Router /admin/auth/blocklist/{jti} [delete]
func (h *AdminAuthHandler) DeleteBlocklistedJTI(c *gin.Context) { /* call blocklist.Delete */ }
```

**Permission enforcement:** use existing permission_service; require slug (e.g., `auth.blocklist.manage`) before executing.

### Routes
**File:** `internal/adapter/left/http/routes/routes.go`
- Register admin routes under admin group with middleware auth + permission check:
```go
adminAuth := admin.Group("/auth")
adminAuth.GET("/blocklist", adminAuthHandler.ListBlocklistedJTIs)
adminAuth.POST("/blocklist", adminAuthHandler.AddBlocklistedJTI)
adminAuth.DELETE("/blocklist/:jti", adminAuthHandler.DeleteBlocklistedJTI)
```

### Session repo by JTI (optional)
**Port:** add method in `session_repo_port.go`:
```go
GetActiveSessionByTokenJTI(ctx context.Context, tx *sql.Tx, jti string) (sessionmodel.SessionInterface, error)
```
**Adapter file:** `internal/adapter/right/mysql/session/get_active_session_by_token_jti.go`
Query: `SELECT ... FROM sessions WHERE token_jti = ? AND revoked = false AND expires_at > UTC_TIMESTAMP()`; map with existing mapper; return `sql.ErrNoRows` when absent.

### Factories/DI
- Inject Redis client already built for cache; pass into `NewAdapter` of blocklist; expose via `AdapterFactory` and into middleware/handlers.
- Ensure lifecycle closes Redis client as today.

### Observability
- Add metrics counters: blocklist_hits, blocklist_miss, blocklist_add, blocklist_delete.
- Log keys: `auth.blocklist.hit`, `auth.blocklist.add`, `auth.blocklist.delete`, `auth.blocklist.error`.

## Execution checklist
- [x] Create port file.
- [x] Create Redis adapter files.
- [x] Wire DI (factory/config) and make blocklist available to middleware/handlers.
- [x] Update auth middleware for typ/jti + blocklist.
- [x] Update refresh handler + service for deviceID propagation.
- [x] Update signout to blocklist current access jti.
- [x] Add admin DTOs/handlers/routes with permission guard.
- [ ] (Optional) Add session repo by JTI + index request to DBA (not implemented in this iteration).
- [x] Run `go fmt` and `go vet`/lint.

## Components & Responsibilities
- **Port:** `TokenBlocklistPort` (right/redis): `Add(ctx, jti, ttl)`, `Exists(ctx, jti)`, `Delete(ctx, jti)`, `List(ctx, page, pageSize)` (for admin list), `Count(ctx)` (optional for paging total).
- **Adapter (Redis):** Key `toq:blocklist:jti:{jti}`; value optional (e.g., timestamp); TTL = access TTL. Leverage existing Redis client pattern from `internal/core/cache/redis_cache.go` (same instrumentation).
- **Middleware auth:**
  - Parse JWT HS256, validate `typ == "access"`, require `jti`.
  - Check blocklist `Exists`; deny if present.
  - (Optional) also check session-by-JTI if implemented.
- **Refresh handler/service:**
  - Require/validate `X-Device-Id` (UUIDv4 recommended); inject into ctx.
  - When rotating, propagate prior session `device_id` into ctx before `CreateTokens` to avoid losing device binding.
  - (Optional) blocklist previous access `jti` when rotating (if available from request).
- **Logout service:**
  - Blocklist current access `jti` (from Authorization header) with TTL=access TTL.
  - Keep existing session/refresh revocation and device-token pruning.
- **Session repo (optional support):** Add `GetActiveSessionByTokenJTI(ctx, tx, jti)`; index recommendation below.
- **Admin endpoints:**
  - **List blocklisted JTIs:** `GET /admin/auth/blocklist` (paginate, show jti, expiresAt).
  - **Revoke/Add JTI:** `POST /admin/auth/blocklist` (body: jti, ttl optional -> default access TTL).
  - **Delete JTI:** `DELETE /admin/auth/blocklist/{jti}`.
  - Protected via existing permission system (require admin role/permission key to be defined with permission_service).

## DBA Notes (no migrations in this task)
- Add composite index to `sessions(token_jti, revoked, expires_at)` to support `GetActiveSessionByTokenJTI` (filter: `token_jti = ? AND revoked = false AND expires_at > UTC_TIMESTAMP()`).
- Ensure `refresh_hash` remains indexed (unique/BTREE) — already used today.
- No DB changes for Redis blocklist.

## Sequence (phased-friendly)
1. **Ports/Adapters**
   - Add `TokenBlocklistPort` interface (right/redis).
   - Implement Redis adapter with `Add/Exists/Delete/List/Count`; reuse Redis client style from `redis_cache.go` (RESP2, tracing/metrics).
   - Wire in factory/DI; expose via `globalService` if shared helpers needed.

2. **Middleware**
   - Enforce `typ == access`, require `jti`.
   - Call blocklist `Exists`; on hit => 401.
   - Maintain current claims extraction for `infos`.

3. **Refresh Flow**
   - Handler: validate `X-Device-Id` (UUIDv4); propagate into ctx.
   - Service: before `CreateTokens`, set `DeviceIDKey` from previous session; optional blocklist of previous access `jti` (if provided).

4. **Logout Flow**
   - Extract current access `jti` from Authorization header; add to blocklist (TTL=access TTL).
   - Keep revocation of session/refresh and device-token removal.

5. **Session Repo (optional but recommended)**
   - Add `GetActiveSessionByTokenJTI`; use new index.
   - Use in middleware as secondary check (after blocklist) if desired.

6. **Admin Endpoints**
   - Add handler(s) under admin routes (one file per method per guide): list, add, delete JTI.
   - DTOs for request/response with Swagger annotations.
   - Permission guard via existing permission_service; define required permission slug (e.g., `auth.blocklist.manage`).

7. **Observability**
   - Metrics: counters for blocklist hits/misses/adds/deletes; histogram optional.
   - Logs: `auth.blocklist.hit`, `auth.blocklist.add`, `auth.blocklist.delete`.

8. **Testing/Validation (manual for now)**
   - Signin, get access, logout → access must be denied immediately after blocklist add.
   - Refresh with device ID → new session retains device_id; logout should prune tokens by device.
   - Admin list/add/delete flows; permission enforcement.

## Notes/Constraints
- Keep one public method per file (per project guide).
- Comments/Godoc in English; no Swagger JSON edit—use annotations.
- No JWT secret change to env in this iteration.
- Multiple sessions per user remain allowed.

