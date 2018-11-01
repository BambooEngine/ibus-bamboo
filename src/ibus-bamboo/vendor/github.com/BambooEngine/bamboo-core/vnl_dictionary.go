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

import "strings"

var vnl_firstConsonants = [3]string{
	"b d đ g gh m n nh p ph r s t tr v z",
	"c h k kh qu th",
	"ch gi l ng ngh x",
}

var vnl_vowels = [6]string{
	"e ê i ua ue uê uy y",
	"a ie iê oa uye uyê ye yê",
	"a â ă e o oo ô ơ oe u ư ua uâ uo uô ưo uơ ươ",
	"oa oă",
	"uo uơ",
	"ai ao au âu ay ây eo eu êu ia ieu iêu iu oai oao oay oeo oi ôi ơi ua ưa uay uây ui ưi uoi uôi ươi uou ưou uơu ươu uu ưu uya uyu yeu yêu",
}

var vnl_lastConsonants = [3]string{
	"ch nh",
	"c ng",
	"m n p t",
}

var vnl_firstConsonant_vowel = [3][]uint{
	{0, 1, 2, 3, 5},
	{0, 1, 2, 3, 5},
	{0, 1, 2, 3, 5},
}

var vnl_vowel_lastConsonant = [6][]uint{
	{0, 2},
	{0, 1, 2},
	{1, 2},
	{2},
	// ignore
}

var vnl_dictionary map[string]bool

func vnl_isValid(i1, i2, i3 int) bool {
	if i1 > len(vnl_firstConsonant_vowel) || i2 >= len(vnl_vowel_lastConsonant) {
		return false
	}
	var isVowelValid = false
	var isLastConsonantsValid = false
	if i1 < 0 {
		isVowelValid = true
	} else {
		for _, j := range vnl_firstConsonant_vowel[i1] {
			if int(j) == i2 {
				isVowelValid = true
			}
		}
	}
	if i3 < 0 {
		return isVowelValid
	}
	for _, j := range vnl_vowel_lastConsonant[i2] {
		if int(j) == i3 {
			isLastConsonantsValid = true
		}
	}
	return isVowelValid && isLastConsonantsValid
}

func init() {
	vnl_dictionary = map[string]bool{}
	for i1, firstConsonants := range vnl_firstConsonants {
		for _, c1 := range strings.Split(firstConsonants, " ") {
			for i2, vowels := range vnl_vowels {
				for _, v := range strings.Split(vowels, " ") {
					for i3, lastConsonants := range vnl_lastConsonants {
						for _, c2 := range strings.Split(lastConsonants, " ") {
							if vnl_isValid(i1, i2, i3) {
								word := c1 + v + c2
								vnl_dictionary[word] = true
							}
							if vnl_isValid(-1, i2, i3) {
								word := v + c2
								vnl_dictionary[word] = true
							}
							if vnl_isValid(i1, i2, -1) {
								word := c1 + v
								vnl_dictionary[word] = true
							}
						}
					}
				}
			}
		}
	}
}

func isVnlFirstConsonant(str string) bool {
	for _, line := range vnl_firstConsonants {
		for _, consonant := range strings.Split(line, " ") {
			if str == consonant {
				return true
			}
		}
	}
	return false
}

func isVnlVowel(str string) bool {
	for _, line := range vnl_vowels {
		for _, vowel := range strings.Split(line, " ") {
			if str == vowel {
				return true
			}
		}
	}
	return false
}

func isVnlLastConsonant(str string) bool {
	for _, line := range vnl_lastConsonants {
		for _, consonant := range strings.Split(line, " ") {
			if str == consonant {
				return true
			}
		}
	}
	return false
}

func LookupVnlDictionary(word string) bool {
	if isVnlFirstConsonant(word) || isVnlVowel(word) {
		return true
	}
	_, found := vnl_dictionary[word]
	if found {
		return true
	}
	for w := range vnl_dictionary {
		if strings.Contains(w, word) {
			return true
		}
	}
	return false
}
