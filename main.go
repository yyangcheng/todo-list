package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/yyangc/todo-list/config"
	"github.com/yyangc/todo-list/data"
	"github.com/yyangc/todo-list/handlers"
	"github.com/yyangc/todo-list/libs"
	"github.com/yyangc/todo-list/middlewares"
	"io"
	"os"
)

func main() {
	// initialize log
	l := initLogger()

	// connect redis
	rdb, err := libs.InitRedis()
	if err != nil {
		l.Fatalln(err)
	}
	// connect mysql
	mdb, err := libs.InitMysql()
	if err != nil {
		l.Fatalln(err)
	}

	// initialize data
	db, err := data.New(l, rdb, mdb)
	if err != nil {
		l.Fatalln(err)
	}
	defer db.Close()

	// initialize Handler
	h := handlers.New(l, db)

	// routers
	mw := middlewares.New(l, db)
	r := gin.Default()
	r.POST("/login", h.Login)

	r.Use(mw.Auth())
	{
		r.GET("/todo/list", h.GetAllList)

		r.POST("/todo/list", h.CreateList)
		r.POST("/todo/list/:listId/item", h.CreateItem)
		r.POST("/logout", h.Logout)

		r.PUT("/todo/list/:listId/item/:itemId", h.UpdateItem)

		r.DELETE("/todo/list/:listId/item/:itemId", h.DeleteItem)
	}

	r.Run(":" + config.Env.ListenPort)
}

func initLogger() *logrus.Logger {
	l := logrus.New()
	l.SetLevel(logrus.DebugLevel)
	l.SetFormatter(&logrus.JSONFormatter{})
	l.SetReportCaller(true)
	mw := io.MultiWriter(os.Stdout)
	l.SetOutput(mw)
	return l
}
