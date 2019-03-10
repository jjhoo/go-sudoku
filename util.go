// Copyright (c) 2019 Jani J. Hakala <jjhakala@gmail.com> Jyväskylä, Finland
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
package main

import (
	_ "fmt"
	"sort"
)

func dedupeInt8(ns []int8) []int8 {
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

func numberCounts(nums []int8) []numCount {
	nnums := dedupeInt8(nums)

	counts := []numCount{}

	for _, n := range nnums {
		var count int8 = 0
		for _, nn := range nums {
			if n == nn {
				count++
			}
		}
		counts = append(counts, numCount{num: n, count: count})
	}

	return counts
}

