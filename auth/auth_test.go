package auth_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/VitoNaychev/food-app/auth"
	"github.com/VitoNaychev/food-app/httperrors"
	"github.com/VitoNaychev/food-app/msgtypes"
	"github.com/VitoNaychev/food-app/testutil"
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

type DummyVerifier struct {
	shouldError bool
	shouldFail  bool
}

func (d *DummyVerifier) DoesSubjectExist(id int) (bool, error) {
	if d.shouldError {
		return false, auth.ErrMissingSubject
	}

	if d.shouldFail {
		return false, nil
	}

	return true, nil
}

func TestAuthenticationMW(t *testing.T) {
	dummyVerifier := &DummyVerifier{false, false}
	dummyHandler := auth.AuthenticationMW(DummyHandler, dummyVerifier, secretKey)

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

		var errorResponse httperrors.ErrorResponse
		decoder := json.NewDecoder(response.Body)
		decoder.DisallowUnknownFields()
		decoder.Decode(&errorResponse)
		if errorResponse.Error == "" {
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

	t.Run("returns Internal Server Error on error from verifier", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		response := httptest.NewRecorder()

		want := 10
		dummyJWT, _ := auth.GenerateJWT(secretKey, expiresAt, want)
		request.Header.Add("Token", dummyJWT)

		dummyVerifier.shouldError = true
		dummyHandler(response, request)
		dummyVerifier.shouldError = false

		assertStatus(t, response.Code, http.StatusInternalServerError)
	})

	t.Run("returns Not Found when verifier returns false", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		response := httptest.NewRecorder()

		want := 10
		dummyJWT, _ := auth.GenerateJWT(secretKey, expiresAt, want)
		request.Header.Add("Token", dummyJWT)

		dummyVerifier.shouldFail = true
		dummyHandler(response, request)
		dummyVerifier.shouldFail = false

		assertStatus(t, response.Code, http.StatusNotFound)
		assertErrorResponse(t, response.Body, auth.ErrSubjectNotFound)
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

func StubVerifyJWT(jwt string) (msgtypes.AuthResponse, error) {
	if jwt == "invalidJWT" {
		return msgtypes.AuthResponse{Status: msgtypes.INVALID, ID: 0}, nil
	} else if jwt == "10" {
		return msgtypes.AuthResponse{Status: msgtypes.NOT_FOUND, ID: 0}, nil
	} else {
		id, _ := strconv.Atoi(jwt)
		return msgtypes.AuthResponse{Status: msgtypes.OK, ID: id}, nil
	}
}

func TestAuthMiddleware(t *testing.T) {
	handler := auth.RemoteAuthenticationMW(DummyHandler, StubVerifyJWT)

	t.Run("returns Unauthorized on missing JWT", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		handler(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusUnauthorized)
	})

	t.Run("returns Unauthorized on invalid JWT", func(t *testing.T) {
		customerJWT := "invalidJWT"
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		request.Header.Add("Token", customerJWT)

		response := httptest.NewRecorder()

		handler(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusUnauthorized)
		testutil.AssertErrorResponse(t, response.Body, auth.ErrInvalidToken)
	})

	t.Run("returns Not Found on nonexistent customer", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		request.Header.Add("Token", strconv.Itoa(10))
		response := httptest.NewRecorder()

		handler(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusNotFound)
		testutil.AssertErrorResponse(t, response.Body, auth.ErrSubjectNotFound)
	})

	t.Run("returns Accepted on authentic customer", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		request.Header.Add("Token", strconv.Itoa(1))
		response := httptest.NewRecorder()

		handler(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusAccepted)
	})
}

func assertErrorResponse(t testing.TB, body io.Reader, expetedError error) {
	t.Helper()

	var errorResponse httperrors.ErrorResponse
	json.NewDecoder(body).Decode(&errorResponse)

	if errorResponse.Error != expetedError.Error() {
		t.Errorf("got error %q want %q", errorResponse.Error, expetedError.Error())
	}
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Fatalf("got status %d, want %d", got, want)
	}
}
