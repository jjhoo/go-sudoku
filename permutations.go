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

type Permutation struct {
	ajs    []int
	length int

	visitFlag bool
	visit     func([]int)
}

func NewPermutation(slice interface{}, visitf func([]int)) Permutation {
	tmp := Permutation{visit: visitf, visitFlag: true}

	xs := reflect.ValueOf(slice)
	tmp.length = xs.Len()

	ajs := make([]int, tmp.length+1)

	for i := 0; i <= tmp.length; i++ {
		ajs[i] = i
	}
	ajs = append([]int{0}, ajs...)
	tmp.ajs = ajs

	return tmp
}

func (p *Permutation) Visit() {
	n := p.length + 1
	p.visit(p.ajs[1:n])
	p.visitFlag = false
}

// Essentially a translation of implemention found in
// https://github.com/jjhoo/sudoku-newlisp/blob/master/sudoku.lsp
func (p *Permutation) Next() bool {
	if p.visitFlag {
		return true
	}

	// L2
	j := p.length - 1
	cont := true

	for cont {
		if p.ajs[j] >= p.ajs[j+1] {
			j--
		} else if p.ajs[j] < p.ajs[j+1] {
			cont = false
		} else if j == 0 {
			cont = false
		}
	}

	if j == 0 {
		return false
	}

	// L3
	l := p.length
	if p.ajs[j] >= p.ajs[l] {
		cont = true
		for cont {
			l--
			if p.ajs[j] < p.ajs[l] {
				cont = false
			}
		}
	}
	p.ajs[j], p.ajs[l] = p.ajs[l], p.ajs[j]

	// L4
	k := j + 1
	l = p.length

	for k < l {
		p.ajs[k], p.ajs[l] = p.ajs[l], p.ajs[k]
		k++
		l--
	}
	return true
}
