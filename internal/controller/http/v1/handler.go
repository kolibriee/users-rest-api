package v1

import "github.com/kolibriee/users-rest-api/internal/service"

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}
