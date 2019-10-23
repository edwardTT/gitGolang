package main

import (
	"net"
	"time"
	"strconv"
	"gitGolang/protocolTest/protocolServer/protocol"
	"fmt"
	"os"
)

//发送100次请求
func send(conn net.Conn)  {
	for i := 0; i < 100; i++ {
		session := GetSession()
		words := "{\"ID\":\""+strconv.Itoa(i)+"\",\"Session\":\""+session+"20170914165908\",\"Meta\":\"golang\",\"Content\":\"message\"}"
		conn.Write(protocol.Enpack([]byte(words)))
		fmt.Println(words)      //打印发送出去的信息
	}
	fmt.Println("send over")
	defer conn.Close()
}
//用当前时间做识别。当前时间的十进制整数
func GetSession() string {
	gs1 := time.Now().Unix()
	gs2 := strconv.FormatInt(gs1, 10)
	return gs2
}

func main() {
	server := "localhost:7373"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil{
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil{
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	fmt.Println("connect success")

	send(conn)

}

// client 端串口输出信息
// https://www.jianshu.com/p/103c927928b1
//bwxzing@bwxzing-docker:~/go/src/gitGolang/protocolTest/client$ go run protocolClient.go
//connect success
//{"ID":"0","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"1","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"2","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"3","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"4","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"5","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"6","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"7","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"8","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"9","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"10","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"11","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"12","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"13","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"14","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"15","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"16","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"17","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"18","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"19","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"20","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"21","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"22","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"23","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"24","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"25","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"26","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"27","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"28","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"29","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"30","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"31","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"32","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"33","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"34","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"35","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"36","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"37","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"38","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"39","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"40","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"41","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"42","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"43","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"44","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"45","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"46","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"47","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"48","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"49","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"50","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"51","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"52","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"53","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"54","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"55","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"56","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"57","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"58","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"59","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"60","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"61","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"62","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"63","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"64","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"65","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"66","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"67","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"68","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"69","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"70","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"71","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"72","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"73","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"74","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"75","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"76","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"77","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"78","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"79","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"80","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"81","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"82","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"83","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"84","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"85","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"86","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"87","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"88","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"89","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"90","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"91","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"92","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"93","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"94","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"95","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"96","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"97","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"98","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//{"ID":"99","Session":"157181796820170914165908","Meta":"golang","Content":"message"}
//send over
