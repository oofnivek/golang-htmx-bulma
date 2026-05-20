package crypto

import (
	"fmt"
	"testing"
)

// Run this test using:
// go test -v ./internal/pkg/crypto -run TestGenerateHash
func TestGenerateHash(t *testing.T) {
	// Put the email and password you want to hash here:
	email := "user@email.com"
	password := "password"

	hash, err := HashUserPassword(email, password, 100000)
	if err != nil {
		t.Fatalf("Failed to generate hash: %v", err)
	}

	fmt.Printf("\n==================================================\n")
	fmt.Printf("Generated ASP.NET Core V3 Owner-Linked hash for '%s' (%s):\n", password, email)
	fmt.Printf("%s\n", hash)
	fmt.Printf("==================================================\n\n")
}

func TestVerifyUserPassword(t *testing.T) {
	// 1. Verify the specific user hash provided (which uses owner-linking)
	hash := "AQAAAAIAAYagAAAAEGNVoWtvxXghKHrLAd3HybnMe+VFv00gZEDvFvv0oRlXfdoGWXoHQkwc/rFuxboleQ=="
	email := "user@email.com"
	password := "password"

	result := VerifyUserPassword(email, hash, password)
	if result != PasswordVerificationSuccess {
		t.Fatalf("Expected PasswordVerificationSuccess, got %d", result)
	}

	// 2. Test fallback for standard password hashed without owner-linking (legacy format)
	plainPassword := "LegacySecret123!"
	plainHash, err := HashPasswordV3(plainPassword, 10000)
	if err != nil {
		t.Fatalf("Failed to generate plain hash: %v", err)
	}

	resultFallback := VerifyUserPassword(email, plainHash, plainPassword)
	if resultFallback != PasswordVerificationSuccessRehash {
		t.Fatalf("Expected PasswordVerificationSuccessRehash, got %d", resultFallback)
	}

	// 3. Test newly generated owner-linked password hash
	newLinkedHash, err := HashUserPassword(email, password, 10000)
	if err != nil {
		t.Fatalf("Failed to generate owner-linked hash: %v", err)
	}

	resultNew := VerifyUserPassword(email, newLinkedHash, password)
	if resultNew != PasswordVerificationSuccess {
		t.Fatalf("Expected PasswordVerificationSuccess for new hash, got %d", resultNew)
	}

	// 4. Test failure case
	resultFailed := VerifyUserPassword(email, newLinkedHash, "WrongPassword!")
	if resultFailed != PasswordVerificationFailed {
		t.Fatalf("Expected PasswordVerificationFailed, got %d", resultFailed)
	}
}
