package handler

import (
	"bytes"
	"encoding/json"
	"filestore-server/mq"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	cmn "filestore-server/common"
	cfg "filestore-server/config"
	dblayer "filestore-server/db"
	"filestore-server/meta"
	"filestore-server/store/oss"
	"filestore-server/util"
)

func init() {
	// 目录已存在
	if _, err := os.Stat(cfg.TempLocalRootDir); err == nil {
		return
	}

	// 尝试创建目录
	err := os.MkdirAll(cfg.TempLocalRootDir, 0744)
	if err != nil {
		log.Println("无法创建临时存储目录，程序将退出")
		os.Exit(1)
	}
}

// UploadHandler ： 处理文件上传
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// 返回上传html页面
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "internel server error")
			return
		}
		io.WriteString(w, string(data))
		// 另一种返回方式:
		// 动态文件使用http.HandleFunc设置，静态文件使用到http.FileServer设置(见main.go)
		// 所以直接redirect到http.FileServer所配置的url
		// http.Redirect(w, r, "/static/view/index.html",  http.StatusFound)
	} else if r.Method == "POST" {
		// 1. 从form表单中获得文件内容句柄
		file, head, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("Failed to get form data, err:%s\n", err.Error())
			return
		}
		defer file.Close()

		// 2. 把文件内容转为[]byte
		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, file); err != nil {
			fmt.Printf("Failed to get file data, err:%s\n", err.Error())
			return
		}

		// 3. 构建文件元信息
		fileMeta := meta.FileMeta{
			FileName: head.Filename,
			FileSha1: util.Sha1(buf.Bytes()), //　计算文件sha1
			FileSize: int64(len(buf.Bytes())),
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		// 4. 将文件写入临时存储位置
		fileMeta.Location = cfg.TempLocalRootDir + fileMeta.FileSha1 // 临时存储地址
		newFile, err := os.Create(fileMeta.Location)
		if err != nil {
			fmt.Printf("Failed to create file, err:%s\n", err.Error())
			return
		}
		defer newFile.Close()

		nByte, err := newFile.Write(buf.Bytes())
		if int64(nByte) != fileMeta.FileSize || err != nil {
			fmt.Printf("Failed to save data into file, writtenSize:%d, err:%s\n", nByte, err.Error())
			return
		}

		// 5. 同步或异步将文件转移到Ceph/OSS
		newFile.Seek(0, 0) // 游标重新回到文件头部
		// if cfg.CurrentStoreType == cmn.StoreCeph {
		// 	// 文件写入Ceph存储
		// 	data, _ := ioutil.ReadAll(newFile)
		// 	cephPath := "/ceph/" + fileMeta.FileSha1
		// 	_ = ceph.PutObject("userfile", cephPath, data)
		// 	fileMeta.Location = cephPath
		// } else if cfg.CurrentStoreType == cmn.StoreOSS {
		if cfg.CurrentStoreType == cmn.StoreOSS {
			// 文件写入OSS存储
			ossPath := "oss/" + fileMeta.FileSha1
			// 判断写入OSS为同步还是异步
			if !cfg.AsyncTransferEnable {
				err = oss.Bucket().PutObject(ossPath, newFile)
				if err != nil {
					fmt.Println(err.Error())
					w.Write([]byte("Upload failed!"))
					return
				}
				fileMeta.Location = ossPath
			} else {
				// 写入异步转移任务队列
				data := mq.TransferData{
					FileHash:      fileMeta.FileSha1,
					CurLocation:   fileMeta.Location,
					DestLocation:  ossPath,
					DestStoreType: cmn.StoreOSS,
				}
				pubData, _ := json.Marshal(data)
				pubSuc := mq.Publish(
					cfg.TransExchangeName,
					cfg.TransOSSRoutingKey,
					pubData,
				)
				if !pubSuc {
					// TODO: 当前发送转移信息失败，稍后重试
				}
			}
		}

		//6.  更新用户文件表记录
		// //meta.UpdateFileMeta(fileMeta)
		_ = meta.UpdateFileMetaDB(fileMeta)

		r.ParseForm()
		username := r.Form.Get("username")
		suc := dblayer.OnUserFileUploadFinished(username, fileMeta.FileSha1,
			fileMeta.FileName, fileMeta.FileSize)
		if suc {
			http.Redirect(w, r, "/static/view/home.html", http.StatusFound)
		} else {
			w.Write([]byte("Upload Failed."))
		}
	}
}

// UploadSucHandler : 上传已完成
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload finished!")
}

// GetFileMetaHandler : 获取文件元信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	filehash := r.Form["filehash"][0]
	//fMeta := meta.GetFileMeta(filehash)
	fMeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if fMeta != nil {
		data, err := json.Marshal(fMeta)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
	} else {
		w.Write([]byte(`{"code":-1,"msg":"no such file"}`))
	}
}

// FileQueryHandler : 查询批量的文件元信息
func FileQueryHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	limitCnt, _ := strconv.Atoi(r.Form.Get("limit"))
	username := r.Form.Get("username")
	//fileMetas, _ := meta.GetLastFileMetasDB(limitCnt)
	userFiles, err := dblayer.QueryUserFileMetas(username, limitCnt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(userFiles)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// DownloadHandler : 文件下载接口
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fsha1 := r.Form.Get("filehash")
	username := r.Form.Get("username")

	fm, _ := meta.GetFileMetaDB(fsha1)
	userFile, err := dblayer.QueryUserFileMeta(username, fsha1)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	f, err := os.Open(fm.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octect-stream")
	// attachment表示文件将会提示下载到本地，而不是直接在浏览器中打开
	w.Header().Set("content-disposition", "attachment; filename=\""+userFile.FileName+"\"")
	w.Write(data)
}

// FileMetaUpdateHandler ： 更新元信息接口(重命名)
func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	opType := r.Form.Get("op")
	fileSha1 := r.Form.Get("filehash")
	username := r.Form.Get("username")
	newFileName := r.Form.Get("filename")

	if opType != "0" || len(newFileName) < 1 {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// 更新用户文件表tbl_user_file中的文件名，tbl_file的文件名不用修改
	_ = dblayer.RenameFileName(username, fileSha1, newFileName)

	// 返回最新的文件信息
	userFile, err := dblayer.QueryUserFileMeta(username, fileSha1)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(userFile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// FileDeleteHandler : 删除文件及元信息
func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	fileSha1 := r.Form.Get("filehash")

	fm, err := meta.GetFileMetaDB(fileSha1)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 删除本地文件
	os.Remove(fm.Location)
	// TODO: 可考虑删除Ceph/OSS上的文件
	// 可以不立即删除，加个超时机制，
	// 比如该文件10天后也没有用户再次上传，那么就可以真正的删除了

	// 删除文件表中的一条记录
	suc := dblayer.DeleteUserFile(username, fileSha1)
	if !suc {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// TryFastUploadHandler : 尝试秒传接口
func TryFastUploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// 1. 解析请求参数
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filename := r.Form.Get("filename")
	filesize, _ := strconv.Atoi(r.Form.Get("filesize"))

	// 2. 从文件表中查询相同hash的文件记录
	fileMeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 3. 查不到记录则返回秒传失败
	if fileMeta == nil {
		resp := util.RespMsg{
			Code: -1,
			Msg:  "秒传失败，请访问普通上传接口",
		}
		w.Write(resp.JSONBytes())
		return
	}

	// 4. 上传过则将文件信息写入用户文件表， 返回成功
	suc := dblayer.OnUserFileUploadFinished(
		username, filehash, filename, int64(filesize))
	if suc {
		resp := util.RespMsg{
			Code: 0,
			Msg:  "秒传成功",
		}
		w.Write(resp.JSONBytes())
		return
	}
	resp := util.RespMsg{
		Code: -2,
		Msg:  "秒传失败，请稍后重试",
	}
	w.Write(resp.JSONBytes())
	return
}

// DownloadURLHandler : 生成文件的下载地址
func DownloadURLHandler(w http.ResponseWriter, r *http.Request) {
	filehash := r.Form.Get("filehash")
	// 从文件表查找记录
	row, _ := dblayer.GetFileMeta(filehash)

	// TODO: 判断文件存在OSS，还是Ceph，还是在本地
	if strings.HasPrefix(row.FileAddr.String, cfg.TempLocalRootDir) {
		username := r.Form.Get("username")
		token := r.Form.Get("token")
		tmpUrl := fmt.Sprintf("http://%s/file/download?filehash=%s&username=%s&token=%s",
			r.Host, filehash, username, token)
		w.Write([]byte(tmpUrl))
	} else if strings.HasPrefix(row.FileAddr.String, "/ceph") {
		// TODO: ceph下载url
	} else if strings.HasPrefix(row.FileAddr.String, "oss/") {
		// oss下载url
		signedURL := oss.DownloadURL(row.FileAddr.String)
		w.Write([]byte(signedURL))
	}
}
