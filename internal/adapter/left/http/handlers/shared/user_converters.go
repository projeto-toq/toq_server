package shared

import (
	"time"

	httpmodels "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/models"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
)

// ConvertHTTPUpdateProfileUserToDomain converte modelo HTTP de update de perfil para modelo de domínio
// Esta função é específica para atualização de perfil e inclui apenas campos seguros para edição
// Campos como ID, email, phone_number, national_id e role têm processos específicos de alteração
// Parâmetros:
//   - httpUser: Modelo HTTP do usuário para update de perfil
//
// Retorna:
//   - usermodel.UserInterface: Modelo de domínio equivalente com apenas campos editáveis
func ConvertHTTPUpdateProfileUserToDomain(httpUser httpmodels.UpdateProfileUser) usermodel.UserInterface {
	// Cria nova instância do modelo de domínio
	user := usermodel.NewUser()

	// Converte apenas os campos que podem ser atualizados via update profile
	// Campos como ID, email, phone, national_id e role são excluídos intencionalmente
	if httpUser.NickName != "" {
		user.SetNickName(httpUser.NickName)
	}
	if httpUser.BornAt != "" {
		// Converte a string de data para time.Time
		if bornAt, err := time.Parse(time.RFC3339, httpUser.BornAt); err == nil {
			user.SetBornAt(bornAt)
		}
	}
	if httpUser.ZipCode != "" {
		user.SetZipCode(httpUser.ZipCode)
	}
	if httpUser.Street != "" {
		user.SetStreet(httpUser.Street)
	}
	if httpUser.Number != "" {
		user.SetNumber(httpUser.Number)
	}
	if httpUser.Complement != "" {
		user.SetComplement(httpUser.Complement)
	}
	if httpUser.Neighborhood != "" {
		user.SetNeighborhood(httpUser.Neighborhood)
	}
	if httpUser.City != "" {
		user.SetCity(httpUser.City)
	}
	if httpUser.State != "" {
		user.SetState(httpUser.State)
	}

	return user
}

// ConvertDomainUserToProfileHTTP converte modelo de domínio de usuário para modelo HTTP completo de perfil
// Esta função inclui todos os campos disponíveis exceto password para resposta do get profile
// Parâmetros:
//   - domainUser: Modelo de domínio do usuário a ser convertido
//
// Retorna:
//   - httpmodels.ProfileUser: Modelo HTTP completo para resposta da API
func ConvertDomainUserToProfileHTTP(domainUser usermodel.UserInterface) httpmodels.ProfileUser {
	// Determina o papel ativo do usuário
	roleString := "Unknown"
	if domainUser.GetActiveRole() != nil {
		roleString = domainUser.GetActiveRole().GetRole().String()
	}

	// Converte datas para formato RFC3339 string
	var bornAtStr, creciValidityStr string

	if !domainUser.GetBornAt().IsZero() {
		bornAtStr = domainUser.GetBornAt().Format(time.RFC3339)
	}

	if !domainUser.GetCreciValidity().IsZero() {
		creciValidityStr = domainUser.GetCreciValidity().Format(time.RFC3339)
	}

	// Constrói o modelo HTTP com campos relevantes para o perfil do usuário
	return httpmodels.ProfileUser{
		ID:            domainUser.GetID(),
		FullName:      domainUser.GetFullName(),
		NickName:      domainUser.GetNickName(),
		NationalID:    domainUser.GetNationalID(),
		CreciNumber:   domainUser.GetCreciNumber(),
		CreciState:    domainUser.GetCreciState(),
		CreciValidity: creciValidityStr,
		BornAt:        bornAtStr,
		PhoneNumber:   domainUser.GetPhoneNumber(),
		Email:         domainUser.GetEmail(),
		ZipCode:       domainUser.GetZipCode(),
		Street:        domainUser.GetStreet(),
		Number:        domainUser.GetNumber(),
		Complement:    domainUser.GetComplement(),
		Neighborhood:  domainUser.GetNeighborhood(),
		City:          domainUser.GetCity(),
		State:         domainUser.GetState(),
		Role:          roleString,
		// Password, DeviceToken e LastActivity são excluídos intencionalmente
	}
}

// ConvertHTTPUserToDomain converte modelo HTTP de usuário para modelo de domínio
// Esta função centraliza a conversão entre as camadas HTTP e de domínio
// Parâmetros:
//   - httpUser: Modelo HTTP do usuário a ser convertido
//
// Retorna:
//   - usermodel.UserInterface: Modelo de domínio equivalente
func ConvertHTTPUserToDomain(httpUser httpmodels.User) usermodel.UserInterface {
	// Cria nova instância do modelo de domínio
	user := usermodel.NewUser()

	// Converte campos individuais, verificando se estão preenchidos
	// Isso evita sobrescrever valores padrão com valores vazios
	if httpUser.ID != 0 {
		user.SetID(httpUser.ID)
	}
	if httpUser.NickName != "" {
		user.SetNickName(httpUser.NickName)
	}
	if httpUser.Email != "" {
		user.SetEmail(httpUser.Email)
	}
	if httpUser.PhoneNumber != "" {
		user.SetPhoneNumber(httpUser.PhoneNumber)
	}
	if httpUser.NationalID != "" {
		user.SetNationalID(httpUser.NationalID)
	}
	if httpUser.BirthDate != "" {
		// Parse da data de nascimento do formato ISO 8601 (YYYY-MM-DD)
		if birthDate, err := time.Parse("2006-01-02", httpUser.BirthDate); err == nil {
			user.SetBornAt(birthDate)
		}
	}
	if httpUser.ZipCode != "" {
		user.SetZipCode(httpUser.ZipCode)
	}
	if httpUser.Password != "" {
		user.SetPassword(httpUser.Password)
	}

	return user
}

// ConvertDomainUserToHTTP converte modelo de domínio de usuário para modelo HTTP
// Esta função prepara os dados do domínio para exposição via API REST
// Parâmetros:
//   - domainUser: Modelo de domínio do usuário a ser convertido
//
// Retorna:
//   - httpmodels.User: Modelo HTTP equivalente para resposta da API
func ConvertDomainUserToHTTP(domainUser usermodel.UserInterface) httpmodels.User {
	// Determina o papel ativo do usuário
	// Se não houver papel ativo, usa "Unknown" como padrão
	roleString := "Unknown"
	if domainUser.GetActiveRole() != nil {
		roleString = domainUser.GetActiveRole().GetRole().String()
	}

	// Constrói o modelo HTTP com dados seguros para exposição
	// A senha nunca é incluída na resposta por segurança
	return httpmodels.User{
		ID:          domainUser.GetID(),
		NickName:    domainUser.GetNickName(),
		Email:       domainUser.GetEmail(),
		PhoneNumber: domainUser.GetPhoneNumber(),
		NationalID:  domainUser.GetNationalID(),
		// Senha nunca é incluída na resposta por motivos de segurança
		Role: roleString,
	}
}
