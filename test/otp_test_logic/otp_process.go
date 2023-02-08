package otp_test_logic

import (
	"crypto/hmac"
	"crypto/subtle"
	"encoding/binary"
	"fmt"
	"math"
)

// Generate
// https://datatracker.ietf.org/doc/html/rfc4226#section-5.3
func (otp OTPConfig) generate(in uint64, opts ...Option) string {
	options := NewOptions()
	for _, opt := range opts {
		opt(options)
	}

	// 8-byte counter value, the moving factor.
	counter := make([]byte, 8)
	binary.BigEndian.PutUint64(counter, uint64(in))

	h := hmac.New(options.HashFunc, otp.Key)
	_, _ = h.Write(counter)
	hash := h.Sum(nil)

	offset := hash[len(hash)-1] & 0xf
	// The dynamic binary code is a 31-bit, unsigned, big-endian integer
	binary := int(hash[offset]&0x7f)<<24 |
		int(hash[offset+1]&0xff)<<16 |
		int(hash[offset+2]&0xff)<<8 |
		int(hash[offset+3])&0xff

	remainder := int64(binary) % int64(math.Pow10(options.Digit))
	code := fmt.Sprintf(fmt.Sprintf("%%0%dd", options.Digit), remainder)
	return code
}

// Validate
// return true if the code matched the generated code
func (otp OTPConfig) validate(in uint64, code string, opts ...Option) bool {
	res := otp.generate(in, opts...)
	return subtle.ConstantTimeCompare([]byte(code), []byte(res)) == 1
}
