package validators

import (
	"regexp"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"github.com/nyaruka/phonenumbers"
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

// ValidateCreciEquality compares two CRECI numbers for equality after normalization using TrimSpaces and RemoveInitialZeroes.
//
// Parameters:
// - creciNumber1: The first CRECI number as a string.
// - creciNumber2: The second CRECI number as a string.
//
// Returns:
// - isEqual: A boolean indicating whether the two CRECI numbers are equal after normalization.
// func ValidateCreciEquality(creciNumber1 string, creciNumber2 string) (isEqual bool) {
// 	creciNumber1 = converters.NormalizeAndTrimString(creciNumber1)
// 	creciNumber2 = converters.NormalizeAndTrimString(creciNumber2)
// 	return creciNumber1 == creciNumber2
// }
