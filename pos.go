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
)

type Pos struct {
	Row    int8
	Column int8
	Box    int8
}

func (p Pos) init(row int8, col int8) Pos {
	p.Row = row
	p.Column = col
	p.Box = (((row - 1) / sudokuBoxes) * sudokuBoxes + ((col - 1) / sudokuBoxes)) + 1

	return p
}

func (p Pos) eqRow(other Pos) bool {
	return p.Row == other.Row
}

func (p Pos) eqColumn(other Pos) bool {
	return p.Column == other.Column
}

func (p Pos) eqBox(other Pos) bool {
	return p.Box == other.Box
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
