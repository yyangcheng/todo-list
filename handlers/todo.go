package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/yyangc/todo-list/data"
	"net/http"
	"strconv"
)

// 取得該用戶之所有代辦事項列表
func (h *Handler) GetAllList(c *gin.Context) {
	account := c.MustGet("account").(uint64)
	list, err := h.db.GetUserAllList(account)
	if err != nil {
		ResERROR(c, http.StatusInternalServerError, err)
		return
	}
	ResJSON(c, http.StatusOK, &Response{Data: list})
}

// 建立代辦事項列表
func (h *Handler) CreateList(c *gin.Context) {
	account := c.MustGet("account").(uint64)
	list := new(data.List)
	if err := c.ShouldBindJSON(list); err != nil {
		ResERROR(c, http.StatusUnprocessableEntity, err)
		return
	}
	h.l.Debug(list)

	if err := h.db.CreateList(list); err != nil {
		ResERROR(c, http.StatusInternalServerError, err)
		return
	}

	listUser := &data.ListUser{
		LId: list.Id,
		UId: account,
	}
	if err := h.db.CreateListUser(listUser); err != nil {
		ResERROR(c, http.StatusInternalServerError, err)
		return
	}

	info := map[string]uint64{
		"id": list.Id,
	}
	ResJSON(c, http.StatusOK, &Response{Data: info})
}

// 建立代辦事項列表之項目
func (h *Handler) CreateItem(c *gin.Context) {
	account := c.MustGet("account").(uint64)
	listItem := new(data.ListItem)
	listItem.LId, _ = strconv.ParseUint(c.Param("listId"), 10, 64)

	if err := c.ShouldBindJSON(listItem); err != nil {
		ResERROR(c, http.StatusUnprocessableEntity, err)
		return
	}

	if err := h.db.CheckTodoAuth(account, listItem.LId); err != nil {
		ResERROR(c, http.StatusUnauthorized, errors.New("Permission denied"))
		return
	}
	if err := h.db.CreateItem(listItem); err != nil {
		ResERROR(c, http.StatusInternalServerError, err)
		return
	}

	info := map[string]uint64{
		"id": listItem.Id,
	}
	ResJSON(c, http.StatusOK, &Response{Data: info})
}

// 更新代辦事項列表之項目
func (h *Handler) UpdateItem(c *gin.Context) {
	account := c.MustGet("account").(uint64)
	listItem := new(data.ListItem)
	listItem.LId, _ = strconv.ParseUint(c.Param("listId"), 10, 64)
	listItem.Id, _ = strconv.ParseUint(c.Param("itemId"), 10, 64)

	if err := c.ShouldBindJSON(listItem); err != nil {
		ResERROR(c, http.StatusUnprocessableEntity, err)
		return
	}
	if err := h.db.CheckTodoAuth(account, listItem.LId); err != nil {
		ResERROR(c, http.StatusUnauthorized, errors.New("Permission denied"))
		return
	}

	if err := h.db.UpdateItem(listItem); err != nil {
		ResERROR(c, http.StatusInternalServerError, err)
		return
	}
	info := map[string]uint64{
		"id": listItem.Id,
	}
	ResJSON(c, http.StatusOK, &Response{Data: info})
}

func (h *Handler) DeleteItem(c *gin.Context) {
	account := c.MustGet("account").(uint64)
	lId, _ := strconv.ParseUint(c.Param("listId"), 10, 64)
	iId, _ := strconv.ParseUint(c.Param("itemId"), 10, 64)

	if err := h.db.CheckTodoAuth(account, lId); err != nil {
		ResERROR(c, http.StatusUnauthorized, errors.New("Permission denied"))
		return
	}
	if err := h.db.DeleteItem(iId); err != nil {
		ResERROR(c, http.StatusInternalServerError, err)
		return
	}
	ResJSON(c, http.StatusOK, &Response{})
}
