package db

import (
	"bytes"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func GetConn(connMap map[string]any) *sqlx.DB {
	var (
		db       *sqlx.DB
		buf      bytes.Buffer
		dbType   = strings.ToLower(connMap["type"].(string))
		user     = connMap["user"].(string)
		password = connMap["password"].(string)
		host     = connMap["host"].(string)
		port     = connMap["port"].(int)
		dbname   = connMap["db"].(string)
	)

	if dbType == "mysql" {
		mysqlStr := "%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local"
		buf.WriteString(fmt.Sprintf(mysqlStr, user, password, host, port, dbname))
		db = sqlx.MustConnect("mysql", buf.String())
	}
	if dbType == "postgresql" {
		pgStr := "postgres://%s:%s@%s:%d/%s"
		buf.WriteString(fmt.Sprintf(pgStr, user, password, host, port, dbname))
		db = sqlx.MustConnect("sqlx.MustConnect", buf.String())
	}
	// if dbType == "oracle" {

	// }

	return db
}

func GetByID(db *sqlx.DB) (map[string]interface{}, error) {
	// var rs map[string]string

	var rs = make(map[string]interface{})
	// err := db.QueryRowx("SELECT * FROM user WHERE id = ?", id).StructScan(rs)

	// sql := "select * from sys_dept"
	showCreateTable := "show create table sys_dept"
	rows, err := db.Queryx(showCreateTable)

	for rows.Next() {
		err = rows.MapScan(rs)
		//fmt.Println("GetByID", rs)
	}

	return rs, err
}

type Table struct {
	Table string `json:"table" db:"Table"`
	DDL   string `json:"ddl" db:"Create Table"`
}

func GetTable(db *sqlx.DB) (table Table, err error) {
	err = db.Get(&table, sqlGetTable)
	// err will be nil on success
	return table, err
}

const sqlGetTable = `SHOW CREATE TABLE sys_dept`
