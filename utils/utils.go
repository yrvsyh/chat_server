package utils

import (
	"fmt"
	"os"
)

func Bytes2BinStr(b []byte) string {
	var ret string
	for _, n := range b {
		ret += fmt.Sprintf("%08b ", n)
	}
	return ret
}

func FileExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		} else if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return true
}
