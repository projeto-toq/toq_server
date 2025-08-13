Starting: /home/giulio/go/bin/dlv dap --listen=127.0.0.1:36279 --log-dest=3 from /home/giulio/go-code/toq_server/cmd
DAP server listening at: 127.0.0.1:36279
Type 'dlv help' for list of commands.
2025/08/13 14:31:13 INFO 🚀 Starting TOQ Server initialization version=2.0.0
2025/08/13 14:31:13 INFO ✅ Environment configuration loaded successfully
LOG.LEVEL setted to DEBUG
{"time":"2025-08-13T14:31:13.924129331-03:00","level":"DEBUG","msg":"log configured to console"}
{"time":"2025-08-13T14:31:13.924139177-03:00","level":"INFO","msg":"✅ Logging system initialized"}
{"time":"2025-08-13T14:31:13.924159984-03:00","level":"INFO","msg":"🔧 TOQ API Server starting","version":"2.0.0"}
{"time":"2025-08-13T14:31:13.925365066-03:00","level":"INFO","msg":"Database answered the ping. MySql connection successfuly!"}
{"time":"2025-08-13T14:31:13.94455315-03:00","level":"INFO","msg":"✅ Database connection established"}
{"time":"2025-08-13T14:31:13.946081284-03:00","level":"INFO","msg":"✅ OpenTelemetry initialized (tracing + metrics)"}
{"time":"2025-08-13T14:31:13.949176146-03:00","level":"INFO","msg":"Activity tracker initialized successfully"}
{"time":"2025-08-13T14:31:13.949228723-03:00","level":"INFO","msg":"✅ Activity tracking system initialized"}
{"time":"2025-08-13T14:31:13.949356792-03:00","level":"INFO","msg":"Server listening on","Addr:":{"IP":"::","Port":50051,"Zone":""}}
{"time":"2025-08-13T14:31:13.952682328-03:00","level":"INFO","msg":"✅ gRPC server configured with TLS and interceptors"}
{"time":"2025-08-13T14:31:13.952730547-03:00","level":"INFO","msg":"Starting dependency injection using Factory Pattern"}
{"time":"2025-08-13T14:31:13.952773576-03:00","level":"DEBUG","msg":"Validating factory configuration"}
{"time":"2025-08-13T14:31:13.952818229-03:00","level":"DEBUG","msg":"Factory configuration validation successful"}
{"time":"2025-08-13T14:31:13.952840136-03:00","level":"INFO","msg":"Creating storage adapters"}
{"time":"2025-08-13T14:31:13.952863586-03:00","level":"INFO","msg":"Creating storage adapters"}
{"time":"2025-08-13T14:31:13.962024083-03:00","level":"INFO","msg":"Redis cache connected successfully","url":"redis://localhost:6379/0"}
{"time":"2025-08-13T14:31:13.962704614-03:00","level":"INFO","msg":"Successfully created all storage adapters"}
{"time":"2025-08-13T14:31:13.962724263-03:00","level":"DEBUG","msg":"Assigning storage adapters to config"}
{"time":"2025-08-13T14:31:13.962737217-03:00","level":"INFO","msg":"Creating repository adapters"}
{"time":"2025-08-13T14:31:13.962750012-03:00","level":"INFO","msg":"Creating repository adapters"}
{"time":"2025-08-13T14:31:13.962812078-03:00","level":"INFO","msg":"Successfully created all repository adapters"}
{"time":"2025-08-13T14:31:13.96283131-03:00","level":"DEBUG","msg":"Assigning repository adapters to config"}
{"time":"2025-08-13T14:31:13.96284493-03:00","level":"INFO","msg":"Creating validation adapters"}
{"time":"2025-08-13T14:31:13.962857898-03:00","level":"INFO","msg":"Creating validation adapters"}
{"time":"2025-08-13T14:31:13.96287617-03:00","level":"INFO","msg":"Successfully created all validation adapters"}
{"time":"2025-08-13T14:31:13.962890132-03:00","level":"DEBUG","msg":"Assigning validation adapters to config"}
{"time":"2025-08-13T14:31:13.962903439-03:00","level":"INFO","msg":"Creating external service adapters"}
{"time":"2025-08-13T14:31:13.962916553-03:00","level":"INFO","msg":"Creating external service adapters"}
{"time":"2025-08-13T14:31:13.966586007-03:00","level":"INFO","msg":"Successfully created all external service adapters"}
{"time":"2025-08-13T14:31:13.966611922-03:00","level":"DEBUG","msg":"Assigning external service adapters to config"}
{"time":"2025-08-13T14:31:13.966622426-03:00","level":"INFO","msg":"Initializing services"}
{"time":"2025-08-13T14:31:13.966631446-03:00","level":"DEBUG","msg":"Initializing Global Service"}
{"time":"2025-08-13T14:31:13.966642483-03:00","level":"DEBUG","msg":"GlobalService injected into RedisCache"}
{"time":"2025-08-13T14:31:13.96665102-03:00","level":"DEBUG","msg":"GlobalService injected into Redis cache"}
{"time":"2025-08-13T14:31:13.966673448-03:00","level":"DEBUG","msg":"Initializing Complex Handler"}
{"time":"2025-08-13T14:31:13.966680392-03:00","level":"DEBUG","msg":"Initializing Listing Handler"}
{"time":"2025-08-13T14:31:13.966713676-03:00","level":"DEBUG","msg":"Initializing User Handler"}
{"time":"2025-08-13T14:31:13.966747501-03:00","level":"INFO","msg":"All services initialized successfully"}
{"time":"2025-08-13T14:31:13.966756547-03:00","level":"INFO","msg":"Dependency injection completed successfully using Factory Pattern"}
{"time":"2025-08-13T14:31:13.966764327-03:00","level":"INFO","msg":"✅ Dependency injection completed via Factory Pattern"}
{"time":"2025-08-13T14:31:13.966772461-03:00","level":"INFO","msg":"✅ Activity tracker linked with user service"}
{"time":"2025-08-13T14:31:13.97221988-03:00","level":"INFO","msg":"✅ Database schema verified"}
{"time":"2025-08-13T14:31:13.972256541-03:00","level":"INFO","msg":"✅ Background workers initialized"}
{"time":"2025-08-13T14:31:13.972271289-03:00","level":"INFO","msg":"CRECI validation routine started"}
{"time":"2025-08-13T14:31:13.9722879-03:00","level":"INFO","msg":"memory cache cleaner routine started"}
{"time":"2025-08-13T14:31:13.972290134-03:00","level":"INFO","msg":"🌟 TOQ Server ready to serve","services":4,"methods":69,"version":"2.0.0"}
{"time":"2025-08-13T14:31:13.972332776-03:00","level":"INFO","msg":"Activity batch worker started","interval":30000000000}
{"time":"2025-08-13T14:31:13.972338634-03:00","level":"INFO","msg":"🚀 Starting gRPC server"}
{"time":"2025-08-13T14:31:13.972344302-03:00","level":"INFO","msg":"session cleaner routine started"}
{"time":"2025-08-13T14:31:43.973964628-03:00","level":"DEBUG","msg":"Processing active users","count":1}
{"time":"2025-08-13T14:31:43.976913799-03:00","level":"DEBUG","msg":"Batch updated user activities","affected_rows":0,"batch_size":1}
{"time":"2025-08-13T14:31:43.977144754-03:00","level":"DEBUG","msg":"Successfully updated user activities","count":1}
{"time":"2025-08-13T14:31:46.712245153-03:00","level":"INFO","msg":"Request received:","Method:":"/grpc.UserService/SignIn"}
{"time":"2025-08-13T14:31:46.712364129-03:00","level":"DEBUG","msg":"Cache lookup","key":"toq_cache:method:/grpc.UserService/SignIn:role:0","method":"/grpc.UserService/SignIn","role":0}
{"time":"2025-08-13T14:31:46.713372968-03:00","level":"DEBUG","msg":"Cache miss","key":"toq_cache:method:/grpc.UserService/SignIn:role:0"}
{"time":"2025-08-13T14:31:46.713448779-03:00","level":"DEBUG","msg":"Fetching from service","service":0,"method":3,"role":0}
{"time":"2025-08-13T14:31:46.714566738-03:00","level":"DEBUG","msg":"LoadGRPCAccess called","service":0,"method":3,"role":0}
{"time":"2025-08-13T14:31:46.716340705-03:00","level":"DEBUG","msg":"LoadGRPCAccess query result","role":0,"entities_count":33}
{"time":"2025-08-13T14:31:46.716380135-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":0,"allowed":true,"target_service":0,"target_method":3}
{"time":"2025-08-13T14:31:46.716421958-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":1,"allowed":true,"target_service":0,"target_method":3}
{"time":"2025-08-13T14:31:46.716444603-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":2,"allowed":true,"target_service":0,"target_method":3}
{"time":"2025-08-13T14:31:46.716455434-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":3,"allowed":true,"target_service":0,"target_method":3}
{"time":"2025-08-13T14:31:46.716465502-03:00","level":"DEBUG","msg":"Found matching privilege","allowed":true}
{"time":"2025-08-13T14:31:46.716751253-03:00","level":"DEBUG","msg":"Privilege fetched from service","service":0,"method":3,"role":0,"allowed":true}
{"time":"2025-08-13T14:31:46.717103347-03:00","level":"DEBUG","msg":"Cache entry stored","key":"toq_cache:method:/grpc.UserService/SignIn:role:0","ttl":900000000000}
{"time":"2025-08-13T14:31:46.717125162-03:00","level":"DEBUG","msg":"Permission check result","method":"/grpc.UserService/SignIn","role":0,"allowed":true,"valid":true}
{"time":"2025-08-13T14:31:46.928681905-03:00","level":"INFO","msg":"Request received:","Method:":"/grpc.UserService/GetProfile"}
{"time":"2025-08-13T14:31:46.929811-03:00","level":"DEBUG","msg":"Cache lookup","key":"toq_cache:method:/grpc.UserService/GetProfile:role:1","method":"/grpc.UserService/GetProfile","role":1}
{"time":"2025-08-13T14:31:46.930465544-03:00","level":"DEBUG","msg":"Cache miss","key":"toq_cache:method:/grpc.UserService/GetProfile:role:1"}
{"time":"2025-08-13T14:31:46.930535194-03:00","level":"DEBUG","msg":"Fetching from service","service":0,"method":13,"role":1}
{"time":"2025-08-13T14:31:46.931640512-03:00","level":"DEBUG","msg":"LoadGRPCAccess called","service":0,"method":13,"role":1}
{"time":"2025-08-13T14:31:46.934167157-03:00","level":"DEBUG","msg":"LoadGRPCAccess query result","role":1,"entities_count":66}
{"time":"2025-08-13T14:31:46.934363851-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":0,"allowed":true,"target_service":0,"target_method":13}
{"time":"2025-08-13T14:31:46.934439933-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":1,"allowed":false,"target_service":0,"target_method":13}
{"time":"2025-08-13T14:31:46.934486929-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":2,"allowed":false,"target_service":0,"target_method":13}
{"time":"2025-08-13T14:31:46.9345305-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":3,"allowed":true,"target_service":0,"target_method":13}
{"time":"2025-08-13T14:31:46.9345784-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":4,"allowed":true,"target_service":0,"target_method":13}
{"time":"2025-08-13T14:31:46.934622943-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":5,"allowed":true,"target_service":0,"target_method":13}
{"time":"2025-08-13T14:31:46.934666444-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":6,"allowed":true,"target_service":0,"target_method":13}
{"time":"2025-08-13T14:31:46.934709983-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":7,"allowed":true,"target_service":0,"target_method":13}
{"time":"2025-08-13T14:31:46.934752951-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":8,"allowed":true,"target_service":0,"target_method":13}
{"time":"2025-08-13T14:31:46.934961707-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":9,"allowed":true,"target_service":0,"target_method":13}
{"time":"2025-08-13T14:31:46.935028341-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":10,"allowed":true,"target_service":0,"target_method":13}
{"time":"2025-08-13T14:31:46.935099324-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":11,"allowed":true,"target_service":0,"target_method":13}
{"time":"2025-08-13T14:31:46.93516102-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":12,"allowed":false,"target_service":0,"target_method":13}
{"time":"2025-08-13T14:31:46.93522485-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":13,"allowed":true,"target_service":0,"target_method":13}
{"time":"2025-08-13T14:31:46.935294141-03:00","level":"DEBUG","msg":"Found matching privilege","allowed":true}
{"time":"2025-08-13T14:31:46.936273406-03:00","level":"DEBUG","msg":"Privilege fetched from service","service":0,"method":13,"role":1,"allowed":true}
{"time":"2025-08-13T14:31:46.937079844-03:00","level":"DEBUG","msg":"Cache entry stored","key":"toq_cache:method:/grpc.UserService/GetProfile:role:1","ttl":900000000000}
{"time":"2025-08-13T14:31:46.937136626-03:00","level":"DEBUG","msg":"Permission check result","method":"/grpc.UserService/GetProfile","role":1,"allowed":true,"valid":true}
{"time":"2025-08-13T14:31:47.027671561-03:00","level":"INFO","msg":"Request received:","Method:":"/grpc.UserService/UpdateOptStatus"}
{"time":"2025-08-13T14:31:47.028841603-03:00","level":"DEBUG","msg":"Cache lookup","key":"toq_cache:method:/grpc.UserService/UpdateOptStatus:role:1","method":"/grpc.UserService/UpdateOptStatus","role":1}
{"time":"2025-08-13T14:31:47.029371417-03:00","level":"DEBUG","msg":"Cache miss","key":"toq_cache:method:/grpc.UserService/UpdateOptStatus:role:1"}
{"time":"2025-08-13T14:31:47.029431802-03:00","level":"DEBUG","msg":"Fetching from service","service":0,"method":29,"role":1}
{"time":"2025-08-13T14:31:47.030469918-03:00","level":"DEBUG","msg":"LoadGRPCAccess called","service":0,"method":29,"role":1}
{"time":"2025-08-13T14:31:47.033312688-03:00","level":"DEBUG","msg":"LoadGRPCAccess query result","role":1,"entities_count":66}
{"time":"2025-08-13T14:31:47.033420911-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":0,"allowed":true,"target_service":0,"target_method":29}
{"time":"2025-08-13T14:31:47.033504781-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":1,"allowed":false,"target_service":0,"target_method":29}
{"time":"2025-08-13T14:31:47.03357728-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":2,"allowed":false,"target_service":0,"target_method":29}
{"time":"2025-08-13T14:31:47.033637837-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":3,"allowed":true,"target_service":0,"target_method":29}
{"time":"2025-08-13T14:31:47.033692406-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":4,"allowed":true,"target_service":0,"target_method":29}
{"time":"2025-08-13T14:31:47.033752176-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":5,"allowed":true,"target_service":0,"target_method":29}
{"time":"2025-08-13T14:31:47.033817579-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":6,"allowed":true,"target_service":0,"target_method":29}
{"time":"2025-08-13T14:31:47.033878998-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":7,"allowed":true,"target_service":0,"target_method":29}
{"time":"2025-08-13T14:31:47.033942106-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":8,"allowed":true,"target_service":0,"target_method":29}
{"time":"2025-08-13T14:31:47.033992985-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":9,"allowed":true,"target_service":0,"target_method":29}
{"time":"2025-08-13T14:31:47.034050366-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":10,"allowed":true,"target_service":0,"target_method":29}
{"time":"2025-08-13T14:31:47.034103661-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":11,"allowed":true,"target_service":0,"target_method":29}
{"time":"2025-08-13T14:31:47.034160614-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":12,"allowed":false,"target_service":0,"target_method":29}
{"time":"2025-08-13T14:31:47.034219175-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":13,"allowed":true,"target_service":0,"target_method":29}
{"time":"2025-08-13T14:31:47.034266064-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":14,"allowed":true,"target_service":0,"target_method":29}
{"time":"2025-08-13T14:31:47.034327525-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":15,"allowed":true,"target_service":0,"target_method":29}
{"time":"2025-08-13T14:31:47.034378376-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":16,"allowed":true,"target_service":0,"target_method":29}
{"time":"2025-08-13T14:31:47.034433645-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":17,"allowed":true,"target_service":0,"target_method":29}
{"time":"2025-08-13T14:31:47.034491812-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":18,"allowed":true,"target_service":0,"target_method":29}
{"time":"2025-08-13T14:31:47.034539355-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":19,"allowed":true,"target_service":0,"target_method":29}
{"time":"2025-08-13T14:31:47.034600803-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":20,"allowed":true,"target_service":0,"target_method":29}
{"time":"2025-08-13T14:31:47.034650424-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":21,"allowed":true,"target_service":0,"target_method":29}
{"time":"2025-08-13T14:31:47.034819863-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":22,"allowed":false,"target_service":0,"target_method":29}
{"time":"2025-08-13T14:31:47.034876894-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":23,"allowed":false,"target_service":0,"target_method":29}
{"time":"2025-08-13T14:31:47.03493719-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":24,"allowed":false,"target_service":0,"target_method":29}
{"time":"2025-08-13T14:31:47.035000341-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":25,"allowed":false,"target_service":0,"target_method":29}
{"time":"2025-08-13T14:31:47.035061689-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":26,"allowed":false,"target_service":0,"target_method":29}
{"time":"2025-08-13T14:31:47.035120687-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":27,"allowed":false,"target_service":0,"target_method":29}
{"time":"2025-08-13T14:31:47.035181992-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":28,"allowed":false,"target_service":0,"target_method":29}
{"time":"2025-08-13T14:31:47.035241738-03:00","level":"DEBUG","msg":"Checking privilege","service":0,"method":29,"allowed":false,"target_service":0,"target_method":29}
{"time":"2025-08-13T14:31:47.035290627-03:00","level":"DEBUG","msg":"Found matching privilege","allowed":false}
{"time":"2025-08-13T14:31:47.036178791-03:00","level":"DEBUG","msg":"Privilege fetched from service","service":0,"method":29,"role":1,"allowed":false}
{"time":"2025-08-13T14:31:47.036892115-03:00","level":"DEBUG","msg":"Cache entry stored","key":"toq_cache:method:/grpc.UserService/UpdateOptStatus:role:1","ttl":900000000000}
{"time":"2025-08-13T14:31:47.036962913-03:00","level":"DEBUG","msg":"Permission check result","method":"/grpc.UserService/UpdateOptStatus","role":1,"allowed":false,"valid":true}
{"time":"2025-08-13T14:31:47.037038866-03:00","level":"WARN","msg":"Usuário não tem acesso a este RPC","method":"/grpc.UserService/UpdateOptStatus","role":1,"allowed":false,"valid":true}
{"time":"2025-08-13T14:32:13.972503424-03:00","level":"INFO","msg":"Creci validation routine ticked"}
{"time":"2025-08-13T14:32:13.972676191-03:00","level":"DEBUG","msg":"Cleaning cache","pattern":"toq_cache:*"}
{"time":"2025-08-13T14:32:13.974038798-03:00","level":"DEBUG","msg":"Processing active users","count":1}
{"time":"2025-08-13T14:32:13.974626481-03:00","level":"INFO","msg":"Cache cleaned","deleted_keys":3}
{"time":"2025-08-13T14:32:13.983127035-03:00","level":"DEBUG","msg":"Batch updated user activities","affected_rows":1,"batch_size":1}
{"time":"2025-08-13T14:32:13.983202929-03:00","level":"DEBUG","msg":"Successfully updated user activities","count":1}
