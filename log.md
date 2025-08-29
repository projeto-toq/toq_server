Starting: /home/toq_admin/go/bin/dlv dap --listen=127.0.0.1:36939 --log-dest=3 from /codigos/go_code/toq_server/cmd
DAP server listening at: 127.0.0.1:36939
Type 'dlv help' for list of commands.
2025/08/29 21:04:29 INFO 🚀 Iniciando TOQ Server Bootstrap version=2.0.0 timestamp=2025-08-29T21:04:29Z
2025/08/29 21:04:29 INFO ▶️ Executando fase phase=Phase01_InitializeContext timestamp=2025-08-29T21:04:29Z
2025/08/29 21:04:29 INFO 🎯 FASE 1: Inicialização de Contexto e Sinais
2025/08/29 21:04:29 INFO ✅ Contexto e sinais inicializados com sucesso
2025/08/29 21:04:29 INFO ✅ Fase concluída phase=Phase01_InitializeContext duration=1.427151ms
2025/08/29 21:04:29 INFO ▶️ Executando fase phase=Phase02_LoadConfiguration timestamp=2025-08-29T21:04:29Z
2025/08/29 21:04:29 INFO 🔍 Iniciando servidor pprof na porta 6060
2025/08/29 21:04:29 INFO ✅ Servidor pprof iniciado em localhost:6060
2025/08/29 21:04:29 INFO 🎯 FASE 2: Carregamento e Validação de Configuração
2025/08/29 21:04:29 INFO Configuration loaded successfully from YAML path=configs/env.yaml
time=2025-08-29T21:04:29.923Z level=INFO msg="Logging system initialized" level=INFO to_file=false add_source=false
time=2025-08-29T21:04:29.923Z level=INFO msg="INFO ✅ Logging inicial baseado em ENV configurado"
time=2025-08-29T21:04:29.923Z level=INFO msg="Logging system initialized" level=INFO to_file=false add_source=false
time=2025-08-29T21:04:29.923Z level=INFO msg="INFO ✅ Logging reconfigurado com prioridade ENV > YAML > defaults"
time=2025-08-29T21:04:29.923Z level=INFO msg="INFO ✅ Configuração carregada e validada com sucesso version=2.0.0"
time=2025-08-29T21:04:29.924Z level=INFO msg="INFO ✅ Fase concluída phase=Phase02_LoadConfiguration duration=2.490007ms"
time=2025-08-29T21:04:29.924Z level=INFO msg="INFO ▶️ Executando fase phase=Phase03_InitializeInfrastructure timestamp=2025-08-29T21:04:29Z"
time=2025-08-29T21:04:29.924Z level=INFO msg="INFO 🎯 FASE 3: Inicialização da Infraestrutura Core"
time=2025-08-29T21:04:29.926Z level=INFO msg="Database connection initialized" uri="toq_user:toq_password@tcp(localhost:3306)/toq_db?parseTime=true&loc=UTC&timeout=30s&readTimeout=30s&writeTimeout=30s"
time=2025-08-29T21:04:29.926Z level=INFO msg="INFO ✅ Conexão com banco de dados estabelecida"
time=2025-08-29T21:04:29.929Z level=INFO msg="Redis cache connected successfully" url=redis://localhost:6379/0
time=2025-08-29T21:04:29.929Z level=INFO msg="INFO ✅ Sistema de cache Redis inicializado com sucesso"
time=2025-08-29T21:04:29.929Z level=INFO msg="OpenTelemetry initialization placeholder - not fully implemented" enabled=true otlp_enabled=true endpoint=http://localhost:14318
time=2025-08-29T21:04:29.929Z level=INFO msg="INFO ✅ OpenTelemetry inicializado (tracing + metrics)"
time=2025-08-29T21:04:29.929Z level=ERROR msg="ActivityTracker não foi criado na Phase 04 - falha na inicialização"
time=2025-08-29T21:04:29.929Z level=INFO msg="ERROR ❌ Fase falhou phase=Phase03_InitializeInfrastructure duration=5.471833ms error=\"[Phase03] activity_tracker: Failed to initialize activity tracker (failed to initialize activity tracker: ActivityTracker não foi inicializado)\""
time=2025-08-29T21:04:29.929Z level=INFO msg="ERROR ❌ Falha na fase de inicialização phase=Phase03_InitializeInfrastructure error=\"[Phase03] activity_tracker: Failed to initialize activity tracker (failed to initialize activity tracker: ActivityTracker não foi inicializado)\""
time=2025-08-29T21:04:29.929Z level=ERROR msg="❌ Falha crítica durante inicialização" error="bootstrap failed at phase Phase03_InitializeInfrastructure: [Phase03] activity_tracker: Failed to initialize activity tracker (failed to initialize activity tracker: ActivityTracker não foi inicializado)"
Process 115761 has exited with status 1
Detaching
dlv dap (115679) exited with code: 0
