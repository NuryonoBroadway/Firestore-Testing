package otp_test_logic

import (
	"crypto/sha1"
	"hash"
	"time"
)

type Options struct {
	Digit     int
	HashFunc  func() hash.Hash
	TimeStep  int
	StartTime int64
}

type CodeInformation struct {
	code    string
	timeout time.Time
}

func NewOptions() *Options {
	return &Options{
		Digit:     6,        // how many digit the code is
		HashFunc:  sha1.New, // default otp hash
		TimeStep:  150,      // default otp timestep
		StartTime: 0,
	}
}

type Option func(*Options)

// WithDigit
// The default digit = 6
func WithDigit(digit int) Option {
	return func(options *Options) {
		options.Digit = digit
	}
}

// WithHashFunc
// The default hash func is sha1.New
// You can use other hash func such as: sha256.New, sha512.New, md5.New
func WithHashFunc(f func() hash.Hash) Option {
	return func(options *Options) {
		options.HashFunc = f
	}
}

// WithTimeStep
// The default time step is 30s.
func WithTimeStep(step int) Option {
	return func(options *Options) {
		options.TimeStep = step
	}
}

// WithStartTime
// The default start time = 0
// The UNIX time to start counting time steps
func WithStartTime(start int64) Option {
	return func(options *Options) {
		options.StartTime = start
	}
}
