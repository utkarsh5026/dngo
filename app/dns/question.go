package dns

import (
	"encoding/binary"
	"fmt"
	"strings"
)

// Question represents a DNS question as defined in RFC 1035.
// A DNS question is used in the question section of a DNS query
// and specifies the domain name and the type of query being made.
//
// Fields:
//
// - Name: The domain name for which the query is being made. This is a fully qualified domain name (FQDN).
//
// - Type: A two-octet code which specifies the type of the query. The following values are defined:
//   - 1: A (IPv4 address)
//   - 2: NS (authoritative name server)
//   - 5: CNAME (canonical name for an alias)
//   - 6: SOA (start of a zone of authority)
//   - 12: PTR (domain name pointer)
//   - 15: MX (mail exchange)
//   - 28: AAAA (IPv6 address)
//   - 255: ANY (any type)
//
// - Class: A two-octet code that specifies the class of the query. The following values are defined:
//   - 1: IN (Internet)
//   - 3: CH (Chaos)
//   - 4: HS (Hesiod)
//   - 255: ANY (any class)
type Question struct {
	Name  string
	Type  uint16
	Class uint16
}

// Marshal encodes the DNS Question into a byte slice.
// It converts the Name, Type, and Class fields of the Question into their byte representations
// and concatenates them into a single byte slice.
//
// Returns:
// - A byte slice containing the encoded DNS Question.
// - An error if any issue occurs during encoding, such as a label exceeding 63 bytes.
func (q *Question) Marshal() ([]byte, error) {
	parts := strings.Split(q.Name, ".")

	byteCount := 0
	for _, part := range parts {
		if len(part) > 63 {
			return nil, fmt.Errorf("label '%s' exceeds 63 bytes", part)
		}
		byteCount += len(part) + 1
	}
	byteCount++ // for the last 0x00

	buffer := make([]byte, byteCount+4)

	offset := 0
	for _, part := range parts {
		buffer[offset] = byte(len(part))
		offset++
		copy(buffer[offset:], part)
		offset += len(part)
	}
	buffer[offset] = 0x00
	offset++

	binary.BigEndian.PutUint16(buffer[offset:offset+2], q.Type)
	offset += 2
	binary.BigEndian.PutUint16(buffer[offset:offset+2], q.Class)

	return buffer, nil
}

// UnmarshalQuestions decodes a byte slice into a slice of DNS Question structs.
// It reads the specified number of questions from the DNS message.
//
// Parameters:
// - dnsMessage: A byte slice containing the encoded DNS message.
// - count: The number of questions to unmarshal.
//
// Returns:
// - A slice of Question structs populated with the decoded values.
//
// - An error if any issue occurs during decoding.
func UnmarshalQuestions(dnsMessage []byte, count uint16) ([]Question, error) {
	offset := HeaderSize
	questions := make([]Question, 0, count)

	for i := 0; i < int(count); i++ {
		// Check if we have enough bytes left to read a question
		typeClassByteCount := 4
		if offset+typeClassByteCount > len(dnsMessage) {
			return nil, fmt.Errorf("incomplete question at offset %d", offset)
		}

		label, bytesRead := parseLabel(dnsMessage[offset:], dnsMessage)
		offset += bytesRead

		fmt.Println("label parsed", label)
		question := Question{
			Name:  label,
			Type:  1,
			Class: 1,
		}

		questions = append(questions, question)
		offset += typeClassByteCount // type + class
	}

	return questions, nil
}

// parseLabel decodes a DNS label from a byte slice.
// It handles both uncompressed and compressed labels as specified in RFC 1035.
//
// Parameters:
// - label: A byte slice containing the encoded DNS label.
// - dnsMessage: A byte slice containing the entire DNS message, used for resolving compressed labels.
//
// Returns:
//
// - A string representing the decoded domain name.
//
// - The number of bytes read from the label.
func parseLabel(label []byte, dnsMessage []byte) (string, int) {
	offset := 0
	var parts []string

	// Check for empty label
	if len(label) == 0 {
		return "", 0
	}

	// Main parsing loop
	for offset < len(label) && label[offset] != 0 {
		if label[offset]&0xC0 == 0xC0 {
			if offset+2 > len(label) {
				break
			}
			pointer := binary.BigEndian.Uint16(label[offset:offset+2]) & 0x3FFF
			compressedLabel, _ := parseLabel(dnsMessage[pointer:], dnsMessage)
			parts = append(parts, compressedLabel)
			offset += 2
			break // We should stop after a compression pointer
		}

		length := int(label[offset])
		if offset+1+length > len(label) {
			break
		}

		labelPart := string(label[offset+1 : offset+1+length])
		parts = append(parts, labelPart)
		offset += length + 1
	}

	// Increment offset for the zero byte, if present
	if offset < len(label) && label[offset] == 0 {
		offset++
	}

	return strings.Join(parts, "."), offset
}
