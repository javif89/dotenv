package dotenv

import "strings"

func cleanString(s string) string {
	s = strings.Trim(s, " ")
	s = strings.Trim(s, "\n")
	s = strings.Trim(s, "\r")
	s = strings.Trim(s, "\"")
	return s
}