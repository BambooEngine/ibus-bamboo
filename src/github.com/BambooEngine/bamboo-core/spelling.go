/*
 * Bamboo - A Vietnamese Input method editor
 * Copyright (C) Luong Thanh Lam <ltlam93@gmail.com>
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
 * A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
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

var spellingTrie = &W{F: false}

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
		AddTrie(spellingTrie, []rune(word), false)
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

func getLastWord(composition []*Transformation, effectiveKeys []rune) []*Transformation {
	for i := len(composition) - 1; i >= 0; i-- {
		var t = composition[i]
		if t.Rule.EffectType == Appending && !unicode.IsLetter(t.Rule.EffectOn) && !inKeyList(effectiveKeys, t.Rule.EffectOn) {
			if i == len(composition)-1 {
				return nil
			}
			return composition[i+1:]
		}
	}
	return composition
}

func getLastSyllable(composition []*Transformation) []*Transformation {
	var ret []*Transformation
	if len(composition) <= 1 {
		return composition
	}
	for i, trans := range composition {
		ret = append(ret, trans)
		if i < len(composition)-1 && composition[i+1].Rule.EffectType != Appending {
			continue
		}
		str := Flatten(ret, VietnameseMode|ToneLess|LowerCase)
		if str == "" {
			continue
		}
		if FindWord(spellingTrie, []rune(str), false) == FindResultNotMatch {
			if i == 0 {
				return getLastSyllable(composition[1:])
			}
			return getLastSyllable(composition[i:])
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
