package models

type Article struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Content  string `json:"content"`
}
