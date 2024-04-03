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
