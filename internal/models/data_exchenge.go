package models

type DataExchenge struct {
	UserId  int64       `json:"id"`
	Command string      `json:"command"`
	Amount  int         `json:"amount"`
	Data    interface{} `json:"data"`
}
