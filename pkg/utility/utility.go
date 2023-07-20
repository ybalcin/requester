package utility

import "strings"

func StrLength(str string) int {
	str = strings.TrimSpace(str)
	return len([]rune(str))
}

func IsStrEmpty(str string) bool {
	return StrLength(str) <= 0
}

func IsSchemeExistInURL(address string) bool {
	if strings.HasPrefix(address, "http://") || strings.HasPrefix(address, "https://") {
		return true
	}

	return false
}
