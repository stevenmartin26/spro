package handler

import (
	"github.com/SawitProRecruitment/UserService/service"
)

type Server struct {
	authService    service.AuthService
	profileService service.ProfileService
}

type NewServerOptions struct {
	AuthService    service.AuthService
	ProfileService service.ProfileService
}

func NewServer(opts NewServerOptions) *Server {
	return &Server{
		authService:    opts.AuthService,
		profileService: opts.ProfileService,
	}
}
