package handlers

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/yyangc/todo-list/data"
	"github.com/yyangc/todo-list/libs"
	"net/http"
)

func (h *Handler) Login(c *gin.Context) {
	u := new(data.User)
	if err := c.ShouldBindJSON(u); err != nil {
		ResERROR(c, http.StatusUnprocessableEntity, err)
		return
	}

	user, err := h.db.GetUserInfo(u.Mail)
	if errors.Is(err, sql.ErrNoRows) {
		ResERROR(c, http.StatusUnauthorized, errors.New("user not found"))
		return
	}
	// 還未處理
	if u.Password != user.Password {
		ResERROR(c, http.StatusUnauthorized, errors.New("wrong password"))
		return
	}

	// create jwt token
	td, err := libs.CreateToken(user.ID)
	if err != nil {
		h.l.Error(err)
		ResERROR(c, http.StatusUnprocessableEntity, err)
		return
	}

	// write uuid into redis
	err = h.db.CreateRedisAuth(user.ID, td)
	if err != nil {
		h.l.Error(err)
		ResERROR(c, http.StatusUnprocessableEntity, err)
		return
	}

	data := map[string]string{
		"accessToken": td.AccessToken,
	}
	ResJSON(c, http.StatusOK, &Response{Data: data})
}

func (h *Handler) Logout(c *gin.Context) {
	metadata, err := libs.ExtractTokenMetadata(c.Request)
	if err != nil {
		ResERROR(c, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}
	if err := h.db.DeleteRedisAuth(metadata); err != nil {
		ResERROR(c, http.StatusUnauthorized, err)
		return
	}
	ResJSON(c, http.StatusOK, &Response{Message: "Successfully logged out"})
}
