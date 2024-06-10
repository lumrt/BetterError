package main

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

func TranslateErrorWithChatGPT(apiKey, errorMsg string) (string, error) {
	client := openai.NewClient(apiKey)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    "Instructor",
					Content: "I am a computer science instructor translating error messages from English to French.",
				},
				{
					Role:    "Student",
					Content: errorMsg,
				},
			},
		},
	)
	if err != nil {
		return "", fmt.Errorf("erreur lors de la demande de traduction avec OpenAI: %v", err)
	}

	if len(resp.Choices) > 0 {
		return resp.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("aucune traduction trouvée dans la réponse")
}

func main() {
	apiKey := ""
	if apiKey == "" {
		fmt.Println("L'environnement 'OPENAI_API_KEY' n'est pas défini.")
		return
	}

	errorMsg := "module not found error"

	translatedError, err := TranslateErrorWithChatGPT(apiKey, errorMsg)
	if err != nil {
		fmt.Println("Erreur:", err)
		return
	}

	fmt.Println("Erreur traduite:", translatedError)
}
