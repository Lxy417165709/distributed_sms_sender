package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type InvokeSituation struct{
	gorm.Model
	UserId string
	InvokeTime time.Time
}

func (*InvokeSituation)TableName() string{
	return "invoke_situation"
}
