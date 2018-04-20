package bench

import (
	"crypto/sha256"
	"io"
	"log"
	"math/rand"
	"sync"
	"time"
)

var sleepTime = 10 * time.Second

func Hash(thread int, size int64) (hashes int) {
	wg := new(sync.WaitGroup)
	wg.Add(thread)
	for i := 0; i < thread; i++ {
		go func() {
			hashes = hashes + s256(size)
			wg.Done()
		}()
	}
	wg.Wait()
	return
}

func s256(size int64) (hashes int) {
	s := sha256.New()
	count := 0

	go func() {
		for {
			r := &io.LimitedReader{
				R: rand.New(rand.NewSource(233)),
				N: size,
			}

			_, err := io.Copy(s, r)
			if err != nil {
				log.Fatalln(err.Error())
			}

			s.Sum(nil)
			s.Reset()
			count++
		}
	}()

	time.Sleep(sleepTime)
	return count
}
