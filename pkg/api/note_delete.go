package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/lomik/share-note-server/pkg/random"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (h *API) NoteDelete(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Filename string `json:"filename"`
	}
	type Response struct {
		Success bool `json:"success"`
	}
	if !h.auth(w, r) {
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.error(w, r, err, http.StatusBadRequest)
		return
	}

	var req Request
	if err = json.Unmarshal(body, &req); err != nil {
		h.error(w, r, err, http.StatusBadRequest)
		return
	}

	if !random.IsToken(noteKeyLength, req.Filename) {
		h.error(w, r, err, http.StatusBadRequest)
		return
	}

	currentNote, _, err := h.notes.Get(req.Filename)
	if err != nil {
		h.error(w, r, err, http.StatusInternalServerError)
		return
	}

	if currentNote == nil {
		h.json(w, r, &Response{
			Success: true,
		})
		return
	}

	if currentNote.UID != r.Header.Get("X-Sharenote-Id") {
		h.error(w, r, errors.New("wrong owner"), http.StatusForbidden)
		return
	}

	if err = h.notes.Delete(req.Filename); err != nil {
		h.error(w, r, err, http.StatusInternalServerError)
		return
	}

	h.logger(r).Info("note deleted", zap.String("filename", req.Filename))

	h.json(w, r, &Response{
		Success: true,
	})
}
