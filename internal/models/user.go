package models

import "time"

const (
	ALLOW int = iota
	BANNED
)

const (
	USERROLE int = iota
	ADMIN
)

const DEFOULTVALUE string = "DEFOULTVALUE"

type User struct {
	Id              int64     `json:"id"`
	UserId          int64     `json:"user_id"`
	UserName        string    `json:"user_name"`
	MessageType     int       `json:"message_type"`
	Command         string    `json:"command"`
	Status          int       `json:"status"`
	HistoryUserName string    `json:"history_user_name"`
	Balance         int       `json:"balance"`
	BalanceAllTime  int       `json:"balance_all_time"`
	AutoPay         bool      `json:"auto_pay"`
	TestUsed        bool      `json:"test_used"`
	LastMsgId       int64     `json:"last_msg_id"`
	ReferalId       int64     `json:"referal_id"`
	TimeCreate      time.Time `json:"time_create"`
	TimeUpdate      time.Time `json:"time_update"`
	dataForComand   string
	ignoreDeletion  bool
	callbackId      string
}

func NewUser() *User {
	return &User{
		Id:              -1,
		UserId:          -1,
		UserName:        DEFOULTVALUE,
		MessageType:     -1,
		Command:         DEFOULTVALUE,
		HistoryUserName: DEFOULTVALUE,
		Balance:         -1,
		BalanceAllTime:  -1,
		AutoPay:         false,
		TestUsed:        false,
		ReferalId:       -1,
		Status:          -1,
		LastMsgId:       -1,
	}
}

// Геттеры
func (u *User) GetId() int64 {
	return u.Id
}
func (u *User) GetUserStats() int {
	return u.Status
}
func (u *User) GetLastMsgId() int64 {
	return u.LastMsgId
}
func (u *User) GetUserId() int64 {
	return u.UserId
}

func (u *User) GetUserName() string {
	return u.UserName
}

func (u *User) GetMessageType() int {
	return u.MessageType
}

func (u *User) GetCommand() string {
	return u.Command
}
func (u *User) GetCommandData() string {
	return u.dataForComand
}

func (u *User) GetHistoryUserName() string {
	return u.HistoryUserName
}

func (u *User) GetBalance() int {
	return u.Balance
}

func (u *User) GetAutoPay() bool {
	return u.AutoPay
}

func (u *User) GetTestUsed() bool {
	return u.TestUsed
}

func (u *User) GetReferalId() int64 {
	return u.ReferalId
}
func (u *User) GetIgnorDeletion() bool {
	return u.ignoreDeletion
}

func (u *User) GetCallBackId() string {
	return u.callbackId
}
func (u *User) SetIgnorDeletion(ignore bool) *User {
	u.ignoreDeletion = ignore
	return u
}

func (u *User) SetCallBackId(cId string) *User {
	u.callbackId = cId
	return u
}

// Сеттеры, возвращающие указатель на саму структуру
func (u *User) SetId(id int64) *User {
	u.Id = id
	return u
}

func (u *User) SetUserId(userId int64) *User {
	u.UserId = userId
	return u
}
func (u *User) SetUserStatus(status int) *User {
	u.Status = status
	return u
}
func (u *User) SetLastMsgId(msgId int64) *User {
	u.LastMsgId = msgId
	return u
}
func (u *User) SetUserName(userName string) *User {
	u.UserName = userName
	return u
}

func (u *User) SetMessageType(messageType int) *User {
	u.MessageType = messageType
	return u
}

func (u *User) SetCommand(command string) *User {
	u.Command = command
	return u
}
func (u *User) SetCommandData(data string) *User {
	u.dataForComand = data
	return u
}

func (u *User) SetHistoryUserName(historyUserName string) *User {
	u.HistoryUserName = historyUserName
	return u
}

func (u *User) SetBalance(balance int) *User {
	u.Balance = balance
	return u
}

func (u *User) SetAutoPay(autoPay bool) *User {
	u.AutoPay = autoPay
	return u
}

func (u *User) SetTestUsed(testUsed bool) *User {
	u.TestUsed = testUsed
	return u
}

func (u *User) SetReferalId(referalId int64) *User {
	u.ReferalId = referalId
	return u
}
