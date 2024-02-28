package otp

import (
	"crypto/sha1"
	"encoding/hex"
	"strconv"
	"strings"
	"testing"
)

type TestCase struct {
	Key   string
	Count int

	// Hexadecimal HMAC-SHA-1(secret, count)
	ExpectedHash string
	TokenHOTP    string
	TokenDecimal string
	TokenHex     string
}

func TestGenerateTOTPHash(t *testing.T) {
	tc, err := getTestCases()
	if err != nil {
		t.Fatal(err)
	}

	for i, c := range tc {
		hmac := hex.EncodeToString(ComputeHMAC(sha1.New, []byte(c.Key), int64(c.Count)))
		failed := hmac != c.ExpectedHash

		//t.Logf("[%d] [%t] expected %s got %s", c.Count, !failed, c.ExpectedHash, hmac)

		if failed {
			t.Errorf("[%d] expected %s got %s", i, c.ExpectedHash, hmac)
		}
	}
}

func TestGenerateTOTP(t *testing.T) {
	tc, err := getTestCases()
	if err != nil {
		t.Fatal(err)
	}

	for i, c := range tc {
		val, err := GenerateCounterBasedOTP(sha1.New, []byte(c.Key), c.Count)
		if err != nil {
			t.Fatal(err)
		}

		//t.Logf("H: %10s V: %10d  HOTP: %10d", val.HEX(), val.Decimal(), val.HOTP())

		if val.HEX() != c.TokenHex {
			t.Fatalf("[%d] Failed to compute token HEX expected %s got %s\n", i, c.TokenHex, val.HEX())
		}

		if val.HOTPString() != c.TokenHOTP {
			t.Fatalf("[%d] Failed to compute token HOTP expected %s got %s\n", i, c.TokenHOTP, val.HOTPString())
		}

		if val.DecimalString() != c.TokenDecimal {
			t.Fatalf("[%d] Failed to compute token HEX expected %s got %s\n", i, c.TokenDecimal, val.DecimalString())
		}
	}
}

func getTestCases() ([]*TestCase, error) {
	tc := make([]*TestCase, 0)

	const testKey = "12345678901234567890"
	const testCases = `
   Count    Hexadecimal HMAC-SHA-1(secret, count)
   0        cc93cf18508d94934c64b65d8ba7667fb7cde4b0
   1        75a48a19d4cbe100644e8ac1397eea747a2d33ab
   2        0bacb7fa082fef30782211938bc1c5e70416ff44
   3        66c28227d03a2d5529262ff016a1e6ef76557ece
   4        a904c900a64b35909874b33e61c5938a8e15ed1c
   5        a37e783d7b7233c083d4f62926c7a25f238d0316
   6        bc9cd28561042c83f219324d3c607256c03272ae
   7        a4fb960c0bc06e1eabb804e5b397cdc4b45596fa
   8        1b3c89f65e6c9e883012052823443f048b4332db
   9        1637409809a679dc698207310c8c7fc07290d9e5
`

	const testResults = `
   Count    Hexadecimal    Decimal        HOTP
   0        4c93cf18       1284755224     755224
   1        41397eea       1094287082     287082
   2         82fef30        137359152     359152
   3        66ef7655       1726969429     969429
   4        61c5938a       1640338314     338314
   5        33c083d4        868254676     254676
   6        7256c032       1918287922     287922
   7         4e5b397         82162583     162583
   8        2823443f        673399871     399871
   9        2679dc69        645520489     520489
`

	firstSkipped := false
	for _, line := range strings.Split(strings.TrimSpace(testCases), "\n") {
		line := strings.TrimSpace(line)
		if len(line) < 0 {
			continue
		}

		if !firstSkipped {
			firstSkipped = true
			continue
		}

		fields := strings.Fields(line)
		index, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, err
		}

		tc = append(tc, &TestCase{
			Key:          testKey,
			Count:        index,
			ExpectedHash: fields[1],
		})
	}

	firstSkipped = false
	for _, line := range strings.Split(strings.TrimSpace(testResults), "\n") {
		line := strings.TrimSpace(line)
		if len(line) < 0 {
			continue
		}

		if !firstSkipped {
			firstSkipped = true
			continue
		}

		fields := strings.Fields(line)
		index, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, err
		}

		c := tc[index]
		c.TokenHex = fields[1]
		c.TokenDecimal = fields[2]
		c.TokenHOTP = fields[3]
	}

	return tc, nil
}
