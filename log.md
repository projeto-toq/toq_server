Starting: /home/toq_admin/go/bin/dlv dap --listen=127.0.0.1:39853 --log-dest=3 from /codigos/go_code/toq_server/cmd
DAP server listening at: 127.0.0.1:39853
Type 'dlv help' for list of commands.
2025/08/29 20:14:48 INFO 🚀 Iniciando TOQ Server Bootstrap version=2.0.0 timestamp=2025-08-29T20:14:48Z
2025/08/29 20:14:48 INFO ▶️ Executando fase phase=Phase01_InitializeContext timestamp=2025-08-29T20:14:48Z
2025/08/29 20:14:48 INFO 🎯 FASE 1: Inicialização de Contexto e Sinais
2025/08/29 20:14:48 INFO ✅ Contexto e sinais inicializados com sucesso
2025/08/29 20:14:48 INFO 🔍 Iniciando servidor pprof na porta 6060
2025/08/29 20:14:48 INFO ✅ Fase concluída phase=Phase01_InitializeContext duration=790.072µs
2025/08/29 20:14:48 INFO ▶️ Executando fase phase=Phase02_LoadConfiguration timestamp=2025-08-29T20:14:48Z
2025/08/29 20:14:48 INFO ✅ Servidor pprof iniciado em localhost:6060
2025/08/29 20:14:48 INFO 🎯 FASE 2: Carregamento e Validação de Configuração
2025/08/29 20:14:48 INFO Configuration loaded successfully from YAML path=configs/env.yaml
time=2025-08-29T20:14:48.602Z level=INFO msg="Logging system initialized" level=INFO to_file=false add_source=false
time=2025-08-29T20:14:48.602Z level=INFO msg="INFO ✅ Logging inicial baseado em ENV configurado"
time=2025-08-29T20:14:48.603Z level=INFO msg="Logging system initialized" level=INFO to_file=false add_source=false
time=2025-08-29T20:14:48.603Z level=INFO msg="INFO ✅ Logging reconfigurado com prioridade ENV > YAML > defaults"
time=2025-08-29T20:14:48.603Z level=INFO msg="INFO ✅ Configuração carregada e validada com sucesso version=2.0.0"
time=2025-08-29T20:14:48.603Z level=INFO msg="INFO ✅ Fase concluída phase=Phase02_LoadConfiguration duration=2.651001ms"
time=2025-08-29T20:14:48.603Z level=INFO msg="INFO ▶️ Executando fase phase=Phase03_InitializeInfrastructure timestamp=2025-08-29T20:14:48Z"
time=2025-08-29T20:14:48.603Z level=INFO msg="INFO 🎯 FASE 3: Inicialização da Infraestrutura Core"
time=2025-08-29T20:14:48.606Z level=INFO msg="Database connection initialized" uri="toq_user:toq_password@tcp(localhost:3306)/toq_db?parseTime=true&loc=UTC&timeout=30s&readTimeout=30s&writeTimeout=30s"
time=2025-08-29T20:14:48.606Z level=INFO msg="INFO ✅ Conexão com banco de dados estabelecida"
time=2025-08-29T20:14:48.608Z level=INFO msg="Redis cache connected successfully" url=redis://localhost:6379/0
time=2025-08-29T20:14:48.608Z level=INFO msg="INFO ✅ Sistema de cache Redis inicializado com sucesso"
time=2025-08-29T20:14:48.608Z level=INFO msg="OpenTelemetry initialization placeholder - not fully implemented" enabled=true otlp_enabled=true endpoint=http://localhost:14318
time=2025-08-29T20:14:48.608Z level=INFO msg="INFO ✅ OpenTelemetry inicializado (tracing + metrics)"
time=2025-08-29T20:14:48.608Z level=INFO msg="Activity tracker initialization deferred to Phase 07"
time=2025-08-29T20:14:48.608Z level=INFO msg="INFO ✅ Activity tracker inicializado"
time=2025-08-29T20:14:48.608Z level=INFO msg="INFO ✅ Infraestrutura core inicializada com sucesso"
time=2025-08-29T20:14:48.608Z level=INFO msg="INFO ✅ Fase concluída phase=Phase03_InitializeInfrastructure duration=5.306921ms"
time=2025-08-29T20:14:48.608Z level=INFO msg="INFO ▶️ Executando fase phase=Phase04_InjectDependencies timestamp=2025-08-29T20:14:48Z"
time=2025-08-29T20:14:48.608Z level=INFO msg="INFO 🎯 FASE 4: Injeção de Dependências via Factory Pattern"
time=2025-08-29T20:14:48.608Z level=INFO msg="Starting dependency injection using Factory Pattern"
time=2025-08-29T20:14:48.608Z level=INFO msg="DEBUG: InjectDependencies method called on config instance"
time=2025-08-29T20:14:48.608Z level=INFO msg="Creating storage adapters"
time=2025-08-29T20:14:48.608Z level=INFO msg="Creating storage adapters"
time=2025-08-29T20:14:48.611Z level=INFO msg="Redis cache connected successfully" url=redis://localhost:6379/0
time=2025-08-29T20:14:48.611Z level=INFO msg="Successfully created all storage adapters"
time=2025-08-29T20:14:48.611Z level=INFO msg="Creating repository adapters"
time=2025-08-29T20:14:48.611Z level=INFO msg="Creating repository adapters"
time=2025-08-29T20:14:48.611Z level=INFO msg="Successfully created all repository adapters"
time=2025-08-29T20:14:48.611Z level=INFO msg="Assigning repository adapters"
time=2025-08-29T20:14:48.611Z level=INFO msg="Repository adapters assigned successfully"
time=2025-08-29T20:14:48.611Z level=INFO msg="Creating validation adapters"
time=2025-08-29T20:14:48.611Z level=INFO msg="Creating validation adapters"
time=2025-08-29T20:14:48.611Z level=INFO msg="Successfully created all validation adapters"
time=2025-08-29T20:14:48.611Z level=INFO msg="Creating external service adapters"
time=2025-08-29T20:14:48.611Z level=INFO msg="Creating external service adapters"
time=2025-08-29T20:14:48.614Z level=INFO msg="Creating S3 adapter" region=us-east-1 bucket=toq-app-media
time=2025-08-29T20:14:48.614Z level=INFO msg="S3 adapter created successfully" bucket=toq-app-media region=us-east-1
time=2025-08-29T20:14:48.614Z level=INFO msg="Successfully created all external service adapters"
time=2025-08-29T20:14:48.614Z level=INFO msg="Initializing all services"
time=2025-08-29T20:14:48.615Z level=INFO msg="All services initialized successfully"
time=2025-08-29T20:14:48.615Z level=INFO msg="Dependency injection completed successfully using Factory Pattern"
time=2025-08-29T20:14:48.615Z level=INFO msg="INFO ✅ Injeção de dependências concluída via Factory Pattern"
time=2025-08-29T20:14:48.617Z level=INFO msg="INFO ✅ Fase concluída phase=Phase04_InjectDependencies duration=8.473049ms"
time=2025-08-29T20:14:48.617Z level=INFO msg="INFO ▶️ Executando fase phase=Phase05_InitializeServices timestamp=2025-08-29T20:14:48Z"
time=2025-08-29T20:14:48.617Z level=INFO msg="INFO 🎯 FASE 5: Inicialização de Serviços"
time=2025-08-29T20:14:48.617Z level=INFO msg="INFO ✅ Serviço inicializado service=GlobalService"
time=2025-08-29T20:14:48.617Z level=INFO msg="INFO ✅ Serviço inicializado service=PermissionService"
time=2025-08-29T20:14:48.618Z level=INFO msg="INFO ✅ Serviço inicializado service=UserService"
time=2025-08-29T20:14:48.618Z level=INFO msg="INFO ✅ Serviço inicializado service=ComplexService"
time=2025-08-29T20:14:48.618Z level=INFO msg="INFO ✅ Serviço inicializado service=ListingService"
time=2025-08-29T20:14:48.618Z level=INFO msg="INFO ✅ Todos os serviços inicializados com sucesso"
time=2025-08-29T20:14:48.619Z level=INFO msg="INFO ✅ Fase concluída phase=Phase05_InitializeServices duration=2.584769ms"
time=2025-08-29T20:14:48.620Z level=INFO msg="INFO ▶️ Executando fase phase=Phase06_ConfigureHandlers timestamp=2025-08-29T20:14:48Z"
time=2025-08-29T20:14:48.621Z level=INFO msg="INFO 🎯 FASE 6: Configuração de Handlers e Rotas"
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)
time=2025-08-29T20:14:48.621Z level=INFO msg="HTTP server initialized" port=:8080 read_timeout=30s write_timeout=30s
time=2025-08-29T20:14:48.621Z level=INFO msg="HTTP server initialization completed"
time=2025-08-29T20:14:48.622Z level=INFO msg="INFO ✅ Servidor HTTP configurado com TLS e middleware"
time=2025-08-29T20:14:48.622Z level=INFO msg="INFO ✅ Handlers HTTP preparados para criação"
time=2025-08-29T20:14:48.622Z level=INFO msg="Creating HTTP handlers"
time=2025-08-29T20:14:48.622Z level=INFO msg="Successfully created all HTTP handlers"
time=2025-08-29T20:14:48.622Z level=INFO msg="✅ HTTP handlers created successfully via factory"
[GIN-debug] GET    /healthz                  --> github.com/giulio-alfieri/toq_server/internal/core/config.(*config).setupBasicRoutes.func1 (1 handlers)
[GIN-debug] GET    /readyz                   --> github.com/giulio-alfieri/toq_server/internal/core/config.(*config).setupBasicRoutes.func2 (1 handlers)
[GIN-debug] GET    /api/v1/ping              --> github.com/giulio-alfieri/toq_server/internal/core/config.(*config).setupBasicRoutes.func3 (1 handlers)
[GIN-debug] GET    /swagger/*any             --> github.com/swaggo/gin-swagger.CustomWrapHandler.func1 (5 handlers)
[GIN-debug] POST   /api/v1/auth/owner        --> github.com/giulio-alfieri/toq_server/internal/core/port/left/http/authhandler.AuthHandlerPort.CreateOwner-fm (5 handlers)
[GIN-debug] POST   /api/v1/auth/realtor      --> github.com/giulio-alfieri/toq_server/internal/core/port/left/http/authhandler.AuthHandlerPort.CreateRealtor-fm (5 handlers)
[GIN-debug] POST   /api/v1/auth/agency       --> github.com/giulio-alfieri/toq_server/internal/core/port/left/http/authhandler.AuthHandlerPort.CreateAgency-fm (5 handlers)
[GIN-debug] POST   /api/v1/auth/signin       --> github.com/giulio-alfieri/toq_server/internal/core/port/left/http/authhandler.AuthHandlerPort.SignIn-fm (5 handlers)
[GIN-debug] POST   /api/v1/auth/refresh      --> github.com/giulio-alfieri/toq_server/internal/core/port/left/http/authhandler.AuthHandlerPort.RefreshToken-fm (5 handlers)
[GIN-debug] POST   /api/v1/auth/password/request --> github.com/giulio-alfieri/toq_server/internal/core/port/left/http/authhandler.AuthHandlerPort.RequestPasswordChange-fm (5 handlers)
[GIN-debug] POST   /api/v1/auth/password/confirm --> github.com/giulio-alfieri/toq_server/internal/core/port/left/http/authhandler.AuthHandlerPort.ConfirmPasswordChange-fm (5 handlers)
[GIN-debug] GET    /api/v1/user/profile      --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func1 (5 handlers)
[GIN-debug] PUT    /api/v1/user/profile      --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func2 (5 handlers)
[GIN-debug] DELETE /api/v1/user/account      --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func3 (5 handlers)
[GIN-debug] GET    /api/v1/user/onboarding   --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func4 (5 handlers)
[GIN-debug] GET    /api/v1/user/roles        --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func5 (5 handlers)
[GIN-debug] GET    /api/v1/user/home         --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func6 (5 handlers)
[GIN-debug] PUT    /api/v1/user/opt-status   --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func7 (5 handlers)
[GIN-debug] POST   /api/v1/user/signout      --> github.com/giulio-alfieri/toq_server/internal/core/port/left/http/userhandler.UserHandlerPort.SignOut-fm (5 handlers)
[GIN-debug] POST   /api/v1/user/photo/upload-url --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func8 (5 handlers)
[GIN-debug] GET    /api/v1/user/profile/thumbnails --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func9 (5 handlers)
[GIN-debug] POST   /api/v1/user/email/request --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func10 (5 handlers)
[GIN-debug] POST   /api/v1/user/email/confirm --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func11 (5 handlers)
[GIN-debug] POST   /api/v1/user/email/resend --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func12 (5 handlers)
[GIN-debug] POST   /api/v1/user/phone/request --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func13 (5 handlers)
[GIN-debug] POST   /api/v1/user/phone/confirm --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func14 (5 handlers)
[GIN-debug] POST   /api/v1/user/phone/resend --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func15 (5 handlers)
[GIN-debug] POST   /api/v1/user/role/alternative --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func16 (5 handlers)
[GIN-debug] POST   /api/v1/user/role/switch  --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func17 (5 handlers)
[GIN-debug] POST   /api/v1/agency/invite-realtor --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func18 (5 handlers)
[GIN-debug] GET    /api/v1/agency/realtors   --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func19 (5 handlers)
[GIN-debug] GET    /api/v1/agency/realtors/:id --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func20 (5 handlers)
[GIN-debug] DELETE /api/v1/agency/realtors/:id --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func21 (5 handlers)
[GIN-debug] POST   /api/v1/realtor/creci/verify --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func22 (5 handlers)
[GIN-debug] POST   /api/v1/realtor/creci/upload-url --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func23 (5 handlers)
[GIN-debug] POST   /api/v1/realtor/invitation/accept --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func24 (5 handlers)
[GIN-debug] POST   /api/v1/realtor/invitation/reject --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func25 (5 handlers)
[GIN-debug] GET    /api/v1/realtor/agency    --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func26 (5 handlers)
[GIN-debug] DELETE /api/v1/realtor/agency    --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func27 (5 handlers)
[GIN-debug] GET    /api/v1/listings          --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func1 (5 handlers)
[GIN-debug] POST   /api/v1/listings          --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func2 (5 handlers)
[GIN-debug] GET    /api/v1/listings/search   --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func3 (5 handlers)
[GIN-debug] GET    /api/v1/listings/options  --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func4 (5 handlers)
[GIN-debug] GET    /api/v1/listings/features/base --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func5 (5 handlers)
[GIN-debug] GET    /api/v1/listings/favorites --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func6 (5 handlers)
[GIN-debug] GET    /api/v1/listings/:id      --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func7 (5 handlers)
[GIN-debug] PUT    /api/v1/listings/:id      --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func8 (5 handlers)
[GIN-debug] DELETE /api/v1/listings/:id      --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func9 (5 handlers)
[GIN-debug] POST   /api/v1/listings/:id/end-update --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func10 (5 handlers)
[GIN-debug] GET    /api/v1/listings/:id/status --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func11 (5 handlers)
[GIN-debug] POST   /api/v1/listings/:id/approve --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func12 (5 handlers)
[GIN-debug] POST   /api/v1/listings/:id/reject --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func13 (5 handlers)
[GIN-debug] POST   /api/v1/listings/:id/suspend --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func14 (5 handlers)
[GIN-debug] POST   /api/v1/listings/:id/release --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func15 (5 handlers)
[GIN-debug] POST   /api/v1/listings/:id/copy --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func16 (5 handlers)
[GIN-debug] POST   /api/v1/listings/:id/share --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func17 (5 handlers)
[GIN-debug] POST   /api/v1/listings/:id/favorite --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func18 (5 handlers)
[GIN-debug] DELETE /api/v1/listings/:id/favorite --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func19 (5 handlers)
[GIN-debug] POST   /api/v1/listings/:id/visit/request --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func20 (5 handlers)
[GIN-debug] GET    /api/v1/listings/:id/visits --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func21 (5 handlers)
[GIN-debug] POST   /api/v1/listings/:id/offers --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func22 (5 handlers)
[GIN-debug] GET    /api/v1/listings/:id/offers --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func23 (5 handlers)
[GIN-debug] GET    /api/v1/visits            --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func24 (5 handlers)
[GIN-debug] DELETE /api/v1/visits/:id        --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func25 (5 handlers)
[GIN-debug] POST   /api/v1/visits/:id/confirm --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func26 (5 handlers)
[GIN-debug] POST   /api/v1/visits/:id/approve --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func27 (5 handlers)
[GIN-debug] POST   /api/v1/visits/:id/reject --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func28 (5 handlers)
[GIN-debug] GET    /api/v1/offers            --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func29 (5 handlers)
[GIN-debug] PUT    /api/v1/offers/:id        --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func30 (5 handlers)
[GIN-debug] DELETE /api/v1/offers/:id        --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func31 (5 handlers)
[GIN-debug] POST   /api/v1/offers/:id/send   --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func32 (5 handlers)
[GIN-debug] POST   /api/v1/offers/:id/approve --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func33 (5 handlers)
[GIN-debug] POST   /api/v1/offers/:id/reject --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func34 (5 handlers)
[GIN-debug] POST   /api/v1/realtors/:id/evaluate --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func35 (5 handlers)
[GIN-debug] POST   /api/v1/owners/:id/evaluate --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func36 (5 handlers)
time=2025-08-29T20:14:48.624Z level=INFO msg="HTTP handlers and routes configured successfully"
time=2025-08-29T20:14:48.624Z level=INFO msg="INFO ✅ Rotas e middlewares configurados"
time=2025-08-29T20:14:48.624Z level=INFO msg="INFO ✅ Health checks configurados"
time=2025-08-29T20:14:48.624Z level=INFO msg="INFO ✅ Handlers e rotas configurados com sucesso"
time=2025-08-29T20:14:48.624Z level=INFO msg="INFO ✅ Fase concluída phase=Phase06_ConfigureHandlers duration=4.258595ms"
time=2025-08-29T20:14:48.624Z level=INFO msg="INFO ▶️ Executando fase phase=Phase07_StartBackgroundWorkers timestamp=2025-08-29T20:14:48Z"
time=2025-08-29T20:14:48.624Z level=INFO msg="INFO 🎯 FASE 7: Inicialização de Background Workers"
time=2025-08-29T20:14:48.624Z level=WARN msg="Activity tracker or wait group not available for goroutine initialization"
time=2025-08-29T20:14:48.624Z level=INFO msg="INFO ✅ Background workers inicializados"
time=2025-08-29T20:14:48.624Z level=WARN msg="Activity tracker or user service not available for connection"
time=2025-08-29T20:14:48.624Z level=INFO msg="INFO ✅ Activity tracker conectado ao user service"
time=2025-08-29T20:14:48.625Z level=INFO msg="Database connection verified successfully"
time=2025-08-29T20:14:48.625Z level=INFO msg="INFO ✅ Schema do banco de dados verificado"
time=2025-08-29T20:14:48.625Z level=INFO msg="INFO ✅ Background workers inicializados com sucesso"
time=2025-08-29T20:14:48.625Z level=INFO msg="INFO ✅ Fase concluída phase=Phase07_StartBackgroundWorkers duration=568.898µs"
time=2025-08-29T20:14:48.625Z level=INFO msg="INFO ▶️ Executando fase phase=Phase08_StartServer timestamp=2025-08-29T20:14:48Z"
time=2025-08-29T20:14:48.625Z level=INFO msg="INFO 🎯 FASE 8: Inicialização Final e Runtime"
time=2025-08-29T20:14:48.625Z level=INFO msg="INFO ✅ Servidor marcado como ready para receber tráfego"
time=2025-08-29T20:14:48.625Z level=INFO msg="INFO 🚀 Iniciando servidor HTTP na porta configurada"
time=2025-08-29T20:14:48.725Z level=INFO msg="INFO ✅ Servidor HTTP iniciado com sucesso"
time=2025-08-29T20:14:48.725Z level=INFO msg="INFO ✅ Monitoramento de saúde em runtime iniciado"
time=2025-08-29T20:14:48.725Z level=INFO msg="INFO 🌟 TOQ Server pronto para servir uptime=126.29943ms"
time=2025-08-29T20:14:48.725Z level=INFO msg="INFO ✅ Fase concluída phase=Phase08_StartServer duration=100.404297ms"
time=2025-08-29T20:14:48.725Z level=INFO msg="INFO 🎉 TOQ Server inicializado com sucesso total_time=126.358591ms"
time=2025-08-29T20:14:57.290Z level=ERROR msg="Error rolling back transaction" error="sql: transaction has already been committed or rolled back"
time=2025-08-29T20:14:57.290Z level=ERROR msg="Error rolling back transaction" error="HTTP 500: Internal Server Error"
time=2025-08-29T20:15:06.731Z level=ERROR msg="mysqlpermissionadapter/Create: error executing statement" error="context canceled"
time=2025-08-29T20:15:06.731Z level=ERROR msg="Error rolling back transaction" error="sql: transaction has already been committed or rolled back"
time=2025-08-29T20:15:06.731Z level=ERROR msg="Error rolling back transaction" error="HTTP 500: Internal Server Error"
time=2025-08-29T20:15:06.731Z level=ERROR msg="Error rolling back transaction" error="sql: transaction has already been committed or rolled back"
time=2025-08-29T20:15:06.731Z level=ERROR msg="Error rolling back transaction" error="HTTP 500: Internal Server Error"
