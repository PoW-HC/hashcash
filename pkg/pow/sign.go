package pow

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	"github.com/PoW-HC/hashcash/pkg/hash"
)

type ExtGeneratorFunc func(*Hashcach) (string, error)
type ExtValidatorFunc func(*Hashcach) error

// SignExt creates signed extension
// See extSum description for hash generating details
func SignExt(secret string, hasher hash.Hasher) ExtGeneratorFunc {
	return func(h *Hashcach) (string, error) {
		ext, err := extSum(h.Resource, secret, h.Bits, h.Rand, h.Date, hasher)
		if err != nil {
			return "", err
		}

		return ext, nil
	}
}

// VerifyExt verify extension from hashcash to validate hashcash was provided by server.
// See extSum description for hash generating details
func VerifyExt(secret string, hasher hash.Hasher) ExtValidatorFunc {
	return func(h *Hashcach) error {
		extSum, err := extSum(h.Resource, secret, h.Bits, h.Rand, h.Date, hasher)
		if err != nil {
			return fmt.Errorf("verify ext sum error: %w", err)
		}

		if h.Ext != extSum {
			return ErrExtInvalid
		}

		return nil
	}
}

// extSum generates hash sum with hasher interface from fields:
//    - resource  - ip address
//    - randBytes - random number
//    - secret    - secret known only on server
//    - time      - timestamp
func extSum(resource, secret string, bits int32, randBytes []byte, t time.Time, hasher hash.Hasher) (string, error) {
	var ext bytes.Buffer
	ext.WriteString(resource)
	ext.Write(randBytes)
	ext.WriteString(secret)
	ext.WriteString(strconv.Itoa(int(t.Unix())))
	ext.WriteString(strconv.Itoa(int(bits)))

	extSum, err := hasher.Hash(ext.String())
	if err != nil {
		return "", fmt.Errorf("calculate hashcash ext hash sum error: %w", err)
	}

	return extSum, nil
}
