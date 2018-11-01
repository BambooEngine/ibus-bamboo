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

func FindAppendingRule(rules []Rule, key rune) Rule {
	var result Rule
	result.EffectType = Appending
	result.Key = key
	result.EffectOn = key
	var applicableRules []Rule
	for _, inputRule := range rules {
		if inputRule.Key == key {
			applicableRules = append(applicableRules, inputRule)
		}
	}
	for _, applicableRule := range applicableRules {
		if applicableRule.EffectType == Appending {
			result.EffectOn = applicableRule.EffectOn
			result.AppendedRules = applicableRule.AppendedRules
			return result
		}
	}
	return result
}

func FindLastAppendingTrans(composition []*Transformation) *Transformation {
	for i := len(composition) - 1; i >= 0; i-- {
		var trans = composition[i]
		if trans.Rule.EffectType == Appending {
			return trans
		}
	}
	return nil
}

func FindNextAppendingTransformation(composition []*Transformation, trans *Transformation) (*Transformation, bool) {
	fromIndex := FindTransformationIndex(composition, trans)
	if fromIndex == -1 {
		return nil, false
	}
	var nextAppendingTrans *Transformation
	found := false
	for i := fromIndex + 1; int(i) < len(composition); i++ {
		if composition[i].Rule.EffectType == Appending {
			nextAppendingTrans = composition[i]
			found = true
		}
	}
	return nextAppendingTrans, found
}
