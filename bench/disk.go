package bench

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

func Write(size int64) (times int) {
	os.Mkdir("temp", 600)
	c := make(chan int)
	go func() {
		ts := 0
		t := time.Now()
		for {
			name := fmt.Sprintf("temp/gobench-%d", ts)
			file, err := os.Create(name)
			if err != nil {
				panic(err)
			}

			_, err = io.Copy(file, &io.LimitedReader{
				R: &ZeroReadWriter{},
				N: size,
			})

			if err != nil {
				panic(err)
			}

			err = file.Close()
			if err != nil {
				panic(err)
			}

			ts++

			if time.Since(t) >= sleepTime {
				c <- ts
				return
			}
		}
	}()

	times = <-c
	return
}

func Read(size int64) (times int) {
	c := make(chan int)
	fs, err := ioutil.ReadDir("temp")
	if err != nil {
		panic(err)
	}

	s := len(fs)

	go func() {
		ts := 0
		t := time.Now()
		for {
			rd := rand.Intn(s)
			name := fmt.Sprintf("temp/gobench-%d", rd)
			file, err := os.Open(name)
			if err != nil {
				panic(err)
			}

			_, err = io.Copy(&ZeroReadWriter{}, &io.LimitedReader{
				R: file,
				N: size,
			})

			if err != nil {
				panic(err)
			}

			err = file.Close()
			if err != nil {
				panic(err)
			}

			ts++
			if time.Since(t) >= sleepTime {
				c <- ts
				return
			}
		}
	}()

	times = <-c
	for _, f := range fs {
		err = os.Remove(fmt.Sprintf("temp/%s", f.Name()))
		if err != nil {
			panic(err)
		}
	}
	err = os.Remove("temp")
	if err != nil {
		panic(err)
	}
	return
}

type ZeroReadWriter struct{}

func (z *ZeroReadWriter) Read(p []byte) (n int, err error) {
	return len(p), nil
}

func (z *ZeroReadWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
