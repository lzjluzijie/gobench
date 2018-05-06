package bench

import (
	"bytes"
	"io/ioutil"
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
		cmd := exec.Command("tracepath", "-b", tr.Host)
		cmd.Stdout = buf
		err = cmd.Run()
	}

	if err != nil {
		return
	}

	result, err := ioutil.ReadAll(buf)
	if err != nil {
		return
	}

	tr.Result = string(result)
	return
}
