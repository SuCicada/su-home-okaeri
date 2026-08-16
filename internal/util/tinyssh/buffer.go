package tinyssh

import (
	"bytes"
	"sync"
)

// syncBuffer 让 exec 的输出可以在 Wait 的 goroutine 之外安全读取。
type syncBuffer struct {
	lock sync.Mutex
	buf  bytes.Buffer
}

func (b *syncBuffer) Write(p []byte) (int, error) {
	b.lock.Lock()
	defer b.lock.Unlock()
	return b.buf.Write(p)
}

func (b *syncBuffer) String() string {
	b.lock.Lock()
	defer b.lock.Unlock()
	return b.buf.String()
}

func (b *syncBuffer) Len() int {
	b.lock.Lock()
	defer b.lock.Unlock()
	return b.buf.Len()
}
