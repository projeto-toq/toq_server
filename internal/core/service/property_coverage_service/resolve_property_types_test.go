package propertycoverageservice

import (
	"testing"

	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
)

func TestSanitizeCoverageNumber(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "empty", input: "", expected: ""},
		{name: "trim spaces", input: " 123 ", expected: "123"},
		{name: "remove inner spaces", input: "149; 189", expected: "149;189"},
		{name: "upper case", input: "apto 12", expected: "APTO12"},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := sanitizeCoverageNumber(tt.input)
			if got != tt.expected {
				t.Fatalf("sanitizeCoverageNumber(%q) = %q, expected %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestNormalizeCoverageInput(t *testing.T) {
	t.Parallel()

	input := propertycoveragemodel.ResolvePropertyTypesInput{
		ZipCode: "06472001",
		Number:  " 149; 189 ",
	}

	zip, number, err := normalizeCoverageInput(input)
	if err != nil {
		t.Fatalf("normalizeCoverageInput unexpected error: %v", err)
	}

	if zip != "06472001" {
		t.Fatalf("expected normalized zip 06472001, got %s", zip)
	}

	if number != "149;189" {
		t.Fatalf("expected sanitized number 149;189, got %s", number)
	}

	_, _, err = normalizeCoverageInput(propertycoveragemodel.ResolvePropertyTypesInput{ZipCode: "123"})
	if err == nil {
		t.Fatalf("expected error for invalid zip code")
	}
}
