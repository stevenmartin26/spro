package service

import (
	"context"

	"github.com/SawitProRecruitment/UserService/common"
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/google/uuid"
)

type AuthService interface {
	Register(ctx context.Context, params generated.RegisterRequest) (generated.RegisterResponse, *common.CustomError)
	Login(ctx context.Context, params generated.LoginRequest) (generated.LoginResponse, *common.CustomError)
}

type ProfileService interface {
	GetProfile(ctx context.Context) (generated.GetProfileResponse, *common.CustomError)
	UpdateProfile(ctx context.Context, params generated.UpdateProfileRequest) *common.CustomError
}

type TokenManager interface {
	GenerateToken(userID uuid.UUID) (string, *common.CustomError)
	ValidateToken(accessToken string) (uuid.UUID, *common.CustomError)
}
