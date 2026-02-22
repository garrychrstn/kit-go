package requests

type IRequestLogin struct {
	UsernameOrEmail string `json:"username_or_email" binding:"required"`
	Password        string `json:"password" binding:"required"`
}

type IRequestRegister struct {
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
	PhoneNumber     string `json:"phone_number" binding:"required"`
	Name            string `json:"name" binding:"required"`
	Username        string `json:"username" binding:"required"`
}

type IPartUser struct {
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
	Name            string `json:"name" binding:"required"`
}
