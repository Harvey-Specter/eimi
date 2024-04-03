package svc

import (
	"fmt"

	"github.com/Harvey-Specter/eimi/db"
)

func GetRecord(cfg map[string]any) {

	conn := db.GetConn(cfg)
	rsMap, err := db.GetByID(conn)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(rsMap)

}

func GetDDL(cfg map[string]any) []db.Table {
	conn := db.GetConn(cfg)
	srcTables := cfg["tables"].([]interface{})
	//srcTables := srcMap["tables"].([]string)
	tables, err := db.GetTables(conn, srcTables)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return tables
}

func ExecCopy(srcMap map[string]any, destMap map[string]any, tables []db.Table) (int, error) {
	srcdb := db.GetConn(srcMap)
	destdb := db.GetConn(destMap)
	affectTableCnt, err := db.SelectAndInsert(srcdb, destdb, tables)
	if err != nil {
		fmt.Println("ExecCopy===", err.Error())
		return 0, err
	}
	return affectTableCnt, nil
}
