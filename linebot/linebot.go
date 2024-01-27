package linebot

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

type LineBot interface {
	ReplyMessageEvent(c *gin.Context)
}

type lineBot struct {
	bot *messaging_api.MessagingApiAPI

	channel_access_token string
	channel_secret       string
}

func NewLineBot(channel_access_token, channel_secret string) (LineBot, error) {
	bot, err := messaging_api.NewMessagingApiAPI(
		channel_access_token,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &lineBot{
		bot:                  bot,
		channel_access_token: channel_access_token,
		channel_secret:       channel_secret,
	}, nil
}

func (l *lineBot) ReplyMessageEvent(c *gin.Context) {
	cb, err := webhook.ParseRequest(l.channel_secret, c.Request)
	if err != nil {
		fmt.Println(err)
	}

	for _, event := range cb.Events {
		e, ok := event.(webhook.MessageEvent)
		if !ok {
			fmt.Println("Unsupported event type")
			return
		}
		massage, ok := e.Message.(webhook.TextMessageContent)
		if !ok {
			fmt.Println("Unsupported message type")
			return
		}

		// userId, ok := e.Source.(webhook.UserSource)
		// if !ok {
		// 	log.Printf("Unsupported source type: %T\n", e.Source)
		// 	return
		// }
		// log.Println(userId)

		messageRequest := &messaging_api.ReplyMessageRequest{
			ReplyToken: e.ReplyToken,
			Messages: []messaging_api.MessageInterface{
				&messaging_api.TextMessage{
					Text: massage.Text,
				},
			},
		}

		if _, err = l.bot.ReplyMessage(messageRequest); err != nil {
			fmt.Println(err)
			return
		}
	}
}
