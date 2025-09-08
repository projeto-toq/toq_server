package usermodel

import permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"

type Tokens struct {
	//the token to be used on request
	AccessToken string `json:"access_token"`
	//the token to refresh the access token whe it expires
	RefreshToken string `json:"refresh_token"`
}

type UserInfos struct {
	ID         int64                          // ID do usuário
	UserRoleID int64                          // ID do UserRole ativo
	RoleStatus permissionmodel.UserRoleStatus // Status detalhado da role
	RoleSlug   permissionmodel.RoleSlug       // Slug textual da role ativa (aditivo, backward-compatible)
}

type JWT struct {
	Secret string
}
