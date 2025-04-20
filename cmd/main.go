package main

import (
	"log"

	"github.com/Venukishore-R/chat_app/endpoints"
	"github.com/Venukishore-R/chat_app/handlers"
	"github.com/Venukishore-R/chat_app/internal/app/services/server"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	ss := server.NewServer()
	handlers := handlers.NewHandler(ss)
	endpoints.Endpoints(r, handlers)

	go ss.Run()
	
	if err := r.Run(":8000"); err != nil {
		log.Fatalf("error starting server at port 8000")
		return
	}
}
