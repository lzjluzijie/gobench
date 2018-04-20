package bench

import (
	"io"
	"math/rand"
	"os"
	"time"
)

func Write(size int64) (d time.Duration, err error) {
	file, err := os.Create("gobench.tmp")
	if err != nil {
		return
	}
	defer file.Close()

	r := &io.LimitedReader{
		R: rand.New(rand.NewSource(233)),
		N: size,
	}

	t := time.Now()
	_, err = io.Copy(file, r)
	d = time.Since(t)

	return
}

func Read(size int64) (d time.Duration, err error) {
	file, err := os.Open("gobench.tmp")
	if err != nil {
		return
	}
	defer file.Close()

	r := &io.LimitedReader{
		R: file,
		N: size,
	}
	t := time.Now()
	io.Copy(new(NWriter), r)
	d = time.Since(t)
	return
}

type NWriter struct{}

func (w *NWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
