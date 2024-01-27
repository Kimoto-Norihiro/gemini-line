package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Kimoto-Norihiro/gemini-line/utils"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func main() {
	setting, err := utils.LoadSetting()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(setting.GeminiApiKey))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-pro")
	cs := model.StartChat()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("You: ")
		scanner.Scan()
		text := scanner.Text()
		fmt.Println("---")

		resp, err := cs.SendMessage(ctx, genai.Text(text))
		if err != nil {
			log.Println("send message error:")
			log.Fatal(err)
		}
		printResponse(resp)
	}
}

func printResponse(resp *genai.GenerateContentResponse) {
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				fmt.Print("Gemini: ")
				fmt.Println(part)
			}
		}
	}
	fmt.Println("---")
}
