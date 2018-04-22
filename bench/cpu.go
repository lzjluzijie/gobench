package bench

import (
	"io"
	"log"
	"math/rand"
	"time"

	"golang.org/x/crypto/sha3"
)

var sleepTime = 10 * time.Second

func Hash(thread int, size int64) (hashes int) {
	c := make(chan int)
	for i := 0; i < thread; i++ {
		go func() {
			hash := sha3.New512()
			hs := 0
			t := time.Now()
			for {
				_, err := io.Copy(hash, &io.LimitedReader{
					R: rand.New(rand.NewSource(233)),
					N: size,
				})

				if err != nil {
					log.Fatalln(err.Error())
				}

				hash.Sum(nil)
				hash.Reset()
				hs++
				if time.Since(t) >= sleepTime {
					c <- hs
					return
				}
			}
		}()
	}

	for i := 0; i < thread; i++ {
		hashes = hashes + <-c
	}
	return
}
