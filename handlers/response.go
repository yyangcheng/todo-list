package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func ResJSON(c *gin.Context, statusCode int, res *Response) {
	if res.Message == "" {
		res.Message = http.StatusText(statusCode)
	}
	c.JSON(statusCode, res)
}

func ResERROR(c *gin.Context, statusCode int, err error) {
	if err != nil {
		ResJSON(c, statusCode, &Response{
			Message: err.Error(),
		})
		return
	}
	ResJSON(c, http.StatusBadRequest, &Response{})
}
