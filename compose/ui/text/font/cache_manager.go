package font

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// CacheManager handles local caching of font files.
type CacheManager struct {
	cacheDir string
}

// NewCacheManager creates a new CacheManager.
// It ensures the cache directory exists.
func NewCacheManager(cacheDir string) (*CacheManager, error) {
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}
	return &CacheManager{cacheDir: cacheDir}, nil
}

// Get returns the path to the cached file if it exists.
func (c *CacheManager) Get(url string) (string, bool) {
	filename := c.generateFilename(url)
	path := filepath.Join(c.cacheDir, filename)
	if _, err := os.Stat(path); err == nil {
		return path, true
	}
	return "", false
}

// Put saves the data to the cache and returns the path.
func (c *CacheManager) Put(url string, data []byte) (string, error) {
	filename := c.generateFilename(url)
	path := filepath.Join(c.cacheDir, filename)
	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write cache file: %w", err)
	}
	return path, nil
}

func (c *CacheManager) generateFilename(url string) string {
	hash := sha256.Sum256([]byte(url))
	return hex.EncodeToString(hash[:])
}
