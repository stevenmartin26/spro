package repository

import (
	"context"

	"github.com/SawitProRecruitment/UserService/common"
	"github.com/SawitProRecruitment/UserService/model"
	"github.com/google/uuid"
)

type UserRepository interface {
	Save(ctx context.Context, user model.User) (uuid.UUID, *common.CustomError)
	GetByPhoneNumber(ctx context.Context, phoneNumber string) (*model.User, *common.CustomError)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*model.User, *common.CustomError)
	Update(ctx context.Context, user model.User) *common.CustomError
}

type LoginLogRepository interface {
	Save(ctx context.Context, log model.LoginLog) (uuid.UUID, *common.CustomError)
}
