/*
 * Bamboo - A Vietnamese Input method editor
 * Copyright (C) Luong Thanh Lam <ltlam93@gmail.com>
 *
 * This software is licensed under the MIT license. For more information,
 * see <https://github.com/BambooEngine/bamboo-core/blob/master/LISENCE>.
 */
package bamboo

import "log"

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
	for i, row := range seq {
		var canvas []rune
		var rowLen = len([]rune(row))
		for ii, char := range []rune(row) {
			if char != ' ' {
				canvas = append(canvas, char)
				if ii < rowLen-1 {
					continue
				}
			}
			var canvasLen = len(canvas)
			if canvasLen < inputLen || (inputIsFull && canvasLen > inputLen) {
				canvas = nil
				continue
			}
			var isMatch = true
			for j, ic := range []rune(input) {
				if ic != canvas[j] && !(!inputIsComplete && AddMarkToTonelessChar(canvas[j], 0) == ic) {
					isMatch = false
				}
			}
			canvas = nil
			if isMatch {
				ret = append(ret, i)
				break
			}
		}
	}
	return ret
}

func isValidCVC(fc, vo, lc string, inputIsFullComplete bool) bool {
	var ret bool
	defer func() {
		return
		log.Printf("fc=%s vo=%s lc=%s ret=%v", fc, vo, lc, ret)
	}()
	var fcIndexes, voIndexes, lcIndexes []int
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
	// first consonant only
	if voIndexes == nil {
		return fcIndexes != nil
	}
	// first consonant + vowel
	if fcIndexes != nil {
		for _, fci := range fcIndexes {
			for _, voi := range voIndexes {
				if isValidCV(fci, voi) {
					ret = true
				}
			}
		}
		if ret == false || lcIndexes == nil {
			return ret
		}
	}
	// vowel + last consonant
	if lcIndexes != nil {
		ret = false
		for _, voi := range voIndexes {
			for _, lci := range lcIndexes {
				if isValidVC(voi, lci) {
					ret = true
				}
			}
		}
	} else {
		// vowel only
		ret = true
	}
	return ret
}

func isValidCV(fc, vo int) bool {
	for _, v := range cvMatrix[fc] {
		if v == vo {
			return true
		}
	}
	return false
}

func isValidVC(vo, lc int) bool {
	for _, t := range vcMatrix[vo] {
		if t == lc {
			return true
		}
	}
	return false
}
