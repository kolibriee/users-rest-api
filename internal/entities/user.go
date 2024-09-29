package entities

type User struct {
	Id           int    `json:"-" db:"id"`
	Role         string `json:"role" db:"role"`
	Name         string `json:"name" binding:"required"`
	Username     string `json:"username" binding:"required"`
	Password     string `json:"password" binding:"required"`
	RegisteredAt string `json:"-" db:"registered_at"`
}

type SignInUserInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
