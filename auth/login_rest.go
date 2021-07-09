package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/fatdes/reap_backend_challenge/log"
	"github.com/fatdes/reap_backend_challenge/restutil"
	"github.com/go-chi/render"
)

var exampleLogin string

func init() {
	bytes, _ := json.Marshal(&login{
		Username: "...",
		Password: "...",
	})
	exampleLogin = string(bytes)
}

type LoginDB interface {
	/// Create a new login
	Create(username string, password string) error

	/// True if the username already exists
	Exists(username string) (bool, error)

	/// True if the username and password matches
	Login(username string, password string) (bool, error)
}

type LoginTokenGenerator interface {
	NewToken(username string) string
}

//go:generate mockgen -destination=../mock/auth/login_mocks.go -package=mock_auth -source login_rest.go LoginDB LoginTokenGenerator
type Login struct {
	DB             LoginDB
	TokenGenerator LoginTokenGenerator
}

func (l *Login) RegisterAndLogin(w http.ResponseWriter, r *http.Request) {
	log := log.GetLogger("login")
	defer func() {
		log.Add("canonical", "yes").Info("RegisterAndLogin")
		log.Sync()
	}()

	data := &login{}
	if err := render.Bind(r, data); err != nil {
		e := &restutil.Error{
			StatusCode:   http.StatusBadRequest,
			ErrorMessage: err.Error(),
		}
		if errors.Is(err, io.EOF) {
			e = &restutil.Error{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("request body should be in json %s", exampleLogin),
			}
		}
		log.Add("actual_error", err).Add("error", e)
		render.Render(w, r, e)
		return
	}

	log.Add("username", data.Username)
	exists, err := l.DB.Exists(data.Username)
	if err != nil {
		e := &restutil.Error{
			StatusCode:   http.StatusInternalServerError,
			ErrorMessage: "We will fix the issues ASAP. Please try again later",
		}
		log.Add("actual_error", err).Add("error", e)
		render.Render(w, r, e)
		return
	}
	log.Add("exists", exists)

	status := http.StatusOK
	if !exists {
		log.Debug("registering new user")
		err := l.DB.Create(data.Username, data.Password)
		if err != nil {
			e := &restutil.Error{
				StatusCode:   http.StatusInternalServerError,
				ErrorMessage: "We will fix the issues ASAP. Please try again later",
			}
			log.Add("actual_error", err).Add("error", e)
			render.Render(w, r, e)
			return
		}
		status = http.StatusCreated
	} else {
		log.Debug("logging in existing user")
		succeed, err := l.DB.Login(data.Username, data.Password)
		if err != nil {
			e := &restutil.Error{
				StatusCode:   http.StatusInternalServerError,
				ErrorMessage: "We will fix the issues ASAP. Please try again later",
			}
			log.Add("actual_error", err).Add("error", e)
			render.Render(w, r, e)
			return
		}
		if !succeed {
			e := &restutil.Error{
				StatusCode:   http.StatusUnauthorized,
				ErrorMessage: "username or password not match",
			}
			log.Add("error", e)
			render.Render(w, r, e)
			return
		}
	}

	log.Debug("generates jwt token")
	token := l.TokenGenerator.NewToken(data.Username)

	w.WriteHeader(status)
	render.JSON(w, r, &loginResponse{Token: token})
}

type login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (data *login) Bind(r *http.Request) error {
	data.Username = strings.TrimSpace(data.Username)
	if l := len(data.Username); l < 4 || l > 20 {
		return errors.New(fmt.Sprintf("\"username\" (length: %d) should be between 4 - 20 characters", l))
	}

	data.Password = strings.TrimSpace(data.Password)
	if l := len(data.Password); l < 4 || l > 20 {
		return errors.New(fmt.Sprintf("\"password\" (length: %d) should be between 4 - 20 characters", l))
	}

	return nil
}
