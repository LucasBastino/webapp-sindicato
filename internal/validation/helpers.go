package validation

import (
	"fmt"
	"strings"
	"time"
)

func removeChars(input, charsToRemove string) string {
	// strings.Map recibe un string y retorna otro segun las condiciones impuestas
	// internamente recorre el string rune por rune y
	// dependiendo las condiciones agrega o no el caracter al string final
	return strings.Map(func(r rune) rune {
		if strings.ContainsRune(charsToRemove, r) {
			return -1
		}
		return r
	}, input)
}

func hasXLength(input string, length int) string {
	if len(input) != length {
		return fmt.Sprintf("Debe contener %d caracteres.", length)
	}
	return ""
}

func isNotLongerThan(input string, limit int) string {
	if len(input) > limit {
		return fmt.Sprintf("No puede contener mas de %d caracteres.", limit)
	}
	return ""
}

func hasAtLeast(input string, limit int) string {
	if len(input) < limit {
		return fmt.Sprintf("Debe contener al menos %d caracteres.", limit)
	}
	return ""
}

func isAlphabetic(input string) string {
	// se incluye Ã por la codificacion
	validChars := "abcdefghijklmnñopqrstuvwxyzáéíóúüÃ"
	return isValidCharacter(input, validChars)
}

func isNumeric(input string) string {
	validChars := "0123456789"
	return isValidCharacter(input, validChars)
}

func isAlphanumeric(input string) string {
	// se incluye Ã por la codificacion
	validChars := "abcdefghijklmnñopqrstuvwxyzáéíóúüÃ0123456789"
	return isValidCharacter(input, validChars)
}

func buildMap(validChars string) map[rune]struct{}{
	charsMap := make(map[rune]struct{}, len(validChars))
	for _, r := range validChars{
		charsMap[r] = struct{}{}
	}
	return charsMap
}

func isValidCharacter(input, validChars string) string {
	input = strings.ToLower(input)
	charsMap := buildMap(validChars)
	for _, r := range input {
		if _, ok := charsMap[r]; !ok  {
			return "Caracter no válido."
		}
	}
	return ""
}

func isValidOption(input string, optionsSet map[string]struct{}) string {
	// chequeo que exista la key dentro del map, si existe todo okay
	// en el "set" lo unico importante es la key, el value no importa
	// y como value tiene un struct vacio porque ocupa 0 bytes,
	// no le pongo bool porque este sí ocupa espacio
	if _, exists := optionsSet[input]; exists {
		return ""
	}
	return "Opción no válida."
}

func isValidDate(input string, min, max time.Time) string {
	date, err := ParseDMY(input)
	if err != nil {
		return "Fecha no válida."
	}
	if date.Before(min) || date.After(max) {
		return "Fecha no válida."
	}
	return ""
}

// ParseDMY accepts 1 or 2 digit day/month (5/6/2026 and 05/06/2026).
func ParseDMY(input string) (time.Time, error) {
	return time.Parse("2/1/2006", strings.TrimSpace(input))
}

// func isBoolean(input string) string {
// 	if input == "true" || input == "false" {
// 		return ""
// 	}
// 	return "El valor ingresado no es válido."
// }
