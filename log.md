Starting: /home/toq_admin/go/bin/dlv dap --listen=127.0.0.1:37741 --log-dest=3 from /codigos/go_code/toq_server/cmd
DAP server listening at: 127.0.0.1:37741
Type 'dlv help' for list of commands.
{"time":"2025-10-23T15:25:33.091910054Z","level":"INFO","msg":"🚀 Iniciando TOQ Server Bootstrap","version":"2.0.0","component":"bootstrap","log_level":"info","log_format":"json","log_output":"stdout"}
{"time":"2025-10-23T15:25:33.092097579Z","level":"INFO","msg":"▶️ Executando fase","phase":"Phase01_InitializeContext","component":"bootstrap","timestamp":"2025-10-23T15:25:33Z"}
{"time":"2025-10-23T15:25:33.092157291Z","level":"INFO","msg":"🎯 FASE 1: Inicialização de Contexto e Sinais"}
{"time":"2025-10-23T15:25:33.092764858Z","level":"INFO","msg":"✅ Contexto e sinais inicializados com sucesso"}
{"time":"2025-10-23T15:25:33.092793729Z","level":"INFO","msg":"✅ Fase concluída","phase":"Phase01_InitializeContext","component":"bootstrap","duration":"698.19µs"}
{"time":"2025-10-23T15:25:33.092811729Z","level":"INFO","msg":"▶️ Executando fase","phase":"Phase02_LoadConfiguration","component":"bootstrap","timestamp":"2025-10-23T15:25:33Z"}
{"time":"2025-10-23T15:25:33.09284353Z","level":"INFO","msg":"🎯 FASE 2: Carregamento e Validação de Configuração"}
{"time":"2025-10-23T15:25:33.092894732Z","level":"INFO","msg":"🔍 Iniciando servidor pprof na porta 6060"}
{"time":"2025-10-23T15:25:33.093003665Z","level":"INFO","msg":"✅ Servidor pprof iniciado em localhost:6060"}
{"time":"2025-10-23T15:25:33.094872458Z","level":"INFO","msg":"Configuration loaded successfully from YAML","path":"configs/env.yaml"}
{"time":"2025-10-23T15:25:33.094913189Z","level":"INFO","msg":"Overrides de ambiente aplicados","environment":"dev","http_port":"127.0.0.1:18080","workers_enabled":false,"telemetry_endpoint":"localhost:4318"}
{"time":"2025-10-23T15:25:33.09496526Z","level":"INFO","msg":"🔐 JWT and token TTL configuration applied"}
{"time":"2025-10-23T15:25:33.09498106Z","level":"INFO","msg":"✅ Configuração carregada e validada com sucesso","version":"2.0.0"}
{"time":"2025-10-23T15:25:33.095027992Z","level":"INFO","msg":"✅ Fase concluída","phase":"Phase02_LoadConfiguration","component":"bootstrap","duration":"2.213202ms"}
{"time":"2025-10-23T15:25:33.095050352Z","level":"INFO","msg":"▶️ Executando fase","phase":"Phase03_InitializeInfrastructure","component":"bootstrap","timestamp":"2025-10-23T15:25:33Z"}
{"time":"2025-10-23T15:25:33.095096774Z","level":"INFO","msg":"🎯 FASE 3: Inicialização da Infraestrutura Core"}
{"time":"2025-10-23T15:25:33.095183036Z","level":"INFO","msg":"Database connection opened successfully"}
{"time":"2025-10-23T15:25:33.095199997Z","level":"INFO","msg":"✅ Conexão com banco de dados estabelecida"}
{"time":"2025-10-23T15:25:33.100398493Z","level":"INFO","msg":"Redis cache connected successfully","url":"redis://localhost:6379/0"}
{"time":"2025-10-23T15:25:33.100458355Z","level":"INFO","msg":"✅ Sistema de cache Redis inicializado com sucesso"}
{"time":"2025-10-23T15:25:33.100897917Z","level":"INFO","msg":"OpenTelemetry tracing initialized","endpoint":"localhost:4318"}
{"time":"2025-10-23T15:25:33.101389202Z","level":"INFO","msg":"OpenTelemetry metrics initialized","endpoint":"localhost:4318"}
{"time":"2025-10-23T15:25:33.104557701Z","level":"INFO","msg":"OpenTelemetry logs initialized","endpoint":"localhost:4318"}
{"time":"2025-10-23T15:25:33.104687744Z","level":"INFO","msg":"OpenTelemetry initialized successfully","tracing_enabled":true,"metrics_enabled":true,"endpoint":"localhost:4318"}
{"time":"2025-10-23T15:25:33.104764306Z","level":"INFO","msg":"✅ OpenTelemetry inicializado (tracing + metrics)"}
{"time":"2025-10-23T15:25:33.104819288Z","level":"INFO","msg":"Creating metrics adapter"}
{"time":"2025-10-23T15:25:33.107850493Z","level":"INFO","msg":"✅ Adapter de métricas Prometheus inicializado"}
{"time":"2025-10-23T15:25:33.107898525Z","level":"INFO","msg":"✅ Infraestrutura core inicializada com sucesso"}
{"time":"2025-10-23T15:25:33.107944996Z","level":"INFO","msg":"✅ Fase concluída","phase":"Phase03_InitializeInfrastructure","component":"bootstrap","duration":"12.890924ms"}
{"time":"2025-10-23T15:25:33.107978118Z","level":"INFO","msg":"▶️ Executando fase","phase":"Phase04_InjectDependencies","component":"bootstrap","timestamp":"2025-10-23T15:25:33Z"}
{"time":"2025-10-23T15:25:33.108036999Z","level":"INFO","msg":"🎯 FASE 4: Injeção de Dependências via Factory Pattern"}
{"time":"2025-10-23T15:25:33.10806483Z","level":"INFO","msg":"Starting dependency injection using Factory Pattern"}
{"time":"2025-10-23T15:25:33.108153092Z","level":"INFO","msg":"Creating storage adapters"}
{"time":"2025-10-23T15:25:33.111110646Z","level":"INFO","msg":"Redis cache connected successfully","url":"redis://localhost:6379/0"}
{"time":"2025-10-23T15:25:33.111231469Z","level":"INFO","msg":"Successfully created all storage adapters"}
{"time":"2025-10-23T15:25:33.111260429Z","level":"INFO","msg":"✅ ActivityTracker criado com sucesso com Redis client"}
{"time":"2025-10-23T15:25:33.11128648Z","level":"INFO","msg":"Creating repository adapters"}
{"time":"2025-10-23T15:25:33.111370782Z","level":"INFO","msg":"Successfully created all repository adapters"}
{"time":"2025-10-23T15:25:33.111397963Z","level":"INFO","msg":"Assigning repository adapters"}
{"time":"2025-10-23T15:25:33.111425875Z","level":"INFO","msg":"Repository adapters assigned successfully"}
{"time":"2025-10-23T15:25:33.111451405Z","level":"INFO","msg":"Creating validation adapters"}
{"time":"2025-10-23T15:25:33.111488466Z","level":"INFO","msg":"Successfully created all validation adapters"}
{"time":"2025-10-23T15:25:33.111512557Z","level":"INFO","msg":"Creating external service adapters"}
{"time":"2025-10-23T15:25:33.11163321Z","level":"INFO","msg":"Creating external service adapters"}
{"time":"2025-10-23T15:25:33.113182443Z","level":"INFO","msg":"fcm.adapter.initialized","request_id":"server_initialization","trace_id":"bc896ce687b005530a3ef64986c2d81d","span_id":"996900e26ed04502"}
{"time":"2025-10-23T15:25:33.113745799Z","level":"INFO","msg":"adapter.s3.creating","request_id":"server_initialization","region":"us-east-1","bucket":"toq-app-media"}
{"time":"2025-10-23T15:25:33.114460729Z","level":"INFO","msg":"adapter.s3.created","request_id":"server_initialization","bucket":"toq-app-media","region":"us-east-1"}
{"time":"2025-10-23T15:25:33.114638535Z","level":"INFO","msg":"Successfully created all external service adapters"}
{"time":"2025-10-23T15:25:33.114674396Z","level":"INFO","msg":"Initializing all services"}
{"time":"2025-10-23T15:25:33.114748938Z","level":"INFO","msg":"All services initialized successfully"}
{"time":"2025-10-23T15:25:33.114779608Z","level":"INFO","msg":"Workers desabilitados; TempBlockCleaner não será inicializado","environment":"dev"}
{"time":"2025-10-23T15:25:33.114808689Z","level":"INFO","msg":"Dependency injection completed successfully using Factory Pattern"}
{"time":"2025-10-23T15:25:33.11483866Z","level":"INFO","msg":"✅ Injeção de dependências concluída via Factory Pattern"}
{"time":"2025-10-23T15:25:33.114913733Z","level":"INFO","msg":"✅ Fase concluída","phase":"Phase04_InjectDependencies","component":"bootstrap","duration":"6.932505ms"}
{"time":"2025-10-23T15:25:33.114980504Z","level":"INFO","msg":"▶️ Executando fase","phase":"Phase05_InitializeServices","component":"bootstrap","timestamp":"2025-10-23T15:25:33Z"}
{"time":"2025-10-23T15:25:33.115044506Z","level":"INFO","msg":"🎯 FASE 5: Inicialização de Serviços"}
{"time":"2025-10-23T15:25:33.115126478Z","level":"INFO","msg":"✅ Serviço inicializado","service":"GlobalService"}
{"time":"2025-10-23T15:25:33.115149299Z","level":"INFO","msg":"✅ Serviço inicializado","service":"PermissionService"}
{"time":"2025-10-23T15:25:33.115161859Z","level":"INFO","msg":"✅ Serviço inicializado","service":"HolidayService"}
{"time":"2025-10-23T15:25:33.115172249Z","level":"INFO","msg":"✅ Serviço inicializado","service":"PhotoSessionService"}
{"time":"2025-10-23T15:25:33.11518411Z","level":"INFO","msg":"✅ Serviço inicializado","service":"ComplexService"}
{"time":"2025-10-23T15:25:33.115197821Z","level":"INFO","msg":"✅ Serviço inicializado","service":"ScheduleService"}
{"time":"2025-10-23T15:25:33.115212491Z","level":"INFO","msg":"✅ Serviço inicializado","service":"ListingService"}
{"time":"2025-10-23T15:25:33.115225262Z","level":"INFO","msg":"✅ Serviço inicializado","service":"UserService"}
{"time":"2025-10-23T15:25:33.115235452Z","level":"INFO","msg":"✅ Todos os serviços inicializados com sucesso"}
{"time":"2025-10-23T15:25:33.115288533Z","level":"INFO","msg":"✅ Fase concluída","phase":"Phase05_InitializeServices","component":"bootstrap","duration":"311.669µs"}
{"time":"2025-10-23T15:25:33.115314874Z","level":"INFO","msg":"▶️ Executando fase","phase":"Phase06_ConfigureHandlers","component":"bootstrap","timestamp":"2025-10-23T15:25:33Z"}
{"time":"2025-10-23T15:25:33.115394766Z","level":"INFO","msg":"🎯 FASE 6: Configuração de Handlers e Rotas"}
{"time":"2025-10-23T15:25:33.115716465Z","level":"INFO","msg":"✅ Servidor HTTP configurado com TLS e middleware"}
{"time":"2025-10-23T15:25:33.115735905Z","level":"INFO","msg":"✅ Handlers HTTP preparados para criação"}
{"time":"2025-10-23T15:25:33.115748946Z","level":"INFO","msg":"Creating HTTP handlers"}
{"time":"2025-10-23T15:25:33.115502609Z","level":"WARN","msg":"Failed to parse IdleTimeout, using default","value":"","error":"time: invalid duration \"\""}
{"time":"2025-10-23T15:25:33.115799247Z","level":"INFO","msg":"Successfully created all HTTP handlers"}
{"time":"2025-10-23T15:25:33.116922449Z","level":"INFO","msg":"✅ Rotas e middlewares configurados"}
{"time":"2025-10-23T15:25:33.117015051Z","level":"INFO","msg":"✅ Health checks configurados"}
{"time":"2025-10-23T15:25:33.117032322Z","level":"INFO","msg":"✅ Handlers e rotas configurados com sucesso"}
{"time":"2025-10-23T15:25:33.117094924Z","level":"INFO","msg":"✅ Fase concluída","phase":"Phase06_ConfigureHandlers","component":"bootstrap","duration":"1.77545ms"}
{"time":"2025-10-23T15:25:33.117155686Z","level":"INFO","msg":"▶️ Executando fase","phase":"Phase07_StartBackgroundWorkers","component":"bootstrap","timestamp":"2025-10-23T15:25:33Z"}
{"time":"2025-10-23T15:25:33.117209147Z","level":"INFO","msg":"🎯 FASE 7: Inicialização de Background Workers"}
{"time":"2025-10-23T15:25:33.117266909Z","level":"INFO","msg":"Workers desabilitados para o ambiente atual; fase 7 vai pular inicialização","environment":"dev"}
{"time":"2025-10-23T15:25:33.120594383Z","level":"INFO","msg":"Database connection verified"}
{"time":"2025-10-23T15:25:33.120658985Z","level":"INFO","msg":"✅ Schema do banco de dados verificado"}
{"time":"2025-10-23T15:25:33.120671015Z","level":"INFO","msg":"✅ Background workers inicializados com sucesso"}
{"time":"2025-10-23T15:25:33.120712516Z","level":"INFO","msg":"✅ Fase concluída","phase":"Phase07_StartBackgroundWorkers","component":"bootstrap","duration":"3.585391ms"}
{"time":"2025-10-23T15:25:33.120739657Z","level":"INFO","msg":"▶️ Executando fase","phase":"Phase08_StartServer","component":"bootstrap","timestamp":"2025-10-23T15:25:33Z"}
{"time":"2025-10-23T15:25:33.120786048Z","level":"INFO","msg":"🎯 FASE 8: Inicialização Final e Runtime"}
{"time":"2025-10-23T15:25:33.120821019Z","level":"INFO","msg":"✅ Servidor marcado como ready para receber tráfego"}
{"time":"2025-10-23T15:25:33.120863971Z","level":"INFO","msg":"🚀 Iniciando servidor HTTP na porta configurada"}
{"time":"2025-10-23T15:25:33.221483637Z","level":"INFO","msg":"✅ Servidor HTTP iniciado com sucesso"}
{"time":"2025-10-23T15:25:33.221554279Z","level":"INFO","msg":"✅ Monitoramento de saúde em runtime iniciado"}
{"time":"2025-10-23T15:25:33.22160755Z","level":"INFO","msg":"🌟 TOQ Server pronto para servir","uptime":129525941}
{"time":"2025-10-23T15:25:33.221635981Z","level":"INFO","msg":"✅ Fase concluída","phase":"Phase08_StartServer","component":"bootstrap","duration":"100.894594ms"}
{"time":"2025-10-23T15:25:33.221658251Z","level":"INFO","msg":"🎉 TOQ Server inicializado com sucesso","component":"bootstrap","total_time":"129.577622ms"}
{"time":"2025-10-23T15:25:54.109802057Z","level":"INFO","msg":"Request received","method":"POST","path":"/api/v2/auth/validate/cep","remote_addr":"127.0.0.1:57278"}
{"time":"2025-10-23T15:25:54.110392664Z","level":"ERROR","msg":"auth.validate.signature.validator_missing","request_id":"3ae4fb91-e4aa-4bc7-bfcc-fce3334cf47c"}
{"time":"2025-10-23T15:25:54.110639381Z","level":"ERROR","msg":"HTTP Error","request_id":"3ae4fb91-e4aa-4bc7-bfcc-fce3334cf47c","method":"POST","path":"/api/v2/auth/validate/cep","status":500,"duration":631638,"size":75,"client_ip":"87.9.226.126","user_agent":"PostmanRuntime/7.49.0","trace_id":"4b8b39ffdf64260e4a2e14cbb796d8e3","span_id":"d784af4a33e272fe","function":"github.com/projeto-toq/toq_server/internal/core/utils.InternalError","file":"/codigos/go_code/toq_server/internal/core/utils/http_errors.go","line":246,"stack":["github.com/projeto-toq/toq_server/internal/core/utils.InternalError (http_errors.go:246)"],"error_code":500,"error_message":"Signature validator not configured.","errors":["HTTP 500: Signature validator not configured."]}
{"time":"2025-10-23T15:34:53.003517391Z","level":"INFO","msg":"Request received","method":"GET","path":"/","remote_addr":"127.0.0.1:36172"}
{"time":"2025-10-23T15:34:53.003747767Z","level":"INFO","msg":"HTTP Response","request_id":"97adff3a-8023-458c-83d4-dec4f16156d8","method":"GET","path":"/","status":404,"duration":94282,"size":-1,"client_ip":"206.168.34.196","user_agent":"Mozilla/5.0 (compatible; CensysInspect/1.1; +https://about.censys.io/)","trace_id":"1d02755d401f3c62bf66439390790944","span_id":"2c60296fef1f96fd"}
Detaching and terminating target process
dlv dap (475127) exited with code: 0
