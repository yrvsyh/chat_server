package utils_test

import (
	"chat_server/utils"
	"testing"
)

func BenchmarkGenMsgID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		utils.GenMsgID()
	}
}
