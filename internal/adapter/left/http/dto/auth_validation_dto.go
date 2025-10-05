package dto

// ValidateCPFRequest representa o payload de validação de CPF
// swagger:model ValidateCPFRequest
// @Description Assinatura da requisição deve ser calculada com METHOD|PATH|timestamp|payload_compacto_sem_campo_hmac
// @Example {"nationalID":"12345678901","bornAt":"1990-01-01","timestamp":1712345678,"hmac":"<signature>"}
type ValidateCPFRequest struct {
	NationalID string `json:"nationalID" binding:"required" example:"12345678901"`
	BornAt     string `json:"bornAt" binding:"required" example:"1990-01-01"`
	Timestamp  int64  `json:"timestamp" binding:"required" example:"1712345678"`
	HMAC       string `json:"hmac" binding:"required" example:"a1b2c3"`
}

// ValidateCNPJRequest representa o payload de validação de CNPJ
// @Example {"nationalID":"12345678000195","timestamp":1712345678,"hmac":"<signature>"}
type ValidateCNPJRequest struct {
	NationalID string `json:"nationalID" binding:"required" example:"12345678000195"`
	Timestamp  int64  `json:"timestamp" binding:"required" example:"1712345678"`
	HMAC       string `json:"hmac" binding:"required" example:"a1b2c3"`
}

// ValidateCEPRequest representa o payload de validação de CEP
// @Example {"postalCode":"12345678","timestamp":1712345678,"hmac":"<signature>"}
type ValidateCEPRequest struct {
	PostalCode string `json:"postalCode" binding:"required" example:"12345678"`
	Timestamp  int64  `json:"timestamp" binding:"required" example:"1712345678"`
	HMAC       string `json:"hmac" binding:"required" example:"a1b2c3"`
}

// ValidationResultResponse indica se a validação foi bem-sucedida
// swagger:model ValidationResultResponse
type ValidationResultResponse struct {
	Valid bool `json:"valid" example:"true"`
}

// CEPValidationResponse representa o retorno completo da busca de CEP
// swagger:model CEPValidationResponse
type CEPValidationResponse struct {
	Valid        bool   `json:"valid" example:"true"`
	PostalCode   string `json:"postalCode" example:"12345678"`
	Street       string `json:"street" example:"Av. Paulista"`
	Complement   string `json:"complement" example:"Apto 101"`
	Neighborhood string `json:"neighborhood" example:"Bela Vista"`
	City         string `json:"city" example:"São Paulo"`
	State        string `json:"state" example:"SP"`
}
