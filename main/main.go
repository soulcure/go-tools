package main

import (
	"../models"
	"../mysql"
	"github.com/astaxie/beego/logs"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
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

func before(ctx iris.Context) {
	shareInformation := "this is a sharable information between handlers"

	requestPath := ctx.Path()
	logrus.Debug("Before the mainHandler: " + requestPath)
	ctx.Values().Set("info", shareInformation)
	ctx.Next() //继续执行下一个handler，在本例中是mainHandler。
}

func after(ctx iris.Context) {
	requestPath := ctx.Path()
	logrus.Debug("After the mainHandler: " + requestPath)
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
	token, err := request.ParseFromRequest(ctx.Request(), request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		})

	if err == nil {
		if token.Valid {
			ctx.Next()
		} else {
			ctx.StatusCode(http.StatusUnauthorized)
			var res models.ProtocolRsp
			res.SetCode(models.TokenExp)
			res.ResponseWriter(ctx)
			logs.Debug("Token is not valid")
		}
	} else {
		ctx.StatusCode(http.StatusUnauthorized)
		var res models.ProtocolRsp
		res.SetCode(models.NotLogin)
		res.ResponseWriter(ctx)

		logs.Debug("Unauthorized access to this resource")
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

	//将“before”处理程序注册为将要执行的第一个处理程序
	//在所有域的路由上。
	//或使用`UseGlobal`注册一个将跨子域触发的中间件。
	app.Use(before)

	//将“after”处理程序注册为将要执行的最后一个处理程序
	//在所有域的路由'处理程序之后。
	app.Done(after)

	// register our routes.
	app.Post("/register", registerHandler)
	app.Post("/login", loginHandler)
	app.Post("/api/update", tokenHandler, updateProfile, after)

	if err := app.Run(iris.Addr(":8080")); err != nil {
		panic(err)
	}

}
