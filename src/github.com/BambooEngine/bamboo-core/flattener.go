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
			if mode&MarkLess != 0 && (effectOn < 'a' || effectOn > 'z') {
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
			if mode&MarkLess != 0 {
				break
			}
			apply_effect(AddMarkToChar, trans)
			break
		case ToneTransformation:
			if mode&ToneLess != 0 {
				break
			}
			apply_effect(AddToneToChar, trans)
			break
		}
	}
	return canvas
}
