package main

import (
	"fmt"
	"log"

	"github.com/Kimoto-Norihiro/gemini-line/utils"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

func main() {
	setting, err := utils.LoadSetting()
	if err != nil {
		fmt.Println(err)
	}

	bot, err := messaging_api.NewMessagingApiAPI(
		setting.ChannelAccessToken,
	)

	r := gin.Default()

	r.POST("/", func(c *gin.Context) {
		cb, err := webhook.ParseRequest(setting.ChannelSecret, c.Request)
		if err != nil {
			fmt.Println(err)
		}

		for _, event := range cb.Events {
			switch e := event.(type) {
			case webhook.MessageEvent:
				switch message := e.Message.(type) {
				case webhook.TextMessageContent:
					switch message.Text {
					case "登録":
						fmt.Println("登録")
					case "一覧":
						fmt.Println("一覧")
					case "削除":
						fmt.Println("削除")
					case "メニュー":
						fmt.Println("メニュー")
					default:
						fmt.Println("default")
					}

					userId, ok := e.Source.(webhook.UserSource)
					if !ok {
						log.Printf("Unsupported source type: %T\n", e.Source)
						return
					}
					log.Printf("userId: %s\n", userId)

					quickMessage := &messaging_api.TextMessage{
						Text: "メニューを選択してください",
						QuickReply: &messaging_api.QuickReply{
							Items: []messaging_api.QuickReplyItem{
								{
									Type: "action",
									Action: &messaging_api.MessageAction{
										Label: "登録", Text: "登録",
									},
								},
								{
									Type: "action",
									Action: &messaging_api.MessageAction{
										Label: "一覧", Text: "一覧",
									},
								},
								{
									Type: "action",
									Action: &messaging_api.MessageAction{
										Label: "削除", Text: "削除",
									},
								},
							},
						},
					}

					messageRequest := &messaging_api.ReplyMessageRequest{
						ReplyToken: e.ReplyToken,
						Messages: []messaging_api.MessageInterface{
							quickMessage,
						},
					}

					if _, err = bot.ReplyMessage(messageRequest); err != nil {
						log.Print(err)
					} else {
						log.Println("Sent text reply.")
					}
				default:
					log.Printf("Unsupported message type: %T\n", message)
				}
			default:
				log.Printf("Unsupported event type: %T\n", event)
			}
		}
	})

	r.Run(":8080")
}
