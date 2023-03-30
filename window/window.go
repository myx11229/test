package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"math/big"
	"time"
	"work/db"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const N = 5

type Window struct {
	mp   map[string]int
	Head int //首块号
	Tail int //尾块号
	Diff int //查询时间长度(30 or 60)
}

type Queue struct {
	Id     int    `json:"id" form:"id"`
	Sender string `json:"sender" form:"sender"`
	Number int    `json:"number" form:"number"`
}

type HT struct {
	Id   int    `json:"id" form:"id"`
	Head string `json:"head" form:"head"`
	Tail int    `json:"tail" form:"tail"`
}

var Window30m Window
var Window60m Window

func (window *Window) UpdateInfo(block *types.Block, op int) {
	for i, tx := range block.Transactions() {
		if i >= N {
			break
		}
		from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
		if err == nil {
			Sender := from.Hex()
			Sender = Sender[:3]
			window.mp[Sender] += op
			var err error
			if window.Diff == 30 {
				_, err = db.SqlDB.Exec("REPLACE INTO map30 (sender,number) VALUE(?,?)", Sender, window.mp[Sender])
			} else {
				_, err = db.SqlDB.Exec("REPLACE INTO map60 (sender,number) VALUE(?,?)", Sender, window.mp[Sender])
			}
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
}

func (window *Window) Build(diff int, NowNumber int) {
	var rows *sql.Rows
	var err error
	window.mp = make(map[string]int)
	window.Head = NowNumber
	window.Tail = NowNumber - 1
	window.Diff = diff

	if diff == 30 {
		rows, err = db.SqlDB.Query("select * from map30")
	} else {
		rows, err = db.SqlDB.Query("select * from map60")
	}
	defer rows.Close()
	if err != nil {
		log.Fatalln(err)
	}
	for rows.Next() {
		var Id int
		var Sender string
		var Value int
		err := rows.Scan(&Id, &Sender, &Value)
		if err != nil {
			log.Fatalln(err)
		}
		window.mp[Sender] = Value
	}

	if diff == 30 {
		rows, err = db.SqlDB.Query("select * from map30ht")
	} else {
		rows, err = db.SqlDB.Query("select * from map60ht")
	}
	defer rows.Close()
	if err != nil {
		log.Fatalln(err)
	}
	for rows.Next() {
		var Id int
		var head int
		var tail int
		err := rows.Scan(&Id, &head, &tail)
		if err != nil {
			log.Fatalln(err)
		}
		window.Head = head
		window.Tail = tail
	}

	client, err := ethclient.Dial("wss://mainnet.infura.io/ws/v3/9aa3d95b3bc440fa88ea12eaa4456161")
	if err != nil {
		log.Fatalln(err)
	}
	block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(window.Tail)))
	if err != nil {
		log.Fatalln(err)
	}
	NowTime := uint64(time.Now().Unix())
	if NowTime-block.Time() > uint64(diff*60) {
		//窗口尾块与同步块号时间间隔过大则直接删除所有记录从同步块开始
		window.mp = make(map[string]int)
		window.Head = NowNumber
		window.Tail = NowNumber - 1
		if window.Diff == 30 {
			_, err = db.SqlDB.Exec("drop table mp30")
			_, err = db.SqlDB.Exec("CREATE TABLE IF NOT EXISTS map30(Id int auto_increment,Sender VARCHAR(20),number int,PRIMARY KEY(Id))")
		} else {
			_, err = db.SqlDB.Exec("drop table mp60")
			_, err = db.SqlDB.Exec("CREATE TABLE IF NOT EXISTS map60(Id int auto_increment,Sender VARCHAR(20),number int,PRIMARY KEY(Id))")
		}
		if err != nil {
			log.Fatalln(err)
		}
	}

	if window.Diff == 30 {
		_, err = db.SqlDB.Exec("REPLACE INTO map30ht (id,head,tail) VALUE(1,?,?)", window.Head, window.Tail)
	} else {
		_, err = db.SqlDB.Exec("REPLACE INTO map60ht (id,head,tail) VALUE(1,?,?)", window.Head, window.Tail)
	}
}

func (window *Window) Update(BlockNumber int) {
	client, err := ethclient.Dial("wss://mainnet.infura.io/ws/v3/9aa3d95b3bc440fa88ea12eaa4456161")
	if err != nil {
		log.Fatalln(err)
	}
	block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(BlockNumber)))
	if err != nil {
		log.Fatalln(err)
	}
	window.UpdateInfo(block, 1)
	window.Tail++
	if window.Diff == 30 {
		_, err = db.SqlDB.Exec("REPLACE INTO map30ht (id,head,tail) VALUE(1,?,?)", window.Head, window.Tail)
	} else {
		_, err = db.SqlDB.Exec("REPLACE INTO map60ht (id,head,tail) VALUE(1,?,?)", window.Head, window.Tail)
	}

	NowTime := uint64(time.Now().Unix())
	for {
		if window.Head > window.Tail+1 {
			break
		}
		block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(window.Head)))
		if err != nil {
			log.Fatalln(err)
		}
		if NowTime-block.Time() <= uint64(window.Diff*60) {
			break
		}
		window.UpdateInfo(block, -1)
		window.Head++
		if window.Diff == 30 {
			_, err = db.SqlDB.Exec("REPLACE INTO map30ht (id,head,tail) VALUE(1,?,?)", window.Head, window.Tail)
		} else {
			_, err = db.SqlDB.Exec("REPLACE INTO map60ht (id,head,tail) VALUE(1,?,?)", window.Head, window.Tail)
		}
	}
}

func WindowUpdate(BlockNumber int) {
	Window30m.Update(BlockNumber)
	Window60m.Update(BlockNumber)
}

func (window *Window) Query() (res string) {
	if window.Head > window.Tail {
		return
	}
	resjson, _ := json.MarshalIndent(window.mp, "", "")
	res = string(resjson)
	return
}

func WindowInit(NowNumber int) {
	Window30m.Build(30, NowNumber)
	Window60m.Build(60, NowNumber)
}
