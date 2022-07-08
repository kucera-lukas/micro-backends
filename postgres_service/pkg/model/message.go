package model

import (
	"time"
)

type Message struct {
	Id       uint32    `db:"id"`
	Data     string    `db:"data"`
	Created  time.Time `db:"created"`
	Modified time.Time `db:"modified"`
}
