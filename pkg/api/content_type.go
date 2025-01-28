package api

import (
	"bytes"
	"mime"
	"strings"
)

func getContentType(filetype string, body []byte) string {
	// excalidraw
	if strings.HasPrefix(filetype, "md/") {
		// detect by content
		if bytes.HasPrefix(body, []byte("<svg")) {
			return mime.TypeByExtension(".svg")
		}
	}

	// "filetype":"image/svg+xml;base64,PHN2ZyB3aWR0aD0iMTYiIG
	if strings.Contains(filetype, ";") {
		mimeType := strings.Split(filetype, ";")[0]
		ext, err := mime.ExtensionsByType(mimeType)
		if err == nil && len(ext) > 0 {
			return mimeType
		}
	}

	return mime.TypeByExtension("." + filetype)
}
