package auth

import (
	"net/http"
	"testing"
)

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name          string
		headers       http.Header
		expectedToken string
		expectedError bool
		errorMsg      string
	}{
		{
			name:          "Valid Bearer Token",
			headers:       http.Header{"Authorization": []string{"Bearer abc123"}},
			expectedToken: "abc123",
			expectedError: false,
		},
		{
			name:          "Missing Authorization Header",
			headers:       http.Header{},
			expectedToken: "",
			expectedError: true,
			errorMsg:      "No Authorization header",
		},
		{
			name:          "Empty Authorization Header",
			headers:       http.Header{"Authorization": []string{""}},
			expectedToken: "",
			expectedError: true,
			errorMsg:      "No Authorization header",
		},
		{
			name:          "Authorization Header Without Bearer Prefix",
			headers:       http.Header{"Authorization": []string{"abc123"}},
			expectedToken: "",
			expectedError: true,
			errorMsg:      "Invalid Authorization header",
		},
		{
			name:          "Authorization Header With Different Prefix",
			headers:       http.Header{"Authorization": []string{"Basic abc123"}},
			expectedToken: "",
			expectedError: true,
			errorMsg:      "Invalid Authorization header",
		},
		{
			name:          "Authorization Header With Extra Parts",
			headers:       http.Header{"Authorization": []string{"Bearer abc123 extra"}},
			expectedToken: "",
			expectedError: true,
			errorMsg:      "Invalid Authorization header",
		},
		{
			name:          "Authorization Header With Only Bearer",
			headers:       http.Header{"Authorization": []string{"Bearer"}},
			expectedToken: "",
			expectedError: true,
			errorMsg:      "Invalid Authorization header",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GetBearerToken(tt.headers)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected an error but got none")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("Expected error message '%s', but got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}

			if token != tt.expectedToken {
				t.Errorf("Expected token '%s', but got '%s'", tt.expectedToken, token)
			}
		})
	}
}
