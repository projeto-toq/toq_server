
Eu preciso que você atue como um engenheiro de software Go sênior, especializado em arquitetura hexagonal e boas práticas de código. Siga as instruções abaixo de forma **ESTRITA** para analisar um problema e propor uma solução em português.

---

### 🛠️ Análise e Solução

**Problema:** Após várias refatorações estou fazendo uma verificação de qualidade. Assim, analise o fluxo de solicitação  de perfil user_handler/update_profile, e verifique se:
- a lógica está correta;
- existem otimizações possíveis;
- a documentação está adequada.
**Atenção** os campos do perfil de usuário, senha, telefone e email tem processo próprio de atualização e não podem ser atualizados por update_profile

---

### REGRAS OBRIGATÓRIAS DE ANÁLISE E PLANEJAMENTO

1.  **Arquitetura e Fluxo de Código**
    * **Arquitetura:** A solução proposta deve seguir estritamente a Arquitetura Hexagonal.
    * **Fluxo de Chamadas:** Mantenha a hierarquia de dependências: `Handlers` → `Services` → `Repositories`.
    * **Injeção de Dependência:** O plano deve contemplar o padrão de factories para injeção de dependências.
    * **Localização de Repositórios:** A solução deve prever que os repositórios residam em `/internal/adapter/right/mysql/`.
    * **Transações SQL:** Todas as transações de banco de dados devem utilizar `global_services/transactions`.

2.  **Tratamento de Erros**
    * A solução deve prever o tratamento de erros conforme o padrão do projeto, utilizando os pacotes `http/http_errors` ou `utils/http_errors`.
    * O plano deve garantir a correta propagação e log de erros.

3.  **Boas Práticas Gerais**
    * **Estilo de Código:** A proposta deve alinhar-se com o Go Best Practices e o Google Go Style Guide.
    * **Separação:** O plano deve manter a clara separação entre arquivos de `domínio`, `interfaces` e suas implementações.
    * **Processo:** Não inclua no plano a geração de scripts de migração de banco de dados ou qualquer tipo de solução temporária.

---

### REGRAS DE DOCUMENTAÇÃO E COMENTÁRIOS
* A documentação da solução deve ser clara e concisa.
* O plano deve prever a documentação das funções em **inglês** e comentários internos **em português**, quando necessário.
* Se aplicável, a solução deve incluir documentação para a API no padrão **Swagger**.

---

### INSTRUÇÕES FINAIS PARA O PLANO
* **Ação:** Não implemente nenhum código. Apenas analise e gere o plano.
* **Análise:** Analise cuidadosamente o problema e os requisitos. Se necessário, solicite informações adicionais. Analise sempre o código e os arquivos de configuração existentes.
* **Plano:** Apresente um plano detalhado para a implementação. O plano deve incluir:
    * Descrição da arquitetura proposta e seu alinhamento com a arquitetura hexagonal.
    * Interfaces a serem criadas (com métodos e assinaturas).
    * Estrutura de diretórios e arquivos sugerida.
    * Ordem das etapas de refatoração para garantir uma transição suave.
* **Qualidade do Plano:** O plano deve ser completo, sem mocks ou soluções temporárias. Se for muito grande, divida-o em etapas que possam ser implementadas separadamente.
* **Acompanhamento:** Sempre informe as etapas já planejadas e as próximas etapas a serem analisadas/planejadas para o acompanhamento do processo.