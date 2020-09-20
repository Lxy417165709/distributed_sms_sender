package model

import "github.com/jinzhu/gorm"

type Message struct{
	gorm.Model
	SenderId string
	ReceiverId string
	Content string
}

func (m *Message)TableName()string {
	return "message"
}
