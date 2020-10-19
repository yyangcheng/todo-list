package data

import (
	"database/sql"
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

// 檢查用戶是否有此代辦事項板之權限
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

// 新增代辦事項成員
func (db *DB) CreateListUser(lu *ListUser) error {
	_, err := db.Ms.Exec("INSERT INTO list_user (l_id, u_id) VALUES (?, ?)", lu.LId, lu.UId)
	if err != nil {
		db.l.Warn(err)
		return err
	}
	return nil
}

// 新增代辦列表的事項
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

// 新增代辦列表的項目
func (db *DB) UpdateItem(l *ListItem) error {
	_, err := db.Ms.Exec("UPDATE list_item SET title = ?, content = ?, status = ? WHERE id = ?",
		l.Title, l.Content, l.Status, l.Id)
	if err != nil {
		db.l.Warn(err)
		return err
	}
	return nil
}

// 刪除代辦列表的項目
func (db *DB) DeleteItem(id uint64) error {
	_, err := db.Ms.Exec("DELETE FROM list_item WHERE id = ?", id)
	if err != nil {
		db.l.Warn(err)
		return err
	}
	return nil
}
