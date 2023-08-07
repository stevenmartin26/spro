package service

import (
	"context"
	"strings"
	"time"

	"github.com/SawitProRecruitment/UserService/common"
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/model"
	"github.com/SawitProRecruitment/UserService/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthServiceImpl struct {
	userRepository     repository.UserRepository
	loginLogRepository repository.LoginLogRepository
	tokenManager       TokenManager
}

func NewAuthServiceImpl(userRepository repository.UserRepository, loginLogRepository repository.LoginLogRepository, tokenManager TokenManager) *AuthServiceImpl {
	return &AuthServiceImpl{
		userRepository:     userRepository,
		loginLogRepository: loginLogRepository,
		tokenManager:       tokenManager,
	}
}

func (s *AuthServiceImpl) Register(ctx context.Context, params generated.RegisterRequest) (generated.RegisterResponse, *common.CustomError) {
	if err := s.validateRegisterRequest(params); err != nil {
		return generated.RegisterResponse{}, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return generated.RegisterResponse{}, common.NewCustomError(common.ErrUnexpectedError, err.Error())
	}

	user := model.User{
		PhoneNumber:  params.PhoneNumber,
		FullName:     params.FullName,
		PasswordHash: string(passwordHash),
	}

	userID, errSave := s.userRepository.Save(ctx, user)
	if errSave != nil {
		return generated.RegisterResponse{}, errSave
	}
	return generated.RegisterResponse{UserId: userID}, nil
}

func (s *AuthServiceImpl) validateRegisterRequest(params generated.RegisterRequest) *common.CustomError {
	errDetails := []string{}

	params.FullName = strings.TrimSpace(params.FullName)
	if err := validateFullName(params.FullName); err != nil {
		errDetails = append(errDetails, err.Details...)
	}

	params.PhoneNumber = strings.TrimSpace(params.PhoneNumber)
	if err := validatePhoneNumber(params.PhoneNumber); err != nil {
		errDetails = append(errDetails, err.Details...)
	}

	params.Password = strings.TrimSpace(params.Password)
	if err := validatePassword(params.Password); err != nil {
		errDetails = append(errDetails, err.Details...)
	}

	if len(errDetails) != 0 {
		return common.NewCustomError(common.ErrInvalidInput, "invalid request params", errDetails...)
	}
	return nil
}

func (s *AuthServiceImpl) Login(ctx context.Context, params generated.LoginRequest) (generated.LoginResponse, *common.CustomError) {
	user, err := s.userRepository.GetByPhoneNumber(ctx, params.PhoneNumber)
	if err != nil {
		if err.ErrType == common.ErrEntityNotFound {
			return generated.LoginResponse{}, common.NewCustomError(common.ErrInvalidInput, "phone number or password is incorrect")
		}
		return generated.LoginResponse{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(params.Password)); err != nil {
		return generated.LoginResponse{}, common.NewCustomError(common.ErrInvalidInput, "phone number or password is incorrect")
	}

	tokenString, err := s.tokenManager.GenerateToken(user.ID)
	if err != nil {
		return generated.LoginResponse{}, err
	}

	loginLog := model.LoginLog{
		UserID:  user.ID,
		LoginAt: time.Now(),
	}
	go s.loginLogRepository.Save(ctx, loginLog)

	return generated.LoginResponse{UserId: user.ID, AccessToken: tokenString}, nil
}
