/* ************************************************************************** */
/*                                                                            */
/*                                                        :::      ::::::::   */
/*   check-js-api.go                                    :+:      :+:    :+:   */
/*                                                    +:+ +:+         +:+     */
/*   By: lucas <lucas@student.42.fr>                +#+  +:+       +#+        */
/*                                                +#+#+#+#+#+   +#+           */
/*   Created: 2024/06/10 14:34:50 by lucas             #+#    #+#             */
/*   Updated: 2024/07/10 19:57:49 by lucas            ###   ########.fr       */
/*                                                                            */
/* ************************************************************************** */

package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/sashabaranov/go-openai"
)

/* implementation of api request to translate*/

type ChatGPTRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatGPTResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

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

func readFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func executeJSScript(script string) (string, error) {
	cmd := exec.Command("node", script)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return stderr.String(), err
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <filename.js>")
		return
	}

	// check api key is inserted
	apiKey := ""
	if apiKey == "" {
		fmt.Println("environnement 'OPENAI_API_KEY' is not defined. go better-py/py_api_src/main.go l.146")
		return
	}

	filename := os.Args[1]
	script, err := readFile(filename)
	if err != nil {
		fmt.Println("Erreur lors de la lecture du fichier:", err)
		return
	}

	output, err := executeJSScript(script)
	if err != nil {
		translatedError, err := TranslateErrorWithChatGPT(apiKey, output)
		if err != nil {
			fmt.Println("Erreur lors de la traduction avec ChatGPT:", err)
		} else {
			fmt.Println("Erreur capturée et traduite :", translatedError)
		}
	} else {
		fmt.Println("Sortie du script :", output)
	}
}
