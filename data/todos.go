package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type List struct {
	Id        uint64         `json:"id,string"`
	Title     string         `json:"title"`
	CreatedDt time.Time      `json:"-"`
	UpdatedDt sql.NullString `json:"-"`
}

type ListUser struct {
	LId       uint64    `json:"l_id,string"`
	UId       uint64    `json:"u_id,string"`
	CreatedDt time.Time `json:"-"`
}

type ListItem struct {
	Id        uint64         `json:"id,string"`
	LId       uint64         `json:"l_id,string"`
	Title     string         `json:"title"`
	Content   string         `json:"content"`
	Status    int8           `json:"status"`
	CreatedDt time.Time      `json:"-"`
	UpdatedDt sql.NullString `json:"-"`
}

var UserLists map[uint64]bool

// 取得用戶之待辦事項板
func (db *DB) GetAuthList(id uint64) (*map[uint64]bool, error) {
	key := fmt.Sprintf("id-%v-lists", id)
	// get auth from redis
	res, err := db.R.Get(key).Result()
	if err == nil {
		json.Unmarshal([]byte(res), &UserLists)
		db.l.Debugln("get list from redis", UserLists)
		return &UserLists, nil
	}

	// get auth from db
	rows, err := db.Ms.Query("SELECT l_id FROM list_user WHERE u_id = ?", id)

	if err != nil {
		db.l.Warn(err)
		return nil, err
	}
	defer rows.Close()
	UserLists := make(map[uint64]bool)
	for rows.Next() {
		var listId uint64
		if err := rows.Scan(&listId); err != nil {
			db.l.Warn(err)
			return nil, err
		}
		UserLists[listId] = true
	}
	db.l.Debugln("get list from db", UserLists)
	db.SaveAuthList(id, &UserLists)
	return &UserLists, nil
}

// 將用戶之待辦事項板存至 redis
func (db *DB) SaveAuthList(id uint64, userLists *map[uint64]bool) error {
	key := fmt.Sprintf("id-%v-lists", id)
	jsonStr, _ := json.Marshal(userLists)
	if _, err := db.R.Set(key, string(jsonStr), time.Minute*30*60).Result(); err != nil {
		db.l.Warn(err)
		return err
	}
	db.l.Info("SaveAuthList", userLists)
	return nil
}

// 檢查用戶是否有此待辦事項板之權限
func (db *DB) CheckTodoAuth(id uint64, listId uint64) error {
	var lu ListUser
	err := db.Ms.QueryRow("SELECT l_id, u_id FROM list_user WHERE l_id = ? and u_id = ?", listId, id).Scan(&lu.LId, &lu.UId)
	return err
}

func (db *DB) GetUserAllList(id uint64) ([]*List, error) {
	rows, err := db.Ms.Query("SELECT lu.l_id, l.title FROM list_user AS lu LEFT JOIN list as l ON lu.l_id = l.id WHERE lu.u_id = ?", id)
	if err != nil {
		db.l.Warn(err)
		return nil, err
	}
	defer rows.Close()
	lists := make([]*List, 0)
	for rows.Next() {
		list := &List{}
		if err := rows.Scan(&list.Id, &list.Title); err != nil {
			db.l.Warn(err)
			return nil, err
		}
		lists = append(lists, list)

	}
	return lists, nil
}

func (db *DB) CreateList(l *List) error {
	res, err := db.Ms.Exec("INSERT INTO list (title) VALUES (?)", l.Title)
	if err != nil {
		db.l.Warn(err)
		return err
	}
	id, _ := res.LastInsertId()
	if err != nil {
		return err
	}
	l.Id = uint64(id)
	return nil
}

// 新增待辦事項成員
func (db *DB) CreateListUser(lu *ListUser) error {
	_, err := db.Ms.Exec("INSERT INTO list_user (l_id, u_id) VALUES (?, ?)", lu.LId, lu.UId)
	if err != nil {
		db.l.Warn(err)
		return err
	}

	list, _ := db.GetAuthList(lu.UId)
	(*list)[lu.LId] = true
	db.SaveAuthList(lu.UId, list)
	return nil
}

// 新增待辦列表的事項
func (db *DB) CreateItem(l *ListItem) error {
	res, err := db.Ms.Exec("INSERT INTO list_item (l_id, title, content) VALUES (?, ?, ?)",
		l.LId, l.Title, l.Content)
	if err != nil {
		db.l.Warn(err)
		return err
	}
	id, _ := res.LastInsertId()
	if err != nil {
		return err
	}
	l.Id = uint64(id)
	return nil
}

// 新增待辦列表的項目
func (db *DB) UpdateItem(l *ListItem) error {
	_, err := db.Ms.Exec("UPDATE list_item SET title = ?, content = ?, status = ? WHERE id = ?",
		l.Title, l.Content, l.Status, l.Id)
	if err != nil {
		db.l.Warn(err)
		return err
	}
	return nil
}

// 刪除待辦列表的項目
func (db *DB) DeleteItem(id uint64) error {
	_, err := db.Ms.Exec("DELETE FROM list_item WHERE id = ?", id)
	if err != nil {
		db.l.Warn(err)
		return err
	}
	return nil
}
