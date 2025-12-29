package types

type IRequestLogin struct {
	UsernameOrEmail string `json:"username_or_email" binding:"required"`
	Password        string `json:"password" binding:"required"`
	Age             int    `json:"age" binding:"required,gt=10,max=150"`
}

type IPartUser struct {
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
	Name            string `json:"name" binding:"required"`
}
