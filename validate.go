package main

import (
	"regexp"
	"strings"
)

// CountryCodeMap maps country codes to their dial codes
// I just made a small selection of the entire set for the example
// See: https://www.internationalinsurance.com/calling-codes/?srsltid=AfmBOoqNh14at9fwXCtDsi71BB1cyLoNY1177PY792OyUdkiyMEQDvE0
// We use this to identify the country its from and later its area code (see below)
var CountryCodeMap = map[string]string{
	"US": "1",   // US
	"MX": "52",  // Mexico
	"ES": "34",  // Spain
	"PT": "351", // Portugal
	"CA": "1",   // Canada
	"GB": "44",  // United Kingdom
	"FR": "33",  // France
	"DE": "49",  // Germany
	"IT": "39",  // Italy
	"JP": "81",  // Japan
}

// AreaCodeMap maps country codes to their area code lengths
// Basically it states that in a given country, the area codes in that country have the corresponding link
// We need to know this for parsing
var AreaCodeMap = map[string]int{
	"US": 3, // US i.e 212 for Manhattan
	"MX": 3, // Mexico
	"ES": 3, // Spain
	"PT": 2, // Portugal
	"CA": 3, // Canada
	"GB": 4, // United Kingdom
	"FR": 1, // France
	"DE": 3, // Germany
	"IT": 3, // Italy
	"JP": 1, // Japan
}

func cleanNumber(phoneNumber string) string {
	return strings.ReplaceAll(strings.TrimPrefix(phoneNumber, "+"), " ", "")
}

func stripPlus(phoneNumber string) string {
	return strings.TrimPrefix(phoneNumber, "+")
}

// For ISO 3166-1 explanation see https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2
func validateCountryCodeMeetsISO_3166_1_alpha_2(countryCode string) bool {
	if len(countryCode) != 2 {
		return false
	}
	// Check that it contains only letters
	for _, char := range countryCode {
		if !((char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z')) {
			return false
		}
	}
	return true
}

func validatePhoneNumberFormat(phoneNumber string) bool {
	// Remove + if present for easier parsing
	cleanNumber := stripPlus(phoneNumber)

	// Empty or only spaces is invalid
	if cleanNumber == "" || strings.TrimSpace(cleanNumber) == "" {
		return false
	}

	// Check if it contains only digits and spaces
	validChars := regexp.MustCompile(`^[0-9 ]+$`)
	return validChars.MatchString(cleanNumber)
}

// Checks if spaces are in valid positions: between country, area, and local only
func validateSpaces(phoneNumber string) bool {
	// Remove + if present for easier parsing
	cleanNumber := stripPlus(phoneNumber)

	// Check for consecutive spaces
	if strings.Contains(cleanNumber, "  ") {
		return false
	}

	// Check for leading or trailing spaces
	if strings.HasPrefix(cleanNumber, " ") || strings.HasSuffix(cleanNumber, " ") {
		return false
	}

	// Split by spaces and check each part contains only digits
	parts := strings.Split(cleanNumber, " ")
	for _, part := range parts {
		if part == "" {
			return false
		}
		if !regexp.MustCompile(`^[0-9]+$`).MatchString(part) {
			return false
		}
	}

	// For numbers with spaces, we should have at most 3 parts: country code, area code, local number
	return len(parts) <= 3
}

func extractCountryCodeFromNumber(phoneNumber string) (string, string, bool) {
    cleanNumber := cleanNumber(phoneNumber)

	// There is a wacky edge case where we want to prefer US for 1 instead of Canada
    
    // Check dial codes in order of length (longest first)
	// This is a hacky way to do it to save time
	// If I had more time I could figure out how to sort the map
    dialCodes := []string{"351", "52", "44", "49", "39", "81", "34", "33", "1"}
    
    for _, dialCode := range dialCodes {
        if strings.HasPrefix(cleanNumber, dialCode) {
            if dialCode == "1" {
                return "US", dialCode, true
            }
            // Find the country for this dial code
            for country, code := range CountryCodeMap {
                if code == dialCode {
                    return country, dialCode, true
                }
            }
        }
    }
    
    return "", "", false
}