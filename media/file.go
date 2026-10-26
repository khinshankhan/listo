package media

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// FileStore keeps images under a directory, one file per key.
type FileStore struct {
	dir    string
	prefix string
}

// NewFileStore stores under dir and serves keys at prefix, eg /media/.
func NewFileStore(dir, prefix string) *FileStore {
	return &FileStore{dir: dir, prefix: prefix}
}

// Put writes to a temp file beside the target and renames it into place, so
// a half-written image is never served.
func (s *FileStore) Put(ctx context.Context, key string, r io.Reader, _ string) error {
	if err := checkKey(key); err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	path := s.path(key)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("putting %s: %w", key, err)
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".put-*")
	if err != nil {
		return fmt.Errorf("putting %s: %w", key, err)
	}
	// a no-op once renamed
	defer os.Remove(tmp.Name())

	if _, err := io.Copy(tmp, r); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("putting %s: %w", key, err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("putting %s: %w", key, err)
	}
	// CreateTemp makes it 0600
	if err := os.Chmod(tmp.Name(), 0o644); err != nil {
		return fmt.Errorf("putting %s: %w", key, err)
	}
	if err := os.Rename(tmp.Name(), path); err != nil {
		return fmt.Errorf("putting %s: %w", key, err)
	}
	return nil
}

func (s *FileStore) Open(ctx context.Context, key string) (io.ReadCloser, error) {
	if err := checkKey(key); err != nil {
		return nil, err
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return os.Open(s.path(key))
}

func (s *FileStore) URL(key string) string {
	return s.prefix + key
}

func (s *FileStore) path(key string) string {
	return filepath.Join(s.dir, filepath.FromSlash(key))
}
