package hash

import (
	"crypto/sha1" //nolint:gosec
	"encoding/hex"
	"fmt"
)

type SHA1 struct {
}

// NewSHA1 sha1 hash function
func NewSHA1() *SHA1 {
	return new(SHA1)
}

func (s *SHA1) Hash(str string) (string, error) {
	// Weak cryptographic primitive. I know, I know..
	sha := sha1.New() //nolint:gosec

	_, err := sha.Write([]byte(str))
	if err != nil {
		return "", fmt.Errorf("sha256 hash error: %w", err)
	}

	return hex.EncodeToString(sha.Sum(nil)), err
}
