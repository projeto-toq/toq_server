# Status da Implementação - Refatoração de Versionamento de Listings

## ✅ Concluído

### Step 1: Remoção de Código Legado
- ✅ Deletados: `create_listing.go`, `update_listing.go`, `get_listing_by_id.go`, `get_listing_by_zip_number.go` do adapter
- ✅ Interface `listing_repository_interface.go` atualizada removendo métodos legados
- ✅ Novos métodos adicionados à interface:
  - `CheckActiveListingExists`
  - `GetListingVersionByAddress`
  - `GetActiveListingVersion`
  - `GetPreviousActiveVersionStatus`
  - `UpdateListingVersion`
  - `CloneListingVersionSatellites`
- ✅ Parâmetros renomeados: `listingID` → `listingVersionID` em satellite operations

### Step 2: Renomeação de Arquivos
- ✅ `start_listing.go` → `create_listing.go` (service)
- ✅ `start_listing.go` → `create_listing_handler.go` (handler)
- ✅ Métodos renomeados: `StartListing` → `CreateListing`
- ✅ Input/Output types renomeados: `StartListingInput` → `CreateListingInput`

### Step 3: Validações de Unicidade
- ✅ Criado: `check_active_listing_exists.go` - valida se usuário já tem listing ativo
- ✅ Criado: `get_listing_version_by_address.go` - valida endereço duplicado
- ✅ `create_listing.go` service atualizado com validações integradas

### Repository Files Criados
- ✅ `check_active_listing_exists.go`
- ✅ `get_listing_version_by_address.go`
- ✅ `get_active_listing_version.go`
- ✅ `get_previous_active_version_status.go`
- ✅ `update_listing_version.go`
- ✅ `clone_listing_version_satellites.go`

### Service Files Criados
- ✅ `create_draft_version.go` (parcial - necessita integração completa)

## 🚧 Pendente

### Step 4: Endpoint POST /listings/versions/draft
- ⏳ Handler: `create_draft_version_handler.go` - CRIAR
- ⏳ DTOs em `listing_dto.go` - ADICIONAR
  - `CreateDraftVersionRequest`
  - `CreateDraftVersionResponse`
- ⏳ Port interface: adicionar método `CreateDraftVersion` em `listing_handler_port.go`
- ⏳ Router: adicionar rota `POST /listings/versions/draft`

### Step 5: Refatorar update_listing
- ⏳ Adicionar campo `VersionID int64` em `UpdateListingInput` (`update_listing_input.go`)
- ⏳ Atualizar `update_listing.go` service para:
  - Buscar versão via `GetListingVersionByID`
  - Validar `status == StatusDraft`
  - Chamar `UpdateListingVersion` ao invés de `UpdateListing`
- ⏳ Atualizar `update_listing_handler.go` para aceitar `versionId` no body
- ⏳ Refatorar satellite update methods em repository:
  - `update_features.go` - alterar parâmetro para `listingVersionID`
  - `update_exchange_places.go` - idem
  - `update_financing_blockers.go` - idem
  - `update_guarantees.go` - idem
  - `delete_listing_features.go` - idem
  - `delete_listing_exchange_places.go` - idem
  - `delete_listing_financing_blokers.go` - idem
  - `delete_listing_guarantees.go` - idem

### Step 6: Refatorar promote_listing_version
- ⏳ Atualizar `promote_listing_version.go` service para:
  - Adicionar lógica condicional `version == 1` vs `version > 1`
  - Para v1: `StatusDraft → StatusPendingAvailability` + criar agenda
  - Para v>1: buscar status anterior e aplicar
  - Usar `listingIdentityId` ao criar agenda

### Step 7: Atualizar Satellite Entities
- ⏳ Entities (renomear `ListingID` → `ListingVersionID`):
  - `entity/features_entity.go`
  - `entity/exchange_place_entity.go`
  - `entity/financing_blocker_entity.go`
  - `entity/guarantee_entity.go`
- ⏳ Interfaces domain (já possuem métodos mas precisam validação):
  - `feature_interface.go`
  - `exchange_place_interface.go`
  - `financing_blocker_interface.go`
  - `guarantee_interface.go`
- ⏳ Converters:
  - `converters/listing_entity_to_domain.go` - atualizar conversões

### Step 8: Atualizar Documentação
- ⏳ `docs/procedimento_de_criação_de_novo_anuncio.md`:
  - Atualizar step 2 (POST /listings - validações 409)
  - Atualizar step 3 (PUT /listings - requer versionId)
  - Adicionar step 3.5 (POST /listings/versions/draft)
  - Atualizar step 4 (promoção v1 vs v>1)

### Interfaces a Atualizar
- ⏳ `listing_service.go` - adicionar métodos:
  - `CreateListing` (renomeado de StartListing)
  - `CreateDraftVersion`
- ⏳ `listing_handler_port.go` - adicionar método:
  - `CreateDraftVersionHandler`
- ⏳ `listing_handlers.go` - atualizar referências
- ⏳ Router principal - atualizar rotas

### Validações Necessárias
- ⏳ Executar `make lint` e corrigir erros
- ⏳ Verificar imports em todos os arquivos
- ⏳ Testar compilação completa
- ⏳ Validar que não há referências a métodos/arquivos deletados

## 📝 Notas de Implementação

### Decisões Tomadas (conforme respostas do usuário)
1. ✅ Validação de unicidade: 1 listing ativo + 1 draft por listing identity
2. ✅ Cascade delete mantido no schema SQL
3. ✅ Frontend fará refatoração (sem compatibilidade retroativa)
4. ✅ Erro 400 "Invalid version ID" para versionId inexistente
5. ✅ Renomear parâmetros ao invés de criar novos métodos
6. ✅ Schedule service já usa `ListingIdentityID` corretamente
7. ✅ Domain models já possuem `ListingVersionID()` methods

### Schema SQL
✅ Schema validado e correto:
- Tabelas satélite usam `listing_version_id`
- Tabelas agenda/visitas/bookings referenciam `listing_identities.id`
- `active_version_id` é `INT UNSIGNED NULL`

### Próximos Passos Imediatos
1. Criar handler `create_draft_version_handler.go`
2. Adicionar DTOs em `listing_dto.go`
3. Atualizar `update_listing.go` e `update_listing_input.go`
4. Refatorar satellite update methods
5. Atualizar `promote_listing_version.go`
6. Atualizar entities e converters
7. Atualizar interfaces service/handler/port
8. Atualizar router
9. Atualizar documentação
10. Executar `make lint` e corrigir

### Arquivos que Precisam de Atenção Especial
- `create_draft_version.go` - método `getDraftVersion` está simplificado, precisa usar adapter corretamente
- Todos os satellite update/delete methods - precisam atualizar queries para usar `listing_version_id`
- `promote_listing_version.go` - adicionar lógica v1 vs v>1
- DTOs - adicionar novos types necessários
- Router - adicionar nova rota e atualizar referências
