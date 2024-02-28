package otp

import (
	"crypto/sha1"
	"encoding/base32"
	"hash"
)

// GenerateTOTP produces a One-Time Password (OTP) code using a base32 key.
// It uses a 30s interval and SHA-1.
func GenerateTOTP(secret string) (result *BinCodeWithExpiration, err error) {
	k, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return result, err
	}
	return GenerateTimeBasedOTP(sha1.New, k, 30)
}

// GenerateTOTPFromBase32Key produces a One-Time Password (OTP) code using a base32 key.
//
// h:      Hash function to be used for hashing. (e.g., sha1.New)
// secret: Base32 secret key for OTP generation.
// period: Time interval in seconds for which the OTP code is valid.
func GenerateTOTPFromBase32Key(h func() hash.Hash, secret string, period int) (result *BinCodeWithExpiration, err error) {
	k, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return result, err
	}
	return GenerateTimeBasedOTP(h, k, period)
}

// GenerateTimeBasedOTP produces a One-Time Password (OTP) code.
// Parameters:
//
//	h:      Algorithm for OTP generation (e.g., sha1.New, sha256.New, sha521.New).
//	key:    Secret key for OTP generation.
//	period: Time interval in seconds for which the OTP code is valid.
func GenerateTimeBasedOTP(h func() hash.Hash, key []byte, period int) (result *BinCodeWithExpiration, err error) {
	cnt, expireAt := GetCurrentTimeBasedCounter(period)
	hmacResult := ComputeHMAC(h, key, cnt)

	result = &BinCodeWithExpiration{
		BinCode:  ExtractBinCode(hmacResult),
		ExpireAt: expireAt,
	}

	return result, nil
}

// GenerateCounterBasedOTP generates a One-Time Password (OTP) using a counter-based algorithm.
func GenerateCounterBasedOTP(h func() hash.Hash, key []byte, counter int) (result BinCode, err error) {
	hmacResult := ComputeHMAC(h, key, int64(counter))
	return ExtractBinCode(hmacResult), nil
}
