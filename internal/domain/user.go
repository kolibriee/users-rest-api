package domain

type User struct {
    ID       int64  `json:"id" db:"id"`
    Name     string `json:"name" db:"name"`
    Username string `json:"username" db:"username"`
    Email    string `json:"email" db:"email"`
    Password string `json:"password,omitempty" db:"password_hash"`
    City     string `json:"city" db:"city"`
    Role     string `json:"role" db:"role"`
}

type UpdateUserInput struct {
    Name     *string `json:"name"`
    Username *string `json:"username"`
    Email    *string `json:"email"`
    Password *string `json:"password"`
    City     *string `json:"city"`
    Role     *string `json:"role"`
}

type SignUpInput struct {
    Name     string `json:"name" binding:"required"`
    Username string `json:"username" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
    City     string `json:"city" binding:"required"`
}

type SignInInput struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required,min=6"`
}
