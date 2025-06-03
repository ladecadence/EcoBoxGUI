package models

type User struct {
	Result int    `json:"resultado"`
	Code   int    `json:"codigo"`
	ID     int    `json:"id"`
	Name   string `json:"nombre"`
}
