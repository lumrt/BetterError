/* ************************************************************************** */
/*                                                                            */
/*                                                        :::      ::::::::   */
/*   main.go                                            :+:      :+:    :+:   */
/*                                                    +:+ +:+         +:+     */
/*   By: lucas <lucas@student.42.fr>                +#+  +:+       +#+        */
/*                                                +#+#+#+#+#+   +#+           */
/*   Created: 2024/06/05 13:29:38 by lucas             #+#    #+#             */
/*   Updated: 2024/06/10 13:47:35 by lucas            ###   ########.fr       */
/*                                                                            */
/* ************************************************************************** */

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
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

func translateErrorWithChatGPT(apiKey, errorMsg string) (string, error) {

	url := "https://api.openai.com/v1/chat/completions"

	// request structure
	requestBody := ChatGPTRequest{
		Model: "gpt-3.5",
		Messages: []Message{
			{Role: "system", Content: "You are a helpful assistant that translates Python error messages from English to French."},
			{Role: "user", Content: errorMsg},
		},
	}

	// from go struct to JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	// POST http request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	// init request header
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// request output
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// reading resp body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// from JSON to go struct
	var chatGPTResponse ChatGPTResponse
	if err := json.Unmarshal(body, &chatGPTResponse); err != nil {
		return "", err
	}

	// extraction of the first traduction in the struct
	if len(chatGPTResponse.Choices) > 0 {
		return chatGPTResponse.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no translation found in the response")
}

var errorTranslations = map[string]string{
	"division by zero":          "division par zéro",
	"name is not defined":       "objet non défini n'est pas défini",
	"invalid syntax":            "Erreur de Syntaxe",
	"SyntaxError: expected ':'": "Tu as oublié ':'",
}

func translateErrorMessage(errorMsg string) string {
	for en, fr := range errorTranslations {
		if strings.Contains(errorMsg, en) {
			return strings.ReplaceAll(errorMsg, en, fr)
		}
	}
	return errorMsg
}

func readFile(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func executePythonScript(script string) (string, error) {
	cmd := exec.Command("python3", "-c", script)
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
		fmt.Println("Usage: go run main.go <filename.py>")
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

	output, err := executePythonScript(script)
	if err != nil {
		translatedError, err := translateErrorWithChatGPT(apiKey, output)
		if err != nil {
			fmt.Println("Erreur lors de la traduction avec ChatGPT:", err)
		} else {
			fmt.Println("Erreur capturée et traduite :", translatedError)
		}
	} else {
		fmt.Println("Sortie du script :", output)
	}
}
