package models

type Task struct {
	ID        int    `json:"id" db:"id"`
	Content   string `json:"content" db:"content"`
	Completed bool   `json:"completed" db:"completed"`
	CreatedAt string `json:"-" db:"created_at"`
}
