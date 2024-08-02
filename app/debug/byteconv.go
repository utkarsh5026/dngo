package debug

import (
	"fmt"
	"strings"
)

// BytesToHex converts a byte slice into its hexadecimal string representation.
func BytesToHex(data []byte) string {

	var builder strings.Builder
	for _, b := range data {
		builder.WriteString(ToHex(b))
		builder.WriteString(" ")
	}
	return builder.String()
}

func ToHex(b byte) string {
	return fmt.Sprintf("%02x", b)
}
