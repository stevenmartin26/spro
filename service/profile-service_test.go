package service_test

import (
	"context"
	"testing"

	"github.com/SawitProRecruitment/UserService/common"
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/model"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/service"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type ProfileServiceTestSuite struct {
	suite.Suite
	ctrl           *gomock.Controller
	tokenManager   *service.MockTokenManager
	userRepository *repository.MockUserRepository
	sut            *service.ProfileServiceImpl
}

func (s *ProfileServiceTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.tokenManager = service.NewMockTokenManager(s.ctrl)
	s.userRepository = repository.NewMockUserRepository(s.ctrl)
	s.sut = service.NewProfileServiceImpl(s.userRepository, s.tokenManager)
}

func (s *ProfileServiceTestSuite) AfterTest(suiteName, testName string) {
	s.ctrl.Finish()
}

func TestProfileServiceImpl(t *testing.T) {
	suite.Run(t, new(ProfileServiceTestSuite))
}

func (s *ProfileServiceTestSuite) TestGetProfileGivenInvalidAccessTokenShouldReturnUnauthorizedError() {
	result, err := s.sut.GetProfile(context.Background())

	s.Equal(common.ErrUnauthorized, err.ErrType)
	s.Equal(generated.GetProfileResponse{}, result)
}

func (s *ProfileServiceTestSuite) TestGetProfileOnTokenValidationErrorShouldReturnUnauthorizedError() {
	accessToken := "access token"
	ctx := context.WithValue(context.Background(), common.KeyAccessToken, accessToken)

	validateTokenErr := common.NewCustomError(common.ErrUnauthorized, "invalid token")

	s.tokenManager.EXPECT().ValidateToken(gomock.Eq(accessToken)).Return(uuid.Nil, validateTokenErr)

	result, err := s.sut.GetProfile(ctx)

	s.Equal(common.ErrUnauthorized, err.ErrType)
	s.Equal(generated.GetProfileResponse{}, result)
}

func (s *ProfileServiceTestSuite) TestGetProfileOnGetUserErrorShouldReturnError() {
	accessToken := "access token"
	userID := uuid.New()
	ctx := context.WithValue(context.Background(), common.KeyAccessToken, accessToken)

	repoErr := common.NewCustomError(common.ErrUnexpectedError, "database error")

	s.tokenManager.EXPECT().ValidateToken(gomock.Eq(accessToken)).Return(userID, nil)
	s.userRepository.EXPECT().GetByUserID(gomock.Eq(ctx), gomock.Eq(userID)).Return(nil, repoErr)

	result, err := s.sut.GetProfile(ctx)

	s.Equal(repoErr.ErrType, err.ErrType)
	s.Equal(generated.GetProfileResponse{}, result)
}

func (s *ProfileServiceTestSuite) TestGetProfileShouldReturnProfileFromRepository() {
	accessToken := "access token"
	userID := uuid.New()
	ctx := context.WithValue(context.Background(), common.KeyAccessToken, accessToken)

	user := model.User{
		ID:           userID,
		PhoneNumber:  "+628",
		FullName:     "full",
		PasswordHash: "hash",
	}

	expectedResult := generated.GetProfileResponse{
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
	}

	s.tokenManager.EXPECT().ValidateToken(gomock.Eq(accessToken)).Return(userID, nil)
	s.userRepository.EXPECT().GetByUserID(gomock.Eq(ctx), gomock.Eq(userID)).Return(&user, nil)

	result, err := s.sut.GetProfile(ctx)

	s.Nil(err)
	s.Equal(expectedResult, result)
}
