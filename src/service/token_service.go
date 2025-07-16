package service

import (
	"time"

	"github.com/amirazad1/creditor/config"
	constants "github.com/amirazad1/creditor/constant"
	"github.com/amirazad1/creditor/pkg/logging"
	"github.com/amirazad1/creditor/pkg/service_errors"
	"github.com/amirazad1/creditor/service/dto"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type TokenService struct {
	logger logging.Logger
	cfg    *config.Config
}

type tokenDto struct {
	UserId       int
	FirstName    string
	LastName     string
	Position     string	
	Roles        []string
}

func NewTokenService(cfg *config.Config) *TokenService {
	logger := logging.NewLogger(cfg)
	return &TokenService{
		cfg:    cfg,
		logger: logger,
	}
}

func (u *TokenService) GenerateToken(token tokenDto) (*dto.TokenDetail, error) {
	td := &dto.TokenDetail{}
	td.AccessTokenExpireTime = time.Now().Add(u.cfg.JWT.AccessTokenExpireDuration * time.Minute).Unix()
	td.RefreshTokenExpireTime = time.Now().Add(u.cfg.JWT.RefreshTokenExpireDuration * time.Minute).Unix()

	atc := jwt.MapClaims{}

	atc[constants.UserIdKey] = token.UserId
	atc[constants.FirstNameKey] = token.FirstName
	atc[constants.LastNameKey] = token.LastName	
	atc[constants.RolesKey] = token.Roles
	atc[constants.ExpireTimeKey] = td.AccessTokenExpireTime

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atc)

	var err error
	td.AccessToken, err = at.SignedString([]byte(u.cfg.JWT.Secret))

	if err != nil {
		return nil, err
	}

	rtc := jwt.MapClaims{}

	rtc[constants.UserIdKey] = token.UserId
	rtc[constants.FirstNameKey] = token.FirstName
	rtc[constants.LastNameKey] = token.LastName	
	rtc[constants.RolesKey] = token.Roles
	rtc[constants.ExpireTimeKey] = td.RefreshTokenExpireTime

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtc)

	td.RefreshToken, err = rt.SignedString([]byte(u.cfg.JWT.RefreshSecret))

	if err != nil {
		return nil, err
	}

	return td, nil
}

func (u *TokenService) VerifyToken(token string) (*jwt.Token, error) {
	at, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, &service_errors.ServiceError{EndUserMessage: service_errors.UnExpectedError}
		}
		return []byte(u.cfg.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	return at, nil
}

func (u *TokenService) GetClaims(token string) (claimMap map[string]interface{}, err error) {
	claimMap = map[string]interface{}{}

	verifyToken, err := u.VerifyToken(token)
	if err != nil {
		return nil, err
	}
	claims, ok := verifyToken.Claims.(jwt.MapClaims)
	if ok && verifyToken.Valid {
		for k, v := range claims {
			claimMap[k] = v
		}
		return claimMap, nil
	}
	return nil, &service_errors.ServiceError{EndUserMessage: service_errors.ClaimsNotFound}
}

func (s *TokenService) RefreshToken(c *gin.Context) (*dto.TokenDetail, error) {
	refreshToken, err := c.Cookie(constants.RefreshTokenCookieName)
	if err != nil {
		return nil, &service_errors.ServiceError{EndUserMessage: service_errors.InvalidRefreshToken}
	}

	claims, err := s.GetClaims(refreshToken)
	if err != nil {
		return nil, err
	}

	// Convert roles to []string
	rolesInterface, ok := claims[constants.RolesKey].([]interface{})
	if !ok {
		return nil, &service_errors.ServiceError{EndUserMessage: service_errors.InvalidRolesFormat}
	}

	roles := make([]string, len(rolesInterface))
	for i, role := range rolesInterface {
		roles[i], ok = role.(string)
		if !ok {
			return nil, &service_errors.ServiceError{EndUserMessage: service_errors.InvalidRolesFormat}
		}
	}

	tokenDto := tokenDto{
		UserId:       int(claims[constants.UserIdKey].(float64)),
		FirstName:    claims[constants.FirstNameKey].(string),
		LastName:     claims[constants.LastNameKey].(string),		
		Roles:        roles,
	}
	newTokenDetail, err := s.GenerateToken(tokenDto)
	if err != nil {
		return nil, err
	}

	return newTokenDetail, nil
}
