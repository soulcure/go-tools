package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
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

const dbName = "root:123456@tcp(localhost:3306)/nuuinfo?charset=utf8&loc=Local"

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

func Insert(userId, userName, password, email string, gender int) (int64, error) {
	r, err := db.Exec("insert into person(user_id,user_name, password, email,gender)values(?,?, ?, ?,?)",
		userId, userName, password, email, gender)
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
	row := db.QueryRow("select * from person where user_name = ? and password=?", userName, password)
	err := row.Scan(&person.Id, &person.UserId, &person.UserName, &person.Password, &person.Gender, &person.Email)
	if err != nil {
		logrus.Error(err)
	} else {
		logrus.Debug("select user success:", person)
	}

	return person, err
}
