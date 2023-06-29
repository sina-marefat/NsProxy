package dns

import "encoding/binary"

func GenerateErrorResponse(id []byte, error error, ns string) []byte {
	resp := sampleResponse(id, ns)
	switch error {
	case ErrServerError:
		binary.BigEndian.PutUint16(resp[2:4], 0x8183)
	case ErrUnspportedType:
		binary.BigEndian.PutUint16(resp[2:4], 0x8185)
	default:
		binary.BigEndian.PutUint16(resp[2:4], 0x8181)
	}
	return resp

}

func sampleResponse(id []byte, ns string) []byte {
	// DNS header (12 bytes)
	header := make([]byte, 12)

	// ID (2 bytes)
	copy(header[0:2], id)

	// Question count (2 bytes)
	binary.BigEndian.PutUint16(header[4:6], 1) // One question

	// Answer count (2 bytes)
	binary.BigEndian.PutUint16(header[6:8], 0) // No answers

	// Authority count (2 bytes)
	binary.BigEndian.PutUint16(header[8:10], 0) // No authority records

	// Additional count (2 bytes)
	binary.BigEndian.PutUint16(header[10:12], 0) // No additional records

	// Question section
	questionSection := make([]byte, 4+len(ns))

	// Name
	questionSection[0] = byte(len(ns))
	copy(questionSection[1:], []byte(ns))

	// Type (A record)
	binary.BigEndian.PutUint16(questionSection[1+len(ns):], 1)

	// Class (IN)
	binary.BigEndian.PutUint16(questionSection[3+len(ns):], 1)

	// Concatenate the header and question section
	response := append(header, questionSection...)

	return response
}
