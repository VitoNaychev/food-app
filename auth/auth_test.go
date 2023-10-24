package auth

import (
	"testing"
	"time"
)

func TestJWTVerification(t *testing.T) {
	secretKey := []byte("mySecretKey")

	t.Run("returns Token on valid JWT ", func(t *testing.T) {
		jwtString, _ := GenerateJWT(secretKey, time.Now().Add(time.Second), 0)

		_, err := VerifyJWT(jwtString, secretKey)
		if err != nil {
			t.Errorf("did not expect error, got %v", err)
		}
	})

	t.Run("returns error on invalid JWT", func(t *testing.T) {
		jwtString, _ := GenerateJWT(secretKey, time.Now().Add(time.Second), 0)

		jwtByteArr := []byte(jwtString)
		if jwtByteArr[10] == 'A' {
			jwtByteArr[10] = 'B'
		} else {
			jwtByteArr[10] = 'A'
		}
		jwtString = string(jwtByteArr)

		_, err := VerifyJWT(jwtString, secretKey)
		if err == nil {
			t.Errorf("did not get error but expected one")
		}
	})
}
