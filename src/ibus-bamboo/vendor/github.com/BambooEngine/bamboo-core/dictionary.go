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
	"strings"
)

var firstConsonants = [3]string{
	"b d đ g gh m n nh p ph r s t tr v z",
	"c h k kh qu th",
	"ch gi l ng ngh x",
}

var vowels = [6]string{
	"ê i ua uê uy y",
	"a iê oa uyê yê",
	"â ă e o oo ô ơ oe u ư uâ uô ươ",
	"oă",
	"uơ",
	"ai ao au âu ay ây eo êu ia iêu iu oai oao oay oeo oi ôi ơi ưa uây ui ưi uôi ươi ươu ưu uya uyu yêu",
}

var lastConsonants = [3]string{
	"ch nh",
	"c ng",
	"m n p t",
}

var firstConsonant_vowel = [3][]uint{
	{0, 1, 2, 5},
	{0, 1, 2, 3, 4, 5},
	{0, 1, 2, 3, 5},
}

var vowel_lastConsonant = [6][]uint{
	{0, 2},
	{0, 1, 2},
	{1, 2},
	{1, 2},
}

var dictionary map[string][]Sound

func isValid(i1, i2, i3 int) bool {
	if i1 > len(firstConsonant_vowel) || i2 >= len(vowel_lastConsonant) {
		return false
	}
	var isVowelValid = false
	var isLastConsonantsValid = false
	if i1 < 0 {
		isVowelValid = true
	} else {
		for _, j := range firstConsonant_vowel[i1] {
			if int(j) == i2 {
				isVowelValid = true
			}
		}
	}
	if i3 < 0 {
		return isVowelValid
	}
	for _, j := range vowel_lastConsonant[i2] {
		if int(j) == i3 {
			isLastConsonantsValid = true
		}
	}
	return isVowelValid && isLastConsonantsValid
}

func isFirstConsonant(str string) bool {
	for _, line := range firstConsonants {
		for _, consonant := range strings.Split(line, " ") {
			if str == consonant {
				return true
			}
		}
	}
	return false
}

func isVowelSound(str string) bool {
	for _, line := range vowels {
		for _, vowel := range strings.Split(line, " ") {
			if str == vowel {
				return true
			}
		}
	}
	return false
}

func isVowelString(str string) bool {
	var isVowels = true
	for _, chr := range []rune(str) {
		if !isVowel(chr) {
			isVowels = false
		}
	}
	return isVowels
}

func isLastConsonant(str string) bool {
	for _, line := range lastConsonants {
		for _, consonant := range strings.Split(line, " ") {
			if str == consonant {
				return true
			}
		}
	}
	return false
}

func attachSound(str string, s Sound) []Sound {
	var sounds []Sound
	for _ = range []rune(str) {
		sounds = append(sounds, s)
	}
	return sounds
}

func init() {
	dictionary = map[string][]Sound{}
	for i1, firstConsonants := range firstConsonants {
		for _, c1 := range strings.Split(firstConsonants, " ") {
			for i2, vowels := range vowels {
				for _, v := range strings.Split(vowels, " ") {
					var sounds = attachSound(v, VowelSound)
					dictionary[v] = sounds
					for i3, lastConsonants := range lastConsonants {
						for _, c2 := range strings.Split(lastConsonants, " ") {
							if isValid(i1, i2, i3) {
								word := c1 + v + c2
								var sounds []Sound
								sounds = append(sounds, attachSound(c1, FirstConsonantSound)...)
								sounds = append(sounds, attachSound(v, VowelSound)...)
								sounds = append(sounds, attachSound(c2, LastConsonantSound)...)
								dictionary[word] = sounds
							}
							if isValid(-1, i2, i3) {
								word := v + c2
								var sounds []Sound
								sounds = append(sounds, attachSound(v, VowelSound)...)
								sounds = append(sounds, attachSound(c2, LastConsonantSound)...)
								dictionary[word] = sounds
							}
							if isValid(i1, i2, -1) {
								word := c1 + v
								var sounds []Sound
								sounds = append(sounds, attachSound(c1, FirstConsonantSound)...)
								sounds = append(sounds, attachSound(v, VowelSound)...)
								dictionary[word] = sounds
							}
						}
					}
				}
			}
		}
	}
}

func LookupDictionary(word string) (bool, []Sound) {
	sounds, found := dictionary[word]
	return found, sounds
}
