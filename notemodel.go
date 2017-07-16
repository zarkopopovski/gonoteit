package main

type NoteModel struct {
	Id     string `json:"id"`
	UserID string `json:"user_id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	Date   string `json:"date_created"`
}
