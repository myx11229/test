package main

import (
	"context"
	"log"
	"math/big"
	"time"

	. "work/models"
	. "work/window"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
)

func epochToHumanReadable(epoch int64) time.Time {
	return time.Unix(epoch, 0)
}

func Chain() {

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	_ = viper.ReadInConfig()
	StartNumber := viper.GetInt("startnumber")

	client, err := ethclient.Dial("wss://mainnet.infura.io/ws/v3/9aa3d95b3bc440fa88ea12eaa4456161")
	if err != nil {
		log.Fatalln(err)
	}

	NowNumber := StartNumber - 1
	Number, err := CompleteNumber()
	if err != nil {
		log.Fatalln(err)
	}
	if NowNumber < Number {
		NowNumber = Number
	}
	NowNumber++
	WindowInit(NowNumber)

	/*header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("new blocknumber:")
	fmt.Println(header.Number.String())*/

	for {
		block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(NowNumber)))
		if err != nil {
			log.Fatal(err)
		}

		GasLimit := int(block.GasLimit())
		Timestamp := int(block.Time())
		Timestamp_Readable := epochToHumanReadable(int64(block.Time()))
		Timestamp_Readable_String := Timestamp_Readable.Format("2006-01-02 15:04:05")

		Block := BlockData{
			Number:             NowNumber,
			GasLimit:           GasLimit,
			Timestamp:          Timestamp,
			Timestamp_Readable: Timestamp_Readable_String,
		}
		Block.InsertBlock() //将块信息插入到表block中

		for i, tx := range block.Transactions() {
			if i >= N {
				break
			}
			from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
			if err == nil {
				Sender := from.Hex() //发送者地址 0x……
				Sender = Sender[:3]
				//Timestamp同区块
				trans := Trans{
					Sender:    Sender,
					Timestamp: Timestamp_Readable_String,
				}
				trans.InsertTrans() //将交易信息插入到表trans中
			}
		}
		WindowUpdate(NowNumber)
		NowNumber++
	}
}
