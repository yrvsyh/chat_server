package chat

import (
	"sync"
	"time"
)

var seq uint8 = 0
var lock sync.Mutex

func GenMsgID() int64 {
	lock.Lock()
	defer lock.Unlock()

	now := time.Now().UnixMicro()
	id := now<<8 + int64(seq)
	return id
}
