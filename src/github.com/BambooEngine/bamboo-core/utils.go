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

var Vowels = []rune("aàáảãạăằắẳẵặâầấẩẫậeèéẻẽẹêềếểễệiìíỉĩịoòóỏõọôồốổỗộơờớởỡợuùúủũụưừứửữựyỳýỷỹỵ")

func IsVowel(chr rune) bool {
	isVowel := false
	for _, v := range Vowels {
		if v == chr {
			isVowel = true
		}
	}
	return isVowel
}

func HasVowel(seq []rune) bool {
	for _, s := range seq {
		if IsVowel(s) {
			return true
		}
	}
	return false
}

func FindVowelPosition(chr rune) int {
	for pos, v := range Vowels {
		if v == chr {
			return pos
		}
	}
	return -1
}

var marksMaps = map[rune]string{
	'a': "aâă__",
	'â': "aâă__",
	'ă': "aâă__",
	'e': "eê___",
	'ê': "eê___",
	'o': "oô_ơ_",
	'ô': "oô_ơ_",
	'ơ': "oô_ơ_",
	'u': "u__ư_",
	'ư': "u__ư_",
	'd': "d___đ",
	'đ': "d___đ",
}

func FindMarkPosition(chr rune) int {
	if str, found := marksMaps[chr]; found {
		for pos, v := range []rune(str) {
			if v == chr {
				return pos
			}
		}
	}
	return -1
}

func FindMarkFromChar(chr rune) (Mark, bool) {
	var pos = FindMarkPosition(chr)
	if pos >= 0 {
		return Mark(pos), true
	}
	return 0, false
}

func RemoveMarkFromChar(chr rune) rune {
	if str, found := marksMaps[chr]; found {
		marks := []rune(str)
		if len(marks) > 0 {
			return marks[0]
		}
	}
	return chr
}

func AddMarkToChar(chr rune, mark uint8) rune {
	var result rune
	tone := FindToneFromChar(chr)
	chr = AddToneToChar(chr, 0)
	if str, found := marksMaps[chr]; found {
		marks := []rune(str)
		if marks[mark] != '_' {
			result = marks[mark]
		}
	}
	result = AddToneToChar(result, uint8(tone))
	return result
}

func IsAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func findIndexRune(chars []rune, r rune) int {
	for i, c := range chars {
		if c == r {
			return i
		}
	}
	return -1
}

func inKeyMap(keys []rune, key rune) bool {
	for _, k := range keys {
		if k == key {
			return true
		}
	}
	return false
}

func FindToneFromChar(chr rune) Tone {
	pos := FindVowelPosition(chr)
	if pos == -1 {
		return TONE_NONE
	}
	return Tone(pos % 6)
}

func AddToneToChar(chr rune, tone uint8) rune {
	pos := FindVowelPosition(chr)
	if pos > -1 {
		current_tone := pos % 6
		offset := int(tone) - current_tone
		return Vowels[pos+offset]
	} else {
		return chr
	}
}

func RemoveToneFromWord(word string) string {
	var chars = []rune(word)
	for i, c := range chars {
		if IsVowel(c) {
			chars[i] = AddToneToChar(c, 0)
		}
	}
	return string(chars)
}
