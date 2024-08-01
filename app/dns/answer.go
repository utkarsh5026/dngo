package dns

import "encoding/binary"

// Answer represents a DNS answer as defined in RFC 1035.
// A DNS answer is used in the answer section of a DNS response
// and provides information about the query that was made.
//
// Fields:
//
// - Name: The domain name to which this resource record pertains. This is a fully qualified domain name (FQDN).
//
// - Type: A two-octet code which specifies the type of the resource record. The following values are defined:
//   - 1: A (IPv4 address)
//   - 2: NS (authoritative name server)
//   - 5: CNAME (canonical name for an alias)
//   - 6: SOA (start of a zone of authority)
//   - 12: PTR (domain name pointer)
//   - 15: MX (mail exchange)
//   - 28: AAAA (IPv6 address)
//   - 255: ANY (any type)
//
// - Class: A two-octet code that specifies the class of the resource record. The following values are defined:
//   - 1: IN (Internet)
//   - 3: CH (Chaos)
//   - 4: HS (Hesiod)
//   - 255: ANY (any class)
//
// - TTL: A 32-bit unsigned integer that specifies the time interval (in seconds) that the resource record may be cached
// before it should be discarded.
//
// - RDLength: An unsigned 16-bit integer that specifies the length in octets of the RData field.
//
// - RData: A variable length string of octets that describes the resource.
// The format of this information varies according to the Type and Class of the resource record.
type Answer struct {
	Name     string
	Type     uint16
	Class    uint16
	TTL      uint32
	RDLength uint16
	RData    []byte
}

func (a *Answer) Marshal() []byte {
	encodedName := EncodeLabel(a.Name)

	encoded := make([]byte, len(encodedName)+10+len(a.RData))
	offset := 0

	copy(encoded[offset:], encodedName)
	offset += len(encodedName)

	binary.BigEndian.PutUint16(encoded[offset:], a.Type)
	offset += 2

	binary.BigEndian.PutUint16(encoded[offset:], a.Class)
	offset += 2

	binary.BigEndian.PutUint32(encoded[offset:], a.TTL)
	offset += 4

	binary.BigEndian.PutUint16(encoded[offset:], a.RDLength)
	offset += 2

	copy(encoded[offset:], a.RData)

	return encoded
}
