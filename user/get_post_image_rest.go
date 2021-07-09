package user

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/fatdes/reap_backend_challenge/log"
	"github.com/fatdes/reap_backend_challenge/restutil"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type PostImageGetter interface {
	/// Get image of specific post ID
	GetPostImage(id string) ([]byte, error)
}

type Get struct {
	PostImageGetter PostImageGetter
}

func (g *Get) GetPostImage(w http.ResponseWriter, r *http.Request) {
	log := log.GetLogger("user")
	defer func() {
		log.Add("canonical", "yes").Info("GetPostImage")
		log.Sync()
	}()

	ctx := r.Context()
	username := ctx.Value("user").(string)
	if strings.TrimSpace(username) == "" {
		render.Render(w, r, &restutil.Error{
			StatusCode:   http.StatusUnauthorized,
			ErrorMessage: "You are not authorized to list any posts",
		})
		return
	}

	log.Add("username", username)

	postID := chi.URLParam(r, "post_id")
	if len(postID) == 0 {
		render.Render(w, r, &restutil.Error{
			StatusCode:   http.StatusBadRequest,
			ErrorMessage: "Missing post_id in URL",
		})
		return
	}

	image, err := g.PostImageGetter.GetPostImage(postID)
	if err != nil {
		log.Add("actual_error", err)
		render.Render(w, r, &restutil.Error{
			StatusCode:   http.StatusBadRequest,
			ErrorMessage: "Fail to get post image, Please contact support.",
		})
		return
	}

	render.PlainText(w, r, base64.StdEncoding.EncodeToString(image))
}
