package crypto

import (
	"fmt"
	"testing"
)

// Run this test using:
// go test -v ./internal/pkg/crypto -run TestGenerateHash
func TestGenerateHash(t *testing.T) {
	// Put the password you want to hash here:
	password := "password"

	hash, err := HashPasswordV3(password, 100000)
	if err != nil {
		t.Fatalf("Failed to generate hash: %v", err)
	}

	fmt.Printf("\n==================================================\n")
	fmt.Printf("Generated ASP.NET Core V3 hash for '%s':\n", password)
	fmt.Printf("%s\n", hash)
	fmt.Printf("==================================================\n\n")
}
