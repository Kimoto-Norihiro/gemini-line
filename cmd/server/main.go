package main

import (
	"fmt"

	"github.com/Kimoto-Norihiro/gemini-line/linebot"
	"github.com/Kimoto-Norihiro/gemini-line/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	setting, err := utils.LoadSetting()
	if err != nil {
		fmt.Println(err)
	}

	linebot, err := linebot.NewLineBot(setting.ChannelAccessToken, setting.ChannelSecret, setting.GeminiApiKey)
	if err != nil {
		fmt.Println(err)
	}

	r := gin.Default()

	r.POST("/", linebot.ReplyMessageEvent)
	r.Run(":8080")
}
