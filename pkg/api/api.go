package api

import (
	"bytes"
	"embed"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/lomik/share-note-server/pkg/blobstore"
	"github.com/lomik/share-note-server/pkg/config"
	"github.com/lomik/share-note-server/pkg/keystore"
	"github.com/lomik/share-note-server/pkg/models"
	"go.uber.org/zap"
)

//go:embed templates/*
var templatesFS embed.FS

type Options struct {
	Config *config.Config
	Users  keystore.Store[string, *models.User]
	Themes keystore.Store[string, string]
	Media  keystore.Store[*models.Media, string]
	Blob   blobstore.Store
	Notes  keystore.Store[string, *models.Note]
}

type API struct {
	cfg    *config.Config
	tpl    *template.Template
	users  keystore.Store[string, *models.User]
	themes keystore.Store[string, string]
	media  keystore.Store[*models.Media, string]
	blob   blobstore.Store
	notes  keystore.Store[string, *models.Note]
}

func New(opts Options) *API {
	funcMap := template.FuncMap{
		"attr": func(s string) template.HTMLAttr {
			return template.HTMLAttr(s)
		},
		"safe": func(s string) template.HTML {
			return template.HTML(s)
		},
	}
	tpl, err := template.New("").Funcs(funcMap).ParseFS(templatesFS, "templates/*.go.html")
	if err != nil {
		log.Fatal(err)
	}

	return &API{
		cfg:    opts.Config,
		tpl:    tpl,
		users:  opts.Users,
		themes: opts.Themes,
		media:  opts.Media,
		blob:   opts.Blob,
		notes:  opts.Notes,
	}
}

func (h *API) json(w http.ResponseWriter, r *http.Request, obj interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if r.Header.Get("Origin") != "" {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Vary", "Origin")
	}

	resp, err := json.Marshal(obj)
	if err != nil {
		h.error(w, r, err, http.StatusInternalServerError)
		return
	}

	io.Copy(w, bytes.NewReader(resp))
}

func (h *API) render(w http.ResponseWriter, r *http.Request, tpl string, obj interface{}) {
	w.Header().Add("Content-Type", "text/html")
	err := h.tpl.ExecuteTemplate(w, tpl, obj)
	if err != nil {
		h.error(w, r, err, http.StatusInternalServerError)
		return
	}
}

func (h *API) status(w http.ResponseWriter, _ *http.Request, code int) {
	http.Error(w, http.StatusText(code), code)
}

func (h *API) error(w http.ResponseWriter, r *http.Request, err error, code int) {
	h.logger(r).Error(http.StatusText(code), zap.Error(err))
	h.status(w, r, code)
}

func (h *API) logger(r *http.Request) *zap.Logger {
	ret := zap.L()
	if r.Header.Get("X-Sharenote-Id") != "" {
		ret = ret.With(zap.String("uid", r.Header.Get("X-Sharenote-Id")))
	}
	if r.Header.Get("X-Sharenote-Hash") != "" {
		ret = ret.With(zap.String("hash", r.Header.Get("X-Sharenote-Hash")))
	}
	return ret
}
