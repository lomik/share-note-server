package api

import (
	"net/http"

	"github.com/lomik/share-note-server/pkg/models"
	"github.com/lomik/share-note-server/pkg/random"
	"go.uber.org/zap"
)

func (h *API) GetKey(w http.ResponseWriter, r *http.Request) {
	uid := r.URL.Query().Get("id")
	if uid == "" {
		h.status(w, r, http.StatusBadRequest)
		return
	}

	currentUser, _, err := h.users.Get(uid)
	if err != nil {
		h.error(w, r, err, http.StatusInternalServerError)
		return
	}

	if currentUser == nil && !h.cfg.AllowCreateApiKey {
		h.render(w, r, "deny_create_key.go.html", nil)
		return
	}
	if currentUser != nil && !h.cfg.AllowChangeApiKey {
		h.render(w, r, "deny_change_key.go.html", nil)
		return
	}

	apiKey := random.UID()

	h.logger(r).Info("update or create apiKey", zap.String("uid", uid))
	if _, err = h.users.Set(uid, &models.User{
		UID:    uid,
		ApiKey: apiKey,
	}); err != nil {
		h.error(w, r, err, http.StatusInternalServerError)
		return
	}

	h.render(w, r, "get_key.go.html", map[string]interface{}{
		"apiKey": apiKey,
	})
}
