package message

import (
	"bytes"
	"testing"
)

func TestTopicBytes(t *testing.T) {
	bytes := topicBytes(5)
	if bytes[0] != 5 {
		t.Errorf("First byte was %d; want 5", bytes[0])
	}
	if bytes[1] != 0 {
		t.Errorf("First byte was %d; want 0", bytes[1])
	}
	bytes = topicBytes(65000)
	if bytes[0] != 232 {
		t.Errorf("First byte was %d; want 5", bytes[0])
	}
	if bytes[1] != 253 {
		t.Errorf("First byte was %d; want 0", bytes[1])
	}
}

func TestCursorBytes(t *testing.T) {
	_bytes := cursorBytes(18440000000000000000)
	correctBytes := []byte{0, 0, 52, 250, 76, 10, 232, 255}
	if len(_bytes) != 8 {
		t.Errorf("Length was %d; want 8", len(_bytes))
	}
	if bytes.Compare(_bytes, correctBytes) != 0 {
		t.Errorf("Values were %b; want %b", _bytes, correctBytes)
	}
}

func TestPayloadLengthBytes(t *testing.T) {
	_bytes, err := payloadLengthBytes(0)
	if len(_bytes) != 1 {
		t.Errorf("Length was %d; want 1", len(_bytes))
	}
	if _bytes[0] != 0 {
		t.Errorf("First byte was %b; want %b", _bytes[0], 0)
	}

	_bytes, err = payloadLengthBytes(60)
	if len(_bytes) != 1 {
		t.Errorf("Length was %d; want 1", len(_bytes))
	}
	if _bytes[0] != 60 {
		t.Errorf("First byte was %b; want %b", _bytes[0], 60)
	}

	_bytes, err = payloadLengthBytes(128 + 3)
	correctBytes := []byte{0b10000011, 0b1}
	if len(_bytes) != 2 {
		t.Errorf("Length was %d; want 2", len(_bytes))
	}
	if bytes.Compare(_bytes, correctBytes) != 0 {
		t.Errorf("Values were %b; want %b", _bytes, correctBytes)
	}

	_bytes, err = payloadLengthBytes(268_435_455)
	correctBytes = []byte{0b11111111, 0b11111111, 0b11111111, 0b01111111}
	if len(_bytes) != 4 {
		t.Errorf("Length was %d; want 4", len(_bytes))
	}
	if bytes.Compare(_bytes, correctBytes) != 0 {
		t.Errorf("Values were %b; want %b", _bytes, correctBytes)
	}

	_bytes, err = payloadLengthBytes(268_435_456)
	correctBytes = []byte{0b0}
	expectedError := "message: payload length 2.68435456e+08 over limit of 2^28 bytes"
	if err != nil {
		if err.Error() != expectedError {
			t.Errorf("Error was '%s'; want '%s'", err.Error(), expectedError)
		}
	} else {
		t.Errorf("Values were %b; want error", _bytes)
	}
	if bytes.Compare(_bytes, correctBytes) != 0 {
		t.Errorf("Values were %b; want %b", _bytes, correctBytes)
	}
}
