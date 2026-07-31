package detectors

import (
	"io/fs"
	"path/filepath"
	"strings"
)

func WalkFiles(root string, visit func(path string) error) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		name := d.Name()
		if d.IsDir() {
			switch name {
			case ".git", "node_modules", ".venv", "vendor", "dist", "build":
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(path, ".py") || strings.HasSuffix(path, ".json") || strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
			return visit(path)
		}
		return nil
	})
}
