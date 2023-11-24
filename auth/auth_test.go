package auth_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/VitoNaychev/auth"
	"github.com/VitoNaychev/errorresponse"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("very-secret-key")
var expiresAt = time.Second

func DummyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
}

func TestJWTVerification(t *testing.T) {

	t.Run("returns Token on valid JWT ", func(t *testing.T) {
		jwtString, _ := auth.GenerateJWT(secretKey, expiresAt, 0)

		_, err := auth.VerifyJWT(jwtString, secretKey)
		if err != nil {
			t.Errorf("did not expect error, got %v", err)
		}
	})

	t.Run("returns error on invalid JWT", func(t *testing.T) {
		jwtString, _ := auth.GenerateJWT(secretKey, expiresAt, 0)

		jwtByteArr := []byte(jwtString)
		if jwtByteArr[10] == 'A' {
			jwtByteArr[10] = 'B'
		} else {
			jwtByteArr[10] = 'A'
		}
		jwtString = string(jwtByteArr)

		_, err := auth.VerifyJWT(jwtString, secretKey)
		if err == nil {
			t.Errorf("did not get error but expected one")
		}
	})
}

func TestAuthenticationMiddleware(t *testing.T) {
	dummyHandler := auth.AuthenticationMiddleware(DummyHandler, secretKey)

	t.Run("returns Unauthorized on missing JWT", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		response := httptest.NewRecorder()

		dummyHandler(response, request)

		assertStatus(t, response.Code, http.StatusUnauthorized)
		assertErrorResponse(t, response.Body, auth.ErrMissingToken)
	})

	t.Run("returns Unauthorized on invalid JWT", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		response := httptest.NewRecorder()

		request.Header.Add("Token", "thisIsAnInvalidJWT")

		dummyHandler(response, request)

		assertStatus(t, response.Code, http.StatusUnauthorized)

		var errorResponse errorresponse.ErrorResponse
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
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresAt)),
		})

		tokenString, _ := token.SignedString(secretKey)
		request.Header.Add("Token", tokenString)

		dummyHandler(response, request)

		assertStatus(t, response.Code, http.StatusUnauthorized)
		assertErrorResponse(t, response.Body, auth.ErrMissingSubject)
	})

	t.Run("returns Unauthorized on noninteger Subject in Token", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		response := httptest.NewRecorder()

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
			Subject:   "notAnIntegerSubject",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresAt)),
		})

		tokenString, _ := token.SignedString(secretKey)
		request.Header.Add("Token", tokenString)

		dummyHandler(response, request)

		assertStatus(t, response.Code, http.StatusUnauthorized)
		assertErrorResponse(t, response.Body, auth.ErrNonIntegerSubject)
	})

	t.Run("returns Token's Subject on valid JWT", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		response := httptest.NewRecorder()

		want := 10
		dummyJWT, _ := auth.GenerateJWT(secretKey, expiresAt, want)
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

func assertErrorResponse(t testing.TB, body io.Reader, expetedError error) {
	t.Helper()

	var errorResponse errorresponse.ErrorResponse
	json.NewDecoder(body).Decode(&errorResponse)

	if errorResponse.Message != expetedError.Error() {
		t.Errorf("got error %q want %q", errorResponse.Message, expetedError.Error())
	}
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Fatalf("got status %d, want %d", got, want)
	}
}
