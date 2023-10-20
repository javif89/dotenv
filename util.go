package dotenv

import (
	"os"
	"strings"
)

func cleanString(s string) string {
	s = strings.Trim(s, " ")
	s = strings.Trim(s, "\n")
	s = strings.Trim(s, "\r")
	s = strings.Trim(s, "\"")
	return s
}

func standardizeKey(s string) string {
	s = cleanString(s)
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ToUpper(s)

	return s
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func formatValueForPrint(s string) string {
	s = cleanString(s)
	// If it's multiple words add quotes
	if strings.Contains(s, " ") || strings.Contains(s, "-") || strings.Contains(s, "+") {
		s = "\"" + s + "\""
	}
	return s
}