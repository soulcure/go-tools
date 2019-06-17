package main

import (
	"dev/models"
	"dev/mysql"
	"dev/redis"
	"dev/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	SecretKey = "welcome to go server"
)

func init() {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
}

func notFound(ctx iris.Context) {
	ctx.StatusCode(http.StatusNotFound)
	var res models.ProtocolRsp
	res.Code = models.NoFoundErrCode
	res.Msg = models.NoFoundErr
	res.ResponseWriter(ctx)
}

//当出现错误的时候，再试一次
func internalServerError(ctx iris.Context) {
	ctx.StatusCode(http.StatusRequestTimeout)
	var res models.ProtocolRsp
	res.Code = models.UnknownErrCode
	res.Msg = models.UnknownErr
	res.ResponseWriter(ctx)
}

func test(ctx iris.Context) {
	ctx.Application().Logger().Info("Request path: %s", ctx.Path())
	ctx.Application().Logger().Infof("Request path: %+v", ctx)
	ctx.Application().Logger().Debug("Request path: %+v", ctx)
	var res models.ProtocolRsp
	res.Code = models.OK
	res.Msg = models.SUCCESS
	res.ResponseWriter(ctx)
}

//用户注册处理函数
func registerHandler(ctx iris.Context) {
	username := ctx.FormValue("username")
	email := ctx.FormValue("email")
	mobile := ctx.FormValue("mobile")
	iso := ctx.FormValue("iso")
	password := ctx.FormValue("password")

	ctx.Application().Logger().Info("Request path: %s", ctx.Path())
	ctx.Application().Logger().Infof("Request path: %+v", ctx)
	ctx.Application().Logger().Debug("Request path: %+v", ctx)

	if checkRegisterFormat(ctx, username, email, mobile, iso, password) {
		userId := uuid.Must(uuid.NewV4()).String()
		logrus.Debug("user register uuid:", userId)
		if _, err := mysql.RegisterInsert(userId, username, email, mobile, iso, password); err == nil {
			logrus.Debug("user register success")
			var res models.ProtocolRsp
			res.Code = models.OK
			res.Msg = models.SUCCESS
			res.Data = &models.RegisterRsp{Uuid: userId, UserName: username, Email: email, PassWord: password}
			res.ResponseWriter(ctx)

		} else {
			var res models.ProtocolRsp
			res.Code = models.RegisterErrCode
			res.Msg = err.Error()
			res.ResponseWriter(ctx)
		}

	}

}

func loginHandler(ctx iris.Context) {
	username := ctx.FormValue("username")
	email := ctx.FormValue("email")
	password := ctx.FormValue("password")

	if checkLoginFormat(ctx, username, email, password) {
		if account, err := mysql.AccountLogin(username, email, password); err == nil {
			var name string
			if username != "" {
				name = username
			} else {
				name = email
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"name": name,
				"exp":  time.Now().Add(time.Hour * 72).Unix(),
			})

			if token, err := token.SignedString([]byte(SecretKey)); err == nil {
				logrus.Debug(username, "  set Token:", token)

				if _, e := redis.SetStruct(account.Uuid, account); e == nil {
					var res models.ProtocolRsp
					res.Code = models.OK
					res.Msg = models.SUCCESS
					res.Data = &models.LoginRsp{Token: token, Uuid: account.Uuid}
					res.ResponseWriter(ctx)
					return
				}

			}

		}
	}
}

func tokenHandler(ctx iris.Context) {
	tokenString := ctx.GetHeader("token")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err == nil {
		if token.Valid {
			logrus.Debug("Token is valid")
			ctx.Next()
		} else {
			ctx.StatusCode(http.StatusUnauthorized)
			var res models.ProtocolRsp
			res.Code = models.TokenExpCode
			res.Msg = models.TokenExpiredErr
			res.ResponseWriter(ctx)
			logrus.Error("Token is not valid")
		}
	} else {
		ctx.StatusCode(http.StatusUnauthorized)
		var res models.ProtocolRsp
		res.Code = models.NotLoginCode
		res.Msg = models.TokenErr
		res.ResponseWriter(ctx)

		logrus.Error("Unauthorized access to this resource")
	}

}

func updateProfile(ctx iris.Context) {
	userId := ctx.FormValue("userId")
	email := ctx.FormValue("email")
	genderStr := ctx.FormValue("gender")
	gender, err := strconv.Atoi(genderStr)
	if err != nil {
		gender = 0
		logrus.Error(err)
	}

	if userId != "" && email != "" {
		if err := mysql.Update(userId, gender, email); err == nil {
			logrus.Debug("user update profile success")
			var e error

			user := &mysql.Account{}
			e = redis.GetStruct(userId, user)

			logrus.Debug("user profile:", user)

			user.Email = email
			_, e = redis.SetStruct(userId, user)

			if e == nil {
				var res models.ProtocolRsp
				res.Code = models.OK
				res.Msg = models.SUCCESS
				res.ResponseWriter(ctx)
				return
			}
		}

	}
}

// Get a filename based on the date, just for the sugar.
func todayFilename() string {
	today := time.Now().Format("2006-01-02")
	return today + ".txt"
}

func newLogFile() *os.File {
	filename := todayFilename()
	// Open the file, this will append to the today's file if server restarted.
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	return f
}

func main() {
	f := newLogFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("close log file error: %s", err)
		}
	}()

	// the rest of the code stays the same.
	app := iris.New()

	app.Logger().SetOutput(f)

	app.OnErrorCode(iris.StatusNotFound, notFound)
	app.OnErrorCode(iris.StatusInternalServerError, internalServerError)

	// register our routes.
	app.Get("/test", test)
	app.Post("/register", registerHandler)
	app.Post("/login", loginHandler)
	app.Post("/api/update", tokenHandler, updateProfile)

	app.Logger().SetLevel("debug")

	config := iris.WithConfiguration(iris.YAML("./config/iris.yml"))

	if err := app.Run(iris.Server(&http.Server{Addr: ":9090"}), config); err != nil {
		logrus.Error(err)
	}

	/*if err := app.Run(iris.Addr("119.23.74.49:7654"), config); err != nil {
		logrus.Error(err)
	}*/

}

func checkRegisterFormat(ctx iris.Context, username, email, mobile, iso, password string) bool {
	if username == "" {
		var res models.ProtocolRsp
		res.Code = models.RegisterErrCode
		res.Msg = models.RegisterUserNameEmptyErr
		res.ResponseWriter(ctx)
		return false
	} else if !utils.IsUserName(username) {
		var res models.ProtocolRsp
		res.Code = models.RegisterErrCode
		res.Msg = models.RegisterUserNameFormatErr
		res.ResponseWriter(ctx)
		return false
	}
	if email == "" {
		var res models.ProtocolRsp
		res.Code = models.RegisterErrCode
		res.Msg = models.RegisterEmailEmptyErr
		res.ResponseWriter(ctx)
		return false
	} else if !utils.IsEmail(email) {
		var res models.ProtocolRsp
		res.Code = models.RegisterErrCode
		res.Msg = models.RegisterEmailFormatErr
		res.ResponseWriter(ctx)
		return false
	}

	if mobile == "" || iso == "" {
		var res models.ProtocolRsp
		res.Code = models.RegisterErrCode
		res.Msg = models.RegisterMobileEmptyErr
		res.ResponseWriter(ctx)
		return false
	} else if !utils.IsMobile(mobile, iso) {
		var res models.ProtocolRsp
		res.Code = models.RegisterErrCode
		res.Msg = models.RegisterMobileFormatErr
		res.ResponseWriter(ctx)
		return false
	}

	if password == "" {
		var res models.ProtocolRsp
		res.Code = models.RegisterErrCode
		res.Msg = models.RegisterPassWordEmptyErr
		res.ResponseWriter(ctx)
		return false
	} else if !utils.IsPwd(password) {
		var res models.ProtocolRsp
		res.Code = models.RegisterErrCode
		res.Msg = models.RegisterPassWordFormatErr
		res.ResponseWriter(ctx)
		return false
	}

	return true
}

func checkLoginFormat(ctx iris.Context, username, email, password string) bool {
	if password == "" {
		var res models.ProtocolRsp
		res.Code = models.LoginErrCode
		res.Msg = models.LoginErrPassWordEmptyErr
		res.ResponseWriter(ctx)
		return false

	} else if !utils.IsPwd(password) {
		var res models.ProtocolRsp
		res.Code = models.LoginErrCode
		res.Msg = models.LoginErrPassWordFormatErr
		res.ResponseWriter(ctx)
		return false
	}

	if username == "" || email == "" {
		var res models.ProtocolRsp
		res.Code = models.LoginErrCode
		res.Msg = models.LoginErrUserNameOrEmailEmptyErr
		res.ResponseWriter(ctx)
		return false
	} else if username == "" && !utils.IsEmail(email) {
		var res models.ProtocolRsp
		res.Code = models.LoginErrCode
		res.Msg = models.LoginEmailFormatErr
		res.ResponseWriter(ctx)
		return false
	} else if email == "" && !utils.IsUserName(username) {
		var res models.ProtocolRsp
		res.Code = models.LoginErrCode
		res.Msg = models.LoginUserNameFormatErr
		res.ResponseWriter(ctx)
		return false
	}

	return true
}
