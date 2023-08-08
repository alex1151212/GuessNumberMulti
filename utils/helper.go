package utils

// 驗證值是否為不重複的4位數
func ValidateNumber(number string) bool {
	if len(number) != 4 {
		return false
	}

	// 檢查是否有重複的數字
	seen := make(map[rune]bool)
	for _, digit := range number {
		if seen[digit] {
			return false
		}
		seen[digit] = true
	}

	return true
}
