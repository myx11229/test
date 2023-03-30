package api

import (
	"net/http"
	"strconv"

	"work/models"

	"github.com/gin-gonic/gin"
)

func IndexUsers(c *gin.Context) {
	c.String(http.StatusOK, "It works")
}

func CompleteNumber_R(c *gin.Context) {
	//1. 输出系统当前已经完成同步了的区块高度
	res, _ := models.CompleteNumber()
	c.JSON(http.StatusOK, gin.H{
		"CompleteNumber": res,
	})
}

func QueryBlockByNumber_R(c *gin.Context) {
	//2. 根据请求的区块高度，从数据库返回该区块的数据
	number_string := c.Param("number")
	number, _ := strconv.Atoi(number_string)
	res, _ := models.QueryBlockByNumber(number)
	c.JSON(http.StatusOK, gin.H{
		"result": res,
	})
}

func Query30m_R(c *gin.Context) {
	//3. 返回当前时刻起，半小时内，各类地址的交易发送量
	res, _ := models.Query30m()
	c.JSON(http.StatusOK, gin.H{
		"result": res,
	})
}

func Query60m_R(c *gin.Context) {
	res, _ := models.Query60m()
	//4. 返回当前时刻起，一小时内，各类地址的交易发送量
	c.JSON(http.StatusOK, gin.H{
		"result": res,
	})
}

func QueryAll_R(c *gin.Context) {
	//5. 输出自系统运行以来，以首字母作为区分的各类地址的总的交易发送量
	res, _ := models.QueryAll()
	c.JSON(http.StatusOK, gin.H{
		"result": res,
	})
}
