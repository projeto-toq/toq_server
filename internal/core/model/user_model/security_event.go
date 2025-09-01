package usermodel

import "time"

// SecurityEventType representa os tipos de eventos de segurança
type SecurityEventType string

const (
	SecurityEventSigninAttempt      SecurityEventType = "signin_attempt"
	SecurityEventSigninSuccess      SecurityEventType = "signin_success"
	SecurityEventSigninFailure      SecurityEventType = "signin_failure"
	SecurityEventUserBlocked        SecurityEventType = "user_blocked"
	SecurityEventUserUnblocked      SecurityEventType = "user_unblocked"
	SecurityEventInvalidCredentials SecurityEventType = "invalid_credentials"
	SecurityEventNoActiveRoles      SecurityEventType = "no_active_roles"
)

// SecurityEventResult representa o resultado de um evento de segurança
type SecurityEventResult string

const (
	SecurityEventResultSuccess SecurityEventResult = "success"
	SecurityEventResultFailure SecurityEventResult = "failure"
	SecurityEventResultBlocked SecurityEventResult = "blocked"
)

// SecurityEvent representa um evento de segurança no sistema
type SecurityEvent struct {
	UserID     *int64                 `json:"userId,omitempty"`
	NationalID string                 `json:"nationalId,omitempty"`
	EventType  SecurityEventType      `json:"eventType"`
	Result     SecurityEventResult    `json:"result"`
	IPAddress  string                 `json:"ipAddress,omitempty"`
	UserAgent  string                 `json:"userAgent,omitempty"`
	Timestamp  time.Time              `json:"timestamp"`
	Details    map[string]interface{} `json:"details,omitempty"`
	ErrorType  *SigninErrorType       `json:"errorType,omitempty"`
	Reason     string                 `json:"reason,omitempty"`
}

// NewSecurityEvent cria um novo evento de segurança
func NewSecurityEvent(eventType SecurityEventType, result SecurityEventResult) *SecurityEvent {
	return &SecurityEvent{
		EventType: eventType,
		Result:    result,
		Timestamp: time.Now().UTC(),
		Details:   make(map[string]interface{}),
	}
}

// WithUserID adiciona o ID do usuário ao evento
func (se *SecurityEvent) WithUserID(userID int64) *SecurityEvent {
	se.UserID = &userID
	return se
}

// WithNationalID adiciona o CPF/CNPJ ao evento
func (se *SecurityEvent) WithNationalID(nationalID string) *SecurityEvent {
	se.NationalID = nationalID
	return se
}

// WithIPAddress adiciona o endereço IP ao evento
func (se *SecurityEvent) WithIPAddress(ip string) *SecurityEvent {
	se.IPAddress = ip
	return se
}

// WithUserAgent adiciona o user agent ao evento
func (se *SecurityEvent) WithUserAgent(userAgent string) *SecurityEvent {
	se.UserAgent = userAgent
	return se
}

// WithErrorType adiciona o tipo de erro ao evento
func (se *SecurityEvent) WithErrorType(errorType SigninErrorType) *SecurityEvent {
	se.ErrorType = &errorType
	return se
}

// WithReason adiciona uma razão ao evento
func (se *SecurityEvent) WithReason(reason string) *SecurityEvent {
	se.Reason = reason
	return se
}

// WithDetail adiciona um detalhe específico ao evento
func (se *SecurityEvent) WithDetail(key string, value interface{}) *SecurityEvent {
	se.Details[key] = value
	return se
}

// GetUserContext retorna informações de contexto do usuário
func (se *SecurityEvent) GetUserContext() map[string]interface{} {
	context := make(map[string]interface{})

	if se.UserID != nil {
		context["userID"] = *se.UserID
	}

	if se.NationalID != "" {
		context["nationalID"] = se.NationalID
	}

	if se.IPAddress != "" {
		context["ipAddress"] = se.IPAddress
	}

	if se.UserAgent != "" {
		context["userAgent"] = se.UserAgent
	}

	return context
}

// IsFailureEvent verifica se é um evento de falha
func (se *SecurityEvent) IsFailureEvent() bool {
	return se.Result == SecurityEventResultFailure
}

// IsBlockingEvent verifica se é um evento de bloqueio
func (se *SecurityEvent) IsBlockingEvent() bool {
	return se.EventType == SecurityEventUserBlocked || se.Result == SecurityEventResultBlocked
}

// IsSigninEvent verifica se é um evento relacionado a signin
func (se *SecurityEvent) IsSigninEvent() bool {
	return se.EventType == SecurityEventSigninAttempt ||
		se.EventType == SecurityEventSigninSuccess ||
		se.EventType == SecurityEventSigninFailure
}
