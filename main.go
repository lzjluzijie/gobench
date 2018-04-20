package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"encoding/json"

	"runtime"

	"github.com/juju/loggo"
	"github.com/lzjluzijie/gobench/bench"
)

var MB = int64(1024 * 1024)
var hashSize = 1 * MB
var downloadSize = 100 * MB

var logger = loggo.GetLogger("main")

func main() {
	logger.SetLogLevel(loggo.INFO)
	//System info
	info, err := bench.GetInfo()
	j, err := json.Marshal(info)
	if err != nil {
		logger.Errorf(err.Error())
	}
	logger.Infof(string(j))

	//memory test
	d, err := bench.Rand(1024, 1000)
	if err != nil {
		logger.Errorf(err.Error())
	}
	logger.Infof("Read rand 1KB: %.2fµs", float64(d)/1000/1000)
	d, err = bench.Rand(1024*1024, 100)
	if err != nil {
		logger.Errorf(err.Error())
	}
	logger.Infof("Read rand 1MB: %.2fµs", float64(d)/1000/100)

	//disk test
	times := bench.Write(1024)
	logger.Infof("Write 1KB: %dKB/s", times/10)
	times = bench.Read(1024)
	logger.Infof("Read 1KB: %dKB/s", times/10)
	times = bench.Write(1 * MB)
	logger.Infof("Write 1MB: %dMB/s", times/10)
	times = bench.Read(1 * MB)
	logger.Infof("Read 1MB: %dMB/s", times/10)
	times = bench.Write(10 * MB)
	logger.Infof("Write 10MB: %dMB/s", times*10/10)
	times = bench.Read(10 * MB)
	logger.Infof("Read 10MB: %dMB/s", times*10/10)
	time.Sleep(5 * time.Second)
	err = os.Remove("gobench.tmp")
	if err != nil {
		logger.Errorf(err.Error())
	}

	//CPU test
	threads := runtime.NumCPU()
	hashes := bench.Hash(threads, hashSize)
	logger.Infof("CPU benchmark: %d hashes with %d threads in 10s", hashes, threads)

	//speed test
	t := time.Now()
	resp, err := http.Get(fmt.Sprintf("http://www2.unicomtest.com:8080/download?downloadSize=%d", downloadSize))
	if err != nil {
		logger.Errorf(err.Error())
	}

	_, err = ioutil.ReadAll(&io.LimitedReader{
		R: resp.Body,
		N: downloadSize,
	})

	if err != nil {
		logger.Errorf(err.Error())
	}

	logger.Infof("北京联通: time %.2fs, speed %.2fMB/s", time.Since(t).Seconds(), float64(downloadSize/MB)/time.Since(t).Seconds())
}
