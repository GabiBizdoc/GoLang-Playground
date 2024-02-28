package otp

import (
	"fmt"
	"strconv"
	"time"
)

type BinCodeWithExpiration struct {
	BinCode
	ExpireAt time.Time
}

func NewBinCodeWithExpiration(binCode BinCode, expireAt time.Time) *BinCodeWithExpiration {
	return &BinCodeWithExpiration{BinCode: binCode, ExpireAt: expireAt}
}

type BinCode int

func (b BinCode) HOTP() int {
	return int(b) % 1_000_000
}

func (b BinCode) TakeDigits(n int) int {
	if n <= 0 {
		return int(b)
	}
	return int(b) % pow(10, n)
}

func (b BinCode) TakeDigitsString(n int) string {
	return fmt.Sprintf("%0*d", n, b.TakeDigits(n))
}

func (b BinCode) HOTPString() string {
	return strconv.Itoa(b.HOTP())
}

func (b BinCode) Decimal() int {
	return int(b)
}

func (b BinCode) DecimalString() string {
	return strconv.Itoa(int(b))
}

func (b BinCode) HEX() string {
	return strconv.FormatInt(int64(b), 16)
}

// todo: use a faster algorithm for power
// pow calculates the power of a raised to the exponent n.
func pow(a, n int) int {
	p := a

	if n < 0 {
		panic("invalid exponent: n must be a non-negative integer")
	}

	if n == 0 {
		return 1
	}

	for n > 1 {
		n -= 1
		p *= a
	}

	return p
}
