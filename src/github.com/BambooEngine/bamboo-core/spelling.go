/*
 * Bamboo - A Vietnamese Input method editor
 * Copyright (C) Luong Thanh Lam <ltlam93@gmail.com>
 *
 * This software is licensed under the MIT license. For more information,
 * see <https://github.com/BambooEngine/bamboo-core/blob/master/LICENCE>.
 */

package bamboo

var firstConsonantSeqs = []string{
	"b d đ g gh m n nh p ph r s t tr v z",
	"c h k kh qu th",
	"ch gi l ng ngh x",
	"đ l",
}

var vowelSeqs = []string{
	"ê i ua uê uy y",
	"a iê oa uyê yê",
	"â ă e o oo ô ơ oe u ư uâ uô ươ",
	"oă",
	"uơ",
	"ai ao au âu ay ây eo êu ia iêu iu oai oao oay oeo oi ôi ơi ưa uây ui ưi uôi ươi ươu ưu uya uyu yêu",
	"ă",
}

var lastConsonantSeqs = []string{
	"ch nh",
	"c ng",
	"m n p t",
	"k",
}

var cvMatrix = [4][]int{
	{0, 1, 2, 5},
	{0, 1, 2, 3, 4, 5},
	{0, 1, 2, 3, 5},
	{6},
}

var vcMatrix = [7][]int{
	{0, 2},
	{0, 1, 2},
	{1, 2},
	{1, 2},
	{},
	{},
	{3},
}

func lookup(seq []string, input string, inputIsFull, inputIsComplete bool) []int {
	var ret []int
	var inputLen = len([]rune(input))
	for index, row := range seq {
		var i = 0
		var rows = append([]rune(row), ' ')
		for j, char := range rows {
			if char != ' ' {
				continue
			}
			var canvas = rows[i:j]
			i = j + 1
			if len(canvas) < inputLen || (inputIsFull && len(canvas) > inputLen) {
				continue
			}
			var isMatch = true
			for k, ic := range []rune(input) {
				if ic != canvas[k] && !(!inputIsComplete && AddMarkToTonelessChar(canvas[k], 0) == ic) {
					isMatch = false
					break
				}
			}
			if isMatch {
				ret = append(ret, index)
				break
			}
		}
	}
	return ret
}

func isValidCVC(fc, vo, lc string, inputIsFullComplete bool) bool {
	var ret bool
	var fcIndexes, voIndexes, lcIndexes []int
	// log.Printf("fc=%s vo=%s lc=%s ret=%v", fc, vo, lc, ret)
	if fc != "" {
		if fcIndexes = lookup(firstConsonantSeqs, fc, inputIsFullComplete || vo != "", true); fcIndexes == nil {
			return false
		}
	}
	if vo != "" {
		if voIndexes = lookup(vowelSeqs, vo, inputIsFullComplete || lc != "", inputIsFullComplete); voIndexes == nil {
			return false
		}
	}
	if lc != "" {
		if lcIndexes = lookup(lastConsonantSeqs, lc, inputIsFullComplete, true); lcIndexes == nil {
			return false
		}
	}
	if voIndexes == nil {
		// first consonant only
		return fcIndexes != nil
	}
	if fcIndexes != nil {
		// first consonant + vowel
		if ret = isValidCV(fcIndexes, voIndexes); !ret || lcIndexes == nil {
			return ret
		}
	}
	if lcIndexes != nil {
		// vowel + last consonant
		ret = isValidVC(voIndexes, lcIndexes)
	} else {
		// vowel only
		ret = true
	}
	return ret
}

func isValidCV(fcIndexes, voIndexes []int) bool {
	for _, fc := range fcIndexes {
		for _, c := range cvMatrix[fc] {
			for _, vo := range voIndexes {
				if c == vo {
					return true
				}
			}
		}
	}
	return false
}

func isValidVC(voIndexes, lcIndexes []int) bool {
	for _, vo := range voIndexes {
		for _, c := range vcMatrix[vo] {
			for _, lc := range lcIndexes {
				if c == lc {
					return true
				}
			}
		}
	}
	return false
}
