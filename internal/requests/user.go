package requests

import (
	"example.com/expenses-tracker/internal/util"
	"github.com/go-playground/validator/v10"
)

type LoginRequest struct {
	EmailAddress string `json:"email_address" binding:"required,email"`
	Password     string `json:"password" binding:"required"`
}

func (l LoginRequest) FormatValidationMessages(errors validator.ValidationErrors) map[string]string {
	return util.FormatValidationMessages(l, errors)
}

type CreateUserRequest struct {
	EmailAddress    string `json:"email_address" binding:"required,email"`
	FirstName       string `json:"first_name" binding:"required,alpha"`
	LastName        string `json:"last_name" binding:"required,alpha"`
	Password        string `json:"password" binding:"required,validpassword"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,validpassword"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
}
