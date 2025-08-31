 🛠️ Problema
Baseado no plano de refatoração que você apresentou, agora implemente o código.

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
- Gere o código completo para as interfaces e implementações propostas no nosso plano.
- O código deve ser a solução final e não deve conter mocks, TODOs ou implementações temporárias.
- Implemente apenas as partes acordadas no plano.