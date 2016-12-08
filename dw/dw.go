package dw

import (
	"github.com/Bruinxs/uio/dw/replace"
	"github.com/Bruinxs/util"
	"io"
	"sync"
	"time"
)

type DateWriter struct {
	writer   io.WriteCloser
	replacer replace.Replacer
	lock     sync.Locker
}

func NewDateWriter(replacer replace.Replacer, delay time.Duration) *DateWriter {
	if replacer == nil {
		panic("NewDateWriter arg replacer is nil")
	}

	dateWriter := &DateWriter{replacer: replacer, lock: &sync.Mutex{}}
	go dateWriter.runReset(delay)
	return dateWriter
}

func (dw *DateWriter) Write(p []byte) (int, error) {
	dw.lock.Lock()
	defer dw.lock.Unlock()
	if dw.writer == nil {
		var err error
		now := time.Now().Local()
		dw.writer, err = dw.replacer.Replace(now.Year(), int(now.Month()), now.Day())
		if err != nil {
			return 0, err
		}
	}

	return dw.writer.Write(p)
}

func (dw *DateWriter) Close() error {
	dw.lock.Lock()
	defer dw.lock.Unlock()
	if dw.writer != nil {
		return dw.writer.Close()
	}
	return nil
}

func (dw *DateWriter) runReset(delay time.Duration) {
	util.Alarm(int64(delay), int64(24*time.Hour), -1, func(i int) bool {
		dw.lock.Lock()
		defer dw.lock.Unlock()

		if dw.writer != nil {
			dw.writer.Close()
		}
		dw.writer = nil

		return false
	})
}
