package handlers

import (
	"github.com/sirupsen/logrus"
	"github.com/yyangc/todo-list/data"
)

type Handler struct {
	l  *logrus.Logger
	db *data.DB
}

func New(l *logrus.Logger, db *data.DB) *Handler {
	return &Handler{l: l, db: db}
}
