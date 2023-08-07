package service

import (
	"context"
	"strings"

	"github.com/SawitProRecruitment/UserService/common"
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/model"
	"github.com/SawitProRecruitment/UserService/repository"
)

type ProfileServiceImpl struct {
	userRepository repository.UserRepository
	tokenManager   TokenManager
}

func NewProfileServiceImpl(userRepository repository.UserRepository, tokenManager TokenManager) *ProfileServiceImpl {
	return &ProfileServiceImpl{
		userRepository: userRepository,
		tokenManager:   tokenManager,
	}
}

func (s *ProfileServiceImpl) GetProfile(ctx context.Context) (generated.GetProfileResponse, *common.CustomError) {
	accessToken, ok := ctx.Value(common.KeyAccessToken).(string)
	if !ok {
		return generated.GetProfileResponse{}, common.NewCustomError(common.ErrUnauthorized, "invalid access token")
	}

	userID, err := s.tokenManager.ValidateToken(accessToken)
	if err != nil {
		return generated.GetProfileResponse{}, common.NewCustomError(common.ErrUnauthorized, err.Message)
	}

	user, err := s.userRepository.GetByUserID(ctx, userID)
	if err != nil {
		return generated.GetProfileResponse{}, err
	}

	return generated.GetProfileResponse{
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
	}, nil
}

func (s *ProfileServiceImpl) UpdateProfile(ctx context.Context, params generated.UpdateProfileRequest) *common.CustomError {
	accessToken, ok := ctx.Value(common.KeyAccessToken).(string)
	if !ok {
		return common.NewCustomError(common.ErrUnauthorized, "invalid access token")
	}

	userID, err := s.tokenManager.ValidateToken(accessToken)
	if err != nil {
		return err
	}

	if err := s.validateUpdateProfileRequest(params); err != nil {
		return err
	}

	user := model.User{
		ID:          userID,
		FullName:    params.FullName,
		PhoneNumber: params.PhoneNumber,
	}

	if err := s.userRepository.Update(ctx, user); err != nil {
		return err
	}
	return nil
}

func (s *ProfileServiceImpl) validateUpdateProfileRequest(params generated.UpdateProfileRequest) *common.CustomError {
	errDetails := []string{}

	if params.PhoneNumber == "" && params.FullName == "" {
		return common.NewCustomError(common.ErrInvalidInput, "at least one phone number or full name is required")
	}

	if params.PhoneNumber != "" {
		params.PhoneNumber = strings.TrimSpace(params.PhoneNumber)
		if err := validatePhoneNumber(params.PhoneNumber); err != nil {
			errDetails = append(errDetails, err.Details...)
		}
	}

	if params.FullName != "" {
		params.FullName = strings.TrimSpace(params.FullName)
		if err := validateFullName(params.FullName); err != nil {
			errDetails = append(errDetails, err.Details...)
		}
	}

	if len(errDetails) != 0 {
		return common.NewCustomError(common.ErrInvalidInput, "invalid request params", errDetails...)
	}
	return nil
}
