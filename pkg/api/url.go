package api

func (h *API) urlTheme(key string) string {
	return h.cfg.ServerURL + "/theme/" + key
}

func (h *API) urlMedia(key string) string {
	return h.cfg.ServerURL + "/media/" + key
}

func (h *API) urlNote(key string) string {
	return h.cfg.ServerURL + "/" + key
}
