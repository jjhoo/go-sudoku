// Copyright (c) 2019-2022 Jani J. Hakala <jjhakala@gmail.com>, Finland
//
//  This program is free software: you can redistribute it and/or modify
//  it under the terms of the GNU Affero General Public License as
//  published by the Free Software Foundation, version 3 of the
//  License.
//
//  This program is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//  GNU Affero General Public License for more details.
//
//  You should have received a copy of the GNU Affero General Public License
//  along with this program.  If not, see <http://www.gnu.org/licenses/>.
//

package sudoku

import (
	_ "fmt"
	"sort"
)

func dedupeInt8(ns []int8) int8List {
	if len(ns) <= 1 {
		return ns
	}

	prev := ns[0]
	res := []int8{prev}

	for _, n := range ns[1:] {
		if n != prev {
			res = append(res, n)
			prev = n
		}
	}

	return res
}

func sortInt8(array []int8) {
	sort.Slice(array, func(i, j int) bool {
		return array[i] < array[j]
	})
}

type numCount struct {
	num   int8
	count int8
}

func numberCounts(nums int8List) numCountList {
	nnums := make(int8List, len(nums))
	copy(nnums, nums)

	sortInt8(nnums)
	unums := dedupeInt8(nnums)

	counts := unums.MapNumCount(func(n int8) numCount {
		count := nums.Reduce(0, func(acc, nn int8) int8 {
			if n == nn {
				acc++
			}
			return acc
		})
		return numCount{num: n, count: count}
	})

	return counts
}
