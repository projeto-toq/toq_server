🛠️ Problema
estamos com estas  funcções necessiando implementação real:
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

✅ Requisitos obrigatórios para qualquer revisão, refatoração ou correção:
- Adoção das melhores práticas de desenvolvimento em Go (Go Best Practices, Google Go Style Guide).
- Implementação seguindo arquitetura hexagonal.
- Injeção de dependência nos services via factory na inicialização.
- Adapters inicializados uma única vez na inicialização, com seus respectivos ports injetados.
- Interfaces separadas das implementações, cada uma em seu próprio arquivo.
- Separação clara entre arquivos de domínio (domain) e interfaces.
- Handlers devem chamar services injetados, que por sua vez chamam repositórios injetados.
- Implementação efetiva (sem uso de mocks ou código temporário).
- Manutenção da consistência no padrão de desenvolvimento entre funções.
- Tratamento de erros sempre utilizando utils/http_errors.
- Remoção completa de código legado após a refatoração.
- Eventuais alterações no DB são feitas por MySQL Workbench, não crie/altere scripts para migração de dados/tabelas.
- Erros devem ser logados no momento do erro e transformados em utils/http_errors e retornados para o chamador.
- Chamadores intermediários apenas repassam o erro sem logging ou recriação do erro.
- Todo erro deve ser verificado.

📌 Instruções finais
- Não implemente nenhum código.
- Analise cuidadosamente o problema e os requisitose solicite informações adicionais se necessário.
- Apresente um plano detalhado para a refatoração. O plano deve incluir:
  - Uma descrição da arquitetura proposta e como ela se alinha com a arquitetura hexagonal.
  - As interfaces que precisarão ser criadas (com seus métodos e assinaturas).
  - A estrutura de diretórios e arquivos sugerida.
  - A ordem das etapas de refatoração para garantir uma transição suave e sem quebras.
- Certifique-se de que o plano esteja completo e não inclua mocks ou soluções temporárias.
- Apenas apresente o plano, sem gerar o código.
