package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/lomik/share-note-server/pkg/models"
	"github.com/lomik/share-note-server/pkg/random"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const noteKeyLength = 16

type NoteCreateResponse struct {
	Success bool   `json:"success"`
	URL     string `json:"url"`
}

func (h *API) NoteCreate(w http.ResponseWriter, r *http.Request) {
	if !h.auth(w, r) {
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.error(w, r, err, http.StatusBadRequest)
		return
	}

	var n models.Note
	err = json.Unmarshal(body, &n)
	if err != nil {
		h.error(w, r, err, http.StatusBadRequest)
		return
	}

	if !random.IsToken(noteKeyLength, n.Filename) {
		n.Filename = random.Token(noteKeyLength)
	}
	n.UID = r.Header.Get("X-Sharenote-Id")

	currentNote, _, err := h.notes.Get(n.Filename)
	if err != nil {
		h.error(w, r, err, http.StatusInternalServerError)
		return
	}
	if currentNote != nil && currentNote.UID != n.UID {
		h.error(w, r, errors.New("note owner mismatch"), http.StatusForbidden)
		return
	}

	noteKey, err := h.notes.Set(n.Filename, &n)
	if err != nil {
		h.error(w, r, err, http.StatusInternalServerError)
		return
	}

	h.logger(r).Info("note saved", zap.String("note", noteKey))

	h.json(w, r, &NoteCreateResponse{
		URL:     h.urlNote(n.Filename),
		Success: true,
	})
}
