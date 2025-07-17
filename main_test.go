package main

import (
	"testing"
)

func TestProcessNumberWithCountryCode(t *testing.T) {
	tests := []struct {
		name        string
		phoneNumber string
		dialCode    string
		countryCode string
		expectError bool
		errorKey    string
		expected    *PhoneNumberResponse
	}{
		{
			name:        "Valid US number",
			phoneNumber: "+12125690123",
			dialCode:    "1",
			countryCode: "US",
			expectError: false,
			expected: &PhoneNumberResponse{
				PhoneNumber:      "+12125690123",
				CountryCode:      "US",
				AreaCode:         "212",
				LocalPhoneNumber: "5690123",
			},
		},
		{
			name:        "Valid Mexico number",
			phoneNumber: "+526313118150",
			dialCode:    "52",
			countryCode: "MX",
			expectError: false,
			expected: &PhoneNumberResponse{
				PhoneNumber:      "+526313118150",
				CountryCode:      "MX",
				AreaCode:         "631",
				LocalPhoneNumber: "3118150",
			},
		},
		{
			name:        "Valid Spain number",
			phoneNumber: "34915872200",
			dialCode:    "34",
			countryCode: "ES",
			expectError: false,
			expected: &PhoneNumberResponse{
				PhoneNumber:      "+34915872200",
				CountryCode:      "ES",
				AreaCode:         "915",
				LocalPhoneNumber: "872200",
			},
		},
		{
			name:        "Unsupported country",
			phoneNumber: "+99912345",
			dialCode:    "999",
			countryCode: "XX",
			expectError: true,
			errorKey:    "countryCode",
		},
		{
			name:        "Valid number with spaces",
			phoneNumber: "+1 212 5690123",
			dialCode:    "1",
			countryCode: "US",
			expectError: false,
			expected: &PhoneNumberResponse{
				PhoneNumber:      "+12125690123",
				CountryCode:      "US",
				AreaCode:         "212",
				LocalPhoneNumber: "5690123",
			},
		},
		{
			name:        "Short number without area code",
			phoneNumber: "+81123456",
			dialCode:    "81",
			countryCode: "JP",
			expectError: false,
			expected: &PhoneNumberResponse{
				PhoneNumber:      "+81123456",
				CountryCode:      "JP",
				AreaCode:         "1",
				LocalPhoneNumber: "23456",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, errorResp := processNumberWithCountryCode(tt.phoneNumber, tt.dialCode, tt.countryCode)

			if tt.expectError {
				if errorResp == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if _, exists := errorResp.Error[tt.errorKey]; !exists {
					t.Errorf("Expected error key %s, got %v", tt.errorKey, errorResp.Error)
				}
			} else {
				if errorResp != nil {
					t.Errorf("Expected no error but got %v", errorResp.Error)
					return
				}
				if result.PhoneNumber != tt.expected.PhoneNumber {
					t.Errorf("Expected phoneNumber %s, got %s", tt.expected.PhoneNumber, result.PhoneNumber)
				}
				if result.CountryCode != tt.expected.CountryCode {
					t.Errorf("Expected countryCode %s, got %s", tt.expected.CountryCode, result.CountryCode)
				}
				if result.AreaCode != tt.expected.AreaCode {
					t.Errorf("Expected areaCode %s, got %s", tt.expected.AreaCode, result.AreaCode)
				}
				if result.LocalPhoneNumber != tt.expected.LocalPhoneNumber {
					t.Errorf("Expected localPhoneNumber %s, got %s", tt.expected.LocalPhoneNumber, result.LocalPhoneNumber)
				}
			}
		})
	}
}

func TestProcessNumberWithoutCountryCode(t *testing.T) {
	tests := []struct {
		name        string
		phoneNumber string
		countryCode string
		expectError bool
		errorKey    string
		errorValue  string
		expected    *PhoneNumberResponse
	}{
		{
			name:        "Valid with US country code",
			phoneNumber: "2125690123",
			countryCode: "US",
			expectError: false,
			expected: &PhoneNumberResponse{
				PhoneNumber:      "+12125690123",
				CountryCode:      "US",
				AreaCode:         "212",
				LocalPhoneNumber: "5690123",
			},
		},
		{
			name:        "Valid with MX country code",
			phoneNumber: "6313118150",
			countryCode: "MX",
			expectError: false,
			expected: &PhoneNumberResponse{
				PhoneNumber:      "+526313118150",
				CountryCode:      "MX",
				AreaCode:         "631",
				LocalPhoneNumber: "3118150",
			},
		},
		{
			name:        "Missing country code",
			phoneNumber: "2125690123",
			countryCode: "",
			expectError: true,
			errorKey:    "countryCode",
			errorValue:  "required value is missing",
		},
		{
			name:        "Invalid country code format",
			phoneNumber: "2125690123",
			countryCode: "ESP",
			expectError: true,
			errorKey:    "countryCode",
			errorValue:  "invalid format",
		},
		{
			name:        "Unsupported country code",
			phoneNumber: "2125690123",
			countryCode: "XX",
			expectError: true,
			errorKey:    "countryCode",
			errorValue:  "unsupported country",
		},
		{
			name:        "Valid with lowercase country code",
			phoneNumber: "2125690123",
			countryCode: "us",
			expectError: false,
			expected: &PhoneNumberResponse{
				PhoneNumber:      "+12125690123",
				CountryCode:      "US",
				AreaCode:         "212",
				LocalPhoneNumber: "5690123",
			},
		},
		{
			name:        "Valid with spaces",
			phoneNumber: "212 569 0123",
			countryCode: "US",
			expectError: false,
			expected: &PhoneNumberResponse{
				PhoneNumber:      "+12125690123",
				CountryCode:      "US",
				AreaCode:         "212",
				LocalPhoneNumber: "5690123",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, errorResp := processNumberWithoutCountryCode(tt.phoneNumber, tt.countryCode)

			if tt.expectError {
				if errorResp == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if errorResp.Error[tt.errorKey] != tt.errorValue {
					t.Errorf("Expected error %s: %s, got %s: %s", tt.errorKey, tt.errorValue, tt.errorKey, errorResp.Error[tt.errorKey])
				}
			} else {
				if errorResp != nil {
					t.Errorf("Expected no error but got %v", errorResp.Error)
					return
				}
				if result.PhoneNumber != tt.expected.PhoneNumber {
					t.Errorf("Expected phoneNumber %s, got %s", tt.expected.PhoneNumber, result.PhoneNumber)
				}
				if result.CountryCode != tt.expected.CountryCode {
					t.Errorf("Expected countryCode %s, got %s", tt.expected.CountryCode, result.CountryCode)
				}
				if result.AreaCode != tt.expected.AreaCode {
					t.Errorf("Expected areaCode %s, got %s", tt.expected.AreaCode, result.AreaCode)
				}
				if result.LocalPhoneNumber != tt.expected.LocalPhoneNumber {
					t.Errorf("Expected localPhoneNumber %s, got %s", tt.expected.LocalPhoneNumber, result.LocalPhoneNumber)
				}
			}
		})
	}
}

func TestParsePhoneNumber(t *testing.T) {
	tests := []struct {
		name        string
		phoneNumber string
		countryCode string
		expectError bool
		errorKey    string
		errorValue  string
		expected    *PhoneNumberResponse
	}{
		{
			name:        "Valid US number with country code",
			phoneNumber: "+12125690123",
			countryCode: "",
			expectError: false,
			expected: &PhoneNumberResponse{
				PhoneNumber:      "+12125690123",
				CountryCode:      "US",
				AreaCode:         "212",
				LocalPhoneNumber: "5690123",
			},
		},
		{
			name:        "Valid Mexico number with spaces",
			phoneNumber: "+52 631 3118150",
			countryCode: "",
			expectError: false,
			expected: &PhoneNumberResponse{
				PhoneNumber:      "+526313118150",
				CountryCode:      "MX",
				AreaCode:         "631",
				LocalPhoneNumber: "3118150",
			},
		},

		{
			name:        "Valid with provided country code",
			phoneNumber: "631 311 8150",
			countryCode: "MX",
			expectError: false,
			expected: &PhoneNumberResponse{
				PhoneNumber:      "+526313118150",
				CountryCode:      "MX",
				AreaCode:         "631",
				LocalPhoneNumber: "3118150",
			},
		},
		{
			name:        "Invalid format with letters",
			phoneNumber: "abc123",
			countryCode: "",
			expectError: true,
			errorKey:    "phoneNumber",
			errorValue:  "invalid format",
		},
		{
			name:        "Invalid space placement",
			phoneNumber: "351 21 094 2000",
			countryCode: "",
			expectError: true,
			errorKey:    "phoneNumber",
			errorValue:  "invalid space placement",
		},
		{
			name:        "Missing country code",
			phoneNumber: "631 311 8150",
			countryCode: "",
			expectError: true,
			errorKey:    "countryCode",
			errorValue:  "required value is missing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, errorResp := parsePhoneNumber(tt.phoneNumber, tt.countryCode)

			if tt.expectError {
				if errorResp == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if errorResp.Error[tt.errorKey] != tt.errorValue {
					t.Errorf("Expected error %s: %s, got %s: %s", tt.errorKey, tt.errorValue, tt.errorKey, errorResp.Error[tt.errorKey])
				}
			} else {
				if errorResp != nil {
					t.Errorf("Expected no error but got %v", errorResp.Error)
					return
				}
				if result.PhoneNumber != tt.expected.PhoneNumber {
					t.Errorf("Expected phoneNumber %s, got %s", tt.expected.PhoneNumber, result.PhoneNumber)
				}
				if result.CountryCode != tt.expected.CountryCode {
					t.Errorf("Expected countryCode %s, got %s", tt.expected.CountryCode, result.CountryCode)
				}
				if result.AreaCode != tt.expected.AreaCode {
					t.Errorf("Expected areaCode %s, got %s", tt.expected.AreaCode, result.AreaCode)
				}
				if result.LocalPhoneNumber != tt.expected.LocalPhoneNumber {
					t.Errorf("Expected localPhoneNumber %s, got %s", tt.expected.LocalPhoneNumber, result.LocalPhoneNumber)
				}
			}
		})
	}
}
