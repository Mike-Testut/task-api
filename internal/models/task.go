package models

type Task struct {
	ID        int    `json:"id" db:"id"`
	UserID    int    `json:"-" db:"user_id"`
	Content   string `json:"content" db:"content"`
	Completed bool   `json:"completed" db:"completed"`
	CreatedAt string `json:"-" db:"created_at"`
}
