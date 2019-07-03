package db

import (
	"fmt"
	mydb "gitGolang/edwardfilesever/lesson4/db/mysql"
)

//UserSignUp : 通过用户名及密码完成user表的注册操作
func UserSignUp(username string, passwd string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into tbl_user (`user_name`,`user_pwd`) values (?,?)")
	if err != nil {
		fmt.Println("Failed to inser, err" + err.Error())
		return false
	}
	fmt.Println("UserSignUp insert done")
	defer stmt.Close()

	ret, err := stmt.Exec(username, passwd)
	if err != nil {
		fmt.Println("Failed  to insert. Err " + err.Error())
		return false
	}
	if rowsAffected, err := ret.RowsAffected(); nil == err && rowsAffected > 0 {
		return true
	}
	return false
}
