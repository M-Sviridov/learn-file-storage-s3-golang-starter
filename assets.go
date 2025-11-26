package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (cfg apiConfig) ensureAssetsDir() error {
	if _, err := os.Stat(cfg.assetsRoot); os.IsNotExist(err) {
		return os.Mkdir(cfg.assetsRoot, 0o755)
	}
	return nil
}

func getAssetPath(mediaType string) string {
	data := make([]byte, 32)
	_, err := rand.Read(data)
	if err != nil {
		os.Exit(1)
	}
	prefix := base64.RawURLEncoding.EncodeToString(data)
	ext := mediaTypeToExt(mediaType)
	return fmt.Sprintf("%s%s", prefix, ext)
}

func (cfg apiConfig) getAssetDiskPath(assetPath string) string {
	return filepath.Join(cfg.assetsRoot, assetPath)
}

func (cfg apiConfig) getAssetURL(assetPath string) string {
	return fmt.Sprintf("http://localhost:%s/assets/%s", cfg.port, assetPath)
}

func mediaTypeToExt(mediaType string) string {
	splitMediaType := strings.Split(mediaType, "/")
	if len(splitMediaType) != 2 {
		return ".bin"
	}

	return "." + splitMediaType[1]
}
