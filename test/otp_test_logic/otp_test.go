package otp_test_logic

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type totpCase struct {
	T        time.Time
	HashFunc func() hash.Hash
	Digit    int
	TimeStep int
	Code     string
	Secret   string
}

var (
	// https://datatracker.ietf.org/doc/html/rfc6238#appendix-B
	totpCases = []totpCase{
		{time.Unix(59, 0), sha1.New, 8, 30, "94287082", "12345678901234567890"},
		{time.Unix(1111111109, 0), sha1.New, 8, 30, "07081804", "12345678901234567890"},
		{time.Unix(1111111111, 0), sha1.New, 8, 30, "14050471", "12345678901234567890"},
		{time.Unix(1234567890, 0), sha1.New, 8, 30, "89005924", "12345678901234567890"},
		{time.Unix(2000000000, 0), sha1.New, 8, 30, "69279037", "12345678901234567890"},
		{time.Unix(20000000000, 0), sha1.New, 8, 30, "65353130", "12345678901234567890"},
		{time.Unix(59, 0), sha256.New, 8, 30, "46119246", "12345678901234567890123456789012"},
		{time.Unix(1111111109, 0), sha256.New, 8, 30, "68084774", "12345678901234567890123456789012"},
		{time.Unix(1111111111, 0), sha256.New, 8, 30, "67062674", "12345678901234567890123456789012"},
		{time.Unix(1234567890, 0), sha256.New, 8, 30, "91819424", "12345678901234567890123456789012"},
		{time.Unix(2000000000, 0), sha256.New, 8, 30, "90698825", "12345678901234567890123456789012"},
		{time.Unix(20000000000, 0), sha256.New, 8, 30, "77737706", "12345678901234567890123456789012"},
		{time.Unix(59, 0), sha512.New, 8, 30, "90693936", "1234567890123456789012345678901234567890123456789012345678901234"},
		{time.Unix(1111111109, 0), sha512.New, 8, 30, "25091201", "1234567890123456789012345678901234567890123456789012345678901234"},
		{time.Unix(1111111111, 0), sha512.New, 8, 30, "99943326", "1234567890123456789012345678901234567890123456789012345678901234"},
		{time.Unix(1234567890, 0), sha512.New, 8, 30, "93441116", "1234567890123456789012345678901234567890123456789012345678901234"},
		{time.Unix(2000000000, 0), sha512.New, 8, 30, "38618901", "1234567890123456789012345678901234567890123456789012345678901234"},
		{time.Unix(20000000000, 0), sha512.New, 8, 30, "47863826", "1234567890123456789012345678901234567890123456789012345678901234"},
	}
)

// func Test_OTP_Generate(t *testing.T) {
// 	for _, tt := range totpCases {
// 		otp := NewOTP([]byte(tt.Secret))
// 		code := otp.Generate(
// 			WithDigit(tt.Digit),
// 			WithHashFunc(tt.HashFunc),
// 		)
// 		assert.Equal(t, tt.Code, code.code)
// 	}
// }

// func Test_OTP_Validate(t *testing.T) {
// 	for _, tt := range totpCases {
// 		otp := NewOTP([]byte(tt.Secret))
// 		assert.Equal(t, true, otp.Validate(
// 			"adnan",
// 			tt.Code,
// 			WithDigit(tt.Digit),
// 			WithHashFunc(tt.HashFunc),
// 		))
// 	}
// }

// func Test_OTP_Timeout(t *testing.T) {
// 	otp := NewOTP([]byte("12345678901234567890"))
// 	code := otp.Generate(
// 		WithDigit(6),
// 		WithHashFunc(sha1.New),
// 		WithTimeStep(1),
// 	)

// 	time.Sleep(3 * time.Second)
// 	assert.Equal(t, false, otp.Validate(
// 		code.code,
// 		WithDigit(6),
// 		WithHashFunc(sha1.New),
// 		WithTimeStep(1),
// 	))
// }

func Test_OTP_Not_Timeout(t *testing.T) {
	key := fmt.Sprintf("%v%v", strings.ToUpper("adnan"), "081215775661")
	otp := NewOTP([]byte(key))
	code := otp.Generate(
		WithDigit(6),
		WithHashFunc(sha1.New),
		WithTimeStep(10),
	)

	assert.Equal(t, true, otp.Validate(
		code.code,
		WithDigit(6),
		WithHashFunc(sha1.New),
		WithTimeStep(10),
	))
}

func Test_OTP_Invalid_Credential(t *testing.T) {
	key := fmt.Sprintf("%v%v", strings.ToUpper("adnan"), "081215775661")
	otp := NewOTP([]byte(key))
	code := otp.Generate(
		WithDigit(6),
		WithHashFunc(sha1.New),
		WithTimeStep(10),
	)

	key = fmt.Sprintf("%v%v", strings.ToUpper("wrong"), "081215775661")
	otp = NewOTP([]byte(key))
	assert.Equal(t, false, otp.Validate(
		code.code,
		WithDigit(6),
		WithHashFunc(sha1.New),
		WithTimeStep(10),
	))
}

// func Test_OTP_Different_Timestep(t *testing.T) {
// 	otp := NewOTP([]byte("12345678901234567890"))
// 	code := otp.Generate(
// 		WithDigit(6),
// 		WithHashFunc(sha1.New),
// 		WithTimeStep(6),
// 	)

// 	assert.Equal(t, false, otp.Validate(
// 		code.code,
// 		WithDigit(6),
// 		WithHashFunc(sha1.New),
// 		WithTimeStep(5),
// 	))
// }

// func Test_OTP_Different_Digit(t *testing.T) {
// 	otp := NewOTP([]byte("12345678901234567890"))
// 	code := otp.Generate(
// 		WithDigit(10),
// 		WithHashFunc(sha1.New),
// 		WithTimeStep(5),
// 	)

// 	assert.Equal(t, false, otp.Validate(
// 		code.code,
// 		WithDigit(6),
// 		WithHashFunc(sha1.New),
// 		WithTimeStep(5),
// 	))
// }

// func Test_OTP_Different_Hash(t *testing.T) {
// 	otp := NewOTP([]byte("12345678901234567890"))
// 	code := otp.Generate(
// 		WithDigit(10),
// 		WithHashFunc(sha1.New),
// 		WithTimeStep(5),
// 	)

// 	assert.Equal(t, false, otp.Validate(
// 		code.code,
// 		WithDigit(6),
// 		WithHashFunc(sha256.New),
// 		WithTimeStep(5),
// 	))
// }
