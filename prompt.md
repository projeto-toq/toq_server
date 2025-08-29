🛠️ Problema
existe um dualidade sobre slugs no projeto entre root e admin. em alguns lugares é usado root como role em outros admin.
- faça um scan no projeto e padronize como admin.
- revise os arquivos csv para que tenhamos os 6 base roles do base_permission_roles mas sem usuários criados, a exceção do admin, que é criado na criação do banco

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

📌 Instruções finais
- Não implemente nada até que eu autorize.
- Analise cuidadosamente a solicitação e o código atual, e apresente um plano detalhado de implementação antes de qualquer alteração.
