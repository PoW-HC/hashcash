package pow

import (
	"time"
)

type options struct {
	validateExtFunc      ExtValidatorFunc
	challengeExpDuration time.Duration
}

type Options func(*options)

func WithValidateExtFunc(callback ExtValidatorFunc) Options {
	return func(pow *options) {
		pow.validateExtFunc = callback
	}
}
func WithChallengeExpDuration(callback time.Duration) Options {
	return func(pow *options) {
		pow.challengeExpDuration = callback
	}
}

func setDefaultOptions(o *options) {
	if o.challengeExpDuration == 0 {
		o.challengeExpDuration = defaultChallengeDuration
	}
}
