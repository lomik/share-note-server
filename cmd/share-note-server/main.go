package main

import (
	"flag"
	"net/http"
	"path/filepath"

	"github.com/lomik/share-note-server/pkg/api"
	"github.com/lomik/share-note-server/pkg/blobstore"
	"github.com/lomik/share-note-server/pkg/config"
	"github.com/lomik/share-note-server/pkg/keystore"
	"github.com/lomik/share-note-server/pkg/models"
	"github.com/lomik/share-note-server/pkg/static"
	"go.uber.org/zap"
)

func main() {
	configFn := flag.String("config", "config.yaml", "Config filename")
	flag.Parse()

	// configure logging
	zapConfig := zap.NewProductionConfig()
	zapConfig.DisableStacktrace = true
	logger := zap.Must(zapConfig.Build())
	defer zap.RedirectStdLog(logger)()
	defer zap.ReplaceGlobals(logger)()

	cfg, err := config.LoadFromFile(*configFn)
	if err != nil {
		zap.L().Fatal(
			"can't load config from file",
			zap.String("filename", *configFn),
			zap.Error(err),
		)
	}

	mux := http.NewServeMux()

	// init storages

	api := api.New(api.Options{
		Config: cfg,
		Users:  keystore.New[string, *models.User](filepath.Join(cfg.Data, "user")),
		Themes: keystore.New[string, string](filepath.Join(cfg.Data, "theme")),
		Media:  keystore.New[*models.Media, string](filepath.Join(cfg.Data, "media")),
		Blob:   blobstore.New(filepath.Join(cfg.Data, "blob")),
		Notes:  keystore.New[string, *models.Note](filepath.Join(cfg.Data, "note")),
	})

	fs := http.FileServer(http.FS(static.FS))

	mux.HandleFunc("OPTIONS /", api.Options)
	mux.HandleFunc("GET /v1/account/get-key", api.GetKey)
	mux.HandleFunc("POST /v1/file/upload", api.MediaUpload)
	mux.HandleFunc("POST /v1/file/check-files", api.MediaCheck)
	mux.HandleFunc("POST /v1/file/create-note", api.NoteCreate)
	mux.HandleFunc("POST /v1/file/delete", api.NoteDelete)
	mux.Handle("GET /favicon.ico", fs)
	mux.Handle("GET /static/{media}", http.StripPrefix("/static/", fs))
	mux.HandleFunc("GET /media/{media}", api.MediaView)
	mux.HandleFunc("GET /theme/{theme}", api.MediaViewTheme)
	mux.HandleFunc("GET /{note}", api.NoteView)

	http.ListenAndServe(cfg.Listen, mux)
}
