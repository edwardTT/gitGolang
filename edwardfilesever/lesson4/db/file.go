package db

import (
	"fmt"
	mydb "gitGolang/edwardfilesever/lesson4/db/mysql"
)

//OnFileUploadFinished :文件上传完成
func OnFileUploadFinished(filehash string, filename string,
	filesize int64, fileaddr string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into tbl_file (`file_sha1`,`file_name`,`file_size`," +
			"`file_addr`, `status`) values(?,?,?,?,1)")
	if err != nil {
		fmt.Println("Failed to prepare statment, err:" + err.Error())
		return false
	}
	ret, err := stmt.Exec(filehash, filename, filesize, fileaddr)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 {
			fmt.Printf("File with hash1:%s has been uploaded before", filehash)
		}
		return true
	}
	return false
}
