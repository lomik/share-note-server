package api

import (
	"mime"
	"net/http"

	"github.com/pkg/errors"
)

func (h *API) MediaView(w http.ResponseWriter, r *http.Request) {
	media, blobKey, err := h.media.Unhash(r.PathValue("media"))
	if err != nil {
		h.error(w, r, err, http.StatusNotFound)
		return
	}

	if media == nil {
		h.error(w, r, errors.New("media not found"), http.StatusNotFound)
		return
	}

	body, err := h.blob.Get(blobKey)
	if err != nil {
		h.error(w, r, err, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", getContentType(media.Filetype, body))
	w.Write(body)
}

func (h *API) MediaViewTheme(w http.ResponseWriter, r *http.Request) {
	uid, content, err := h.themes.Unhash(r.PathValue("theme"))
	if err != nil {
		h.error(w, r, err, http.StatusNotFound)
		return
	}

	if uid == "" {
		h.error(w, r, errors.New("theme not found"), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", mime.TypeByExtension(".css"))
	w.Write([]byte(content))
}
