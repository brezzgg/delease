//go:build !windows
package decoder

func Decode(data []byte) []byte {
	return data
}
