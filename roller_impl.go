package uio

import (
	"io"
	"os"
	"path/filepath"
	"time"
)

// DateFileRoller is a roller to roll by date
type DateFileRoller struct {
	fileDir         string
	fileNameFormat  string
	currectFileName string
	currectFile     *os.File
	nextFileName    string
}

// NewDateFileRoller return a new DateFileRoller
// format is time format such as 2006-01-02.log
func NewDateFileRoller(dir, format string) *DateFileRoller {
	os.MkdirAll(dir, 0771)
	return &DateFileRoller{fileDir: dir, fileNameFormat: format}
}

func (d *DateFileRoller) newFileName() string {
	return time.Now().Format(d.fileNameFormat)
}

func (d *DateFileRoller) createFile(fileName string) (*os.File, error) {
	return os.OpenFile(filepath.Join(d.fileDir, fileName), os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
}

func (d *DateFileRoller) exist(fileName string) bool {
	_, err := os.Stat(filepath.Join(d.fileDir, fileName))
	if err == nil {
		return true
	}
	return os.IsExist(err)
}

// NeedRoll decide when to roll
func (d *DateFileRoller) NeedRoll() bool {
	d.nextFileName = d.newFileName()
	return d.currectFileName != d.nextFileName
}

// Roll return new writer
func (d *DateFileRoller) Roll() (_ io.Writer, err error) {
	if d.currectFile != nil {
		err = d.currectFile.Close()
		if err != nil {
			return nil, err
		}
	}

	d.currectFileName = d.nextFileName
	d.currectFile, err = d.createFile(d.currectFileName)
	return d.currectFile, err
}

// Close close the file if be holded
func (d *DateFileRoller) Close() error {
	if d.currectFile != nil {
		err := d.currectFile.Close()
		if err != nil {
			return err
		}
		d.currectFile = nil
	}
	return nil
}
