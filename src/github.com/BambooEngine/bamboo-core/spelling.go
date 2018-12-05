package bamboo

import (
	"log"
	"strings"
	"unicode"
)

var firstConsonantSeq = [3]string{
	"b d đ g gh m n nh p ph r s t tr v z",
	"c h k kh qu th",
	"ch gi l ng ngh x",
}

var vowelSeq = [6]string{
	"ê i ua uê uy y",
	"a iê oa uyê yê",
	"â ă e o oo ô ơ oe u ư uâ uô ươ",
	"oă",
	"uơ",
	"ai ao au âu ay ây eo êu ia iêu iu oai oao oay oeo oi ôi ơi ưa uây ui ưi uôi ươi ươu ưu uya uyu yêu",
}

var lastConsonantSeq = [3]string{
	"ch nh",
	"c ng",
	"m n p t",
}

var cvMatrix = [3][]uint{
	{0, 1, 2, 5},
	{0, 1, 2, 3, 4, 5},
	{0, 1, 2, 3, 5},
}

var vcMatrix = [6][]uint{
	{0, 2},
	{0, 1, 2},
	{1, 2},
	{1, 2},
}

var spellingTrie = &W{F: false}

func attachSound(str string, s Sound) []Sound {
	var sounds []Sound
	for _ = range []rune(str) {
		sounds = append(sounds, s)
	}
	return sounds
}

func buildCV(consonants []string, vowels []string) []string {
	var ret []string
	for _, c := range consonants {
		for _, v := range vowels {
			ret = append(ret, c+v)
			var sounds []Sound
			sounds = append(sounds, attachSound(c, FirstConsonantSound)...)
			sounds = append(sounds, attachSound(v, VowelSound)...)
			AddTrie(spellingTrie, []rune(c+v), false, sounds)
		}
	}
	return ret
}

func generateVowels() []string {
	var ret []string
	for _, vRow := range vowelSeq {
		for _, v := range strings.Split(vRow, " ") {
			ret = append(ret, v)
			var sounds []Sound
			sounds = append(sounds, attachSound(v, VowelSound)...)
			AddTrie(spellingTrie, []rune(v), false, sounds)
		}
	}
	return ret
}

func buildVC(vowels []string, consonants []string) []string {
	var ret []string
	for _, v := range vowels {
		for _, c := range consonants {
			ret = append(ret, v+c)
			var sounds []Sound
			sounds = append(sounds, attachSound(v, VowelSound)...)
			sounds = append(sounds, attachSound(c, LastConsonantSound)...)
			AddTrie(spellingTrie, []rune(v+c), false, sounds)
		}
	}
	return ret
}

func buildCVC(cs1 []string, vs1 []string, cs2 []string) []string {
	var ret []string
	for _, c1 := range cs1 {
		for _, v := range vs1 {
			for _, c2 := range cs2 {
				ret = append(ret, c1+v+c2)
				var sounds []Sound
				sounds = append(sounds, attachSound(c1, FirstConsonantSound)...)
				sounds = append(sounds, attachSound(v, VowelSound)...)
				sounds = append(sounds, attachSound(c2, LastConsonantSound)...)
				AddTrie(spellingTrie, []rune(c1+v+c2), false, sounds)
			}
		}
	}
	return ret
}

func init() {
	generateVowels()
	generateCV()
	generateVC()
	generateCVC()
	// fix gi+? vs g+i? confusion
	buildCV([]string{"g"}, []string{"i"})
}

func generateCV() []string {
	var ret []string
	for cRow, vRows := range cvMatrix {
		for _, vRow := range vRows {
			var consonants = strings.Split(firstConsonantSeq[cRow], " ")
			var vowels = strings.Split(vowelSeq[vRow], " ")
			ret = append(ret, buildCV(consonants, vowels)...)
		}
	}
	return ret
}

func generateVC() []string {
	var ret []string
	for vRow, cRows := range vcMatrix {
		for _, cRow := range cRows {
			var vowels = strings.Split(vowelSeq[vRow], " ")
			var consonants = strings.Split(lastConsonantSeq[cRow], " ")
			ret = append(ret, buildVC(vowels, consonants)...)
		}
	}
	return ret
}

func generateCVC() []string {
	var ret []string
	for c1Row, vRows := range cvMatrix {
		for _, vRow := range vRows {
			for _, c2Row := range vcMatrix[vRow] {
				var cs1 = strings.Split(firstConsonantSeq[c1Row], " ")
				var vowels = strings.Split(vowelSeq[vRow], " ")
				var cs2 = strings.Split(lastConsonantSeq[c2Row], " ")
				ret = append(ret, buildCVC(cs1, vowels, cs2)...)
			}
		}
	}
	return ret
}

func GetLastCombination(composition []*Transformation) []*Transformation {
	var ret []*Transformation
	if len(composition) <= 1 {
		return composition
	}
	for i, trans := range composition {
		ret = append(ret, trans)
		str := Flatten(ret, VietnameseMode|NoTone)
		if str == "" {
			continue
		}
		if res, _ := FindWord(spellingTrie, []rune(str), false); res == FindResultNotMatch {
			if i == 0 {
				return GetLastCombination(composition[1:])
			}
			return GetLastCombination(composition[i:])
		}
	}
	return ret
}

// only appending trans has sound
func getCombinationWithSound(composition []*Transformation) ([]*Transformation, []Sound) {
	var lastComb = getAppendingComposition(composition)
	if len(lastComb) <= 0 {
		return lastComb, GenerateDumpSoundFromTonelessWord("")
	}
	var str = Flatten(lastComb, VietnameseMode|NoTone|LowerCase)
	if res, sounds := FindWord(spellingTrie, []rune(str), false); res != FindResultNotMatch {
		return lastComb, sounds
	}
	return lastComb, GenerateDumpSoundFromTonelessWord(str)
}

func getCompositionBySound(composition []*Transformation, sound Sound) []*Transformation {
	var lastComb, sounds = getCombinationWithSound(composition)
	if len(lastComb) != len(sounds) {
		log.Println("Something is wrong with the length of sounds")
		return lastComb
	}
	var ret []*Transformation
	for i, s := range sounds {
		if s == sound {
			ret = append(ret, lastComb[i])
		}
	}
	return ret
}

func getSpellingMatchResult(composition []*Transformation, mode Mode) (uint8, []Sound) {
	if len(composition) <= 0 {
		return FindResultMatchFull, []Sound{}
	}
	if mode&NoTone != 0 {
		str := Flatten(composition, NoTone|LowerCase)
		var chars = []rune(str)
		if len(chars) <= 1 {
			return FindResultMatchFull, GenerateDumpSoundFromTonelessWord(str)
		}
		return FindWord(spellingTrie, chars, false)
	}
	return FindResultNotMatch, []Sound{}
}

func isSpellingCorrect(composition []*Transformation, mode Mode) bool {
	res, _ := getSpellingMatchResult(composition, mode)
	return res == FindResultMatchFull
}

func isSpellingLikelyCorrect(composition []*Transformation, mode Mode) bool {
	res, _ := getSpellingMatchResult(composition, mode)
	return res == FindResultMatchPrefix
}

func GetSoundMap(composition []*Transformation) map[*Transformation]Sound {
	var soundMap = map[*Transformation]Sound{}
	var lastComb, sounds = getCombinationWithSound(composition)
	if len(sounds) <= 0 || len(sounds) != len(lastComb) {
		log.Println("Something is wrong with the length of sounds")
		return soundMap
	}
	for i, trans := range lastComb {
		soundMap[trans] = sounds[i]
	}
	return soundMap
}

func ParseSoundsFromTonelessWord(word string) []Sound {
	res, sounds := FindWord(spellingTrie, []rune(word), false)
	if res == FindResultNotMatch {
		return GenerateDumpSoundFromTonelessWord(word)
	}
	return sounds
}

func getRightMostVowels(composition []*Transformation) []*Transformation {
	return getCompositionBySound(composition, VowelSound)
}

func getRightMostVowelWithMarks(composition []*Transformation) []*Transformation {
	var vowels = getRightMostVowels(composition)
	return addMarksToComposition(composition, vowels)
}

func GenerateDumpSoundFromTonelessWord(word string) []Sound {
	var sounds []Sound
	for _, c := range []rune(word) {
		if IsVowel(c) {
			sounds = append(sounds, VowelSound)
		} else if unicode.IsLetter(c) {
			sounds = append(sounds, FirstConsonantSound)
		} else {
			sounds = append(sounds, NoSound)
		}
	}
	return sounds
}
