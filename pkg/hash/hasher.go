package hash

//go:generate mockgen -package=mock -destination=./mock/hasher.go github.com/PoW-HC/hashcash/pkg/hash Hasher

// Hasher interface to hash function
type Hasher interface {
	Hash(str string) (string, error)
}
