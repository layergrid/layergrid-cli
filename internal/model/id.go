package model

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"strings"
)

func StableID(parts ...string) string {
	h := sha256.New()
	for _, p := range parts {
		_, _ = h.Write([]byte(p))
		_, _ = h.Write([]byte{0})
	}
	return hex.EncodeToString(h.Sum(nil))[:16]
}

func RelativeLocation(root, path string, line int) Location {
	rel, err := filepath.Rel(root, path)
	if err != nil || strings.HasPrefix(rel, "..") {
		rel = path
	}
	return Location{Path: filepath.ToSlash(rel), Line: line}
}

func Descriptor(parts ...string) string {
	return fmt.Sprintf("sha256:%s", StableID(parts...))
}
