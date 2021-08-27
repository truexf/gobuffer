package gobuffer

import (
	"fmt"
	"os"
	"time"
)

type TimedFileWriter struct {
	dir            string
	filePrefix     string
	datetimeFormat string
	fnFormat       string
}

func (m *TimedFileWriter) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	fn := ""
	if m.fnFormat != "" {
		fn = fmt.Sprintf("%s/%s", m.dir, time.Now().Format(m.fnFormat))
	} else {
		fn = fmt.Sprintf("%s/%s%s", m.dir, m.filePrefix, time.Now().Format(m.datetimeFormat))
	}
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

func NewTimedFileWriterWithGoBuffer2(bufCap int, flushIntervalSecond int, saveDir, filePrefix, datetimeFormat string) (buf *GoBuffer, e error) {
	fileWriter := &TimedFileWriter{dir: saveDir, filePrefix: filePrefix, datetimeFormat: datetimeFormat}
	ret, err := NewGoBuffer(bufCap, fileWriter, flushIntervalSecond)
	if err == nil {
		ret.Start()
	}
	return ret, err
}
