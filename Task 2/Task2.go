package main

import (
	"strings"
)

func WordFrequencyCount(str string) map[string]int {
	i, j, ans := 0, 0, make(map[string]int)
	for j < len(str) {
		// Assuming Numbers are part of words
		for ; j < len(str) && ((65 <= str[j] && str[j] <= 90) || (97 <= str[j] && str[j] <= 122) || (48 <= str[j] && str[j] <= 58)); j++ {
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
