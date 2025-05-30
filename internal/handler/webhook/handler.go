package webhook

import (
	"net/http"
	"strconv"

	"opsalert/internal/service/webhook"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
)

type Handler struct {
	bot     *linebot.Client
	service *webhook.Service
}

func NewHandler(bot *linebot.Client, service *webhook.Service) *Handler {
	return &Handler{
		bot:     bot,
		service: service,
	}
}

// HandleLineWebhook รับ webhook จาก LINE
func (h *Handler) HandleLineWebhook(c *gin.Context) {
	oaID, err := strconv.Atoi(c.Param("oa_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid OA ID"})
		return
	}

	// รับ webhook events
	events, err := h.bot.ParseRequest(c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid webhook request"})
		return
	}

	// ประมวลผลแต่ละ event
	for _, event := range events {
		switch event.Type {
		case linebot.EventTypeMessage:
			// ถ้าเป็นข้อความตอบกลับ
			if message, ok := event.Message.(*linebot.TextMessage); ok {
				// ส่ง webhook ไปยังระบบภายนอก
				webhookData := map[string]interface{}{
					"type":      "message",
					"oa_id":     oaID,
					"user_id":   event.Source.UserID,
					"message":   message.Text,
					"timestamp": event.Timestamp,
				}

				if err := h.service.ForwardWebhook(c.Request.Context(), oaID, webhookData); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to forward webhook"})
					return
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
