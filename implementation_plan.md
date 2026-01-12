# Plano de Implementação — Marketplace de Propostas

## 1. Contexto e Objetivo
Este plano descreve com nível máximo de detalhes as alterações necessárias para atender aos novos requisitos de negócio:
- **Lista de Propostas (owner)**: exibir nome, foto assinada, tempo de cadastro na TOQ e quantidade de propostas aceitas do corretor.
- **Detalhes da Proposta (owner/realtor)**: expor datas de criação, recebimento e resposta da proposta.

O plano segue estritamente o guia `docs/toq_server_go_guide.md`, preservando a Regra de Espelhamento e os padrões de documentação (Godoc + Swagger). Nenhum código deve ser implementado antes da aprovação formal registrada em `prompt_approvall.md`.

## 2. Diagnóstico
### 2.1 Arquitetura atual
- **Handlers/DTOs**: [internal/adapter/left/http/handlers/proposal_handlers](internal/adapter/left/http/handlers/proposal_handlers) já oferecem `/proposals/realtor`, `/proposals/owner` e `/proposals/detail`, porém os DTOs não incluem timeline e os dados do corretor são limitados.
- **Serviço**: [internal/core/service/proposal_service](internal/core/service/proposal_service) centraliza a orquestração, mas não depende de `userService`, o que impede gerar URLs assinadas.
- **Repositório MySQL**: [internal/adapter/right/mysql/proposal](internal/adapter/right/mysql/proposal) calcula `usage_months` e `total_proposals`, porém não retorna contagem de propostas aceitas nem fornece fotos assinadas (depende do serviço).
- **Modelo**: [internal/core/model/proposal_model](internal/core/model/proposal_model) não expõe campos para `AcceptedProposals`/`PhotoURL` em `RealtorSummary`, tampouco getters que facilitem timeline.

### 2.2 Lacunas identificadas
1. **Dados do Corretor**:
   - Falta do campo `AcceptedProposals` na query agregada (`list_realtor_summaries.go`).
   - Ausência de URL assinada (`userService.GetPhotoDownloadURL`).
2. **Timeline**:
   - DTO `ProposalResponse` não possui `createdAt`, `receivedAt` (primeira ação do owner) e `respondedAt` (timestamp final associado ao status).
3. **Injeção de dependência**:
   - `proposalService` precisa receber `userservices.UserServiceInterface` para reutilizar a lógica de assinatura já disponível.
4. **Conversões**:
   - `converters.ProposalDomainToResponse` precisa mapear novos campos (timeline + enriquecimento do corretor).

### 2.3 Impactos e riscos
- **Mudança de contrato HTTP**: novos campos serão adicionados às respostas JSON; é necessário atualizar as anotações Swagger nos handlers existentes para garantir documentação consistente.
- **Desempenho**: geração de URLs assinadas para cada corretor exige cuidado (usar cache em memória do loop e reutilizar contextos com `utils.SetUserInContext`).
- **Telemetria**: chamadas para `userService.GetPhotoDownloadURL` já seguem o padrão de tracing via service.
- **DBA**: nenhuma mudança de schema; apenas consultas enriquecidas.

## 3. Escopo Técnico Detalhado
### 3.1 Domínio (`internal/core/model/proposal_model/realtor_summary.go`)
- Adicionar métodos `AcceptedProposals() int64`, `SetAcceptedProposals(int64)`, `PhotoURL() string`, `SetPhotoURL(string)` na interface e implementação.
- Atualizar comentários para explicar os novos campos (contagem de aceites e URL assinada).

### 3.2 Entidade SQL (`internal/adapter/right/mysql/proposal/entities/realtor_summary_entity.go`)
- Incluir `AcceptedProposals sql.NullInt64` com comentário explicando que a origem é `SUM(status='accepted')`.

### 3.3 Repositório (`internal/adapter/right/mysql/proposal/list_realtor_summaries.go`)
- Ajustar SELECT para calcular `accepted_proposals`:
  ```sql
  SUM(CASE WHEN status = 'accepted' THEN 1 ELSE 0 END) AS accepted_proposals
  ```
- Atualizar `rows.Scan` para preencher o novo campo.
- Garantir que a subquery continue evitando `SELECT *` e mantenha índices funcionais (`realtor_id`, `status`).

### 3.4 Conversor (`internal/adapter/right/mysql/proposal/converters/realtor_summary_entity_to_domain.go`)
- Mapear `entity.AcceptedProposals` para `summary.SetAcceptedProposals()`.

### 3.5 Serviço (`internal/core/service/proposal_service`)
1. **Estrutura e construtor**
   - Adicionar `userService userservices.UserServiceInterface` ao struct.
   - Atualizar `New` para receber `userService` e armazenar; manter backward compatibility para testes (permitir `nil`).
2. **Helpers** (`helpers.go`)
   - Criar `enrichRealtorSummaries(ctx context.Context, cache map[int64]proposalmodel.RealtorSummary) error`:
     - Ignorar se `userService` for `nil`.
     - Iterar sobre o mapa; se `PhotoURL` vazio, gerar contexto impersonado com `utils.SetUserInContext(ctx, usermodel.UserInfos{ID: realtorID})` e chamar `userService.GetPhotoDownloadURL(..., "small")`.
     - Em caso de erro, logar em nível `Debug` e seguir (não bloquear a lista).
   - Criar `generateRealtorPhotoURL(ctx context.Context, userID int64) (string, error)` encapsulando a lógica acima para facilitar testes.
3. **listProposals**
   - Após `loadRealtorSummaryMap`, chamar `enrichRealtorSummaries` antes de montar `ListItem`.
   - Garantir que `ContextWithLogger` é reutilizado para logs dentro do helper.

### 3.6 DTOs e Conversores HTTP
1. **DTO** (`internal/adapter/left/http/dto/proposal_dto.go`)
   - Expandir `ProposalResponse` com campos opcionais `CreatedAt`, `ReceivedAt`, `RespondedAt` (`*time.Time`).
   - Atualizar `ProposalRealtorResponse` para incluir `AccountAgeMonths`, `AcceptedProposals` e `PhotoURL` (todos com comentários + tags `example`).
2. **Converter** (`internal/adapter/left/http/converters/proposal_converter.go`)
   - `ProposalDomainToResponse`: popular novos timestamps usando `proposal.CreatedAt()`, `proposal.FirstOwnerActionAt()` e os `NullTime` de status.
   - `proposalRealtorToResponse`: preencher novos campos a partir de `RealtorSummary`.
   - Adicionar helpers internos `firstNonNilTimePtr` e `nullableTimePtr` se necessário, mantendo a regra “uma função pública por arquivo”.

### 3.7 Serviço → DTO: timeline
- **Data de criação**: `proposal.CreatedAt()`.
- **Data de recebimento**: `proposal.FirstOwnerActionAt()` (caso `Valid=false`, campo omitido).
- **Data de resposta**: primeiro timestamp válido em `accepted_at`, `rejected_at` ou `cancelled_at`.

### 3.8 Injeção de Dependências (`internal/core/config/inject_dependencies.go`)
- Atualizar chamada para `proposalservice.New` adicionando `c.userService`.
- Garantir que `InitProposalService` valida `userService != nil` (log `Warn` se nil, mas permitir continuar – apenas ficamos sem foto assinada em ambientes de teste).

### 3.9 Swagger / Documentação
- Handlers existentes devem receber comentários extras nas respostas para documentar os novos campos, por exemplo:
  ```go
  // @Success 200 {object} dto.ListProposalsResponse "Items now include realtor.photoUrl and timeline timestamps"
  ```
- Após implementação, executar `make swagger` para regenerar `docs/swagger.*` (feito somente após aprovação).

## 4. Estrutura Final Esperada
```
internal/core/model/proposal_model/
  realtor_summary.go          # novos getters/setters + comentários
internal/adapter/right/mysql/proposal/entities/
  realtor_summary_entity.go   # campo AcceptedProposals
internal/adapter/right/mysql/proposal/
  list_realtor_summaries.go   # query com accepted count
  converters/realtor_summary_entity_to_domain.go
internal/core/service/proposal_service/
  proposal_service.go         # userService injection
  helpers.go                  # enrichment helpers
internal/adapter/left/http/dto/proposal_dto.go
internal/adapter/left/http/converters/proposal_converter.go
internal/core/config/inject_dependencies.go
```

## 5. Sequenciamento das Atividades
1. **Atualizar Domínio e DTOs**
   - Modificar `RealtorSummary` e `ProposalResponse`.
   - RODAR `go test ./internal/core/model/...` para garantir compilação.
2. **Ajustar Repositório**
   - Editar entidade, query e converter.
   - Validar com testes unitários específicos do adapter, se existentes.
3. **Modificar Serviço**
   - Introduzir `userService` e helpers.
   - Assegurar que `listProposals` chama o enrichment.
   - Atualizar mocks (se houver) para refletir a nova dependência.
4. **Atualizar Conversores HTTP**
   - Popular timeline e novos campos do corretor.
   - Revisar os handlers para garantir que a documentação mencione os campos.
5. **Injeção de Dependência**
   - Ajustar `InitProposalService` e `proposalservice.New`.
   - Executar `go build ./cmd/...` para garantir wiring correto.
6. **Swagger** (após aprovação + codificação)
   - `make swagger` e commit do diff gerado.

## 6. Validações e Testes Recomendados
- **Unitários**: serviços de proposta (`list_proposals`, `enrichRealtorSummaries`), converter DTO.
- **Integração**: endpoint `/api/v2/proposals/owner` e `/api/v2/proposals/detail` via testes de rota (se existir suíte).
- **Manuais**:
  1. Criar uma proposta via `/proposals`.
  2. Listar como owner e verificar:
     - `realtor.photoUrl` preenchido quando o corretor possui foto.
     - `realtor.acceptedProposals` refletindo histórico.
     - `proposal.createdAt`, `receivedAt`, `respondedAt` consistentes.
  3. Aceitar a proposta e revalidar timeline.

## 7. Comunicação e Aprovação
- Submeter este plano para aprovação em `prompt_approvall.md` antes de qualquer alteração de código.
- Após aprovação, seguir estritamente a ordem sequencial e documentar cada etapa no PR com referência a este arquivo.

---
Este documento foi gerado por GitHub Copilot (GPT-5.1-Codex) em 12/01/2026 para orientar times presentes e futuros na implementação segura dos novos requisitos do fluxo de propostas.
