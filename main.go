/* ************************************************************************** */
/*                                                                            */
/*                                                        :::      ::::::::   */
/*   main.go                                            :+:      :+:    :+:   */
/*                                                    +:+ +:+         +:+     */
/*   By: lucas <lucas@student.42.fr>                +#+  +:+       +#+        */
/*                                                +#+#+#+#+#+   +#+           */
/*   Created: 2024/06/05 13:29:38 by lucas             #+#    #+#             */
/*   Updated: 2024/06/05 14:30:44 by lucas            ###   ########.fr       */
/*                                                                            */
/* ************************************************************************** */


import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "os/exec"
    "strings"
)

/* implementation of api request to translate*/

// type ChatGPTRequest struct {
//     Model    string    `json:"model"`
//     Messages []Message `json:"messages"`
// }

// type Message struct {
//     Role    string `json:"role"`
//     Content string `json:"content"`
// }


// type ChatGPTResponse struct {
//     Choices []struct {
//         Message Message `json:"message"`
//     } `json:"choices"`
// }





// TO DO !
func translateErrorWithChatGPT(apiKey, errorMsg string) (string, error) {
    url := "https://api.openai.com/v1/chat/completions"
    requestBody := ChatGPTRequest{
        Model: "gpt-4",
        Messages: []Message{
            {Role: "system", Content: "You are a helpful assistant that translates Python error messages from English to French."},
            {Role: "user", Content: errorMsg},
        },
    }



	
	return "", fmt.Errorf("no translation found in the response")
}

var errorTranslations = map[string]string{
    "division by zero": "division par zéro",
    "name is not defined": "objet non défini n'est pas défini",
    "invalid syntax": "Erreur de Syntaxe",
	"SyntaxError: expected ':'":"Tu as oublié ':'",
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
	apiKey := os.Getenv("OPENAI_API_KEY")
    if apiKey == "" {
        fmt.Println("L'environnement 'OPENAI_API_KEY' n'est pas défini.")
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
        translatedError := translateErrorMessage(output)
        fmt.Println("Erreur capturée et traduite :", translatedError)
    } else {
        fmt.Println(output)
    }
}
