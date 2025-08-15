package crecimodel

import "time"

type creci struct {
	creciNumber   string
	creciState    string
	creciValidity time.Time
}

func (c *creci) GetCreciNumber() string {
	return c.creciNumber
}

func (c *creci) SetCreciNumber(creciNumber string) {
	c.creciNumber = creciNumber
}

func (c *creci) GetCreciState() string {
	return c.creciState
}

func (c *creci) SetCreciState(creciState string) {
	c.creciState = creciState
}

func (c *creci) GetCreciValidity() time.Time {
	return c.creciValidity
}

func (c *creci) SetCreciValidity(creciValidity time.Time) {
	c.creciValidity = creciValidity
}

// list of implemented creci states (for now only SP is implemented
var ImplementedCreciStates = []string{"SP"}

var ValidStates = map[string]bool{
	"AC": true, "AL": true, "AP": true, "AM": true, "BA": true, "CE": true, "DF": true, "ES": true, "GO": true,
	"MA": true, "MT": true, "MS": true, "MG": true, "PA": true, "PB": true, "PR": true, "PE": true, "PI": true,
	"RJ": true, "RN": true, "RS": true, "RO": true, "RR": true, "SC": true, "SP": true, "SE": true, "TO": true,
}

const (
	// FaceMatchThreshold define o limiar mínimo de similaridade de cosseno (0..1) entre dois vetores faciais
	// 0.85 é um ponto inicial intermediário; poderá ser ajustado empiricamente em produção.
	FaceMatchThreshold = float32(0.85)
)
