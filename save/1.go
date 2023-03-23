package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	_ "github.com/go-sql-driver/mysql"
)

const N = 5

func main() {
	client, err := ethclient.Dial("wss://mainnet.infura.io/ws/v3/9aa3d95b3bc440fa88ea12eaa4456161")
	if err != nil {
		log.Fatal(err)
	}

	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("mysql", "root:myx12345@tcp(127.0.0.1:3306)/test")
	db.Ping()
	defer db.Close()
	if err != nil {
		fmt.Println("database link fail")
		log.Fatalln(err)
	}

	fmt.Println("database link")
	/*
		block{
			Number int key
			gasLimit int
			timestamp int
			timestamp_readable string
		}
	*/
	/*_, err = db.Exec("CREATE TABLE data(Number BIGINT NOT NULL,gasLimit BIGINT,timestamp BIGINT,timestamp_readable VARCHAR(20),PRIMARY KEY(Number));")
	if err != nil {
		fmt.Println("create fail")
		log.Fatalln(err)
	}
	fmt.Println("create")*/

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case header := <-headers:
			//fmt.Println(header.Hash().Hex())

			block, err := client.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("find block")
			Number := block.Number()
			Number_String := Number.String()
			Number_Int, err := strconv.Atoi(Number_String)
			if err != nil {
				fmt.Println("转化错误")
				log.Fatalln(err)
			}
			GasLimit := block.GasLimit()
			Timestamp := block.Time()
			Timestamp_Readable := epochToHumanReadable(int64(block.Time()))
			//Timestamp_Readable_String = Timestamp_Readable.String()
			Timestamp_Readable_String := Timestamp_Readable.Format("2006-01-02 15:04:05")
			_, err = db.Query("INSERT INTO data VALUES(?,?,?,?)", Number_Int, GasLimit, Timestamp, Timestamp_Readable_String)
			if err != nil {
				fmt.Println("insert fail")
				log.Fatalln(err)
			}
			//fmt.Println("insert")
			for i, tx := range block.Transactions() {

				if i >= N {
					break
				}

				//fmt.Println(i)
				from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
				if err == nil {
					fmt.Println(from.Hex())
					s := from.Hex() //发送者地址 0x……
					s = s[2:len(s)]
				}

			}

		}
	}

}

func epochToHumanReadable(epoch int64) time.Time {
	return time.Unix(epoch, 0)
}

/*func QueryNumber() {

}

func QueryBlockData() {

}

func Query30m() {

}

func Query60m() {

}*/
