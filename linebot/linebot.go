package linebot

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
	"google.golang.org/api/option"
)

type LineBot interface {
	ReplyMessageEvent(c *gin.Context)
}

type lineBot struct {
	bot *messaging_api.MessagingApiAPI

	channel_access_token string
	channel_secret       string

	geminiApiKey string
}

func NewLineBot(channel_access_token, channel_secret, geminiApiKey string) (LineBot, error) {
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
		geminiApiKey:         geminiApiKey,
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
		message, ok := e.Message.(webhook.TextMessageContent)
		if !ok {
			fmt.Println("Unsupported message type")
			return
		}

		client, err := genai.NewClient(c, option.WithAPIKey(l.geminiApiKey))
		if err != nil {
			fmt.Println(err)
			return
		}
		defer client.Close()
		model := client.GenerativeModel("gemini-pro")
		cs := model.StartChat()

		resp, err := cs.SendMessage(c, genai.Text(message.Text))
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, cand := range resp.Candidates {
			if cand.Content != nil {
				for _, part := range cand.Content.Parts {
					partStr, ok := part.(genai.Text)
					if !ok {
						continue
					}
					messageRequest := &messaging_api.ReplyMessageRequest{
						ReplyToken: e.ReplyToken,
						Messages: []messaging_api.MessageInterface{
							&messaging_api.TextMessage{
								Text: string(partStr),
							},
						},
					}
					if _, err = l.bot.ReplyMessage(messageRequest); err != nil {
						fmt.Println(err)
						return
					}
				}
			}
		}
	}
}

func stringsToContentReq(messageStrings []string, replyToken string) *messaging_api.ReplyMessageRequest {
	var messages []messaging_api.MessageInterface
	for _, messageStr := range messageStrings {
		messages = append(messages, &messaging_api.TextMessage{
			Text: messageStr,
		})
	}

	return &messaging_api.ReplyMessageRequest{
		ReplyToken: replyToken,
		Messages:   messages,
	}
}
