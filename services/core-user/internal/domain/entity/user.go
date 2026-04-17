package entity

import (
	"strconv"
	"time"
)

type UserID int64

type User struct {
	ID             UserID
	Username       string
	Email          string
	HashedPassword string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (id UserID) String() string {
	return strconv.FormatInt(int64(id), 10)
}
