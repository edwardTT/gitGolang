package main

import (
	"fmt"
	"gitGolang/edwardfilesever/lesson2/handler"
	"net/http"
)

func main() {
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/suc", handler.UploadSucHandler)
	fmt.Println("HandleFunc")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Failed  start, err:%s", err.Error())
	}
	fmt.Println("listen")
}
