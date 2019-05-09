package handler

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// UploadHandler : 处理上传文件
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//返回上传html页面
		fmt.Println(" Handler Method GET")
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, " go to server fail")
			fmt.Println("Read index.html file")
			return
		}
		io.WriteString(w, string(data))
	} else if r.Method == "POST" {
		//接收文件及存储到本地的文件
		file, head, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("Failed to get data. err. %s \n", err.Error())
			return
		}
		defer file.Close()

		newFile, err := os.Create("/tmp/" + head.Filename)
		if err != nil {
			fmt.Printf("Failed to creat files. err:%s\n", err.Error())
		}
		defer newFile.Close()

		_, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("Failed to save data to new file.err: %s\n", err.Error())
			return
		}
		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}

//UploadSucHandler  :上传成功
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Uplaod finished")

}
