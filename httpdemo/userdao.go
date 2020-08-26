package httpdemo

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
)

const (
	insertUserinfoSql        = "INSERT INTO `userinfo`(username, password, departname, created) VALUES (?, ?, ?, ?)"
	insertUserdetailSql      = "INSERT INTO `userdetail`(uid, intro, profile) VALUES (?, ?, ?)"
	getUserinfoByUsernameSql = `SELECT uid, username, password, departname, created FROM userinfo WHERE username = ?`
	getUserinfoByUidSql      = `SELECT uid, username, password, departname, created FROM userinfo WHERE uid = ?`
	getUserdetailByUidSql    = `SELECT uid, intro, profile FROM userdetail WHERE uid = ?`
)

type UserDao struct {
	db *sqlx.DB
}

func NewUserDao() (*UserDao, error) {
	userdao := UserDao{}

	db, err := sqlx.Open("mysql", UserDataSource)
	if err != nil {
		log.Println("failed to sqlx.Open(), err: ", err)
		return nil, err
	}

	userdao.db = db

	return &userdao, nil
}

func (ud *UserDao) CreateUserinfo(ctx context.Context, userinfo *Userinfo) error {
	result, err := ud.db.ExecContext(ctx, insertUserinfoSql,
		userinfo.Username, userinfo.Password, userinfo.Departname, userinfo.Created)
	if err != nil {
		log.Println("failed to ud.db.ExecContext(insertUserinfoSql), err: ", err)
		return err
	}

	_, err = result.RowsAffected()
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		} else {
			log.Println("failed to result.RowsAffected(), err: ", err)
			return err
		}
	} else {
		lastInsertID, err := result.LastInsertId()
		if err != nil {
			log.Println("failed to result.LastInsertId(), err: ", err)
			return err
		} else {
			userinfo.Uid = int32(lastInsertID)
		}
		return nil
	}
}

func (ud *UserDao) CreateUserdetail(ctx context.Context, userdetail *Userdetail) error {
	result, err := ud.db.ExecContext(ctx, insertUserdetailSql,
		userdetail.Uid, userdetail.Intro, userdetail.Profile)
	if err != nil {
		log.Println("failed to ud.db.ExecContext(), err: ", err)
		return err
	}

	_, err = result.RowsAffected()
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		} else {
			log.Println("failed to result.RowsAffectId(), err: ", err)
			return err
		}
	} else {
		return nil
	}
}

func (ud *UserDao) GetUserinfoByUid(ctx context.Context, uid uint64) (*Userinfo, error) {
	userinfo := Userinfo{}
	err := ud.db.GetContext(ctx, &userinfo, getUserinfoByUidSql, uid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			log.Println("failed to get userinfo, err: ", err)
			return nil, err
		}
	} else {
		return &userinfo, nil
	}
}

func (ud *UserDao) GetUserinfoByUsername(ctx context.Context, username string) (*Userinfo, error) {
	userinfo := Userinfo{}
	err := ud.db.GetContext(ctx, &userinfo, getUserinfoByUsernameSql, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			log.Println("failed to get userinfo, err: ", err)
			return nil, err
		}
	} else {
		return &userinfo, nil
	}
}

func (ud *UserDao) GetUserdetailByUid(ctx context.Context, uid uint32) (*Userdetail, error) {
	userdetail := Userdetail{}
	err := ud.db.GetContext(ctx, &userdetail, getUserdetailByUidSql, uid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			log.Println("failed to get userdetail, err: ", err)
			return nil, err
		}
	} else {
		return &userdetail, nil
	}
}
