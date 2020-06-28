package message

import (
	"encoding/binary"
	"fmt"
)

func topicBytes(length uint16) []byte {
	array := make([]byte, 2)
	binary.LittleEndian.PutUint16(array, length)
	return array
}

func cursorBytes(cursor uint64) []byte {
	array := make([]byte, 8)
	binary.LittleEndian.PutUint64(array, cursor)
	return array
}

// PayloadOverLengthError Error indicates the length of the payload was over the allowed length.
type PayloadOverLengthError uint64

func (f PayloadOverLengthError) Error() string {
	return fmt.Sprintf("message: payload length %g over limit of 2^28 bytes", float64(f))
}

const continuationBit = uint8(1 << 7)

func payloadLengthBytes(lengthOfRest uint64) ([]byte, error) {
	// Quick exist for 0 length with at least 1 byte.
	// The other logic would result in an empty array.
	if lengthOfRest == 0 {
		return make([]byte, 1), nil
	}
	array := make([]byte, 4)
	index := 0
	for bit := lengthOfRest; bit > 0; bit = bit >> 7 {
		if index < 4 {
			if index > 0 {
				// Set continuation bit in previous byte
				array[index-1] = array[index-1] | continuationBit
			}
			// Mask first 7 bits and set as byte
			value := bit & 0b1111111
			array[index] = uint8(value)
		}
		index++
	}
	// Check if length required more than 4 bytes to fill
	if index > 4 {
		return make([]byte, 1), PayloadOverLengthError(lengthOfRest)
	}
	return array[0:index], nil
}
