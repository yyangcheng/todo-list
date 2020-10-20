package libs

import (
	"database/sql"
	"fmt"
	"github.com/yyangc/todo-list/config"
)

func InitMysql() (*sql.DB, error) {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=true&multiStatements=true",
		config.Env.MySQL.User,
		config.Env.MySQL.Password,
		config.Env.MySQL.Host,
		config.Env.MySQL.Port,
		config.Env.MySQL.Name)
	db, err := sql.Open("mysql", connString)
	if err != nil {
		return nil, err
	}
	return db, nil
}
