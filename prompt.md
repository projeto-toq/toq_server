🛠️ Problema
Erro recebido ao tentar signin
{"time":"2025-08-30T13:59:16.800798592Z","level":"DEBUG","source":{"function":"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/middlewares.RequestIDMiddleware.func1","file":"/codigos/go_code/toq_server/internal/adapter/left/http/middlewares/request_id_middleware.go","line":30},"msg":"Request ID generated","request_id":"d58174aa-f989-42d3-9919-1f32e7f9a3e5","path":"/api/v1/auth/signin"}
{"time":"2025-08-30T13:59:16.803695567Z","level":"ERROR","source":{"function":"github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/user.(*UserAdapter).Read","file":"/codigos/go_code/toq_server/internal/adapter/right/mysql/user/basic_read.go","line":20},"msg":"Error preparing statement on mysqluseradapter/Read","error":"Error 1054 (42S22): Unknown column 'active' in 'where clause'"}
{"time":"2025-08-30T13:59:16.803911233Z","level":"ERROR","source":{"function":"github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/user.(*UserAdapter).GetActiveUserRolesByUserID","file":"/codigos/go_code/toq_server/internal/adapter/right/mysql/user/get_active_user_roles_by_userid.go","line":24},"msg":"mysqluseradapter/GetActiveUserRolesByUserID: error executing Read","error":"Error 1054 (42S22): Unknown column 'active' in 'where clause'"}
{"time":"2025-08-30T13:59:16.804602403Z","level":"WARN","source":{"function":"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/middlewares.StructuredLoggingMiddleware.func1","file":"/codigos/go_code/toq_server/internal/adapter/left/http/middlewares/structured_logging_middleware.go","line":107},"msg":"HTTP Error","request_id":"d58174aa-f989-42d3-9919-1f32e7f9a3e5","method":"POST","path":"/api/v1/auth/signin","status":401,"duration":3603066,"size":72,"client_ip":"179.110.194.42","user_agent":"PostmanRuntime/7.45.0"}

✅ Requisitos obrigatórios para qualquer revisão, refatoração ou correção
- Adoção das melhores práticas de desenvolvimento em Go
- Go Best Practices
- Google Go Style Guide
- Implementação seguindo arquitetura hexagonal
- Injeção de dependência nos services via factory na inicialização
- Adapters inicializados uma única vez na inicialização, com seus respectivos ports injetados
- Interfaces separadas das implementações, cada uma em seu próprio arquivo
- Separação clara entre arquivos de domínio (domain) e interfaces
- Handlers devem chamar services injetados, que por sua vez chamam repositórios injetados
- Implementação efetiva (sem uso de mocks)
- Manutenção da consistência no padrão de desenvolvimento entre funções
- Tratamento de erros sempre utilizando utils/http_errors
- Remoção completa de código legado após a refatoração, dado que estamos em fase ativa de desenvolvimento
- Eventuais alterações no DB são feitas por MySQL Workbench, não crie/altere scripts para migração de dados/tabelas.
- Erros devem ser logados no momento do erro etransformados em utils/http_errors e retornados para a chamador
- chamadores intermediários apenas repassam o erro sem logging ou recriação do erro
- Todo erro deve ser verificado.

📌 Instruções finais
- Não implemente nada até que eu autorize.
- Analise cuidadosamente a solicitação e o código atual e busque a causa raiz e pssíveis correções
