package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type DBConfig struct {
	Mysql DBInfo `yaml:"mysql"`
}

//数据库账号配置
type DBInfo struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
	Charset  string `yaml:"charset"`
}

type Account struct {
	AccountId int    `db:"account_id" redis:"account_id,omitempty"`
	Uuid      string `db:"uuid" redis:"uuid"`
	UserName  string `db:"username" redis:"username"`
	Mobile    string `db:"mobile" redis:"mobile"`
	Email     string `db:"email" redis:"email"`
	Iso       string `db:"iso" redis:"iso"`
	Password  string `db:"password" redis:"password"`
}

var (
	db       *sql.DB
	dbConfig DBConfig
)

func init() {
	path := "./conf/db.yml"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("db conf file does not exist")
	}

	data, _ := ioutil.ReadFile(path)
	if err := yaml.Unmarshal(data, &dbConfig); err != nil {
		log.Panic("db conf yaml Unmarshal error ")
	}

	dbName := getConnURL(&dbConfig.Mysql)

	database, err := sql.Open("mysql", dbName)
	if err != nil {
		log.Panic("mysql can not connect")
		return
	}
	db = database
	log.Print("mysql connect at ", dbName)
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

func getConnURL(info *DBInfo) (url string) {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s",
		info.User, info.Password, info.Host, info.Port, info.Database, info.Charset)
}
