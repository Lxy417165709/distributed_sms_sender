package controller

import (
	"distributed/sms/src/model"
	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SendMessage(c *gin.Context) {
	var messageForm model.MessageForm
	if err := c.ShouldBindJSON(&messageForm); err != nil {
		logs.Error(err)
		c.JSON(http.StatusBadRequest, &model.Response{
			Err: err.Error(),
		})
	}
	resultOfSend := SingleSmsSender.Send(&model.Message{
		ReceiverId: messageForm.ReceiverId,
		SenderId:   messageForm.SenderId,
		Content:    messageForm.Content,
	})
	if resultOfSend.IsSuccess == false {
		c.JSON(http.StatusBadRequest, &model.Response{
			Err: resultOfSend.Msg,
		})
		return
	}
	c.JSON(http.StatusOK, &model.Response{
		Msg: resultOfSend.Msg,
	})
}
