package gemini

import (
	"context"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type Gemini interface{
	GenContent(ctx context.Context, chatText string) ([]string, error)
}

type gemini struct {
	client *genai.Client
}

func NewGemini(apiKey string, ctx context.Context) (Gemini, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	defer client.Close()

	return &gemini{
		client: client,
	}, nil
}

func (g *gemini) GenContent(ctx context.Context, chatText string) ([]string, error) {
	model := g.client.GenerativeModel("gemini-pro")
	cs := model.StartChat()

	resp, err := cs.SendMessage(ctx, genai.Text(chatText))
	if err != nil {
		return nil, err
	}
	respStr := contentResToStrings(resp)

	return respStr, nil
}

func contentResToStrings(resp *genai.GenerateContentResponse) []string {
	var res []string
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				partStr, ok := part.(*genai.Text)
				if !ok {
					continue
				}
				res = append(res, string(*partStr))
			}
		}
	}
	return res
}
