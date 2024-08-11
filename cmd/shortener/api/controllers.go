package api

import (
	"github.com/IgorViskov/go_33_shortener/internal/config"
	"github.com/IgorViskov/go_33_shortener/internal/shs"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type mainController struct {
	path    string
	service *shs.ShortenerService
	config  *config.AppConfig
}

func (c mainController) Get(w http.ResponseWriter, req *http.Request) {
	short := req.URL.Path[1:]
	moved, err := c.service.UnShort(short)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, req, moved, http.StatusTemporaryRedirect)
}

func (c mainController) Post(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	u, okValidate := validateURL(string(body))
	if !okValidate {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	shorted, err := c.service.Short(u)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	redirect := c.config.RedirectAddress
	redirect.Path = shorted
	_, ok := w.Write([]byte(redirect.String()))
	if ok != nil {
		http.Error(w, ok.Error(), http.StatusInternalServerError)
		return
	}
}

func (c mainController) GetPath() string {
	return c.path
}

func NewMainController(config *config.AppConfig) *mainController {
	return &mainController{
		path:    "/",
		service: shs.NewShortenerService(),
		config:  config,
	}
}

func validateURL(u string) (string, bool) {
	if len(strings.TrimSpace(u)) == 0 {
		return "", false
	}
	_, err := url.Parse(u)
	if err != nil {
		return "", false
	}
	return u, true
}
