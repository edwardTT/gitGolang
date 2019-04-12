package handler

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// UploadHandler : 上传文件
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//返回上传html页面
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, " go to server fail")
			fmt.Println("Read index.html file")
			return
		}
		io.WriteString(w, string(data))
	} else if r.Method == "POST" {
		//接收文件及存储到本地的文件
	}

}
