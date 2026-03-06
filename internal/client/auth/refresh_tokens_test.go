package auth

import (
	"encoding/hex"
	"testing"
)

func TestMakeRefreshToken(t *testing.T) {
	tests := []struct {
		name        string
		isErr       bool
		expectedErr string
	}{
		{
			name:        "success case basic",
			isErr:       false,
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, errRaw := MakeRefreshToken()
			err := ""
			if errRaw != nil {
				err = errRaw.Error()
			}

			if tt.isErr && err == "" {
				t.Errorf("MakeRefreshToken test %v failed: expected error but got none", tt.name)
				return
			}

			if !tt.isErr && err != "" {
				t.Errorf("MakeRefreshToken test %v failed: unexpected error: %v", tt.name, err)
				return
			}

			if err != tt.expectedErr {
				t.Errorf("MakeRefreshToken test %v failed:\n Expected Error: %v\n Returned Error: %v\n",
					tt.name, tt.expectedErr, err)
				return
			}

			// If no error expected, validate the token properties
			if !tt.isErr {
				// Check token length (32 bytes hex encoded = 64 characters)
				if len(result) != 64 {
					t.Errorf("MakeRefreshToken test %v failed: expected token length 64, got %v", tt.name, len(result))
				}

				// Verify it's valid hex
				_, decodeErr := hex.DecodeString(result)
				if decodeErr != nil {
					t.Errorf("MakeRefreshToken test %v failed: token is not valid hex: %v", tt.name, decodeErr)
				}
			}
		})
	}
}

func TestMakeRefreshTokenUniqueness(t *testing.T) {
	// Generate multiple tokens and ensure they're unique
	tokens := make(map[string]bool)
	iterations := 100

	for i := 0; i < iterations; i++ {
		token, err := MakeRefreshToken()
		if err != nil {
			t.Fatalf("Unexpected error generating token: %v", err)
		}

		if tokens[token] {
			t.Errorf("Duplicate token generated: %v", token)
		}

		tokens[token] = true
	}

	if len(tokens) != iterations {
		t.Errorf("Expected %v unique tokens, got %v", iterations, len(tokens))
	}
}
