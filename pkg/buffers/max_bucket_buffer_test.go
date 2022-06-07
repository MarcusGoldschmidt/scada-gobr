package buffers

import (
	"testing"
)

func TestSimpleMaxBucketBuffer(t *testing.T) {
	bBuffer := NewMaxBucketBuffer(2, 20)

	_, _ = bBuffer.Write([]byte("123456"))
	_, _ = bBuffer.Write([]byte("123"))
	_, _ = bBuffer.Write([]byte("1234"))

	if len(bBuffer.buckets) != 3 {
		t.Error("Expected 3 buckets")
	}

	dest := make([]byte, bBuffer.Len())

	_, _ = bBuffer.Read(dest)

	if string(dest) != "1234561231234" {
		t.Error("Expected '1234561231234'")
	}
}
