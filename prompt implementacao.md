 🛠️ Problema
Baseado no plano de refatoração que você apresentou, agora implemente o código.

✅ Requisitos OBRIGATÓRIOS a serem respeitados
1. Padrões de Arquitetura e Código
Código dever simples e eficiente.
Arquitetura Hexagonal: A implementação deve seguir a arquitetura hexagonal.
Fluxo de Dependências: O fluxo de chamadas deve ser Handlers → Services → Repositórios, todos com dependências injetadas.
Boas Práticas: Adotar as melhores práticas de desenvolvimento em Go, incluindo o Go Best Practices e o Google Go Style Guide.
Separação de Responsabilidades: Manter a separação clara entre arquivos de domínio, interfaces e suas respectivas implementações.

2. Injeção de Dependência
Padrão de Injeção: A injeção de dependência deve ser feita através de factories.
Estrutura de Repositórios: Os repositórios devem estar em /internal/adapter/right/mysql/.
Inicialização Única: Os adapters e services devem ser inicializados uma única vez na inicialização da aplicação.

3. Tratamento e Propagação de Erros
Padrão de Erros: Todos os erros devem ser tratados usando o pacote utils/http_errors.
Propagação:
Erros devem ser logados e transformados em utils/http_errors no ponto onde ocorrem.
Chamadores intermediários devem apenas repassar o erro, sem logar ou recriar.
Verificação: Toda função que pode retornar um erro deve ter sua resposta verificada.

4. Processo de Desenvolvimento
Sem Código Temporário: Implementações devem ser efetivas, sem a utilização de mocks ou código temporário.
Remoção de Legado: O código legado deve ser completamente removido após a refatoração.
Consistência: Manter a consistência no padrão de desenvolvimento entre todas as funções e arquivos.
Banco de Dados: Alterações de DB devem ser feitas manualmente via MySQL Workbench. Não criar scripts de migração.
Compatibilidade: Não é necessária retrocompatibilidade com versões anteriores.


📌 Instruções finais
- Não implemente nenhum código.
- Analise cuidadosamente o problema e os requisitose solicite informações adicionais se necessário.
- Analise sempre o código existente e não assuma nada sem verificar antes.
- Apresente um plano detalhado para a refatoração. O plano deve incluir:
  - Uma descrição da arquitetura proposta e como ela se alinha com a arquitetura hexagonal.
  - As interfaces que precisarão ser criadas (com seus métodos e assinaturas).
  - A estrutura de diretórios e arquivos sugerida.
  - A ordem das etapas de refatoração para garantir uma transição suave e sem quebras.
- Certifique-se de que o plano esteja completo e não inclua mocks ou soluções temporárias.
- Apenas apresente o plano, sem gerar o código.
