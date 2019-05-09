package main

import (
	"fmt"
	"gitGolang/edwardfilesever/lensone1_handler_get/handler"
	"net/http"
)

func main() {
	http.HandleFunc("/file/upload", handler.UploadHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Failed  start, err:%s", err.Error())
	}
}
