Starting: /home/toq_admin/go/bin/dlv dap --listen=127.0.0.1:34527 --log-dest=3 from /codigos/go_code/toq_server/cmd
DAP server listening at: 127.0.0.1:34527
Type 'dlv help' for list of commands.
{"time":"2025-09-02T11:20:15.356512519Z","level":"INFO","msg":"🚀 Iniciando TOQ Server Bootstrap","version":"2.0.0","component":"bootstrap","log_level":"info","log_format":"json","log_output":"stdout"}
{"time":"2025-09-02T11:20:15.35664068Z","level":"INFO","msg":"▶️ Executando fase","phase":"Phase01_InitializeContext","component":"bootstrap","timestamp":"2025-09-02T11:20:15Z"}
{"time":"2025-09-02T11:20:15.356739482Z","level":"INFO","msg":"🎯 FASE 1: Inicialização de Contexto e Sinais"}
{"time":"2025-09-02T11:20:15.35733231Z","level":"INFO","msg":"✅ Contexto e sinais inicializados com sucesso"}
{"time":"2025-09-02T11:20:15.35736535Z","level":"INFO","msg":"✅ Fase concluída","phase":"Phase01_InitializeContext","component":"bootstrap","duration":"724.49µs"}
{"time":"2025-09-02T11:20:15.35738708Z","level":"INFO","msg":"▶️ Executando fase","phase":"Phase02_LoadConfiguration","component":"bootstrap","timestamp":"2025-09-02T11:20:15Z"}
{"time":"2025-09-02T11:20:15.357417331Z","level":"INFO","msg":"🎯 FASE 2: Carregamento e Validação de Configuração"}
{"time":"2025-09-02T11:20:15.357505232Z","level":"INFO","msg":"🔍 Iniciando servidor pprof na porta 6060"}
{"time":"2025-09-02T11:20:15.357539402Z","level":"INFO","msg":"✅ Servidor pprof iniciado em localhost:6060"}
{"time":"2025-09-02T11:20:15.359222325Z","level":"INFO","msg":"Configuration loaded successfully from YAML","path":"configs/env.yaml"}
{"time":"2025-09-02T11:20:15.359323156Z","level":"INFO","msg":"✅ Configuração carregada e validada com sucesso","version":"2.0.0"}
{"time":"2025-09-02T11:20:15.359392907Z","level":"INFO","msg":"✅ Fase concluída","phase":"Phase02_LoadConfiguration","component":"bootstrap","duration":"2.003637ms"}
{"time":"2025-09-02T11:20:15.359415677Z","level":"INFO","msg":"▶️ Executando fase","phase":"Phase03_InitializeInfrastructure","component":"bootstrap","timestamp":"2025-09-02T11:20:15Z"}
{"time":"2025-09-02T11:20:15.359513928Z","level":"INFO","msg":"🎯 FASE 3: Inicialização da Infraestrutura Core"}
{"time":"2025-09-02T11:20:15.362180114Z","level":"INFO","msg":"Database connection initialized","uri":"toq_user:toq_password@tcp(localhost:3306)/toq_db?parseTime=true&loc=UTC&timeout=30s&readTimeout=30s&writeTimeout=30s"}
{"time":"2025-09-02T11:20:15.362523608Z","level":"INFO","msg":"✅ Conexão com banco de dados estabelecida"}
{"time":"2025-09-02T11:20:15.364466304Z","level":"INFO","msg":"Redis cache connected successfully","url":"redis://localhost:6379/0"}
{"time":"2025-09-02T11:20:15.364534685Z","level":"INFO","msg":"✅ Sistema de cache Redis inicializado com sucesso"}
{"time":"2025-09-02T11:20:15.364819968Z","level":"INFO","msg":"OpenTelemetry tracing initialized","endpoint":"localhost:4318"}
{"time":"2025-09-02T11:20:15.365124172Z","level":"INFO","msg":"OpenTelemetry metrics initialized","endpoint":"localhost:4318"}
{"time":"2025-09-02T11:20:15.365169603Z","level":"INFO","msg":"OpenTelemetry initialized successfully","tracing_enabled":true,"metrics_enabled":true,"endpoint":"localhost:4318"}
{"time":"2025-09-02T11:20:15.365187803Z","level":"INFO","msg":"OpenTelemetry initialized successfully"}
{"time":"2025-09-02T11:20:15.365199143Z","level":"INFO","msg":"✅ OpenTelemetry inicializado (tracing + metrics)"}
{"time":"2025-09-02T11:20:15.365210944Z","level":"INFO","msg":"Creating metrics adapter"}
{"time":"2025-09-02T11:20:15.365364696Z","level":"INFO","msg":"✅ Adapter de métricas Prometheus inicializado"}
{"time":"2025-09-02T11:20:15.365385946Z","level":"INFO","msg":"✅ Infraestrutura core inicializada com sucesso"}
{"time":"2025-09-02T11:20:15.365406476Z","level":"INFO","msg":"✅ Fase concluída","phase":"Phase03_InitializeInfrastructure","component":"bootstrap","duration":"5.988479ms"}
{"time":"2025-09-02T11:20:15.365427776Z","level":"INFO","msg":"▶️ Executando fase","phase":"Phase04_InjectDependencies","component":"bootstrap","timestamp":"2025-09-02T11:20:15Z"}
{"time":"2025-09-02T11:20:15.365467357Z","level":"INFO","msg":"🎯 FASE 4: Injeção de Dependências via Factory Pattern"}
{"time":"2025-09-02T11:20:15.365497567Z","level":"INFO","msg":"Starting dependency injection using Factory Pattern"}
{"time":"2025-09-02T11:20:15.365511288Z","level":"INFO","msg":"Creating storage adapters"}
{"time":"2025-09-02T11:20:15.367422813Z","level":"INFO","msg":"Redis cache connected successfully","url":"redis://localhost:6379/0"}
{"time":"2025-09-02T11:20:15.367498794Z","level":"INFO","msg":"Successfully created all storage adapters"}
{"time":"2025-09-02T11:20:15.367515504Z","level":"INFO","msg":"✅ ActivityTracker criado com sucesso com Redis client"}
{"time":"2025-09-02T11:20:15.367525684Z","level":"INFO","msg":"Creating repository adapters"}
{"time":"2025-09-02T11:20:15.367581615Z","level":"INFO","msg":"Successfully created all repository adapters"}
{"time":"2025-09-02T11:20:15.367594985Z","level":"INFO","msg":"Assigning repository adapters"}
{"time":"2025-09-02T11:20:15.367606075Z","level":"INFO","msg":"Repository adapters assigned successfully"}
{"time":"2025-09-02T11:20:15.367615745Z","level":"INFO","msg":"Creating validation adapters"}
{"time":"2025-09-02T11:20:15.367627996Z","level":"INFO","msg":"Successfully created all validation adapters"}
{"time":"2025-09-02T11:20:15.367638816Z","level":"INFO","msg":"Creating external service adapters"}
{"time":"2025-09-02T11:20:15.367647506Z","level":"INFO","msg":"Creating external service adapters"}
{"time":"2025-09-02T11:20:15.369067305Z","level":"INFO","msg":"Creating S3 adapter","region":"us-east-1","bucket":"toq-app-media"}
{"time":"2025-09-02T11:20:15.369616662Z","level":"INFO","msg":"S3 adapter created successfully","bucket":"toq-app-media","region":"us-east-1"}
{"time":"2025-09-02T11:20:15.369644802Z","level":"INFO","msg":"Successfully created all external service adapters"}
{"time":"2025-09-02T11:20:15.369684143Z","level":"INFO","msg":"Initializing all services"}
{"time":"2025-09-02T11:20:15.369703273Z","level":"INFO","msg":"All services initialized successfully"}
{"time":"2025-09-02T11:20:15.369716473Z","level":"INFO","msg":"✅ TempBlockCleanerWorker initialized"}
{"time":"2025-09-02T11:20:15.369725963Z","level":"INFO","msg":"Dependency injection completed successfully using Factory Pattern"}
{"time":"2025-09-02T11:20:15.369736363Z","level":"INFO","msg":"✅ Injeção de dependências concluída via Factory Pattern"}
{"time":"2025-09-02T11:20:15.369805444Z","level":"INFO","msg":"✅ Fase concluída","phase":"Phase04_InjectDependencies","component":"bootstrap","duration":"4.374698ms"}
{"time":"2025-09-02T11:20:15.369864955Z","level":"INFO","msg":"▶️ Executando fase","phase":"Phase05_InitializeServices","component":"bootstrap","timestamp":"2025-09-02T11:20:15Z"}
{"time":"2025-09-02T11:20:15.369913816Z","level":"INFO","msg":"🎯 FASE 5: Inicialização de Serviços"}
{"time":"2025-09-02T11:20:15.369942136Z","level":"INFO","msg":"✅ Serviço inicializado","service":"GlobalService"}
{"time":"2025-09-02T11:20:15.369955566Z","level":"INFO","msg":"✅ Serviço inicializado","service":"PermissionService"}
{"time":"2025-09-02T11:20:15.369967926Z","level":"INFO","msg":"✅ Serviço inicializado","service":"UserService"}
{"time":"2025-09-02T11:20:15.369978887Z","level":"INFO","msg":"✅ Serviço inicializado","service":"ComplexService"}
{"time":"2025-09-02T11:20:15.369990067Z","level":"INFO","msg":"✅ Serviço inicializado","service":"ListingService"}
{"time":"2025-09-02T11:20:15.370014337Z","level":"INFO","msg":"✅ Todos os serviços inicializados com sucesso"}
{"time":"2025-09-02T11:20:15.370041007Z","level":"INFO","msg":"✅ Fase concluída","phase":"Phase05_InitializeServices","component":"bootstrap","duration":"208.032µs"}
{"time":"2025-09-02T11:20:15.370058648Z","level":"INFO","msg":"▶️ Executando fase","phase":"Phase06_ConfigureHandlers","component":"bootstrap","timestamp":"2025-09-02T11:20:15Z"}
{"time":"2025-09-02T11:20:15.370085898Z","level":"INFO","msg":"🎯 FASE 6: Configuração de Handlers e Rotas"}
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

{"time":"2025-09-02T11:20:15.370129539Z","level":"INFO","msg":"HTTP server initialized","port":":8080","read_timeout":"30s","write_timeout":"30s"}
{"time":"2025-09-02T11:20:15.37020114Z","level":"INFO","msg":"HTTP server initialization completed"}
{"time":"2025-09-02T11:20:15.37021376Z","level":"INFO","msg":"✅ Servidor HTTP configurado com TLS e middleware"}
{"time":"2025-09-02T11:20:15.3702236Z","level":"INFO","msg":"✅ Handlers HTTP preparados para criação"}
{"time":"2025-09-02T11:20:15.37023431Z","level":"INFO","msg":"Creating HTTP handlers"}
{"time":"2025-09-02T11:20:15.37024605Z","level":"INFO","msg":"Successfully created all HTTP handlers"}
{"time":"2025-09-02T11:20:15.37025618Z","level":"INFO","msg":"✅ HTTP handlers created successfully via factory"}
[GIN-debug] GET    /swagger/*any             --> github.com/swaggo/gin-swagger.CustomWrapHandler.func1 (6 handlers)
[GIN-debug] POST   /api/v1/auth/owner        --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/handlers/auth_handlers.(*AuthHandler).CreateOwner-fm (6 handlers)
[GIN-debug] POST   /api/v1/auth/realtor      --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/handlers/auth_handlers.(*AuthHandler).CreateRealtor-fm (6 handlers)
[GIN-debug] POST   /api/v1/auth/agency       --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/handlers/auth_handlers.(*AuthHandler).CreateAgency-fm (6 handlers)
[GIN-debug] POST   /api/v1/auth/signin       --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/handlers/auth_handlers.(*AuthHandler).SignIn-fm (6 handlers)
[GIN-debug] POST   /api/v1/auth/refresh      --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/handlers/auth_handlers.(*AuthHandler).RefreshToken-fm (6 handlers)
[GIN-debug] POST   /api/v1/auth/password/request --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/handlers/auth_handlers.(*AuthHandler).RequestPasswordChange-fm (6 handlers)
[GIN-debug] POST   /api/v1/auth/password/confirm --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/handlers/auth_handlers.(*AuthHandler).ConfirmPasswordChange-fm (6 handlers)
[GIN-debug] GET    /api/v1/user/profile      --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func1 (8 handlers)
[GIN-debug] PUT    /api/v1/user/profile      --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func2 (8 handlers)
[GIN-debug] DELETE /api/v1/user/account      --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func3 (8 handlers)
[GIN-debug] GET    /api/v1/user/onboarding   --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func4 (8 handlers)
[GIN-debug] GET    /api/v1/user/roles        --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func5 (8 handlers)
[GIN-debug] GET    /api/v1/user/home         --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func6 (8 handlers)
[GIN-debug] PUT    /api/v1/user/opt-status   --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func7 (8 handlers)
[GIN-debug] POST   /api/v1/user/signout      --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/handlers/user_handlers.(*UserHandler).SignOut-fm (8 handlers)
[GIN-debug] POST   /api/v1/user/photo/upload-url --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func8 (8 handlers)
[GIN-debug] GET    /api/v1/user/profile/thumbnails --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func9 (8 handlers)
[GIN-debug] POST   /api/v1/user/email/request --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func10 (8 handlers)
[GIN-debug] POST   /api/v1/user/email/confirm --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func11 (8 handlers)
[GIN-debug] POST   /api/v1/user/email/resend --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func12 (8 handlers)
[GIN-debug] POST   /api/v1/user/phone/request --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func13 (8 handlers)
[GIN-debug] POST   /api/v1/user/phone/confirm --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func14 (8 handlers)
[GIN-debug] POST   /api/v1/user/phone/resend --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func15 (8 handlers)
[GIN-debug] POST   /api/v1/user/role/alternative --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func16 (8 handlers)
[GIN-debug] POST   /api/v1/user/role/switch  --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func17 (8 handlers)
[GIN-debug] POST   /api/v1/agency/invite-realtor --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func18 (8 handlers)
[GIN-debug] GET    /api/v1/agency/realtors   --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func19 (8 handlers)
[GIN-debug] GET    /api/v1/agency/realtors/:id --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func20 (8 handlers)
[GIN-debug] DELETE /api/v1/agency/realtors/:id --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func21 (8 handlers)
[GIN-debug] POST   /api/v1/realtor/creci/verify --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func22 (8 handlers)
[GIN-debug] POST   /api/v1/realtor/creci/upload-url --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func23 (8 handlers)
[GIN-debug] POST   /api/v1/realtor/invitation/accept --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func24 (8 handlers)
[GIN-debug] POST   /api/v1/realtor/invitation/reject --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func25 (8 handlers)
[GIN-debug] GET    /api/v1/realtor/agency    --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func26 (8 handlers)
[GIN-debug] DELETE /api/v1/realtor/agency    --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterUserRoutes.func27 (8 handlers)
[GIN-debug] GET    /api/v1/listings          --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func1 (8 handlers)
[GIN-debug] POST   /api/v1/listings          --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func2 (8 handlers)
[GIN-debug] GET    /api/v1/listings/search   --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func3 (8 handlers)
[GIN-debug] GET    /api/v1/listings/options  --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func4 (8 handlers)
[GIN-debug] GET    /api/v1/listings/features/base --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func5 (8 handlers)
[GIN-debug] GET    /api/v1/listings/favorites --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func6 (8 handlers)
[GIN-debug] GET    /api/v1/listings/:id      --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func7 (8 handlers)
[GIN-debug] PUT    /api/v1/listings/:id      --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func8 (8 handlers)
[GIN-debug] DELETE /api/v1/listings/:id      --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func9 (8 handlers)
[GIN-debug] POST   /api/v1/listings/:id/end-update --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func10 (8 handlers)
[GIN-debug] GET    /api/v1/listings/:id/status --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func11 (8 handlers)
[GIN-debug] POST   /api/v1/listings/:id/approve --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func12 (8 handlers)
[GIN-debug] POST   /api/v1/listings/:id/reject --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func13 (8 handlers)
[GIN-debug] POST   /api/v1/listings/:id/suspend --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func14 (8 handlers)
[GIN-debug] POST   /api/v1/listings/:id/release --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func15 (8 handlers)
[GIN-debug] POST   /api/v1/listings/:id/copy --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func16 (8 handlers)
[GIN-debug] POST   /api/v1/listings/:id/share --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func17 (8 handlers)
[GIN-debug] POST   /api/v1/listings/:id/favorite --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func18 (8 handlers)
[GIN-debug] DELETE /api/v1/listings/:id/favorite --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func19 (8 handlers)
[GIN-debug] POST   /api/v1/listings/:id/visit/request --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func20 (8 handlers)
[GIN-debug] GET    /api/v1/listings/:id/visits --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func21 (8 handlers)
[GIN-debug] POST   /api/v1/listings/:id/offers --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func22 (8 handlers)
[GIN-debug] GET    /api/v1/listings/:id/offers --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func23 (8 handlers)
[GIN-debug] GET    /api/v1/visits            --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func24 (8 handlers)
[GIN-debug] DELETE /api/v1/visits/:id        --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func25 (8 handlers)
[GIN-debug] POST   /api/v1/visits/:id/confirm --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func26 (8 handlers)
[GIN-debug] POST   /api/v1/visits/:id/approve --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func27 (8 handlers)
[GIN-debug] POST   /api/v1/visits/:id/reject --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func28 (8 handlers)
[GIN-debug] GET    /api/v1/offers            --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func29 (8 handlers)
[GIN-debug] PUT    /api/v1/offers/:id        --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func30 (8 handlers)
[GIN-debug] DELETE /api/v1/offers/:id        --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func31 (8 handlers)
[GIN-debug] POST   /api/v1/offers/:id/send   --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func32 (8 handlers)
[GIN-debug] POST   /api/v1/offers/:id/approve --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func33 (8 handlers)
[GIN-debug] POST   /api/v1/offers/:id/reject --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func34 (8 handlers)
[GIN-debug] POST   /api/v1/realtors/:id/evaluate --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func35 (8 handlers)
[GIN-debug] POST   /api/v1/owners/:id/evaluate --> github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes.RegisterListingRoutes.func36 (8 handlers)
[GIN-debug] GET    /healthz                  --> github.com/giulio-alfieri/toq_server/internal/core/config.(*config).setupBasicRoutes.func2 (7 handlers)
[GIN-debug] GET    /readyz                   --> github.com/giulio-alfieri/toq_server/internal/core/config.(*config).setupBasicRoutes.func3 (7 handlers)
[GIN-debug] GET    /metrics                  --> go:interface { GetMetrics(*github.com/gin-gonic/gin.Context) }.GetMetrics-fm (6 handlers)
[GIN-debug] GET    /api/v1/ping              --> github.com/giulio-alfieri/toq_server/internal/core/config.(*config).setupBasicRoutes.func4 (6 handlers)
{"time":"2025-09-02T11:20:15.37175944Z","level":"INFO","msg":"HTTP handlers and routes configured successfully"}
{"time":"2025-09-02T11:20:15.37177958Z","level":"INFO","msg":"✅ Rotas e middlewares configurados"}
{"time":"2025-09-02T11:20:15.371789431Z","level":"INFO","msg":"✅ Health checks configurados"}
{"time":"2025-09-02T11:20:15.371798561Z","level":"INFO","msg":"✅ Handlers e rotas configurados com sucesso"}
{"time":"2025-09-02T11:20:15.371836741Z","level":"INFO","msg":"✅ Fase concluída","phase":"Phase06_ConfigureHandlers","component":"bootstrap","duration":"1.773733ms"}
{"time":"2025-09-02T11:20:15.371859901Z","level":"INFO","msg":"▶️ Executando fase","phase":"Phase07_StartBackgroundWorkers","component":"bootstrap","timestamp":"2025-09-02T11:20:15Z"}
{"time":"2025-09-02T11:20:15.371893752Z","level":"INFO","msg":"🎯 FASE 7: Inicialização de Background Workers"}
{"time":"2025-09-02T11:20:15.371920952Z","level":"INFO","msg":"Activity tracker batch worker started"}
{"time":"2025-09-02T11:20:15.371931312Z","level":"INFO","msg":"Temp block cleaner worker started"}
{"time":"2025-09-02T11:20:15.371940753Z","level":"INFO","msg":"✅ Background workers inicializados"}
{"time":"2025-09-02T11:20:15.371950633Z","level":"INFO","msg":"Activity tracker connected to user service"}
{"time":"2025-09-02T11:20:15.371959153Z","level":"INFO","msg":"✅ Activity tracker conectado ao user service"}
{"time":"2025-09-02T11:20:15.372199446Z","level":"INFO","msg":"TempBlockCleanerWorker started"}
{"time":"2025-09-02T11:20:15.372749953Z","level":"INFO","msg":"Activity batch worker started","interval":30000000000}
{"time":"2025-09-02T11:20:15.372807264Z","level":"INFO","msg":"Database connection verified successfully"}
{"time":"2025-09-02T11:20:15.372820784Z","level":"INFO","msg":"✅ Schema do banco de dados verificado"}
{"time":"2025-09-02T11:20:15.372831174Z","level":"INFO","msg":"✅ Background workers inicializados com sucesso"}
{"time":"2025-09-02T11:20:15.372850455Z","level":"INFO","msg":"✅ Fase concluída","phase":"Phase07_StartBackgroundWorkers","component":"bootstrap","duration":"988.803µs"}
{"time":"2025-09-02T11:20:15.372870595Z","level":"INFO","msg":"▶️ Executando fase","phase":"Phase08_StartServer","component":"bootstrap","timestamp":"2025-09-02T11:20:15Z"}
{"time":"2025-09-02T11:20:15.372896355Z","level":"INFO","msg":"🎯 FASE 8: Inicialização Final e Runtime"}
{"time":"2025-09-02T11:20:15.372918165Z","level":"INFO","msg":"✅ Servidor marcado como ready para receber tráfego"}
{"time":"2025-09-02T11:20:15.372937466Z","level":"INFO","msg":"🚀 Iniciando servidor HTTP na porta configurada"}
{"time":"2025-09-02T11:20:15.47309041Z","level":"INFO","msg":"✅ Servidor HTTP iniciado com sucesso"}
{"time":"2025-09-02T11:20:15.473161681Z","level":"INFO","msg":"✅ Monitoramento de saúde em runtime iniciado"}
{"time":"2025-09-02T11:20:15.473176791Z","level":"INFO","msg":"🌟 TOQ Server pronto para servir","uptime":116546381}
{"time":"2025-09-02T11:20:15.473205972Z","level":"INFO","msg":"✅ Fase concluída","phase":"Phase08_StartServer","component":"bootstrap","duration":"100.332667ms"}
{"time":"2025-09-02T11:20:15.473234082Z","level":"INFO","msg":"🎉 TOQ Server inicializado com sucesso","component":"bootstrap","total_time":"116.603832ms"}
{"time":"2025-09-02T11:20:40.026411154Z","level":"ERROR","msg":"mysqlpermissionadapter/ReadRow: error getting column count","error":"driver: bad connection"}
{"time":"2025-09-02T11:20:40.026750819Z","level":"ERROR","msg":"Failed to get user role for temp block check","userID":3,"error":"driver: bad connection"}
{"time":"2025-09-02T11:20:40.026774519Z","level":"ERROR","msg":"Failed to check if user is temporarily blocked","userID":3,"error":"HTTP 500: Failed to check user status"}
{"time":"2025-09-02T11:20:40.027017342Z","level":"ERROR","msg":"Error rolling back transaction","error":"invalid connection"}
{"time":"2025-09-02T11:20:40.027083263Z","level":"ERROR","msg":"Error rolling back transaction","error":"rollback tx: invalid connection"}
{"time":"2025-09-02T11:20:40.027325606Z","level":"ERROR","source":{"function":"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/middlewares.StructuredLoggingMiddleware.func1","file":"/codigos/go_code/toq_server/internal/adapter/left/http/middlewares/structured_logging_middleware.go","line":126},"msg":"HTTP Error","request_id":"674b439b-e4f3-4e84-9739-ef72554aee54","method":"POST","path":"/api/v1/auth/signin","status":500,"duration":4628222,"size":46,"client_ip":"179.110.194.42","user_agent":"PostmanRuntime/7.45.0"}
