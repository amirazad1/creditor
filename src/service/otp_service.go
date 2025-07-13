package service

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/amirazad1/creditor/common"
	"github.com/amirazad1/creditor/config"
	constant "github.com/amirazad1/creditor/constant"
	"github.com/amirazad1/creditor/infra/cache"
	"github.com/amirazad1/creditor/pkg/logging"
	"github.com/amirazad1/creditor/pkg/service_errors"
)

type OtpService struct {
	logger      logging.Logger
	cfg         *config.Config
	redisClient *redis.Client
}

type otpDto struct {
	Value string
	Used  bool
}

func NewOtpService(cfg *config.Config) *OtpService {
	logger := logging.NewLogger(cfg)
	redis := cache.GetRedis()
	return &OtpService{logger: logger, cfg: cfg, redisClient: redis}
}

func (u *OtpService) SendOtp(mobileNumber string) error {
	otp := common.GenerateOtp()
	err := u.SetOtp(mobileNumber, otp)
	if err != nil {
		return err
	}
	return nil
}

func (u *OtpService) SetOtp(mobileNumber string, otp string) error {
	key := fmt.Sprintf("%s:%s", constant.RedisOtpDefaultKey, mobileNumber)
	val := &otpDto{
		Value: otp,
		Used:  false,
	}

	res, err := cache.Get[otpDto](u.redisClient, key)
	if err == nil && !res.Used {
		return &service_errors.ServiceError{EndUserMessage: service_errors.OtpExists}
	} else if err == nil && res.Used {
		return &service_errors.ServiceError{EndUserMessage: service_errors.OtpUsed}
	}
	err = cache.Set(u.redisClient, key, val, u.cfg.Otp.ExpireTime*time.Second)
	if err != nil {
		return err
	}
	return nil
}

func (u *OtpService) ValidateOtp(mobileNumber string, otp string) error {
	key := fmt.Sprintf("%s:%s", constant.RedisOtpDefaultKey, mobileNumber)
	res, err := cache.Get[otpDto](u.redisClient, key)
	if err != nil {
		return err
	} else if res.Used {
		return &service_errors.ServiceError{EndUserMessage: service_errors.OtpUsed}
	} else if !res.Used && res.Value != otp {
		return &service_errors.ServiceError{EndUserMessage: service_errors.OtpNotValid}
	} else if !res.Used && res.Value == otp {
		res.Used = true
		err = cache.Set(u.redisClient, key, res, u.cfg.Otp.ExpireTime*time.Second)
		if err != nil {
			return err
		}
	}
	return nil
}
