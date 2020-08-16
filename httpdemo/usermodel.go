package httpdemo

type Userinfo struct {
	Uid        int32  `db:"uid"`
	Username   string `db:"username"`
	Departname string `db:"departname"`
	Created    int64  `db:"created"`
}

type Userdetail struct {
	Uid     int32  `db:"uid"`
	Intro   string `db:"intro"`
	Profile string `db:"profile"`
}
