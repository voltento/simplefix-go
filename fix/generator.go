package fix

import (
	"bytes"
	"fmt"
	"strconv"
)

// TimeLayout is a time layout.
const TimeLayout = "20060102-15:04:05.000"

const (
	// CountOfSOHSymbols is the number of SOH delimiters in a full message.
	// Deprecated: should not be used, count SOH symbols by yourself
	CountOfSOHSymbols = 3
	// CountOfSOHSymbolsWithoutBody is the number of SOH delimiters excluding the body.
	// Deprecated: should not be used, count SOH symbols by yourself
	CountOfSOHSymbolsWithoutBody = 2
)

// Delimiter is a delimiter.
var Delimiter = []byte{1}

// DelimiterChar is a delimiter char.
const DelimiterChar = 1

func joinBody(values ...[]byte) []byte {
	return bytes.Join(values, Delimiter)
}

// nolint
func makeGroup(entries []map[string][]byte, tags []string) []byte {
	var groupItems [][]byte
	for _, entry := range entries {
		for _, tag := range tags {
			groupItems = append(groupItems, bytes.Join([][]byte{[]byte(tag), entry[tag]}, []byte{61}))
		}
	}

	return bytes.Join(groupItems, Delimiter)
}

// CalcCheckSum ...
func CalcCheckSum(body []byte) []byte {
	var sum int
	for _, b := range body {
		sum += int(b)
	}
	sum += int(byte(1))

	return []byte(fmt.Sprintf("%03s", strconv.Itoa(sum%256)))
}

// CalcCheckSumOptimized ...
func CalcCheckSumOptimized(bytes []byte) []byte {
	var sum int
	for _, b := range bytes {
		sum += int(b)
	}
	sum += int(byte(1))
	n := sum % 256
	return []byte{byte('0' + n/100), byte('0' + (n/10)%10), byte('0' + n%10)}
}

// CalcCheckSumOptimizedFromBuffer ...
func CalcCheckSumOptimizedFromBuffer(buffer *bytes.Buffer) []byte {
	var sum int
	for _, b := range buffer.Bytes() {
		sum += int(b)
	}
	n := (sum + int(byte(1))) % 256
	return []byte{byte('0' + n/100), byte('0' + (n/10)%10), byte('0' + n%10)}
}
