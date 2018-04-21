package bench

import (
	"math/rand"
	"time"
)

func Rand(size int64, times int) (d time.Duration) {
	s := make([]byte, size)
	t := time.Now()
	for i := 0; i < times; i++ {
		rand.Read(s)
	}
	return time.Since(t)
}
