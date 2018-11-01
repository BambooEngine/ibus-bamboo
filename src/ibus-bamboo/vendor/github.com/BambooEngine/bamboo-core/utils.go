/*
 * Bamboo - A Vietnamese Input method editor
 * Copyright (C) 2018 Luong Thanh Lam <ltlam93@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program. If not, see <http://www.gnu.org/licenses/>.
 *
 */

package bamboo

import (
	"regexp"
	"strings"
	"unicode"
)

func copyRunes(r []rune) []rune {
	t := make([]rune, len(r))
	copy(t, r)

	return t
}

func isVowel(chr rune) bool {
	isVowel := false
	for _, v := range vowelSeq {
		if v == chr {
			isVowel = true
		}
	}
	return isVowel
}

func FindVowelPosition(chr rune) int {
	for pos, v := range vowelSeq {
		if v == chr {
			return pos
		}
	}
	return -1
}


func FindMissingRuleForUo(composition []*Transformation, isSuperKey bool) (Rule, bool) {
	var rule Rule
	if len(composition) < 2 {
		return rule, false
	}
	var target rune
	var full = strings.ToLower(Flatten(composition, NoTone|LowerCase))

	if !isSuperKey {
		var reg = regexp.MustCompile(`(h|th|kh)uơ\p{L}*`)
		if reg.MatchString(full) {
			target = 'u'
		}
	} else {
		var reg = regexp.MustCompile(`^(h|th|kh)uo$`)
		if reg.MatchString(full) {
			return rule, false
		}
		var vowels = GetRightMostVowelWithMarks(composition)
		var str = Flatten(vowels, NoTone|LowerCase)
		if strings.Contains(str, "uo") {
			target = 'o'
		}
	}
	if target > 0 {
		rule = Rule{
			Key:        rune(0),
			EffectType: MarkTransformation,
			Effect:     MARK_HORN,
			EffectOn:   target,
		}
		return rule, true
	}
	return rule, false
}

func FindIndexRune(chars []rune, r rune) int {
	for i, c := range chars {
		if c == r {
			return i
		}
	}
	return -1
}

func SplitStringToWords(str string) []string {
	var words []string
	var word []rune
	var seq = []rune(str)
	for i, r := range seq {
		word = append(word, r)
		// todo: need to check if r is a space
		if i+1 < len(seq)-1 && isVowel(r) && !isVowel(seq[i+1]) {
			words = append(words, string(word))
			word = []rune{}
		}
	}
	words = append(words, string(word))
	return words
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