package main

import (
	"testing"
	"time"
)

func TestJWTVerification(t *testing.T) {
	secretKey := []byte("mySecretKey")

	t.Run("test verify with valid JWT", func(t *testing.T) {
		jwtString, _ := generateJWT(secretKey, time.Now().Add(time.Second), 0)

		_, err := verifyJWT(jwtString, secretKey)
		if err != nil {
			t.Errorf("did not expect error, got %v", err)
		}
	})

	t.Run("test verify with invalid JWT", func(t *testing.T) {
		jwtString, _ := generateJWT(secretKey, time.Now().Add(time.Second), 0)

		jwtByteArr := []byte(jwtString)
		if jwtByteArr[10] == 'A' {
			jwtByteArr[10] = 'B'
		} else {
			jwtByteArr[10] = 'A'
		}
		jwtString = string(jwtByteArr)

		_, err := verifyJWT(jwtString, secretKey)
		if err == nil {
			t.Errorf("did not get error but expected one")
		}
	})
}
