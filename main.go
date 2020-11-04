package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/yyangc/todo-list/config"
	"github.com/yyangc/todo-list/data"
	"github.com/yyangc/todo-list/handlers"
	"github.com/yyangc/todo-list/libs"
	"github.com/yyangc/todo-list/middlewares"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	// initialize log
	l := initLogger()

	// connect redis
	rdb, err := libs.InitRedis()
	if err != nil {
		l.Warn(err)
	}
	// connect mysql
	mdb, err := libs.InitMysql()
	if err != nil {
		l.Warn(err)
	}

	// initialize data
	db, err := data.New(l, rdb, mdb)
	if err != nil {
		l.Warn(err)
	}
	defer db.Close()

	// initialize Handler
	h := handlers.New(l, db)

	// routers
	mw := middlewares.New(l, db)
	r := gin.Default()
	r.POST("/login", h.Login)
	r.POST("/refresh", h.Refresh)

	r.Use(mw.Auth())
	{
		r.GET("/todo/list", h.GetAllList)

		r.POST("/todo/list", h.CreateList)
		r.POST("/todo/list/:listId/item", h.CreateItem)
		r.POST("/logout", h.Logout)

		r.PUT("/todo/list/:listId/item/:itemId", h.UpdateItem)

		r.DELETE("/todo/list/:listId/item/:itemId", h.DeleteItem)
	}

	srv := &http.Server{
		Addr:         config.Env.ListenHost + ":" + config.Env.ListenPort,
		Handler:      r,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			l.Error(err.Error())
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, os.Kill)

	sig := <-sigChan
	l.Info("Received terminate, graceful shutdown", sig)
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	if err := srv.Shutdown(tc); err != nil {
		l.Warning("HTTP server Shutdown: %v", err.Error())
	}
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
