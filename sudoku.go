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

type Cell struct {
	Row    int8
	Column int8
	Value  int8
}

type Sudoku struct {
	Solved     []Cell
	Candidates []Cell
}

type cell_predicate func(Cell) bool

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
		return c.Row == row
	})
}

func (s Sudoku) getColumn(col int8) []Cell {
	return filter(s.Solved, func(c Cell) bool {
		return c.Column == col
	})
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

		s.Solved = append(s.Solved, Cell{Row: row, Column: column, Value: ascii})

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

	for _, cell := range s.Solved {
		if cell.Value == 0 {
			for i := 1; i < 10; i++ {
				s.Candidates = append(s.Candidates,
					Cell{Row: cell.Row, Column: cell.Column, Value: int8(i)})
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
}
