package pyutil

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func Walk(root string, include, exclude []string, visit func(path string, lines []string) error) error {
	return filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			switch d.Name() {
			case ".git", ".venv", "node_modules", "vendor", "dist", "build":
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".py") {
			return nil
		}
		if !Included(root, path, include, exclude) {
			return nil
		}
		lines, err := ReadLines(path)
		if err != nil {
			return err
		}
		return visit(path, lines)
	})
}

func Included(root, path string, include, exclude []string) bool {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		rel = path
	}
	rel = filepath.ToSlash(rel)
	for _, pattern := range exclude {
		if matchPattern(pattern, rel) {
			return false
		}
	}
	if len(include) == 0 {
		return true
	}
	for _, pattern := range include {
		if matchPattern(pattern, rel) {
			return true
		}
	}
	return false
}

func matchPattern(pattern, rel string) bool {
	pattern = filepath.ToSlash(strings.TrimPrefix(pattern, "./"))
	if strings.HasSuffix(pattern, "/**") {
		return strings.HasPrefix(rel, strings.TrimSuffix(pattern, "/**")+"/")
	}
	if strings.HasPrefix(pattern, "**/") {
		suffix := strings.TrimPrefix(pattern, "**/")
		if ok, _ := filepath.Match(suffix, filepath.Base(rel)); ok {
			return true
		}
		return strings.HasSuffix(rel, "/"+suffix)
	}
	if ok, _ := filepath.Match(pattern, rel); ok {
		return true
	}
	if ok, _ := filepath.Match(pattern, filepath.Base(rel)); ok {
		return true
	}
	return rel == pattern
}

func ReadLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func HasAny(lines []string, needles ...string) bool {
	joined := strings.ToLower(strings.Join(lines, "\n"))
	for _, needle := range needles {
		if strings.Contains(joined, strings.ToLower(needle)) {
			return true
		}
	}
	return false
}

func AssignmentName(line, fallback string) string {
	re := regexp.MustCompile(`^\s*([A-Za-z_][A-Za-z0-9_]*)\s*=`)
	if m := re.FindStringSubmatch(line); len(m) == 2 {
		return m[1]
	}
	return fallback
}

func KeywordString(line, key, fallback string) string {
	re := regexp.MustCompile(key + `\s*=\s*["']([^"']+)["']`)
	if m := re.FindStringSubmatch(line); len(m) == 2 {
		return m[1]
	}
	return fallback
}

func Block(lines []string, start, window int) string {
	end := start + window
	if end > len(lines) {
		end = len(lines)
	}
	return strings.Join(lines[start:end], "\n")
}
