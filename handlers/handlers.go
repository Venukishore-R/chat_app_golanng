package handlers

import (
	"fmt"
	"io"
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
	Messages(ctx *gin.Context)

	Health(ctx *gin.Context)
}

func (h *Handler) Health(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "OK...",
	})
}

type Request struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

func (h *Handler) Join(ctx *gin.Context) {
	var req *Request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "id field missing in query",
			"error":   "Bad request",
		})
		return
	}

	// id := ctx.Query("id")
	if req.Id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "id field missing in query",
			"error":   "Bad request",
		})
		return
	}

	client := h.ServerService.JoinRoom(req.Id)
	if client == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Failed to join room, User already joined the room",
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
	var req *Request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "id field missing in query",
			"error":   "Bad request",
		})
		return
	}

	// id := ctx.Query("id")

	if req.Id == "" || req.Message == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "id field missing in query",
			"error":   "Bad request",
		})
		return
	}

	if _, exists := h.ServerService.Client[req.Id]; !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Failed to leave room",
			"error":   fmt.Errorf("user not yet joined chat room"),
		})
	}

	if err := h.ServerService.LeaveRoom(req.Id); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": err.Error(),
			"error":   "Client not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": fmt.Sprintf("Client %s left", req.Id),
	})
}

func (h *Handler) SendMessage(ctx *gin.Context) {
	var req *Request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "id field missing in query",
			"error":   "Bad request",
		})
		return
	}

	// id := ctx.Query("id")
	// msg := ctx.Query("msg")

	id := req.Id
	msg := req.Message

	if id == "" || msg == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "id/msg field missing in query",
			"error":   "Bad request",
		})
		return
	}

	if _, exists := h.ServerService.Client[id]; !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Failed to send message",
			"error":   fmt.Errorf("user not yet joined chat room"),
		})
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

	ctx.Writer.Header().Set("Content-type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	ctx.Writer.Flush()

	c := ctx.Request.Context()

	// var messages []string
	// timeout := time.After(10 * time.Second)

	for {
		// select {
		// case msg := <-msgClient.MsgChan:
		// 	messages = append(messages, msg)
		// 	if len(messages) >= 10 {
		// 		ctx.JSON(200, gin.H{"messages": messages})
		// 		return
		// 	}
		// case <-timeout:
		// 	if len(messages) == 0 {
		// 		ctx.JSON(204, gin.H{"message": "No new messages"})
		// 	} else {
		// 		ctx.JSON(200, gin.H{"messages": messages})
		// 	}
		// 	return
		// case <-ctx.Done():
		// 	ctx.JSON(499, gin.H{"error": "client disconnected"})
		// 	return

		select {
		case msg, ok := <-msgClient.MsgChan:
			if !ok {
				return
			}
			fmt.Fprintf(ctx.Writer, "data: %s\n\n", msg)
			ctx.Writer.Flush()
		case <-c.Done():
			fmt.Fprintf(ctx.Writer, ": ping\n\n")
			ctx.Writer.Flush()
			return
		case <-time.After(30 * time.Second):
			fmt.Fprintf(ctx.Writer, ": ping\n\n")
			ctx.Writer.Flush()

		}
	}
}

func (h *Handler) Messages(ctx *gin.Context) {
	var req *Request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "id field missing in query",
			"error":   "Bad request",
		})
		return
	}

	// id := ctx.Query("id")
	id := req.Id

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

	// ctx.Writer.Header().Set("Content-type", "text/event-stream")
	// ctx.Writer.Header().Set("Cache-control", "no-cache")
	// ctx.Writer.Header().Set("Connection", "keep-alive")
	// ctx.Writer.Flush()

	ctx.Stream(func(w io.Writer) bool {
		// for {
		select {
		case msg, ok := <-msgClient.MsgChan:
			if !ok {
				return false
			}
			ctx.SSEvent("message", msg)
			return true
		case <-ctx.Request.Context().Done():
			return false
		}
		// }
	})
}
