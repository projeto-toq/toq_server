package converters

import (
	"time"

	dto "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
)

// ToUserProfileResponse converts a domain user to a stable API response DTO.
// Campos opcionais são preenchidos de forma segura para evitar panics (ex.: active role nulo).
func ToUserProfileResponse(u usermodel.UserInterface) dto.UserProfileResponse {
	data := dto.UserProfileData{
		ID:          u.GetID(),
		Email:       u.GetEmail(),
		PhoneNumber: u.GetPhoneNumber(),
		FullName:    u.GetFullName(),
		NickName:    u.GetNickName(),
		NationalID:  u.GetNationalID(),
		ZipCode:     u.GetZipCode(),
		Street:      u.GetStreet(),
		City:        u.GetCity(),
		State:       u.GetState(),
	}

	// born_at format YYYY-MM-DD (omit if zero)
	if t := u.GetBornAt(); !t.IsZero() {
		data.BornAt = t.Format("2006-01-02")
	}

	// Active role (nil-safe)
	if ar := u.GetActiveRole(); ar != nil {
		data.ActiveRole = toActiveRoleDTO(ar)
	}

	return dto.UserProfileResponse{Data: data}
}

func toActiveRoleDTO(ar permissionmodel.UserRoleInterface) *dto.ActiveRoleDTO {
	out := &dto.ActiveRoleDTO{
		ID:     ar.GetID(),
		Active: ar.GetIsActive(),
	}
	if r := ar.GetRole(); r != nil {
		out.Role = &dto.RoleDTO{
			ID:           r.GetID(),
			Name:         r.GetName(),
			Slug:         r.GetSlug(),
			Description:  r.GetDescription(),
			IsSystemRole: r.GetIsSystemRole(),
			IsActive:     r.GetIsActive(),
		}
	}
	return out
}

// parseDate is kept here for potential future needs (not used now),
// demonstrating how to manage time safely if inputs change in the future.
func parseDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
}
