package database

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Gamemode int64

const (
	Countdown Gamemode = iota
	Amount
)

type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Name     string             `json:"name" bson:"name"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"-" bson:"password"`
}

type Run struct {
	UserID         primitive.ObjectID `bson:"user_id" json:"user_id"`
	Gamemode       Gamemode           `bson:"game_mode" json:"gamemode"`
	TimeMs         int64              `bson:"time_ms" json:"time_ms"`
	CorrectWords   int8               `bson:"correct_words" json:"correct_words"`
	IncorrectWords int8               `bson:"incorrect_words" json:"incorrect_words"`
	Seed           int64              `bson:"seed" json:"seed"`
	Date           time.Time          `bson:"date"`
}
