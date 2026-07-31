package scan

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func BenchmarkRunManyFiles(b *testing.B) {
	root := b.TempDir()
	for i := 0; i < 1000; i++ {
		path := filepath.Join(root, fmt.Sprintf("file_%04d.py", i))
		if err := os.WriteFile(path, []byte("print('hello')\n"), 0o644); err != nil {
			b.Fatal(err)
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := Run(root, Options{}); err != nil {
			b.Fatal(err)
		}
	}
}
