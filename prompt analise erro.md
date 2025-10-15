### Resumo e Refatoração: Engenheiro de Software Go Sênior

Este documento descreve as instruções para atuar como um engenheiro de software Go sênior, focando na análise de um problema e na proposição de uma solução detalhada, seguindo a arquitetura hexagonal e boas práticas de código, garantindo o alinhamento com os padrões de arquitetura, tratamento de erros, observabilidade e documentação. Toda a interação deve ser feita em português.

---

**Problemas:**
Voce estava implementando o plano criado por voce mesmo e aprovado por mim, para resolver os problemas abaixo, quando um erro apagou o arquivo internal/adapter/left/http/handlers/listing_handlers/get_all_listings.go e em seguida voce se perdeu na implemtação do plano. Parte foi implementado e parte não.

é necessário que voce verifique o que já foi feito, me parece que estava no final da implementação, e termine o que falta.

Os problemas são:
- O endpoint GET /admin/users que efetuar busca de usuários com filtros, somente aceita que os parâmetros  sejam passados na integra e não aceita buscas parciais, como like ou *abc*.
- O mesmo acontece com get /admin/roles que busca por roles.
- o mesmo acontece com GET /listings.
- é necessário que estes endpoints tenham forma de buscas parciais para strings e intervalos para datas e números.
- o endpoint GET /admin/user/pending está sem paginação, retornando todos os usuários pendentes de uma vez só.

O plano que voce havia traçado para resolver estes problemas é o seguinte:
Plano Detalhado

Mapear filtros e utilidades

Definir quais campos string aceitarão busca parcial (incluindo tradução de * → %).
Criar helper em utils para sanitizar padrões (TrimSpace, substituição segura de curingas).
Documentar limites padrão de paginação para reuso em handlers/serviços.
Atualizar DTOs e Handlers HTTP

dto.AdminListUsersRequest: adicionar campos opcionais idFrom/idTo, bornAtFrom/bornAtTo, lastActivityFrom/lastActivityTo (ISO 8601), garantir conversão de * em todos os campos string.
get_admin_users.go: parsear datas com time.Parse, validar faixas, popular um novo userservices.ListUsersInput estendido e atualizar Swagger (@Param) conforme a regra (documentação via comentários).
dto.AdminListRolesRequest: incluir description, slug parcial, idFrom/idTo; ajustar handler get_admin_roles.go para tratar novas entradas.
dto.GetAllListingsRequest: acrescentar filtros de texto (cidade, bairro, código, status, título) e de intervalo (createdFrom/To se disponível via auditoria, minLandSize/maxLandSize, minRent/maxRent etc.).
get_all_listings.go: substituir resposta 501 por chamada ao serviço, tratar parsing numérico (strconv.ParseInt/ParseFloat) e atualizar comentários Swagger.
Estender contratos dos serviços

Expandir userservices.ListUsersInput/ListUsersOutput para carregar as novas faixas e padrões; ajustar ListUsers para normalizar defaults e repassar tudo ao repositório.
permissionservice.ListRolesInput: incluir novos campos e validar intervalos antes da chamada ao repositório.
Criar listingservices.ListListingsInput/ListListingsOutput e adicionar GetAllListings à interface (listing_service.go), garantindo tracing (utils.GenerateTracer) e transação read-only do globalService.
Ajustar portas e adapters MySQL

user_repository.ListUsersFilter: adicionar campos de faixa/datas; refatorar list_users.go para montar cláusulas LIKE em todos os strings e BETWEEN/>= <= para intervalos (usando helper de padrões).
permission_repository.RoleListFilter: idem para slug/descriptions e range de IDs; atualizar list_roles.go.
Introduzir listing_repository.ListListingsFilter + ListListings e implementar em mysql/listing: SQL com projeção enxuta (campos usados no DTO), COUNT(DISTINCT ...) para paginação e joins necessários; reutilizar conversores existentes para montar domínios.
Injeção e conversões

Atualizar implementações de serviço (internal/core/service/*) para consumir novos filtros, garantindo utils.SetSpanError em falhas e logs slog apenas quando necessário.
Ajustar pontos de conversão (ex.: criar toListingResponse reutilizável em listing_handlers para mapear domínio → DTO, incluindo timestamps caso consultados na auditoria).
Documentação e Swagger

Revisar comentários nos handlers/DTOs com descrições em inglês dos novos parâmetros.
Após implementação, rodar make swagger (sem editar arquivos gerados manualmente) para publicar mudanças.
Validação Manual

Exercitar /admin/users, /admin/roles, /listings com combinações de LIKE, * e intervalos (datas/números) verificando paginação e contagem.
Monitorar logs/traces para confirmar que spans adicionais não surgem em handlers e que erros propagam pelos utilitários padrão.
Acompanhamento

Etapas analisadas: handlers admin/listing, serviços user/permission/listing, adapters MySQL correspondentes, DTOs e schema SQL.
Próximos focos na execução: atualizar contratos/ports conforme o plano, implementar a nova query em mysql/listing, e revisar factories/outros pontos que compilarão com as novas assinaturas antes de gerar Swagger.

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
- Não execute git status, git diff nem go test.

---

### **Regras de Documentação e Comentários**

- A documentação da solução deve ser clara e concisa.
- A documentação das funções deve ser em **inglês**.
- Os comentários internos devem ser em **português**.
- A API deve ser documentada com **Swagger**, usando anotações diretamente no código, em inglês e não alterando swagger.yaml/json manualmente.