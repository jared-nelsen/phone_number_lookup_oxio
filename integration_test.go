package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/v1/phone-numbers", phoneNumberHandler)
	return r
}

func TestPhoneNumberHandlerIntegration(t *testing.T) {
	router := setupTestRouter()

	tests := []struct {
		name           string
		phoneNumber    string
		countryCode    string
		expectedStatus int
		expectError    bool
		errorKey       string
		errorValue     string
		expectedResult *PhoneNumberResponse
	}{
		{
			name:           "Valid US number with +",
			phoneNumber:    "+12125690123",
			countryCode:    "",
			expectedStatus: http.StatusOK,
			expectError:    false,
			expectedResult: &PhoneNumberResponse{
				PhoneNumber:      "+12125690123",
				CountryCode:      "US",
				AreaCode:         "212",
				LocalPhoneNumber: "5690123",
			},
		},
		{
			name:           "Valid Mexico number with spaces",
			phoneNumber:    "+52 631 3118150",
			countryCode:    "",
			expectedStatus: http.StatusOK,
			expectError:    false,
			expectedResult: &PhoneNumberResponse{
				PhoneNumber:      "+526313118150",
				CountryCode:      "MX",
				AreaCode:         "631",
				LocalPhoneNumber: "3118150",
			},
		},
		{
			name:           "Valid Spain number without +",
			phoneNumber:    "34 915 872200",
			countryCode:    "",
			expectedStatus: http.StatusOK,
			expectError:    false,
			expectedResult: &PhoneNumberResponse{
				PhoneNumber:      "+34915872200",
				CountryCode:      "ES",
				AreaCode:         "915",
				LocalPhoneNumber: "872200",
			},
		},
		{
			name:           "Valid with provided country code",
			phoneNumber:    "631 311 8150",
			countryCode:    "MX",
			expectedStatus: http.StatusOK,
			expectError:    false,
			expectedResult: &PhoneNumberResponse{
				PhoneNumber:      "+526313118150",
				CountryCode:      "MX",
				AreaCode:         "631",
				LocalPhoneNumber: "3118150",
			},
		},
		{
			name:           "Invalid - too many spaces",
			phoneNumber:    "351 21 094 2000",
			countryCode:    "",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
			errorKey:       "phoneNumber",
			errorValue:     "invalid space placement",
		},
		{
			name:           "Invalid - missing country code",
			phoneNumber:    "631 311 8150",
			countryCode:    "",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
			errorKey:       "countryCode",
			errorValue:     "required value is missing",
		},
		{
			name:           "Invalid - wrong country code format",
			phoneNumber:    "631 311 8150",
			countryCode:    "ESP",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
			errorKey:       "countryCode",
			errorValue:     "invalid format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := url.Values{}
			if tt.phoneNumber != "" {
				params.Add("phoneNumber", tt.phoneNumber)
			}
			if tt.countryCode != "" {
				params.Add("countryCode", tt.countryCode)
			}

			req, err := http.NewRequest("GET", "/v1/phone-numbers?"+params.Encode(), nil)
			if err != nil {
				t.Fatalf("Could not create request: %v", err)
			}

			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectError {
				var errorResp ErrorResponse
				if err := json.Unmarshal(w.Body.Bytes(), &errorResp); err != nil {
					t.Fatalf("Could not parse error response: %v", err)
				}

				if errorResp.Error[tt.errorKey] != tt.errorValue {
					t.Errorf("Expected error %s: %s, got %s: %s",
						tt.errorKey, tt.errorValue, tt.errorKey, errorResp.Error[tt.errorKey])
				}

				expectedPhoneNumber := tt.phoneNumber
				if tt.phoneNumber == "" {
					expectedPhoneNumber = ""
				}
				if errorResp.PhoneNumber != expectedPhoneNumber {
					t.Errorf("Expected error phoneNumber %s, got %s", expectedPhoneNumber, errorResp.PhoneNumber)
				}
			} else {
				var result PhoneNumberResponse
				if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
					t.Fatalf("Could not parse success response: %v", err)
				}

				if result.PhoneNumber != tt.expectedResult.PhoneNumber {
					t.Errorf("Expected phoneNumber %s, got %s", tt.expectedResult.PhoneNumber, result.PhoneNumber)
				}
				if result.CountryCode != tt.expectedResult.CountryCode {
					t.Errorf("Expected countryCode %s, got %s", tt.expectedResult.CountryCode, result.CountryCode)
				}
				if result.AreaCode != tt.expectedResult.AreaCode {
					t.Errorf("Expected areaCode %s, got %s", tt.expectedResult.AreaCode, result.AreaCode)
				}
				if result.LocalPhoneNumber != tt.expectedResult.LocalPhoneNumber {
					t.Errorf("Expected localPhoneNumber %s, got %s", tt.expectedResult.LocalPhoneNumber, result.LocalPhoneNumber)
				}
			}
		})
	}
}
