package models

type Media struct {
	UID        string `json:"uid"`
	Filetype   string `json:"filetype"`
	Hash       string `json:"hash"`
	ByteLength int    `json:"byteLength"`
}
