package user

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/fatdes/reap_backend_challenge/log"
	"github.com/fatdes/reap_backend_challenge/restutil"
	"github.com/go-chi/render"
)

type PostLister interface {
	/// List posts with specific username and sort by created at
	ListPost(username string, createdAtSortOrder SortOrder) ([]*Post, error)
}

type List struct {
	PostLister PostLister
}

func (l *List) ListPost(w http.ResponseWriter, r *http.Request) {
	log := log.GetLogger("user")
	defer func() {
		log.Add("canonical", "yes").Info("ListPost")
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

	filterUsername := r.URL.Query().Get("username")
	createdAtSortOrder := Descending
	if strings.ToLower(r.URL.Query().Get("created_at_sort_order")) == "asc" {
		createdAtSortOrder = Ascending
	}

	postList, err := l.PostLister.ListPost(filterUsername, createdAtSortOrder)
	if err != nil {
		log.Add("actual_error", err)
		render.Render(w, r, &restutil.Error{
			StatusCode:   http.StatusBadRequest,
			ErrorMessage: "Fail to list posts, Please contact support.",
		})
		return
	}

	postResponseList := make([]*postResponse, len(postList))
	for i, p := range postList {
		postResponseList[i] = &postResponse{
			URL:         fmt.Sprintf("/v1/user/post/%s/image", p.ID),
			Username:    p.Username,
			Description: p.Description,
			CreatedAt:   p.CreatedAt,
		}
	}

	render.JSON(w, r, &listPostResponse{PostList: postResponseList})
}

type listPostResponse struct {
	PostList []*postResponse `json:"post_list"`
}

type postResponse struct {
	URL         string    `json:"url"`
	Username    string    `json:"username"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
