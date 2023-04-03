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
	_, err = SqlDB.Exec("CREATE TABLE IF NOT EXISTS block(Number BIGINT NOT NULL,gasLimit BIGINT,timestamp BIGINT,timestamp_readable VARCHAR(20),PRIMARY KEY(Number))")
	if err != nil {
		log.Fatalln(err)
	}

	_, err = SqlDB.Exec("CREATE TABLE IF NOT EXISTS trans(nonce BIGINT,gasprice VARCHAR(50),gas BIGINT,from1 VARCHAR(50),to1 VARCHAR(50),value1 VARCHAR(50),PRIMARY KEY(nonce, from1))")
	if err != nil {
		log.Fatalln(err)
	}

	_, err = SqlDB.Exec("CREATE TABLE IF NOT EXISTS map30(Sender VARCHAR(20),Value INT,PRIMARY KEY(Sender))")
	if err != nil {
		log.Fatalln(err)
	}
	_, err = SqlDB.Exec("CREATE TABLE IF NOT EXISTS map30ht(Id INT auto_increment,head INT,tail INT,PRIMARY KEY(Id))")
	if err != nil {
		log.Fatalln(err)
	}

	_, err = SqlDB.Exec("CREATE TABLE IF NOT EXISTS map60(Sender VARCHAR(20),Value INT,PRIMARY KEY(Sender))")
	if err != nil {
		log.Fatalln(err)
	}
	_, err = SqlDB.Exec("CREATE TABLE IF NOT EXISTS map60ht(Id INT auto_increment,head INT,tail INT,PRIMARY KEY(Id))")
	if err != nil {
		log.Fatalln(err)
	}
}
