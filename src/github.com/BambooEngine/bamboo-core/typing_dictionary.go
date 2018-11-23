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

var typingFirstConsonantSeq = [3]string{
	"b d đ g gh m n nh p ph r s t tr v z",
	"c h k kh qu th",
	"ch gi l ng ngh x",
}

var typingVowelSeq = [6]string{
	"e ê i ua ue uê uy y",
	"a ie iê oa uye uyê ye yê",
	"a â ă e o oo ô ơ oe u ư ua uâ uo uô ưo uơ ươ",
	"oa oă",
	"uo uơ",
	"ai ao au âu ay ây eo eu êu ia ieu iêu iu oai oao oay oeo oi ôi ơi ua ưa uay uây ui ưi uoi uơi ưoi uôi ươi uou ưou uơu ươu uu ưu uya uyu yeu yêu",
}

var typingLastConsonantSeq = [3]string{
	"ch nh",
	"c ng",
	"m n p t",
}

var typingCVMatrix = [3][]uint{
	{0, 1, 2, 3, 5},
	{0, 1, 2, 3, 5},
	{0, 1, 2, 3, 5},
}

var typingVCMatrix = [6][]uint{
	{0, 2},
	{0, 1, 2},
	{1, 2},
	{1, 2},
	// ignore
}

var typingDictionary map[string]bool

func isValidTypingCVC(i1, i2, i3 int) bool {
	if i1 > len(typingCVMatrix) || i2 >= len(typingVCMatrix) {
		return false
	}
	var isVowelValid = false
	var isLastConsonantsValid = false
	if i1 < 0 {
		isVowelValid = true
	} else {
		for _, j := range typingCVMatrix[i1] {
			if int(j) == i2 {
				isVowelValid = true
			}
		}
	}
	if i3 < 0 {
		return isVowelValid
	}
	for _, j := range typingVCMatrix[i2] {
		if int(j) == i3 {
			isLastConsonantsValid = true
		}
	}
	return isVowelValid && isLastConsonantsValid
}

func init() {
	typingDictionary = map[string]bool{}
	for i1, firstConsonants := range typingFirstConsonantSeq {
		for _, c1 := range strings.Split(firstConsonants, " ") {
			for i2, vowels := range typingVowelSeq {
				for _, v := range strings.Split(vowels, " ") {
					for i3, lastConsonants := range typingLastConsonantSeq {
						for _, c2 := range strings.Split(lastConsonants, " ") {
							if isValidTypingCVC(i1, i2, i3) {
								word := c1 + v + c2
								typingDictionary[word] = true
							}
							if isValidTypingCVC(-1, i2, i3) {
								word := v + c2
								typingDictionary[word] = true
							}
							if isValidTypingCVC(i1, i2, -1) {
								word := c1 + v
								typingDictionary[word] = true
							}
						}
					}
				}
			}
		}
	}
}

func LookupTypingDictionary(word string) bool {
	if isTypingFirstConsonantSeq(word) || isTypingVowelSeq(word) {
		return true
	}
	_, found := typingDictionary[word]
	if found {
		return true
	}
	for w := range typingDictionary {
		if strings.Contains(w, word) {
			return true
		}
	}
	return false
}
