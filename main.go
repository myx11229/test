package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

const N = 5

var db *sql.DB

type BlockData struct {
	Number             int    `json:"number" form:"number"`
	GasLimit           int    `json:"gaslimit" form:"gaslimit"`
	Timestamp          int    `json:"timestamp" form:"timestamp"`
	Timestamp_Readable string `json:"timestamp_readable" form:"timestamp_readable"`
}

type Trans struct {
	Id        int    `json:"id" form:"id"`
	Sender    string `json:"sender" form:"sender"`
	Timestamp string `json:"timestamp" form:"timestamp"`
}

var CompleteNumber int

func main() {

	viper.SetConfigFile("/home/pa/1/config")
	//viper.SetConfigName("st")
	viper.SetConfigType("json")
	StartNumber := viper.GetInt("startnumber")
	router := gin.Default()
	//router.GET("/", handler)
	client, err := ethclient.Dial("wss://mainnet.infura.io/ws/v3/9aa3d95b3bc440fa88ea12eaa4456161")
	if err != nil {
		log.Fatal(err)
	}

	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal(err)
	}

	/*db, err := sql.Open("mysql", "root:myx12345@tcp(127.0.0.1:3306)/test")
	db.Ping()
	defer db.Close()
	if err != nil {
		fmt.Println("database link fail")
		log.Fatalln(err)
	}*/
	db, _ = sql.Open("mysql", "root:myx12345@tcp(127.0.0.1:3306)/test")
	db.Ping()
	defer db.Close()

	//fmt.Println("database link")

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS block(Number BIGINT NOT NULL,gasLimit BIGINT,timestamp BIGINT,timestamp_readable VARCHAR(20),PRIMARY KEY(Number));")
	if err != nil {
		log.Fatalln(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS trans(Id int auto_increment,Sender VARCHAR(20),timestamp TIMESTAMP,PRIMARY KEY(Id));")
	if err != nil {
		log.Fatalln(err)
	}

	//fmt.Println("create")

	router.GET("/CompleteNumber", func(c *gin.Context) {
		//1. 输出系统当前已经完成同步了的区块高度
		c.JSON(http.StatusOK, gin.H{
			"CompleteNumber": QueryNumber(),
		})
	})

	router.GET("/BlockNumber/:number", func(c *gin.Context) {
		//2. 根据请求的区块高度，从数据库返回该区块的数据
		number_string := c.Param("number")
		number, _ := strconv.Atoi(number_string)
		res, _ := QueryBlockByNumber(number)
		c.JSON(http.StatusOK, gin.H{
			"result": res,
		})
	})

	router.GET("/Query30m", func(c *gin.Context) {
		//3. 返回当前时刻起，半小时内，各类地址的交易发送量
		res, _ := Query30m()
		//res := make(map[string]int)
		c.JSON(http.StatusOK, gin.H{
			"result": res,
		})
	})

	router.GET("/Query60m", func(c *gin.Context) {
		//4. 返回当前时刻起，一小时内，各类地址的交易发送量
		res, _ := Query60m()
		//res := make(map[string]int)
		c.JSON(http.StatusOK, gin.H{
			"result": res,
		})
	})

	router.GET("/QueryAll", func(c *gin.Context) {
		//5. 输出自系统运行以来，以首字母作为区分的各类地址的总的交易发送量
		res, _ := QueryAll()
		//res := make(map[string]int)
		c.JSON(http.StatusOK, gin.H{
			"result": res,
		})
	})

	go func() {
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
					fmt.Println("高度转化错误")
					log.Fatalln(err)
				}

				if Number_Int < StartNumber {
					continue
				}

				GasLimit := int(block.GasLimit())
				Timestamp := int(block.Time())
				Timestamp_Readable := epochToHumanReadable(int64(block.Time()))
				//Timestamp_Readable_String = Timestamp_Readable.String()
				Timestamp_Readable_String := Timestamp_Readable.Format("2006-01-02 15:04:05")

				Block := BlockData{
					Number:             Number_Int,
					GasLimit:           GasLimit,
					Timestamp:          Timestamp,
					Timestamp_Readable: Timestamp_Readable_String,
				}
				Block.InsertBlock() //将块信息插入到表block中

				for i, tx := range block.Transactions() {

					if i >= N {
						break
					}

					//fmt.Println(i)
					from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
					if err == nil {
						//fmt.Println(from.Hex())
						Sender := from.Hex() //发送者地址 0x……
						Sender = Sender[:3]
						//Timestamp同区块
						trans := Trans{
							Sender: Sender,
							//Timestamp: Timestamp,
							Timestamp: Timestamp_Readable_String,
						}
						trans.InsertTrans() //将交易信息插入到表trans中
					}

				}
				CompleteNumber = Number_Int

			}
		}
	}()

	router.Run(":8801")

	//endless.ListenAndServe(":8801", router)
}

func epochToHumanReadable(epoch int64) time.Time {
	return time.Unix(epoch, 0)
}

func (data *BlockData) InsertBlock() { //插入块信息到表block中
	_, err := db.Query("INSERT INTO block VALUE(?,?,?,?)", data.Number, data.GasLimit, data.Timestamp, data.Timestamp_Readable)
	if err != nil {
		fmt.Println("Insert Block Fail")
		log.Fatalln(err)
	}
	//fmt.Println("Insert Block Success")
}

func (data *Trans) InsertTrans() { //插入交易信息到表trans中
	//fmt.Println(data.Sender)
	//fmt.Println(data.Timestamp)
	_, err := db.Query("INSERT INTO trans (sender, timestamp) VALUE(?,?)", data.Sender, data.Timestamp)
	if err != nil {
		fmt.Println("Insert Trans Fail")
		log.Fatalln(err)
	}
	//fmt.Println("Insert Trans Success")
}

func QueryNumber() int { //1. 输出系统当前已经完成同步了的区块高度
	//fmt.Println(CompleteNumber)
	return CompleteNumber
}

func QueryBlockByNumber(number int) (block BlockData, err error) {
	//2. 根据请求的区块高度，从数据库返回该区块的数据
	block = BlockData{}
	err = db.QueryRow("select number,gaslimit,timestamp,timestamp_readable from block where number= ?", number).Scan(&block.Number, &block.GasLimit, &block.Timestamp, &block.Timestamp_Readable)
	return
}

func Query30m() (_ string, err error) { //3. 返回当前时刻起，半小时内，各类地址的交易发送量
	return QueryXmin(30)
}

func Query60m() (_ string, err error) { //4. 返回当前时刻起，一小时内，各类地址的交易发送量
	return QueryXmin(60)
}

func QueryXmin(X int) (res string, err error) {
	//返回当前时刻起，X分钟内，各类地址的交易发送量
	mp := make(map[string]int)
	rows, _ := db.Query("select id,sender,timestamp from trans where timestamp>CURRENT_TIMESTAMP-INTERVAL ? MINUTE", X)
	for rows.Next() {
		var Id int
		var Sender string
		var Timestamp string
		err = rows.Scan(&Id, &Sender, &Timestamp)
		mp[Sender]++
		//fmt.Println("Sender:")
		//fmt.Println(Sender)
	}
	resjson, _ := json.MarshalIndent(mp, "", "")
	res = string(resjson)
	return
}

func QueryAll() (_ string, err error) { //5. 输出自系统运行以来，以首字母作为区分的各类地址的总的交易发送量
	return QueryXmin(10000000)
}
