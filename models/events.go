package models

type Event struct {
	ID    int    `json:"id" db:"id"`
	Title string `json:"title" db:"title"`
}

type Group struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type User struct {
	ID       int    `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
	Email    string `json:"email" db:"email"`
}

type Category struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}
