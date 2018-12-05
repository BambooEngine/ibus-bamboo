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
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package bamboo

var downLvlMap = map[rune][]rune{
	'á': {'a'},
	'à': {'a'},
	'ả': {'a'},
	'ã': {'a'},
	'ạ': {'a'},

	'ă': {'a'},
	'ắ': {'ă', 'á', 'a'},
	'ằ': {'ă', 'à', 'a'},
	'ẳ': {'ă', 'ả', 'a'},
	'ẵ': {'ă', 'ã', 'a'},
	'ặ': {'ă', 'ạ', 'a'},

	'â': {'a'},
	'ấ': {'â', 'á', 'a'},
	'ầ': {'â', 'à', 'a'},
	'ẩ': {'â', 'ả', 'a'},
	'ẫ': {'â', 'ã', 'a'},
	'ậ': {'â', 'ạ', 'a'},

	'é': {'e'},
	'è': {'e'},
	'ẻ': {'e'},
	'ẽ': {'e'},
	'ẹ': {'e'},

	'ê': {'e'},
	'ế': {'ê', 'é', 'e'},
	'ề': {'ê', 'è', 'e'},
	'ể': {'ê', 'ẻ', 'e'},
	'ễ': {'ê', 'ẽ', 'e'},
	'ệ': {'ê', 'ẹ', 'e'},

	'í': {'i'},
	'ì': {'i'},
	'ỉ': {'i'},
	'ĩ': {'i'},
	'ị': {'i'},

	'ó': {'o'},
	'ò': {'o'},
	'ỏ': {'o'},
	'õ': {'o'},
	'ọ': {'o'},

	'ô': {'o'},
	'ố': {'ô', 'ó', 'o'},
	'ồ': {'ô', 'ò', 'o'},
	'ổ': {'ô', 'ỏ', 'o'},
	'ỗ': {'ô', 'õ', 'o'},
	'ộ': {'ô', 'ọ', 'o'},

	'ơ': {'o'},
	'ớ': {'ơ', 'ó', 'o'},
	'ờ': {'ơ', 'ò', 'o'},
	'ở': {'ơ', 'ỏ', 'o'},
	'ỡ': {'ơ', 'õ', 'o'},
	'ợ': {'ơ', 'ọ', 'o'},

	'ú': {'u'},
	'ù': {'u'},
	'ủ': {'u'},
	'ũ': {'u'},
	'ụ': {'u'},

	'ư': {'u'},
	'ứ': {'ư', 'ú', 'u'},
	'ừ': {'ư', 'ù', 'u'},
	'ử': {'ư', 'ủ', 'u'},
	'ữ': {'ư', 'ũ', 'u'},
	'ự': {'ư', 'ụ', 'u'},

	'ý': {'y'},
	'ỳ': {'y'},
	'ỷ': {'y'},
	'ỹ': {'y'},
	'ỵ': {'y'},

	'đ': {'d'},
}
