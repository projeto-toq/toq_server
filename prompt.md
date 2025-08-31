🛠️ Problema
apesar de claramente dizer que - Implementação efetiva (sem uso de mocks) - inssitemente prompts enviados geram códigos como:
// getUserPermissionsFromCache busca permissões do cache Redis
func (p *permissionServiceImpl) getUserPermissionsFromCache(_ context.Context, cacheKey string) ([]permissionmodel.PermissionInterface, error) {
	// TODO: Implementar cache otimizado para o novo sistema
	// Por enquanto, sempre retorna cache miss para forçar busca no banco
	slog.Debug("Cache temporarily disabled for user permissions", "cache_key", cacheKey)
	return nil, fmt.Errorf("cache miss - using database")
}

// setUserPermissionsInCache armazena permissões no cache Redis
func (p *permissionServiceImpl) setUserPermissionsInCache(_ context.Context, cacheKey string, permissions []permissionmodel.PermissionInterface) error {
	// TODO: Implementar cache otimizado para o novo sistema
	// Por enquanto, não faz cache para simplificar a migração
	slog.Debug("Cache temporarily disabled for storing user permissions", "cache_key", cacheKey, "count", len(permissions))
	return nil
}

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
- Analise cuidadosamente a solicitação e o código atual, e apresente a melhor forma para fazer a refatoração buscando simplicade e melhores práticas