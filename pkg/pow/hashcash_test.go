package pow

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomBytes(t *testing.T) {
	assert.Greater(t, len(randomBytes()), 0)
}

func TestInitHashcash(t *testing.T) {
	var bits int32 = 3
	var resource = "127.0.0.1"
	a := assert.New(t)

	actual, err := InitHashcash(bits, resource, nil)
	a.Nil(err)
	a.Equal(bits, actual.Bits)
	a.Equal(resource, actual.Resource)
}
