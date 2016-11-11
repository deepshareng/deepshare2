package pkg

import "strings"

const (
	//TODO: how to define this separator
	escapeSeparator = "|"
)

func EncodeStringSlice(strs []string) string {
	return strings.Join(strs, escapeSeparator)
}

func DecodeStringSlice(str string) []string {
	return strings.Split(str, escapeSeparator)
}
