# OTP Generator for Go
This Go package provides a simple implementation for generating One-Time Passwords (OTP) based on the 
HMAC-based One-Time Password Algorithm (HOTP) and
Time-Based One-Time Password Algorithm (TOTP) 
as described in [RFC 4226](https://datatracker.ietf.org/doc/html/rfc4226) 
and [RFC 6238](https://datatracker.ietf.org/doc/html/rfc6238).

## Usage

```go
package main

import (
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/GabiBizdoc/golang-playground/pgk/otp"
)

func main() {
	secret := "<your base32 key"
	code, _ := otp.GenerateTOTP(secret)
	token := code.TakeDigitsString(6)
	expireIn := time.Until(code.ExpireAt)
	
	fmt.Println(token, expireIn)
	
	var rawSecret []byte
	bincode, _ := otp.GenerateTimeBasedOTP(sha1.New, rawSecret, 30)
	token = bincode.TakeDigitsString(6)
	
	fmt.Println(token)
}
```