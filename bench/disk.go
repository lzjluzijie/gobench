package bench

import (
	"io"
	"os"
	"time"
)

func Write(size int64) (times int) {
	go func() {
		for {
			file, err := os.Create("gobench.tmp")
			if err != nil {
				panic(err)
			}

			r := &io.LimitedReader{
				R: new(ZeroReadWriter),
				N: size,
			}

			_, err = io.Copy(file, r)
			if err != nil {
				panic(err)
			}
			file.Close()
			times++
		}
	}()

	time.Sleep(sleepTime)
	return
}

func Read(size int64) (times int) {
	go func() {
		for {
			file, err := os.Open("gobench.tmp")
			if err != nil {
				panic(err)
			}

			r := &io.LimitedReader{
				R: file,
				N: size,
			}

			io.Copy(new(ZeroReadWriter), r)
			file.Close()
			times++
		}
	}()

	time.Sleep(sleepTime)
	return
}

type ZeroReadWriter struct{}

func (z *ZeroReadWriter) Read(p []byte) (n int, err error) {
	return len(p), nil
}

func (z *ZeroReadWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
