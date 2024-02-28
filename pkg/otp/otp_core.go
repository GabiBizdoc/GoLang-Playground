package otp

import (
	"crypto/hmac"
	"encoding/binary"
	"hash"
	"time"
)

func ComputeHMAC(h func() hash.Hash, key []byte, counter int64) []byte {
	counterBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(counterBytes, uint64(counter))

	hmacSha := hmac.New(h, key)
	hmacSha.Write(counterBytes)
	hmacResult := hmacSha.Sum(nil)

	return hmacResult
}

func ExtractBinCode(hmacResult []byte) BinCode {
	offset := int(hmacResult[len(hmacResult)-1]) & 0xf
	binCode := int(hmacResult[offset]&0x7f)<<24 |
		int(hmacResult[offset+1]&0xff)<<16 |
		int(hmacResult[offset+2]&0xff)<<8 |
		int(hmacResult[offset+3]&0xff)

	return BinCode(binCode)
}

func GetCurrentTimeBasedCounter(interval int) (int64, time.Time) {
	x := int64(interval) * 1000
	ct := time.Now()

	unixTime := ct.UTC().UnixMilli()
	cnt := unixTime / x

	expireAt := (cnt + 1) * x
	return cnt, time.UnixMilli(expireAt)
}
