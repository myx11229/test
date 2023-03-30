package main

import (
	"work/db"
)

func main() {
	defer db.SqlDB.Close()
	router := initRouter()

	go func() {
		Chain()
	}()

	router.Run(":8810")
}
