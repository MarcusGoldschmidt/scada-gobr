package buffers

import (
	"bytes"
	"sync"
)

type MaxBuffer struct {
	sync.RWMutex
	buffer  *bytes.Buffer
	maxSize uint64
}

func NewMaxBuffer(maxSize ByteSize) *MaxBuffer {
	buffer := bytes.NewBuffer(make([]byte, 0, maxSize))

	return &MaxBuffer{
		buffer:  buffer,
		maxSize: uint64(maxSize),
		RWMutex: sync.RWMutex{},
	}
}

func (b *MaxBuffer) WriteString(value string) (int, error) {
	b.Lock()
	defer b.Unlock()

	if b.buffer.Len()+len(value) >= int(b.maxSize) {
		b.buffer.Reset()
	}

	return b.buffer.WriteString(value)
}

func (b *MaxBuffer) Write(p []byte) (n int, err error) {
	return b.WriteString(string(p))
}

func (b *MaxBuffer) ReadAll() string {
	b.RLock()
	defer b.RUnlock()

	return b.buffer.String()
}

func (b *MaxBuffer) Read(lines int) string {
	b.RLock()
	defer b.RUnlock()

	result := ""

	for {
		if lines <= 0 {
			break
		}

		readString, err := b.buffer.ReadString('\n')
		if err != nil {
			break
		}

		result = result + readString

		lines--
	}

	return result
}

func (b *MaxBuffer) Len() int {
	b.RLock()
	defer b.RUnlock()

	return b.buffer.Len()
}
