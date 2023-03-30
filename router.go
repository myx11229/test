package main

import (
	. "work/api"

	"github.com/gin-gonic/gin"
)

func initRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/", IndexUsers)

	users := router.Group("")
	{
		users.GET("/CompleteNumber", CompleteNumber_R)
		users.GET("/BlockNumber/:number", QueryBlockByNumber_R)
		users.GET("/Query30m", Query30m_R)
		users.GET("/Query60m", Query60m_R)
		users.GET("/QueryAll", QueryAll_R)
	}
	return router
}
