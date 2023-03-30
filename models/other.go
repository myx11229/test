package models

import (
	"encoding/json"
	"work/db"
	. "work/window"
)

func CompleteNumber() (res int, err error) { //1. 输出系统当前已经完成同步了的区块高度
	rows, _ := db.SqlDB.Query("select * from block order by number DESC limit 1")
	defer rows.Close()
	for rows.Next() {
		var GasLimit int
		var Timestamp int
		var Timestamp_Readable string
		err = rows.Scan(&res, &GasLimit, &Timestamp, &Timestamp_Readable)
	}
	return
}

func QueryBlockByNumber(number int) (block BlockData, err error) {
	//2. 根据请求的区块高度，从数据库返回该区块的数据
	block = BlockData{}
	err = db.SqlDB.QueryRow("select number,gaslimit,timestamp,timestamp_readable from block where number= ?", number).Scan(&block.Number, &block.GasLimit, &block.Timestamp, &block.Timestamp_Readable)
	return
}

func Query30m() (_ string, err error) { //3. 返回当前时刻起，半小时内，各类地址的交易发送量
	return Window30m.Query(), nil
}

func Query60m() (_ string, err error) { //4. 返回当前时刻起，一小时内，各类地址的交易发送量
	return Window60m.Query(), nil
}

func QueryAll() (res string, err error) {
	//5. 输出自系统运行以来，以首字母作为区分的各类地址的总的交易发送量
	mp := make(map[string]int)
	rows, _ := db.SqlDB.Query("select id,sender,timestamp from trans")
	defer rows.Close()
	for rows.Next() {
		var Id int
		var Sender string
		var Timestamp string
		err = rows.Scan(&Id, &Sender, &Timestamp)
		mp[Sender]++
	}
	resjson, _ := json.MarshalIndent(mp, "", "")
	res = string(resjson)
	return
}
