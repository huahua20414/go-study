package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	//server := web.RegisterRoutes()
	server := gin.Default()
	server.GET("/login", func(c *gin.Context) {
		c.String(200, "成功")
	})
	server.Run(":8080")
}
