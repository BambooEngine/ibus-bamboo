/*
 * Bamboo - A Vietnamese Input method editor
 * Copyright (C) Luong Thanh Lam <ltlam93@gmail.com>
 *
 * This software is licensed under the MIT license. For more information,
 * see <https://github.com/BambooEngine/bamboo-core/blob/master/LICENCE>.
 */
package bamboo

import (
	"unicode"
)

func Flatten(composition []*Transformation, mode Mode) string {
	return string(getCanvas(composition, mode))
}

func getCanvas(composition []*Transformation, mode Mode) []rune {
	var canvas []rune
	var appendingMap = map[*Transformation][]*Transformation{}
	var appendingList []*Transformation
	for _, trans := range composition {
		if mode&EnglishMode != 0 {
			if trans.Rule.Key == 0 {
				// ignore virtual key
				continue
			}
			appendingList = append(appendingList, trans)
		} else if trans.Rule.EffectType == Appending {
			appendingList = append(appendingList, trans)
		} else if trans.Target != nil {
			appendingMap[trans.Target] = append(appendingMap[trans.Target], trans)
		}
	}
	for _, appendingTrans := range appendingList {
		var chr rune
		var transList = appendingMap[appendingTrans]
		if mode&EnglishMode != 0 {
			chr = appendingTrans.Rule.Key
		} else {
			chr = appendingTrans.Rule.EffectOn
			for _, trans := range transList {
				switch trans.Rule.EffectType {
				case MarkTransformation:
					if trans.Rule.Effect == uint8(MARK_RAW) {
						chr = appendingTrans.Rule.Key
					} else {
						chr = AddMarkToChar(chr, trans.Rule.Effect)
					}
				case ToneTransformation:
					chr = AddToneToChar(chr, trans.Rule.Effect)
				}
			}
		}
		if mode&ToneLess != 0 {
			chr = AddToneToChar(chr, 0)
		}
		if mode&MarkLess != 0 {
			chr = AddMarkToChar(chr, 0)
		}
		if mode&LowerCase != 0 {
			chr = unicode.ToLower(chr)
		} else if appendingTrans.IsUpperCase {
			chr = unicode.ToUpper(chr)
		}
		canvas = append(canvas, chr)
	}
	return canvas
}
