package gobuffer

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"
)

type TimedFileWriter struct {
	dir            string
	filePrefix     string
	datetimeFormat string
	fnFormat       string
	latestFileName string
	linedContent   bool
}

func (m *TimedFileWriter) SetContentLined(v bool) {
	m.linedContent = v
}

func (m *TimedFileWriter) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}

	writePos := 0
	if m.linedContent && m.latestFileName != "" {
		oldFd, err := os.OpenFile(m.latestFileName, os.O_RDWR|os.O_APPEND, 0666)
		if err != nil {
			return 0, err
		}
		defer oldFd.Close()
		oldFd.Seek(0, io.SeekEnd)
		linePos := bytes.Index(p, []byte{'\n'})
		if linePos >= 0 {
			oldFd.Write(p[:linePos+1])
			writePos = linePos + 1
		} else {
			oldFd.Write(p)
			return len(p), nil
		}

		m.latestFileName = ""
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

	if len(p) > writePos && p[len(p)-1] != '\n' {
		m.latestFileName = fn
	}
	fd.Write(p[writePos:])
	return len(p), nil

}

func NewTimedFileWriterWithGoBuffer(bufCap int, flushIntervalSecond int, saveDir, fileNameFormat string) (buf *GoBuffer, e error) {
	fileWriter := &TimedFileWriter{dir: saveDir, fnFormat: fileNameFormat}
	fileWriter.SetContentLined(false)
	ret, err := NewGoBuffer(bufCap, fileWriter, flushIntervalSecond)
	if err == nil {
		ret.Start()
	}
	return ret, err
}

func NewTimedFileWriterWithGoBuffer2(bufCap int, flushIntervalSecond int, saveDir, filePrefix, datetimeFormat string) (buf *GoBuffer, e error) {
	fileWriter := &TimedFileWriter{dir: saveDir, filePrefix: filePrefix, datetimeFormat: datetimeFormat}
	fileWriter.SetContentLined(false)
	ret, err := NewGoBuffer(bufCap, fileWriter, flushIntervalSecond)
	if err == nil {
		ret.Start()
	}
	return ret, err
}

func NewLinedTimedFileWriterWithGoBuffer(bufCap int, flushIntervalSecond int, saveDir, filePrefix, datetimeFormat string) (buf *GoBuffer, e error) {
	fileWriter := &TimedFileWriter{dir: saveDir, filePrefix: filePrefix, datetimeFormat: datetimeFormat}
	fileWriter.SetContentLined(true)
	ret, err := NewGoBuffer(bufCap, fileWriter, flushIntervalSecond)
	if err == nil {
		ret.Start()
	}
	return ret, err
}
