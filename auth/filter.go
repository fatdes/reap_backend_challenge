package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/fatdes/reap_backend_challenge/log"
	"github.com/fatdes/reap_backend_challenge/restutil"
	"github.com/go-chi/render"
)

type TokenVerifier interface {
	VerifyToken(tokenString string) (string, error)
}

type AuthFilter struct {
	Verifier TokenVerifier
}

func (f *AuthFilter) AccessTokenFilter() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			log := log.GetLogger("AccessTokenFilter")

			token, err := parseBearerToken(r)
			if err != nil {
				log.Debugw("unable to parse bearer token", "error", err)
				render.Render(w, r, &restutil.Error{
					StatusCode:   http.StatusBadRequest,
					ErrorMessage: "Invalid bearer token",
				})
				return
			}

			username, err := f.Verifier.VerifyToken(token)
			if err != nil {
				log.Debugw("invalid bearer token", "error", err)
				render.Render(w, r, &restutil.Error{
					StatusCode:   http.StatusBadRequest,
					ErrorMessage: "Invalid bearer token",
				})
				return
			}

			requestCtx := context.WithValue(r.Context(), "user", username)
			next.ServeHTTP(w, r.WithContext(requestCtx))
		}

		return http.HandlerFunc(fn)
	}
}

func parseBearerToken(r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		return "", errors.New("bearer token not found")
	}

	bearerToken := strings.TrimPrefix(auth, "Bearer ")
	return bearerToken, nil
}
