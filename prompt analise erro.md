### Resumo e Refatoração: Engenheiro de Software Go Sênior

Este documento descreve as instruções para atuar como um engenheiro de software Go sênior, focando na análise de um problema e na proposição de uma solução detalhada, seguindo a arquitetura hexagonal e boas práticas de código, garantindo o alinhamento com os padrões de arquitetura, tratamento de erros, observabilidade e documentação. Toda a interação deve ser feita em português.

---

**Problemas:**
O endpoint /user/photo/download-url está retornando para o usuário o url
{
    "signedUrl": "https://toq-app-media.s3.us-east-1.amazonaws.com/2/photo/original.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Checksum-Mode=ENABLED&X-Amz-Credential=ASIAQ3EGR6UWYYLJDK5G%2F20251013%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20251013T145652Z&X-Amz-Expires=3600&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEJ%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEaCXVzLWVhc3QtMSJIMEYCIQCYexn2IypiM00I85r3R0med9j03GduynSGsdoixwb9CwIhAJ9gcHQBwrZw259LzWOXAhXHdz8t9J3mEmpILdjHnKQzKr0FCEgQABoMMDU4MjY0MjUzNzQxIgwNNJ%2FaOdbdglsVo2MqmgVGdOOTbwnaBrQUq7wVUVNRLxNWdV7YuncMa4AVKXNb6%2FQ2AMKUFePn4sD2qDGzqGLYlSV7IktyYttyISrmT4FRv3GTsWKP4ikHQQjS5Sb7upmGYlaXN7JUmby%2FHc1K8vCo6%2Fu3tVG2mDazgpAysdc4QAVv%2BPSPLK9XsFp8u0d4Xw4ixgLT0zrr1dAQ%2BZGoXbqKqKvCkyE6Ck25sowGCIvjp7uLqh3mjt%2FA0Xw2%2FDiFUBYM6CtlaovmZhNAeeCcAoZWBfJNbVcdLZNSLGi%2FR%2BnlWjl9CFbeeoOb%2Fo80pzfaPCZAQTTOivUh98WCDKK4VJiAqyZCOql2DXKVVnMFhGplP6gRDiBBuEW1nKUJScDxvMfrLXPJmBOzeDckVi0wrxcKZ%2Bm9ibgMeNtSQPuo3n6ld2oB0bv5qdWVzSD0E%2Fa2cexxQb5FHLMJ%2Bxa50nE102E6QTpoj1HPUMXVowgYTXtiIDZELywmKgh1PqLhlkXoQFHxtojsniSvDVmLmH%2BbH4h9b1kFkDSWUTPisQ1MZxAX%2B7ShSYlo7hQ%2F2O%2BLXm%2FjL4SNHEl%2FnLnQKc1mRqRGnvGCZNj7HuKYrIHvDwf58eiKwFN%2FB5HEr4CtpSuNNuWiinAeMkKGOoHms14LXdfbABErqmpAC8i%2Bx2bDfLhCaoItaS3wP1h0%2F4s7CyFxu2Zg339THyRYqH23KAMlh6ZpnXipQ%2F60kCaBStHRANXXwdbCUJk2ARL%2Bu%2BkCzMXA6Y0cqrVL0bIiO%2BT6x7lXNkL%2FYSDFmTbAxRzMLyO6Xjo7oznbloHnKJDPNDDqwPpc8dKcyJmskWpcB9PwjqdD6H0PRUjEhwlQfFeftZ08i8snPgGBy0Yim%2FVQeCSqOY0hgdYRqJOTlPgaqAmuh3ow6Z60xwY6sAEiiqjbnyYxepniOQBt7O2yq9iurOQ7DxVLo5E8IzPgVDVnjNKieb1aQ1tTkDwwGQWY3RCu1hEntGzBZFsP1SECb5fC%2BLKNejdj1S9%2BVF9UuPhnEJtWhVhq2I%2BevOIvRWI6yqOh1%2Bdc5tcsdWssYuaFNniuogexxVPOPqIj2szYsTCpP%2BSmsKIwM7dyizltFZAPto5Yw%2Fioc7qWASy%2BNA6ngK9LcD0RPQ1fi0CQOdu5Jw%3D%3D&X-Amz-SignedHeaders=host&x-id=GetObject&X-Amz-Signature=37f7b367ee7c516a8931698ac2c7f0b8a33b683f165cde1908b830063c6311ae"
}

entretanto ao colocar esta url no navegador, receno esta mensagem:
<Error>
<Code>InvalidAccessKeyId</Code>
<Message>The AWS Access Key Id you provided does not exist in our records.</Message>
<AWSAccessKeyId>ASIAQ3EGR6UWYYLJDK5G</AWSAccessKeyId>
<RequestId>C8F369ZYEPK84GT3</RequestId>
<HostId>cV0GMYdDFvmgB6BHYNN7bDU1bArka9Qr5AaaFZq12zMqUPG/mZvvlcq7up6ZGtIupk2a7c2rskI=</HostId>
</Error>

**Solicitação:** Analise o problema, **leia o código** envolvido, **ache a causa raiz** e proponha um plano detalhado para a implementação da solução.

### **Instruções para a Proposição do Plano**

- **Ação:** Apenas a análise e a geração do plano são solicitadas. **Nenhum código deve ser implementado**.
- **Análise:** O problema e os requisitos devem ser analisados cuidadosamente. O código e arquivos de configuração existentes devem ser revisados para um plano preciso. Não faça suposições e confirme todos os detalhes necessários.
- **Plano:** Um plano detalhado deve ser apresentado, incluindo a descrição da arquitetura proposta, as interfaces, a estrutura de diretórios e a ordem de execução das etapas.
- **Qualidade do Plano 1:** O plano deve ser completo, sem o uso de _mocks_ ou soluções temporárias. Caso seja extenso, deve ser dividido em etapas implementáveis.
- **Acompanhamento:** As etapas já planejadas e as próximas a serem analisadas devem ser sempre informadas para acompanhamento.
- **Ambiente:** O plano deve considerar que estamos em ambiente de desvolvimento, portanto não deve haver back compatibility, migração de dados, preocupação com janela de manutenção ou _downtime_.
- **Testes:** O plano **NÃO**deve incluir a criação/alteração de testes unitários e de integração para garantir a qualidade do código.
- **Documentação:** A documentação Swagger/docs deve ser criada por comentários em DTO/Handler e execuçÃo de make swagger. Sem alterações manuais no swagger.yaml/json.
---

### **Regras Obrigatórias de Análise e Planejamento**

#### 1. Arquitetura e Fluxo de Código
- **Arquitetura:** A solução deve seguir estritamente a **Arquitetura Hexagonal**.
- **Fluxo de Chamadas:** As chamadas de função devem seguir a hierarquia `Handlers` → `Services` → `Repositories`.
- **Injeção de Dependência:** O padrão de _factories_ deve ser usado para a injeção de dependências.
- **Localização de Repositórios:** Os repositórios devem ser localizados em `/internal/adapter/right/mysql/` e deve fazer uso dos convertess para mapear entidades de banco de dados para entidades e vice versa.
- **Transações SQL:** Todas as transações de banco de dados devem utilizar `global_services/transactions`.


#### 2. Tratamento de Erros e Observabilidade

- **Tracing:**
  - Iniciar _tracing_ com `utils.GenerateTracer(ctx)` em métodos públicos de **Services**, **Repositories** e em **Workers/Go routines**.
  - Evitar _spans_ duplicados em **Handlers HTTP**, pois o `TelemetryMiddleware` já inicia o _tracing_.
  - Chamar a função de finalização (`defer spanEnd()`) e usar `utils.SetSpanError` para marcar erros.

- **Logging:**
  - Usar `slog` para _logs_ de domínio e segurança.
    - `slog.Info`: Eventos esperados do domínio.
    - `slog.Warn`: Condições anômalas ou falhas não fatais.
    - `slog.Error`: Falhas internas de infraestrutura.
  - Evitar _logs_ excessivos em **Repositórios (adapters)**.
  - **Handlers** não devem gerar _logs_ de acesso, pois o `StructuredLoggingMiddleware` já faz isso.

- **Tratamento de Erros:**
  - **Repositórios (Adapters):** Retornam erros "puros" (`error`).
  - **Serviços (Core):** Propagam erros de domínio usando `utils.WrapDomainErrorWithSource(derr)` e criam novos erros com `utils.NewHTTPErrorWithSource(...)`.
  - **Handlers (HTTP):** Usam `http_errors.SendHTTPErrorObj(c, err)` para converter erros em JSON.

#### 3. Boas Práticas Gerais
- **Estilo de Código:** A proposta deve seguir as **Go Best Practices** e o **Google Go Style Guide**.
- **Separação:** Manter a clara separação entre arquivos de **domínio**, **interfaces** e suas implementações.
- **Processo:** O plano não deve incluir a geração de _scripts_ de migração ou soluções temporárias.

---

### **Regras de Documentação e Comentários**

- A documentação da solução deve ser clara e concisa.
- A documentação das funções deve ser em **inglês**.
- Os comentários internos devem ser em **português**.
- A API deve ser documentada com **Swagger**, usando anotações diretamente no código e não alterando swagger.yaml/json manualmente.