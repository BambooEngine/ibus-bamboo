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
			if mode&MarkLess != 0 && (chr < 'a' || chr > 'z') {
				chr = RemoveMarkFromChar(chr)
			}
			for _, trans := range transList {
				switch trans.Rule.EffectType {
				case MarkTransformation:
					if mode&MarkLess != 0 {
						break
					}
					if trans.Rule.Effect == uint8(MARK_RAW) {
						chr = appendingTrans.Rule.Key
					} else {
						chr = AddMarkToChar(chr, trans.Rule.Effect)
					}
				case ToneTransformation:
					if mode&ToneLess != 0 {
						break
					}
					chr = AddToneToChar(chr, trans.Rule.Effect)
				}
			}
		}
		if mode&ToneLess != 0 {
			chr = AddToneToChar(chr, 0)
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
