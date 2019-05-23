package main

import (
	"fmt"
	"gitGolang/edwardfilesever/lesson4/handler"
	"net/http"
)

func main() {
	//静态资源处理
	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("./static"))))

	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/suc", handler.UploadSucHandler)
	http.HandleFunc("/file/meta", handler.GetFileMetaHandler)
	http.HandleFunc("/file/query", handler.FileQueryHandler)
	http.HandleFunc("/file/download", handler.DownloadHandler)

	http.HandleFunc("/user/signup", handler.SignupHandler)
	fmt.Println("Start HandleFunc")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Failed  start, err:%s", err.Error())
	}
	fmt.Println("listen done")
}
