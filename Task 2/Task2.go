package Task2

import (
	"strings"
)

func WordFrequencyCount(str string) map[string]int {
	i, j, ans := 0, 0, make(map[string]int)
	for j < len(str) {
		// Assuming Numbers are part of words
		for ; j < len(str) && ((65 <= str[j] && str[j] <= 90) || (97 <= str[j] && str[j] <= 122) || (48 <= str[j] && str[j] < 58)); j++ {
		}
		if i != j {
			ans[strings.ToLower(str[i:j])]++
		} else {
			j++
		}
		i = j
	}
	return ans
}

func IsPalindrome(str string) bool {
	i, j := 0, len(str)-1
	for i < j {
		for ; i < j && !((65 <= str[i] && str[i] <= 90) || (97 <= str[i] && str[i] <= 122) || (48 <= str[i] && str[i] <= 58)); i++ {
		}
		for ; j > i && !((65 <= str[j] && str[j] <= 90) || (97 <= str[j] && str[j] <= 122) || (48 <= str[j] && str[j] <= 58)); j-- {
		}
		if !strings.EqualFold(str[i:i+1], str[j:j+1]) {
			return false
		}
		i++
		j--
	}
	return true
}
