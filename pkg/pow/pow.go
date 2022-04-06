package pow

import (
	"context"
	"fmt"
	"time"

	"github.com/PoW-HC/hashcash/pkg/hash"
)

const (
	zero rune = 48 // ASCII code for number zero

	defaultChallengeDuration = 120 * time.Second
	maxHashSize              = 1 << 7 // max length of sha512
)

var (
	ErrMaxIterationsExceeded = fmt.Errorf("max iterations exceeded")
	ErrWrongResource         = fmt.Errorf("wrong resource")
	ErrChallengeExpired      = fmt.Errorf("challenge expired")
	ErrWrongChallenge        = fmt.Errorf("wrong challenge")
)

// POW proof of work class
type POW struct {
	s     hash.Hasher
	zeros []rune

	options options
}

// New constructor
func New(s hash.Hasher, opts ...Options) *POW {
	p := &POW{
		s:     s,
		zeros: make([]rune, maxHashSize),
	}

	for i := range p.zeros {
		p.zeros[i] = zero
	}

	for i := range opts {
		opts[i](&p.options)
	}

	setDefaultOptions(&p.options)

	return p
}

// Compute time waster. Do all useless load.
func (p *POW) Compute(ctx context.Context, h *Hashcach, max int64) (*Hashcach, error) {
	if max > 0 {
		for h.Counter <= max {
			if err := ctx.Err(); err != nil {
				break
			}

			hashString, err := p.s.Hash(h.String())
			if err != nil {
				return nil, fmt.Errorf("calculate pow hash sum error: %w", err)
			}

			if isHashCorrect(hashString, p.zeros, int(h.Bits)) {
				return h, nil
			}

			h.Counter++
		}
	}

	return nil, ErrMaxIterationsExceeded
}

// Verify that hashcash correct and provided by server
func (p *POW) Verify(h *Hashcach, resource string) error {
	if h == nil || h.Resource != resource {
		return ErrWrongResource
	}

	if h.Date.Add(p.options.challengeExpDuration).Before(time.Now()) {
		return ErrChallengeExpired
	}

	hashString, err := p.s.Hash(h.String())
	if err != nil {
		return fmt.Errorf("calculate pow hash sum error: %w", err)
	}

	if !isHashCorrect(hashString, p.zeros, int(h.Bits)) {
		return ErrWrongChallenge
	}

	if p.options.validateExtFunc != nil {
		err = p.options.validateExtFunc(h)
		if err != nil {
			return fmt.Errorf("validation extension error: %w", err)
		}
	}

	return nil
}

func isHashCorrect(hash string, zeroHash []rune, zerosCount int) bool {
	if zerosCount > len(hash) || zerosCount > len(zeroHash) {
		return false
	}

	return hash[:zerosCount] == string(zeroHash[:zerosCount])
}
