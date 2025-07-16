package dto

import dto_user "github.com/amirazad1/creditor/service/dto"

type GetOtpRequest struct {
	MobileNumber string `json:"mobileNumber" binding:"required,mobile,min=11,max=11"`
}

type RegisterUserByUsernameRequest struct {
	FirstName string `json:"firstName" binding:"required,min=3"`
	LastName  string `json:"lastName" binding:"required,min=6"`
	Username  string `json:"username" binding:"required,min=5"`
	Email     string `json:"email" binding:"min=6,email"`
	Password  string `json:"password" binding:"required,password,min=6"`
}

type RegisterLoginByMobileRequest struct {
	MobileNumber string `json:"mobileNumber" binding:"required,mobile,min=11,max=11"`
	Otp          string `json:"otp" binding:"required,min=6,max=6"`
}

type LoginByUsernameRequest struct {
	Username string `json:"username" binding:"required,min=5"`
	Password string `json:"password" binding:"required,min=6"`
}

func (from RegisterUserByUsernameRequest) ToRegisterUserByUsername() dto_user.RegisterUserByUsername {
	return dto_user.RegisterUserByUsername{
		Username:  from.Username,
		FirstName: from.FirstName,
		LastName:  from.LastName,
		Email:     from.Email,
		Password:  from.Password,
	}
}
