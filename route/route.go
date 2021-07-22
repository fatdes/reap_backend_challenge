package route

import (
	"net/http"
	"os"
	"time"

	"github.com/fatdes/reap_backend_challenge/auth"
	"github.com/fatdes/reap_backend_challenge/user"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RealIP)

	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"}, // should be configurable instead of allow everyone
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	dbURL := os.Getenv("DATABASE_URL")
	authDB := &auth.DBPGX{URL: dbURL}
	token := &auth.Token{}

	authFilter := &auth.AuthFilter{
		Verifier: token,
	}

	r.Route("/v1/auth", func(r chi.Router) {
		r.Post("/login", (&auth.Login{DB: authDB, TokenGenerator: token}).RegisterAndLogin)
	})

	postDB := &user.PostDBPGX{URL: os.Getenv("DATABASE_URL")}

	r.Route("/v1/user", func(r chi.Router) {
		r.Use(authFilter.AccessTokenFilter())
		r.Post("/post", (&user.Create{PostCreator: postDB}).CreatePost)
		r.Get("/post", (&user.List{PostLister: postDB}).ListPost)
		r.Get("/post/{post_id}/image", (&user.Get{PostImageGetter: postDB}).GetPostImage)
	})

	return r
}
