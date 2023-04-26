package main

type Gamemode int64

const (
	Countdown Gamemode = iota
	Amount
)

type User struct {
	Name     string `json:"name" bson:"name"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"-" bson:"password"`
	Runs     []Run  `json:"runs" bson:"runs"`
}

type Run struct {
	Gamemode       Gamemode `bson:"gameMode" json:"gameMode"`
	TimeMs         int64    `bson:"timeMs" json:"timeMs"`
	CorrectWords   int8     `bson:"correctWords" json:"correctWords"`
	IncorrectWords int8     `bson:"incorrectWords" json:"incorrectWords"`
	Seed           int64    `bson:"seed" json:"seed"`
}
