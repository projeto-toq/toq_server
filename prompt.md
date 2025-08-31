🛠️ Problema
Alertas do projeto
[{
	"resource": "/codigos/go_code/toq_server/go.mod",
	"owner": "_generated_diagnostic_collection_name_#2",
	"severity": 4,
	"message": "go.opentelemetry.io/otel/trace should be direct",
	"source": "go mod tidy",
	"startLineNumber": 117,
	"startColumn": 41,
	"endLineNumber": 117,
	"endColumn": 52,
	"origin": "extHost2"
}]
[{
	"resource": "/codigos/go_code/toq_server/internal/adapter/left/http/middlewares/structured_logging_middleware.go",
	"owner": "go-golangci-lint",
	"severity": 4,
	"message": "ineffectual assignment to logLevel (ineffassign)",
	"source": "go-golangci-lint",
	"startLineNumber": 82,
	"startColumn": 3,
	"endLineNumber": 82,
	"endColumn": 4,
	"origin": "extHost2"
}]

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
- Analise cuidadosamente a solicitação e o código atual, descubra a causa raiz e proponha a solução