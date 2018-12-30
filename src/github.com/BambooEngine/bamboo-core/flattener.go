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
	"log"
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
	if mode&LowerCase != 0 {
		return string(canvas)
	}
	return f.toUpper(composition, canvas, mode)
}

func (f *BambooFlattener) toUpper(composition []*Transformation, canvas []rune, mode Mode) string {
	if mode&VietnameseMode != 0 {
		for _, trans := range composition {
			if trans.Rule.EffectType == Appending {
				if int(trans.Dest) >= len(canvas) {
					log.Println("Something is wrong with dest of trans")
					continue
				}
				if trans.IsUpperCase {
					canvas[trans.Dest] = unicode.ToUpper(canvas[trans.Dest])
				}
			}
		}
		return string(canvas)
	}
	for _, trans := range composition {
		if int(trans.Dest) >= len(canvas) {
			log.Println("Something is wrong with dest of trans")
			continue
		}
		if trans.IsUpperCase {
			canvas[trans.Dest] = unicode.ToUpper(canvas[trans.Dest])
		}
	}
	return string(canvas)
}

func (f *BambooFlattener) GetCanvas(composition []*Transformation, mode Mode) []rune {
	var canvas []rune
	apply_effect := func(callback func(rune, uint8) rune, trans *Transformation) {
		if trans.Target == nil || len(canvas) <= int(trans.Target.Dest) {
			//log.Println("There's something wrong with canvas [nhoawfng]")
			return
		}
		index := trans.Target.Dest
		canvas[index] = callback(canvas[index], trans.Rule.Effect)
	}
	for _, trans := range composition {
		trans.Dest = 0
		if trans.IsDeleted {
			continue
		}
		if mode&EnglishMode != 0 {
			if trans.Rule.Key > 0 {
				trans.Dest = uint(len(canvas))
				canvas = append(canvas, trans.Rule.Key)
			}
			// ignore virtual key
			continue
		}
		if trans.Rule.EffectType == Appending {
			trans.Dest = uint(len(canvas))
			var effectOn = trans.Rule.EffectOn
			if mode&NoMark != 0 && (effectOn < 'a' || effectOn > 'z') {
				effectOn = RemoveMarkFromChar(effectOn)
			}
			canvas = append(canvas, effectOn)
		}
	}
	if mode&EnglishMode != 0 || len(canvas) == 0 {
		return canvas
	}
	for _, trans := range composition {
		if trans.IsDeleted || trans.Target == nil {
			continue
		}
		switch trans.Rule.EffectType {
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
