package httpdemo

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
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
	result, err := ud.db.ExecContext(ctx, InsertUserinfoSql,
		userinfo.Username, userinfo.Departname, userinfo.Created)
	if err != nil {
		log.Println("failed to ud.db.ExecContext(InsertUserinfoSql), err: ", err)
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
	result, err := ud.db.ExecContext(ctx, InsertUserdetailSql,
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
