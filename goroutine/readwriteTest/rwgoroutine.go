package main

import (
	"fmt"
	"strings"
	"time"
)

type LogProcess struct {
	rc          chan string
	wc          chan string
	path        string
	influxBDDsn string
}

func (l *LogProcess) ReadFromFile() {
	// 读取模块
	line := "message"
	l.rc <- line
}

func (l *LogProcess) Process() {
	//解析模块
	data := <-l.rc
	l.wc <- strings.ToUpper(data)
}

func (l *LogProcess) WriteToInfluxBD() {
	//写入模块
	fmt.Println(<-l.wc)
}

func main() {
	lp := &LogProcess{
		rc:          make(chan string),
		wc:          make(chan string),
		path:        "/tmp/edward.log",
		influxBDDsn: "edward.liu",
	}

	go lp.ReadFromFile()
	go lp.Process()
	go lp.WriteToInfluxBD()

	time.Sleep(1 * time.Second)
}
