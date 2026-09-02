package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"github.com/royanqodri/Login-Gateway-API/config"
	entityTransaction "github.com/royanqodri/Login-Gateway-API/model/entity/transaction"
	"github.com/royanqodri/Login-Gateway-API/model/notification"
	repositoryTransaction "github.com/royanqodri/Login-Gateway-API/repository/transaction"
	"github.com/royanqodri/Login-Gateway-API/util"
)

var (
	rdb *redis.Client
	ctx context.Context

	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	channels = make(map[string]map[*websocket.Conn]bool)
	mutex    = &sync.Mutex{}
)

func InitWebsocketRedis() {
	rdb = util.RedisClient
	ctx = util.Ctx
}

// --- Subscriber ---
func HandleWebsocketConnections(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v\n", err)
		return
	}
	defer ws.Close()

	// Determine channel from query or JSON
	channel := c.Query("channel")
	if channel == "" {
		ws.WriteMessage(websocket.TextMessage, []byte("Missing channel"))
		return
	}

	// Register subscriber
	mutex.Lock()
	if channels[channel] == nil {
		channels[channel] = make(map[*websocket.Conn]bool)
	}
	channels[channel][ws] = true
	mutex.Unlock()

	log.Printf("Client subscribed to channel: %s\n", channel)

	// Heartbeat ping
	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := ws.WriteControl(websocket.PingMessage, nil, time.Now().Add(10*time.Second)); err != nil {
					log.Printf("Ping error for channel %s: %v\n", channel, err)
					return
				}
			case <-done:
				return
			}
		}
	}()

	// Keep connection alive
	for {
		if _, _, err := ws.ReadMessage(); err != nil {
			log.Printf("Client disconnected from channel %s: %v\n", channel, err)
			break
		}
	}

	// Stop heartbeat goroutine
	close(done)

	// Cleanup
	mutex.Lock()
	if subs, ok := channels[channel]; ok {
		delete(subs, ws)
		if len(subs) == 0 {
			delete(channels, channel)
		}
	}
	mutex.Unlock()
}

// --- Publisher ---
func HandleWebsocketPublish(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v\n", err)
		return
	}

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Printf("Publisher closed connection normally")
			} else {
				log.Printf("Publisher disconnected or error: %v\n", err)
			}
			ws.Close()
			return
		}

		var publish notification.NotificationMessage
		if err := json.Unmarshal(msg, &publish); err != nil {
			log.Printf("Invalid JSON from publisher: %v\n", err)
			continue
		}

		publish.Timestamp = time.Now().Format("2006-01-02 15:04:05")
		key := fmt.Sprintf("notification:%s", publish.Channel)
		publishJSON, _ := json.Marshal(publish)

		if _, err := rdb.LPush(ctx, key, publishJSON).Result(); err != nil {
			log.Printf("Error saving notification to Redis List: %v\n", err)
		}
		_ = rdb.LTrim(ctx, key, 0, 999)
		_ = rdb.Expire(ctx, key, 48*time.Hour)

		// GET notification_tolerance_time from mst_setting default 15 min if not found
		const defaultToleranceMinutes = 15
		toleranceMinutes := defaultToleranceMinutes

		// customerRepo := repository.NewMstCustomerRepository()
		// dataCustomer, errCustomer := customerRepo.GetByCustomerNo(c, nil, publish.CustomerNo)
		// if errCustomer != nil {
		// 	log.Printf("Warning: customer '%s' not found: %v → using default tolerance %d min", publish.CustomerNo, errCustomer, defaultToleranceMinutes)
		// 	dataCustomer = entity.MstCustomer{}
		// }

		// settingRepo := repository.NewMstSettingRepository()
		// settingRequest := request.MstSettingGetRequest{
		// 	Parameter:  "notification_tolerance_time",
		// 	Parameters: []string{"notification_tolerance_time"},
		// }

		// settingMap, errSetting := settingRepo.GetByParamsINKeyValue(c, nil, dataCustomer, settingRequest)
		// if errSetting != nil {
		// 	log.Printf("Warning: failed to fetch notification_tolerance_time: %v → using default %d min", errSetting, defaultToleranceMinutes)
		// } else if val, ok := settingMap["notification_tolerance_time"]; ok && val != "" {
		// 	if parsed, errParse := strconv.Atoi(val); errParse == nil && parsed > 0 {
		// 		toleranceMinutes = parsed
		// 	} else {
		// 		log.Printf("Warning: invalid notification_tolerance_time value '%s' → using default %d min", val, defaultToleranceMinutes)
		// 	}
		// } else {
		// 	log.Printf("notification_tolerance_time not set in mst_setting → using default %d min", defaultToleranceMinutes)
		// }

		// =============================================
		// Insert into t_notifications
		// =============================================

		currentUser := c.GetString("user")
		publishUser := publish.Username
		fallbackUser := "system"
		timeLog := util.GetFormattedDateTime(publish.TimeLog)

		entity := entityTransaction.TNotifications{
			CustomerNo:     publish.CustomerNo,
			DateLog:        util.GetFormattedDate(publish.DateLog),
			TimeLog:        timeLog,
			Shift:          publish.Shift,
			ShiftSequence:  publish.ShiftSequence,
			EquipmentNo:    publish.EquipmentNo,
			EquipmentType:  publish.EquipmentType,
			EquipmentModel: publish.EquipmentModel,
			Username:       publish.Username,
			Name:           publish.Name,
			Fleet:          publish.Fleet,
			State:          publish.State,
			Reason:         publish.Reason,
			Activity:       publish.Activity,
			StatusActivity: publish.StatusActivity,
			Latitude:       publish.Latitude,
			Longitude:      publish.Longitude,
			Altitude:       publish.Altitude,
			Bearing:        publish.Bearing,
			Channel:        publish.Channel,
			Category:       publish.Category,
			Title:          publish.Title,
			Content:        publish.Content,
			Type:           publish.Type,
			Event:          publish.Event,
			Site:           publish.Site,
			InsertBy:       currentUser,
			InsertTime:     util.GetTimeNowByLoc(),
			UpdateBy:       currentUser,
			UpdateTime:     util.GetTimeNowByLoc(),
		}

		if currentUser == "" && publishUser != "" {
			entity.InsertBy = publishUser
			entity.UpdateBy = publishUser
			log.Printf("Warning: no authenticated user → using publish '%s' for notification", publishUser)
		}

		if currentUser == "" && publishUser == "" {
			entity.InsertBy = fallbackUser
			entity.UpdateBy = fallbackUser
			log.Printf("Warning: no authenticated user → using fallback '%s' for notification", fallbackUser)
		}

		repo := repositoryTransaction.NewTNotificationsRepository()

		err = repo.Save(nil, nil, []entityTransaction.TNotifications{entity})
		if err != nil {
			log.Printf("KRITIS: FAILED save to DB | customer=%s | type=%s | title=%s | err=%v", publish.CustomerNo, publish.Type, publish.Title, err)

			_ = ws.WriteMessage(websocket.TextMessage, []byte("NACK: database error"))
			continue
		}

		log.Printf("SUKSES save to DB → %s | %s | %s | channel=%s", publish.CustomerNo, publish.Type, publish.Title, publish.Channel)

		// CHECK tolerance: only broadcast if time_log is within notification_tolerance_time from now
		elapsed := time.Since(timeLog)
		tolerance := time.Duration(toleranceMinutes) * time.Minute

		if elapsed < 0 {
			log.Printf("SKIP broadcast: time_log '%s' is in the future | channel=%s", timeLog.Format("2006-01-02 15:04:05"), publish.Channel)

			_ = ws.WriteMessage(websocket.TextMessage, []byte("SKIP: time_log is in the future, saved as historical only"))
			continue
		}

		if elapsed > tolerance {
			log.Printf("SKIP broadcast: time_log '%s' is %v ago, exceeds tolerance %v | channel=%s", publish.TimeLog, elapsed.Round(time.Second), tolerance, publish.Channel)

			_ = ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("SKIP: time_log expired (%s ago > %d min tolerance), saved as historical only", elapsed.Round(time.Second).String(), toleranceMinutes)))
			continue
		}

		mutex.Lock()
		if subscribers, ok := channels[publish.Channel]; ok {
			for conn := range subscribers {
				if err := conn.WriteMessage(websocket.TextMessage, publishJSON); err != nil {
					log.Printf("Error sending to subscriber on channel %s: %v\n", publish.Channel, err)
					conn.Close()
					delete(subscribers, conn)
				}
			}
			if len(subscribers) == 0 {
				delete(channels, publish.Channel)
			}
		}
		mutex.Unlock()

		if err := ws.WriteMessage(websocket.TextMessage, []byte("ACK: broadcasted to channel "+publish.Channel)); err != nil {
			log.Printf("Error sending ACK to publisher: %v\n", err)
		}
	}
}

// --- Get history from Redis ---
func GetNotifications(channel string) ([]notification.NotificationMessage, error) {
	key := fmt.Sprintf("notification:%s", channel)
	values, err := rdb.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var history []notification.NotificationMessage
	for _, v := range values {
		var msg notification.NotificationMessage
		if err := json.Unmarshal([]byte(v), &msg); err != nil {
			log.Printf("Error unmarshaling notification: %v\n", err)
			continue
		}
		history = append(history, msg)
	}

	sort.Slice(history, func(i, j int) bool {
		ti, err1 := time.Parse("2006-01-02 15:04:05", history[i].Timestamp)
		tj, err2 := time.Parse("2006-01-02 15:04:05", history[j].Timestamp)
		if err1 != nil || err2 != nil {
			return history[i].Timestamp > history[j].Timestamp
		}
		return ti.After(tj)
	})

	return history, nil
}

// --- Send notification to server ---
func SendNotification(notificationMessage notification.NotificationMessage) error {
	notificationMessage.Timestamp = time.Now().Format("2006-01-02 15:04:05")
	serverURL := config.Get().WebSocketServer

	conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to WebSocket server: %w", err)
	}
	defer conn.Close()

	notificationMessage.Channel = fmt.Sprintf("%s:%s:%s",
		notificationMessage.CustomerNo,
		notificationMessage.Site,
		notificationMessage.Channel,
	)

	message, err := json.Marshal(notificationMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	_, response, err := conn.ReadMessage()
	if err != nil {
		return fmt.Errorf("failed to read acknowledgment: %w", err)
	}

	log.Printf("ACK from server: %s\n", string(response))
	return nil
}

// --- Cleanup all connections ---
func CleanupWebSocket() {
	mutex.Lock()
	defer mutex.Unlock()
	for channel, subscribers := range channels {
		for conn := range subscribers {
			conn.Close()
		}
		delete(channels, channel)
	}
	log.Println("Cleaned up WebSocket connections")
}
