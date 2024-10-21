package models

type Code struct {
	Path     string `json:"path"`
	Lang     string `json:"lang"`
	TaskName string `json:"taskname"`
	UserName string `json:"username"`
}
