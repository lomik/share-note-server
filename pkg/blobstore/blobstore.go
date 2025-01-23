package blobstore

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"tailscale.com/atomicfile"
)

type Store interface {
	Save(body []byte) (string, error)
	Get(hash string) ([]byte, error)
	Exists(hash string) (bool, error)
}

type filestore struct {
	root string
}

func New(root string) Store {
	return &filestore{root: root}
}

func validate(hash string) error {
	if len(hash) != 64 {
		return fmt.Errorf("invalid hash %#v", hash)
	}

	if strings.Trim(hash, "0123456789abcdef") != "" {
		return fmt.Errorf("invalid hash %#v", hash)
	}

	return nil
}

func (s *filestore) filename(hash string) string {
	return filepath.Join(s.root, hash[:2], hash[2:4], hash[4:])
}

func (s *filestore) Save(body []byte) (string, error) {
	hash := fmt.Sprintf("%x", sha256.Sum256(body))

	exists, err := s.Exists(hash)
	if err != nil {
		return "", err
	}
	if exists {
		return hash, nil
	}

	fn := s.filename(hash)

	// create dirs
	if err = os.MkdirAll(filepath.Dir(fn), 0755); err != nil {
		return "", err
	}

	// write file
	if err = atomicfile.WriteFile(fn, body, 0644); err != nil {
		return "", err
	}
	return hash, nil
}

func (s *filestore) Get(hash string) ([]byte, error) {
	if err := validate(hash); err != nil {
		return nil, err
	}
	return os.ReadFile(s.filename(hash))
}

func (s *filestore) Exists(hash string) (bool, error) {
	if err := validate(hash); err != nil {
		return false, err
	}

	fn := s.filename(hash)
	stat, err := os.Stat(fn)

	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return false, err
	}

	if err == nil { // file exists
		if !stat.Mode().IsRegular() {
			return false, fmt.Errorf("%#v is not a regular file", fn)
		}
		// file exists and ok
		return true, nil
	}

	return false, nil
}

func (s *filestore) SaveJson(obj interface{}) (string, error) {
	body, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return s.Save(body)
}

func (s *filestore) GetJson(hash string, obj interface{}) error {
	body, err := s.Get(hash)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, obj)
}
