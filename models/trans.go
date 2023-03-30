package models

import (
	"log"
	"work/db"
)

type Trans struct {
	Id        int    `json:"id" form:"id"`
	Sender    string `json:"sender" form:"sender"`
	Timestamp string `json:"timestamp" form:"timestamp"`
}

func (data *Trans) InsertTrans() { //插入交易信息到表trans中
	rows, err := db.SqlDB.Query("REPLACE INTO trans (sender, timestamp) VALUE(?,?)", data.Sender, data.Timestamp)
	defer rows.Close()
	if err != nil {
		log.Fatalln(err)
	}
}
