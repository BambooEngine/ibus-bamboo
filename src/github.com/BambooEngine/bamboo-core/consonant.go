package bamboo

import "strings"

func isFirstConsonant(str string) bool {
	for _, line := range firstConsonantSeq {
		for _, consonant := range strings.Split(line, " ") {
			if str == consonant {
				return true
			}
		}
	}
	return false
}

func isLastConsonant(str string) bool {
	for _, line := range lastConsonantSeq {
		for _, consonant := range strings.Split(line, " ") {
			if str == consonant {
				return true
			}
		}
	}
	return false
}

func isTypingFirstConsonantSeq(str string) bool {
	for _, line := range typingFirstConsonantSeq {
		for _, consonant := range strings.Split(line, " ") {
			if str == consonant {
				return true
			}
		}
	}
	return false
}

func isTypingVowelSeq(str string) bool {
	for _, line := range typingVowelSeq {
		for _, vowel := range strings.Split(line, " ") {
			if str == vowel {
				return true
			}
		}
	}
	return false
}

func isTypingLastConsonant(str string) bool {
	for _, line := range typingLastConsonantSeq {
		for _, consonant := range strings.Split(line, " ") {
			if str == consonant {
				return true
			}
		}
	}
	return false
}
