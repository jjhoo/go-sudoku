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
	"reflect"
)

type combination struct {
	cjs    []int
	length int
	koo    int
	j      int
	k      int

	visitFlag bool
}

// Knuth, algorithm T
func Combination(slice interface{}, koo int) combination {
	tmp := combination{visitFlag: true, koo: koo, j: koo, k: koo}

	xs := reflect.ValueOf(slice)
	tmp.length = xs.Len()

	cjs := make([]int, koo+1)

	for i := 0; i <= koo; i++ {
		cjs[i] = i - 1
	}

	cjs = append(cjs, tmp.length, 0)
	tmp.cjs = cjs

	return tmp
}

func (c *combination) visit() []int {
	n := c.k + 1
	c.visitFlag = false

	return c.cjs[1:n]
}

// Essentially a translation of implemention found in
// https://github.com/jjhoo/sudoku-newlisp/blob/master/sudoku.lsp
func (c *combination) Next() []int {
	if c.visitFlag {
		return c.visit()
	}

	if c.j > 0 {
		//  T6
		x := c.j
		c.cjs[c.j] = x
		c.j--

		c.visitFlag = true
		return c.visit()
	}

	// T3
	if (c.cjs[1] + 1) < c.cjs[2] {
		c.cjs[1]++

		c.visitFlag = true
		return c.visit()
	}

	// T4
	c.j = 2
	cont := true
	x := -1

	for cont {
		c.cjs[c.j-1] = c.j - 2
		x = c.cjs[c.j] + 1

		if x == c.cjs[c.j+1] {
			c.j++
		} else {
			cont = false
		}
	}

	// T5
	if c.j > c.k {
		return nil
	}

	// T6
	c.cjs[c.j] = x
	c.j--

	c.visitFlag = true
	return c.visit()
}
