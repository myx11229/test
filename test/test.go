package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	//"gopkg.in/gin-gonic/gin.v1"
)

var db *sql.DB

type Person struct {
	Id        int    `json:"id" form:"id"`
	Name      string `json:"name" form:"name"`
	Telephone string `json:"telephone" form:"telephone"`
}

func main() {

	router := gin.Default()
	/*db, err := sql.Open("mysql", "root:myx12345@tcp(127.0.0.1:3306)/test")
	if err != nil {
		log.Fatalln(err)
	}*/
	db, _ = sql.Open("mysql", "root:myx12345@tcp(127.0.0.1:3306)/test")
	defer db.Close()

	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(20)

	if err := db.Ping(); err != nil {
		log.Fatalln(err)
	}

	//增加一条记录
	router.POST("/add", func(c *gin.Context) {
		name := c.Request.FormValue("name")
		telephone := c.Request.FormValue("telephone")
		person := Person{
			Name:      name,
			Telephone: telephone,
		}
		id := person.Create()
		msg := fmt.Sprintf("insert successful %d", id)
		c.JSON(http.StatusOK, gin.H{
			"msg": msg,
		})
	})

	router.GET("/users", func(c *gin.Context) {
		rs, _ := getRows()
		c.JSON(http.StatusOK, gin.H{
			"list": rs,
		})
	})

	router.GET("/users/:id", func(c *gin.Context) {
		id_string := c.Param("id")
		id, _ := strconv.Atoi(id_string)
		rs, _ := getRow(id)
		c.JSON(http.StatusOK, gin.H{
			"result": rs,
		})
	})

	router.POST("/users/update", func(c *gin.Context) {
		ids := c.Request.FormValue("id")
		id, _ := strconv.Atoi(ids)
		telephone := c.Request.FormValue("telephone")
		person := Person{
			Id:        id,
			Telephone: telephone,
		}
		row := person.Update()
		msg := fmt.Sprintf("updated successful %d", row)
		c.JSON(http.StatusOK, gin.H{
			"msg": msg,
		})
	})

	//删除一条记录
	router.POST("/users/del", func(c *gin.Context) {
		ids := c.Request.FormValue("id")
		id, _ := strconv.Atoi(ids)
		row := Delete(id)
		msg := fmt.Sprintf("delete successful %d", row)
		c.JSON(http.StatusOK, gin.H{
			"msg": msg,
		})
	})

	router.Run(":8801")
}

// 插入
func (person *Person) Create() int64 {
	rs, err := db.Exec("INSERT into users (name, telephone) value (?,?)", person.Name, person.Telephone)
	if err != nil {
		log.Fatal(err)
	}
	id, err := rs.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	return id
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// 查询所有记录
func getRows() (persons []Person, err error) {
	rows, err := db.Query("select id,name,telephone from users")
	for rows.Next() {
		person := Person{}
		err := rows.Scan(&person.Id, &person.Name, &person.Telephone)
		if err != nil {
			log.Fatal(err)
		}
		persons = append(persons, person)
	}
	rows.Close()
	return
}

func getRow(id int) (person Person, err error) {
	person = Person{}
	err = db.QueryRow("select id,name,telephone from users where id = ?", id).Scan(&person.Id, &person.Name, &person.Telephone)
	return
}

// 修改
func (person *Person) Update() int64 {
	rs, err := db.Exec("update users set telephone = ? where id = ?", person.Telephone, person.Id)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := rs.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	return rows
}

func Delete(id int) int64 {
	rs, err := db.Exec("delete from users where id = ?", id)
	if err != nil {
		log.Fatal()
	}
	rows, err := rs.RowsAffected()
	if err != nil {
		log.Fatal()
	}
	return rows
}
