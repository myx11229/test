package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var SqlDB *sql.DB

func init() {
	var err error
	SqlDB, err = sql.Open("mysql", "root:myx12345@tcp(127.0.0.1:3306)/test")
	if err != nil {
		log.Fatalln(err.Error())
	}
	err = SqlDB.Ping()
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = SqlDB.Exec("CREATE TABLE IF NOT EXISTS block(Number BIGINT NOT NULL,gasLimit BIGINT,timestamp BIGINT,timestamp_readable VARCHAR(20),PRIMARY KEY(Number));")
	if err != nil {
		log.Fatalln(err)
	}

	_, err = SqlDB.Exec("CREATE TABLE IF NOT EXISTS trans(Id int auto_increment,Sender VARCHAR(20),timestamp TIMESTAMP,PRIMARY KEY(Id));")
	if err != nil {
		log.Fatalln(err)
	}

	_, err = SqlDB.Exec("CREATE TABLE IF NOT EXISTS map30(Id int auto_increment,Sender VARCHAR(20),number int,PRIMARY KEY(Id))")
	if err != nil {
		log.Fatalln(err)
	}
	_, err = SqlDB.Exec("CREATE TABLE IF NOT EXISTS map30ht(Id int auto_increment,head int,tail int,PRIMARY KEY(Id))")
	if err != nil {
		log.Fatalln(err)
	}

	_, err = SqlDB.Exec("CREATE TABLE IF NOT EXISTS map60(Id int auto_increment,Sender VARCHAR(20),number int,PRIMARY KEY(Id))")
	if err != nil {
		log.Fatalln(err)
	}
	_, err = SqlDB.Exec("CREATE TABLE IF NOT EXISTS map60ht(Id int auto_increment,head int,tail int,PRIMARY KEY(Id))")
	if err != nil {
		log.Fatalln(err)
	}
}
