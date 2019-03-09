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

import "fmt"

type Cell struct {
	Row    int8
	Column int8
	Value  int8
}

func initGrid(grids string) ([]Cell, error) {
	if len(grids) != 81 {
		return nil, fmt.Errorf("Grid '%s' has invalid size", grids)
	}

	var grid []Cell

	var row int8 = 1
	var column int8 = 1

	zero := int8('0')

	for i, c := range grids {
		ascii := int8(c) - zero

		grid = append(grid, Cell{Row: row, Column: column, Value: ascii})

		if (i+1)%9 == 0 {
			row += 1
			column = 1
		} else {
			column += 1
		}
	}

	return grid, nil
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

	grid, err := initGrid(grid1)
	if err == nil {
		fmt.Println(grid)
	}
}
