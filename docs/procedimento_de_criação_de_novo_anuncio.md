Procedimento de criação de novo anuncio:

## Conceito de Versionamento

O sistema utiliza **versionamento de listings** para preservar o histórico e permitir edições não-destrutivas:

- **Listing Identity** (`listing_identities`): Representa o imóvel único, identificado por UUID. Contém metadados compartilhados (user_id, code, active_version_id).
- **Listing Version** (`listing_versions`): Cada alteração cria uma nova versão vinculada à identity. Versões draft podem ser promovidas à ativa, mantendo o histórico completo.
- **Versão Ativa**: Apenas uma versão por identity está ativa por vez. Mudanças de status (pendências, aprovações) aplicam-se à versão ativa.
- **Fluxo de Edição**: Para alterar um listing já ativo, crie uma nova versão draft, valide-a e promova via `POST /listings/versions/promote`.

## Fluxo de Criação

1 - POST `/listings/options` - Buscar as opções de imóvel possíveis e dados completos do condomínio (se houver) no cep/número

2 - POST `/listings` - Cria o anuncio com as informações básicas com `StatusDraft`
	2.1 - Cria automaticamente a **listing identity** (UUID) e a primeira **versão** (v1)
	2.2 - Valida se já existe listing ativo/publicado no mesmo endereço (retorna 409 Conflict se houver duplicidade)
	2.3 - Utilizar POST `/auth/validate/cep` para obter o endereço completo permitindo ao usuário ajustes de complemento e bairro
	
3 - PUT `/listings` - Quantos necessários para preencher todos os dados do anuncio. Neste momento nenhuma validação é feita, apenas grava os dados informados.
	3.1 - **REQUER** campos `listingIdentityId` (int64) e `listingVersionId` (int64) no body (obrigatórios para identificar o listing e qual versão está sendo editada)
	3.2 - Valida se a versão está em `StatusDraft` (retorna 409 Conflict caso contrário)
	3.3 - Atualiza a versão draft atual (v1 ou versão draft criada posteriormente)
	3.4 - Utilizar GET `/listings/catalog` para obter Available categories: property_owner, property_delivered, who_lives, transaction_type, installment_plan, financing_blocker, visit_type, accompanying_type, guarantee_type, land_terrain_type, warehouse_sector.
	3.5 - Utilizar GET `/listings/features/base` para obter as features possíveis de serem incluídas
	3.6 - Utilizar POST `/listings/options` para obter os dados do condomínio (tamanhos, torres, etc) a partir do CEP e número.
	
	3.7 - **Campos do Body do UpdateListingRequest**:
	
	**Campos Obrigatórios (sempre)**:
	- `listingIdentityId` (int64) - ID da identity do listing
	- `listingVersionId` (int64) - ID da versão sendo editada
	
	**Campos Opcionais (Optional[T] - podem ser omitidos, nulos ou com valor)**:
	
	*Informações Básicas*:
	- `owner` (string) - Slug do catálogo property_owner (ex: "myself", "family", "third_party")
	- `title` (string) - Título do anúncio
	- `description` (string) - Descrição detalhada do imóvel
	- `features` (array) - Array de objetos {featureId: int64, quantity: int}
	
	*Dimensões e Características Físicas*:
	- `landSize` (float64) - Área total do terreno em m²
	- `corner` (bool) - Se o imóvel é de esquina
	- `nonBuildable` (float64) - Área não edificável em m²
	- `buildable` (float64) - Área edificável em m²
	- `delivered` (string) - Slug do catálogo property_delivered (ex: "furnished", "unfurnished", "semi_furnished")
	- `whoLives` (string) - Slug do catálogo who_lives (ex: "owner", "tenant", "vacant")
	
	*Transação e Valores*:
	- `transaction` (string) - Slug do catálogo transaction_type (ex: "sale", "rent", "both")
	- `sellNet` (float64) - Valor líquido de venda
	- `rentNet` (float64) - Valor líquido de aluguel
	- `condominium` (float64) - Valor do condomínio
	- `annualTax` (float64) - IPTU anual (mutuamente exclusivo com monthlyTax)
	- `monthlyTax` (float64) - IPTU mensal (mutuamente exclusivo com annualTax)
	- `annualGroundRent` (float64) - Laudêmio anual (mutuamente exclusivo com monthlyGroundRent)
	- `monthlyGroundRent` (float64) - Laudêmio mensal (mutuamente exclusivo com annualGroundRent)
	
	*Permuta*:
	- `exchange` (bool) - Se aceita permuta
	- `exchangePercentual` (float64) - Percentual de permuta aceito
	- `exchangePlaces` (array) - Array de objetos {neighborhood: string, city: string, state: string}
	
	*Financiamento*:
	- `installment` (string) - Slug do catálogo installment_plan
	- `financing` (bool) - Se aceita financiamento
	- `financingBlockers` (array) - Array de slugs do catálogo financing_blocker
	
	*Garantias (para locação)*:
	- `guarantees` (array) - Array de objetos {priority: int, guarantee: string (slug do catálogo)}
	
	*Visitação*:
	- `visit` (string) - Slug do catálogo visit_type (ex: "owner", "client", "both")
	- `accompanying` (string) - Slug do catálogo accompanying_type (ex: "assistant", "broker")
	
	*Inquilino (quando whoLives = "tenant")*:
	- `tenantName` (string) - Nome do inquilino
	- `tenantEmail` (string) - Email do inquilino
	- `tenantPhone` (string) - Telefone no formato E.164 (ex: "+5511912345678")
	
	*Campos Específicos por Tipo de Imóvel*:
	
	**Casa em Construção (256)**:
	- `completionForecast` (string) - Previsão de conclusão no formato YYYY-MM-DD (ex: "2026-06-01")
	
	**Terrenos (16=Urbano, 32=Rural, 64=Industrial, 128=Comercial, 512=Residencial)**:
	- `landBlock` (string) - Quadra/Bloco (ex: "A", "B1")
	- `landLot` (string) - Lote (ex: "15", "23-A") - **obrigatório para Comercial(64) e Residencial(512)**
	- `landFront` (float64) - Frente do terreno em metros
	- `landSide` (float64) - Lateral do terreno em metros
	- `landBack` (float64) - Fundo do terreno em metros
	- `landTerrainType` (string) - Slug do catálogo land_terrain_type (ex: "flat", "uphill", "downhill", "slight_uphill", "slight_downhill") - **obrigatório para Comercial(64) e Residencial(512)**
	- `hasKmz` (bool) - Se possui arquivo KMZ - **obrigatório para Comercial(64) e Residencial(512)**
	- `kmzFile` (string) - URL do arquivo KMZ - **obrigatório para Comercial(64) se hasKmz=true**
	
	**Prédio (1024)**:
	- `buildingFloors` (int16) - Número total de andares do prédio
	
	**Apartamento (1), Sala (2), Laje Corporativa (4)**:
	- `unitTower` (string) - Identificação da torre (ex: "Torre A", "Bloco B")
	- `unitFloor` (int16) - Andar da unidade (ex: 5, 12)
	- `unitNumber` (string) - Número/identificação da unidade (ex: "502", "1201-A")
	
	**Galpão/Industrial/Logístico (2048)**:
	- `warehouseManufacturingArea` (float64) - Área de produção/manufatura em m²
	- `warehouseSector` (string) - Slug do catálogo warehouse_sector (ex: "manufacturing", "industrial", "logistics")
	- `warehouseHasPrimaryCabin` (bool) - Se possui cabine primária de energia
	- `warehouseCabinKva` (float64) - Potência da cabine em KVA - **obrigatório se warehouseHasPrimaryCabin=true**
	- `warehouseGroundFloor` (float64) - Pé direito do piso térreo em metros
	- `warehouseFloorResistance` (float64) - Resistência do piso em kg/m²
	- `warehouseZoning` (string) - Zoneamento (ex: "ZI-1", "ZI-2", "ZL-3")
	- `warehouseHasOfficeArea` (bool) - Se possui área de escritórios
	- `warehouseOfficeArea` (float64) - Área de escritórios em m² - **obrigatório se warehouseHasOfficeArea=true**
	- `warehouseAdditionalFloors` (array) - Andares adicionais além do térreo, array de objetos:
	  - `floorName` (string) - Nome do andar (ex: "Mezanino", "Segundo Piso")
	  - `floorOrder` (int) - Ordem do andar (1=primeiro acima do térreo, 2=segundo, etc)
	  - `floorHeight` (float64) - Pé direito em metros
	
	**Loja (8)**:
	- `storeHasMezzanine` (bool) - Se possui mezanino
	- `storeMezzanineArea` (float64) - Área do mezanino em m² - **obrigatório se storeHasMezzanine=true**
	
	**Exemplo de Body Completo**:
	```json
	{
	  "listingIdentityId": 1024,
	  "listingVersionId": 5001,
	  "owner": "myself",
	  "title": "Apartamento 3 dormitórios com piscina",
	  "description": "Apartamento amplo com vista panorâmica",
	  "features": [
	    {"featureId": 101, "quantity": 3},
	    {"featureId": 205, "quantity": 2}
	  ],
	  "landSize": 423.5,
	  "corner": true,
	  "nonBuildable": 12.75,
	  "buildable": 410.75,
	  "delivered": "furnished",
	  "whoLives": "owner",
	  "transaction": "sale",
	  "sellNet": 1200000,
	  "rentNet": null,
	  "condominium": 1200.5,
	  "annualTax": 3400.75,
	  "exchange": true,
	  "exchangePercentual": 50,
	  "exchangePlaces": [
	    {"neighborhood": "Vila Mariana", "city": "São Paulo", "state": "SP"}
	  ],
	  "financing": true,
	  "guarantees": [
	    {"priority": 1, "guarantee": "security_deposit"}
	  ],
	  "visit": "both",
	  "accompanying": "assistant",
	  "unitTower": "Torre B",
	  "unitFloor": 5,
	  "unitNumber": "502"
	}
	```

3.5 - POST `/listings/versions/draft` - Cria nova versão draft a partir da versão ativa (para editar listing já publicado)
	3.5.1 - **REQUER** campo `listingIdentityId` (int64) no body
	3.5.2 - **Status permitidos para criar draft**: `StatusPendingAvailability` (2), `StatusPendingPhotoScheduling` (3), `StatusPendingPhotoConfirmation` (4), `StatusPhotosScheduled` (5), `StatusPendingPhotoProcessing` (6), `StatusPendingAdminReview` (8), `StatusSuspended` (14)
	3.5.3 - **Status Published (10)**: Retorna 409 Conflict - "Cannot create draft from published listing"
	3.5.4 - **Status UnderOffer/UnderNegotiation (11/12) ou StatusRejectedByOwner (9)**: Retorna 423 Locked - "Listing is locked for draft creation"
	3.5.5 - **Status Closed/Expired/Archived (13/15/16)**: Retorna 410 Gone - "Listing is finalized"
	3.5.6 - Copia todos os campos mutáveis e entidades satélite (features, exchange_places, guarantees, financing_blockers) da versão ativa
	3.5.7 - Apenas 1 draft pode coexistir com 1 versão ativa por identity

4 - POST `/listings/versions/promote` - Efetua todas as validações e caso esteja tudo bem, promove a versão draft para ativa
	4.1 - **REQUER** campos `listingIdentityId` (int64) e `versionId` (int64) no body
	4.2 - **Se for a primeira versão (v1)**: Muda o status para `StatusPendingAvailability` e cria a agenda básica do imóvel
	4.3 - **Se for uma versão posterior (v>1)**: Herda o status da versão ativa anterior (preserva o ciclo de vida do listing)
	4.4 - Atualiza o campo `active_version_id` na listing identity para apontar para a nova versão ativa
	
	4.5 - **Regras de Validação do Promote (campos obrigatórios)**:
	
	**Campos Básicos (obrigatórios para qualquer tipo de imóvel)**:
	- `code` - Código do listing
	- `version` - Número da versão
	- `zipCode` - CEP
	- `street` - Logradouro
	- `number` - Número
	- `city` - Cidade
	- `state` - Estado
	- `title` - Título do anúncio
	- `listingType` - Tipo(s) de imóvel (bitmask)
	- `owner` - Dono do imóvel (catálogo property_owner)
	- `buildable` - Área edificável
	- `delivered` - Status de entrega (catálogo property_delivered)
	- `whoLives` - Quem mora (catálogo who_lives)
	- `description` - Descrição do imóvel
	- `transaction` - Tipo de transação (catálogo transaction_type)
	- `visit` - Tipo de visita (catálogo visit_type)
	- `accompanying` - Tipo de acompanhamento (catálogo accompanying_type)
	- IPTU: **Exatamente um** dos campos deve estar preenchido (nunca ambos):
	  - `annualTax` - IPTU anual **OU**
	  - `monthlyTax` - IPTU mensal
	- Laudêmio: **Opcional**, mas se informado apenas um (nunca ambos):
	  - `annualGroundRent` - Laudêmio anual **OU**
	  - `monthlyGroundRent` - Laudêmio mensal
	- `features` - Pelo menos 1 feature cadastrada (features_count > 0)
	
	**Validações Condicionais por Tipo de Transação**:
	
	*Se transaction = "sale" ou "both" (venda)*:
	- `saleNet` (sellNet) - Valor líquido de venda
	- `exchange` - Flag de permuta (true/false)
	- Se `exchange = true`:
	  - `exchangePercentual` - Percentual de permuta
	  - `exchangePlaces` - Pelo menos 1 local de permuta cadastrado (exchange_places_count > 0)
	- `financing` - Flag de financiamento (true/false)
	- Se `financing = false`:
	  - `financingBlockers` - Pelo menos 1 impeditivo cadastrado (financing_blockers_count > 0)
	
	*Se transaction = "rent" ou "both" (locação)*:
	- `rentNet` - Valor líquido de aluguel
	- `guarantees` - Pelo menos 1 garantia cadastrada (guarantees_count > 0)
	
	**Validações Condicionais por Tipo de Imóvel**:
	
	*Se Apartamento (1) ou Laje Corporativa (4)*:
	- `condominium` - Valor do condomínio
	
	*Se Terrenos (16, 32, 64, 128, 512)*:
	- `landSize` - Área do terreno
	- `corner` - Flag de esquina
	
	*Se whoLives = "tenant" (inquilino)*:
	- `tenantName` - Nome do inquilino
	- `tenantPhone` - Telefone do inquilino (formato E.164)
	- `tenantEmail` - Email do inquilino
	
	**Validações Específicas por Tipo de Imóvel (bitmask)**:
	
	*Prédio (code: 256)*:
	- `completionForecast` - Previsão de conclusão no formato YYYY-MM
	
	*Terreno Residencial (code: 64) ou Terreno Comercial (code: 128)*:
	- `landBlock` - Quadra/Bloco
	- `landLot` - Número do lote
	- `landTerrainType` - Tipo do terreno (catálogo land_terrain_type)
	- `hasKmz` - Flag indicando se possui arquivo KMZ
	- Se `hasKmz = true` **E é Terreno Comercial (128)**:
	  - `kmzFile` - Caminho/URL do arquivo KMZ
	
	*Apartamento (code: 1), Loja (code: 2), Laje Corporativa (code: 4)*:
	- `unitTower` - Torre/Bloco da unidade
	- `unitFloor` - Andar da unidade
	- `unitNumber` - Número da unidade
	
	*Galpão/Industrial/Logístico (code: 512)*:
	- `warehouseManufacturingArea` - Área de manufatura/produção
	- `warehouseSector` - Setor do galpão (catálogo warehouse_sector)
	- `warehouseHasPrimaryCabin` - Flag de cabine primária
	- Se `warehouseHasPrimaryCabin = true`:
	  - `warehouseCabinKva` - Potência da cabine em KVA
	- `warehouseGroundFloor` - Pé direito do piso térreo
	- `warehouseFloorResistance` - Resistência do piso em kg/m²
	- `warehouseZoning` - Zoneamento
	- `warehouseHasOfficeArea` - Flag de área de escritórios
	- Se `warehouseHasOfficeArea = true`:
	  - `warehouseOfficeArea` - Área de escritórios em m²
	
	*Loja (code: 2)*:
	- `storeHasMezzanine` - Flag de mezanino
	- Se `storeHasMezzanine = true`:
	  - `storeMezzanineArea` - Área do mezanino em m²
	
	**Observações Importantes**:
	- Os códigos de tipo de imóvel são bitmask, um listing pode ter múltiplos tipos simultaneamente
	- Códigos válidos: 1=Apartamento, 2=Loja, 4=Laje, 8=Sala, 16=Casa, 32=Casa na Planta, 64=Terreno Residencial, 128=Terreno Comercial, 256=Prédio, 512=Galpão
	- Durante o `PUT /listings` nenhuma validação é feita, os dados são apenas gravados
	- Todas as validações acima são executadas apenas no `POST /listings/versions/promote`
	- Se alguma validação falhar, o promote retorna 400 Bad Request com mensagem específica do campo faltante
	- Campos opcionais no DTO podem ser omitidos, enviados como null, ou com valor durante o update
	- Catálogos disponíveis: property_owner, property_delivered, who_lives, transaction_type, installment_plan, financing_blocker, visit_type, accompanying_type, guarantee_type, land_terrain_type, warehouse_sector

## Endpoints de Versionamento

- **POST** `/listings/versions/draft` - Cria nova versão draft a partir da versão ativa atual (body: `{"listingIdentityId": 1024}`)
- **POST** `/listings/versions` - Lista todas as versões de um listing (body: `{"listingIdentityId": 1024, "includeDeleted": false}`)
- **POST** `/listings/versions/promote` - Promove versão draft para ativa (body: `{"listingIdentityId": 1024, "versionId": 5001}`)
- **DELETE** `/listings/versions/discard` - Descarta versão draft não promovida (body: `{"listingIdentityId": 1024, "versionId": 5001}`)

## Hierarquia de Validação

Todos os endpoints de listing seguem o seguinte padrão de validação para garantir segurança e consistência:

1. **Validar `listingIdentityId`**: Verificar se o ID da identity foi fornecido e é válido
2. **Buscar Identity**: Localizar o registro `listing_identity` correspondente no banco
3. **Validar Ownership**: Comparar `identity.user_id` com o `user_id` do contexto autenticado
   - Se divergir: registrar log de auditoria `unauthorized_<operation>_attempt` com campos `listing_identity_id`, `listing_version_id` (se aplicável), `requester_user_id`, `owner_user_id`
   - Retornar HTTP 403 Forbidden com mensagem "not the listing owner"
4. **Buscar Version**: Localizar a versão específica (`versionId`) ou a versão ativa conforme o endpoint
5. **Verificar Relacionamento**: Confirmar que `version.identity_id == input.listingIdentityId`
   - Se divergir: retornar HTTP 400 Bad Request com mensagem "version does not belong to this listing"
6. **Validar Regras de Negócio**: Executar validações específicas do endpoint (status, campos obrigatórios, etc.)

Esta hierarquia previne ambiguidade na identificação de listings e garante que usuários não possam acessar ou modificar listings de terceiros.

5 - GET/POST/PUT/DELETE `/schedules/listing/**` altera a agenda básica do imóvel, através de bloqueios semanais para definir quando o proprietário autoriza visitas

6 - POST `/schedules/listing/finish` confirma fim da criação da agenda do imóvel e altera o status para `StatusPendingPhotoScheduling`
	6.1 - GET `/schedules/owner/summary` apresenta a agenda consolidada do proprietário, caso tenha mais de um imóvel.

7 - GET `/listings/photo-session/slots` apresenta as disponibilidades de fotografo para a sessão de fotos do imóvel. Usuário seleciona uma.
	7.1 - **REQUER** query param `listingIdentityId` (int64) para identificar o listing
	7.2 - GET/POST/PUT/DELETE de `/photographer/service-area/**` permite que o fotografo defina a sua área de atuação (cidade e estado)
	7.3 - GET/POST/PUT/DELETE de `/photographer/agenda/time-off/**` permite que o fotografo bloqueie horários de sua agenda, não permitido que proprietário agende sessão de fotos

8 - POST `/listings/photo-session/reserve` solicita o slot escolhido pelo usuário. Status muda para `StatusPendingPhotoConfirmation`
	8.1 - **REQUER** campos `listingIdentityId` (int64) e `slotId` (int64) no body
	8.2 - O fotografo é avisado por push notification

9 - POST `/photographer/sessions/status` 0 fotografo confirma o aceite ou a recusa da sessão de fotos solicitada ==> esta etapa está configurada para não ser executada e ser autoaprovada pelo sistema
	9.1 - Caso aceite, o status do listing muda para `StatusPhotosScheduled`
	9.2 - Caso recuse, o status do listing volta para `StatusPendingPhotoScheduling`, permitindo que o usuário escolha outro slot
	9.3 - Em ambos os casos o proprietário é avisado por push notification

10 - POST `/photographer/sessions/status` o fotografo confirma a realização da sessão de fotos
	10.1 - O status do listing muda para `StatusPendingPhotoProcessing`
	10.2 - O proprietário é avisado por push notification

11 - ainda pendentes.....  passar a aprovação do owner, aprovar e publicar.

____________________________

O envio de push por FCM ao aprovar o cadastro do corretor está funcionando, entretanto, faltam algumas coisas:
1) o sininho da home tem que indicar com badget numérico, quantas mensagens não lidas existe;
2) Ao chegar nova mensagem, verificar novamente o status da conta, e fazer refresh da home, pois a aprovação permite acesso a home.
3) Ao clicar no sininho, deve abrir uma tela com todas as notificações, indicando lidas e não lidas e opção de deletar 1 a uma e todas, e marcar como lidas 1 a uma ou todas, e ao ler marcar como lidas
4) Algumas mensagens devem redirecionar para uma página espécífica da aplicaçÃo. Precisamos decidir juntos o que devo enviar na msg para voce poder implementar isso

____________________________


Possíveis estados do listing:
// StatusDraft: O anúncio está sendo criado pelo proprietário e permanece invisível ao público.
StatusDraft ListingStatus = iota + 1
// StatusPendingAvailability: Anúncio criado, aguardando criação de agenda de disponibilidades do imóvel
StatusPendingAvailability	// StatusPendingPhotoScheduling: Anúncio criado e aguardando o agendamento da sessão de fotos.
// StatusPendingPhotoScheduling: Anúncio criado e aguardando o agendamento da sessão de fotos.
StatusPendingPhotoScheduling
// StatusPendingPhotoConfirmation: Solicitado fotos para slot disponível. ag confirmação
StatusPendingPhotoConfirmation
// StatusPhotosScheduled: Sessão de fotos agendada, aguardando execução.
StatusPhotosScheduled
// StatusPendingPhotoProcessing: Sessão concluída, aguardando tratamento e upload das fotos.
StatusPendingPhotoProcessing
// StatusPendingOwnerApproval: Materiais revisados e aguardando aprovação final do proprietário.
StatusPendingOwnerApproval
// StatusRejectedByOwner: Versão final reprovada pelo proprietário (ex.: não aprovou as fotos).
StatusRejectedByOwner
// StatusPendingAdminReview: Proprietário aprovou, aguardando revisão do time administrativo antes da publicação.
StatusPendingAdminReview
// StatusPublished: Anúncio ativo e visível publicamente.
StatusPublished
// StatusUnderOffer: Anúncio publicado que recebeu uma ou mais propostas.
StatusUnderOffer
// StatusUnderNegotiation: Uma proposta foi aceita e a negociação está em andamento.
StatusUnderNegotiation
// StatusClosed: O imóvel foi comercializado (vendido ou alugado) e o processo foi encerrado.
StatusClosed
// StatusSuspended: Anúncio pausado temporariamente pelo proprietário ou administrador.
StatusSuspended
// StatusExpired: Prazo de validade do anúncio encerrou sem negociação concluída.
StatusExpired
// StatusArchived: Anúncio removido do catálogo e mantido apenas para histórico.
StatusArchived
// StatusNeedsRevision: Anúncio reprovado e aguardando ajustes antes de retornar ao fluxo de criação.
StatusNeedsRevision
____________________________

Validações do listing:

**NOTA**: As validações abaixo foram migradas para a seção "4.5 - Regras de Validação do Promote" no fluxo de criação acima. Esta seção está mantida por compatibilidade mas recomenda-se consultar a seção 4.5 para a documentação mais completa e estruturada.

**Resumo das validações por categoria**:


	// Campos básicos obrigatórios para qualquer anúncio em draft.
	if data.Code == 0 {
		return utils.BadRequest("Listing code is required")
	}
	if data.Version == 0 {
		return utils.BadRequest("Listing version is required")
	}
	if strings.TrimSpace(data.ZipCode) == "" {
		return utils.BadRequest("Zip code is required")
	}
	if !data.Street.Valid || strings.TrimSpace(data.Street.String) == "" {
		return utils.BadRequest("Street is required")
	}
	if !data.Number.Valid || strings.TrimSpace(data.Number.String) == "" {
		return utils.BadRequest("Number is required")
	}
	if !data.City.Valid || strings.TrimSpace(data.City.String) == "" {
		return utils.BadRequest("City is required")
	}
	if !data.State.Valid || strings.TrimSpace(data.State.String) == "" {
		return utils.BadRequest("State is required")
	}
	if data.ListingType == 0 {
		return utils.BadRequest("Property type is required")
	}
	if !data.Owner.Valid {
		return utils.BadRequest("Property owner is required")
	}
	if !data.Buildable.Valid {
		return utils.BadRequest("Buildable size is required")
	}
	if !data.Delivered.Valid {
		return utils.BadRequest("Delivered status is required")
	}
	if !data.WhoLives.Valid {
		return utils.BadRequest("Who lives information is required")
	}
	if !data.Description.Valid || strings.TrimSpace(data.Description.String) == "" {
		return utils.BadRequest("Description is required")
	}
	if !data.Transaction.Valid {
		return utils.BadRequest("Transaction type is required")
	}
	if !data.Visit.Valid {
		return utils.BadRequest("Visit type is required")
	}
	if !data.Accompanying.Valid {
		return utils.BadRequest("Accompanying type is required")
	}
	if !data.AnnualTax.Valid {
		return utils.BadRequest("Annual tax is required")
	}
	if data.FeaturesCount == 0 {
		return utils.BadRequest("Listing must include features")
	}

	// Regras condicionais para o tipo de transação.
	txnValue := uint8(data.Transaction.Int16)
	txnCatalog, err := ls.listingRepository.GetCatalogValueByNumeric(ctx, tx, listingmodel.CatalogCategoryTransactionType, txnValue)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.BadRequest("Transaction type is invalid")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("listing.end_update.transaction_catalog_error", "err", err, "listing_id", data.ListingID, "transaction_id", txnValue)
		return utils.InternalError("")
	}

	slug := strings.ToLower(strings.TrimSpace(txnCatalog.Slug()))
	needsSaleValidation := slug == "sale" || slug == "both"
	needsRentValidation := slug == "rent" || slug == "both"

	if needsSaleValidation {
		// Quando a transação envolve venda, validamos preço líquido, permuta e barreiras de financiamento.
		if !data.SaleNet.Valid {
			return utils.BadRequest("Sale net value is required")
		}
		if !data.Exchange.Valid {
			return utils.BadRequest("Exchange flag is required")
		}
		if data.Exchange.Valid && data.Exchange.Int16 == 1 {
			if !data.ExchangePercentual.Valid {
				return utils.BadRequest("Exchange percentual is required when exchange is enabled")
			}
			if data.ExchangePlacesCount == 0 {
				return utils.BadRequest("Exchange places are required when exchange is enabled")
			}
		}
		if !data.Financing.Valid {
			return utils.BadRequest("Financing flag is required")
		}
		if data.Financing.Int16 == 0 && data.FinancingBlockersCount == 0 {
			return utils.BadRequest("Financing blockers are required when financing is disabled")
		}
	}

	if needsRentValidation {
		// Nas locações exigimos valor líquido e garantias cadastradas para prosseguir.
		if !data.RentNet.Valid {
			return utils.BadRequest("Rent net value is required")
		}
		if data.GuaranteesCount == 0 {
			return utils.BadRequest("Guarantees are required for rent transactions")
		}
	}

	// Regras específicas por tipo de imóvel.
	propertyOptions := ls.DecodePropertyTypes(ctx, data.ListingType)
	if len(propertyOptions) == 0 {
		return utils.BadRequest("Property type is invalid")
	}
	// Cada option representa um bit ativo na máscara do tipo; usamos isso para derivar validações adicionais.
	needsCondominium := false
	needsLandData := false
	for _, option := range propertyOptions {
		switch option.Code {
		case 1, 4:
			needsCondominium = true
		case 16, 32, 64, 128:
			needsLandData = true
		}
	}

	if needsCondominium && !data.Condominium.Valid {
		return utils.BadRequest("Condominium value is required for the selected property type")
	}

	if needsLandData {
		if !data.LandSize.Valid {
			return utils.BadRequest("Land size is required for the selected property type")
		}
		if !data.Corner.Valid {
			return utils.BadRequest("Corner information is required for the selected property type")
		}
	}

	// Regras adicionais quando quem mora é inquilino.
	if data.WhoLives.Valid {
		whoLivesValue := uint8(data.WhoLives.Int16)
		whoLivesCatalog, catalogErr := ls.listingRepository.GetCatalogValueByNumeric(ctx, tx, listingmodel.CatalogCategoryWhoLives, whoLivesValue)
		if catalogErr != nil {
			if errors.Is(catalogErr, sql.ErrNoRows) {
				return utils.BadRequest("Who lives value is invalid")
			}
			utils.SetSpanError(ctx, catalogErr)
			logger.Error("listing.end_update.wholives_catalog_error", "err", catalogErr, "listing_id", data.ListingID, "who_lives_id", whoLivesValue)
			return utils.InternalError("")
		}

		if strings.ToLower(strings.TrimSpace(whoLivesCatalog.Slug())) == "tenant" {
			if !data.TenantName.Valid || strings.TrimSpace(data.TenantName.String) == "" {
				return utils.BadRequest("Tenant name is required when tenant lives in the property")
			}
			if !data.TenantPhone.Valid || strings.TrimSpace(data.TenantPhone.String) == "" {
				return utils.BadRequest("Tenant phone is required when tenant lives in the property")
			}
			if !data.TenantEmail.Valid || strings.TrimSpace(data.TenantEmail.String) == "" {
				return utils.BadRequest("Tenant email is required when tenant lives in the property")
			}
		}
	}