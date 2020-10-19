package data

import (
	"database/sql"
	"errors"
	"github.com/yyangc/todo-list/libs"
	"strconv"
	"time"
)

type User struct {
	ID        uint64         `json:"id"`
	UserName  string         `json:"username"`
	Password  string         `json:"password"`
	Mail      string         `json:"mail"`
	Status    int8           `json:"status"`
	CreatedDt time.Time      `json:"-"`
	UpdatedDt sql.NullString `json:"-"`
}

func (db *DB) GetUserInfo(m string) (*User, error) {
	user := new(User)
	row := db.Ms.QueryRow("SELECT id, username, password FROM users WHERE mail = ?", m)
	if err := row.Scan(&user.ID, &user.UserName, &user.Password); err != nil {
		db.l.Warn(err)
		return nil, err
	}
	return user, nil
}

func (db *DB) CreateRedisAuth(userId uint64, td *libs.TokenDetails) error {
	at := time.Unix(td.AtExpires, 0)
	now := time.Now()

	err := db.R.Set(string(td.AccessUuid[:]), strconv.Itoa(int(userId)), at.Sub(now)).Err()
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GetRedisAuth(ad *libs.AccessDetails) (uint64, error) {
	userid, err := db.R.Get(ad.AccessUuid).Result()
	if err != nil {
		return 0, err
	}
	userID, _ := strconv.ParseUint(userid, 10, 64)
	if ad.UserId != userID {
		return 0, errors.New("unauthorized")
	}
	return userID, nil
}

func (db *DB) DeleteRedisAuth(ad *libs.AccessDetails) error {
	_, err := db.R.Del(ad.AccessUuid).Result()
	if err != nil {
		return err
	}
	return nil
}
