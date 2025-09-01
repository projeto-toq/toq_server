package usermodel

import "net/http"

// SigninErrorType representa os diferentes tipos de erro no processo de signin
type SigninErrorType int

const (
	SigninErrorInvalidCredentials SigninErrorType = iota
	SigninErrorUserBlocked
	SigninErrorNoActiveRoles
	SigninErrorInternalError
	SigninErrorInvalidRequest
)

// SigninError representa um erro específico do processo de signin
type SigninError struct {
	Type       SigninErrorType
	Message    string
	StatusCode int
	ShouldLog  bool
	UserID     *int64 // Para contexto de logging
}

// NewSigninError cria um novo erro de signin
func NewSigninError(errorType SigninErrorType, message string, userID *int64) *SigninError {
	se := &SigninError{
		Type:      errorType,
		Message:   message,
		UserID:    userID,
		ShouldLog: true,
	}

	switch errorType {
	case SigninErrorInvalidCredentials:
		se.StatusCode = http.StatusUnauthorized
	case SigninErrorUserBlocked:
		se.StatusCode = http.StatusLocked // 423
	case SigninErrorNoActiveRoles:
		se.StatusCode = http.StatusForbidden
	case SigninErrorInternalError:
		se.StatusCode = http.StatusInternalServerError
	case SigninErrorInvalidRequest:
		se.StatusCode = http.StatusBadRequest
		se.ShouldLog = false // Request inválido não precisa de log detalhado
	default:
		se.StatusCode = http.StatusInternalServerError
	}

	return se
}

// Error implementa a interface error
func (se *SigninError) Error() string {
	return se.Message
}

// GetErrorType retorna o tipo do erro
func (se *SigninError) GetErrorType() SigninErrorType {
	return se.Type
}

// GetMessage retorna a mensagem do erro
func (se *SigninError) GetMessage() string {
	return se.Message
}

// GetStatusCode retorna o código de status HTTP
func (se *SigninError) GetStatusCode() int {
	return se.StatusCode
}

// ShouldLogError indica se o erro deve ser logado
func (se *SigninError) ShouldLogError() bool {
	return se.ShouldLog
}

// GetUserID retorna o ID do usuário relacionado ao erro (se disponível)
func (se *SigninError) GetUserID() *int64 {
	return se.UserID
}

// IsUserBlockedError verifica se é um erro de usuário bloqueado
func (se *SigninError) IsUserBlockedError() bool {
	return se.Type == SigninErrorUserBlocked
}

// IsCredentialsError verifica se é um erro de credenciais inválidas
func (se *SigninError) IsCredentialsError() bool {
	return se.Type == SigninErrorInvalidCredentials
}

// IsInternalError verifica se é um erro interno
func (se *SigninError) IsInternalError() bool {
	return se.Type == SigninErrorInternalError
}
