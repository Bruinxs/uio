package replace

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type FileReplacer struct {
	dir      string
	filename string
}

func NewFileReplacer(filename string) *FileReplacer {
	if !strings.Contains(filename, "%v") {
		panic("[NewFileReplacer] filename must contains a symbol '%v'")
	}
	fr := &FileReplacer{}
	fr.dir = filepath.Dir(filename)
	fr.filename = filename
	return fr
}

func (fr *FileReplacer) Replace(year int, month int, day int) (io.WriteCloser, error) {
	_, err := os.Stat(fr.dir)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(fr.dir, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	return os.OpenFile(fmt.Sprintf(fr.filename, fmt.Sprintf("%v-%v-%v", year, month, day)), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
}
