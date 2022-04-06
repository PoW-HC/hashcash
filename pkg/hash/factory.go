package hash

import (
	"fmt"
	"strings"

	"github.com/PoW-HC/hashcash/pkg/hash/hasher"
)

func NewHasher(hasherName string) (Hasher, error) {
	h, err := hasher.ParseHasher(strings.ToUpper(hasherName))
	if err != nil {
		return nil, fmt.Errorf("parse hasher by name error: %w", err)
	}

	switch h {
	case hasher.HasherSHA1:
		return NewSHA1(), nil
	case hasher.HasherSHA256:
		return NewSHA256(), nil
	case hasher.HasherSHA512:
		return NewSHA512(), nil
	default:
		return NewSHA256(), nil
	}
}
