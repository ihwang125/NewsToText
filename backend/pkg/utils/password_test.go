package utils

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword123"

	hashedPassword, err := HashPassword(password)

	if err != nil {
		t.Errorf("Unexpected error hashing password: %v", err)
	}

	if hashedPassword == "" {
		t.Errorf("Expected non-empty hashed password")
	}

	if hashedPassword == password {
		t.Errorf("Hashed password should not be the same as original password")
	}

	// Verify that hashing the same password twice produces different results
	hashedPassword2, err := HashPassword(password)
	if err != nil {
		t.Errorf("Unexpected error hashing password again: %v", err)
	}

	if hashedPassword == hashedPassword2 {
		t.Errorf("Hashing the same password twice should produce different results due to salt")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "testpassword123"
	wrongPassword := "wrongpassword"

	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password for testing: %v", err)
	}

	tests := []struct {
		name     string
		password string
		hash     string
		expected bool
	}{
		{
			name:     "Correct password",
			password: password,
			hash:     hashedPassword,
			expected: true,
		},
		{
			name:     "Wrong password",
			password: wrongPassword,
			hash:     hashedPassword,
			expected: false,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     hashedPassword,
			expected: false,
		},
		{
			name:     "Empty hash",
			password: password,
			hash:     "",
			expected: false,
		},
		{
			name:     "Invalid hash",
			password: password,
			hash:     "invalid-hash",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckPasswordHash(tt.password, tt.hash)

			if result != tt.expected {
				t.Errorf("Expected %v but got %v", tt.expected, result)
			}
		})
	}
}