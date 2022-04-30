package utils

import (
	"fmt"
	"os"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
)

const epoch = int64(1640966400000)

var (
	lock     sync.Mutex
	sequence int64 = 0
	lastTime int64 = 0
)

func GenMsgID() int64 {
	lock.Lock()
	defer lock.Unlock()

	now := time.Now().UTC().UnixMilli() - epoch
	if now == lastTime {
		sequence = (sequence + 1) & (int64(-1) >> 42)
		if sequence == 0 {
			for now <= lastTime {
				now = time.Now().UTC().UnixMilli() - epoch
			}
		}
	} else {
		sequence = 0
	}

	lastTime = now
	id := now<<23>>1 + sequence

	return id
}

func GenUUID() string {
	return uuid.NewV4().String()
}

func Bytes2BinStr(b []byte) string {
	var ret string
	for _, n := range b {
		ret += fmt.Sprintf("%08b ", n)
	}
	return ret
}

func FileExist(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		} else if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return !info.IsDir()
}
