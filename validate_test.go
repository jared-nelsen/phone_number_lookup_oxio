package main

import (
	"testing"
)

func TestValidateCountryCodeMeetsISO_3166_1_alpha_2(t *testing.T) {
	tests := []struct {
		name        string
		countryCode string
		expected    bool
	}{
		{"Valid US", "US", true},
		{"Valid MX", "MX", true},
		{"Valid ES", "ES", true},
		{"Valid lowercase us", "us", true},
		{"Invalid 3 chars", "ESP", false},
		{"Invalid 1 char", "E", false},
		{"Invalid empty", "", false},
		{"Invalid numeric", "12", false},
		{"Invalid special chars", "U$", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validateCountryCodeMeetsISO_3166_1_alpha_2(tt.countryCode)
			if result != tt.expected {
				t.Errorf("validateCountryCodeMeetsISO_3166_1_alpha_2(%s) = %v, expected %v", tt.countryCode, result, tt.expected)
			}
		})
	}
}

func TestValidatePhoneNumberFormat(t *testing.T) {
	tests := []struct {
		name        string
		phoneNumber string
		expected    bool
	}{
		{"Valid with +", "+12125690123", true},
		{"Valid without +", "12125690123", true},
		{"Valid with spaces", "1 212 569 0123", true},
		{"Valid with + and spaces", "+1 212 569 0123", true},
		{"Invalid with letters", "abc123", false},
		{"Invalid with special chars", "123-456-7890", false},
		{"Invalid with dots", "123.456.7890", false},
		{"Invalid mixed valid/invalid", "123abc456", false},
		{"Valid only digits", "1234567890", true},
		{"Valid only spaces and digits", "123 456 7890", true},
		{"Invalid empty", "", false},
		{"Invalid only +", "+", false},
		{"Invalid only spaces", "   ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validatePhoneNumberFormat(tt.phoneNumber)
			if result != tt.expected {
				t.Errorf("validatePhoneNumberFormat(%s) = %v, expected %v", tt.phoneNumber, result, tt.expected)
			}
		})
	}
}

func TestValidateSpaces(t *testing.T) {
	tests := []struct {
		name        string
		phoneNumber string
		expected    bool
	}{
		{"Valid no spaces", "12125690123", true},
		{"Valid with +", "+12125690123", true},
		{"Valid single space", "1 2125690123", true},
		{"Valid two spaces", "1 212 5690123", true},
		{"Valid three parts", "1 212 569", true},
		{"Invalid four parts", "1 212 569 0123", false},
		{"Invalid consecutive spaces", "1  212 5690123", false},
		{"Invalid leading space", " 1 212 5690123", false},
		{"Invalid trailing space", "1 212 5690123 ", false},
		{"Invalid empty part", "1  212", false},
		{"Invalid non-digit in part", "1 abc 212", false},
		{"Valid with + and spaces", "+1 212 5690123", true},
		{"Invalid with + and four parts", "+1 212 569 0123", false},
		{"Valid edge case - single digit parts", "1 2 3", true},
		{"Valid edge case - country code only", "1", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validateSpaces(tt.phoneNumber)
			if result != tt.expected {
				t.Errorf("validateSpaces(%s) = %v, expected %v", tt.phoneNumber, result, tt.expected)
			}
		})
	}
}

func TestExtractCountryCodeFromNumber(t *testing.T) {
	tests := []struct {
		name             string
		phoneNumber      string
		expectedCountry  string
		expectedDialCode string
		expectedFound    bool
	}{
		{"US number", "+12125690123", "US", "1", true},
		{"US number no +", "12125690123", "US", "1", true},
		{"Mexico number", "+526313118150", "MX", "52", true},
		{"Spain number", "+34915872200", "ES", "34", true},
		{"Portugal number", "+351210942000", "PT", "351", true},
		{"Number with spaces", "+1 212 569 0123", "US", "1", true},
		{"Unsupported country", "+99912345", "", "", false},
		{"Invalid format", "abc123", "", "", false},
		{"Empty string", "", "", "", false},
		{"Just +", "+", "", "", false},
		{"Number without country code", "2125690123", "", "", false},
		{"Germany number", "+49301234567", "DE", "49", true},
		{"Japan number", "+81312345678", "JP", "81", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			country, dialCode, found := extractCountryCodeFromNumber(tt.phoneNumber)
			if country != tt.expectedCountry || dialCode != tt.expectedDialCode || found != tt.expectedFound {
				t.Errorf("extractCountryCodeFromNumber(%s) = (%s, %s, %v), expected (%s, %s, %v)",
					tt.phoneNumber, country, dialCode, found, tt.expectedCountry, tt.expectedDialCode, tt.expectedFound)
			}
		})
	}
}

func TestCleanNumber(t *testing.T) {
	tests := []struct {
		name        string
		phoneNumber string
		expected    string
	}{
		{"With + and spaces", "+1 212 569 0123", "12125690123"},
		{"Only +", "+12125690123", "12125690123"},
		{"Only spaces", "1 212 569 0123", "12125690123"},
		{"No + or spaces", "12125690123", "12125690123"},
		{"Empty string", "", ""},
		{"Only +", "+", ""},
		{"Multiple spaces", "1   212   569", "1212569"},
		{"Leading/trailing spaces", " 123 456 ", "123456"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleanNumber(tt.phoneNumber)
			if result != tt.expected {
				t.Errorf("cleanNumber(%s) = %s, expected %s", tt.phoneNumber, result, tt.expected)
			}
		})
	}
}

func TestStripPlus(t *testing.T) {
	tests := []struct {
		name        string
		phoneNumber string
		expected    string
	}{
		{"With +", "+12125690123", "12125690123"},
		{"Without +", "12125690123", "12125690123"},
		{"Only +", "+", ""},
		{"Empty string", "", ""},
		{"With + and spaces", "+1 212 569 0123", "1 212 569 0123"},
		{"Multiple + (edge case)", "++123", "+123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stripPlus(tt.phoneNumber)
			if result != tt.expected {
				t.Errorf("stripPlus(%s) = %s, expected %s", tt.phoneNumber, result, tt.expected)
			}
		})
	}
}
