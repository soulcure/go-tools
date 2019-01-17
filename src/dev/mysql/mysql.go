package mysql

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"log"
)

type Person struct {
	Id       int64  `db:"id" redis:"id,omitempty"`
	UserId   string `db:"user_id" redis:"user_id"`
	UserName string `db:"user_name" redis:"user_name"`
	Password string `db:"password" redis:"password"`
	Gender   int    `db:"gender" redis:"email"`
	Email    string `db:"email" redis:"gender"`
}

var db *sqlx.DB

func init() {
	database, err := sqlx.Open("mysql", "huxin001:youmai@2018@tcp(120.24.37.50:9906)/charge?charset=utf8&loc=Local")
	if err != nil {
		logrus.Error(err)
		return
	}
	db = database
	log.Print("mysql init success")
}

func Insert(userId, userName, password, email string, gender int) (int64, error) {
	r, err := db.Exec("insert into person(user_id,user_name, password, email,gender)values(?,?, ?, ?,?)", userId, userName, password, email, gender)
	if err != nil {
		logrus.Error(err)
		return 0, err
	} else {
		logrus.Debug("Insert success:", userName)
		return r.LastInsertId()
	}

}

func Update(userId string, gender int, email string) error {
	_, err := db.Exec("update person set gender=?,email=? where user_id=?", gender, email, userId)
	if err != nil {
		logrus.Error(err)
	} else {
		logrus.Debug("update success:", userId)
	}

	return err
}

func Select(userName, password string) (Person, error) {
	var person Person
	err := db.QueryRow("select * from person where user_name = ? and password=?", userName, password).Scan(&person.Id, &person.UserId, &person.UserName, &person.Password, &person.Gender, &person.Email)
	if err != nil {
		logrus.Error(err)
	} else {
		logrus.Debug("select user success:", person)
	}

	return person, err
}
