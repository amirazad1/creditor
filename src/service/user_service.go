package service

import (
	"github.com/amirazad1/creditor/common"
	"github.com/amirazad1/creditor/config"
	constants "github.com/amirazad1/creditor/constant"
	"github.com/amirazad1/creditor/domain/model"
	"github.com/amirazad1/creditor/infra/persistance/database"
	"github.com/amirazad1/creditor/pkg/logging"
	"github.com/amirazad1/creditor/pkg/service_errors"
	dto "github.com/amirazad1/creditor/service/dto"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const userFilterExp string = "username = ?"
const countFilterExp string = "count(*) > 0"

type UserService struct {
	logger       logging.Logger
	cfg          *config.Config
	otpService   *OtpService
	tokenService *TokenService
	database     *gorm.DB
}

func NewUserService(cfg *config.Config) *UserService {
	logger := logging.NewLogger(cfg)
	return &UserService{
		cfg:          cfg,
		logger:       logger,
		otpService:   NewOtpService(cfg),
		tokenService: NewTokenService(cfg),
		database:     database.GetDb(),
	}
}

// Login by username
func (u *UserService) LoginByUsername(username string, password string) (*dto.TokenDetail, error) {
	user, err := u.FetchUserInfo(username, password)

	if err != nil {
		return nil, err
	}
	tokenDto := tokenDto{UserId: user.Id, FirstName: user.FirstName, LastName: user.LastName}

	if len(*user.UserRoles) > 0 {
		for _, ur := range *user.UserRoles {
			tokenDto.Roles = append(tokenDto.Roles, ur.Role.Name)
		}
	}

	token, err := u.tokenService.GenerateToken(tokenDto)

	if err != nil {
		return nil, err
	}
	return token, nil

}

func (r *UserService) CreateUser(u model.User) (model.User, error) {

	roleId, err := r.GetDefaultRole()
	if err != nil {
		r.logger.Error(logging.Postgres, logging.DefaultRoleNotFound, err.Error(), nil)
		return u, err
	}
	tx := r.database.Begin()
	err = tx.Create(&u).Error
	if err != nil {
		tx.Rollback()
		r.logger.Error(logging.Postgres, logging.Rollback, err.Error(), nil)
		return u, err
	}
	err = tx.Create(&model.UserRole{RoleId: roleId, UserId: u.Id}).Error
	if err != nil {
		tx.Rollback()
		r.logger.Error(logging.Postgres, logging.Rollback, err.Error(), nil)
		return u, err
	}
	tx.Commit()
	return u, nil
}

func (u *UserService) FetchUserInfo(username string, password string) (model.User, error) {
	var user model.User
	err := u.database.
		Model(&model.User{}).
		Where(userFilterExp, username).
		Preload("UserRoles", func(tx *gorm.DB) *gorm.DB {
			return tx.Preload("Role")
		}).
		Find(&user).Error

	if err != nil {
		return user, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return user, err
	}

	return user, nil
}

func (u *UserService) FetchUserInfoForMobile(username string) (model.User, error) {
	var user model.User
	err := u.database.
		Model(&model.User{}).
		Where(userFilterExp, username).
		Preload("UserRoles", func(tx *gorm.DB) *gorm.DB {
			return tx.Preload("Role")
		}).
		Find(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

// Register by username
func (u *UserService) RegisterByUsername(req dto.RegisterUserByUsername) error {
	user := dto.ToUserModel(req)

	exists, err := u.ExistsEmail(req.Email)
	if err != nil {
		return err
	}
	if exists {
		return &service_errors.ServiceError{EndUserMessage: service_errors.EmailExists}
	}
	exists, err = u.ExistsUsername(req.Username)
	if err != nil {
		return err
	}
	if exists {
		return &service_errors.ServiceError{EndUserMessage: service_errors.UsernameExists}
	}

	bp := []byte(req.Password)
	hp, err := bcrypt.GenerateFromPassword(bp, bcrypt.DefaultCost)
	if err != nil {
		u.logger.Error(logging.General, logging.HashPassword, err.Error(), nil)
		return err
	}
	user.Password = string(hp)
	_, err = u.CreateUser(user)
	return err

}

// Register/login by mobile number
func (u *UserService) RegisterAndLoginByMobileNumber(mobileNumber string, otp string) (*dto.TokenDetail, error) {
	err := u.otpService.ValidateOtp(mobileNumber, otp)
	if err != nil {
		return nil, err
	}
	exists, err := u.ExistsMobileNumber(mobileNumber)
	if err != nil {
		return nil, err
	}

	user := model.User{MobileNumber: mobileNumber, Username: mobileNumber}

	if exists {
		user, err = u.FetchUserInfoForMobile(user.Username)
		if err != nil {
			return nil, err
		}

		return u.generateToken(user)
	}

	// Register and login
	bp := []byte(common.GeneratePassword())
	hp, err := bcrypt.GenerateFromPassword(bp, bcrypt.DefaultCost)
	if err != nil {
		u.logger.Error(logging.General, logging.HashPassword, err.Error(), nil)
		return nil, err
	}
	user.Password = string(hp)

	user, err = u.CreateUser(user)
	if err != nil {
		return nil, err
	}
	return u.generateToken(user)
}

func (u *UserService) generateToken(user model.User) (*dto.TokenDetail, error) {
	tokenDto := tokenDto{UserId: user.Id, FirstName: user.FirstName, LastName: user.LastName}

	if len(*user.UserRoles) > 0 {
		for _, ur := range *user.UserRoles {
			tokenDto.Roles = append(tokenDto.Roles, ur.Role.Name)
		}
	}

	token, err := u.tokenService.GenerateToken(tokenDto)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (r *UserService) ExistsEmail(email string) (bool, error) {
	var exists bool
	if err := r.database.Model(&model.User{}).
		Select(countFilterExp).
		Where("email = ?", email).
		Find(&exists).
		Error; err != nil {
		r.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		return false, err
	}
	return exists, nil
}

func (r *UserService) ExistsUsername(username string) (bool, error) {
	var exists bool
	if err := r.database.Model(&model.User{}).
		Select(countFilterExp).
		Where(userFilterExp, username).
		Find(&exists).
		Error; err != nil {
		r.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		return false, err
	}
	return exists, nil
}

func (r *UserService) ExistsMobileNumber(mobileNumber string) (bool, error) {
	var exists bool
	if err := r.database.Model(&model.User{}).
		Select(countFilterExp).
		Where("mobile_number = ?", mobileNumber).
		Find(&exists).
		Error; err != nil {
		r.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		return false, err
	}
	return exists, nil
}

func (r *UserService) GetDefaultRole() (roleId int, err error) {

	if err = r.database.Model(&model.Role{}).
		Select("id").
		Where("name = ?", constants.DefaultRoleName).
		First(&roleId).Error; err != nil {
		return 0, err
	}
	return roleId, nil
}
