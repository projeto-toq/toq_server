package userservices

import (
	"time"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
)

const systemUserTemplateID int64 = 1

// ListUsersInput encapsula filtros aceitos para listagem de usuários via painel admin.
type ListUsersInput struct {
	Page             int
	Limit            int
	RoleName         string
	RoleSlug         string
	RoleStatus       *globalmodel.UserRoleStatus
	IsSystemRole     *bool
	FullName         string
	CPF              string
	Email            string
	PhoneNumber      string
	Deleted          *bool
	IDFrom           *int64
	IDTo             *int64
	BornAtFrom       *time.Time
	BornAtTo         *time.Time
	LastActivityFrom *time.Time
	LastActivityTo   *time.Time
}

// ListUsersOutput descreve o resultado da listagem de usuários admin.
type ListUsersOutput struct {
	Users []usermodel.UserInterface
	Total int64
	Page  int
	Limit int
}

// ListPendingRealtorsOutput aggregates pending realtor data with pagination metadata.
type ListPendingRealtorsOutput struct {
	Realtors []usermodel.UserInterface
	Total    int64
	Page     int
	Limit    int
}

// CreateSystemUserInput representa o payload necessário para criar um usuário de sistema.
type CreateSystemUserInput struct {
	NickName    string
	Email       string
	PhoneNumber string
	CPF         string
	BornAt      time.Time
	RoleSlug    permissionmodel.RoleSlug
	ZipCode     string
	Number      string
}

// UpdateSystemUserInput representa os dados editáveis de um usuário de sistema.
type UpdateSystemUserInput struct {
	UserID      int64
	FullName    string
	Email       string
	PhoneNumber string
}

// DeleteSystemUserInput representa o alvo de exclusão lógica de um usuário de sistema.
type DeleteSystemUserInput struct {
	UserID int64
}

// SystemUserResult agrega os dados principais retornados após criar/atualizar um usuário de sistema.
type SystemUserResult struct {
	UserID   int64
	RoleID   int64
	RoleSlug permissionmodel.RoleSlug
	Email    string
}
