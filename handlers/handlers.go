package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Venukishore-R/chat_app/internal/app/services/server"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	ServerService *server.Server
}

func NewHandler(serverService *server.Server) *Handler {
	return &Handler{
		ServerService: serverService,
	}
}

type RequestHandlers interface {
	Join(ctx *gin.Context)
	Leave(ctx *gin.Context)
	SendMessage(ctx *gin.Context)
	GetMessages(ctx *gin.Context)

	Health(ctx *gin.Context)
}

func (h *Handler) Health(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "OK...",
	})
}

func (h *Handler) Join(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "id field missing in query",
			"error":   "Bad request",
		})
		return
	}

	client := h.ServerService.JoinRoom(id)
	if client == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Failed to join room",
			"error":   "Internal error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": fmt.Sprintf("Client %s joined", client.Id),
	})
}

// Leave removes a client from the server
func (h *Handler) Leave(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "id field missing in query",
			"error":   "Bad request",
		})
		return
	}

	if err := h.ServerService.LeaveRoom(id); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": err.Error(),
			"error":   "Client not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": fmt.Sprintf("Client %s left", id),
	})
}

func (h *Handler) SendMessage(ctx *gin.Context) {
	id := ctx.Query("id")
	msg := ctx.Query("msg")

	if id == "" || msg == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "id/msg field missing in query",
			"error":   "Bad request",
		})
		return
	}

	if err := h.ServerService.SendMessage(id, msg); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Failed to send message",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Message sent successfully",
	})
}

func (h *Handler) GetMessages(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		ctx.JSON(400, gin.H{"error": "id field missing in query"})
		return
	}

	h.ServerService.Mu.Lock()
	msgClient, exists := h.ServerService.Client[id]
	h.ServerService.Mu.Unlock()

	if !exists {
		ctx.JSON(404, gin.H{"error": "client does not exist"})
		return
	}

	var messages []string
	timeout := time.After(10 * time.Second)

	for {
		select {
		case msg := <-msgClient.MsgChan:
			messages = append(messages, msg)
			if len(messages) >= 10 {
				ctx.JSON(200, gin.H{"messages": messages})
				return
			}
		case <-timeout:
			if len(messages) == 0 {
				ctx.JSON(204, gin.H{"message": "No new messages"})
			} else {
				ctx.JSON(200, gin.H{"messages": messages})
			}
			return
		case <-ctx.Done():
			ctx.JSON(499, gin.H{"error": "client disconnected"})
			return
		}
	}
}
