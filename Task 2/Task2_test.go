package Task2

import (
	"testing"
)

func TestWordFrequencyCount(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  map[string]int
	}{{
		name:  "Simple Sentence",
		input: "Hello world this is a test",
		want: map[string]int{
			"hello": 1,
			"world": 1,
			"this":  1,
			"is":    1,
			"a":     1,
			"test":  1,
		},
	},
		{
			name:  "Case Insensitivity",
			input: "Go is fun, GO is fast!",
			want: map[string]int{
				"go":   2,
				"is":   2,
				"fun":  1,
				"fast": 1,
			},
		},
		{
			name:  "Repeated Words",
			input: "apple banana apple orange banana apple",
			want: map[string]int{
				"apple":  3,
				"banana": 2,
				"orange": 1,
			},
		},
		{
			name:  "Including Numbers",
			input: "version2 is the 1st release for all 3 platforms",
			want: map[string]int{
				"version2":  1,
				"is":        1,
				"the":       1,
				"1st":       1,
				"release":   1,
				"for":       1,
				"all":       1,
				"3":         1,
				"platforms": 1,
			},
		},
		{
			name:  "Multiple Delimiters",
			input: "one, two... three-four! five.",
			want: map[string]int{
				"one":   1,
				"two":   1,
				"three": 1,
				"four":  1,
				"five":  1,
			},
		},
		{
			name:  "Leading and Trailing Delimiters",
			input: "  ...start and end...  ",
			want: map[string]int{
				"start": 1,
				"and":   1,
				"end":   1,
			},
		},
		{
			name:  "Empty String",
			input: "",
			want:  map[string]int{},
		},
		{
			name:  "String with Only Delimiters",
			input: "!@#$%^&*()_+=-`~[]\\{}|;':\",./<>?",
			want:  map[string]int{},
		},
		{
			name:  "Known Limitation: Non-ASCII characters",
			input: "Go is a café, not a καφές.",
			want: map[string]int{
				"go":  1,
				"is":  1,
				"a":   2,
				"caf": 1,
				"not": 1,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			correct, got := true, WordFrequencyCount(tc.input)
			for str := range got {
				correct = correct && got[str] == tc.want[str]
			}
			for str := range tc.want {
				correct = correct && got[str] == tc.want[str]
			}
			if !correct {
				t.Errorf("WordFrequencyCount(%q) = %v; want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestIsPalindrome(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "Simple Palindrome - Odd Length",
			input: "racecar",
			want:  true,
		},
		{
			name:  "Simple Palindrome - Even Length",
			input: "noon",
			want:  true,
		},
		{
			name:  "Case-Insensitive Palindrome",
			input: "RaceCar",
			want:  true,
		},
		{
			name:  "Palindrome with Punctuation and Spaces",
			input: "A man, a plan, a canal: Panama",
			want:  true,
		},
		{
			name:  "Palindrome with Numbers",
			input: "12321",
			want:  true,
		},
		{
			name:  "Mixed Alphanumeric Palindrome",
			input: "Was it a car or a cat I saw?",
			want:  true,
		},
		{
			name:  "Simple Non-Palindrome",
			input: "hello",
			want:  false,
		},
		{
			name:  "Almost a Palindrome",
			input: "racecars",
			want:  false,
		},
		{
			name:  "Empty String",
			input: "",
			want:  true,
		},
		{
			name:  "Single Alphanumeric Character",
			input: "a",
			want:  true,
		},
		{
			name:  "Single Non-Alphanumeric Character",
			input: "!",
			want:  true,
		},
		{
			name:  "String with Only Delimiters",
			input: ".,;! ",
			want:  true,
		},
		{
			name:  "Known Limitation: Non-ASCII characters are skipped",
			input: "madamé",
			want:  true,
		},
		{
			name:  "Known Limitation: Unicode palindrome fails",
			input: "上海海上",
			want:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := IsPalindrome(tc.input)
			if got != tc.want {
				t.Errorf("IsPalindrome(%q) = %v; want %v", tc.input, got, tc.want)
			}
		})
	}
}
