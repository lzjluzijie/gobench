package bench

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type SpeedTest struct {
	Name     string
	URL      string
	Size     int64
	Duration time.Duration
	Speed    float64

	finished bool
}

func NewSpeedTest(name, url string, size int64) (st *SpeedTest) {
	return &SpeedTest{
		Name: fmt.Sprintf("SpeedTest: %s %dMB", name, size/1048576),
		URL:  url,
		Size: size,
	}
}

func (st *SpeedTest) Do() (err error) {
	t := time.Now()
	resp, err := http.Get(st.URL)
	if err != nil {
		return
	}

	_, err = io.Copy(&ZeroReadWriter{}, &io.LimitedReader{
		R: resp.Body,
		N: st.Size,
	})

	if err != nil {
		return
	}

	st.Duration = time.Since(t)
	st.Speed = float64(st.Size/1048576) / time.Since(t).Seconds()
	st.finished = true
	return
}

func (st *SpeedTest) Result() (result string) {
	if !st.finished {
		err := st.Do()
		if err != nil {
			return ""
		}
	}
	return fmt.Sprintf("%s: time %.2fs, speed %.2fMB/s", st.Name, st.Duration.Seconds(), st.Speed)
}
