package dns

import (
	"errors"
)

const (
	queryHeaderSize = 12
)

var errInvalidDNSResponse = errors.New("invalid DNS response")
var errInvalidDNSRequest = errors.New("invalid DNS request")

const (
	A     = 1
	NS    = 2
	CNAME = 5
	SOA   = 6
	PTR   = 12
	MX    = 15
	TXT   = 16
	AAAA  = 28
)

type dnsMessage struct {
	NsName      string
	NsType      int
}

func parseRequest(message []byte) (dnsMessage, error) {
	var parsedMessage dnsMessage
	name, err := parseDNSName(message, queryHeaderSize)
	if err != nil {
		return parsedMessage, err
	}
	t, err := getNsType(message)
	if err != nil {
		return parsedMessage, err
	}

	parsedMessage = dnsMessage{
		NsName: name,
		NsType: t,
	}

	return parsedMessage, nil
}

func getNsType(query []byte) (int, error) {
	if len(query) < 25 {
		return 0, errInvalidDNSRequest
	}

	return int(query[24]), nil
}

func parseDNSName(response []byte, startPos int) (string, error) {
	nameParts := make([]string, 0)

	pos := startPos

	// Read each label in the name
	for {
		// Get the label length
		labelLen := int(response[pos])

		if labelLen == 0 {
			// The label length is zero, indicating the end of the name
			break
		}

		if labelLen >= 0xC0 {
			// The label length has high-order bits set, indicating a pointer to a compressed name
			offset := int(response[pos+1]) + (labelLen&0x3F)<<8
			name, err := parseDNSName(response, offset)
			if err != nil {
				return "", err
			}

			nameParts = append(nameParts, name)

			// Move past the pointer
			pos += 2

			// The compressed name should be the last part, so we exit the loop here
			break
		}

		// Read the label bytes and append them to the name parts
		labelBytes := response[pos+1 : pos+1+labelLen]
		nameParts = append(nameParts, string(labelBytes))

		// Move to the next label
		pos += labelLen + 1
	}

	return joinDNSNameParts(nameParts), nil
}

func joinDNSNameParts(parts []string) string {
	joined := ""
	for _, part := range parts {
		if joined == "" {
			joined = part
		} else {
			joined += "." + part
		}
	}
	return joined
}
