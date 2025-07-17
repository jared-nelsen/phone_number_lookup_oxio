package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type PhoneNumberResponse struct {
	PhoneNumber      string `json:"phoneNumber"`
	CountryCode      string `json:"countryCode"`
	AreaCode         string `json:"areaCode"`
	LocalPhoneNumber string `json:"localPhoneNumber"`
}

type ErrorResponse struct {
	PhoneNumber string            `json:"phoneNumber"`
	Error       map[string]string `json:"error"`
}

func processNumberWithCountryCode(phoneNumber, dialCode, countryCode string) (*PhoneNumberResponse, *ErrorResponse) {
	cleanNumber := cleanNumber(phoneNumber)

	// Remove country code from number
	remainingNumber := cleanNumber[len(dialCode):]

	// Get area code length for this country
	areaCodeLength, exists := AreaCodeMap[countryCode]
	if !exists {
		return nil, &ErrorResponse{
			PhoneNumber: phoneNumber,
			Error:       map[string]string{"countryCode": "unsupported country"},
		}
	}

	// Extract area code and local number
	var areaCode, localNumber string
	if areaCodeLength > 0 && len(remainingNumber) > areaCodeLength {
		areaCode = remainingNumber[:areaCodeLength]
		localNumber = remainingNumber[areaCodeLength:]
	} else {
		areaCode = ""
		localNumber = remainingNumber
	}

	// Reconstruct full phone number with +
	fullPhoneNumber := "+" + dialCode + remainingNumber

	return &PhoneNumberResponse{
		PhoneNumber:      fullPhoneNumber,
		CountryCode:      countryCode,
		AreaCode:         areaCode,
		LocalPhoneNumber: localNumber,
	}, nil
}

func processNumberWithoutCountryCode(phoneNumber, countryCode string) (*PhoneNumberResponse, *ErrorResponse) {
	cleanNumber := cleanNumber(phoneNumber)

	// Phone number didn't  have country code, so we need countryCode provided
	if countryCode == "" {
		return nil, &ErrorResponse{
			PhoneNumber: phoneNumber,
			Error:       map[string]string{"countryCode": "required value is missing"},
		}
	}

	// Validate provided country code
	if !validateCountryCodeMeetsISO_3166_1_alpha_2(countryCode) {
		return nil, &ErrorResponse{
			PhoneNumber: phoneNumber,
			Error:       map[string]string{"countryCode": "invalid format"},
		}
	}

	// Get dial code for country
	dialCode, exists := CountryCodeMap[strings.ToUpper(countryCode)]
	if !exists {
		return nil, &ErrorResponse{
			PhoneNumber: phoneNumber,
			Error:       map[string]string{"countryCode": "unsupported country"},
		}
	}

	// Get area code length for this country
	areaCodeLength, exists := AreaCodeMap[strings.ToUpper(countryCode)]
	if !exists {
		return nil, &ErrorResponse{
			PhoneNumber: phoneNumber,
			Error:       map[string]string{"countryCode": "unsupported country"},
		}
	}

	// Extract area code and local number
	var areaCode, localNumber string
	if areaCodeLength > 0 && len(cleanNumber) > areaCodeLength {
		areaCode = cleanNumber[:areaCodeLength]
		localNumber = cleanNumber[areaCodeLength:]
	} else {
		areaCode = ""
		localNumber = cleanNumber
	}

	// Reconstruct full phone number with +
	fullPhoneNumber := "+" + dialCode + cleanNumber

	return &PhoneNumberResponse{
		PhoneNumber:      fullPhoneNumber,
		CountryCode:      strings.ToUpper(countryCode),
		AreaCode:         areaCode,
		LocalPhoneNumber: localNumber,
	}, nil
}

func parsePhoneNumber(phoneNumber, countryCode string) (*PhoneNumberResponse, *ErrorResponse) {
	// Validate the phone number format
	if !validatePhoneNumberFormat(phoneNumber) {
		return nil, &ErrorResponse{
			PhoneNumber: phoneNumber,
			Error:       map[string]string{"phoneNumber": "invalid format"},
		}
	}

	// Validate spaces before any other processing
	if !validateSpaces(phoneNumber) {
		return nil, &ErrorResponse{
			PhoneNumber: phoneNumber,
			Error:       map[string]string{"phoneNumber": "invalid space placement"},
		}
	}

	// Try to extract country code from numbr
	extractedCountryCode, dialCode, hasCountryCodeInNumber := extractCountryCodeFromNumber(phoneNumber)

	if hasCountryCodeInNumber {
		return processNumberWithCountryCode(phoneNumber, dialCode, extractedCountryCode)
	} else {
		return processNumberWithoutCountryCode(phoneNumber, countryCode)
	}
}

func phoneNumberHandler(c *gin.Context) {
	phoneNumber := c.Query("phoneNumber")
	countryCode := c.Query("countryCode")

	if phoneNumber == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			PhoneNumber: "",
			Error:       map[string]string{"phoneNumber": "required parameter is missing"},
		})
		return
	}

	result, errorResp := parsePhoneNumber(phoneNumber, countryCode)
	if errorResp != nil {
		c.JSON(http.StatusBadRequest, errorResp)
		return
	}

	c.JSON(http.StatusOK, result)
}

func main() {
	r := gin.Default()

	// Add the phone numbers endpoint
	r.GET("/v1/phone-numbers", phoneNumberHandler)

	// Start server on port 8080
	fmt.Println("Server starting on port 8080...")
	r.Run(":8080")
}
