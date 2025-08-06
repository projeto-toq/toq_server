package converters

import (
	"database/sql"
	"strings"
	"time"
)

// NormalizeAndTrimString takes a string, removes any leading and trailing spaces,
// and removes any leading zeroes from the string. It returns the normalized string.
//
// Parameters:
//
//	s - the input string to be normalized and trimmed.
//
// Returns:
//
//	A string with leading and trailing spaces removed and leading zeroes removed.
func NormalizeAndTrimString(s string) string {
	return RemoveInitialZeroes(RemoveSpaces(s))
}

// RemoveInitialZeroes removes all leading '0' characters from the input string.
// It iterates through the string and returns a substring starting from the first
// non-'0' character. If the string consists entirely of '0' characters, it returns
// an empty string.
//
// Parameters:
//
//	s - the input string from which leading '0' characters are to be removed.
//
// Returns:
//
//	A substring of the input string starting from the first non-'0' character,
//	or an empty string if the input string contains only '0' characters.
func RemoveInitialZeroes(s string) string {
	for i, c := range s {
		if c != '0' {
			return s[i:]
		}
	}
	return ""
}

// RemoveSpaces removes all spaces from the input string and returns the resulting string.
// It uses a strings.Builder to efficiently build the output string.
//
// Parameters:
//   - s: The input string from which spaces will be removed.
//
// Returns:
//
//	A new string with all spaces removed from the input string.
func RemoveSpaces(s string) string {
	var result strings.Builder
	result.Grow(len(s))
	for _, char := range s {
		if char != ' ' {
			result.WriteRune(char)
		}
	}
	return result.String()
}

// StringToNullString converts a string to a sql.NullString.
// If the input string is empty, it returns a sql.NullString with Valid set to false.
// Otherwise, it returns a sql.NullString with the input string and Valid set to true.
//
// Parameters:
//   - value: The input string to be converted.
//
// Returns:
//
//	A sql.NullString representing the input string.
func StringToNullString(value string) sql.NullString {
	if value == "" {
		return sql.NullString{
			String: "",
			Valid:  false,
		}
	}
	return sql.NullString{
		String: value,
		Valid:  true,
	}
}

func BytesToNullString(value []byte) sql.NullString {
	if len(value) == 0 {
		return sql.NullString{
			String: "",
			Valid:  false,
		}
	}
	return sql.NullString{
		String: string(value),
		Valid:  true,
	}
}

// RemoveAllButDigits removes all characters from the input string except digits (0-9).
//
// Parameters:
//   - s: The input string from which non-digit characters will be removed.
//
// Returns:
//
//	A new string containing only the digit characters from the input string.
func RemoveAllButDigits(s string) string {
	var result strings.Builder
	result.Grow(len(s))
	for _, char := range s {
		if char >= '0' && char <= '9' {
			result.WriteRune(char)
		}
	}
	return result.String()
}

// RemoveAllButDigitsAndPlusSign removes all characters from the input string
// except for digits (0-9) and the plus sign (+).
//
// Parameters:
//
//	s - The input string to be processed.
//
// Returns:
//
//	A new string containing only digits and plus signs from the input string.
func RemoveAllButDigitsAndPlusSign(s string) string {

	var result strings.Builder
	result.Grow(len(s))
	for _, char := range s {
		if (char >= '0' && char <= '9') || char == '+' {
			result.WriteRune(char)
		}
	}
	return result.String()
}

// TrimSpaces removes all leading and trailing spaces from the input string.
//
// Parameters:
//   - s: The input string from which leading and trailing spaces will be removed.
//
// Returns:
//
//	A new string with all leading and trailing spaces removed from the input string.
func TrimSpaces(s string) string {
	return strings.TrimSpace(s)
}

func StrngToTime(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}
	return time.Parse("2006-01-02", s)
}
