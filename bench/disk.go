package bench

import (
	"io"
	"io/ioutil"
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
		N: int64(size),
	}

	t := time.Now()
	_, err = io.Copy(file, r)
	d = time.Since(t)

	return
}

func Read() (d time.Duration, err error) {
	file, err := os.Open("gobench.tmp")
	if err != nil {
		return
	}
	defer file.Close()

	t := time.Now()
	_, err = ioutil.ReadAll(file)
	d = time.Since(t)
	return
}
