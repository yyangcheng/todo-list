package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/yyangc/todo-list/config"
	"github.com/yyangc/todo-list/data"
	"github.com/yyangc/todo-list/libs"
	"net/http"
	"strconv"
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
		"accessToken":  td.AccessToken,
		"refreshToken": td.RefreshToken,
	}
	ResJSON(c, http.StatusOK, &Response{Data: data})
}

func (h *Handler) Refresh(c *gin.Context) {
	mapToken := map[string]string{}
	if err := c.ShouldBindJSON(&mapToken); err != nil {
		ResERROR(c, http.StatusUnprocessableEntity, err)
		return
	}
	refreshToken := mapToken["refreshToken"]

	//verify the token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Env.JWT.RefreshSecret), nil
	})

	//if there is an error, the token must have expired
	if err != nil {
		h.l.Warnln("the error: ", err)
		ResERROR(c, http.StatusUnauthorized, errors.New("Refresh token expired"))
		return
	}

	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		ResERROR(c, http.StatusUnauthorized, err)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if !ok || !token.Valid {
		ResERROR(c, http.StatusUnauthorized, errors.New("refresh expired"))
		return
	}

	refreshUuid, ok := claims["refreshUuid"].(string)
	if !ok {
		ResERROR(c, http.StatusUnprocessableEntity, err)
	}

	userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["userId"]), 10, 64)
	if err != nil {
		ResERROR(c, http.StatusUnprocessableEntity, errors.New("Error occurred"))
		return
	}

	//Delete the previous Refresh Token
	if err = h.db.DeleteRedis(refreshUuid); err != nil {
		ResERROR(c, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	//Create new pairs of refresh and access tokens
	ts, err := libs.CreateToken(userId)
	if err != nil {
		ResERROR(c, http.StatusForbidden, err)
		return
	}

	//save the tokens metadata to redis
	err = h.db.CreateRedisAuth(userId, ts)
	if err != nil {
		ResERROR(c, http.StatusForbidden, err)
		return
	}

	tokens := map[string]string{
		"accessToken":  ts.AccessToken,
		"refreshToken": ts.RefreshToken,
	}

	ResJSON(c, http.StatusCreated, &Response{Data: tokens})
}

func (h *Handler) Logout(c *gin.Context) {
	metadata, err := libs.ExtractTokenMetadata(c.Request)
	if err != nil {
		ResERROR(c, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}
	if err := h.db.DeleteRedis(metadata.AccessUuid); err != nil {
		ResERROR(c, http.StatusUnauthorized, err)
		return
	}
	ResJSON(c, http.StatusOK, &Response{Message: "Successfully logged out"})
}
