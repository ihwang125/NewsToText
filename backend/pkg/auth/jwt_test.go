package auth

import (
	"testing"
	"time"
)

func TestJWTManager_GenerateToken(t *testing.T) {
	jwtManager := NewJWTManager("test-secret-key")

	userID := uint(123)
	email := "test@example.com"

	token, err := jwtManager.GenerateToken(userID, email)

	if err != nil {
		t.Errorf("Unexpected error generating token: %v", err)
	}

	if token == "" {
		t.Errorf("Expected non-empty token")
	}

	// Validate the generated token
	claims, err := jwtManager.ValidateToken(token)
	if err != nil {
		t.Errorf("Failed to validate generated token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Expected UserID %d but got %d", userID, claims.UserID)
	}

	if claims.Email != email {
		t.Errorf("Expected email %s but got %s", email, claims.Email)
	}
}

func TestJWTManager_ValidateToken(t *testing.T) {
	jwtManager := NewJWTManager("test-secret-key")

	tests := []struct {
		name    string
		setup   func() string
		wantErr bool
	}{
		{
			name: "Valid token",
			setup: func() string {
				token, _ := jwtManager.GenerateToken(123, "test@example.com")
				return token
			},
			wantErr: false,
		},
		{
			name: "Invalid token format",
			setup: func() string {
				return "invalid-token"
			},
			wantErr: true,
		},
		{
			name: "Empty token",
			setup: func() string {
				return ""
			},
			wantErr: true,
		},
		{
			name: "Token with wrong secret",
			setup: func() string {
				wrongManager := NewJWTManager("wrong-secret")
				token, _ := wrongManager.GenerateToken(123, "test@example.com")
				return token
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := tt.setup()
			claims, err := jwtManager.ValidateToken(token)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if claims == nil {
				t.Errorf("Expected claims but got nil")
			}
		})
	}
}

func TestJWTManager_RefreshToken(t *testing.T) {
	jwtManager := NewJWTManager("test-secret-key")

	// Generate initial token
	originalToken, err := jwtManager.GenerateToken(123, "test@example.com")
	if err != nil {
		t.Fatalf("Failed to generate original token: %v", err)
	}

	// Wait a moment to ensure different timestamps
	time.Sleep(time.Millisecond * 10)

	// Refresh the token
	refreshedToken, err := jwtManager.RefreshToken(originalToken)
	if err != nil {
		t.Errorf("Failed to refresh token: %v", err)
	}

	if refreshedToken == originalToken {
		t.Errorf("Refreshed token should be different from original")
	}

	// Validate the refreshed token
	claims, err := jwtManager.ValidateToken(refreshedToken)
	if err != nil {
		t.Errorf("Failed to validate refreshed token: %v", err)
	}

	if claims.UserID != 123 {
		t.Errorf("Expected UserID 123 but got %d", claims.UserID)
	}

	if claims.Email != "test@example.com" {
		t.Errorf("Expected email test@example.com but got %s", claims.Email)
	}
}