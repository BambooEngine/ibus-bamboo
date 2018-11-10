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

func Encode(charsetName string, input string) string {
	if charsetName == "Unicode" {
		return input
	}
	var output string
	if charset, found := charsetDefinitions[charsetName]; found {
		for _, chr := range []rune(input) {
			if out, found := charset[chr]; found {
				output = output + out
			} else {
				output = output + string(chr)
			}
		}
	} else {
		output = input
	}
	return output
}

func GetCharsetNames() []string {
	var names []string
	names = append(names, "Unicode")
	for cs, _ := range charsetDefinitions {
		names = append(names, cs)
	}
	return names
}
