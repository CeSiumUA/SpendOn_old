package models

import "time"

type Login struct {
	UserName string
	Password string
}

type LoginResult struct {
	Token      string
	ExpireDate time.Time
}

type DbLogin struct {
	Id    string
	Login string
}
