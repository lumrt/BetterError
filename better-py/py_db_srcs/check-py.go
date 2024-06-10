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
	"division by zero":                "division par zéro",
	"SyntaxError: expected '([^']+)'": "Erreur de syntaxe, un '%s' est attendu",
	"name '([^']+)' is not defined":   "le nom '%s' n'est pas défini, as-tu initialisé toutes\n tes variables ?",
	"invalid syntax":                  "syntaxe invalide",
	"list index out of range":         "index de liste hors de portée, tu parcours ta liste trop loin !",
	"tuple index out of range":        "index de tuple hors de portée, ton tuples n'est pas si grand !",
	"string index out of range":       "index de chaîne hors de portée, tu parcours ta chaine trop loin !",
	"zero division error":             "erreur de division par zéro, tu n'as pas le droit, verifie que tout tes diviseurs sont bien != 0",
	"key error":                       "erreur de clé",
	"type error":                      "erreur de type, on additionne les bananes avec les bananes, les int avec les int, les strings avec les strings...",
	"value error":                     "erreur de valeur",
	"attribute error":                 "erreur d'attribut",
	"import error":                    "erreur d'importation",
	"module not found error":          "erreur de module non trouvé",
	"indentation error":               "erreur d'indentation, regarde si tes instructions sont bien à 1 tab de tes operations",
	"tab error":                       "erreur de tabulation, ne les melangent pas avec les espaces",
	"unbound local error":             "erreur de variable locale non liée",
	"recursion error":                 "erreur de récursion",
	"memory error":                    "erreur de mémoire",
	"overflow error":                  "erreur de dépassement",
	"EOFError":                        "erreur de fin de fichier",
	"OSError":                         "erreur du système d'exploitation",
	"FileNotFoundError":               "fichier non trouvé",
	"IsADirectoryError":               "c'est un répertoire",
	"NotADirectoryError":              "ce n'est pas un répertoire",
	"PermissionError":                 "erreur de permission",
	"TimeoutError":                    "erreur de délai d'attente",
	"object '([^']+)' is not defined": "l'objet '%s' n'est pas défini",
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

	filename := os.Args[1]
	script, err := readFile(filename)
	if err != nil {
		fmt.Println("Erreur lors de la lecture du fichier:", err)
		return
	}
	output, err := executePythonScript(script)
	if err != nil {
		translatedError := translateErrorMessage(output)
		fmt.Println("Erreur capturée et traduite : ", translatedError)
	} else {
		fmt.Println("Sortie du script : ", output)
	}
}
