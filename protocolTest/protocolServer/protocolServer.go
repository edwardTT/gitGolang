package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"gitGolang/protocolTest/protocolServer/protocol"
)

func main() {
	netListen, err := net.Listen("tcp", "localhost:7373")
	CheckErr(err)
	defer netListen.Close()

	Log("Waiting for client ...")       //启动后，等待客户端访问。
	for{
		conn, err := netListen.Accept()     //监听客户端
		if err != nil {
			Log(conn.RemoteAddr().String(), "发了了错误：", err)
			continue
		}
		Log(conn.RemoteAddr().String(), "tcp connection success")
		go handleConnection(conn)
	}
}

//连接处理
func handleConnection(conn net.Conn) {
	//缓冲区，存储被截断的数据
	tmpBuffer := make([]byte, 0)
	//接收解包
	readerChannel := make(chan []byte, 10000)
	go reader(readerChannel)

	buffer := make([]byte, 1024)
	for{
		n, err := conn.Read(buffer)
		if err != nil{
			Log(conn.RemoteAddr().String(), "connection error: ", err)
			return
		}
		tmpBuffer = protocol.Depack(append(tmpBuffer, buffer[:n]...))
		readerChannel <- tmpBuffer      //接收的信息写入通道

	}
	defer conn.Close()
}

//获取通道数据
func reader(readerchannel chan []byte) {
	for{
		select {
		case data := <-readerchannel:
			Log(string(data))       //打印通道内的信息
		}
	}
}

//日志处理
func Log(v ...interface{}) {
	log.Println(v...)
}

//错误处理
func CheckErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

// server 端串口信息
// https://www.jianshu.com/p/103c927928b1
//bwxzing@bwxzing-docker:~/go/src/gitGolang/protocolTest/protocolServer$ go run protocolServer.go
//2019/10/23 16:05:59 Waiting for client ...
//2019/10/23 16:06:08 127.0.0.1:52512 tcp connection success
//2019/10/23 16:06:08 {"ID":"1","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"2","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"3","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"4","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"5","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"6","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"7","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"14","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"15","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"16","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"17","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"18","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"19","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"20","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"21","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"22","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"23","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"24","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"25","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"26","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"27","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"28","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"29","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"30","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"31","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"33","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"34","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"35","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"36","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"38","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"39","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"41","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"42","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"43","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"44","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"45","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"46","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"47","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"49","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"51","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"53","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"54","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"55","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"56","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"57","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"58","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"60","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"61","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"62","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"63","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"64","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"65","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"66","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"67","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"68","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"71","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"72","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"73","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"74","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"75","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"76","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"78","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"79","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"80","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"81","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"82","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"83","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"85","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"86","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"87","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"88","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"89","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"90","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"91","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"92","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"93","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"94","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"95","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"96","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"97","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"98","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 {"ID":"99","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//2019/10/23 16:06:08 127.0.0.1:52512 connection error:  EOF

