package bamboo

import (
	"log"
	"unicode"
)

func isSpellingCorrect(composition []*Transformation, mode Mode) bool {
	if len(composition) <= 1 {
		return true
	}
	if mode&NoMark != 0 && mode&NoTone != 0 {
		str := Flatten(composition, NoMark|NoTone|LowerCase)
		if len([]rune(str)) <= 1 {
			return true
		}
		ok := LookupTypingDictionary(str)
		return ok
	}
	if mode&NoTone != 0 {
		str := Flatten(composition, NoTone|LowerCase)
		var chars = []rune(str)
		if len(chars) <= 1 {
			return true
		}
		ok, _ := LookupDictionary(str)
		return ok
	}
	return false
}

func isSpellingLikelyCorrect(composition []*Transformation, mode Mode) bool {
	if len(composition) <= 1 {
		return true
	}
	if isSpellingCorrect(composition, mode) {
		return true
	}
	if mode&NoTone != 0 {
		str := Flatten(composition, NoTone|LowerCase)
		if len([]rune(str)) <= 1 {
			return true
		}
		return LookupTypingDictionary(str)
	}
	return false
}

func GetLastSoundGroup(composition []*Transformation, sound Sound) []*Transformation {
	var appendingTransformations = getAppendingComposition(composition)
	if len(appendingTransformations) == 0 {
		return appendingTransformations
	}
	var str = Flatten(composition, NoTone|LowerCase)
	var sounds = ParseSoundsFromTonelessWord(str)
	if len(appendingTransformations) != len(sounds) {
		log.Println("The length of appending composition and sounds is not the same.")
		return appendingTransformations
	}
	var j = 0
	var i = 0
	for i = len(sounds) - 1; i >= 0; i-- {
		if sounds[i] == sound {
			j++
			if i == 0 || sounds[i-1] != sound {
				break
			}
		}
	}
	if i < 0 {
		i = 0
	}
	return appendingTransformations[i:(j + i)]
}

func GetLastCombination(composition []*Transformation) []*Transformation {
	var combinations = separateComposition(composition)
	if len(combinations) > 0 {
		return combinations[len(combinations)-1]
	}
	return composition
}

func GetSoundMap(composition []*Transformation) map[*Transformation]Sound {
	var soundMap = map[*Transformation]Sound{}
	var appendingTransformations = getAppendingComposition(composition)
	var str = Flatten(composition, NoTone|LowerCase)
	var sounds = ParseSoundsFromTonelessWord(str)
	if len(sounds) <= 0 || len(sounds) != len(appendingTransformations) {
		log.Println("Something wrong with the length of sounds")
		return soundMap
	}
	for i, trans := range appendingTransformations {
		soundMap[trans] = sounds[i]
	}
	return soundMap
}

func ParseSoundsFromTonelessString(str string) []Sound {
	var sounds []Sound
	for _, word := range SplitStringToWords(str) {
		sounds = append(sounds, ParseSoundsFromTonelessWord(word)...)
	}
	return sounds
}

func ParseSoundsFromLastTonelessWord(word string) []Sound {
	words := SplitStringToWords(word)
	if len(words) > 1 {
		return ParseSoundsFromTonelessWord(words[len(words)-1])
	} else {
		return ParseSoundsFromTonelessWord(word)
	}
}

func ParseSoundsFromTonelessWord(word string) []Sound {
	var sounds []Sound
	if found, sounds := LookupDictionary(word); found {
		return sounds
	}
	for _, c := range []rune(word) {
		if isVowel(c) {
			sounds = append(sounds, VowelSound)
		} else if unicode.IsLetter(c) {
			sounds = append(sounds, FirstConsonantSound)
		} else {
			sounds = append(sounds, NoSound)
		}
	}
	return sounds
}
