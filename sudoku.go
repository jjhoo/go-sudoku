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

type cell_predicate func(Cell) bool

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

func filter(cells []Cell, pred cell_predicate) []Cell {
	res := []Cell{}

	for _, cell := range cells {
		if pred(cell) {
			res = append(res, cell)
		}
	}

	return res
}

func (s Sudoku) getRow(row int8) []Cell {
	return filter(s.Solved, func(c Cell) bool {
		return c.Pos.Row == row
	})
}

func (s Sudoku) getColumn(col int8) []Cell {
	return filter(s.Solved, func(c Cell) bool {
		return c.Pos.Column == col
	})
}

func (s Sudoku) getBox(box int8) []Cell {
	b := numToBox(box)
	return filter(s.Solved, func(c Cell) bool {
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

func (s *Sudoku) initGrid(grids string) error {
	if len(grids) != 81 {
		return fmt.Errorf("Grid '%s' has invalid size", grids)
	}

	var row int8 = 1
	var column int8 = 1

	zero := int8('0')
	s.Solved = []Cell{}

	for i, c := range grids {
		ascii := int8(c) - zero

		if ascii > 0 {
			s.Solved = append(s.Solved, Cell{}.init(row, column, ascii))
		}

		if (i+1)%9 == 0 {
			row += 1
			column = 1
		} else {
			column += 1
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

	for _, cell := range s.Solved {
		if cell.Value != 0 {

		}
	}
	fmt.Println(s.Candidates)
}

func (s Sudoku) printGrid() {
	fmt.Println(s.Solved)
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

func main() {
	var grid1 string = "700600008800030000090000310006740005005806900400092100087000020000060009600008001"

	err := printGrid(grid1)
	if err != nil {
		fmt.Printf("%s\n", err)
	}

	var s Sudoku

	err = s.initGrid(grid1)
	if err == nil {
		s.printGrid()
	}

	// s.initCandidates()
	fmt.Println(s.getRow(1))
	fmt.Println(s.getColumn(1))
	fmt.Println(s.getBox(1))
}
