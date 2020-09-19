package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"net/http"
)

const port = 8080

func main() {
	r := gin.Default()
	r.GET("/test",func(c *gin.Context){
		c.JSON(http.StatusOK,"Running go http server success. :)")
	})
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		logs.Error("Running go http server failed. :|")
		return
	}
}
