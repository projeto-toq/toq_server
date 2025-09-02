# Documentação da Cadeia de Middlewares

## 📋 Visão Geral

O TOQ Server implementa uma cadeia de middlewares seguindo as melhores práticas do Go e do framework Gin. A ordem de execução é crítica para o funcionamento correto do sistema.

## 🔄 Ordem de Execução dos Middlewares

### Middlewares Globais (Aplicados a todas as rotas)

```
Request → 1. RequestID → 2. Recovery → 3. StructuredLogging → 4. CORS → 5. Telemetry → [Route Specific] → Handler
```

#### 1. RequestIDMiddleware
- **Posição:** Primeiro middleware (crítico)
- **Responsabilidade:** Gera UUID único para cada requisição
- **Importância:** Base para tracing e correlação de logs
- **Implementação:** `middlewares.RequestIDMiddleware()`

#### 2. Recovery
- **Posição:** Segunda
- **Responsabilidade:** Captura panics e previne crash do servidor
- **Implementação:** `gin.Recovery()`
- **Comportamento:** Converte panics em respostas HTTP 500

#### 3. StructuredLoggingMiddleware
- **Posição:** Terceira
- **Responsabilidade:** Log estruturado JSON com separação stdout/stderr
- **Características:**
  - INFO/DEBUG → stdout
  - WARN/ERROR → stderr
  - Campos: request_id, method, path, status, duration, user_id
- **Implementação:** `middlewares.StructuredLoggingMiddleware()`

#### 4. CORSMiddleware
- **Posição:** Quarta
- **Responsabilidade:** Configuração de headers CORS
- **Implementação:** `middlewares.CORSMiddleware()`

#### 5. TelemetryMiddleware
- **Posição:** Quinta
- **Responsabilidade:** OpenTelemetry tracing
- **Implementação:** `middlewares.TelemetryMiddleware()`

### Middlewares Específicos (Aplicados apenas em rotas protegidas)

```
[Middlewares Globais] → 6. AuthMiddleware → 7. PermissionMiddleware → Handler
```

#### 6. AuthMiddleware
- **Aplicação:** Grupos de rotas autenticadas
- **Responsabilidade:** 
  - Validação JWT
  - Extração de informações do usuário
  - Tracking de atividade via ActivityTracker
- **Dependências:** ActivityTracker
- **Implementação:** `middlewares.AuthMiddleware(activityTracker)`

#### 7. PermissionMiddleware
- **Aplicação:** Grupos de rotas que requerem permissões específicas
- **Responsabilidade:**
  - Verificação de permissões HTTP
  - Validação de roles
  - Controle de acesso granular
- **Dependências:** PermissionService
- **Implementação:** `middlewares.PermissionMiddleware(permissionService)`

## 🛣️ Mapeamento de Rotas e Middlewares

### Rotas Públicas (Apenas middlewares globais)
```
/api/v2/auth/*
/healthz
/readyz
/swagger/*
```

### Rotas Protegidas (Middlewares globais + Auth + Permission)
```
/api/v2/user/*       → AuthMiddleware + PermissionMiddleware
/api/v2/agency/*     → AuthMiddleware + PermissionMiddleware
/api/v2/realtor/*    → AuthMiddleware + PermissionMiddleware
/api/v2/listings/*   → AuthMiddleware + PermissionMiddleware
/api/v2/visits/*     → AuthMiddleware + PermissionMiddleware
/api/v2/offers/*     → AuthMiddleware + PermissionMiddleware
/api/v2/realtors/*   → AuthMiddleware + PermissionMiddleware
/api/v2/owners/*     → AuthMiddleware + PermissionMiddleware
```

## 🔧 Injeção de Dependências

### Padrão Factory
```go
// Em routes/routes.go
func SetupRoutes(
    router *gin.Engine,
    handlers *factory.HTTPHandlers,
    activityTracker *goroutines.ActivityTracker,
    permissionService permissionservice.PermissionServiceInterface,
)
```

### Aplicação nos Grupos
```go
// Exemplo para rotas de usuário
user := router.Group("/user")
user.Use(middlewares.AuthMiddleware(activityTracker))
user.Use(middlewares.PermissionMiddleware(permissionService))
```

## 📊 Logging Estruturado

### Formato JSON
```json
{
  "time": "2025-08-30T10:00:00Z",
  "level": "INFO",
  "msg": "HTTP Request",
  "request_id": "uuid-here",
  "method": "POST",
  "path": "/api/v2/user/profile",
  "status": 200,
  "duration": "15ms",
  "size": 1024,
  "client_ip": "192.168.1.1",
  "user_agent": "Mozilla/5.0...",
  "user_id": 12345,
  "user_role": "Owner",
  "profile_complete": true
}
```

### Separação de Streams
- **stdout:** INFO, DEBUG (operações normais)
- **stderr:** WARN, ERROR (problemas e erros)

## 🔒 Contexto e Segurança

### Fluxo de Contexto
1. **RequestIDMiddleware** → Adiciona Request ID ao contexto Gin
2. **AuthMiddleware** → Adiciona informações do usuário ao contexto
3. **Handlers** → Utilizam Context Utils para acessar informações

### Context Utils
```go
// Extração de informações
userInfo, err := utils.GetUserInfoFromGinContext(c)
requestID := utils.GetRequestIDFromGinContext(c)

// Validação de autenticação
userInfo, err := utils.RequireUserInContext(ctx)
isAuth := utils.IsAuthenticatedContext(ctx)
```

## 🎯 Melhores Práticas Implementadas

### Go Best Practices
- ✅ Interfaces separadas das implementações
- ✅ Injeção de dependência via factory
- ✅ Error handling com utils/http_errors
- ✅ Structured logging com slog

### Google Go Style Guide
- ✅ Nomes descritivos de funções e variáveis
- ✅ Documentação adequada
- ✅ Tratamento consistente de erros
- ✅ Organização clara de pacotes

### Arquitetura Hexagonal
- ✅ Middlewares como adapters left-side
- ✅ Services como core business logic
- ✅ Separação clara de responsabilidades

## 🚨 Considerações Importantes

### Performance
- Middlewares são executados em ordem para cada requisição
- StructuredLoggingMiddleware adiciona overhead mínimo
- AuthMiddleware utiliza cache Redis via ActivityTracker

### Segurança
- Todas as rotas protegidas DEVEM passar por Auth + Permission
- Tokens JWT são validados a cada requisição
- Activity tracking para auditoria

### Monitoramento
- Request ID permite correlação entre logs
- OpenTelemetry para tracing distribuído
- Logs estruturados facilitam análise

## 📈 Exemplos de Uso

### Adicionando Nova Rota Protegida
```go
protected := router.Group("/new-endpoint")
protected.Use(middlewares.AuthMiddleware(activityTracker))
protected.Use(middlewares.PermissionMiddleware(permissionService))
{
    protected.GET("/data", handler.GetData)
}
```

### Adicionando Nova Rota Pública
```go
// Apenas middlewares globais são aplicados automaticamente
auth.POST("/new-public", handler.NewPublicEndpoint)
```

---

**Última atualização:** 30 de Agosto de 2025  
**Versão:** 2.0.0  
**Arquitetura:** Hexagonal com Gin Framework
