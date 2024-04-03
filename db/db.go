package db

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

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
	var rs = make(map[string]interface{})
	sql := "select * from sys_dept"
	rows, err := db.Queryx(sql)
	for rows.Next() {
		err = rows.MapScan(rs)
		fmt.Println("GetByID", rs)
	}
	return rs, err
}

func SelectAndInsert(dbsrc *sqlx.DB, dbdest *sqlx.DB, tables []Table) (int, error) {
	cnt := 0
	for _, t := range tables {
		tableName := t.Table
		createSql := t.DDL
		dbdest.MustExec("drop table if exists " + t.Table)
		dbdest.MustExec(createSql)
		seletSql := "select * from " + tableName
		rows, err := dbsrc.Queryx(seletSql)
		if err != nil {
			fmt.Println(err.Error())
			return 0, err
		}
		dataList := []map[string]interface{}{}
		for rows.Next() {
			var dataMap = make(map[string]interface{})
			err = rows.MapScan(dataMap)
			if err != nil {
				fmt.Println(err.Error())
				return 0, err
			}
			// fmt.Println(tableName, "rows.Next", dataMap)
			dataList = append(dataList, dataMap)
		}
		if len(dataList) > 0 {
			batchInsertSQL := genInsertSql(tableName, dataList)
			fmt.Println("batchInsertSQL===", batchInsertSQL)
		}
		// insertStr := "%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local"
		//insertStr := "insert %s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local"
		//buf.WriteString(fmt.Sprintf(mysqlStr, user, password, host, port, dbname))
		//db = sqlx.MustConnect("mysql", buf.String())

	}
	return cnt, nil
}
func genInsertSql(tableName string, dataList []map[string]interface{}) string {
	fmt.Println("tableName=======", tableName)
	insertStrSlice := []string{"insert into", tableName, "("}
	valuesStrSlice := []string{}
	for i, dataMap := range dataList {
		fmt.Println("dataMap=====", dataMap)
		colSlice := []string{}
		// insert into tName (a,b,c,d) values ('11','3','12345',"@@@")
		valueSlice := []string{}
		for k, v := range dataMap {
			if i == 0 {
				colSlice = append(colSlice, k)
			}
			fmt.Println("k=====", k)
			insertValue := interface2String(v)
			if insertValue == "null" {
				valueSlice = append(valueSlice, interface2String(v))
			} else {
				valueSlice = append(valueSlice, strings.Join([]string{"'", interface2String(v), "'"}, ""))
			}
		}
		if i == 0 {
			insertStrSlice = append(insertStrSlice, strings.Join(colSlice, ","), ")", "values")
		}
		valuesStrSlice = append(valuesStrSlice, strings.Join([]string{"(", strings.Join(valueSlice, ","), ")"}, ""))

	}
	insertStrSlice = append(insertStrSlice, strings.Join(valuesStrSlice, ","))
	insertStr := strings.Join(insertStrSlice, " ")

	return insertStr
}

type Table struct {
	Table string `json:"table" db:"Table"`
	DDL   string `json:"ddl" db:"Create Table"`
}

func interface2String(i interface{}) string {
	fmt.Println("i=======", i)
	fmt.Printf("Type of i is %T\n", i)
	if i == nil {
		return "null"
	}
	switch i.(type) {
	case string:
		return i.(string)
		//fmt.Println("string", i.(string))
		//break
	case int64:
		return strconv.FormatInt(i.(int64), 10)
	case int:
		return strconv.Itoa(i.(int))
		// fmt.Println("int", i.(int))
		// break
	case float64:
		str := fmt.Sprintf("%f", i)
		return str
		// fmt.Println("float64", i.(float64))
		// break
	case []uint8:
		return string(i.([]uint8))
	case time.Time:
		if t, ok := i.(time.Time); ok {
			return t.Format("2006-01-02 15:04:05")
		}
		return "badTime"

	default:
		return i.(string)
	}

}
func GetTables(db *sqlx.DB, tableNames []interface{}) (tables []Table, err error) {
	tableList := []Table{}
	table := Table{}
	for _, v := range tableNames {
		sqlGetTable := "SHOW CREATE TABLE " + v.(string)
		err = db.Get(&table, sqlGetTable)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		tableList = append(tableList, table)
	}
	return tableList, err
}
