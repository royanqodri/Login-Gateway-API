package controller

import (
	"log"
	"net/http"

	"github.com/royanqodri/Login-Gateway-API/model/notification"
	_ "github.com/royanqodri/Login-Gateway-API/model/response"
	"github.com/royanqodri/Login-Gateway-API/websocket"

	// serviceTransaction "github.com/royanqodri/Login-Gateway-API/service/transaction"

	"github.com/gin-gonic/gin"
)

type WebSocketController struct {
}

func NewWebSocketController() *WebSocketController {
	return &WebSocketController{}
}

// Subscribe godoc
// @Summary Subscribe to Notification WebSocket
// @Description Get Subscribe Notification
// @Tags Websocket
// @Produce  json
// @Param channel query string true "Channel format. Contoh nilai yang valid:\n1. customer_no:site:master-data\n2. customer_no:site:opr-dispatch\n3. customer_no:site:opr-report\n4. customer_no:site:opr-fleet\n5. customer_no:site:plant-oem\n6. customer_no:site:all"
// @Success 101 {object} notification.NotificationMessage "WebSocket connection established"
// @Failure 400 {object} response.MainResponse
// @Failure 500 {object} response.MainResponse
// @Router /mintegra/notification/sub [get]
func (ctrl *WebSocketController) Subscribe(c *gin.Context) {
	websocket.HandleWebsocketConnections(c)
}

// Publish godoc
// @Summary Publish Notification WebSocket
// @Description Get Publish Notification
// @Tags Websocket
// @Produce  json
// @Success 101 {object} notification.NotificationMessage "WebSocket connection established"
// @Failure 400 {object} response.MainResponse
// @Failure 500 {object} response.MainResponse
// @Router /mintegra/notification/pub [get]
func (ctrl *WebSocketController) Publish(c *gin.Context) {
	websocket.HandleWebsocketPublish(c)
}

// Sos godoc
// @Summary SOS Notification
// @Description Send SOS Notification to WebSocket
// @Tags Websocket
// @Accept  json
// @Produce  json
// @Param notification body notification.NotificationMessage true "Data SOS Notification"
// @Success 200 {object} map[string]string "SOS Notification sent!"
// @Failure 400 {object} map[string]string "Invalid JSON format"
// @Failure 500 {object} map[string]string "Failed to send notification"
// @Router /mintegra/notification/sos [post]
func (ctrl *WebSocketController) Sos(c *gin.Context) {
	var notification notification.NotificationMessage

	if err := c.ShouldBindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	// Send the notification and handle errors
	err := websocket.SendNotification(notification)
	if err != nil {
		log.Println("Failed to send notification:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to send notification",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "SOS Notification sent!",
	})
}

// Sos godoc
// @Summary POST Notification
// @Description Send POST Notification to WebSocket
// @Tags Websocket
// @Accept  json
// @Produce  json
// @Param notification body notification.NotificationMessage true "Data POST Notification"
// @Success 200 {object} map[string]string "POST Notification sent!"
// @Failure 400 {object} map[string]string "Invalid JSON format"
// @Failure 500 {object} map[string]string "Failed to send notification"
// @Router /mintegra/notification/post [post]
func (ctrl *WebSocketController) Post(c *gin.Context) {
	var notification notification.NotificationMessage

	if err := c.ShouldBindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	err := websocket.SendNotification(notification)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to send notification",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Notification sent!",
	})
}

// Test godoc
// @Summary Test Notification
// @Description Test Notification to WebSocket
// @Tags Websocket
// @Accept  json
// @Produce  json
// @Param notification body notification.NotificationMessage true "Data Test Notification"
// @Success 200 {object} map[string]string "Notification sent!"
// @Failure 400 {object} map[string]string "Invalid JSON format"
// @Failure 500 {object} map[string]string "Failed to send notification"
// @Router /mintegra/notification/test [post]
func (ctrl *WebSocketController) Test(c *gin.Context) {
	var notification notification.NotificationMessage

	if err := c.ShouldBindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	// Send the notification and handle errors
	err := websocket.SendNotification(notification)
	if err != nil {
		log.Println("Failed to send notification:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to send notification",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Notification sent!",
	})
}

func (ctrl *WebSocketController) GetNotifications(c *gin.Context) {
	channel := c.Query("channel")
	if channel == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "channel is required"})
		return
	}

	history, err := websocket.GetNotifications(channel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch notifications"})
		return
	}

	c.JSON(http.StatusOK, history)
}
