package model

import (
	"time"
)

type Message struct {
	Id       int       `db:"id"`
	Data     string    `db:"data"`
	Created  time.Time `db:"created"`
	Modified time.Time `db:"modified"`
}
