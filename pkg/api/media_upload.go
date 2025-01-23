package api

import (
	"crypto/sha1"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/lomik/share-note-server/pkg/models"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

/*
  "X-Sharenote-Bytelength": []string{
    "367891",
  },
  "X-Sharenote-Hash": []string{
    "a65a7f7382243186310c3dd90bccd6faf3a581f8",
  },
  "X-Sharenote-Id": []string{
    "045767078ab9e7623ef3678a21005000",
  },
  "X-Sharenote-Key": []string{
    "9129988369dcd0df99141bfddae962ff8d9f1f1145b5065e22ef84b340f46109",
  },
	"X-Sharenote-Filetype": []string{
		"css",
	},
	"X-Sharenote-Nonce": []string{
		"1737478909399",
	},
*/

func (h *API) MediaUpload(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Success bool   `json:"success"`
		URL     string `json:"url"`
	}

	if !h.auth(w, r) {
		return
	}

	uid := r.Header.Get("X-Sharenote-Id")
	filetype := r.Header.Get("X-Sharenote-Filetype")
	byteLength, err := strconv.Atoi(r.Header.Get("X-Sharenote-Bytelength"))
	if err != nil {
		h.error(w, r, err, http.StatusBadRequest)
		return
	}
	hash := r.Header.Get("X-Sharenote-Hash")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.error(w, r, err, http.StatusBadRequest)
		return
	}

	if hash != fmt.Sprintf("%x", sha1.Sum(body)) {
		h.error(w, r, errors.New("invalid hash"), http.StatusBadRequest)
		return
	}

	// css failed this validation
	// if len(body) != byteLength {
	// 	http.Error(w, "Invalid X-Sharenote-Bytelength", http.StatusBadRequest)
	// 	return
	// }

	// save theme. always to same file
	if filetype == "css" {
		themeKey, err := h.themes.Set(uid, string(body))
		if err != nil {
			h.error(w, r, errors.New("can't save theme"), http.StatusInternalServerError)
			return
		}

		h.logger(r).Info("theme saved",
			zap.String("themeKey", themeKey),
		)

		h.json(w, r, &Response{
			URL:     h.urlTheme(themeKey),
			Success: true,
		})
		return
	}

	// save other
	blobKey, err := h.blob.Save(body)
	if err != nil {
		h.error(w, r, err, http.StatusInternalServerError)
		return
	}

	media := models.Media{
		UID:        uid,
		Filetype:   filetype,
		Hash:       hash,
		ByteLength: byteLength,
	}

	mediaKey, err := h.media.Set(&media, blobKey)
	if err != nil {
		h.error(w, r, err, http.StatusInternalServerError)
		return
	}

	h.logger(r).Info("media saved",
		zap.String("mediaKey", mediaKey),
		zap.String("hash", media.Hash),
		zap.Int("byteLength", media.ByteLength),
		zap.String("filetype", media.Filetype),
		zap.String("blobKey", blobKey),
	)

	h.json(w, r, &Response{
		URL:     h.urlMedia(mediaKey),
		Success: true,
	})
}
