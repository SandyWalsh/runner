package library

import (
	"io"
	"sync"
)

type safeBuffer struct {
	mtx  sync.RWMutex
	data []byte
}

func (b *safeBuffer) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	b.mtx.Lock()
	defer b.mtx.Unlock()

	b.data = append(b.data, p...)
	return len(p), nil
}

func (b *safeBuffer) NewReader() io.Reader {
	return &safeReader{sb: b}
}

type safeReader struct {
	sb  *safeBuffer
	off int
}

func (r *safeReader) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	r.sb.mtx.RLock()
	defer r.sb.mtx.RUnlock()
	n = copy(p, r.sb.data[r.off:])
	if n == 0 {
		return 0, io.EOF
	}
	r.off += n
	return n, nil
}

var _ io.Writer = (*safeBuffer)(nil)
var _ io.Reader = (*safeReader)(nil)
