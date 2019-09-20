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

	"github.com/jjhoo/go-sudoku"
)

func main() {
	// grid1 := "700600008800030000090000310006740005005806900400092100087000020000060009600008001"
	// grid1 := "014600300050000007090840100000400800600050009007009000008016030300000010009008570"
	// grid1 := "000921003009000060000000500080403006007000800500700040003000000020000700800195000"
	grid1 := "000040700500780020070002006810007900460000051009600078900800010080064009002050000"

	err := sudoku.PrintGrid(grid1)
	if err != nil {
		fmt.Printf("%s\n", err)
	}

	s, err := sudoku.NewSudoku(grid1)
	if err == nil {
		s.PrintGrid()
	} else {
		fmt.Println(err)
		panic(err)
	}

	s.PrintGrid()
	// fmt.Println(s.Candidates)
	s.Solve()
	s.PrintGrid()
	fmt.Println(s.GetGridString())
}
