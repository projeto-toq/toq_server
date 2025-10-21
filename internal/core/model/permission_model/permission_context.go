package permissionmodel

// PermissionContext contém o contexto para avaliação de permissões
type PermissionContext struct {
	UserID     int64                  `json:"user_id"`
	UserRoleID int64                  `json:"user_role_id"`
	RoleStatus UserRoleStatus         `json:"role_status"`
	RoleSlug   RoleSlug               `json:"role_slug"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// NewPermissionContext cria um novo contexto de permissão
func NewPermissionContext(userID, userRoleID int64) *PermissionContext {
	return &PermissionContext{
		UserID:     userID,
		UserRoleID: userRoleID,
		RoleStatus: StatusActive, // Mantém por compatibilidade, mas será ignorado
		RoleSlug:   "",
		Metadata:   make(map[string]interface{}),
	}
}

// AddMetadata adiciona metadados ao contexto
func (pc *PermissionContext) AddMetadata(key string, value interface{}) *PermissionContext {
	pc.Metadata[key] = value
	return pc
}

// WithMetadata define todos os metadados
func (pc *PermissionContext) WithMetadata(metadata map[string]interface{}) *PermissionContext {
	pc.Metadata = metadata
	return pc
}

// SetRoleSlug define o role slug associado ao contexto
func (pc *PermissionContext) SetRoleSlug(roleSlug RoleSlug) *PermissionContext {
	pc.RoleSlug = roleSlug
	return pc
}

// GetRoleSlug retorna o role slug associado ao contexto
func (pc *PermissionContext) GetRoleSlug() RoleSlug {
	return pc.RoleSlug
}

// HasRole verifica se o usuário possui um role específico
func (pc *PermissionContext) HasRole(roleSlug RoleSlug) bool {
	// Implementar lógica baseada no UserRoleID
	// Esta função substituirá a verificação por strings
	return pc.UserRoleID > 0
}

// // IsActive verifica se o role do usuário está ativo
// func (pc *PermissionContext) IsActive() bool {
// 	return pc.UserRoleID > 0
// }
