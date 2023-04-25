package main

type Gamemode int64

const (
	Countdown Gamemode = iota
	Amount
)

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Runs     []Run  `json:"runs"`
}

type Run struct {
	Gamemode       Gamemode `json:"gameMode"`
	TimeMs         int64    `json:"timeMs"`
	CorrectWords   int8     `json:"correctWords"`
	IncorrectWords int8     `json:"incorrectWords"`
	Seed           int64    `json:"seed"`
}
