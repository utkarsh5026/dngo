package dns

import "strings"

// EncodeLabel encodes a domain name label into the DNS label format.
// The DNS label format is a sequence of labels where each label is prefixed
// by its length and the sequence is terminated by a zero-length label.
//
// Parameters:
// - label: The domain name label to encode. This should be a fully qualified domain name (FQDN).
//
// Returns:
// - A byte slice containing the encoded label in DNS format.
func EncodeLabel(label string) []byte {
	parts := strings.Split(label, ".")
	byteCount := 0

	for _, part := range parts {
		byteCount += len(part) + 1
	}

	byteCount++

	buffer := make([]byte, byteCount)
	offset := 0

	for _, part := range parts {
		buffer[offset] = byte(len(part))
		offset++
		copy(buffer[offset:], []byte(part))
		offset += len(part)
	}

	buffer[offset] = 0x00
	offset++

	return buffer
}
