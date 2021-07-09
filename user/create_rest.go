package user

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/fatdes/reap_backend_challenge/log"
	"github.com/fatdes/reap_backend_challenge/restutil"
	"github.com/go-chi/render"
)

type PostCreator interface {
	/// Create a new post
	/// Return the id associated
	CreatePost(post *Post) (string, error)
}

type Create struct {
	PostCreator PostCreator
}

func (c *Create) CreatePost(w http.ResponseWriter, r *http.Request) {
	log := log.GetLogger("user")
	defer func() {
		log.Add("canonical", "yes").Info("CreatePost")
		log.Sync()
	}()

	ctx := r.Context()
	username := ctx.Value("user").(string)
	if strings.TrimSpace(username) == "" {
		render.Render(w, r, &restutil.Error{
			StatusCode:   http.StatusUnauthorized,
			ErrorMessage: "You are not authorized to create any posts",
		})
		return
	}

	log.Add("username", username)

	if err := r.ParseMultipartForm(100); err != nil {
		log.Add("actual_error", err)
		render.Render(w, r, &restutil.Error{
			StatusCode:   http.StatusBadRequest,
			ErrorMessage: "Fail to upload image, Please contact support.",
		})
		return
	}

	imageFile, imageHeader, err := r.FormFile("image")
	if err != nil {
		log.Add("actual_error", err)
		render.Render(w, r, &restutil.Error{
			StatusCode:   http.StatusBadRequest,
			ErrorMessage: "Fail to upload image, please contact support.",
		})
		return
	}

	if l := imageHeader.Size; l > 300*1000 {
		err := &restutil.Error{
			StatusCode:   http.StatusBadRequest,
			ErrorMessage: fmt.Sprintf("\"filesize\" (length: %d) should be between 0 - 300k bytes", l),
		}
		log.Add("error", err)
		render.Render(w, r, err)
		return
	}

	description := r.MultipartForm.Value["description"][0]
	if l := len(description); l < 0 || l > 50 {
		err := &restutil.Error{
			StatusCode:   http.StatusBadRequest,
			ErrorMessage: fmt.Sprintf("\"description\" (length: %d) should be between 0 - 50 characters", l),
		}
		log.Add("error", err)
		render.Render(w, r, err)
		return
	}

	var image bytes.Buffer
	io.Copy(&image, imageFile)

	now := time.Now()
	id, err := c.PostCreator.CreatePost(&Post{
		Username:    username,
		Image:       image.Bytes(),
		Description: description,
		CreatedAt:   now,
	})
	if err != nil {
		log.Add("actual_error", err)
		render.Render(w, r, &restutil.Error{
			StatusCode:   http.StatusBadRequest,
			ErrorMessage: "Please contact support.",
		})
		return
	}

	render.JSON(w, r, &createPostResponse{
		URL:         fmt.Sprintf("/v1/user/post/%s/iamge", id),
		Username:    username,
		Description: description,
		CreatedAt:   now,
	})
}

type createPostResponse struct {
	URL         string    `json:"url"`
	Username    string    `json:"username"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
