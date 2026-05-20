package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"hash"

	"golang.org/x/crypto/pbkdf2"
)

const (
	prfHMACSHA256 = 1
	prfHMACSHA512 = 2
)

// HashPasswordV3 mimics the Microsoft.AspNetCore.Identity.PasswordHasher implementation (Version 3).
// PRF=2 (HMAC-SHA512), matching the ASP.NET Core Identity format.
func HashPasswordV3(password string, iterations int) (string, error) {
	const (
		saltSize   = 16
		subkeySize = 32
		prf        = prfHMACSHA512
	)

	// 1. Generate a random salt
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// 2. Derive the subkey using PBKDF2 with HMAC-SHA512
	subkey := pbkdf2.Key([]byte(password), salt, iterations, subkeySize, sha512.New)

	// 3. Assemble the binary payload:
	// [0] Version | [1-4] PRF | [5-8] Iterations | [9-12] SaltSize | [13..] Salt | [...] Subkey
	buffer := make([]byte, 13+saltSize+subkeySize)

	buffer[0] = 0x01 // Format Version 3

	binary.BigEndian.PutUint32(buffer[1:], prf)
	binary.BigEndian.PutUint32(buffer[5:], uint32(iterations))
	binary.BigEndian.PutUint32(buffer[9:], uint32(saltSize))

	copy(buffer[13:], salt)
	copy(buffer[13+saltSize:], subkey)

	// 4. Return as Base64 string
	return base64.StdEncoding.EncodeToString(buffer), nil
}

// VerifyPasswordV3 checks a plaintext password against a stored ASP.NET Core Identity
// Version 3 hash. Returns nil if the password matches, an error otherwise.
// Supports both PRF=1 (HMAC-SHA256) and PRF=2 (HMAC-SHA512).
func VerifyPasswordV3(password, encodedHash string) error {
	// 1. Decode base64
	data, err := base64.StdEncoding.DecodeString(encodedHash)
	if err != nil {
		return errors.New("invalid hash encoding")
	}

	// 2. Minimum length: 1 (version) + 4 (prf) + 4 (iterations) + 4 (saltSize) + 1 (salt) + 1 (subkey)
	if len(data) < 15 {
		return errors.New("hash too short")
	}

	// 3. Check version byte
	if data[0] != 0x01 {
		return errors.New("unsupported hash version")
	}

	// 4. Parse header fields
	prf := binary.BigEndian.Uint32(data[1:5])
	iterations := int(binary.BigEndian.Uint32(data[5:9]))
	saltSize := int(binary.BigEndian.Uint32(data[9:13]))

	// 5. Validate lengths
	if 13+saltSize >= len(data) {
		return errors.New("hash data is malformed")
	}
	salt := data[13 : 13+saltSize]
	expectedSubkey := data[13+saltSize:]
	subkeySize := len(expectedSubkey)

	// 6. Select hash function based on PRF
	var hashFunc func() hash.Hash
	switch prf {
	case prfHMACSHA256:
		hashFunc = sha256.New
	case prfHMACSHA512:
		hashFunc = sha512.New
	default:
		return errors.New("unsupported PRF algorithm")
	}

	// 7. Re-derive the subkey using the extracted parameters
	actualSubkey := pbkdf2.Key([]byte(password), salt, iterations, subkeySize, hashFunc)

	// 8. Compare in constant time to prevent timing attacks
	if subtle.ConstantTimeCompare(actualSubkey, expectedSubkey) != 1 {
		return errors.New("password does not match")
	}

	return nil
}

// PasswordVerificationResult mirrors the PasswordVerificationResult enum in ASP.NET Core Identity.
type PasswordVerificationResult int

const (
	PasswordVerificationFailed        PasswordVerificationResult = 0
	PasswordVerificationSuccess       PasswordVerificationResult = 1
	PasswordVerificationSuccessRehash PasswordVerificationResult = 2
)

// HashUserPassword hashes a password linked to the user's email, using standard ASP.NET Core V3 PBKDF2 format.
func HashUserPassword(email, password string, iterations int) (string, error) {
	return HashPasswordV3(ownerLinkedPassword(email, password), iterations)
}

// VerifyUserPassword verifies a hashed password against a provided plaintext password.
// It first attempts verification using the owner-linked format (email + "\n" + password).
// If that fails, it falls back to verifying the plaintext password directly (supporting legacy hashes).
func VerifyUserPassword(email, hashedPassword, providedPassword string) PasswordVerificationResult {
	// 1. Verify using owner-linked password
	err := VerifyPasswordV3(ownerLinkedPassword(email, providedPassword), hashedPassword)
	if err == nil {
		return PasswordVerificationSuccess
	}

	// 2. Fall back to standard password (old method)
	err = VerifyPasswordV3(providedPassword, hashedPassword)
	if err == nil {
		return PasswordVerificationSuccessRehash
	}

	return PasswordVerificationFailed
}

// ownerLinkedPassword concatenates email, a newline, and password.
func ownerLinkedPassword(email, password string) string {
	return email + "\n" + password
}
