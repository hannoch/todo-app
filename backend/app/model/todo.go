package model

type Todo struct {
	Id           int    `json:"id" orm:"id,primary"`
	Title        string `json:"title" orm:"title"`
	Description  string `json:"description" orm:"description"`
	CategoryId   int    `json:"category_id" orm:"category_id"`
	DocumentLink string `json:"document_link" orm:"document_link"`
	DueDate      string `json:"due_date" orm:"due_date"`
	ReminderTime string `json:"reminder_time" orm:"reminder_time"`
	IsImportant  bool   `json:"is_important" orm:"is_important"`
	IsDaily      bool   `json:"is_daily" orm:"is_daily"`
	IsCompleted  bool   `json:"is_completed" orm:"is_completed"`
}
