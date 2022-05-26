package handler

import (
	"github.com/aknrdlt/final-golang/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	api := router.Group("/api", h.userIdentity)
	{
		lists := api.Group("/books")
		{
			lists.POST("/", h.createBook)
			lists.GET("/", h.getAllBooks)
			lists.GET("/:id", h.getBookById)
			lists.PUT("/:id", h.updateBook)
			lists.DELETE("/:id", h.deleteBook)
		}
	}

	return router
}
