package pow

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/PoW-HC/hashcash/pkg/hash/mock"
)

func TestExtSum(t *testing.T) {
	var (
		resource       = "resource1"
		secret         = "secret"
		rand           = []byte{10}
		bits     int32 = 5
		date           = time.Unix(1648762844, 0)
		expected       = fmt.Sprintf("%s%s%s%d%d", resource, rand, secret, date.Unix(), bits)
		a              = assert.New(t)
		hasher         = (mock.HasherMockParams{
			HashTimes:  1,
			HashReq:    gomock.Eq(expected),
			HashRes:    expected,
			HashResErr: nil,
		}).NewHasher(gomock.NewController(t))
	)

	actual, err := extSum(resource, secret, bits, rand, date, hasher)
	a.Nil(err)
	a.Equal(expected, actual)
}

func TestExtSumErr(t *testing.T) {
	var (
		resource          = "resource"
		secret            = "secret"
		rand              = []byte{10}
		bits        int32 = 5
		date              = time.Unix(1648762844, 0)
		expected          = fmt.Sprintf("%s%s%s%d%d", resource, rand, secret, date.Unix(), bits)
		expectedErr       = fmt.Errorf("expected error")
		a                 = assert.New(t)
		hasher            = (mock.HasherMockParams{
			HashTimes:  1,
			HashReq:    gomock.Eq(expected),
			HashRes:    "",
			HashResErr: expectedErr,
		}).NewHasher(gomock.NewController(t))
	)

	actual, err := extSum(resource, secret, bits, rand, date, hasher)
	a.Empty(actual)
	a.ErrorIs(err, expectedErr)
}

func TestVerifyExt(t *testing.T) {
	hasherErr := fmt.Errorf("expected error")
	ctrl := gomock.NewController(t)

	for _, tCase := range []struct {
		name        string
		hashcash    *Hashcach
		secret      string
		hasherMock  mock.HasherMockParams
		expectedErr error
	}{
		{
			name: "positive",
			hashcash: &Hashcach{
				Resource: "resource",
				Rand:     []byte{10},
				Date:     time.Unix(1648762844, 0),
				Ext:      "resource\nsecret16487628440",
			},
			secret: "secret",
			hasherMock: mock.HasherMockParams{
				HashTimes:  1,
				HashReq:    gomock.Eq("resource\nsecret16487628440"),
				HashRes:    "resource\nsecret16487628440",
				HashResErr: nil,
			},
			expectedErr: nil,
		},
		{
			name: "wrong ext",
			hashcash: &Hashcach{
				Resource: "resource",
				Rand:     []byte{10},
				Date:     time.Unix(1648762844, 0),
				Ext:      "wrong",
			},
			secret: "secret",
			hasherMock: mock.HasherMockParams{
				HashTimes:  1,
				HashReq:    gomock.Eq("resource\nsecret16487628440"),
				HashRes:    "resource\nsecret16487628440",
				HashResErr: nil,
			},
			expectedErr: ErrExtInvalid,
		},
		{
			name: "wrong hasher response",
			hashcash: &Hashcach{
				Resource: "resource",
				Rand:     []byte{10},
				Date:     time.Unix(1648762844, 0),
				Ext:      "resource\nsecret16487628440",
			},
			secret: "secret",
			hasherMock: mock.HasherMockParams{
				HashTimes:  1,
				HashReq:    gomock.Eq("resource\nsecret16487628440"),
				HashRes:    "wrong",
				HashResErr: nil,
			},
			expectedErr: ErrExtInvalid,
		},
		{
			name: "wrong hasher response",
			hashcash: &Hashcach{
				Resource: "resource",
				Rand:     []byte{10},
				Date:     time.Unix(1648762844, 0),
				Ext:      "resource\nsecret16487628440",
			},
			secret: "secret",
			hasherMock: mock.HasherMockParams{
				HashTimes:  1,
				HashReq:    gomock.Eq("resource\nsecret16487628440"),
				HashRes:    "wrong",
				HashResErr: nil,
			},
			expectedErr: ErrExtInvalid,
		},
		{
			name: "hasher error",
			hashcash: &Hashcach{
				Resource: "resource",
				Rand:     []byte{10},
				Date:     time.Unix(1648762844, 0),
				Ext:      "resource\nsecret16487628440",
			},
			secret: "secret",
			hasherMock: mock.HasherMockParams{
				HashTimes:  1,
				HashReq:    gomock.Eq("resource\nsecret16487628440"),
				HashRes:    "",
				HashResErr: hasherErr,
			},
			expectedErr: hasherErr,
		},
	} {
		t.Run(tCase.name, func(t *testing.T) {
			var (
				a      = assert.New(t)
				hasher = tCase.hasherMock.NewHasher(ctrl)
			)

			err := VerifyExt(tCase.secret, hasher)(tCase.hashcash)
			a.ErrorIs(err, tCase.expectedErr)
		})
	}
}

func TestSignExt(t *testing.T) {
	hasherErr := fmt.Errorf("expected error")
	_ = hasherErr
	ctrl := gomock.NewController(t)

	for _, tCase := range []struct {
		name        string
		hashcash    *Hashcach
		secret      string
		hasherMock  mock.HasherMockParams
		expected    string
		expectedErr error
	}{
		{
			name: "positive",
			hashcash: &Hashcach{
				Resource: "resource",
				Rand:     []byte{10},
				Date:     time.Unix(1648762844, 0),
			},
			secret: "secret",
			hasherMock: mock.HasherMockParams{
				HashTimes:  1,
				HashReq:    gomock.Eq("resource\nsecret16487628440"),
				HashRes:    "resource\nsecret16487628440",
				HashResErr: nil,
			},
			expected: "resource\nsecret16487628440",
		},
		{
			name: "hasher error",
			hashcash: &Hashcach{
				Resource: "resource",
				Rand:     []byte{10},
				Date:     time.Unix(1648762844, 0),
				Ext:      "resource\nsecret16487628440",
			},
			secret: "secret",
			hasherMock: mock.HasherMockParams{
				HashTimes:  1,
				HashReq:    gomock.Eq("resource\nsecret16487628440"),
				HashRes:    "",
				HashResErr: hasherErr,
			},
			expectedErr: hasherErr,
		},
	} {
		t.Run(tCase.name, func(t *testing.T) {
			var (
				a      = assert.New(t)
				hasher = tCase.hasherMock.NewHasher(ctrl)
			)

			ext, err := SignExt(tCase.secret, hasher)(tCase.hashcash)
			a.ErrorIs(err, tCase.expectedErr)
			a.Equal(tCase.expected, ext)
		})
	}
}
