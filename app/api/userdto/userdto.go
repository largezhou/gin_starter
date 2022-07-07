package userdto

type CreateDto struct {
	Username        string  `binding:"required"`
	Password        string  `binding:"required"`
	ConfirmPassword string  `binding:"required,eqfield=Password"`
	Email           *string `binding:"omitempty,email"`
}
