package main

import (
	"log"

	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/lzjluzijie/gobench/bench"
)

var MB = int64(1024 * 1024)
var hashSize = 1 * MB
var downloadSize = 100 * MB

func main() {
	//CPU test
	threads := runtime.NumCPU()
	hashes := bench.Hash(threads, hashSize)
	log.Printf("CPU benchmark %d hashes with %d threads in 10s", hashes, threads)

	//disk test
	d, err := bench.Write(1 * MB)
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Printf("Write 1MB speed %.2fMB/s", 1/d.Seconds())
	d, err = bench.Read()
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Printf("Read 1MB speed %.2fMB/s", 1/d.Seconds())

	d, err = bench.Write(16 * MB)
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Printf("Write 16MB speed %.2fMB/s", 16/d.Seconds())
	d, err = bench.Read()
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Printf("Read 16MB speed %.2fMB/s", 16/d.Seconds())
	d, err = bench.Write(256 * MB)
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Printf("Write 256MB speed %.2fMB/s", 256/d.Seconds())
	d, err = bench.Read()
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Printf("Read 256MB speed %.2fMB/s", 256/d.Seconds())
	err = os.Remove("gobench.tmp")
	if err != nil {
		log.Fatalln(err.Error())
	}

	//speed test
	t := time.Now()
	resp, err := http.Get(fmt.Sprintf("http://www2.unicomtest.com:8080/download?downloadSize=%d", downloadSize))
	if err != nil {
		log.Fatalln(err.Error())
	}

	_, err = ioutil.ReadAll(&io.LimitedReader{
		R: resp.Body,
		N: downloadSize,
	})

	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Printf("北京联通 time:%.2fs, speed:%.2fMB/s", time.Since(t).Seconds(), float64(downloadSize/MB)/time.Since(t).Seconds())
}
