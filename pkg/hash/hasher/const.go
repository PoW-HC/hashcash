package hasher

//go:generate go-enum -f=$GOFILE --marshal

// Hasher is an enumeration of hash types that are allowed.
/* ENUM(
	SHA1,
	SHA256,
	SHA512,
)
*/
type Hasher int
