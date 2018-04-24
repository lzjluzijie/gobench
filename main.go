package main

import (
	"encoding/json"
	"os"
	"runtime"

	"log"

	"github.com/lzjluzijie/gobench/bench"
	"github.com/urfave/cli"
)

var MB = int64(1024 * 1024)
var hashSize = 1 * MB

var sts = []*bench.SpeedTest{
	bench.NewSpeedTest("北京联通", "http://www2.unicomtest.com:8080/download?size=10485760", 10*MB),
	bench.NewSpeedTest("上海联通", "http://211.95.17.50:8080/download?size=10485760", 10*MB),
	bench.NewSpeedTest("北京电信", "http://st1.bjtelecom.net:8080/download?size=10485760", 10*MB),
	bench.NewSpeedTest("广州电信", "http://gzspeedtest.com:8080/download?size=10485760", 10*MB),
	bench.NewSpeedTest("深圳移动", "http://speedtest3.gd.chinamobile.com:8080/download?size=10485760", 10*MB),
	bench.NewSpeedTest("北京移动", "http://speedtest.bmcc.com.cn:8080/download?size=10485760", 10*MB),
	bench.NewSpeedTest("东京Linode", "http://speedtest.tokyo.linode.com/100MB-tokyo.bin", 100*MB),
}

var app = cli.NewApp()

func init() {
	app.Name = "GoBench"
	app.Author = "Halulu"
	app.Version = "0.1.0"
	app.Usage = "A simple benchmark tool"
	app.Action = func(c *cli.Context) (err error) {
		if err = info(c); err != nil {
			return
		}
		if err = cpu(c); err != nil {
			return
		}
		if err = memory(c); err != nil {
			return
		}
		if err = disk(c); err != nil {
			return
		}
		if err = speed(c); err != nil {
			return
		}
		return
	}
	app.Commands = []cli.Command{
		{
			Name:   "info",
			Usage:  "Print system info",
			Action: info,
		},
		{
			Name:   "cpu",
			Usage:  "Run cpu benchmark",
			Action: cpu,
		},
		{
			Name:   "memory",
			Usage:  "Run memory benchmark",
			Action: memory,
		},
		{
			Name:   "disk",
			Usage:  "Run disk benchmark",
			Action: disk,
		},
		{
			Name:   "speed",
			Usage:  "Run speed test",
			Action: speed,
		},
	}
}

func info(c *cli.Context) (err error) {
	//System info
	info, err := bench.GetInfo()
	if err != nil {
		log.Println(err.Error())
	}
	j, err := json.Marshal(info)
	if err != nil {
		return
	}
	log.Println(string(j))
	return
}

func cpu(c *cli.Context) (err error) {
	//CPU test
	threads := runtime.NumCPU()
	hashes := bench.Hash(threads, hashSize)
	log.Printf("CPU benchmark: %d hashes with %d threads in 10s", hashes, threads)
	return
}

func memory(c *cli.Context) (err error) {
	//memory test
	d := bench.Rand(1024, 1000)
	log.Printf("Read rand 1KB: %.2fµs", float64(d)/1000/1000)
	d = bench.Rand(1024*1024, 100)
	log.Printf("Read rand 1MB: %.2fµs", float64(d)/1000/100)
	return
}

func disk(c *cli.Context) (err error) {
	//disk test
	times := bench.Write(1024)
	log.Printf("Write 1KB: %dfiles in 10s", times)
	//times = bench.Read(1024)
	//logger.Infof("Read 1KB: %dfiles in 10s", times)
	times = bench.Write(1 * MB)
	log.Printf("Write 1MB: %dfiles in 10s", times)
	//times = bench.Read(1 * MB)
	//logger.Infof("Read 1MB: %dfiles in 10s", times)
	times = bench.Write(10 * MB)
	log.Printf("Write 10MB: %dfiles in 10s", times*10/10)
	//times = bench.Read(10 * MB)
	//logger.Infof("Read 10MB: %dfiles in 10s", times*10/10)
	return
}

func speed(c *cli.Context) (err error) {
	//speed test
	for _, st := range sts {
		err := st.Do()
		if err != nil {
			log.Printf("%s speed test error: %s", st.Name, err.Error())
			continue
		}

		log.Printf(st.Result())
	}
	return
}

func main() {
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf(err.Error())
	}
}
