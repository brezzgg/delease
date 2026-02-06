//go:build windows
package decoder

import (
	"golang.org/x/sys/windows"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func Decode(data []byte) []byte {
	var chmap *charmap.Charmap
	cp, _ := windows.GetConsoleOutputCP()
	switch cp {
	case 65001:
		return data
	case 866:
		chmap = charmap.CodePage866
	case 1251:
		chmap = charmap.Windows1251
	case 850:
		chmap = charmap.CodePage850
	default:
		chmap = charmap.Windows1251
	}
	decoder := chmap.NewDecoder()
	decoded, _, err := transform.Bytes(decoder, data)
	if err != nil {
		return data
	}
	return decoded
}
