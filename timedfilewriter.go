package gobuffer

import (
	"fmt"
	"os"
	"time"
)

type TimedFileWriter struct {
	dir      string
	fnFormat string
}

func (m *TimedFileWriter) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	fn := fmt.Sprintf("%s/%s", m.dir, time.Now().Format(m.fnFormat))
	fd, e := os.OpenFile(fn, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if e != nil {
		return 0, e
	}
	defer fd.Close()
	return fd.Write(p)
}

func NewTimedFileWriterWithGoBuffer(bufCap int, flushIntervalSecond int, saveDir, fileNameFormat string) (buf *GoBuffer, e error) {
	fileWriter := &TimedFileWriter{dir: saveDir, fnFormat: fileNameFormat}
	ret, err := NewGoBuffer(bufCap, fileWriter, flushIntervalSecond)
	if err == nil {
		ret.Start()
	}
	return ret, err
}
