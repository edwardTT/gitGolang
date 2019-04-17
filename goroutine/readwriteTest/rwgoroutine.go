package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Reader interface {
	Read(rc chan []byte)
}

type Writer interface {
	Write(wc chan *Message)
}

type LogProcess struct {
	rc    chan []byte
	wc    chan *Message
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

type Message struct {
	TimeLocal                    time.Time
	BytesSent                    int
	Path, Method, Scheme, Status string
	UpstreamTime, RequestTime    float64
}

func (w *WriteToInfluxDB) Write(wc chan *Message) {
	//写入模块
	for v := range wc {
		fmt.Println(v)
	}

}

//Process 文件解析方法
func (l *LogProcess) Process() {
	// 解析模块

	/**
	172.0.0.12 - - [04/Mar/2018:13:49:52 +0000] http "GET /foo?query=t HTTP/1.0" 200 2133 "-" "KeepAliveClient" "-" 1.005 1.854
	**/

	r := regexp.MustCompile(`([\d\.]+)\s+([^ \[]+)\s+([^ \[]+)\s+\[([^\]]+)\]\s+([a-z]+)\s+\"([^"]+)\"\s+(\d{3})\s+(\d+)\s+\"([^"]+)\"\s+\"(.*?)\"\s+\"([\d\.-]+)\"\s+([\d\.-]+)\s+([\d\.-]+)`)

	loc, _ := time.LoadLocation("Asia/Shanghai")
	for v := range l.rc {
		ret := r.FindStringSubmatch(string(v))
		if len(ret) != 14 {
			log.Println("FindStringSubmatch fail:", string(v))
			continue
		}

		message := &Message{}
		t, err := time.ParseInLocation("02/Jan/2006:15:04:05 +0000", ret[4], loc)
		if err != nil {
			log.Println("ParseInLocation fail:", err.Error(), ret[4])
			continue
		}
		message.TimeLocal = t

		byteSent, _ := strconv.Atoi(ret[8])
		message.BytesSent = byteSent

		// GET /foo?query=t HTTP/1.0
		reqSli := strings.Split(ret[6], " ")
		if len(reqSli) != 3 {
			log.Println("strings.Split fail", ret[6])
			continue
		}
		message.Method = reqSli[0]

		u, err := url.Parse(reqSli[1])
		if err != nil {
			log.Println("url parse fail:", err)
			continue
		}
		message.Path = u.Path

		message.Scheme = ret[5]
		message.Status = ret[7]

		upstreamTime, _ := strconv.ParseFloat(ret[12], 64)
		requestTime, _ := strconv.ParseFloat(ret[13], 64)
		message.UpstreamTime = upstreamTime
		message.RequestTime = requestTime

		l.wc <- message
	}
}

func main() {
	r := &ReadFromFile{
		path: "./access.log",
	}
	w := &WriteToInfluxDB{
		influxBDDsn: "edward.liu",
	}
	//使用go fmt 格式化后
	lp := &LogProcess{
		rc:    make(chan []byte),
		wc:    make(chan *Message),
		read:  r,
		write: w,
	}

	go lp.read.Read(lp.rc)
	go lp.Process()
	go lp.write.Write(lp.wc)

	time.Sleep(60 * time.Second)
}
