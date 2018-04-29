package bench

import (
	"io"
	"log"
	"math/rand"
	"time"

	"fmt"

	"golang.org/x/crypto/sha3"
)

type SHA3Bench struct {
	Name     string
	Thread   int
	Size     int64
	Duration time.Duration
	Hashes   int
	Speed    float64

	finished bool
}

func NewSHA3Bench(thread int, size int64, duration time.Duration) (b *SHA3Bench) {
	return &SHA3Bench{
		Name:     fmt.Sprintf("CPU: sha3-512 %dMB", size/1048576),
		Thread:   thread,
		Size:     size,
		Duration: duration,
	}
}

func (b *SHA3Bench) Do() (err error) {
	thread := b.Thread
	size := b.Size
	duration := b.Duration
	hashes := 0
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
				if time.Since(t) >= duration {
					c <- hs
					return
				}
			}
		}()
	}

	for i := 0; i < thread; i++ {
		hashes = hashes + <-c
	}

	b.Hashes = hashes
	b.Speed = float64(hashes) / duration.Seconds()
	b.finished = true
	return
}

func (b *SHA3Bench) Result() (result string) {
	if !b.finished {
		err := b.Do()
		if err != nil {
			return fmt.Sprintf("%s err: %s", b.Name, err.Error())
		}
	}
	return fmt.Sprintf("%s: %d hashes in %d second", b.Name, b.Hashes, int(b.Duration.Seconds()))
}
