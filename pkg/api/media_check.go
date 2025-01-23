package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/lomik/share-note-server/pkg/models"
	"go.uber.org/zap"
)

func (h *API) MediaCheck(w http.ResponseWriter, r *http.Request) {
	type File struct {
		Filetype   string `json:"filetype"`
		Hash       string `json:"hash"`
		ByteLength int    `json:"byteLength"`
		Url        string `json:"url"`
	}

	type CSS struct {
		URL string `json:"url"`
	}

	type Request struct {
		Files []struct {
			Filetype   string `json:"filetype"`
			Hash       string `json:"hash"`
			ByteLength int    `json:"byteLength"`
		} `json:"files"`
	}

	type Response struct {
		Success bool        `json:"success"`
		Files   []File      `json:"files"`
		CSS     interface{} `json:"css"` // CSS or false
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

	resp := Response{
		Files: make([]File, 0),
	}

	for i := 0; i < len(req.Files); i++ {
		_, mediaKey, err := h.media.Get(
			&models.Media{
				Filetype:   req.Files[i].Filetype,
				Hash:       req.Files[i].Hash,
				ByteLength: req.Files[i].ByteLength,
				UID:        r.Header.Get("X-Sharenote-Id"),
			},
		)
		if err != nil {
			zap.L().Error("media exists error", zap.Error(err))
		}

		if mediaKey != "" {
			resp.Files = append(resp.Files, File{
				Filetype:   req.Files[i].Filetype,
				Hash:       req.Files[i].Hash,
				ByteLength: req.Files[i].ByteLength,
				Url:        h.urlMedia(mediaKey),
			})
		} else {
			h.logger(r).Info("file not exists",
				zap.String("hash", req.Files[i].Hash),
				zap.Int("byteLength", req.Files[i].ByteLength),
				zap.String("filetype", req.Files[i].Filetype),
			)
		}
	}

	resp.CSS = false
	resp.Success = true

	_, themeKey, err := h.themes.Get(r.Header.Get("X-Sharenote-Id"))
	if err != nil {
		h.error(w, r, err, http.StatusInternalServerError)
		return
	}
	if themeKey != "" {
		resp.CSS = CSS{
			URL: h.urlTheme(themeKey),
		}
	}

	h.json(w, r, resp)
}
