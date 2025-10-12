package utils

import (
	"bytes"
	"encoding/json"
	"strings"
)

// Optional representa um campo opcional em payloads, preservando se o valor foi
// enviado e se ele estava explicitamente como null.
type Optional[T any] struct {
	value   T
	present bool
	null    bool
}

// UnmarshalJSON marca o campo como presente e, quando aplicável, sinaliza que o
// valor recebido é nulo. Para qualquer outro conteúdo, delega a decodificação
// padrão do tipo.
func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	o.present = true
	trimmed := strings.TrimSpace(string(data))
	if len(trimmed) == 0 || trimmed == "null" {
		o.null = true
		var zero T
		o.value = zero
		return nil
	}

	var v T
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&v); err != nil {
		// Se a decodificação falhar, preservamos o estado para permitir tratamento a jusante
		return err
	}

	o.value = v
	o.null = false
	return nil
}

// IsPresent indica se o campo estava presente no payload.
func (o Optional[T]) IsPresent() bool {
	return o.present
}

// IsNull indica se o campo presente foi enviado explicitamente como null.
func (o Optional[T]) IsNull() bool {
	return o.present && o.null
}

// Value retorna o valor decodificado juntamente com um indicador de validade
// (valor presente e não nulo).
func (o Optional[T]) Value() (T, bool) {
	return o.value, o.present && !o.null
}

// ValueOrZero retorna o valor armazenado, mesmo quando null. Útil para campos
// onde o zero-value representa limpeza.
func (o Optional[T]) ValueOrZero() T {
	return o.value
}

// NewOptionalValue cria um Optional marcado como presente com o valor informado.
func NewOptionalValue[T any](value T) Optional[T] {
	return Optional[T]{
		value:   value,
		present: true,
		null:    false,
	}
}

// NewOptionalNull cria um Optional presente sinalizando null explícito.
func NewOptionalNull[T any]() Optional[T] {
	return Optional[T]{
		present: true,
		null:    true,
	}
}
