package model

import (
	"time"
)

type Certification struct {
	RemainSendTimes int
	LastSendTime time.Time
}

type ResultOfSend struct {
	IsSuccess bool
	Msg string
}
