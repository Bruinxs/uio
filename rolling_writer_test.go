package uio

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func tmpDir() string {
	dir := os.Getenv("tmp")
	if dir == "" {
		dir = "/tmp"
	}
	return path.Join(dir, "rolling_wirter_test")
}

func TestDateFileWriter(t *testing.T) {
	now := time.Now()
	dir := tmpDir()
	defer os.RemoveAll(dir)

	w := NewDateFileWriter(dir, "2006-01-02.log")
	defer w.Close()

	_, err := w.Write([]byte("test"))
	if err != nil {
		t.FailNow()
	}

	info, err := os.Stat(filepath.Join(dir, fmt.Sprintf("%04d-%02d-%02d.log", now.Year(), now.Month(), now.Day())))
	if err != nil {
		t.FailNow()
	}
	if info.Size() == 0 {
		t.Error("file size is zero")
	}
}

func BenchmarkDateFileWriter(b *testing.B) {
	dir := tmpDir()
	defer os.RemoveAll(dir)

	w := NewDateFileWriter(dir, "2006-01-02.log")
	defer w.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w.Write([]byte(strings.Repeat("test", 100)))
		}
	})
}
