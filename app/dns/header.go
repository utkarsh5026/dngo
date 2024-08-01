package dns

import "encoding/binary"

const (
	HeaderSize = 12
)

// Header represents the DNS packet header as defined in RFC 1035.
// The DNS header is 12 bytes long and contains various fields that
// provide information about the DNS query or response.
//
// Fields:
//
// - ID: A 16-bit identifier assigned by the program that generates any kind of query. This identifier is copied in the corresponding reply and can be used by the requester to match up replies to outstanding queries.
//
// - QR: A one-bit field that specifies whether this message is a query (0) or a response (1).
//
// - OpCode: A four-bit field that specifies the kind of query in this message. This value is set by the originator of a query and copied into the response. The following values are defined:
//   - 0: Standard query (QUERY)
//   - 1: Inverse query (IQUERY)
//   - 2: Server status request (STATUS)
//   - 3-15: Reserved for future use
//
// - AA: Authoritative Answer - this bit is valid in responses, and specifies that the responding name server is an authority for the domain name in question section.
//
// - TC: TrunCation - specifies that this message was truncated due to length greater than that permitted on the transmission channel.
//
// - RD: Recursion Desired - this bit may be set in a query and is copied into the response. If RD is set, it directs the name server to pursue the query recursively.
//
// - RA: Recursion Available - this bit is set or cleared in a response, and denotes whether recursive query support is available in the name server.
//
// - Z: Reserved for future use. Must be zero in all queries and responses.
//
// - RCode: Response code - this 4-bit field is set as part of responses. The following values are defined:
//   - 0: No error condition
//   - 1: Format error - The name server was unable to interpret the query.
//   - 2: Server failure - The name server was unable to process this query due to a problem with the name server.
//   - 3: Name Error - Meaningful only for responses from an authoritative name server, this code signifies that the domain name referenced in the query does not exist.
//   - 4: Not Implemented - The name server does not support the requested kind of query.
//   - 5: Refused - The name server refuses to perform the specified operation for policy reasons.
//   - 6-15: Reserved for future use
//
// - QDCount: An unsigned 16-bit integer specifying the number of entries in the question section.
//
// - ANCount: An unsigned 16-bit integer specifying the number of resource records in the answer section.
//
// - NSCount: An unsigned 16-bit integer specifying the number of name server resource records in the authority records section.
//
// - ARCount: An unsigned 16-bit integer specifying the number of resource records in the additional records section.
type Header struct {
	ID      uint16
	QR      bool
	OpCode  uint8
	AA      bool
	TC      bool
	RD      bool
	RA      bool
	Z       uint8
	RCode   uint8
	QDCount uint16
	ANCount uint16
	NSCount uint16
	ARCount uint16
}

func (h *Header) Marshal() []byte {
	encoded := make([]byte, HeaderSize) // Assuming HeaderSize is 12

	binary.BigEndian.PutUint16(encoded[0:2], h.ID)

	encoded[2] = uint8(h.OpCode) << 3
	if h.QR {
		encoded[2] |= 1 << 7
	}
	if h.AA {
		encoded[2] |= 1 << 2
	}
	if h.TC {
		encoded[2] |= 1 << 1
	}
	if h.RD {
		encoded[2] |= 1
	}

	encoded[3] = h.Z << 4
	if h.RA {
		encoded[3] |= 1 << 7
	}
	encoded[3] |= h.RCode & 0xF

	binary.BigEndian.PutUint16(encoded[4:6], h.QDCount)
	binary.BigEndian.PutUint16(encoded[6:8], h.ANCount)
	binary.BigEndian.PutUint16(encoded[8:10], h.NSCount)
	binary.BigEndian.PutUint16(encoded[10:12], h.ARCount)

	return encoded
}

func UnmarshalHeader(encoded []byte) *Header {
	return &Header{
		ID:      binary.BigEndian.Uint16(encoded[0:2]),
		QR:      encoded[2]&(1<<7) != 0,
		OpCode:  (encoded[2] >> 3) & 0xF,
		AA:      encoded[2]&(1<<2) != 0,
		TC:      encoded[2]&(1<<1) != 0,
		RD:      encoded[2]&1 != 0,
		RA:      encoded[3]&(1<<7) != 0,
		Z:       (encoded[3] >> 4) & 0xF,
		RCode:   encoded[3] & 0xF,
		QDCount: binary.BigEndian.Uint16(encoded[4:6]),
		ANCount: binary.BigEndian.Uint16(encoded[6:8]),
		NSCount: binary.BigEndian.Uint16(encoded[8:10]),
		ARCount: binary.BigEndian.Uint16(encoded[10:12]),
	}
}
