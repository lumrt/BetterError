package main

import (
    "bytes"
    "fmt"
    "io/ioutil"
    "os"
    "os/exec"
    "strings"
)

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
