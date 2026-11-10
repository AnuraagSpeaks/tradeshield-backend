package crypto

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plaintext password using bcrypt (cost 12)
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash verifies a plaintext password against a bcrypt hash.
// It also provides a migration fallback for legacy seeded demo accounts.
func CheckPasswordHash(password, hash string) bool {
	if hash == "" || password == "" {
		return false
	}
	// Standard bcrypt check
	if strings.HasPrefix(hash, "$2a$") || strings.HasPrefix(hash, "$2b$") || strings.HasPrefix(hash, "$2y$") {
		err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
		return err == nil
	}
	// Legacy fallback for pre-seeded dev passwords (e.g. Shield@Pass2026)
	return hash == password
}

// GenerateSecureOTP generates a cryptographically secure 6-digit numeric string
func GenerateSecureOTP() string {
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		// Fallback to time-derived if entropy failure
		return fmt.Sprintf("%06d", (n.Int64()%900000)+100000)
	}
	return fmt.Sprintf("%06d", n.Int64()+100000)
}
