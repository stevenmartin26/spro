package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/SawitProRecruitment/UserService/common"
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/labstack/echo/v4"
)

func (s *Server) PostApiV1UsersRegister(ctx echo.Context) error {
	var request generated.RegisterRequest
	if err := json.NewDecoder(ctx.Request().Body).Decode(&request); err != nil {
		response := generated.ErrorResponse{
			Message: err.Error(),
		}
		return ctx.JSON(http.StatusBadRequest, response)
	}

	result, err := s.authService.Register(ctx.Request().Context(), request)
	if err != nil {
		return ctx.JSON(constructErrorResponse(err))
	}

	return ctx.JSON(http.StatusCreated, result)
}

func (s *Server) PostApiV1UsersLogin(ctx echo.Context) error {
	var request generated.LoginRequest
	if err := json.NewDecoder(ctx.Request().Body).Decode(&request); err != nil {
		response := generated.ErrorResponse{
			Message: err.Error(),
		}
		return ctx.JSON(http.StatusBadRequest, response)
	}

	result, err := s.authService.Login(ctx.Request().Context(), request)
	if err != nil {
		return ctx.JSON(constructErrorResponse(err))
	}

	return ctx.JSON(http.StatusOK, result)
}

func (s *Server) GetV1UsersProfile(ctx echo.Context, params generated.GetV1UsersProfileParams) error {
	accessToken, errResponse := extractAccessToken(params.Authorization)
	if errResponse != nil {
		ctx.JSON(http.StatusForbidden, errResponse)
	}

	appCtx := context.WithValue(ctx.Request().Context(), common.KeyAccessToken, accessToken)

	result, err := s.profileService.GetProfile(appCtx)
	if err != nil {
		return ctx.JSON(constructErrorResponse(err))
	}
	return ctx.JSON(http.StatusOK, result)
}

func (s *Server) PutV1UsersProfile(ctx echo.Context, params generated.PutV1UsersProfileParams) error {
	accessToken, errResponse := extractAccessToken(params.Authorization)
	if errResponse != nil {
		ctx.JSON(http.StatusForbidden, errResponse)
	}

	var request generated.UpdateProfileRequest
	if err := json.NewDecoder(ctx.Request().Body).Decode(&request); err != nil {
		response := generated.ErrorResponse{
			Message: err.Error(),
		}
		return ctx.JSON(http.StatusBadRequest, response)
	}

	appCtx := context.WithValue(ctx.Request().Context(), common.KeyAccessToken, accessToken)

	err := s.profileService.UpdateProfile(appCtx, request)
	if err != nil {
		return ctx.JSON(constructErrorResponse(err))
	}
	return ctx.JSON(http.StatusOK, nil)
}
