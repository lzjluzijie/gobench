package bench

import (
	"math/rand"
	"time"
)

func Rand(size int64, times int) (d time.Duration, err error) {
	s := make([]byte, size)
	t := time.Now()
	for i := 0; i < times; i++ {
		_, err = rand.Read(s)
	}
	return time.Since(t), err
}
