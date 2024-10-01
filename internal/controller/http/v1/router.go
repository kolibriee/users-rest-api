package v1

import (
	"net/http"

	_ "github.com/kolibriee/users-rest-api/docs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func (h *Handler) InitRouter() http.Handler {
	router := echo.New()
	router.Use(middleware.Logger())
	router.GET("/swagger/*", echoSwagger.WrapHandler)
	admin := router.Group("/admin", h.adminIdentity)
	{
		users := admin.Group("/users")
		{
			users.GET("", h.GetAllUsers)
			users.GET("/:id", h.AdminGetUserByID)
			users.POST("", h.CreateUser)
			users.PUT("/:id", h.AdminUpdateUser)
			users.DELETE("/:id", h.AdminDeleteUser)
		}
	}
	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.SignUp)
		auth.POST("/sign-in", h.SignIn)
		auth.GET("/refresh", h.Refresh)
	}
	api := router.Group("/api")
	{
		users := api.Group("/users", h.userIdentity)
		{
			users.GET("/:id", h.GetUserByID)
			users.PUT("/:id", h.UpdateUser)
			users.DELETE("/:id", h.DeleteUser)
		}
	}
	return router
}
