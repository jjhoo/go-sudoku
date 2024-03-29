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

// Permutation internal state data.
type permutation struct {
	ajs    []int
	length int

	visitFlag bool
}

// newPermutation initialize a permutation generator of length elements.
func newPermutation(length int) *permutation {
	p := permutation{visitFlag: true, length: length}

	ajs := make([]int, p.length+2)
	ajs[0] = 0

	tmp := ajs[1:]

	for i := 0; i <= p.length; i++ {
		tmp[i] = i
	}

	p.ajs = ajs

	return &p
}

func (p *permutation) visit() []int {
	n := p.length + 1
	p.visitFlag = false

	return p.ajs[1:n]
}

// Essentially a translation of implemention found in
// https://github.com/jjhoo/sudoku-newlisp/blob/master/sudoku.lsp

// Next() get indexes of next permutation or nil if generator has been exchausted.
func (p *permutation) next() []int {
	if p.visitFlag {
		return p.visit()
	}

	// L2
	j := p.length - 1

	for {
		if p.ajs[j] >= p.ajs[j+1] {
			j--
		} else if p.ajs[j] < p.ajs[j+1] {
			break
		} else if j == 0 {
			break
		}
	}

	if j == 0 {
		p.visitFlag = false
		return nil
	}

	// L3
	l := p.length
	if p.ajs[j] >= p.ajs[l] {
		for {
			l--
			if p.ajs[j] < p.ajs[l] {
				break
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

	p.visitFlag = true
	return p.visit()
}
