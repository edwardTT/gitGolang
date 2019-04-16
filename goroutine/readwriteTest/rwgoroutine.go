package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type Reader interface {
	Read(rc chan []byte)
}

type Writer interface {
	Write(wc chan string)
}

type LogProcess struct {
	rc    chan []byte
	wc    chan string
	read  Reader
	write Writer
}

type ReadFromFile struct {
	path string
}

func (r *ReadFromFile) Read(rc chan []byte) {
	// 读取模块
	f, err := os.Open(r.path)
	if err != nil {
		panic(fmt.Sprintf("open file failed. error %s", err.Error()))
	}
	// 从文件末尾开始逐行读取文件内容
	f.Seek(0, 2)

	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadBytes('\n')
		if err == io.EOF {
			time.Sleep(500 * time.Millisecond)
			continue
		} else if err != nil {
			panic(fmt.Sprintf("ReadByte error: %v ", err.Error()))
		}
		rc <- line[:len(line)-1]

	}

}

type WriteToInfluxDB struct {
	influxBDDsn string
}

func (w *WriteToInfluxDB) Write(wc chan string) {
	//写入模块
	for v := range wc {
		fmt.Println(v)
	}

}

//Process 文件解析方法
func (l *LogProcess) Process() {
	//解析模块
	for v := range l.rc {
		l.wc <- strings.ToUpper(string(v))
	}

}

func main() {
	r := &ReadFromFile{
		path: "./test.log",
	}
	w := &WriteToInfluxDB{
		influxBDDsn: "edward.liu",
	}
	//使用go fmt 格式化后
	lp := &LogProcess{
		rc:    make(chan []byte),
		wc:    make(chan string),
		read:  r,
		write: w,
	}

	go lp.read.Read(lp.rc)
	go lp.Process()
	go lp.write.Write(lp.wc)

	time.Sleep(60 * time.Second)
}
