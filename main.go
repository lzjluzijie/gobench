package main

import (
	"encoding/json"
	"os"
	"runtime"

	"log"

	"time"

	"github.com/lzjluzijie/gobench/bench"
	"github.com/urfave/cli"
)

var sts = []*bench.SpeedTest{
	bench.NewSpeedTest("北京联通", "http://www2.unicomtest.com:8080/download?size=10485760", 10485760),
	bench.NewSpeedTest("上海联通", "http://211.95.17.50:8080/download?size=10485760", 10485760),
	bench.NewSpeedTest("北京电信", "http://st1.bjtelecom.net:8080/download?size=10485760", 10485760),
	bench.NewSpeedTest("广州电信", "http://gzspeedtest.com:8080/download?size=10485760", 10485760),
	bench.NewSpeedTest("深圳移动", "http://speedtest3.gd.chinamobile.com:8080/download?size=10485760", 10485760),
	bench.NewSpeedTest("北京移动", "http://speedtest.bmcc.com.cn:8080/download?size=10485760", 10485760),
	bench.NewSpeedTest("CacheFly", "http://cachefly.cachefly.net/100mb.test", 104857600),
	bench.NewSpeedTest("东京Linode", "http://speedtest.tokyo.linode.com/100MB-tokyo.bin", 104857600),
	bench.NewSpeedTest("洛杉矶Psychz", "http://lg.lax.psychz.net/200MB.test", 209715200),
}

var app = cli.NewApp()

func init() {
	app.Name = "GoBench"
	app.Version = "0.3.0"
	app.Usage = "A simple benchmark tool"
	app.Description = "See https://github.com/lzjluzijie/gobench"
	app.Author = "Halulu"
	app.Email = "lzjluzijie@gmail.com"
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
		if err = speedtest(c); err != nil {
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
			Name:   "speedtest",
			Usage:  "Run speedtest",
			Action: speedtest,
		},
		{
			Name:   "traceroute",
			Usage:  "Run traceroute",
			Action: traceroute,
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
	b := bench.NewSHA3Bench(threads, 1*1048576, 10*time.Second)
	log.Printf(b.Result())
	return
}

func memory(c *cli.Context) (err error) {
	//memory test
	b := bench.NewMemoryBench(1024, 10000)
	log.Printf(b.Result())
	b = bench.NewMemoryBench(1024*1024, 500)
	log.Printf(b.Result())
	return
}

func disk(c *cli.Context) (err error) {
	//disk test
	b := bench.NewDiskBench(1024, 10*time.Second)
	log.Println(b.Result())
	b = bench.NewDiskBench(1024*1024, 10*time.Second)
	log.Println(b.Result())
	return
}

func speedtest(c *cli.Context) (err error) {
	//speedtest
	for _, st := range sts {
		log.Printf(st.Result())
	}
	return
}

func traceroute(c *cli.Context) (err error) {
	tr := bench.NewTraceRoute("北京联通", "www2.unicomtest.com")
	err = tr.Do()
	return
}

func main() {
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf(err.Error())
	}
}
