package handler

import (
	"fmt"
	dblayer "gitGolang/edwardfilesever/lesson4/db"
	"gitGolang/edwardfilesever/lesson4/util"
	"io/ioutil"
	"net/http"
)

const (
	pwd_salt = "#890"
)

//SignupHandler : 处理用户注册请求
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("SignupHandler r.Method %v \n", r.Method)
	if r.Method == http.MethodGet {
		fmt.Println("SignupHandler GET")
		data, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	} else if r.Method == http.MethodPost {
		fmt.Println("Reday to read the username and password")
		r.ParseForm()
		username := r.Form.Get("username")
		passwd := r.Form.Get("password")
		if len(username) < 3 || len(passwd) < 5 {
			w.Write([]byte("Invalid parameter"))
			return
		}

		encPasswd := util.Sha1([]byte(passwd + pwd_salt))
		fmt.Println("Reday to run UserSignUp")
		suc := dblayer.UserSignUp(username, encPasswd)
		fmt.Println("UserSignUp done")
		if suc {
			w.Write([]byte("SUCCESS"))
		} else {
			w.Write([]byte("Failed"))
		}
	}

}
