package utility

func StrLength(str string) int {
	return len([]rune(str))
}

func IsStrEmpty(str string) bool {
	return StrLength(str) <= 0
}
