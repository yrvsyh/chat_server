package utils

import "fmt"

func Bytes2BinStr(b []byte) string {
	var ret string
	for _, n := range b {
		ret += fmt.Sprintf("%08b ", n)
	}
	return ret
}
