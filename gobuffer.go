package gobuffer

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"time"
)

type GoBufferErrorBlockTooBig struct {
}

func (m *GoBufferErrorBlockTooBig) Error() string {
	return "the block be going to write to buffer is too big"
}

type GoBufferErrorBufferPaused struct {
}

func (m *GoBufferErrorBufferPaused) Error() string {
	return "gobuffer was paused"
}

type GoBuffer struct {
	sync.Mutex
	buf1          *bytes.Buffer
	buf2          *bytes.Buffer
	activeBuf     *bytes.Buffer
	out           io.Writer
	flushInterval time.Duration
	stopped       int
	freeBufChan   chan int
}

const _max_buf_cap int = 1024 * 1024 * 1024

func NewGoBuffer(bufCap int, output io.Writer, flushIntervalSecond int) (buf *GoBuffer, e error) {
	if bufCap > _max_buf_cap || bufCap < 128 || output == nil || flushIntervalSecond < 1 {
		return nil, fmt.Errorf("NewGoBuffer fail, invalid params")
	}
	ret := new(GoBuffer)
	ret.out = output
	ret.stopped = 1
	ret.flushInterval = time.Second * time.Duration(flushIntervalSecond)
	ret.freeBufChan = make(chan int, 1)
	ret.freeBufChan <- 1
	ret.buf1 = new(bytes.Buffer)
	ret.buf1.Grow(bufCap)
	ret.buf2 = new(bytes.Buffer)
	ret.buf2.Grow(bufCap)
	ret.activeBuf = ret.buf1
	return ret, nil
}

func (m *GoBuffer) Write(p []byte) (n int, err error) {
	if m.stopped != 0 {
		return 0, new(GoBufferErrorBufferPaused)
	}
	if len(p) == 0 {
		return 0, nil
	}
	if len(p) > m.activeBuf.Cap() {
		return 0, new(GoBufferErrorBlockTooBig)
	}

	m.Lock()
	defer m.Unlock()
	return m.write(p)
}

func (m *GoBuffer) write(p []byte) (n int, err error) {
	retn := 0
	remain := m.activeBuf.Cap() - m.activeBuf.Len()
	doSwap := false
	if remain == 0 {
		doSwap = true
	} else if remain < len(p) {
		ptmp := p[:remain]
		pos := 0
		doSwap = true
		for {
			nn, e := m.activeBuf.Write(ptmp[pos:])
			retn += nn
			pos += nn
			remain -= nn
			if e != nil {
				return retn, e
			}
			if remain <= 0 || retn == len(p) {
				break
			}
		}
	} else {
		ptmp := p
		pos := 0
		remain := len(p)
		for {
			nn, e := m.activeBuf.Write(ptmp[pos:])
			retn += nn
			pos += nn
			remain -= nn
			if e != nil {
				return retn, e
			}
			if remain <= 0 || retn == len(p) {
				break
			}
		}
	}

	if doSwap {
		if e := m.swap(); e != nil {
			return retn, e
		}
	} else {
		return retn, nil
	}

	dataRemain := len(p) - retn
	ptmp := p
	pos := retn

	for {
		nn, e := m.activeBuf.Write(ptmp[pos:])
		retn += nn
		pos += nn
		dataRemain -= nn
		if e != nil {
			return retn, e
		}
		if dataRemain <= 0 || retn == len(p) {
			break
		}
	}

	return retn, nil
}

func (m *GoBuffer) swap() error {
	<-m.freeBufChan
	go m.flush(m.activeBuf)
	if m.activeBuf == m.buf1 {
		m.activeBuf = m.buf2
	} else {
		m.activeBuf = m.buf1
	}
	m.activeBuf.Reset()
	return nil
}

func (m *GoBuffer) flush(buf *bytes.Buffer) bool {
	defer func() {
		m.freeBufChan <- 1
	}()
	flushBufPos := 0
	remain := buf.Len()
	for {
		n, e := m.out.Write(buf.Bytes()[flushBufPos:])
		remain -= n
		flushBufPos += n
		if remain <= 0 {
			buf.Reset()
			return true
		}
		if n == 0 || e != nil {
			break
		}
	}
	return remain == 0
}

func (m *GoBuffer) Flush() error {
	m.Lock()
	defer m.Unlock()
	return m.swap()
}

func (m *GoBuffer) Start() {
	go func() {
		for {
			select {
			case <-time.After(m.flushInterval):
				//fmt.Println("flushinterval")
				m.Flush()
			}
		}
	}()
	m.stopped = 0
}

func (m *GoBuffer) Pause() {
	m.stopped = 1
}

func (m *GoBuffer) Continue() {
	m.stopped = 0
}
