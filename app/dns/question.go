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
		copy(buffer[offset:], []byte(part))
		offset += len(part)
	}
	buffer[offset] = 0x00
	offset++

	binary.BigEndian.PutUint16(buffer[offset:offset+2], q.Type)
	offset += 2
	binary.BigEndian.PutUint16(buffer[offset:offset+2], q.Class)

	return buffer, nil
}
