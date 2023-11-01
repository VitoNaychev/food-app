package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/VitoNaychev/bt-customer-svc/handlers"
	"github.com/VitoNaychev/bt-customer-svc/testutil"
	"github.com/golang-jwt/jwt/v5"
)

func DummyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
}

var testEnv handlers.TestEnv

func TestMain(m *testing.M) {
	testEnv = handlers.LoadTestEnviornment()

	code := m.Run()
	os.Exit(code)
}

func TestJWTVerification(t *testing.T) {
	t.Run("returns Token on valid JWT ", func(t *testing.T) {
		jwtString, _ := GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, 0)

		_, err := VerifyJWT(jwtString, testEnv.SecretKey)
		if err != nil {
			t.Errorf("did not expect error, got %v", err)
		}
	})

	t.Run("returns error on invalid JWT", func(t *testing.T) {
		jwtString, _ := GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, 0)

		jwtByteArr := []byte(jwtString)
		if jwtByteArr[10] == 'A' {
			jwtByteArr[10] = 'B'
		} else {
			jwtByteArr[10] = 'A'
		}
		jwtString = string(jwtByteArr)

		_, err := VerifyJWT(jwtString, testEnv.SecretKey)
		if err == nil {
			t.Errorf("did not get error but expected one")
		}
	})
}

func TestAuthenticationMiddleware(t *testing.T) {
	dummyHandler := AuthenticationMiddleware(DummyHandler, testEnv.SecretKey)

	t.Run("returns Unauthorized on missing JWT", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		response := httptest.NewRecorder()

		dummyHandler(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusUnauthorized)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrMissingToken)
	})

	t.Run("returns Unauthorized on invalid JWT", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		response := httptest.NewRecorder()

		request.Header.Add("Token", "thisIsAnInvalidJWT")

		dummyHandler(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusUnauthorized)

		var errorResponse handlers.ErrorResponse
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
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(testEnv.ExpiresAt)),
		})

		tokenString, _ := token.SignedString(testEnv.SecretKey)
		request.Header.Add("Token", tokenString)

		dummyHandler(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusUnauthorized)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrMissingSubject)
	})

	t.Run("returns Unauthorized on noninteger Subject in Token", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		response := httptest.NewRecorder()

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
			Subject:   "notAnIntegerSubject",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(testEnv.ExpiresAt)),
		})

		tokenString, _ := token.SignedString(testEnv.SecretKey)
		request.Header.Add("Token", tokenString)

		dummyHandler(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusUnauthorized)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrNonIntegerSubject)
	})

	t.Run("returns Token's Subject on valid JWT", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		response := httptest.NewRecorder()

		want := 10
		dummyJWT, _ := GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, want)
		request.Header.Add("Token", dummyJWT)

		dummyHandler(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusAccepted)

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
