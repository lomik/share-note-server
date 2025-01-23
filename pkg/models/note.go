package models

type Note struct {
	UID      string `json:"uid"`
	Filename string `json:"filename"`
	Filetype string `json:"filetype"`
	Hash     string `json:"hash"`
	Template struct {
		Elements []struct {
			Element string   `json:"element"`
			Classes []string `json:"classes"`
			Style   string   `json:"style"`
		} `json:"elements"`
		Encrypted   bool   `json:"encrypted"`
		Content     string `json:"content"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Width       string `json:"width"`
		MathJax     bool   `json:"mathJax"`
	} `json:"template"`
}
