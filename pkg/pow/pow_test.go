package pow

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/PoW-HC/hashcash/pkg/hash/mock"
)

var zeros = make([]rune, maxHashSize)

func TestMain(m *testing.M) {
	for i := range zeros {
		zeros[i] = zero
	}
	os.Exit(m.Run())
}

func TestIsHashCorrect(t *testing.T) {
	for _, tCase := range []struct {
		name     string
		hash     string
		zeros    int
		expected bool
	}{
		{
			name:     "positive",
			hash:     "00000e89df98a05e524fdcd29d8040d64d0259e2d5109ca1998e567a3c1c1c68",
			zeros:    5,
			expected: true,
		},
		{
			name:     "wrong 5 zeros",
			hash:     "00000e89df98a05e524fdcd29d8040d64d0259e2d5109ca1998e567a3c1c1c68",
			zeros:    6,
			expected: false,
		},
		{
			name:     "wrong 0",
			hash:     "d59d15c9a1842bc4563897803799e94f1f242d7e7e8c618f047e068211543998",
			zeros:    5,
			expected: false,
		},
		{
			name:     "too short",
			hash:     "0000",
			zeros:    6,
			expected: false,
		},
	} {
		t.Run(tCase.name, func(t *testing.T) {
			actual := isHashCorrect(tCase.hash, zeros, tCase.zeros)
			assert.Equal(t, tCase.expected, actual)
		})
	}
}

func TestPowCompute(t *testing.T) {
	hasherErr := fmt.Errorf("expected error")
	ctrl := gomock.NewController(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	deadCTX, deadCancel := context.WithCancel(context.Background())
	deadCancel()

	for _, tCase := range []struct {
		name        string
		ctx         context.Context
		hashcash    *Hashcach
		max         int64
		hasherMock  mock.HasherMockParams
		expected    *Hashcach
		expectedErr error
	}{
		{
			name: "positive",
			ctx:  ctx,
			hashcash: &Hashcach{
				Bits:     5,
				Resource: "resource",
				Rand:     []byte{10},
				Date:     time.Unix(1648762844, 0),
				Ext:      "resource\nsecret1648762844",
			},
			max: 1,
			hasherMock: mock.HasherMockParams{
				HashTimes:  1,
				HashReq:    gomock.Eq("0:5:1648762844:resource:resource\nsecret1648762844:Cg==:MA=="),
				HashRes:    "00000e89df98a05e524fdcd29d8040d64d0259e2d5109ca1998e567a3c1c1c68",
				HashResErr: nil,
			},
			expected: &Hashcach{
				Bits:     5,
				Resource: "resource",
				Rand:     []byte{10},
				Date:     time.Unix(1648762844, 0),
				Ext:      "resource\nsecret1648762844",
				Counter:  0,
			},
			expectedErr: nil,
		},
		{
			name: "hasher error",
			ctx:  ctx,
			hashcash: &Hashcach{
				Resource: "resource",
				Rand:     []byte{10},
				Date:     time.Unix(1648762844, 0),
				Ext:      "resource\nsecret1648762844",
			},
			max: 1,
			hasherMock: mock.HasherMockParams{
				HashTimes:  1,
				HashReq:    gomock.Eq("0:0:1648762844:resource:resource\nsecret1648762844:Cg==:MA=="),
				HashRes:    "",
				HashResErr: hasherErr,
			},
			expected:    nil,
			expectedErr: hasherErr,
		},
		{
			name: "deadline exceeded",
			ctx:  ctx,
			hashcash: &Hashcach{
				Bits:     5,
				Resource: "resource",
				Rand:     []byte{10},
				Date:     time.Unix(1648762844, 0),
				Ext:      "resource\nsecret1648762844",
			},
			max: 1,
			hasherMock: mock.HasherMockParams{
				HashTimes:  2,
				HashReq:    gomock.Any(),
				HashRes:    "d59d15c9a1842bc4563897803799e94f1f242d7e7e8c618f047e068211543998",
				HashResErr: nil,
			},
			expected:    nil,
			expectedErr: ErrMaxIterationsExceeded,
		},
		{
			name: "dead ctx",
			ctx:  deadCTX,
			hashcash: &Hashcach{
				Bits:     5,
				Resource: "resource",
				Rand:     []byte{10},
				Date:     time.Unix(1648762844, 0),
				Ext:      "resource\nsecret1648762844",
			},
			max:         1,
			expected:    nil,
			expectedErr: ErrMaxIterationsExceeded,
		},
	} {
		t.Run(tCase.name, func(t *testing.T) {
			var (
				a      = assert.New(t)
				hasher = tCase.hasherMock.NewHasher(ctrl)
				pow    = New(hasher)
			)

			actual, err := pow.Compute(tCase.ctx, tCase.hashcash, tCase.max)
			a.Equal(tCase.expected, actual)
			a.ErrorIs(err, tCase.expectedErr)
		})
	}
}

func TestPowVerify(t *testing.T) {
	hasherErr := fmt.Errorf("expected error")
	ctrl := gomock.NewController(t)
	now := time.Now()

	for _, tCase := range []struct {
		name        string
		powOptions  []Options
		hashcash    *Hashcach
		resource    string
		hasherMock  mock.HasherMockParams
		expectedErr error
	}{
		{
			name: "positive",
			hashcash: &Hashcach{
				Bits:     5,
				Resource: "resource",
				Rand:     []byte{10},
				Date:     now,
				Ext:      "resource\nsecret1648762844",
			},
			resource: "resource",
			hasherMock: mock.HasherMockParams{
				HashTimes:  1,
				HashReq:    gomock.Any(),
				HashRes:    "00000e89df98a05e524fdcd29d8040d64d0259e2d5109ca1998e567a3c1c1c68",
				HashResErr: nil,
			},
			expectedErr: nil,
		},
		{
			name: "positive validate ext",
			hashcash: &Hashcach{
				Bits:     5,
				Resource: "resource",
				Rand:     []byte{10},
				Date:     now,
				Ext:      "resource\nsecret1648762844",
			},
			powOptions: []Options{
				WithValidateExtFunc(func(h *Hashcach) error {
					assert.Equal(
						t,
						&Hashcach{
							Bits:     5,
							Resource: "resource",
							Rand:     []byte{10},
							Date:     now,
							Ext:      "resource\nsecret1648762844",
						},
						h,
					)
					return nil
				}),
			},
			resource: "resource",
			hasherMock: mock.HasherMockParams{
				HashTimes:  1,
				HashReq:    gomock.Any(),
				HashRes:    "00000e89df98a05e524fdcd29d8040d64d0259e2d5109ca1998e567a3c1c1c68",
				HashResErr: nil,
			},
			expectedErr: nil,
		},
		{
			name: "positive duration",
			hashcash: &Hashcach{
				Bits:     5,
				Resource: "resource",
				Rand:     []byte{10},
				Date:     now.Add(50 * time.Second),
				Ext:      "resource\nsecret1648762844",
			},
			powOptions: []Options{
				WithChallengeExpDuration(time.Minute),
			},
			resource: "resource",
			hasherMock: mock.HasherMockParams{
				HashTimes:  1,
				HashReq:    gomock.Any(),
				HashRes:    "00000e89df98a05e524fdcd29d8040d64d0259e2d5109ca1998e567a3c1c1c68",
				HashResErr: nil,
			},
			expectedErr: nil,
		},
		{
			name: "wrong resource",
			hashcash: &Hashcach{
				Resource: "resource",
			},
			resource:    "resource2",
			expectedErr: ErrWrongResource,
		},
		{
			name: "challenge expired",
			hashcash: &Hashcach{
				Bits:     5,
				Resource: "resource",
				Rand:     []byte{10},
				Date:     time.Unix(1648762844, 0),
				Ext:      "resource\nsecret1648762844",
			},
			resource:    "resource",
			expectedErr: ErrChallengeExpired,
		},
		{
			name: "hasher error",
			hashcash: &Hashcach{
				Bits:     5,
				Resource: "resource",
				Rand:     []byte{10},
				Date:     time.Now(),
				Ext:      "resource\nsecret1648762844",
			},
			resource: "resource",
			hasherMock: mock.HasherMockParams{
				HashTimes:  1,
				HashReq:    gomock.Any(),
				HashRes:    "",
				HashResErr: hasherErr,
			},
			expectedErr: hasherErr,
		},
		{
			name: "wrong hash",
			hashcash: &Hashcach{
				Bits:     5,
				Resource: "resource",
				Rand:     []byte{10},
				Date:     time.Now(),
				Ext:      "resource\nsecret1648762844",
			},
			resource: "resource",
			hasherMock: mock.HasherMockParams{
				HashTimes:  1,
				HashReq:    gomock.Any(),
				HashRes:    "d59d15c9a1842bc4563897803799e94f1f242d7e7e8c618f047e068211543998",
				HashResErr: nil,
			},
			expectedErr: ErrWrongChallenge,
		},
		{
			name: "validate ext error",
			hashcash: &Hashcach{
				Bits:     5,
				Resource: "resource",
				Rand:     []byte{10},
				Date:     time.Now(),
				Ext:      "resource\nsecret1648762844",
			},
			powOptions: []Options{
				WithValidateExtFunc(func(h *Hashcach) error {
					return hasherErr
				}),
			},
			resource: "resource",
			hasherMock: mock.HasherMockParams{
				HashTimes:  1,
				HashReq:    gomock.Any(),
				HashRes:    "00000e89df98a05e524fdcd29d8040d64d0259e2d5109ca1998e567a3c1c1c68",
				HashResErr: nil,
			},
			expectedErr: hasherErr,
		},
		{
			name: "error duration",
			hashcash: &Hashcach{
				Bits:     5,
				Resource: "resource",
				Rand:     []byte{10},
				Date:     now.Add(-2 * time.Minute),
				Ext:      "resource\nsecret1648762844",
			},
			powOptions: []Options{
				WithChallengeExpDuration(time.Minute),
			},
			resource:    "resource",
			expectedErr: ErrChallengeExpired,
		},
	} {
		t.Run(tCase.name, func(t *testing.T) {
			var (
				a      = assert.New(t)
				hasher = tCase.hasherMock.NewHasher(ctrl)
				pow    = New(hasher, tCase.powOptions...)
			)

			err := pow.Verify(tCase.hashcash, tCase.resource)
			a.ErrorIs(err, tCase.expectedErr)
		})
	}
}
