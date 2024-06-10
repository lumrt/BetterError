package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var errorTranslations = map[string]string{
	"TypeError: Cannot read property '([^']+)' of undefined": "Erreur de type: impossible de lire la propriété '%s' de undefined",
	"ReferenceError: ([^ ]+) is not defined":                 "Erreur de référence: '%s' n'est pas défini, as-tu initialisé toutes tes variables ?",
	"SyntaxError: Unexpected token '([^']+)'":                "Erreur de syntaxe: jeton inattendu '%s'",
	"SyntaxError: Unexpected end of input":                   "Erreur de syntaxe: fin d'entrée inattendue",
	"TypeError: ([^ ]+) is not a function":                   "Erreur de type: '%s' n'est pas une fonction",
	"RangeError: Maximum call stack size exceeded":           "Erreur de portée: taille maximale de la pile d'appels dépassée",
	"TypeError: Cannot set property '([^']+)' of undefined":  "Erreur de type: impossible de définir la propriété '%s' de undefined",
	"TypeError: Cannot read properties of undefined":         "Erreur de type: impossible de lire les propriétés de undefined",
	"TypeError: ([^ ]+) is not iterable":                     "Erreur de type: '%s' n'est pas itérable",
	"TypeError: Assignment to constant variable":             "Erreur de type: affectation à une variable constante",
	"SyntaxError: Invalid or unexpected token":               "Erreur de syntaxe: jeton invalide ou inattendu",
	"TypeError: Cannot convert undefined or null to object":  "Erreur de type: impossible de convertir undefined ou null en objet",
	"TypeError: ([^ ]+) is not a constructor":                "Erreur de type: '%s' n'est pas un constructeur",
	"SyntaxError: missing ([^ ]+) after argument list":       "Erreur de syntaxe: '%s' manquante après la liste d'arguments",
	"SyntaxError: Unexpected string":                         "Erreur de syntaxe: chaîne inattendue",
}

func translateErrorMessage(errorMsg string) string {
	for en, fr := range errorTranslations {
		re := regexp.MustCompile(en)
		if re.MatchString(errorMsg) {
			matches := re.FindStringSubmatch(errorMsg)
			if len(matches) > 1 {
				translatedMsg := fmt.Sprintf(fr, matches[1])
				return strings.ReplaceAll(errorMsg, matches[0], translatedMsg)
			}
			return strings.ReplaceAll(errorMsg, en, fr)
		}
	}
	return errorMsg
}

func readFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func executeJavaScript(script string) (string, error) {
	cmd := exec.Command("node", "-e", script)
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
	filename := os.Args[1]
	script, err := readFile(filename)
	if err != nil {
		fmt.Println("Erreur lors de la lecture du fichier:", err)
		return
	}
	output, err := executeJavaScript(script)
	if err != nil {
		translatedError := translateErrorMessage(output)
		fmt.Println("Erreur capturée et traduite : ", translatedError)
	} else {
		fmt.Println("Sortie du script : ", output)
	}
}
