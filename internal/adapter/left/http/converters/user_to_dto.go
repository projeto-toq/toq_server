package converters

import (
	"time"

	dto "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
)

// ToGetProfileResponse converts a domain user to the public GetProfileResponse DTO (camelCase).
// Preenche campos opcionais de forma segura e formata datas em YYYY-MM-DD quando aplicável.
func ToGetProfileResponse(u usermodel.UserInterface) dto.GetProfileResponse {
	user := dto.UserResponse{
		ID:           u.GetID(),
		FullName:     u.GetFullName(),
		NickName:     u.GetNickName(),
		NationalID:   u.GetNationalID(),
		CreciNumber:  u.GetCreciNumber(),
		CreciState:   u.GetCreciState(),
		BornAt:       parseDate(u.GetBornAt()),
		PhoneNumber:  u.GetPhoneNumber(),
		Email:        u.GetEmail(),
		ZipCode:      u.GetZipCode(),
		Street:       u.GetStreet(),
		Number:       u.GetNumber(),
		Complement:   u.GetComplement(),
		Neighborhood: u.GetNeighborhood(),
		City:         u.GetCity(),
		State:        u.GetState(),
		LastActivity: formatDateTime(u.GetLastActivityAt()),
	}

	// Active role (nil-safe)
	if ar := u.GetActiveRole(); ar != nil {
		// Preferir slug se disponível
		roleSlug := ""
		if r := ar.GetRole(); r != nil {
			roleSlug = r.GetSlug()
		}
		user.ActiveRole = dto.UserRoleResponse{
			ID:         ar.GetID(),
			UserID:     ar.GetUserID(),
			BaseRoleID: ar.GetRoleID(),
			Role:       roleSlug,
			Active:     ar.GetIsActive(),
			Status:     ar.GetStatus().String(),
		}
	}

	return dto.GetProfileResponse{User: user}
}

// parseDate safely formats a date as YYYY-MM-DD or returns empty string if zero.
func parseDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
}

// formatDateTime returns RFC3339 timestamp or empty string if zero.
func formatDateTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}
