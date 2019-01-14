package mysql

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type Person struct {
	UserId   int    `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
	Gender   int    `db:"gender"`
	Email    string `db:"email"`
}

var db *sqlx.DB

//mysqluser = "huxin001"
//mysqlpass = "youmai@2018"
//mysqlurls = "120.24.37.50:9906"
//mysqldb = "/charge?charset=utf8&loc=Local"

func init() {
	//huxin001:youmai@2018@tcp(120.24.37.50:9906)/charge?charset=utf8&loc=Local
	database, err := sqlx.Open("mysql", "huxin001:youmai@2018@tcp(120.24.37.50:9906)/charge?charset=utf8&loc=Local")
	if err != nil {
		logrus.Error(err)
		return
	}
	db = database
}

func Insert(username, password, email string, gender int) bool {
	r, err := db.Exec("insert into person(username, password, email,gender)values(?, ?, ?,?)", username, password, email, gender)
	if err != nil {
		logrus.Error(err)
		return false
	}
	id, err := r.LastInsertId()
	if err != nil {
		logrus.Error(err)
		return false
	}

	logrus.Debug("insert success:", id)
	return true
}

func Update(gender int, email, username string) bool {
	_, err := db.Exec("update person set gender=?,email=? where username=?", gender, email, username)
	if err != nil {
		logrus.Error(err)
		return false
	}

	logrus.Debug("update success:", username)
	return true
}

func Select(username, password string) (Person, error) {
	var person Person
	err := db.QueryRow("select * from person where username = ?,password=?", username, password).Scan(&person.UserId, &person.Username, &person.Password, &person.Gender, &person.Email)
	if err != nil {
		logrus.Error(err)
	}
	return person, err
}
