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
	"unicode"
)

type Flattener interface {
	Flatten([]*Transformation, Mode) string
}

type BambooFlattener struct {
}

func Flatten(composition []*Transformation, mode Mode) string {
	var flattener Flattener = new(BambooFlattener)
	if mode&LowerCase != 0 {
		return strings.ToLower(flattener.Flatten(composition, mode))
	}
	return flattener.Flatten(composition, mode)
}

func (f *BambooFlattener) Flatten(composition []*Transformation, mode Mode) string {
	canvas := f.GetCanvas(composition, mode)
	for _, trans := range composition {
		if trans.Rule.EffectType == Appending {
			if trans.IsUpperCase {
				canvas[trans.Dest] = unicode.ToUpper(canvas[trans.Dest])
			}
		}
	}
	return string(canvas)
}

func (f *BambooFlattener) GetCanvas(composition []*Transformation, mode Mode) []rune {
	var canvas []rune
	apply_effect := func(callback func(rune, Mark) rune, trans *Transformation) {
		if trans.Target == nil || len(canvas) <= int(trans.Target.Dest) {
			return
		}
		index := trans.Target.Dest
		charWithEffect := callback(canvas[index], trans.Rule.Effect)
		// Double typing an affect key undoes it. Btw, we're playing
		// fast-and-loose here by replying on the fact that TONE_NONE equals
		// MARK_NONE and equals 0.
		if charWithEffect == canvas[index] {
			canvas[index] = callback(canvas[index], TONE_NONE)
		} else {
			canvas[index] = charWithEffect
		}
	}
	for _, trans := range composition {
		if mode&EnglishMode != 0 {
			if trans.Rule.Key > 0 {
				canvas = append(canvas, trans.Rule.Key)
			}
			// ignore virtual key
			continue
		}
		switch trans.Rule.EffectType {
		case Appending:
			trans.Dest = uint(len(canvas))
			var effectOn = trans.Rule.EffectOn
			if mode&NoMark != 0 && (effectOn < 'a' || effectOn > 'z') {
				effectOn = RemoveMarkFromChar(effectOn)
			}
			canvas = append(canvas, effectOn)
			break
		case MarkTransformation:
			if mode&NoMark != 0 {
				break
			}
			apply_effect(AddMarkToChar, trans)
			break
		case ToneTransformation:
			if mode&NoTone != 0 {
				break
			}
			apply_effect(AddToneToChar, trans)
			break
		}
	}
	return canvas
}

