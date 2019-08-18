/*
 * Bamboo - A Vietnamese Input method editor
 * Copyright (C) Luong Thanh Lam <ltlam93@gmail.com>
 *
 * This software is licensed under the MIT license. For more information,
 * see <https://github.com/BambooEngine/bamboo-core/blob/master/LISENCE>.
 */
package bamboo

import (
	"regexp"
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

var spellingTrie = &Node{Full: false}

func buildCV(consonants []string, vowels []string) []string {
	var ret []string
	for _, c := range consonants {
		for _, v := range vowels {
			ret = append(ret, c+v)
		}
	}
	return ret
}

func generateVowels() []string {
	var ret []string
	for _, vRow := range vowelSeq {
		for _, v := range strings.Split(vRow, " ") {
			ret = append(ret, v)
		}
	}
	return ret
}

func buildVC(vowels []string, consonants []string) []string {
	var ret []string
	for _, v := range vowels {
		for _, c := range consonants {
			ret = append(ret, v+c)
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
			}
		}
	}
	return ret
}

func init() {
	for _, word := range GenerateDictionary() {
		AddTrie(spellingTrie, []rune(word), false, false)
	}
}

func AddDictionaryToSpellingTrie(dictionary map[string]bool) {
	for word := range dictionary {
		AddTrie(spellingTrie, []rune(word), true, false)
	}
}

func GenerateDictionary() []string {
	var words = generateVowels()
	words = append(words, generateCV()...)
	words = append(words, generateVC()...)
	words = append(words, generateCVC()...)
	return words
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

var regGI = regexp.MustCompile(`^(qu|gi)(\p{L}+)`)

func ParseSoundsFromWord(word string) []Sound {
	var sounds []Sound
	var chars = []rune(word)
	if len(chars) == 0 {
		return nil
	}
	var suffix string
	if regGI.MatchString(word) {
		subs := regGI.FindStringSubmatch(word)
		if len(subs) == 3 {
			var seq = []rune(subs[2])
			if IsVowel(seq[0]) {
				sounds = append(sounds, FirstConsonantSound)
				sounds = append(sounds, FirstConsonantSound)
				suffix = subs[2]
				sounds = append(sounds, ParseDumpSoundsFromWord(suffix)...)
				return sounds
			} else {
				return ParseDumpSoundsFromWord(word)
			}
		}
	} else {
		sounds = ParseDumpSoundsFromWord(word)
	}
	return sounds
}

func ParseDumpSoundsFromWord(word string) []Sound {
	var sounds []Sound
	var hadVowel bool
	for _, c := range []rune(word) {
		if IsVowel(c) {
			sounds = append(sounds, VowelSound)
			hadVowel = true
		} else if unicode.IsLetter(c) {
			if hadVowel {
				sounds = append(sounds, LastConsonantSound)
			} else {
				sounds = append(sounds, FirstConsonantSound)
			}
		} else {
			sounds = append(sounds, NoSound)
		}
	}
	return sounds
}
