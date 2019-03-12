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

type Cell struct {
	Value int8
	Pos   Pos
}

type CellNumbers struct {
	Pos     Pos
	Numbers []int8
}

type cellGetter func(x int8) CellList
type cellFinder func(cells CellList) finderResult
type cellPredicate func(Cell) bool

func getCellNumbers(pos Pos, cands CellList) CellNumbers {
	nums := cands.FilterMapInt8(
		func(c Cell) int8 { return c.Value },
		func(c Cell) bool { return c.Pos == pos })

	return CellNumbers{Pos: pos, Numbers: nums}
}

func (c Cell) init(row int8, col int8, value int8) Cell {
	c.Pos = Pos{}.init(row, col)
	c.Value = value

	return c
}

func (c Cell) eqPos(other Cell) bool {
	return c.Pos == other.Pos
}

func (c Cell) eqValue(other Cell) bool {
	return c.Value == other.Value
}

func validateSet(cells []Cell) bool {
	nums := make(map[int8]Pos)

	for _, cell := range cells {
		if _, ok := nums[cell.Value]; ok {
			return false
		}
		nums[cell.Value] = cell.Pos
	}

	return true
}

func (c Cell) less(other *Cell) bool {
	if c.Pos.less(&other.Pos) {
		return true
	}

	if c.Pos == other.Pos {
		return c.Value < other.Value
	}

	return false
}

func (c Cell) inBox(others CellList) bool {
	for _, other := range others {
		if c.Pos.Box != other.Pos.Box {
			return false
		}
	}
	return true
}

func (c Cell) inColumn(others CellList) bool {
	return others.All(func(other Cell) bool {
		return c.Pos.Column == other.Pos.Column
	})
}

func (c Cell) inRow(others CellList) bool {
	return others.All(func(other Cell) bool {
		return c.Pos.Row == other.Pos.Row
	})
}

func (c Cell) inLine(others CellList) bool {
	return c.inRow(others) || c.inColumn(others)
}

func cellPositions(cells CellList) []Pos {
	n := len(cells)

	res := make([]Pos, n)

	for i := 0; i < n; i++ {
		res[i] = cells[i].Pos
	}

	return res
}

func dedupeCells(cells CellList) CellList {
	if len(cells) <= 1 {
		return cells
	}

	prev := cells[0]
	res := CellList{prev}

	for _, cell := range cells[1:] {
		if cell != prev {
			res = append(res, cell)
			prev = cell
		}
	}

	return res
}

func sortCells(array CellList) {
	sort.Slice(array, func(i, j int) bool {
		return array[i].less(&array[j])
	})
}

func numbers(cells CellList) []int8 {
	res := cells.MapInt8(func(cell Cell) int8 { return cell.Value })

	return res
}

func uniqueNumbers(cells CellList) []int8 {
	res := numbers(cells)
	sortInt8(res)
	res = dedupeInt8(res)

	return res
}

func uniqueCells(cells CellList) CellList {
	sortCells(cells)
	res := dedupeCells(cells)

	return res
}

