package api

import (
	"fmt"
	"io"
	"net/http"

	"github.com/k0kubun/pp"
)

func (h *API) Default(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Headers", "*")
	pp.Println(r.Method)

	pp.Println(r.Header)
	pp.Println(r.URL.String())

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	pp.Println(string(body))

	fmt.Fprintln(w, `{"success":true, "files":[]}`)
}
