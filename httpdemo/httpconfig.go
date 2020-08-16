package httpdemo

const (
	UserDataSource      = "root:1234@tcp(127.0.0.1:3306)@web?timeout=90s&collation=utf8mb4_bin&parseTime=True"
	Address             = "localhost:4000"
	TemplateDir         = "./templates"
	TestDir             = "/Users/yanghongjie/coding/go_work/go-web/test"
	UploadFile2Fullpath = "/Users/yanghongjie/coding/henry_data/upload4.txt"
	UploadFile2Name     = "upload4.txt"
	UploadUrl           = "http://localhost:4000/upload"

	InsertUserinfoSql   = "INSERT INTO `userinfo`(username, departname, created) VALUES (?, ?, ?)"
	InsertUserdetailSql = "INSERT INTO `userdetail`(uid, intro, profile) VALUES (?, ?, ?)"
)
