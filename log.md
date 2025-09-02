Starting: /home/toq_admin/go/bin/dlv dap --listen=127.0.0.1:43169 --log-dest=3 from /codigos/go_code/toq_server/cmd
DAP server listening at: 127.0.0.1:43169
Type 'dlv help' for list of commands.
{"time":"2025-09-02T12:11:11.515495214Z","level":"INFO","msg":"🚀 Iniciando TOQ Server Bootstrap","version":"2.0.0","component":"bootstrap","log_level":"info","log_format":"json","log_output":"stdout"}
{"time":"2025-09-02T12:11:11.515632403Z","level":"INFO","msg":"▶️ Executando fase","phase":"Phase01_InitializeContext","component":"bootstrap","timestamp":"2025-09-02T12:11:11Z"}
{"time":"2025-09-02T12:11:11.515694052Z","level":"INFO","msg":"🎯 FASE 1: Inicialização de Contexto e Sinais"}
{"time":"2025-09-02T12:11:11.516277547Z","level":"INFO","msg":"✅ Contexto e sinais inicializados com sucesso"}
{"time":"2025-09-02T12:11:11.516312526Z","level":"INFO","msg":"✅ Fase concluída","phase":"Phase01_InitializeContext","component":"bootstrap","duration":"687.834µs"}
{"time":"2025-09-02T12:11:11.516334196Z","level":"INFO","msg":"▶️ Executando fase","phase":"Phase02_LoadConfiguration","component":"bootstrap","timestamp":"2025-09-02T12:11:11Z"}
{"time":"2025-09-02T12:11:11.516365916Z","level":"INFO","msg":"🎯 FASE 2: Carregamento e Validação de Configuração"}
{"time":"2025-09-02T12:11:11.516774952Z","level":"INFO","msg":"🔍 Iniciando servidor pprof na porta 6060"}
{"time":"2025-09-02T12:11:11.516915881Z","level":"INFO","msg":"✅ Servidor pprof iniciado em localhost:6060"}
{"time":"2025-09-02T12:11:11.518378228Z","level":"INFO","msg":"Configuration loaded successfully from YAML","path":"configs/env.yaml"}
{"time":"2025-09-02T12:11:11.518479057Z","level":"INFO","msg":"✅ Configuração carregada e validada com sucesso","version":"2.0.0"}
{"time":"2025-09-02T12:11:11.518529526Z","level":"INFO","msg":"✅ Fase concluída","phase":"Phase02_LoadConfiguration","component":"bootstrap","duration":"2.19293ms"}
{"time":"2025-09-02T12:11:11.518556146Z","level":"INFO","msg":"▶️ Executando fase","phase":"Phase03_InitializeInfrastructure","component":"bootstrap","timestamp":"2025-09-02T12:11:11Z"}
{"time":"2025-09-02T12:11:11.518590816Z","level":"INFO","msg":"🎯 FASE 3: Inicialização da Infraestrutura Core"}
{"time":"2025-09-02T12:11:11.521299721Z","level":"INFO","msg":"Database connection initialized","uri":"toq_user:toq_password@tcp(localhost:3306)/toq_db?parseTime=true&loc=UTC&timeout=30s&readTimeout=30s&writeTimeout=30s"}
{"time":"2025-09-02T12:11:11.521348371Z","level":"INFO","msg":"✅ Conexão com banco de dados estabelecida"}
{"time":"2025-09-02T12:11:11.523398512Z","level":"INFO","msg":"Redis cache connected successfully","url":"redis://localhost:6379/0"}
{"time":"2025-09-02T12:11:11.523451242Z","level":"INFO","msg":"✅ Sistema de cache Redis inicializado com sucesso"}
{"time":"2025-09-02T12:11:11.523935737Z","level":"INFO","msg":"OpenTelemetry tracing initialized","endpoint":"localhost:4318"}
{"time":"2025-09-02T12:11:11.524297734Z","level":"INFO","msg":"OpenTelemetry metrics initialized","endpoint":"localhost:4318"}
{"time":"2025-09-02T12:11:11.524330024Z","level":"INFO","msg":"OpenTelemetry initialized successfully","tracing_enabled":true,"metrics_enabled":true,"endpoint":"localhost:4318"}
{"time":"2025-09-02T12:11:11.524348043Z","level":"INFO","msg":"OpenTelemetry initialized successfully"}
{"time":"2025-09-02T12:11:11.524360193Z","level":"INFO","msg":"✅ OpenTelemetry inicializado (tracing + metrics)"}
{"time":"2025-09-02T12:11:11.524372043Z","level":"INFO","msg":"Creating metrics adapter"}
{"time":"2025-09-02T12:11:11.524535142Z","level":"INFO","msg":"✅ Adapter de métricas Prometheus inicializado"}
{"time":"2025-09-02T12:11:11.524556302Z","level":"INFO","msg":"✅ Infraestrutura core inicializada com sucesso"}
{"time":"2025-09-02T12:11:11.524577961Z","level":"INFO","msg":"✅ Fase concluída","phase":"Phase03_InitializeInfrastructure","component":"bootstrap","duration":"6.020195ms"}
{"time":"2025-09-02T12:11:11.524599701Z","level":"INFO","msg":"▶️ Executando fase","phase":"Phase04_InjectDependencies","component":"bootstrap","timestamp":"2025-09-02T12:11:11Z"}
{"time":"2025-09-02T12:11:11.524630781Z","level":"INFO","msg":"🎯 FASE 4: Injeção de Dependências via Factory Pattern"}
{"time":"2025-09-02T12:11:11.524658831Z","level":"INFO","msg":"Starting dependency injection using Factory Pattern"}
{"time":"2025-09-02T12:11:11.52471125Z","level":"INFO","msg":"Creating storage adapters"}
{"time":"2025-09-02T12:11:11.530769836Z","level":"INFO","msg":"Redis cache connected successfully","url":"redis://localhost:6379/0"}
{"time":"2025-09-02T12:11:11.530874685Z","level":"INFO","msg":"Successfully created all storage adapters"}
{"time":"2025-09-02T12:11:11.530891595Z","level":"INFO","msg":"✅ ActivityTracker criado com sucesso com Redis client"}
{"time":"2025-09-02T12:11:11.530903204Z","level":"INFO","msg":"Creating repository adapters"}
{"time":"2025-09-02T12:11:11.531087953Z","level":"INFO","msg":"Successfully created all repository adapters"}
{"time":"2025-09-02T12:11:11.531247771Z","level":"INFO","msg":"Assigning repository adapters"}
{"time":"2025-09-02T12:11:11.531260051Z","level":"INFO","msg":"Repository adapters assigned successfully"}
{"time":"2025-09-02T12:11:11.531270151Z","level":"INFO","msg":"Creating validation adapters"}
{"time":"2025-09-02T12:11:11.531285331Z","level":"INFO","msg":"Successfully created all validation adapters"}
{"time":"2025-09-02T12:11:11.531466279Z","level":"INFO","msg":"Creating external service adapters"}
{"time":"2025-09-02T12:11:11.531478739Z","level":"INFO","msg":"Creating external service adapters"}
{"time":"2025-09-02T12:11:11.534927888Z","level":"INFO","msg":"Creating S3 adapter","region":"us-east-1","bucket":"toq-app-media"}
{"time":"2025-09-02T12:11:11.535992378Z","level":"INFO","msg":"S3 adapter created successfully","bucket":"toq-app-media","region":"us-east-1"}
{"time":"2025-09-02T12:11:11.536302735Z","level":"INFO","msg":"Successfully created all external service adapters"}
{"time":"2025-09-02T12:11:11.536480884Z","level":"INFO","msg":"Initializing all services"}
{"time":"2025-09-02T12:11:11.537026879Z","level":"INFO","msg":"All services initialized successfully"}
{"time":"2025-09-02T12:11:11.537216717Z","level":"INFO","msg":"✅ TempBlockCleanerWorker initialized"}
{"time":"2025-09-02T12:11:11.537742332Z","level":"INFO","msg":"Dependency injection completed successfully using Factory Pattern"}
{"time":"2025-09-02T12:11:11.537767272Z","level":"INFO","msg":"✅ Injeção de dependências concluída via Factory Pattern"}
{"time":"2025-09-02T12:11:11.537814522Z","level":"INFO","msg":"✅ Fase concluída","phase":"Phase04_InjectDependencies","component":"bootstrap","duration":"13.207711ms"}
{"time":"2025-09-02T12:11:11.537850511Z","level":"INFO","msg":"▶️ Executando fase","phase":"Phase05_InitializeServices","component":"bootstrap","timestamp":"2025-09-02T12:11:11Z"}
{"time":"2025-09-02T12:11:11.5380162Z","level":"INFO","msg":"🎯 FASE 5: Inicialização de Serviços"}
{"time":"2025-09-02T12:11:11.538132129Z","level":"INFO","msg":"✅ Serviço inicializado","service":"GlobalService"}
{"time":"2025-09-02T12:11:11.538148939Z","level":"INFO","msg":"✅ Serviço inicializado","service":"PermissionService"}
{"time":"2025-09-02T12:11:11.538162038Z","level":"INFO","msg":"✅ Serviço inicializado","service":"UserService"}
{"time":"2025-09-02T12:11:11.538175288Z","level":"INFO","msg":"✅ Serviço inicializado","service":"ComplexService"}
{"time":"2025-09-02T12:11:11.538188658Z","level":"INFO","msg":"✅ Serviço inicializado","service":"ListingService"}
{"time":"2025-09-02T12:11:11.538200068Z","level":"INFO","msg":"✅ Todos os serviços inicializados com sucesso"}
{"time":"2025-09-02T12:11:11.538221788Z","level":"INFO","msg":"✅ Fase concluída","phase":"Phase05_InitializeServices","component":"bootstrap","duration":"367.677µs"}
{"time":"2025-09-02T12:11:11.538240248Z","level":"INFO","msg":"▶️ Executando fase","phase":"Phase06_ConfigureHandlers","component":"bootstrap","timestamp":"2025-09-02T12:11:11Z"}
{"time":"2025-09-02T12:11:11.538266838Z","level":"INFO","msg":"🎯 FASE 6: Configuração de Handlers e Rotas"}
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

{"time":"2025-09-02T12:11:11.538316917Z","level":"INFO","msg":"HTTP server initialized","port":":8080","read_timeout":"30s","write_timeout":"30s"}
{"time":"2025-09-02T12:11:11.538337797Z","level":"INFO","msg":"HTTP server initialization completed"}
{"time":"2025-09-02T12:11:11.538354717Z","level":"INFO","msg":"✅ Servidor HTTP configurado com TLS e middleware"}
{"time":"2025-09-02T12:11:11.538364967Z","level":"INFO","msg":"✅ Handlers HTTP preparados para criação"}
{"time":"2025-09-02T12:11:11.538375847Z","level":"INFO","msg":"Creating HTTP handlers"}
{"time":"2025-09-02T12:11:11.538391546Z","level":"INFO","msg":"Successfully created all HTTP handlers"}
{"time":"2025-09-02T12:11:11.538406906Z","level":"INFO","msg":"✅ HTTP handlers created successfully via factory"}
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
{"time":"2025-09-02T12:11:11.540121051Z","level":"INFO","msg":"HTTP handlers and routes configured successfully"}
{"time":"2025-09-02T12:11:11.540142481Z","level":"INFO","msg":"✅ Rotas e middlewares configurados"}
{"time":"2025-09-02T12:11:11.54015389Z","level":"INFO","msg":"✅ Health checks configurados"}
{"time":"2025-09-02T12:11:11.54016367Z","level":"INFO","msg":"✅ Handlers e rotas configurados com sucesso"}
{"time":"2025-09-02T12:11:11.54019646Z","level":"INFO","msg":"✅ Fase concluída","phase":"Phase06_ConfigureHandlers","component":"bootstrap","duration":"1.953542ms"}
{"time":"2025-09-02T12:11:11.54022319Z","level":"INFO","msg":"▶️ Executando fase","phase":"Phase07_StartBackgroundWorkers","component":"bootstrap","timestamp":"2025-09-02T12:11:11Z"}
{"time":"2025-09-02T12:11:11.540257969Z","level":"INFO","msg":"🎯 FASE 7: Inicialização de Background Workers"}
{"time":"2025-09-02T12:11:11.540287989Z","level":"INFO","msg":"Activity tracker batch worker started"}
{"time":"2025-09-02T12:11:11.540299679Z","level":"INFO","msg":"Temp block cleaner worker started"}
{"time":"2025-09-02T12:11:11.540310489Z","level":"INFO","msg":"✅ Background workers inicializados"}
{"time":"2025-09-02T12:11:11.540320829Z","level":"INFO","msg":"Activity tracker connected to user service"}
{"time":"2025-09-02T12:11:11.540330239Z","level":"INFO","msg":"✅ Activity tracker conectado ao user service"}
{"time":"2025-09-02T12:11:11.540421508Z","level":"INFO","msg":"TempBlockCleanerWorker started"}
{"time":"2025-09-02T12:11:11.540937983Z","level":"INFO","msg":"Activity batch worker started","interval":30000000000}
{"time":"2025-09-02T12:11:11.541113582Z","level":"INFO","msg":"Database connection verified successfully"}
{"time":"2025-09-02T12:11:11.541131792Z","level":"INFO","msg":"✅ Schema do banco de dados verificado"}
{"time":"2025-09-02T12:11:11.541142481Z","level":"INFO","msg":"✅ Background workers inicializados com sucesso"}
{"time":"2025-09-02T12:11:11.541214121Z","level":"INFO","msg":"✅ Fase concluída","phase":"Phase07_StartBackgroundWorkers","component":"bootstrap","duration":"985.781µs"}
{"time":"2025-09-02T12:11:11.54124558Z","level":"INFO","msg":"▶️ Executando fase","phase":"Phase08_StartServer","component":"bootstrap","timestamp":"2025-09-02T12:11:11Z"}
{"time":"2025-09-02T12:11:11.541414709Z","level":"INFO","msg":"🎯 FASE 8: Inicialização Final e Runtime"}
{"time":"2025-09-02T12:11:11.541437259Z","level":"INFO","msg":"✅ Servidor marcado como ready para receber tráfego"}
{"time":"2025-09-02T12:11:11.541818285Z","level":"INFO","msg":"🚀 Iniciando servidor HTTP na porta configurada"}
{"time":"2025-09-02T12:11:11.642816229Z","level":"INFO","msg":"✅ Servidor HTTP iniciado com sucesso"}
{"time":"2025-09-02T12:11:11.642899778Z","level":"INFO","msg":"✅ Monitoramento de saúde em runtime iniciado"}
{"time":"2025-09-02T12:11:11.642918178Z","level":"INFO","msg":"🌟 TOQ Server pronto para servir","uptime":127303785}
{"time":"2025-09-02T12:11:11.642946557Z","level":"INFO","msg":"✅ Fase concluída","phase":"Phase08_StartServer","component":"bootstrap","duration":"101.703876ms"}
{"time":"2025-09-02T12:11:11.642974877Z","level":"INFO","msg":"🎉 TOQ Server inicializado com sucesso","component":"bootstrap","total_time":"127.360884ms"}
{"time":"2025-09-02T12:11:17.482819477Z","level":"INFO","msg":"Security event logged","eventType":"signin_success","result":"success","timestamp":"2025-09-02T12:11:17.482793378Z","userID":3,"nationalID":"04679654805","ipAddress":"179.110.194.42","userAgent":"PostmanRuntime/7.45.0"}
{"time":"2025-09-02T12:11:17.482913917Z","level":"INFO","msg":"User signed in successfully","userID":3}
{"time":"2025-09-02T12:11:17.491279812Z","level":"INFO","msg":"HTTP Request","request_id":"8f2d3373-147e-4d0a-bb23-341e52b0843c","method":"POST","path":"/api/v1/auth/signin","status":200,"duration":155965866,"size":599,"client_ip":"179.110.194.42","user_agent":"PostmanRuntime/7.45.0"}
{"time":"2025-09-02T12:11:27.887776344Z","level":"INFO","msg":"Security event logged","eventType":"signin_success","result":"success","timestamp":"2025-09-02T12:11:27.887760174Z","userID":4,"nationalID":"05377401808","ipAddress":"179.110.194.42","userAgent":"PostmanRuntime/7.45.0"}
{"time":"2025-09-02T12:11:27.887855683Z","level":"INFO","msg":"User signed in successfully","userID":4}
{"time":"2025-09-02T12:11:27.89635195Z","level":"INFO","msg":"HTTP Request","request_id":"0a660bde-3c72-45d5-8b9d-3f656f8f8c0f","method":"POST","path":"/api/v1/auth/signin","status":200,"duration":157769319,"size":599,"client_ip":"179.110.194.42","user_agent":"PostmanRuntime/7.45.0"}
{"time":"2025-09-02T12:12:10.545310181Z","level":"ERROR","msg":"mysqlpermissionadapter/Read: error preparing statement","error":"Error 1054 (42S22): Unknown column 'p.slug' in 'field list'"}
{"time":"2025-09-02T12:12:10.546002126Z","level":"ERROR","msg":"Error checking permission","userID":4,"method":"POST","path":"/api/v1/user/signout","error":"HTTP 500: Internal Server Error"}
{"time":"2025-09-02T12:12:10.546314253Z","level":"ERROR","source":{"function":"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/middlewares.StructuredLoggingMiddleware.func1","file":"/codigos/go_code/toq_server/internal/adapter/left/http/middlewares/structured_logging_middleware.go","line":126},"msg":"HTTP Error","request_id":"e445b262-6e8e-43ea-90a7-f11ffda49d48","method":"POST","path":"/api/v1/user/signout","status":500,"duration":14006608,"size":48,"client_ip":"179.110.194.42","user_agent":"PostmanRuntime/7.45.0","user_id":4,"user_role_id":4,"role_status":"pending_both"}
