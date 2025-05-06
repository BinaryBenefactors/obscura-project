package utils

import (
	"fmt"
	"path/filepath"
	"time"
)

var (
	AllowedTypes = []string{
		"image/jpeg",
		"image/png",
		"video/mp4",
		"video/quicktime",
	}

	AllowedExtensions = []string{
		".jpg",
		".jpeg",
		".png",
		".mp4",
		".mov",
	}
)

func IsAllowedFileType(fileType string) bool {
	for _, t := range AllowedTypes {
		if fileType == t {
			return true
		}
	}
	return false
}

func IsAllowedExtension(ext string) bool {
	for _, e := range AllowedExtensions {
		if ext == e {
			return true
		}
	}
	return false
}

func GenerateFilename(original string) string {
	ext := filepath.Ext(original)
	return fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
}

func GenerateFileID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
