package auth_test

import (
	"strings"
	"testing"

	"github.com/fatdes/reap_backend_challenge/auth"
)

func Test_NewToken(t *testing.T) {
	tests := map[string]struct {
		username      string
		expectedToken string
	}{
		"New Token": {
			username:      "1234",
			expectedToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0In0.GcmUH71YjQSuMVO6VbwK4ipmztNGuckeP6XlO7XLhcM",
		},
	}

	for name, test := range tests {
		tokenGenerator := &auth.Token{}

		actual := tokenGenerator.NewToken(test.username)

		if strings.TrimSpace(actual) != test.expectedToken {
			t.Fatalf("[%s] Expect: %s, but got: %s", name, test.expectedToken, actual)
		}
	}
}

func Test_VerifyToken(t *testing.T) {
	tests := map[string]struct {
		token            string
		expectedUsername string
	}{
		"New Token": {
			token:            "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0In0.GcmUH71YjQSuMVO6VbwK4ipmztNGuckeP6XlO7XLhcM",
			expectedUsername: "1234",
		},
	}

	for name, test := range tests {
		tokenGenerator := &auth.Token{}

		actual, err := tokenGenerator.VerifyToken(test.token)

		if err != nil {
			t.Fatalf("[%s] Expect: nil, but got: %s", name, err)
		}

		if strings.TrimSpace(actual) != test.expectedUsername {
			t.Fatalf("[%s] Expect: %s, but got: %s", name, test.expectedUsername, actual)
		}
	}
}
