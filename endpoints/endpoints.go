package endpoints

import (
	"github.com/Venukishore-R/chat_app/handlers"
	"github.com/gin-gonic/gin"
)

func Endpoints(r *gin.Engine, reqHandlers handlers.RequestHandlers) {
	r.GET("/join", reqHandlers.Join)
	r.GET("/leave", reqHandlers.Leave)
	r.GET("/send", reqHandlers.SendMessage)
	r.GET("/messages", reqHandlers.GetMessages)
	r.GET("/health", reqHandlers.Health)
}
