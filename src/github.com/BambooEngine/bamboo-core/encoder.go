/*
 * Bamboo - A Vietnamese Input method editor
 * Copyright (C) Luong Thanh Lam <ltlam93@gmail.com>
 *
 * This software is licensed under the MIT license. For more information,
 * see <https://github.com/BambooEngine/bamboo-core/blob/master/LICENSE>.
 */

package bamboo

const UNICODE = "Unicode"

func Encode(charsetName string, input string) string {
	if charsetName == UNICODE {
		return input
	}
	var output string
	if charset, found := charsetDefinitions[charsetName]; found {
		for _, chr := range input {
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
	names = append(names, UNICODE)
	for cs := range charsetDefinitions {
		names = append(names, cs)
	}
	return names
}
