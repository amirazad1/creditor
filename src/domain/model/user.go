package model

type User struct {
	BaseModel
	Username     string `gorm:"type:string;size:50;not null;unique"`
	FirstName    string `gorm:"type:string;size:20;null"`
	LastName     string `gorm:"type:string;size:50;null"`
	MobileNumber string `gorm:"type:string;size:11;null;unique;default:null"`
	Email        string `gorm:"type:string;size:64;null;unique;default:null"`
	Password     string `gorm:"type:string;size:64;not null"`
	Position     string `gorm:"null"`
	Description  string `gorm:"null"`
	UserRoles    *[]UserRole
}

type Role struct {
	BaseModel
	Name        string `gorm:"type:string;size:20;not null,unique"`
	Description string `gorm:"not null"`
	UserRoles   *[]UserRole
}

type UserRole struct {
	BaseModel
	User   User `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
	Role   Role `gorm:"foreignKey:RoleId;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
	UserId int
	RoleId int
}
