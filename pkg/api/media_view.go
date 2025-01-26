package api

import (
	"bytes"
	"mime"
	"net/http"
	"strings"

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

	var contentType string

	// excalidraw
	if strings.HasPrefix(media.Filetype, "md/") {
		// detect by content
		if bytes.HasPrefix(body, []byte("<svg")) {
			contentType = mime.TypeByExtension(".svg")
		}
	}

	if contentType == "" {
		contentType = mime.TypeByExtension("." + media.Filetype)
	}

	w.Header().Set("Content-Type", contentType)
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
