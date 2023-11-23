package internal

import (
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func ToTitleCase(s string) string {
	t := cases.Title(language.English)
	return t.String(s)
}

// Función para extraer el valor asociado a una clave en la respuesta de whois
func ExtractValue(response, key string) (string, error) {

	// Ejemplo simple:
	startIndex := strings.Index(response, key)
	if startIndex == -1 {
		return "", fmt.Errorf("key not found: %s", key)
	}

	startIndex += len(key)
	endIndex := strings.Index(response[startIndex:], "\n") + startIndex

	if endIndex == -1 {
		return "", fmt.Errorf("failed to extract value for key: %s", key)
	}

	return strings.TrimSpace(response[startIndex:endIndex]), nil
}
