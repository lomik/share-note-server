package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/lomik/share-note-server/pkg/models"
)

func (h *API) NoteView(w http.ResponseWriter, r *http.Request) {
	type TemplateVars struct {
		Note      *models.Note
		ServerURL string
		Elements  map[string]string
		Theme     string
	}

	w.Header().Set("Content-Type", "text/html")

	currentNote, _, err := h.notes.Get(r.PathValue("note"))
	if err != nil {
		h.error(w, r, err, http.StatusNotFound)
		return
	}
	if currentNote == nil {
		h.render(w, r, "note_not_found.go.html", nil)
		return
	}

	var data TemplateVars
	data.Note = currentNote
	data.ServerURL = h.cfg.ServerURL

	data.Elements = map[string]string{}
	for _, element := range currentNote.Template.Elements {
		if element.Element == "" {
			continue
		}
		data.Elements[element.Element] = fmt.Sprintf(
			"class=\"%s\" style=\"%s\"",
			strings.Join(element.Classes, " "),
			element.Style,
		)
	}

	_, themeKey, err := h.themes.Get(currentNote.UID)
	if err != nil {
		h.error(w, r, err, http.StatusInternalServerError)
		return
	}

	if themeKey != "" {
		data.Theme = h.urlTheme(themeKey)
	}

	h.render(w, r, "note.go.html", data)
}
