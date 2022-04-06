package hash

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
)

type SHA512 struct {
}

// NewSHA512 sha512 hash function
func NewSHA512() *SHA512 {
	return new(SHA512)
}

func (s *SHA512) Hash(str string) (string, error) {
	sha := sha512.New()

	_, err := sha.Write([]byte(str))
	if err != nil {
		return "", fmt.Errorf("sha512 hash error: %w", err)
	}

	return hex.EncodeToString(sha.Sum(nil)), err
}
