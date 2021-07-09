package auth_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fatdes/reap_backend_challenge/auth"
	mock_auth "github.com/fatdes/reap_backend_challenge/mock/auth"
	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	. "github.com/golang/mock/gomock"
)

func Test_Login_REST_Errors(t *testing.T) {
	nothingToMock := func(_ *mock_auth.MockLoginDB, _ *mock_auth.MockLoginTokenGenerator) {}

	tests := map[string]struct {
		setupMock          func(*mock_auth.MockLoginDB, *mock_auth.MockLoginTokenGenerator)
		request            func() *http.Request
		expectedStatusCode int
		expectedBody       string
	}{
		"Empty Username": {
			nothingToMock,
			func() *http.Request {
				req := httptest.NewRequest("POST", "/login", strings.NewReader(`{ "username": "" }`))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			http.StatusBadRequest,
			`{"ErrorMessage":"\"username\" (length: 0) should be between 4 - 20 characters"}`,
		},
		"Empty Username After Trim": {
			nothingToMock,
			func() *http.Request {
				req := httptest.NewRequest("POST", "/login", strings.NewReader(`{ "username": "    " }`))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			http.StatusBadRequest,
			`{"ErrorMessage":"\"username\" (length: 0) should be between 4 - 20 characters"}`,
		},
		"Username Length < 4": {
			nothingToMock,
			func() *http.Request {
				req := httptest.NewRequest("POST", "/login", strings.NewReader(`{ "username": "123" }`))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			http.StatusBadRequest,
			`{"ErrorMessage":"\"username\" (length: 3) should be between 4 - 20 characters"}`,
		},
		"Username Length > 20": {
			nothingToMock,
			func() *http.Request {
				req := httptest.NewRequest("POST", "/login", strings.NewReader(`{ "username": "123456789012345678901" }`))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			http.StatusBadRequest,
			`{"ErrorMessage":"\"username\" (length: 21) should be between 4 - 20 characters"}`,
		},
		"Empty Password": {
			nothingToMock,
			func() *http.Request {
				req := httptest.NewRequest("POST", "/login", strings.NewReader(`{ "username": "1234", "password": "" }`))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			http.StatusBadRequest,
			`{"ErrorMessage":"\"password\" (length: 0) should be between 4 - 20 characters"}`,
		},
		"Empty Password After Trim": {
			nothingToMock,
			func() *http.Request {
				req := httptest.NewRequest("POST", "/login", strings.NewReader(`{ "username": "1234", "password": "   " }`))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			http.StatusBadRequest,
			`{"ErrorMessage":"\"password\" (length: 0) should be between 4 - 20 characters"}`,
		},
		"Password Length < 4": {
			nothingToMock,
			func() *http.Request {
				req := httptest.NewRequest("POST", "/login", strings.NewReader(`{ "username": "1234", "password": "abc" }`))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			http.StatusBadRequest,
			`{"ErrorMessage":"\"password\" (length: 3) should be between 4 - 20 characters"}`,
		},
		"Password Length > 20": {
			nothingToMock,
			func() *http.Request {
				req := httptest.NewRequest("POST", "/login", strings.NewReader(`{ "username": "1234", "password": "123456789012345678901" }`))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			http.StatusBadRequest,
			`{"ErrorMessage":"\"password\" (length: 21) should be between 4 - 20 characters"}`,
		},
		"No Request Body": {
			nothingToMock,
			func() *http.Request {
				req := httptest.NewRequest("POST", "/login", strings.NewReader(``))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			http.StatusBadRequest,
			`{"ErrorMessage":"request body should be in json {\"username\":\"...\",\"password\":\"...\"}"}`,
		},
		"Password Not Match": {
			func(db *mock_auth.MockLoginDB, tokenGenerator *mock_auth.MockLoginTokenGenerator) {
				db.EXPECT().Exists(Eq("1234")).Return(true, nil)
				db.EXPECT().Login(Eq("1234"), Any()).Return(false, nil)
			},
			func() *http.Request {
				req := httptest.NewRequest("POST", "/login", strings.NewReader(`{ "username": "1234", "password": "abcde" }`))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			http.StatusUnauthorized,
			`{"ErrorMessage":"username or password not match"}`,
		},
		"DB Error When Check Exists": {
			func(db *mock_auth.MockLoginDB, tokenGenerator *mock_auth.MockLoginTokenGenerator) {
				db.EXPECT().Exists(Eq("1234")).Return(false, errors.New("OPS..."))
			},
			func() *http.Request {
				req := httptest.NewRequest("POST", "/login", strings.NewReader(`{ "username": "1234", "password": "abcde" }`))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			http.StatusInternalServerError,
			`{"ErrorMessage":"We will fix the issues ASAP. Please try again later"}`,
		},
		"DB Error When Create": {
			func(db *mock_auth.MockLoginDB, tokenGenerator *mock_auth.MockLoginTokenGenerator) {
				db.EXPECT().Exists(Eq("1234")).Return(false, nil)
				db.EXPECT().Create(Eq("1234"), Any()).Return(errors.New("OPS..."))
			},
			func() *http.Request {
				req := httptest.NewRequest("POST", "/login", strings.NewReader(`{ "username": "1234", "password": "abcde" }`))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			http.StatusInternalServerError,
			`{"ErrorMessage":"We will fix the issues ASAP. Please try again later"}`,
		},
		"DB Error When Login": {
			func(db *mock_auth.MockLoginDB, tokenGenerator *mock_auth.MockLoginTokenGenerator) {
				db.EXPECT().Exists(Eq("1234")).Return(true, nil)
				db.EXPECT().Login(Eq("1234"), Any()).Return(false, errors.New("OPS..."))
			},
			func() *http.Request {
				req := httptest.NewRequest("POST", "/login", strings.NewReader(`{ "username": "1234", "password": "abcde" }`))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			http.StatusInternalServerError,
			`{"ErrorMessage":"We will fix the issues ASAP. Please try again later"}`,
		},
	}

	for name, test := range tests {
		mockctrl := gomock.NewController(t)
		db := mock_auth.NewMockLoginDB(mockctrl)
		tokenGenerator := mock_auth.NewMockLoginTokenGenerator(mockctrl)
		test.setupMock(db, tokenGenerator)

		w := httptest.NewRecorder()

		router := chi.NewRouter()
		router.Post("/login", (&auth.Login{DB: db, TokenGenerator: tokenGenerator}).RegisterAndLogin)
		router.ServeHTTP(w, test.request())

		if w.Code != test.expectedStatusCode {
			t.Fatalf("[%s] Expect: %d, but got: %d", name, test.expectedStatusCode, w.Code)
		}
		if strings.TrimSpace(w.Body.String()) != test.expectedBody {
			t.Fatalf("[%s] Expect: %s, but got: %s", name, test.expectedBody, w.Body.String())
		}
		mockctrl.Finish()
	}
}

func Test_Login_REST_Succeed(t *testing.T) {
	tests := map[string]struct {
		setupMock          func(*mock_auth.MockLoginDB, *mock_auth.MockLoginTokenGenerator)
		request            func() *http.Request
		expectedStatusCode int
		expectedBody       string
	}{
		"Register New User And Login Succeed": {
			func(db *mock_auth.MockLoginDB, tokenGenerator *mock_auth.MockLoginTokenGenerator) {
				db.EXPECT().Exists(Eq("1234")).Return(false, nil)
				db.EXPECT().Create(Eq("1234"), Eq("abcd")).Return(nil)
				tokenGenerator.EXPECT().NewToken(Eq("1234")).Return("this is a token")
			},
			func() *http.Request {
				req := httptest.NewRequest("POST", "/login", strings.NewReader(`{ "username": "1234", "password": "abcd" }`))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			http.StatusCreated,
			`{"token":"this is a token"}`,
		},
		"Login Existing User Succeed": {
			func(db *mock_auth.MockLoginDB, tokenGenerator *mock_auth.MockLoginTokenGenerator) {
				db.EXPECT().Exists(Eq("1234")).Return(true, nil)
				db.EXPECT().Login(Eq("1234"), Eq("abcd")).Return(true, nil)
				tokenGenerator.EXPECT().NewToken(Eq("1234")).Return("this is a token")
			},
			func() *http.Request {
				req := httptest.NewRequest("POST", "/login", strings.NewReader(`{ "username": "1234", "password": "abcd" }`))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			http.StatusOK,
			`{"token":"this is a token"}`,
		},
	}

	for name, test := range tests {
		mockctrl := gomock.NewController(t)
		db := mock_auth.NewMockLoginDB(mockctrl)
		tokenGenerator := mock_auth.NewMockLoginTokenGenerator(mockctrl)
		test.setupMock(db, tokenGenerator)

		w := httptest.NewRecorder()

		router := chi.NewRouter()
		router.Post("/login", (&auth.Login{DB: db, TokenGenerator: tokenGenerator}).RegisterAndLogin)
		router.ServeHTTP(w, test.request())

		if w.Code != test.expectedStatusCode {
			t.Fatalf("[%s] Expect: %d, but got: %d", name, test.expectedStatusCode, w.Code)
		}
		if strings.TrimSpace(w.Body.String()) != test.expectedBody {
			t.Fatalf("[%s] Expect: %s, but got: %s", name, test.expectedBody, w.Body.String())
		}
		mockctrl.Finish()
	}
}
