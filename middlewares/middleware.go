package middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/yyangc/todo-list/data"
	"github.com/yyangc/todo-list/handlers"
	"github.com/yyangc/todo-list/libs"
	"net/http"
)

type Middleware struct {
	l  *logrus.Logger
	db *data.DB
}

func New(l *logrus.Logger, db *data.DB) *Middleware {
	return &Middleware{l: l, db: db}
}

func (mw *Middleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		metadata, err := libs.ExtractTokenMetadata(c.Request)
		mw.l.Debugln("metadata : ", metadata)
		if err != nil {
			mw.l.Warnln("metadata err: ", err)
			handlers.ResERROR(c, http.StatusUnauthorized, errors.New("unauthorized"))
			c.Abort()
			return
		}
		uid, err := mw.db.GetRedisAuth(metadata)
		if err != nil {
			mw.l.Warnln("GetRedisAuth err : ", err)
			mw.l.Warnln("GetRedisAuth uid : ", uid)
			handlers.ResERROR(c, http.StatusUnauthorized, errors.New("unauthorized"))
			c.Abort()
			return
		}
		c.Set("account", uid)
		c.Next()
	}
}
