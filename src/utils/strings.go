package utils

import (
	"fmt"
	"strings"
)

var specialChars = [...]string{"\\", "+", "*", "$", "(", ")", "^", "]", "[", "?", "."}

func AddEscapeString(target string) string {

	for _, sc := range specialChars {
		escapeStr := fmt.Sprintf("\\%s", sc)

		target = strings.Replace(target, sc, escapeStr, -1)
	}

	return target
}
