package services

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/llms"
)

type FlareService struct {
	llm        llms.Model
	userPrompt string
	temp       float64
}

func NewFlareService(llm llms.Model, userPrompt string, temp float64) *FlareService {
	return &FlareService{
		llm:        llm,
		userPrompt: userPrompt,
		temp:       temp,
	}
}

func (m *FlareService) GenerateFlares(ctx context.Context, input string, previousResponses []string) (string, error) {
	var messages []llms.MessageContent
	prompt := m.userPrompt
	if input != "" {
		prompt = m.userPrompt + "\nThe character details are as follows: \n" + input
	}
	newMessage := llms.MessageContent{
		Role:  llms.ChatMessageTypeHuman,
		Parts: []llms.ContentPart{llms.TextPart(prompt)},
	}
	messages = append(messages, newMessage)
	if len(previousResponses) > 0 {
		previous := "Here are your previous responses. Create a new one that is different:\n"
		for _, msg := range previousResponses {
			previous += msg + "\n"
		}
		messages = append(messages, llms.MessageContent{
			Role:  llms.ChatMessageTypeHuman,
			Parts: []llms.ContentPart{llms.TextPart(previous)},
		})
	}
	response, err := m.llm.GenerateContent(ctx, messages, llms.WithTemperature(m.temp))
	if err != nil {
		return "", fmt.Errorf("failed to generate insult: %w", err)
	}
	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from LLM")
	}
	content := response.Choices[0].Content
	messages = append(messages, llms.MessageContent{
		Role:  llms.ChatMessageTypeAI,
		Parts: []llms.ContentPart{llms.TextPart(content)},
	})
	return content, nil
}
