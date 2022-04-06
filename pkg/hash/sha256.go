package hash

import (
	"encoding/hex"
	"fmt"

	"github.com/minio/sha256-simd"
)

type SHA256 struct {
}

// NewSHA256 sha256 hash function
func NewSHA256() *SHA256 {
	return new(SHA256)
}

func (s *SHA256) Hash(str string) (string, error) {
	sha := sha256.New()

	_, err := sha.Write([]byte(str))
	if err != nil {
		return "", fmt.Errorf("sha256 hash error: %w", err)
	}

	return hex.EncodeToString(sha.Sum(nil)), err
}
