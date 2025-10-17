package validators

import (
	"errors"
	"regexp"
	"strings"

	"github.com/nyaruka/phonenumbers"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func ValidateNoSpecialCharacters(str string) error {
	// validate if the string has only letters (including accented characters) and spaces
	re := regexp.MustCompile(`^[a-zA-ZÀ-ÿ\s'-]+$`)
	if !re.MatchString(str) {
		return utils.ValidationError("text", "contains special characters")
	}
	return nil
}

func ValidateOnlyNumbers(str string) error {
	// validate if the string has only numbers
	re := regexp.MustCompile(`^[0-9]+$`)
	if !re.MatchString(str) {
		return utils.ValidationError("number", "must contain only numeric characters")
	}
	return nil
}

// OnlyDigits returns a string containing only ASCII digits [0-9].
// It preserves leading zeros and removes any non-digit characters.
func OnlyDigits(s string) string {
	// Fast path: if already matches only digits, return as-is
	reFull := regexp.MustCompile(`^[0-9]+$`)
	if reFull.MatchString(s) {
		return s
	}
	// Remove everything that is not a digit
	reStrip := regexp.MustCompile(`[^0-9]`)
	return reStrip.ReplaceAllString(s, "")
}

// ValidateE164 validates if the given phone number is in E.164 format.
// It parses the phone number and checks if it is valid and formatted correctly.
//
// Parameters:
//
//	phoneNumber (string): The phone number to validate.
//
// Returns:
//
//	error: Returns an error if the phone number is invalid or not in E.164 format.
func ValidateE164(phoneNumber string) error {
	// Parse the phone number
	num, err := phonenumbers.Parse(phoneNumber, "")
	if err != nil {
		return utils.ValidationError("phone_number", "invalid phone number format")
	}

	// Check if the phone number is valid and in E.164 format
	if !phonenumbers.IsValidNumber(num) || phonenumbers.Format(num, phonenumbers.E164) != phoneNumber {
		return utils.ValidationError("phone_number", "must be in valid E.164 format")
	}

	return nil
}

// NormalizeToE164 parses and formats a phone number into E.164. Returns normalized value or
// a validation error if it's not a valid number.
func NormalizeToE164(phoneNumber string) (string, error) {
	num, err := phonenumbers.Parse(phoneNumber, "")
	if err != nil {
		return "", utils.ValidationError("phone_number", "invalid phone number format")
	}
	if !phonenumbers.IsValidNumber(num) {
		return "", utils.ValidationError("phone_number", "invalid phone number")
	}
	return phonenumbers.Format(num, phonenumbers.E164), nil
}

func ValidateEmail(email string) error {
	// validate if the email is in a valid format
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !re.MatchString(email) {
		return utils.ValidationError("email", "invalid email format")
	}
	return nil
}

func ValidateCode(code string) error {
	// validate if the code has only letters and numbers and is 6 characters long
	re := regexp.MustCompile(`^[a-zA-Z0-9]{6}$`)
	if !re.MatchString(code) {
		return utils.ValidationError("code", "must be 6 alphanumeric characters")
	}
	return nil
}

var (
	cepStrictPattern   = regexp.MustCompile(`^[0-9]{8}$`)
	creciNumberPattern = regexp.MustCompile(`^[0-9]+-F$`)
	brazilianUFs       = map[string]struct{}{
		"AC": {}, "AL": {}, "AP": {}, "AM": {}, "BA": {}, "CE": {}, "DF": {}, "ES": {}, "GO": {},
		"MA": {}, "MT": {}, "MS": {}, "MG": {}, "PA": {}, "PB": {}, "PR": {}, "PE": {}, "PI": {},
		"RJ": {}, "RN": {}, "RS": {}, "RO": {}, "RR": {}, "SC": {}, "SP": {}, "SE": {}, "TO": {},
	}
)

var ErrInvalidCEPFormat = errors.New("invalid CEP format")

// NormalizeCEP validates and enforces CEP format (eight digits) returning the trimmed value.
// Returns ErrInvalidCEPFormat when the input does not comply.
func NormalizeCEP(value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if !cepStrictPattern.MatchString(trimmed) {
		return "", ErrInvalidCEPFormat
	}
	return trimmed, nil
}

// ValidateCreciNumber normaliza e valida o formato do número CRECI.
// Quando required for true, o valor é obrigatório; caso contrário, retorna string vazia se ausente.
func ValidateCreciNumber(field, value string, required bool) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		if required {
			return "", utils.ValidationError(field, "Creci number is required")
		}
		return "", nil
	}

	if !creciNumberPattern.MatchString(trimmed) {
		return "", utils.ValidationError(field, "Creci number must be numeric and end with -F")
	}

	return trimmed, nil
}

// ValidateCreciState normaliza e valida o estado CRECI (UF brasileira).
// Quando required for true, o valor é obrigatório; caso contrário, retorna string vazia se ausente.
func ValidateCreciState(field, value string, required bool) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		if required {
			return "", utils.ValidationError(field, "Creci state is required")
		}
		return "", nil
	}

	upper := strings.ToUpper(trimmed)
	if _, ok := brazilianUFs[upper]; !ok {
		return "", utils.ValidationError(field, "Creci state must be a valid Brazilian UF")
	}

	return upper, nil
}
