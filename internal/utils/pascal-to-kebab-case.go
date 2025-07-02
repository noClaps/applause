package utils

import "strings"

func PascalToKebabCase(pascal string) string {
	kebab := strings.ToLower(pascal[0:1])
	for _, c := range pascal[1:] {
		if strings.ToUpper(string(c)) == string(c) {
			kebab += "-"
		}
		kebab += strings.ToLower(string(c))
	}
	return kebab
}
