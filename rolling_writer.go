package uio

import (
	"io"
	"sync"
)

// Roller define when to roll and how to roll,
// the roller should be close last writer when roll next writer
// call `NeedRoll` firstly should be return true
type Roller interface {
	io.Closer
	NeedRoll() bool
	Roll() (io.Writer, error)
}

// RollingWriter is a rolling holder
// it should be creater by call `NewRollingWriter`
type RollingWriter struct {
	roller Roller
	writer io.Writer
	lock   *sync.Mutex
}

// NewRollingWriter return a RollingWriter
func NewRollingWriter(roller Roller) *RollingWriter {
	return &RollingWriter{roller: roller, lock: &sync.Mutex{}}
}

func (r *RollingWriter) Write(p []byte) (n int, err error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.roller.NeedRoll() {
		r.writer, err = r.roller.Roll()
		if err != nil {
			return -1, err
		}
	}
	return r.writer.Write(p)
}

// Close just call roller Close safely
func (r *RollingWriter) Close() error {
	r.lock.Lock()
	defer r.lock.Unlock()

	return r.roller.Close()
}

// =============================================================================================
//      rolling writer instance
// =============================================================================================

// NewDateFileWriter roll wirte to file by date
func NewDateFileWriter(dir, format string) *RollingWriter {
	return NewRollingWriter(NewDateFileRoller(dir, format))
}
