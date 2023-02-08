package otp_test_logic

import (
	"time"
)

type OTPConfig struct {
	Key []byte
}

func NewOTP(key []byte) OTPConfig {
	return OTPConfig{Key: key}
}

// Generate
// https://datatracker.ietf.org/doc/html/rfc4226#section-5.3
func (otp OTPConfig) Generate(opts ...Option) CodeInformation {
	options := NewOptions()
	for _, opt := range opts {
		opt(options)
	}

	t := time.Now()
	counter := (t.Unix() - options.StartTime) / int64(options.TimeStep)
	code := NewOTP(otp.Key).generate(uint64(counter), opts...)
	return CodeInformation{
		code:    code,
		timeout: time.Now().Add(time.Duration(options.TimeStep)),
	}
}

// Validate
// return true if the code matched the generated code
func (otp OTPConfig) Validate(code string, opts ...Option) bool {
	options := NewOptions()
	for _, opt := range opts {
		opt(options)
	}

	t := time.Now()
	counter := (t.Unix() - options.StartTime) / int64(options.TimeStep)
	return otp.validate(uint64(counter), code, opts...)
}
