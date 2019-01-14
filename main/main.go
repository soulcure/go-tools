package main

import (
	"../models"
	"../mysql"
	"github.com/astaxie/beego/logs"
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

const (
	SecretKey = "welcome to go server"
)

func notFound(ctx iris.Context) {
	ctx.StatusCode(http.StatusNotFound)
}

//当出现错误的时候，再试一次
func internalServerError(ctx iris.Context) {
	ctx.StatusCode(http.StatusRequestTimeout)
}

func registerHandler(ctx iris.Context) {
	username := ctx.FormValue("username")
	password := ctx.FormValue("password")
	email := ctx.FormValue("email")
	genderStr := ctx.FormValue("gender")
	gender, err := strconv.Atoi(genderStr)
	if err != nil {
		gender = 0
		logrus.Error(err)
	}

	if username != "" && password != "" && email != "" && genderStr != "" {
		if ok := mysql.Insert(username, password, email, gender); ok {
			logrus.Debug("user register success")

			var res models.ProtocolRsp
			res.SetCode(models.OK)
			res.ResponseWriter(ctx)
		} else {
			var res models.ProtocolRsp
			res.SetCode(models.UserNameErr)
			res.ResponseWriter(ctx)
		}

	} else {
		var res models.ProtocolRsp
		res.SetCode(models.RegisterErr)
		res.ResponseWriter(ctx)
	}

}

func loginHandler(ctx iris.Context) {
	/*db := redis.New(service.Config{
		Network:     "tcp",
		Addr:        "127.0.0.1:6379",
		Password:    "",
		Database:    "",
		MaxIdle:     0,
		MaxActive:   10,
		IdleTimeout: service.DefaultRedisIdleTimeout,
		Prefix:      ""}) // optionally configure the bridge between your redis server

	// use go routines to query the database
	// db.Async(true)
	// close connection when control+C/cmd+C
	iris.RegisterOnInterrupt(func() {
		if err := db.Close(); err != nil {
			logrus.Error(err)

		}
	})

	sess := sessions.New(sessions.Config{Cookie: "sessionscookieid", Expires: 45 * time.Minute})

	sess.UseDatabase(db)*/

	username := ctx.FormValue("username")
	password := ctx.FormValue("password")

	if username != "" && password != "" {
		if _, err := mysql.Select(username, password); err == nil {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"name": username,
				"exp":  time.Now().Add(time.Hour * 72).Unix(),
			})

			if t, err := token.SignedString([]byte(SecretKey)); err == nil {
				logs.Debug(username, "set Token:", t)
				//session
				//s := sess.Start(ctx)
				//s.Set(username, t)

				var res models.ProtocolRsp
				res.SetCode(models.OK)
				res.Data = &models.LoginRsp{Token: t}
				res.ResponseWriter(ctx)
				return
			}

		}
	}

	logrus.Error("user login fail")

	var res models.ProtocolRsp
	res.SetCode(models.LoginErr)
	res.ResponseWriter(ctx)

}

func tokenHandler(ctx iris.Context) {
	tokenString := ctx.GetHeader("token")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err == nil {
		if token.Valid {
			logs.Debug("Token is valid")
			ctx.Next()
		} else {
			ctx.StatusCode(http.StatusUnauthorized)
			var res models.ProtocolRsp
			res.SetCode(models.TokenExp)
			res.ResponseWriter(ctx)
			logs.Error("Token is not valid")
		}
	} else {
		ctx.StatusCode(http.StatusUnauthorized)
		var res models.ProtocolRsp
		res.SetCode(models.NotLogin)
		res.ResponseWriter(ctx)

		logs.Error("Unauthorized access to this resource")
	}

}

func updateProfile(ctx iris.Context) {
	username := ctx.FormValue("username")
	email := ctx.FormValue("email")
	genderStr := ctx.FormValue("gender")
	gender, err := strconv.Atoi(genderStr)
	if err != nil {
		gender = 0
		logrus.Error(err)
	}

	if email != "" && genderStr != "" {
		if ok := mysql.Update(gender, email, username); ok {
			logrus.Debug("user register success")
			var res models.ProtocolRsp
			res.SetCode(models.OK)
			res.ResponseWriter(ctx)
		}

	} else {
		var res models.ProtocolRsp
		res.SetCode(models.ParamErr)
		res.ResponseWriter(ctx)
	}

}

func main() {
	// the rest of the code stays the same.
	app := iris.New()

	app.OnErrorCode(iris.StatusNotFound, notFound)
	app.OnErrorCode(iris.StatusInternalServerError, internalServerError)

	// register our routes.
	app.Post("/register", registerHandler)
	app.Post("/login", loginHandler)
	app.Post("/api/update", tokenHandler, updateProfile)

	if err := app.Run(iris.Addr(":8080")); err != nil {
		panic(err)
	}

}
