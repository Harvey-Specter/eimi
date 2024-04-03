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
	rows, err := db.Queryx("select * from sys_dept")

	for rows.Next() {
		err = rows.MapScan(rs)
		fmt.Println("GetByID", rs)
	}

	return rs, err
}
