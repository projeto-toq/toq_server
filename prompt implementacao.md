
Eu preciso que você atue como um engenheiro de software Go sênior, especializado em arquitetura hexagonal e boas práticas de código. Siga as instruções abaixo de forma **ESTRITA** para implementar o plano proposto.

## INSTRUÇÕES PARA GERAÇÃO E IMPLEMENTAÇÃO DE CÓDIGO

* **Ação:** Gere e implemente **exclusivamente** o código Go para as interfaces e funções acordadas no plano de refatoração.
* **Qualidade:** O código deve ser a solução **final e completa**. Não inclua mocks, `TODOs` ou qualquer tipo de implementação temporária.
* **Escopo:** Implemente **somente** as partes que foram definidas no plano. Não adicione funcionalidades extras ou códigos que não sejam estritamente necessários para a solução.
* **Simplicidade:** Mantenha o código simples e eficiente, conforme as regras de boas práticas do projeto.
* **Acompanhamento:** Sempre informe etapas executadas e etapas a serem executadas para acompanhar o andamento.
**Contexto de Desenvolvimento:** Estamos em um ambiente de desenvolvimento. Não é necessário gerar scripts de migração de banco de dados, manter retrocompatibilidade (`backward compatibility`) ou lidar com migração de dados. As alterações no banco de dados devem ser feitas manualmente.

---

## REGRAS OBRIGATÓRIAS DE DESENVOLVIMENTO EM GO

### 1. Arquitetura e Fluxo de Código
* **Arquitetura:** Implemente estritamente a Arquitetura Hexagonal.
* **Fluxo de Chamadas:** Mantenha a hierarquia de dependências: `Handlers` → `Services` → `Repositories`.
* **Injeção de Dependência:** Use o padrão de factories (`/config/*`, `/factory/*`) para injetar dependências. Inicialize `adapters` e `services` **uma única vez** no início da aplicação.
* **Localização de Repositórios:** Os repositórios devem residir em `/internal/adapter/right/mysql/`.
* **Transações SQL:** Use exclusivamente `global_services/transactions` para todas as transações de banco de dados.

### 2. Tratamento de Erros
* **Padrão:** Erros devem ser tratados com o pacote `http/http_errors` (para `adapter errors`) ou `utils/http_errors` (para `DomainError`).
* **Propagação:** Logue e transforme o erro **apenas no ponto de origem**. Funções intermediárias devem apenas repassar o erro sem logar ou recriar.
* **Verificação:** Sempre verifique o retorno de erro de qualquer função.

### 3. Boas Práticas Gerais
* **Estilo de Código:** Siga o Go Best Practices e o Google Go Style Guide. Mantenha o código simples, eficiente e consistente.
* **Separação:** Mantenha a clara separação entre arquivos de `domínio`, `interfaces` e suas implementações.
* **Processo:** Não use mocks ou código temporário. O código legado deve ser completamente removido.

---

## REGRAS DE DOCUMENTAÇÃO E COMENTÁRIOS

* **Documentação de Funções:** Documente todas as funções e métodos exportados usando o padrão do Go. A documentação deve ser clara, concisa e **escrita em inglês**.
* **Comentários Internos:** Use comentários internos **em português** para explicar trechos de código complexos ou lógicas específicas que não são imediatamente óbvias.
* **Documentação da API:** Onde aplicável, inclua documentação para a API usando o **padrão Swagger** nos `handlers`.

---

