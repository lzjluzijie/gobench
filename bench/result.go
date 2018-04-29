package bench

import "time"

type Result struct {
	Info
	SHA3Bench
	MemoryBenches []MemoryBench
	SpeedTests    []SpeedTest
	Time          time.Time
}
