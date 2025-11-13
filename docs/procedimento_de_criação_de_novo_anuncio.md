Procedimento de criação de novo anuncio:

## Conceito de Versionamento

O sistema utiliza **versionamento de listings** para preservar o histórico e permitir edições não-destrutivas:

- **Listing Identity** (`listing_identities`): Representa o imóvel único, identificado por UUID. Contém metadados compartilhados (user_id, code, active_version_id).
- **Listing Version** (`listing_versions`): Cada alteração cria uma nova versão vinculada à identity. Versões draft podem ser promovidas à ativa, mantendo o histórico completo.
- **Versão Ativa**: Apenas uma versão por identity está ativa por vez. Mudanças de status (pendências, aprovações) aplicam-se à versão ativa.
- **Fluxo de Edição**: Para alterar um listing já ativo, crie uma nova versão draft, valide-a e promova via `POST /listings/versions/promote`.

## Fluxo de Criação

1 - POST `/listings/options` - Buscar as opções de imovel possíveis no cep/numero
2 - POST `/listings` - Cria o anuncio com as informações básicas com `StatusDraft`
	2.1 - Cria automaticamente a **listing identity** (UUID) e a primeira **versão** (v1)
	2.2 - Utilizar POST `/auth/validate/cep` para obter o endereço completo permitindo ao usuário ajustes de complemento e bairro
3 - PUT `/listings` - quantos necessários para preencher todos os dados do anuncio. Neste momento nenhuma validação é feita, apenas grava os dados informados.
	3.1 - Atualiza a versão draft atual (v1 ou versão draft criada posteriormente)
	3.2 - Utilizar GET `/listings/catalog` para obter Available categories: property_owner, property_delivered, who_lives, transaction_type, installment_plan, financing_blocker, visit_type, accompanying_type, guarantee_type.
	3.3 - Utilizar GET `/listings/features/base` para obter as features possíveis de serem incluídas
	3.4 - Utilizar GET `/complex/sizes` para obter os tamanhos das plantas padrão para edificios. Mas o usuário pode digitar o que quiser
4 - POST `/listings/versions/promote` - Efetua todas as validações e caso esteja tudo bem, promove a versão draft para ativa
	4.1 - Se for a primeira versão (v1), muda o status para `StatusPendingAvailability` e cria a agenda básica do imóvel
	4.2 - Se for uma versão posterior, mantém o status da versão ativa anterior (preserva o ciclo de vida do listing)
	4.3 - Abaixo as regras de validação atuais, informando campos obrigatório por situação

## Endpoints de Versionamento

- **POST** `/listings/versions` - Lista todas as versões de um listing (body: `{"listingIdentityId": <id>, "includeDeleted": false}`)
- **POST** `/listings/versions/promote` - Promove versão draft para ativa
- **DELETE** `/listings/versions/discard` - Descarta versão draft não promovida
5 - GET/POST/PUT/DELETE `/schedules/listing/**` altera a agenda básica do imóvel, através de bloqueios semanais para definir quando o proprietário autoriza visitas
6 - POST `/schedules/listing/finish` confirma fim da criação da agenda do imóvel e altera o status para `StatusPendingPhotoScheduling`
	6.1 - GET `/schedules/owner/summary` apresenta a agenda consolidada do proprietário, caso tenha mais de um imóvel.
7 - GET `/listings/photo-session/slots` apresenta as disponibilidades de fotografo para a sessão de fotos do imóvel. Usuário seleciona uma.
	7.1 - GET/POST/PUT/DELETE de `/photographer/service-area/**` permite que o fotografo defina a sua área de atuação (cidade e estado)
	7.2 - GET/POST/PUT/DELETE de `/photographer/agenda/time-off/**` permite que o fotografo bloqueie horários de sua agenda, não permitido que proprietário agende sessão de fotos
8 - POST `/listings/photo-session/reserve` solicita o slot escolhido pelo usuário.  Status muda para `StatusPendingPhotoConfirmation`
	8.1 - O fotografo é avisado por push notification
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