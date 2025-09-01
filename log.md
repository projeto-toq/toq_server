Starting: /home/toq_admin/go/bin/dlv dap --listen=127.0.0.1:40961 --log-dest=3 from /codigos/go_code/toq_server/cmd
DAP server listening at: 127.0.0.1:40961
Type 'dlv help' for list of commands.
{"time":"2025-09-01T13:11:20.056163503Z","level":"INFO","msg":"🚀 Iniciando TOQ Server Bootstrap","version":"2.0.0","component":"bootstrap","log_level":"info","log_format":"json","log_output":"stdout"}
{"time":"2025-09-01T13:11:20.056334908Z","level":"INFO","msg":"▶️ Executando fase","phase":"Phase01_InitializeContext","component":"bootstrap","timestamp":"2025-09-01T13:11:20Z"}
{"time":"2025-09-01T13:11:20.05641017Z","level":"INFO","msg":"🎯 FASE 1: Inicialização de Contexto e Sinais"}
{"time":"2025-09-01T13:11:20.057104087Z","level":"INFO","msg":"✅ Contexto e sinais inicializados com sucesso"}
{"time":"2025-09-01T13:11:20.057149978Z","level":"INFO","msg":"✅ Fase concluída","phase":"Phase01_InitializeContext","component":"bootstrap","duration":"817.061µs"}
{"time":"2025-09-01T13:11:20.057173359Z","level":"INFO","msg":"▶️ Executando fase","phase":"Phase02_LoadConfiguration","component":"bootstrap","timestamp":"2025-09-01T13:11:20Z"}
{"time":"2025-09-01T13:11:20.057243641Z","level":"INFO","msg":"🎯 FASE 2: Carregamento e Validação de Configuração"}
{"time":"2025-09-01T13:11:20.057810475Z","level":"INFO","msg":"🔍 Iniciando servidor pprof na porta 6060"}
{"time":"2025-09-01T13:11:20.057853126Z","level":"INFO","msg":"✅ Servidor pprof iniciado em localhost:6060"}
{"time":"2025-09-01T13:11:20.059443286Z","level":"INFO","msg":"Configuration loaded successfully from YAML","path":"configs/env.yaml"}
{"time":"2025-09-01T13:11:20.059536798Z","level":"INFO","msg":"✅ Configuração carregada e validada com sucesso","version":"2.0.0"}
{"time":"2025-09-01T13:11:20.05959793Z","level":"INFO","msg":"✅ Fase concluída","phase":"Phase02_LoadConfiguration","component":"bootstrap","duration":"2.422161ms"}
{"time":"2025-09-01T13:11:20.05962389Z","level":"INFO","msg":"▶️ Executando fase","phase":"Phase03_InitializeInfrastructure","component":"bootstrap","timestamp":"2025-09-01T13:11:20Z"}
{"time":"2025-09-01T13:11:20.059680602Z","level":"INFO","msg":"🎯 FASE 3: Inicialização da Infraestrutura Core"}
{"time":"2025-09-01T13:11:20.063882777Z","level":"INFO","msg":"Database connection initialized","uri":"toq_user:toq_password@tcp(localhost:3306)/toq_db?parseTime=true&loc=UTC&timeout=30s&readTimeout=30s&writeTimeout=30s"}
{"time":"2025-09-01T13:11:20.06398277Z","level":"INFO","msg":"✅ Conexão com banco de dados estabelecida"}
{"time":"2025-09-01T13:11:20.066462822Z","level":"INFO","msg":"Redis cache connected successfully","url":"redis://localhost:6379/0"}
{"time":"2025-09-01T13:11:20.066505343Z","level":"INFO","msg":"✅ Sistema de cache Redis inicializado com sucesso"}
{"time":"2025-09-01T13:11:20.06678211Z","level":"INFO","msg":"OpenTelemetry tracing initialized","endpoint":"localhost:4318"}
{"time":"2025-09-01T13:11:20.067068258Z","level":"INFO","msg":"OpenTelemetry metrics initialized","endpoint":"localhost:4318"}
{"time":"2025-09-01T13:11:20.067093618Z","level":"INFO","msg":"OpenTelemetry initialized successfully","tracing_enabled":true,"metrics_enabled":true,"endpoint":"localhost:4318"}
{"time":"2025-09-01T13:11:20.067106638Z","level":"INFO","msg":"OpenTelemetry initialized successfully"}
{"time":"2025-09-01T13:11:20.067115209Z","level":"INFO","msg":"✅ OpenTelemetry inicializado (tracing + metrics)"}
{"time":"2025-09-01T13:11:20.067139579Z","level":"INFO","msg":"Creating metrics adapter"}
{"time":"2025-09-01T13:11:20.067300153Z","level":"INFO","msg":"✅ Adapter de métricas Prometheus inicializado"}
{"time":"2025-09-01T13:11:20.067315664Z","level":"INFO","msg":"✅ Infraestrutura core inicializada com sucesso"}
{"time":"2025-09-01T13:11:20.067332004Z","level":"INFO","msg":"✅ Fase concluída","phase":"Phase03_InitializeInfrastructure","component":"bootstrap","duration":"7.706684ms"}
{"time":"2025-09-01T13:11:20.067351665Z","level":"INFO","msg":"▶️ Executando fase","phase":"Phase04_InjectDependencies","component":"bootstrap","timestamp":"2025-09-01T13:11:20Z"}
{"time":"2025-09-01T13:11:20.067435447Z","level":"INFO","msg":"🎯 FASE 4: Injeção de Dependências via Factory Pattern"}
{"time":"2025-09-01T13:11:20.067464707Z","level":"INFO","msg":"Starting dependency injection using Factory Pattern"}
{"time":"2025-09-01T13:11:20.067476438Z","level":"INFO","msg":"Creating storage adapters"}
{"time":"2025-09-01T13:11:20.072105954Z","level":"INFO","msg":"Redis cache connected successfully","url":"redis://localhost:6379/0"}
{"time":"2025-09-01T13:11:20.072166606Z","level":"INFO","msg":"Successfully created all storage adapters"}
{"time":"2025-09-01T13:11:20.072198257Z","level":"INFO","msg":"✅ ActivityTracker criado com sucesso com Redis client"}
{"time":"2025-09-01T13:11:20.072209167Z","level":"INFO","msg":"Creating repository adapters"}
{"time":"2025-09-01T13:11:20.072263728Z","level":"INFO","msg":"Successfully created all repository adapters"}
{"time":"2025-09-01T13:11:20.072276618Z","level":"INFO","msg":"Assigning repository adapters"}
{"time":"2025-09-01T13:11:20.072302929Z","level":"INFO","msg":"Repository adapters assigned successfully"}
{"time":"2025-09-01T13:11:20.072312699Z","level":"INFO","msg":"Creating validation adapters"}
{"time":"2025-09-01T13:11:20.07232826Z","level":"INFO","msg":"Successfully created all validation adapters"}
{"time":"2025-09-01T13:11:20.07233979Z","level":"INFO","msg":"Creating external service adapters"}
{"time":"2025-09-01T13:11:20.072368581Z","level":"INFO","msg":"Creating external service adapters"}
{"time":"2025-09-01T13:11:20.077263723Z","level":"INFO","msg":"Creating S3 adapter","region":"us-east-1","bucket":"toq-app-media"}
{"time":"2025-09-01T13:11:20.078684799Z","level":"INFO","msg":"S3 adapter created successfully","bucket":"toq-app-media","region":"us-east-1"}
{"time":"2025-09-01T13:11:20.07872851Z","level":"INFO","msg":"Successfully created all external service adapters"}
{"time":"2025-09-01T13:11:20.078759841Z","level":"INFO","msg":"Initializing all services"}
{"time":"2025-09-01T13:11:20.078776141Z","level":"INFO","msg":"All services initialized successfully"}
{"time":"2025-09-01T13:11:20.078789332Z","level":"INFO","msg":"✅ TempBlockCleanerWorker initialized"}
{"time":"2025-09-01T13:11:20.078799762Z","level":"INFO","msg":"Dependency injection completed successfully using Factory Pattern"}
{"time":"2025-09-01T13:11:20.078831943Z","level":"INFO","msg":"✅ Injeção de dependências concluída via Factory Pattern"}
{"time":"2025-09-01T13:11:20.080931606Z","level":"INFO","msg":"✅ Fase concluída","phase":"Phase04_InjectDependencies","component":"bootstrap","duration":"13.56066ms"}
{"time":"2025-09-01T13:11:20.081046388Z","level":"INFO","msg":"▶️ Executando fase","phase":"Phase05_InitializeServices","component":"bootstrap","timestamp":"2025-09-01T13:11:20Z"}
{"time":"2025-09-01T13:11:20.08111237Z","level":"INFO","msg":"🎯 FASE 5: Inicialização de Serviços"}
{"time":"2025-09-01T13:11:20.081146831Z","level":"INFO","msg":"✅ Serviço inicializado","service":"GlobalService"}
{"time":"2025-09-01T13:11:20.081164171Z","level":"INFO","msg":"✅ Serviço inicializado","service":"PermissionService"}
{"time":"2025-09-01T13:11:20.081200092Z","level":"INFO","msg":"✅ Serviço inicializado","service":"UserService"}
{"time":"2025-09-01T13:11:20.081218333Z","level":"INFO","msg":"✅ Serviço inicializado","service":"ComplexService"}
{"time":"2025-09-01T13:11:20.081231323Z","level":"INFO","msg":"✅ Serviço inicializado","service":"ListingService"}
{"time":"2025-09-01T13:11:20.081261924Z","level":"INFO","msg":"✅ Todos os serviços inicializados com sucesso"}
{"time":"2025-09-01T13:11:20.081283574Z","level":"INFO","msg":"✅ Fase concluída","phase":"Phase05_InitializeServices","component":"bootstrap","duration":"237.066µs"}
{"time":"2025-09-01T13:11:20.081328776Z","level":"INFO","msg":"▶️ Executando fase","phase":"Phase06_ConfigureHandlers","component":"bootstrap","timestamp":"2025-09-01T13:11:20Z"}
{"time":"2025-09-01T13:11:20.081362146Z","level":"INFO","msg":"🎯 FASE 6: Configuração de Handlers e Rotas"}
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)
{"time":"2025-09-01T13:11:20.081436118Z","level":"INFO","msg":"HTTP server initialized","port":":8080","read_timeout":"30s","write_timeout":"30s"}
{"time":"2025-09-01T13:11:20.081454099Z","level":"INFO","msg":"HTTP server initialization completed"}
{"time":"2025-09-01T13:11:20.081465209Z","level":"INFO","msg":"✅ Servidor HTTP configurado com TLS e middleware"}
{"time":"2025-09-01T13:11:20.08149661Z","level":"INFO","msg":"✅ Handlers HTTP preparados para criação"}
{"time":"2025-09-01T13:11:20.0815086Z","level":"INFO","msg":"Creating HTTP handlers"}
{"time":"2025-09-01T13:11:20.08152192Z","level":"INFO","msg":"Successfully created all HTTP handlers"}
{"time":"2025-09-01T13:11:20.081533081Z","level":"INFO","msg":"✅ HTTP handlers created successfully via factory"}
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
{"time":"2025-09-01T13:11:20.085931301Z","level":"INFO","msg":"HTTP handlers and routes configured successfully"}
{"time":"2025-09-01T13:11:20.085956682Z","level":"INFO","msg":"✅ Rotas e middlewares configurados"}
{"time":"2025-09-01T13:11:20.085967512Z","level":"INFO","msg":"✅ Health checks configurados"}
{"time":"2025-09-01T13:11:20.085977262Z","level":"INFO","msg":"✅ Handlers e rotas configurados com sucesso"}
{"time":"2025-09-01T13:11:20.086428164Z","level":"INFO","msg":"✅ Fase concluída","phase":"Phase06_ConfigureHandlers","component":"bootstrap","duration":"5.091718ms"}
{"time":"2025-09-01T13:11:20.086467015Z","level":"INFO","msg":"▶️ Executando fase","phase":"Phase07_StartBackgroundWorkers","component":"bootstrap","timestamp":"2025-09-01T13:11:20Z"}
{"time":"2025-09-01T13:11:20.086494135Z","level":"INFO","msg":"🎯 FASE 7: Inicialização de Background Workers"}
{"time":"2025-09-01T13:11:20.086531046Z","level":"INFO","msg":"Activity tracker batch worker started"}
{"time":"2025-09-01T13:11:20.086584738Z","level":"INFO","msg":"Temp block cleaner worker started"}
{"time":"2025-09-01T13:11:20.086594598Z","level":"INFO","msg":"✅ Background workers inicializados"}
{"time":"2025-09-01T13:11:20.086604768Z","level":"INFO","msg":"Activity tracker connected to user service"}
{"time":"2025-09-01T13:11:20.086615098Z","level":"INFO","msg":"✅ Activity tracker conectado ao user service"}
{"time":"2025-09-01T13:11:20.087011378Z","level":"INFO","msg":"Activity batch worker started","interval":30000000000}
{"time":"2025-09-01T13:11:20.08707338Z","level":"INFO","msg":"TempBlockCleanerWorker started"}
{"time":"2025-09-01T13:11:20.089022849Z","level":"INFO","msg":"Database connection verified successfully"}
{"time":"2025-09-01T13:11:20.08905313Z","level":"INFO","msg":"✅ Schema do banco de dados verificado"}
{"time":"2025-09-01T13:11:20.08906348Z","level":"INFO","msg":"✅ Background workers inicializados com sucesso"}
{"time":"2025-09-01T13:11:20.089088621Z","level":"INFO","msg":"✅ Fase concluída","phase":"Phase07_StartBackgroundWorkers","component":"bootstrap","duration":"2.616595ms"}
{"time":"2025-09-01T13:11:20.089106591Z","level":"INFO","msg":"▶️ Executando fase","phase":"Phase08_StartServer","component":"bootstrap","timestamp":"2025-09-01T13:11:20Z"}
{"time":"2025-09-01T13:11:20.089131312Z","level":"INFO","msg":"🎯 FASE 8: Inicialização Final e Runtime"}
{"time":"2025-09-01T13:11:20.089155952Z","level":"INFO","msg":"✅ Servidor marcado como ready para receber tráfego"}
{"time":"2025-09-01T13:11:20.089174003Z","level":"INFO","msg":"🚀 Iniciando servidor HTTP na porta configurada"}
{"time":"2025-09-01T13:11:20.189493702Z","level":"INFO","msg":"✅ Servidor HTTP iniciado com sucesso"}
{"time":"2025-09-01T13:11:20.189577474Z","level":"INFO","msg":"✅ Monitoramento de saúde em runtime iniciado"}
{"time":"2025-09-01T13:11:20.189599285Z","level":"INFO","msg":"🌟 TOQ Server pronto para servir","uptime":133277058}
{"time":"2025-09-01T13:11:20.189646536Z","level":"INFO","msg":"✅ Fase concluída","phase":"Phase08_StartServer","component":"bootstrap","duration":"100.536205ms"}
{"time":"2025-09-01T13:11:20.189677367Z","level":"INFO","msg":"🎉 TOQ Server inicializado com sucesso","component":"bootstrap","total_time":"133.35544ms"}
{"time":"2025-09-01T13:11:53.00710102Z","level":"INFO","msg":"Role assigned to user successfully","userID":2,"roleID":3,"roleName":"Proprietário"}
{"time":"2025-09-01T13:11:53.088106971Z","level":"INFO","msg":"user folder structure created successfully in S3","userID":2,"bucket":"toq-app-media"}
{"time":"2025-09-01T13:11:53.088215064Z","level":"ERROR","msg":"User has no active role","user_id":2}
{"time":"2025-09-01T13:11:53.091352842Z","level":"WARN","source":{"function":"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/middlewares.StructuredLoggingMiddleware.func1","file":"/codigos/go_code/toq_server/internal/adapter/left/http/middlewares/structured_logging_middleware.go","line":122},"msg":"HTTP Error","request_id":"e535f8c4-b77e-4dbc-b9aa-f814f0b21e6b","method":"POST","path":"/api/v1/auth/owner","status":409,"duration":397561326,"size":76,"client_ip":"179.110.194.42","user_agent":"PostmanRuntime/7.45.0"}
