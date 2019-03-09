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
	"fmt"
	"sort"
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

type Cell struct {
	Value int8
	Pos   Pos
}

type Sudoku struct {
	Solved     []Cell
	Candidates []Cell
}

type finderResult struct {
	Solved     []Cell
	Eliminated []Cell
}

type cellGetter func(x int8) []Cell
type cellFinder func(cells []Cell) finderResult
type cellPredicate func(Cell) bool

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

func filter(cells []Cell, pred cellPredicate) []Cell {
	res := []Cell{}

	for _, cell := range cells {
		if pred(cell) {
			res = append(res, cell)
		}
	}

	return res
}

func remove(cells []Cell, pred cellPredicate) []Cell {
	res := []Cell{}

	for _, cell := range cells {
		if !pred(cell) {
			res = append(res, cell)
		}
	}

	return res
}

func (s Sudoku) getCell(row, col int8) Cell {
	idx := (row-1)*9 + (col - 1)
	return s.Solved[idx]
}

func (s Sudoku) getRow(row int8) []Cell {
	return filter(s.Solved, func(c Cell) bool {
		return c.Value != 0 && c.Pos.Row == row
	})
}

func (s Sudoku) getColumn(col int8) []Cell {
	return filter(s.Solved, func(c Cell) bool {
		return c.Value != 0 && c.Pos.Column == col
	})
}

func (s Sudoku) getBox(box int8) []Cell {
	b := numToBox(box)
	return filter(s.Solved, func(c Cell) bool {
		return c.Value != 0 && c.Pos.Box == b
	})
}

func (s Sudoku) getCandidateCell(row, col int8) []Cell {
	return filter(s.Candidates, func(c Cell) bool {
		return c.Pos.Row == row && c.Pos.Column == col
	})
}

func (s Sudoku) getCandidateRow(row int8) []Cell {
	return filter(s.Candidates, func(c Cell) bool {
		return c.Pos.Row == row
	})
}

func (s Sudoku) getCandidateColumn(col int8) []Cell {
	return filter(s.Candidates, func(c Cell) bool {
		return c.Pos.Column == col
	})
}

func (s Sudoku) getCandidateBox(box int8) []Cell {
	b := numToBox(box)
	return filter(s.Candidates, func(c Cell) bool {
		return c.Pos.Box == b
	})
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

func (c Cell) init(row int8, col int8, value int8) Cell {
	c.Pos = Pos{}.init(row, col)
	c.Value = value

	return c
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

func (c Cell) eqPos(other Cell) bool {
	return c.Pos == other.Pos
}

func (c Cell) eqValue(other Cell) bool {
	return c.Value == other.Value
}

func (s *Sudoku) initGrid(grids string) error {
	if len(grids) != 81 {
		return fmt.Errorf("Grid '%s' has invalid size", grids)
	}

	var row int8 = 1
	var column int8 = 1

	zero := int8('0')
	s.Solved = make([]Cell, 81)

	for i, c := range grids {
		ascii := int8(c) - zero

		idx := (row-1)*9 + (column - 1)
		s.Solved[idx] = Cell{}.init(row, column, ascii)

		if (i+1)%9 == 0 {
			row++
			column = 1
		} else {
			column++
		}
	}

	return nil
}

func (s *Sudoku) initCandidates() {
	s.Candidates = []Cell{}

	for i := 1; i < 10; i++ {
		for j := 1; j < 10; j++ {
			for n := 1; n < 10; n++ {
				s.Candidates = append(
					s.Candidates, Cell{}.init(int8(i), int8(j), int8(n)))
			}
		}
	}

	for _, solved := range s.Solved {
		if solved.Value == 0 {
			continue
		}

		s.Candidates = remove(s.Candidates, func(c Cell) bool {
			return solved.Pos == c.Pos || (solved.Pos.sees(c.Pos) && solved.eqValue(c))
		})
	}
}

func (s Sudoku) printGrid() {
	fmt.Println(filter(s.Solved, func(c Cell) bool { return c.Value != 0 }))
}

func (s Sudoku) ucpos() []Pos {
	res := []Pos{}

	if len(s.Candidates) == 0 {
		return res
	}

	prev := s.Candidates[0].Pos
	res = append(res, prev)

	for _, cell := range s.Candidates[1:] {
		if prev != cell.Pos {
			res = append(res, cell.Pos)
			prev = cell.Pos
		}
	}

	return res
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

func (c Cell) less(other *Cell) bool {
	if c.Pos.less(&other.Pos) {
		return true
	}

	if c.Pos == other.Pos {
		return c.Value < other.Value
	}

	return false
}

func (s *Sudoku) updateSolved(solved []Cell) {
	for _, sol := range solved {
		idx := (sol.Pos.Row-1)*9 + (sol.Pos.Column - 1)
		s.Solved[idx].Value = sol.Value

		s.Candidates = remove(s.Candidates, func(c Cell) bool {
			return sol.Pos == c.Pos || (sol.Pos.sees(c.Pos) && sol.eqValue(c))
		})
	}
}

func (s *Sudoku) updateCandidates(eliminated []Cell) {
	for _, cell := range eliminated {
		s.Candidates = remove(s.Candidates, func(c Cell) bool {
			return cell.Pos == c.Pos
		})
	}
}

// Simple case where there is only one candidate left for a cell
func (s *Sudoku) findSinglesSimple() finderResult {
	poss := s.ucpos()
	found := []Cell{}

	for _, pos := range poss {
		cands := filter(s.Candidates, func(c Cell) bool {
			return pos == c.Pos
		})

		if len(cands) == 1 {
			found = append(found, cands[0])
		}
	}

	found = uniqueCells(found)
	return finderResult{Solved: found, Eliminated: nil}
}

func (s *Sudoku) finder(cf cellFinder) finderResult {
	funs := []cellGetter{s.getCandidateRow, s.getCandidateColumn, s.getCandidateBox}

	found := []Cell{}
	eliminated := []Cell{}

	for _, fun := range funs {
		for i := 1; i < 10; i++ {
			cells := fun(int8(i))

			if len(cells) == 0 {
				continue
			}

			fresult := cf(cells)

			if len(fresult.Solved) > 0 {
				found = append(found, fresult.Solved...)
			}

			if len(fresult.Eliminated) > 0 {
				eliminated = append(eliminated, fresult.Eliminated...)
			}
		}
	}

	found = uniqueCells(found)
	eliminated = uniqueCells(eliminated)

	return finderResult{Solved: found, Eliminated: eliminated}
}

func mapCellInt8(cells []Cell, fn func(Cell) int8) []int8 {
	n := len(cells)

	res := make([]int8, n)

	for i := 0; i < n; i++ {
		res[i] = fn(cells[i])
	}

	return res
}

func cellPositions(cells []Cell) []Pos {
	n := len(cells)

	res := make([]Pos, n)

	for i := 0; i < n; i++ {
		res[i] = cells[i].Pos
	}

	return res
}

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

func dedupeCells(cells []Cell) []Cell {
	if len(cells) <= 1 {
		return cells
	}

	prev := cells[0]
	res := []Cell{prev}

	for _, cell := range cells[1:] {
		if cell != prev {
			res = append(res, cell)
			prev = cell
		}
	}

	return res
}

func sortCells(array []Cell) {
	sort.Slice(array, func(i, j int) bool {
		return array[i].less(&array[j])
	})
}

func uniqueNumbers(cells []Cell) []int8 {
	res := mapCellInt8(cells, func(cell Cell) int8 {
		return cell.Value
	})

	sortInt8(res)
	res = dedupeInt8(res)

	return res
}

func uniqueCells(cells []Cell) []Cell {
	sortCells(cells)
	res := dedupeCells(cells)

	return res
}

// Only one candidate left for a number in row / column / box
func (s *Sudoku) findSingles() finderResult {
	return s.finder(func(cells []Cell) finderResult {
		nums := uniqueNumbers(cells)
		// fmt.Println(nums)

		found := []Cell{}

		for _, n := range nums {
			ncells := filter(cells, func(cell Cell) bool {
				return cell.Value == n
			})

			if len(ncells) == 1 {
				found = append(found, ncells[0])
			}
		}

		found = uniqueCells(found)
		return finderResult{Solved: found, Eliminated: nil}
	})
}

func printGrid(grid string) error {
	if len(grid) != 81 {
		return fmt.Errorf("Grid '%s' has invalid size", grid)
	}

	for i, c := range grid {
		fmt.Printf("%c", c)
		if (i+1)%9 == 0 {
			fmt.Print("\n")
		} else {
			fmt.Print(" ")
		}
	}

	return nil
}

func (s *Sudoku) solve() {
	finders := []func() finderResult{
		s.findSinglesSimple,
		s.findSingles}

	fmt.Println("begin", len(s.Candidates))
	finderCount := len(finders)
	finderIdx := 0

PROGRESS:
	for finderIdx < finderCount {
		// fmt.Println("Finder", finderIdx)

		res := finders[finderIdx]()

		if len(res.Solved) > 0 {
			fmt.Println("Found", res.Solved)
			s.updateSolved(res.Solved)
		}

		if len(res.Eliminated) > 0 {
			fmt.Println("Eliminated", res.Eliminated)
			s.updateCandidates(res.Eliminated)
		}

		if len(res.Solved) != 0 || len(res.Eliminated) != 0 {
			fmt.Println("progress", len(s.Candidates))
			finderIdx = 0
			continue PROGRESS
		}
		finderIdx++
	}
}

func dedupePos(poss []Pos) []Pos {
	if len(poss) <= 1 {
		return poss
	}

	prev := poss[0]
	res := []Pos{prev}

	for _, pos := range poss[1:] {
		if pos != prev {
			res = append(res, pos)
			prev = pos
		}
	}

	return res
}

func ucpos(cells []Cell) []Pos {
	poss := cellPositions(cells)
	poss = dedupePos(poss)

	return poss
}

func main() {
	grid1 := "700600008800030000090000310006740005005806900400092100087000020000060009600008001"

	err := printGrid(grid1)
	if err != nil {
		fmt.Printf("%s\n", err)
	}

	var s Sudoku

	err = s.initGrid(grid1)
	if err == nil {
		s.printGrid()
	}

	fmt.Println(s.getRow(1))
	fmt.Println(s.getColumn(1))
	fmt.Println(s.getBox(1))

	s.initCandidates()

	s.printGrid()
	// fmt.Println(s.Candidates)
	s.solve()
	s.printGrid()

	if false {
		fresult := s.findSinglesSimple()
		if len(fresult.Solved) > 0 {
			s.updateSolved(fresult.Solved)
			s.updateCandidates(fresult.Solved)
		}
		fmt.Println("Found", fresult.Solved)
		// fmt.Println(s.Candidates)

		fresult = s.findSingles()
		if len(fresult.Solved) > 0 {
			s.updateSolved(fresult.Solved)
			s.updateCandidates(fresult.Solved)
		}
		fmt.Println("Found", fresult.Solved)
	}

	test := s.getCandidateRow(1)
	// poss := ucpos(test)
	nums := uniqueNumbers(test)
	foo := NewPermutation(nums, func(idxs []int) {
		out := make([]int8, len(idxs))

		for i, n := range idxs {
			out[i] = nums[n]
		}
		fmt.Println("visit", out)
	})
	fmt.Println("permutation test", nums, foo)

	for foo.Next() {
		foo.Visit()
	}
}
