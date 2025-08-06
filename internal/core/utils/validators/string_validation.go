package validators

import (
	"regexp"

	"github.com/giulio-alfieri/toq_server/internal/core/utils/converters"
	"github.com/nyaruka/phonenumbers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ValidateNoSpecialCharacters(str string) (err error) {
	// validate if the string has only letters (including accented characters) and spaces
	re := regexp.MustCompile(`^[a-zA-ZÀ-ÿ\s'-]+$`)
	if !re.MatchString(str) {
		return status.Error(codes.InvalidArgument, str+" contains special characters")
	}
	return nil
}

func ValidateOnlyNumbers(str string) error {
	// validate if the string has only numbers
	re := regexp.MustCompile(`^[0-9]+$`)
	if !re.MatchString(str) {
		return status.Error(codes.InvalidArgument, str+" contains non-numeric characters")
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
		return status.Error(codes.InvalidArgument, "Invalid phone number")
	}

	// Check if the phone number is valid and in E.164 format
	if !phonenumbers.IsValidNumber(num) || phonenumbers.Format(num, phonenumbers.E164) != phoneNumber {
		return status.Error(codes.InvalidArgument, "Invalid phone number")
	}

	return nil
}

func ValidateEmail(email string) error {
	// validate if the email is in a valid format
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !re.MatchString(email) {
		return status.Error(codes.InvalidArgument, "Invalid email format")
	}
	return nil
}

func ValidateCode(code string) error {
	// validate if the code has only letters and numbers and is 6 characters long
	re := regexp.MustCompile(`^[a-zA-Z0-9]{6}$`)
	if !re.MatchString(code) {
		return status.Error(codes.InvalidArgument, "Invalid code format")
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
func ValidateCreciEquality(creciNumber1 string, creciNumber2 string) (isEqual bool) {
	creciNumber1 = converters.NormalizeAndTrimString(creciNumber1)
	creciNumber2 = converters.NormalizeAndTrimString(creciNumber2)
	return creciNumber1 == creciNumber2
}
