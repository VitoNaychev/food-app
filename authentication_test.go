package bt_customer_svc

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
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

func TestAuthenticationMiddleware(t *testing.T) {
	godotenv.Load("test.env")
	secretKey := []byte(os.Getenv("SECRET"))
	dummyHandler := AuthenticationMiddleware(DummyHandler, secretKey)

	t.Run("returns Unauthorized on missing JWT", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		response := httptest.NewRecorder()

		dummyHandler(response, request)

		assertStatus(t, response.Code, http.StatusUnauthorized)

		var errorResponse ErrorResponse
		json.NewDecoder(response.Body).Decode(&errorResponse)
		assertErrorResponse(t, errorResponse, ErrMissingToken)
	})

	t.Run("returns Unauthorized on invalid JWT", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		response := httptest.NewRecorder()

		request.Header.Add("Token", "thisIsAnInvalidJWT")

		dummyHandler(response, request)

		assertStatus(t, response.Code, http.StatusUnauthorized)

		var errorResponse ErrorResponse
		decoder := json.NewDecoder(response.Body)
		decoder.DisallowUnknownFields()
		decoder.Decode(&errorResponse)
		if errorResponse.Message == "" {
			t.Errorf("expected error message but did not get one")
		}
	})

	t.Run("returns Unauthorized on missing Subject in Token", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		response := httptest.NewRecorder()

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second)),
		})

		tokenString, _ := token.SignedString(secretKey)
		request.Header.Add("Token", tokenString)

		dummyHandler(response, request)

		assertStatus(t, response.Code, http.StatusUnauthorized)

		var errorResponse ErrorResponse
		json.NewDecoder(response.Body).Decode(&errorResponse)
		assertErrorResponse(t, errorResponse, ErrMissingSubject)
	})

	t.Run("returns Unauthorized on noninteger Subject in Token", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		response := httptest.NewRecorder()

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
			Subject:   "notAnIntegerSubject",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second)),
		})

		tokenString, _ := token.SignedString(secretKey)
		request.Header.Add("Token", tokenString)

		dummyHandler(response, request)

		assertStatus(t, response.Code, http.StatusUnauthorized)

		var errorResponse ErrorResponse
		json.NewDecoder(response.Body).Decode(&errorResponse)
		assertErrorResponse(t, errorResponse, ErrNonIntegerSubject)
	})

	t.Run("returns Token's Subject on valid JWT", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		response := httptest.NewRecorder()

		want := 10
		dummyJWT, _ := GenerateJWT(secretKey, time.Now().Add(time.Second), want)
		request.Header.Add("Token", dummyJWT)

		dummyHandler(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)

		subject := request.Header["Subject"]
		if subject == nil {
			t.Fatalf("did not get Subject in request header")
		}

		got, err := strconv.Atoi(subject[0])
		if err != nil {
			t.Fatalf("expected integer Subject, got %q", subject[0])
		}

		if got != want {
			t.Errorf("got %d want %d", got, want)
		}
	})
}
