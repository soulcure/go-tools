package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"log"
)

type Account struct {
	AccountId int    `db:"account_id" redis:"account_id,omitempty"`
	Uuid      string `db:"uuid" redis:"uuid"`
	UserName  string `db:"username" redis:"username"`
	Mobile    string `db:"mobile" redis:"mobile"`
	Email     string `db:"email" redis:"email"`
	Iso       string `db:"iso" redis:"iso"`
	Password  string `db:"password" redis:"password"`
}

const dbName = "root:123456@tcp(localhost:3306)/nuu_db?charset=utf8&loc=Local"

var db *sql.DB

func init() {
	database, err := sql.Open("mysql", dbName)
	if err != nil {
		logrus.Error(err)
		return
	}
	db = database
	log.Print("mysql init success")
}

//插入新注册用户数据
func RegisterInsert(uuid, username, email, mobile, iso, password string) (int64, error) {
	r, err := db.Exec("insert into account(uuid,username,email,mobile,iso,password)values(?,?,?,?,?,?)",
		uuid, username, email, mobile, iso, password)
	if err != nil {
		logrus.Error("mysql register insert " + err.Error())
		return 0, err
	} else {
		logrus.Debug("mysql register insert success:", username)
		return r.LastInsertId()
	}

}

func Update(userId string, gender int, email string) error {
	_, err := db.Exec("update person set gender=?,email=? where user_id=?", gender, email, userId)
	if err != nil {
		logrus.Error("mysql update " + err.Error())
	} else {
		logrus.Debug("update success:", userId)
	}

	return err
}

func AccountLogin(userName, email, password string) (Account, error) {
	var account Account

	var row *sql.Row
	if userName != "" {
		row = db.QueryRow("select * from account where username = ? and password=?", userName, password)
	} else {
		row = db.QueryRow("select * from account where email = ? and password=?", email, password)
	}

	err := row.Scan(&account.AccountId, &account.Uuid, &account.UserName,
		&account.Mobile, &account.Email, &account.Iso, &account.Password)
	if err != nil {
		logrus.Error("mysql select " + err.Error())
	} else {
		logrus.Debug("select user success:", account)
	}

	return account, err
}
