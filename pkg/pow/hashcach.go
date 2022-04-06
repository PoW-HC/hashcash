package pow

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math"
	"math/big"
	mrand "math/rand"
	"strconv"
	"time"
)

const (
	versionV1 = 1
)

var (
	ErrExtInvalid    = fmt.Errorf("extension sum invalid")
	ErrHashcashEmpty = fmt.Errorf("hashcash empty")
)

// Hashcach struct to marshal and unmarshal hashcach to string or proto buf
type Hashcach struct {
	Version  int32
	Bits     int32
	Date     time.Time
	Resource string
	Ext      string
	Rand     []byte
	Counter  int64
}

func NewHashcach(
	version int32,
	bits int32,
	date time.Time,
	resource string,
	ext string,
	rand []byte,
	counter int64,
) *Hashcach {
	return &Hashcach{
		Version:  version,
		Bits:     bits,
		Date:     date,
		Resource: resource,
		Ext:      ext,
		Rand:     rand,
		Counter:  counter,
	}
}

// InitHashcash initiate new hashcash
func InitHashcash(bits int32, resource string, extGenerator ExtGeneratorFunc) (*Hashcach, error) {
	t := time.Now()
	randBytes := randomBytes()

	hc := NewHashcach(
		versionV1,
		bits,
		t,
		resource,
		"",
		randBytes,
		0,
	)

	if extGenerator != nil {
		ext, err := extGenerator(hc)
		if err != nil {
			return nil, fmt.Errorf("ext generator error: %w", err)
		}
		hc.Ext = ext
	}

	return hc, nil
}

// String implements fmt.Stringer interface to get string hashcash
func (h *Hashcach) String() string {
	var buf bytes.Buffer
	buf.WriteString(strconv.Itoa(int(h.Version)))
	buf.WriteString(":")
	buf.WriteString(strconv.Itoa(int(h.Bits)))
	buf.WriteString(":")
	buf.WriteString(strconv.Itoa(int(h.Date.Unix())))
	buf.WriteString(":")
	buf.WriteString(h.Resource)
	buf.WriteString(":")
	buf.WriteString(h.Ext)
	buf.WriteString(":")
	buf.WriteString(base64.StdEncoding.EncodeToString(h.Rand))
	buf.WriteString(":")
	buf.WriteString(base64.StdEncoding.EncodeToString([]byte(strconv.FormatInt(h.Counter, 16))))
	return buf.String()
}

func randomBytes() []byte {
	b, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		// this is fallback in case of unweak rand function fall with error (exceptional case)
		b = big.NewInt(mrand.Int63n(math.MaxInt64)) //nolint:gosec
	}

	return b.Bytes()
}
