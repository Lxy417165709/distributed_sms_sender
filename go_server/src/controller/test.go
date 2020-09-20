package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Test(c *gin.Context){
	c.JSON(http.StatusOK,"Running go http server success. :)")
}
