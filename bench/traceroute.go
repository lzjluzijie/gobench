package bench

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"runtime"
	"time"
)

type TraceRoute struct {
	Name     string
	Host     string
	Result   string
	Duration time.Duration

	finished bool
}

func NewTraceRoute(name, host string) (tr *TraceRoute) {
	return &TraceRoute{
		Name: name,
		Host: host,
	}
}

func (tr *TraceRoute) Do() (err error) {
	buf := &bytes.Buffer{}

	if runtime.GOOS == "windows" {
		cmd := exec.Command("tracert", tr.Host)
		cmd.Stdout = buf
		err = cmd.Run()
	} else {
		cmd := exec.Command("tracepath", fmt.Sprintf("-b %s", tr.Host))
		cmd.Stdout = buf
		err = cmd.Run()
	}

	if err != nil {
		return
	}

	r, err := ioutil.ReadAll(buf)

	log.Println(string(r))
	return
}
