package main

import (
	"fmt"
	"strings"
	"time"
)

type Reader interface {
	Read(rc chan string)
}

type Writer interface {
	Write(wc chan string)
}
type LogProcess struct {
	rc    chan string
	wc    chan string
	read  Reader
	write Writer
}
type ReadFromFile struct {
	path string
}

func (r *ReadFromFile) Read(rc chan string) {
	// 读取模块
	line := "message"
	rc <- line
}

type WriteToInfluxDB struct {
	influxBDDsn string
}

func (w *WriteToInfluxDB) Write(wc chan string) {
	//写入模块
	fmt.Println(<-wc)

}
func (l *LogProcess) Process() {
	//解析模块
	data := <-l.rc
	l.wc <- strings.ToUpper(data)
}

func main() {
	r := &ReadFromFile{
		path: "/tmp/edward.log",
	}
	w := &WriteToInfluxDB{
		influxBDDsn: "edward.liu",
	}
	//使用go fmt 格式化后
	lp := &LogProcess{
		rc:    make(chan string),
		wc:    make(chan string),
		read:  r,
		write: w,
	}

	go lp.read.Read(lp.rc)
	go lp.Process()
	go lp.write.Write(lp.wc)

	time.Sleep(1 * time.Second)
}
