Starting: /home/toq_admin/go/bin/dlv dap --listen=127.0.0.1:34135 --log-dest=3 from /codigos/go_code/toq_server/cmd
DAP server listening at: 127.0.0.1:34135
Type 'dlv help' for list of commands.
{"time":"2025-08-31T12:14:55.232410532Z","level":"INFO","msg":"🚀 Iniciando TOQ Server Bootstrap","version":"2.0.0","component":"bootstrap","log_level":"info","log_format":"json","log_output":"stdout"}
{"time":"2025-08-31T12:14:55.232520642Z","level":"INFO","msg":"▶️ Executando fase","phase":"Phase01_InitializeContext","component":"bootstrap","timestamp":"2025-08-31T12:14:55Z"}
{"time":"2025-08-31T12:14:55.232586943Z","level":"INFO","msg":"🎯 FASE 1: Inicialização de Contexto e Sinais"}
{"time":"2025-08-31T12:14:55.233244947Z","level":"INFO","msg":"✅ Contexto e sinais inicializados com sucesso"}
{"time":"2025-08-31T12:14:55.233275467Z","level":"INFO","msg":"✅ Fase concluída","phase":"Phase01_InitializeContext","component":"bootstrap","duration":"755.925µs"}
{"time":"2025-08-31T12:14:55.233296777Z","level":"INFO","msg":"▶️ Executando fase","phase":"Phase02_LoadConfiguration","component":"bootstrap","timestamp":"2025-08-31T12:14:55Z"}
{"time":"2025-08-31T12:14:55.233328317Z","level":"INFO","msg":"🎯 FASE 2: Carregamento e Validação de Configuração"}
{"time":"2025-08-31T12:14:55.233913711Z","level":"INFO","msg":"🔍 Iniciando servidor pprof na porta 6060"}
{"time":"2025-08-31T12:14:55.234314203Z","level":"INFO","msg":"✅ Servidor pprof iniciado em localhost:6060"}
{"time":"2025-08-31T12:14:55.235378819Z","level":"INFO","msg":"Configuration loaded successfully from YAML","path":"configs/env.yaml"}
{"time":"2025-08-31T12:14:55.23547596Z","level":"INFO","msg":"✅ Configuração carregada e validada com sucesso","version":"2.0.0"}
{"time":"2025-08-31T12:14:55.2355215Z","level":"INFO","msg":"✅ Fase concluída","phase":"Phase02_LoadConfiguration","component":"bootstrap","duration":"2.222063ms"}
{"time":"2025-08-31T12:14:55.23554398Z","level":"INFO","msg":"▶️ Executando fase","phase":"Phase03_InitializeInfrastructure","component":"bootstrap","timestamp":"2025-08-31T12:14:55Z"}
{"time":"2025-08-31T12:14:55.23558104Z","level":"INFO","msg":"🎯 FASE 3: Inicialização da Infraestrutura Core"}
{"time":"2025-08-31T12:14:55.245360296Z","level":"INFO","msg":"Database connection initialized","uri":"toq_user:toq_password@tcp(localhost:3306)/toq_db?parseTime=true&loc=UTC&timeout=30s&readTimeout=30s&writeTimeout=30s"}
{"time":"2025-08-31T12:14:55.245401047Z","level":"INFO","msg":"✅ Conexão com banco de dados estabelecida"}
{"time":"2025-08-31T12:14:55.247426488Z","level":"INFO","msg":"Redis cache connected successfully","url":"redis://localhost:6379/0"}
{"time":"2025-08-31T12:14:55.247459748Z","level":"INFO","msg":"✅ Sistema de cache Redis inicializado com sucesso"}
{"time":"2025-08-31T12:14:55.248136002Z","level":"INFO","msg":"OpenTelemetry tracing initialized","endpoint":"http://localhost:14318"}
{"time":"2025-08-31T12:14:55.248388744Z","level":"ERROR","msg":"❌ Fase falhou","phase":"Phase03_InitializeInfrastructure","component":"bootstrap","duration":"12.834234ms","error":"[Phase03] telemetry: Failed to initialize OpenTelemetry (failed to initialize OpenTelemetry: failed to initialize telemetry: failed to initialize metrics: failed to create OTLP metric exporter: parse \"http://http:%2F%2Flocalhost:14318/v1/metrics\": invalid URL escape \"%2F\")"}
{"time":"2025-08-31T12:14:55.248440894Z","level":"ERROR","msg":"❌ Falha na fase de inicialização","phase":"Phase03_InitializeInfrastructure","error":"[Phase03] telemetry: Failed to initialize OpenTelemetry (failed to initialize OpenTelemetry: failed to initialize telemetry: failed to initialize metrics: failed to create OTLP metric exporter: parse \"http://http:%2F%2Flocalhost:14318/v1/metrics\": invalid URL escape \"%2F\")"}
{"time":"2025-08-31T12:14:55.248476224Z","level":"ERROR","msg":"❌ Falha crítica durante inicialização","error":"bootstrap failed at phase Phase03_InitializeInfrastructure: [Phase03] telemetry: Failed to initialize OpenTelemetry (failed to initialize OpenTelemetry: failed to initialize telemetry: failed to initialize metrics: failed to create OTLP metric exporter: parse \"http://http:%2F%2Flocalhost:14318/v1/metrics\": invalid URL escape \"%2F\")","component":"bootstrap"}
Process 5555 has exited with status 1
Detaching
dlv dap (3411) exited with code: 0
