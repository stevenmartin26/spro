package handler_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SawitProRecruitment/UserService/common"
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/service"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type HTTPHandlerTestSuite struct {
	suite.Suite
	ctrl           *gomock.Controller
	authService    *service.MockAuthService
	profileService *service.MockProfileService
	sut            *handler.Server
}

func (s *HTTPHandlerTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.authService = service.NewMockAuthService(s.ctrl)
	s.profileService = service.NewMockProfileService(s.ctrl)
	s.sut = handler.NewServer(handler.NewServerOptions{
		AuthService:    s.authService,
		ProfileService: s.profileService,
	})
}

func (s *HTTPHandlerTestSuite) AfterTest(suiteName, testName string) {
	s.ctrl.Finish()
}

func TestHTTPHandlerl(t *testing.T) {
	suite.Run(t, new(HTTPHandlerTestSuite))
}

func (s *HTTPHandlerTestSuite) TestPutV1UsersProfileGivenInvalidAuthorizationParamShouldReturnForbidden() {
	e := echo.New()
	r := httptest.NewRequest(http.MethodPut, "/api/v1/users/profile", nil)
	w := httptest.NewRecorder()
	ctx := e.NewContext(r, w)

	params := generated.PutV1UsersProfileParams{
		Authorization: "auth",
	}

	s.sut.PutV1UsersProfile(ctx, params)

	s.Equal(http.StatusForbidden, w.Result().StatusCode)
}

func (s *HTTPHandlerTestSuite) TestPutV1UsersProfileGivenInvalidRequestShouldReturnBadRequest() {
	e := echo.New()
	r := httptest.NewRequest(http.MethodPut, "/api/v1/users/profile", nil)
	w := httptest.NewRecorder()
	ctx := e.NewContext(r, w)

	params := generated.PutV1UsersProfileParams{
		Authorization: "Bearer token",
	}

	s.sut.PutV1UsersProfile(ctx, params)

	s.Equal(http.StatusBadRequest, w.Result().StatusCode)
}

func (s *HTTPHandlerTestSuite) TestPutV1UsersProfileOnInvalidInputErrorShouldReturnBadRequest() {
	request := `
		{
			"phone_number": "+62888888888",
			"full_name": "Jasuke"
		}
	`

	e := echo.New()
	r := httptest.NewRequest(http.MethodPut, "/api/v1/users/profile", bytes.NewReader([]byte(request)))
	w := httptest.NewRecorder()
	ctx := e.NewContext(r, w)

	accessToken := "token"
	params := generated.PutV1UsersProfileParams{
		Authorization: "Bearer " + accessToken,
	}

	getProfileErr := common.NewCustomError(common.ErrInvalidInput, "bad input")

	expectedAppCtx := context.WithValue(r.Context(), common.KeyAccessToken, accessToken)
	expectedRequest := generated.UpdateProfileRequest{
		FullName:    "Jasuke",
		PhoneNumber: "+62888888888",
	}

	s.profileService.EXPECT().UpdateProfile(gomock.Eq(expectedAppCtx), gomock.Eq(expectedRequest)).Return(getProfileErr)

	s.sut.PutV1UsersProfile(ctx, params)

	s.Equal(http.StatusBadRequest, w.Result().StatusCode)
}

func (s *HTTPHandlerTestSuite) TestPutV1UsersProfileOnUnauthorizedErrorShouldReturnForbidden() {
	request := `
		{
			"phone_number": "+62888888888",
			"full_name": "Jasuke"
		}
	`

	e := echo.New()
	r := httptest.NewRequest(http.MethodPut, "/api/v1/users/profile", bytes.NewReader([]byte(request)))
	w := httptest.NewRecorder()
	ctx := e.NewContext(r, w)

	accessToken := "token"
	params := generated.PutV1UsersProfileParams{
		Authorization: "Bearer " + accessToken,
	}

	getProfileErr := common.NewCustomError(common.ErrUnauthorized, "bad access token")

	expectedAppCtx := context.WithValue(r.Context(), common.KeyAccessToken, accessToken)
	expectedRequest := generated.UpdateProfileRequest{
		FullName:    "Jasuke",
		PhoneNumber: "+62888888888",
	}

	s.profileService.EXPECT().UpdateProfile(gomock.Eq(expectedAppCtx), gomock.Eq(expectedRequest)).Return(getProfileErr)

	s.sut.PutV1UsersProfile(ctx, params)

	s.Equal(http.StatusForbidden, w.Result().StatusCode)
}

func (s *HTTPHandlerTestSuite) TestPutV1UsersProfileOnEntityAlreadyExistsErrorShouldReturnConflict() {
	request := `
		{
			"phone_number": "+62888888888",
			"full_name": "Jasuke"
		}
	`

	e := echo.New()
	r := httptest.NewRequest(http.MethodPut, "/api/v1/users/profile", bytes.NewReader([]byte(request)))
	w := httptest.NewRecorder()
	ctx := e.NewContext(r, w)

	accessToken := "token"
	params := generated.PutV1UsersProfileParams{
		Authorization: "Bearer " + accessToken,
	}

	getProfileErr := common.NewCustomError(common.ErrEntityAlreadyExists, "phone number is used")

	expectedAppCtx := context.WithValue(r.Context(), common.KeyAccessToken, accessToken)
	expectedRequest := generated.UpdateProfileRequest{
		FullName:    "Jasuke",
		PhoneNumber: "+62888888888",
	}

	s.profileService.EXPECT().UpdateProfile(gomock.Eq(expectedAppCtx), gomock.Eq(expectedRequest)).Return(getProfileErr)

	s.sut.PutV1UsersProfile(ctx, params)

	s.Equal(http.StatusConflict, w.Result().StatusCode)
}

func (s *HTTPHandlerTestSuite) TestPutV1UsersProfileOnUnexpectedErrorShouldReturnInternalServerError() {
	request := `
		{
			"phone_number": "+62888888888",
			"full_name": "Jasuke"
		}
	`

	e := echo.New()
	r := httptest.NewRequest(http.MethodPut, "/api/v1/users/profile", bytes.NewReader([]byte(request)))
	w := httptest.NewRecorder()
	ctx := e.NewContext(r, w)

	accessToken := "token"
	params := generated.PutV1UsersProfileParams{
		Authorization: "Bearer " + accessToken,
	}

	getProfileErr := common.NewCustomError(common.ErrUnexpectedError, "database error")

	expectedAppCtx := context.WithValue(r.Context(), common.KeyAccessToken, accessToken)
	expectedRequest := generated.UpdateProfileRequest{
		FullName:    "Jasuke",
		PhoneNumber: "+62888888888",
	}

	s.profileService.EXPECT().UpdateProfile(gomock.Eq(expectedAppCtx), gomock.Eq(expectedRequest)).Return(getProfileErr)

	s.sut.PutV1UsersProfile(ctx, params)

	s.Equal(http.StatusInternalServerError, w.Result().StatusCode)
}

func (s *HTTPHandlerTestSuite) TestPutV1UsersProfileWhenUpdateSuccessShouldReturnOK() {
	request := `
		{
			"phone_number": "+62888888888",
			"full_name": "Jasuke"
		}
	`

	e := echo.New()
	r := httptest.NewRequest(http.MethodPut, "/api/v1/users/profile", bytes.NewReader([]byte(request)))
	w := httptest.NewRecorder()
	ctx := e.NewContext(r, w)

	accessToken := "token"
	params := generated.PutV1UsersProfileParams{
		Authorization: "Bearer " + accessToken,
	}

	expectedAppCtx := context.WithValue(r.Context(), common.KeyAccessToken, accessToken)
	expectedRequest := generated.UpdateProfileRequest{
		FullName:    "Jasuke",
		PhoneNumber: "+62888888888",
	}

	s.profileService.EXPECT().UpdateProfile(gomock.Eq(expectedAppCtx), gomock.Eq(expectedRequest)).Return(nil)

	s.sut.PutV1UsersProfile(ctx, params)

	s.Equal(http.StatusOK, w.Result().StatusCode)
}
