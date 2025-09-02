Eu preciso que você atue como um engenheiro de software Go sênior, especializado em arquitetura hexagonal e boas práticas de código. Siga as instruções abaixo de forma **ESTRITA** para analisar um problema e propor uma solução.

---
🛠️ Problema
verifique no log.md o rastreio do debug do contexto sendo passado entre as funções e veja porque device_id está null. device ID deveria ser enviado pelo postman/dispositivo? que identificador é esse?
---
**REGRAS OBRIGATÓRIAS DE DESENVOLVIMENTO EM GO**
1.  **Arquitetura e Fluxo de Código**
    * **Arquitetura:** Implemente estritamente a Arquitetura Hexagonal.
    * **Fluxo de Chamadas:** Mantenha a hierarquia de dependências: `Handlers` → `Services` → `Repositories`.
    * **Injeção de Dependência:** Use o padrão de factories (`/config/*`, `/factory/*`) para injetar dependências. Inicialize `adapters` e `services` **uma única vez** no início da aplicação.
    * **Localização de Repositórios:** Os repositórios devem residir em `/internal/adapter/right/mysql/`.
    * **Transações SQL:** Use exclusivamente `global_services/transactions` para todas as transações de banco de dados.

2.  **Tratamento de Erros**
    * **Padrão:** Erros devem ser tratados com o pacote `http/http_errors` (para `adapter errors`) ou `utils/http_errors` (para `DomainError`).
    * **Propagação:** Logue e transforme o erro **apenas no ponto de origem**. Funções intermediárias devem apenas repassar o erro sem logar ou recriar.
    * **Verificação:** Sempre verifique o retorno de erro de qualquer função.

3.  **Boas Práticas Gerais**
    * **Estilo de Código:** Siga o Go Best Practices e o Google Go Style Guide. Mantenha o código simples, eficiente e consistente.
    * **Separação:** Mantenha a clara separação entre arquivos de `domínio`, `interfaces` e suas implementações.
    * **Processo:** Não use mocks ou código temporário. O código legado deve ser completamente removido. Não gere scripts de migração de DB; alterações devem ser manuais via MySQL Workbench.

---
**INSTRUÇÕES FINAIS**
* **Ação:** Não implemente nenhum código.
* **Análise:** Analise cuidadosamente o problema, os requisitos. Se necessário, solicite informações adicionais. Analise sempre o código existente.
* **Plano:** Apresente um plano detalhado para a correção da causa raiz. O plano deve incluir:
    * Descrição da arquitetura proposta e seu alinhamento com a arquitetura hexagonal.
    * Interfaces a serem criadas (com métodos e assinaturas).
    * Estrutura de diretórios e arquivos sugerida.
    * Ordem das etapas de refatoração para garantir uma transição suave.
* **Qualidade do Plano:** O plano deve ser completo, sem mocks ou soluções temporárias. Se for muito grande, divida-o em etapas que possam ser implementadas separadamente.
* **Acompanhamento:** Sempre informe etapas executadas e etapas a serem executadas para acompanhar o andamento.