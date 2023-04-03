package models

import (
	"log"
	"work/db"
)

type Trans struct {
	Nonce    uint64 `json:"nonce" form:"nonce"`
	Gasprice string `json:"gasprice" form:"gasprice"`
	Gas      uint64 `json:"gas" form:"gas"`
	From     string `json:"from1" form:"from1"`
	To       string `json:"to1" form:"to1"`
	Value    string `json:"value1" form:"value1"`
}

func (trans *Trans) InsertTrans() { //插入交易信息到表trans中
	rows, err := db.SqlDB.Query("REPLACE INTO trans VALUE(?,?,?,?,?,?)", trans.Nonce, trans.Gasprice, trans.Gas, trans.From, trans.To, trans.Value)
	defer rows.Close()
	if err != nil {
		log.Fatalln(err)
	}
}
