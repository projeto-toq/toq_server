🛠️ Problema
após várias refatorações, o processo de inicialização do sistema como um todo sendo realizada no /config e cmd/toq_server está confuso, complexo e instável. várias funções duplicadas e ações endo pulaas por não chamar a rotina certa. 
Exemplo foi sua recente descoberta:
1. Função initializeServices() Inexistente
Localização: /codigos/go_code/toq_server/internal/core/config/inject_dependencies.go:77
Problema: A função é chamada mas não foi implementada
Impacto: Os serviços não são inicializados, incluindo permissionService
2. PermissionService Nil
Localização: role_mapper.go - verificação us.permissionService == nil
Causa: permissionService não é inicializado devido ao problema #1
Sintomas: Erro durante criação de owner quando tenta usar permission service
3. Logging Inadequado de Erros
Localização: Toda implementação de permission service
Problema: Erros são detectados mas não logados adequadamente
Impacto: Dificuldade de debug e monitoramento


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
- Eventuais alterações no DB são feitas por MySQL Workbench, portatno apenas indique as modificação e não crie/altere scripts para migração de dados/tabelas.
- Erros devem ser logados no momento do erro e retornados para a chamador para tratamento,
- Todo erro deve ser verificado.

📌 Instruções finais
- Não implemente nada até que eu autorize.
- Analise cuidadosamente a solicitação e o código atual, e apresente um plano detalhado de implementação para rescrever totalmente a inicillização do sistema para que fique eficiente, limpa, de fácil manutenção, clara.
