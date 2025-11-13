# Plano Faseado — Versionamento de Listings

> Referência cruzada: requisitos definidos em `prompt implementation.md` e guia `docs/toq_server_go_guide.md`.

## Visão Geral

O objetivo é introduzir versionamento integral dos listings, mantendo observabilidade, arquitetura hexagonal e compatibilidade com as integrações internas. A implementação será dividida em fases incrementais, cada uma fechando o ciclo build/lint para permitir rollback controlado.

### Convenções Gerais
- **Build verde ao fim de cada fase** (`make lint`).
- **Documentação**: atualizar Swagger/Godoc/README na fase correspondente.
- **Tracing/Logging**: seguir seções 5, 7 e 9 do guia.
- **InstrumentedAdapter** obrigatório em todos os repositórios MySQL alterados.
- **Transações**: sempre via `globalService` (ou `StartReadOnlyTransaction`).

## Fases

### Fase 1 — Modelos e Interfaces
- Ajustar `listing_model` para separar identidade (`Listing`) de versões (`ListingVersionInterface`).
- Atualizar interfaces de entidades satélites (features, guarantees, financing blockers, exchange places) para referenciar `listing_version_id`.
- Manter métodos legados temporariamente (com TODO) ou adaptadores para não quebrar chamadas existentes até a próxima fase.
- Garantir que `NewListing()` e fábricas retornem estruturas consistentes (UUID ainda opcional nesta fase).

### Fase 2 — Ports e Adapters
- Refatorar `ListingRepoPortInterface` com novos métodos (`CreateListingVersion`, `PromoteListingVersion`, etc.).
- Implementar adapters MySQL seguindo a Regra de Espelhamento, criando novos arquivos (`create_listing_version.go`, `promote_listing_version.go`, etc.).
- Ajustar entidades MySQL para refletir `listing_versions` e novos FKs; preparar conversores.
- Inserir MIGRATION NOTES para DBAs (sem tocar em `scripts/db_creation.sql`).
  - **Notas de Migração:**
    - Criar tabela `listing_identities` com colunas `id` (PK), `listing_uuid` (UUID), `user_id`, `code`, `active_version_id` (FK opcional para `listing_versions`), campos de auditoria e flag `deleted`.
    - Renomear/duplicar a tabela atual de `listings` para `listing_versions`, adicionando coluna obrigatória `listing_identity_id` (FK para `listing_identities`) e flag `deleted`; preservar demais colunas.
    - Atualizar tabelas satélite (`features`, `exchange_places`, `financing_blockers`, `guarantees`) para referenciar `listing_version_id` em vez de `listing_id`, garantindo `ON DELETE CASCADE` consistente.
    - Criar índices: `UNIQUE (listing_uuid)` em `listing_identities` e índice composto `(listing_identity_id, version)` em `listing_versions`.
    - Atualizar views/consultas dependentes para o novo relacionamento `listing_identities ↔ listing_versions`.

### Fase 3 — Services
- Atualizar `StartListing`, `UpdateListing`, `EndUpdateListing` (renomear para `PromoteListingVersion`) e novos serviços (`DiscardDraftVersion`, `ListListingVersions`).
- Implementar clonagem de associações via porta de repositório.
- Garantir auditoria e logging.

### Fase 4 — Handlers, DTOs e Converters
- Ajustar endpoints existentes (`PUT /listings`, `POST /listings`, `POST /listings/end-update`, etc.).
- Expor novos campos (`listingUuid`, `activeVersion`, `draftVersionId`).
- Atualizar converters HTTP ↔ domínio.
- Rodar `make swagger` e validar anotações.

### Fase 5 — Integrações e Serviços Dependentes
- Revisar `schedule_service`, `photo_session_service`, `offer/visit` workflows.
- Garantir que consultem versão ativa correta e, quando necessário, draft explícita.
- Ajustar caches ou adaptadores que persistem `listing_id` para trabalhar com `listing_uuid`/`version`.

### Fase 6 — Documentação e Cleanup
- Atualizar docs (`README`,`procedimento_de_criação_de_novo_anuncio.md`) com novo fluxo de versionamento.
- Remover código legado/tags TODO utilizados como ponte.
- Confirmar build final e executar checklist das seções 13/14 do guia.

## Controle de Progresso
- [x] Fase 1 — Modelos e Interfaces
- [x] Fase 2 — Ports e Adapters
- [x] Fase 3 — Services
- [x] Fase 4 — Handlers & DTOs
- [ ] Fase 5 — Integrações
- [ ] Fase 6 — Documentação & Cleanup

Cada fase concluída deve marcar a checkbox correspondente (commit separado preferencialmente) e incluir um bloco "Notas de Migração" no PR.
