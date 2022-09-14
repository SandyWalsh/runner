package sprinter

import (
	"io"
	"sync"
)

type SafeBuffer struct {
	mtx     sync.RWMutex
	newData []chan int // notify readers that new data is available
	done    bool       // set on close
	data    []byte
}

func (b *SafeBuffer) Close() error {
	b.mtx.Lock()
	b.done = true
	b.mtx.Unlock()

	// inform the readers we're done
	b.mtx.RLock()
	for _, ch := range b.newData {
		ch <- 0
	}
	b.mtx.RUnlock()
	return nil
}

func (b *SafeBuffer) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	b.mtx.Lock()
	b.data = append(b.data, p...)
	b.mtx.Unlock()

	// inform the readers how much new data we have
	b.mtx.RLock()
	for _, ch := range b.newData {
		ch <- len(p)
	}
	b.mtx.RUnlock()
	return len(p), nil
}

func (b *SafeBuffer) NewReader() (io.Reader, <-chan int) {
	b.mtx.Lock()
	defer b.mtx.Unlock()
	b.newData = append(b.newData, make(chan int))
	return &safeReader{sb: b}, b.newData[len(b.newData)-1]
}

type safeReader struct {
	sb  *SafeBuffer
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
		if r.sb.done {
			return 0, io.EOF
		}
		return 0, nil
	}
	r.off += n
	return n, nil
}

var _ io.Writer = (*SafeBuffer)(nil)
var _ io.Reader = (*safeReader)(nil)
