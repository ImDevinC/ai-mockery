package services

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/llms"
)

type MockeryService struct {
	llm              llms.Model
	userPrompt       string
	previousMessages []string
}

func NewMockeryService(llm llms.Model, userPrompt string) *MockeryService {
	return &MockeryService{
		llm:        llm,
		userPrompt: userPrompt,
	}
}

func (m *MockeryService) GenerateInsult(ctx context.Context, input string) (string, error) {
	var messages []llms.MessageContent
	prompt := m.userPrompt + "\nThe character details are as follows: \n" + input
	newMessage := llms.MessageContent{
		Role:  llms.ChatMessageTypeHuman,
		Parts: []llms.ContentPart{llms.TextPart(prompt)},
	}
	messages = append(messages, newMessage)
	if len(m.previousMessages) > 0 {
		previous := "You have already mocked this character with the following insults, so create a new one\n"
		for _, msg := range m.previousMessages {
			previous += msg + "\n"
		}
		messages = append(messages, llms.MessageContent{
			Role:  llms.ChatMessageTypeHuman,
			Parts: []llms.ContentPart{llms.TextPart(previous)},
		})
	}
	response, err := m.llm.GenerateContent(ctx, messages)
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
	m.previousMessages = append(m.previousMessages, content)
	return content, nil
}
