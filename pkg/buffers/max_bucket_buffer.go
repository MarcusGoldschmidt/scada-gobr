package buffers

import (
	"bytes"
)

type MaxBucketBuffer struct {
	buckets       []*bytes.Buffer
	bucketSize    int
	maxSize       uint64
	currentBuffer *bytes.Buffer
}

func NewMaxBucketBuffer(bucketSize int, maxSize ByteSize) *MaxBucketBuffer {
	buckets := []*bytes.Buffer{bytes.NewBuffer(make([]byte, 0, bucketSize))}

	maxSizeInt := uint64(maxSize) * uint64(bucketSize)

	return &MaxBucketBuffer{
		buckets:       buckets,
		currentBuffer: buckets[0],
		bucketSize:    bucketSize,
		maxSize:       maxSizeInt,
	}
}

func (b *MaxBucketBuffer) Write(p []byte) (n int, err error) {
	if b.currentBuffer.Len() >= b.bucketSize {
		newBuffer := bytes.NewBuffer(make([]byte, 0, b.bucketSize))
		b.buckets = append(b.buckets, newBuffer)
		b.currentBuffer = newBuffer

		if b.Len() > b.maxSize {
			b.buckets = b.buckets[1:]
		}
	}

	return b.currentBuffer.Write(p)
}

func (b *MaxBucketBuffer) Read(p []byte) (n int, err error) {
	index := 0
	for _, bucket := range b.buckets {
		copy(p[index:index+bucket.Len()], bucket.Bytes())
		index += bucket.Len()
	}

	return index, nil
}

func (b *MaxBucketBuffer) Len() (sum uint64) {
	for _, bucket := range b.buckets {
		sum += uint64(bucket.Len())
	}

	return
}
