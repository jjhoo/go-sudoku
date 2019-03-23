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
package sudoku

import (
	_ "fmt"
)

type Box struct {
	Row    int8
	Column int8
}

type Pos struct {
	Row    int8
	Column int8
	Box    Box
}

func numToBoxNumber(n int8) int8 {
	var nn int8

	switch n {
	case 1, 2, 3:
		nn = 1
	case 4, 5, 6:
		nn = 2
	case 7, 8, 9:
		nn = 3
	}
	return nn
}

func numToBox(n int8) Box {
	var box = Box{}
	switch n {
	case 1:
		box = Box{Row: 1, Column: 1}
	case 2:
		box = Box{Row: 1, Column: 2}
	case 3:
		box = Box{Row: 1, Column: 3}
	case 4:
		box = Box{Row: 2, Column: 1}
	case 5:
		box = Box{Row: 2, Column: 2}
	case 6:
		box = Box{Row: 2, Column: 3}
	case 7:
		box = Box{Row: 3, Column: 1}
	case 8:
		box = Box{Row: 3, Column: 2}
	case 9:
		box = Box{Row: 3, Column: 3}
	}
	return box
}

func (b Box) init(row int8, col int8) Box {
	b.Row = numToBoxNumber(row)
	b.Column = numToBoxNumber(col)

	return b
}

func (p Pos) init(row int8, col int8) Pos {
	p.Row = row
	p.Column = col
	p.Box = Box{}.init(row, col)

	return p
}

func (p Pos) eqRow(others ...Pos) bool {
	for _, other := range others {
		if p.Row != other.Row {
			return false
		}
	}
	return true
}

func (p Pos) eqColumn(others ...Pos) bool {
	for _, other := range others {
		if p.Column != other.Column {
			return false
		}
	}
	return true
}

func (p Pos) eqBox(others ...Pos) bool {
	for _, other := range others {
		if p.Box != other.Box {
			return false
		}
	}
	return true
}

func (p Pos) sees(other Pos) bool {
	return p.eqRow(other) || p.eqColumn(other) || p.eqBox(other)
}

func (p Pos) less(other *Pos) bool {
	if p.Row < other.Row {
		return true
	}

	if p.Row == other.Row {
		return p.Column < other.Column
	}

	return false
}
