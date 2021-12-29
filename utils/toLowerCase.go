package utils

import "strings"

func ToLowerCase(str string) string {

	var b strings.Builder

	b.WriteString(strings.ToLower(string(str)))

	return b.String()

}
