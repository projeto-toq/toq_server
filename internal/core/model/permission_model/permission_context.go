package permissionmodel

// PermissionContext contém o contexto para avaliação de permissões
type PermissionContext struct {
	UserID    int64                  `json:"user_id"`
	OwnerID   *int64                 `json:"owner_id,omitempty"`
	UserRoles []string               `json:"user_roles"`
	Metadata  map[string]interface{} `json:"metadata"`
}

// NewPermissionContext cria um novo contexto de permissão
func NewPermissionContext(userID int64) *PermissionContext {
	return &PermissionContext{
		UserID:    userID,
		UserRoles: []string{},
		Metadata:  make(map[string]interface{}),
	}
}

// WithOwner define o proprietário do recurso
func (pc *PermissionContext) WithOwner(ownerID int64) *PermissionContext {
	pc.OwnerID = &ownerID
	return pc
}

// WithRoles define os roles do usuário
func (pc *PermissionContext) WithRoles(roles []string) *PermissionContext {
	pc.UserRoles = roles
	return pc
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
