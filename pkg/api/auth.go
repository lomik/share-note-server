package api

import (
	"crypto/sha256"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

func (h *API) auth(w http.ResponseWriter, r *http.Request) bool {
	uid := r.Header.Get("X-Sharenote-Id")
	nonce := r.Header.Get("X-Sharenote-Nonce")
	key := r.Header.Get("X-Sharenote-Key")

	if uid == "" || nonce == "" || key == "" {
		h.error(w, r, errors.New("empty header"), http.StatusForbidden)
		return false
	}

	currentUser, _, err := h.users.Get(uid)
	if err != nil {
		h.error(w, r, err, http.StatusInternalServerError)
		return false
	}

	if currentUser == nil {
		h.error(w, r, errors.New("user not found"), http.StatusForbidden)
		return false
	}

	digest := fmt.Sprintf("%x", sha256.Sum256([]byte(nonce+currentUser.ApiKey)))

	if digest != key {
		h.error(w, r, errors.New("invalid digest"), http.StatusUnauthorized)
		return false
	}

	return true
}
