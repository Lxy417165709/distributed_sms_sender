package model


type MessageForm struct {
	SenderId   string `json:"sender_id"`
	ReceiverId string `json:"receiver_id"`
	Content    string `json:"content"`
}
