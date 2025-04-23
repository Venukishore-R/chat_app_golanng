package endpoints

import (
	"github.com/Venukishore-R/chat_app/handlers"
	"github.com/gin-gonic/gin"
)

func Endpoints(r *gin.Engine, reqHandlers handlers.RequestHandlers) {
	r.POST("/join", reqHandlers.Join)
	r.DELETE("/leave", reqHandlers.Leave)
	r.POST("/send", reqHandlers.SendMessage)
	r.POST("/messages", reqHandlers.Messages)
	
	r.GET("/health", reqHandlers.Health)
}
