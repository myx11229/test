package models

import (
	"log"
	"work/db"
)

type BlockData struct {
	Number             int    `json:"number" form:"number"`
	GasLimit           int    `json:"gaslimit" form:"gaslimit"`
	Timestamp          int    `json:"timestamp" form:"timestamp"`
	Timestamp_Readable string `json:"timestamp_readable" form:"timestamp_readable"`
}

func (data *BlockData) InsertBlock() { //插入块信息到表block中
	rows, err := db.SqlDB.Query("REPLACE INTO block VALUE(?,?,?,?)", data.Number, data.GasLimit, data.Timestamp, data.Timestamp_Readable)
	defer rows.Close()
	if err != nil {
		log.Fatalln(err)
	}
}
