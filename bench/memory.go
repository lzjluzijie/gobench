package bench

import (
	"fmt"
	"math/rand"
	"time"
)

type MemoryBench struct {
	Name     string
	Size     int64
	Times    int
	Duration time.Duration
	Speed    float64

	finished bool
}

func NewMemoryBench(size int64, times int) (b *MemoryBench) {
	return &MemoryBench{
		Name:  fmt.Sprintf("Memory: random read %dKB", size/1024),
		Size:  size,
		Times: times,
	}
}

func (b *MemoryBench) Do() {
	size := b.Size
	times := b.Times
	s := make([]byte, size)
	t := time.Now()
	for i := 0; i < times; i++ {
		rand.Read(s)
	}

	b.Duration = time.Since(t)
	b.Speed = float64(int64(times)*size/1048576) / b.Duration.Seconds()
	b.finished = true
	return
}

func (b *MemoryBench) Result() (result string) {
	if !b.finished {
		b.Do()
	}
	return fmt.Sprintf("%s: %.2fMB/s", b.Name, b.Speed)
}
