package auth

import (
	"net/http"
	"testing"
)

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name           string
		headers        http.Header
		isErr          bool
		expectedResult string
		expectedErr    string
	}{
		{
			name: "success case basic",
			headers: http.Header{
				"Authorization": []string{"Bearer token-test-123"},
				"Content-Type":  []string{"application/json"},
				"User-Agent":    []string{"postflow-cli/1.0"},
				"Accept":        []string{"application/json"},
				"Host":          []string{"api.example.com"},
			},
			isErr:          false,
			expectedResult: "token-test-123",
			expectedErr:    "",
		},
		{
			name: "fail case no auth header",
			headers: http.Header{
				"Content-Type": []string{"application/json"},
				"User-Agent":   []string{"postflow-cli/1.0"},
				"Accept":       []string{"application/json"},
				"Host":         []string{"api.example.com"},
			},
			isErr:          true,
			expectedResult: "",
			expectedErr:    "header not found: authorization",
		},
		{
			name: "fail case no bearer prefix",
			headers: http.Header{
				"Authorization": []string{"token-test-123"},
				"Content-Type":  []string{"application/json"},
				"User-Agent":    []string{"postflow-cli/1.0"},
				"Accept":        []string{"application/json"},
				"Host":          []string{"api.example.com"},
			},
			isErr:          true,
			expectedResult: "",
			expectedErr:    "bearer token not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, errRaw := GetBearerToken(tt.headers)
			err := ""
			if errRaw != nil {
				err = errRaw.Error()
			}

			if result != tt.expectedResult || err != tt.expectedErr {
				t.Errorf("Bearer token test %v failed:\n\nExpected Result:%v\n Returned Result:%v\n Expected Error: %v\n Returned Error: %v\n",
					tt.name, tt.expectedResult, result, tt.expectedErr, err)
			}
		})
	}
}
