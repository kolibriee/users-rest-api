package controller

import (
	"net/http"

	v1 "github.com/kolibriee/users-rest-api/internal/controller/v1"
	"github.com/kolibriee/users-rest-api/internal/service"
)

type RouterInitializer interface {
	InitRouter() http.Handler
}

type Controller struct {
	Handler RouterInitializer
}

func NewController(services *service.Service) *Controller {
	return &Controller{
		Handler: v1.NewHandler(services),
	}
}
