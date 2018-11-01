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

type inputMethodDefinition = map[rune]string

// todo: move to input_method.json
var inputMethodDefinitions = map[string]inputMethodDefinition{
	"Telex": {
		'z': "XoaDauThanh",
		's': "DauSac",
		'f': "DauHuyen",
		'r': "DauHoi",
		'x': "DauNga",
		'j': "DauNang",
		'a': "A_Â",
		'e': "E_Ê",
		'o': "O_Ô",
		'w': "UOA_ƯƠĂ",
		'd': "D_Đ",
	},
	"VNI": {
		'0': "XoaDauThanh",
		'1': "DauSac",
		'2': "DauHuyen",
		'3': "DauHoi",
		'4': "DauNga",
		'5': "DauNang",
		'6': "AEO_ÂÊÔ",
		'7': "UO_ƯƠ",
		'8': "A_Ă",
		'9': "D_Đ",
	},

	"VIQR": {
		'0':  "XoaDauThanh",
		'\'': "DauSac",
		'`':  "DauHuyen",
		'?':  "DauHoi",
		'~':  "DauNga",
		'.':  "DauNang",
		'^':  "AEO_ÂÊÔ",
		'+':  "UO_ƯƠ",
		'*':  "UO_ƯƠ",
		'(':  "A_Ă",
		'\\': "D_Đ",
	},
	"Microsoft layout": {
		'8': "DauSac",
		'5': "DauHuyen",
		'6': "DauHoi",
		'7': "DauNga",
		'9': "DauNang",
		'1': "__ă",
		'!': "_Ă",
		'2': "__â",
		'@': "_Â",
		'3': "__ê",
		'#': "_Ê",
		'4': "__ô",
		'$': "_Ô",
		'0': "__đ",
		')': "_Đ",
		'[': "__ư",
		'{': "_Ư",
		']': "__ơ",
		'}': "_Ơ",
	},

	"Telex 2": {
		'z': "XoaDauThanh",
		's': "DauSac",
		'f': "DauHuyen",
		'r': "DauHoi",
		'x': "DauNga",
		'j': "DauNang",
		'a': "A_Â",
		'e': "E_Ê",
		'o': "O_Ô",
		'w': "UOA_ƯƠĂ__Ư",
		'd': "D_Đ",
		']': "__ư",
		'[': "__ơ",
		'}': "_Ư",
		'{': "_Ơ",
	},
	"Telex + VNI + VIQR": {
		'z':  "XoaDauThanh",
		's':  "DauSac",
		'f':  "DauHuyen",
		'r':  "DauHoi",
		'x':  "DauNga",
		'j':  "DauNang",
		'a':  "A_Â",
		'e':  "E_Ê",
		'o':  "O_Ô",
		'w':  "UOA_ƯƠĂ",
		'd':  "D_Đ",
		'0':  "XoaDauThanh",
		'1':  "DauSac",
		'2':  "DauHuyen",
		'3':  "DauHoi",
		'4':  "DauNga",
		'5':  "DauNang",
		'6':  "AEO_ÂÊÔ",
		'7':  "UO_ƯƠ",
		'8':  "A_Ă",
		'9':  "D_Đ",
		'\'': "DauSac",
		'`':  "DauHuyen",
		'?':  "DauHoi",
		'~':  "DauNga",
		'.':  "DauNang",
		'^':  "AEO_ÂÊÔ",
		'+':  "UO_ƯƠ",
		'*':  "UO_ƯƠ",
		'(':  "A_Ă",
		'\\': "D_Đ",
	},
	"VNI Bàn phím tiếng Pháp": {
		'&':  "XoaDauThanh",
		'é':  "DauSac",
		'"':  "DauHuyen",
		'\'': "DauHoi",
		'(':  "DauNga",
		'-':  "DauNang",
		'è':  "AEO_ÂÊÔ",
		'_':  "UO_ƯƠ",
		'ç':  "A_Ă",
		'à':  "D_Đ",
	},
	"Telex 3": {
		'z': "XoaDauThanh",
		's': "DauSac",
		'f': "DauHuyen",
		'r': "DauHoi",
		'x': "DauNga",
		'j': "DauNang",
		'a': "A_Â",
		'e': "E_Ê",
		'o': "O_Ô",
		'w': "UOA_ƯƠĂ",
		'd': "D_Đ",
		'[': "__ươ",
		'{': "_ƯƠ",
	},
}
