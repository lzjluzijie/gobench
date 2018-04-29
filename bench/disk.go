package bench

import (
	"fmt"
	"io"
	"os"
	"time"
)

type DiskBench struct {
	Name     string
	Size     int64
	Times    int
	Duration time.Duration
	Speed    float64

	finished bool
}

func NewDiskBench(size int64, duration time.Duration) (b *DiskBench) {
	return &DiskBench{
		Name:     fmt.Sprintf("Disk: write %dKB", size/1024),
		Size:     size,
		Duration: duration,
	}
	return
}

func (b *DiskBench) Do() (err error) {
	size := b.Size
	duration := b.Duration
	times := 0
	c := make(chan int)
	for i := 0; i < 4; i++ {
		id := i
		go func() {
			for {
				ts := 0
				t := time.Now()
				for {
					name := fmt.Sprintf("gobench-temp-%d", id)
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

					err = os.Remove(name)

					ts++

					if time.Since(t) >= duration {
						c <- ts
						return
					}
				}
			}
		}()
	}

	for i := 0; i < 4; i++ {
		times = times + <-c
	}

	b.Times = times
	b.Speed = float64(int64(times)*size/1048576) / b.Duration.Seconds()
	b.finished = true
	return
}

func (b *DiskBench) Result() (result string) {
	if !b.finished {
		err := b.Do()
		if err != nil {
			return err.Error()
		}
	}
	return fmt.Sprintf("%s: %.2fMB/s", b.Name, b.Speed)
}

type ZeroReadWriter struct{}

func (z *ZeroReadWriter) Read(p []byte) (n int, err error) {
	return len(p), nil
}

func (z *ZeroReadWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
