package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	coreutils "github.com/giulio-alfieri/toq_server/internal/core/utils"
	validator "github.com/go-playground/validator/v10"
)

// MapBindingError converts binding/validation errors into a DomainError with field-level details.
// Semantics:
// - JSON syntax/type errors -> 400 Invalid request format
// - Validation errors (required, len, min, max, email, etc.) -> 422 Validation failed with details
func MapBindingError(err error) coreutils.DomainError {
	if err == nil {
		return nil
	}

	// JSON decoding errors (syntax/type)
	var synErr *json.SyntaxError
	if ok := errors.As(err, &synErr); ok {
		return coreutils.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}
	var typeErr *json.UnmarshalTypeError
	if ok := errors.As(err, &typeErr); ok {
		return coreutils.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	// Validator field errors → 422 with details
	if verrs, ok := err.(validator.ValidationErrors); ok {
		details := make([]map[string]string, 0, len(verrs))
		for _, fe := range verrs {
			fieldPath := toJSONPathFromStructNamespace(fe.StructNamespace())
			message := messageForTag(fe)
			details = append(details, map[string]string{
				"field":   fieldPath,
				"message": message,
			})
		}
		return coreutils.NewHTTPError(http.StatusUnprocessableEntity, "Validation failed", map[string]any{
			"errors": details,
		})
	}

	// Fallback: treat as invalid request format
	// Some binding errors are generic (e.g., EOF, invalid character). Keep 400.
	return coreutils.NewHTTPError(http.StatusBadRequest, "Invalid request format")
}

// messageForTag maps validator tags to human-readable messages.
func messageForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "email":
		return "invalid email"
	case "len":
		return "must have length = " + fe.Param()
	case "min":
		return "must be at least " + fe.Param() + " characters"
	case "max":
		return "must be at most " + fe.Param() + " characters"
	default:
		return "invalid value"
	}
}

// toJSONPathFromStructNamespace builds a dotted json path from a struct namespace.
// Example: "CreateOwnerRequest.Owner.NickName" -> "owner.nickName".
func toJSONPathFromStructNamespace(ns string) string {
	if ns == "" {
		return ""
	}
	parts := strings.Split(ns, ".")
	if len(parts) <= 1 {
		return lowerFirst(ns)
	}
	// drop root struct name
	parts = parts[1:]
	for i := range parts {
		parts[i] = lowerFirst(parts[i])
	}
	return strings.Join(parts, ".")
}

func lowerFirst(s string) string {
	if s == "" {
		return s
	}
	// ASCII-only lowercasing of first rune is sufficient for JSON tag names here
	b := []byte(s)
	if b[0] >= 'A' && b[0] <= 'Z' {
		b[0] = b[0] - 'A' + 'a'
	}
	return string(b)
}
