package main

import (
	"encoding/json"

	"runtime"

	"github.com/juju/loggo"
	"github.com/lzjluzijie/gobench/bench"
)

var MB = int64(1024 * 1024)
var hashSize = 1 * MB

var logger = loggo.GetLogger("main")

func main() {
	logger.SetLogLevel(loggo.INFO)
	//System info
	info, err := bench.GetInfo()
	if err != nil {
		logger.Errorf(err.Error())
	}
	j, err := json.Marshal(info)
	if err != nil {
		logger.Errorf(err.Error())
	}
	logger.Infof(string(j))

	//memory test
	d := bench.Rand(1024, 1000)
	if err != nil {
		logger.Errorf(err.Error())
	}
	logger.Infof("Read rand 1KB: %.2fµs", float64(d)/1000/1000)
	d = bench.Rand(1024*1024, 100)
	if err != nil {
		logger.Errorf(err.Error())
	}
	logger.Infof("Read rand 1MB: %.2fµs", float64(d)/1000/100)

	//disk test
	times := bench.Write(1024)
	logger.Infof("Write 1KB: %dfiles in 10s", times)
	times = bench.Read(1024)
	logger.Infof("Read 1KB: %dfiles in 10s", times)
	times = bench.Write(1 * MB)
	logger.Infof("Write 1MB: %dfiles in 10s", times)
	times = bench.Read(1 * MB)
	logger.Infof("Read 1MB: %dfiles in 10s", times)
	times = bench.Write(10 * MB)
	logger.Infof("Write 10MB: %dfiles in 10s", times*10/10)
	times = bench.Read(10 * MB)
	logger.Infof("Read 10MB: %dfiles in 10s", times*10/10)

	//CPU test
	threads := runtime.NumCPU()
	hashes := bench.Hash(threads, hashSize)
	logger.Infof("CPU benchmark: %d hashes with %d threads in 10s", hashes, threads)

	//speed test
	sts := []*bench.SpeedTest{
		bench.NewSpeedTest("北京联通", "http://www2.unicomtest.com:8080/download?size=10485760", 10*MB),
		bench.NewSpeedTest("上海联通", "http://211.95.17.50:8080/download?size=10485760", 10*MB),
		bench.NewSpeedTest("北京电信", "http://st1.bjtelecom.net:8080/download?size=10485760", 10*MB),
		bench.NewSpeedTest("广州电信", "http://gzspeedtest.com:8080/download?size=10485760", 10*MB),
		bench.NewSpeedTest("深圳移动", "http://speedtest3.gd.chinamobile.com:8080/download?size=10485760", 10*MB),
		bench.NewSpeedTest("北京移动", "http://speedtest.bmcc.com.cn:8080/download?size=10485760", 10*MB),
		bench.NewSpeedTest("东京Linode", "http://speedtest.tokyo.linode.com/100MB-tokyo.bin", 100*MB),
	}

	for _, st := range sts {
		err = st.Do()
		if err != nil {
			logger.Errorf("%s speed test error: %s", st.Name, err.Error())
			continue
		}

		logger.Infof(st.Result())
	}
}
