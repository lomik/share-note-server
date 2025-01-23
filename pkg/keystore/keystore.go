package keystore

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"tailscale.com/atomicfile"
)

type Store[K, V any] interface {
	Set(key K, value V) (string, error)
	Get(key K) (V, string, error)
	Delete(key K) error
	Unhash(hash string) (K, V, error)
}

type record[K, V any] struct {
	Key   K `json:"key"`
	Value V `json:"value"`
}

type filestore[K, V any] struct {
	root string
}

func New[K, V any](root string) Store[K, V] {
	return &filestore[K, V]{root: root}
}

func (s *filestore[K, V]) filename(hash string) string {
	return filepath.Join(s.root, hash[:2], hash[2:4], hash[4:])
}

func (s *filestore[K, V]) Set(key K, value V) (string, error) {
	keyBody, err := json.Marshal(key)
	if err != nil {
		return "", err
	}

	body, err := json.Marshal(&record[K, V]{
		Key:   key,
		Value: value,
	})
	if err != nil {
		return "", err
	}

	hash := fmt.Sprintf("%x", sha256.Sum256(keyBody))
	fn := s.filename(hash)

	// create dirs
	if err := os.MkdirAll(filepath.Dir(fn), 0755); err != nil {
		return "", err
	}

	// write file
	if err := atomicfile.WriteFile(fn, body, 0644); err != nil {
		return "", err
	}
	return hash, nil
}

func (s *filestore[K, V]) Delete(key K) error {
	keyBody, err := json.Marshal(key)
	if err != nil {
		return err
	}

	hash := fmt.Sprintf("%x", sha256.Sum256(keyBody))
	fn := s.filename(hash)
	return os.Remove(fn)
}

func (s *filestore[K, V]) Get(key K) (V, string, error) {
	var empty V
	keyBody, err := json.Marshal(key)
	if err != nil {
		return empty, "", err
	}
	hash := fmt.Sprintf("%x", sha256.Sum256(keyBody))
	body, err := os.ReadFile(s.filename(hash))
	if os.IsNotExist(err) {
		return empty, "", nil
	}
	if err != nil {
		return empty, "", err
	}

	var rec record[K, V]
	err = json.Unmarshal(body, &rec)
	if err != nil {
		return empty, "", err
	}

	return rec.Value, hash, nil
}

func (s *filestore[K, V]) Unhash(hash string) (K, V, error) {
	var emptyK K
	var emptyV V
	if len(hash) != 64 || strings.Trim(hash, "0123456789abcdef") != "" {
		return emptyK, emptyV, fmt.Errorf("invalid hash")
	}

	body, err := os.ReadFile(s.filename(hash))
	if err != nil {
		return emptyK, emptyV, err
	}

	var rec record[K, V]
	err = json.Unmarshal(body, &rec)
	if err != nil {
		return emptyK, emptyV, err
	}

	return rec.Key, rec.Value, nil
}
