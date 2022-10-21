package buffers

import (
	"testing"
)

func TestSimpleMaxBuffer(t *testing.T) {
	bBuffer := NewMaxBuffer(20)

	_, _ = bBuffer.WriteString("123456")
	_, _ = bBuffer.WriteString("123")
	_, _ = bBuffer.WriteString("1234")

	temp := bBuffer.ReadAll()

	if temp != "1234561231234" {
		t.Errorf("Expected '1234561231234' got %s", temp)
	}
}

func TestResetMaxBuffer(t *testing.T) {
	bBuffer := NewMaxBuffer(2)

	_, _ = bBuffer.WriteString("123456")
	_, _ = bBuffer.WriteString("321")

	temp := bBuffer.ReadAll()

	if temp != "321" {
		t.Errorf("Expected '321' got %s", temp)
	}
}
