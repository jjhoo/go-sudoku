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
	"fmt"
)

func main() {
	// grid1 := "700600008800030000090000310006740005005806900400092100087000020000060009600008001"
	// grid1 := "014600300050000007090840100000400800600050009007009000008016030300000010009008570"
	// grid1 := "000921003009000060000000500080403006007000800500700040003000000020000700800195000"
	grid1 := "000040700500780020070002006810007900460000051009600078900800010080064009002050000"

	err := printGrid(grid1)
	if err != nil {
		fmt.Printf("%s\n", err)
	}

	var s Sudoku

	err = s.initGrid(grid1)
	if err == nil {
		s.printGrid()
	} else {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println(s.getRow(1))
	fmt.Println(s.getColumn(1))
	fmt.Println(s.getBox(1))

	s.initCandidates()

	s.printGrid()
	// fmt.Println(s.Candidates)
	s.Solve()
	s.printGrid()
	fmt.Println(s.getGridString())

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

	if false {
		test := s.getCandidateRow(1)
		// poss := ucpos(test)
		nums := uniqueNumbers(test)
		visitf := func(idxs []int) {
			out := make([]int8, len(idxs))

			for i, n := range idxs {
				out[i] = nums[n]
			}
			fmt.Println("visit", out)
		}

		foo := Permutation(nums)
		fmt.Println("permutation test", nums, foo)

		for {
			idxs := foo.Next()

			if idxs == nil {
				break
			}
			visitf(idxs)
		}

		bar := Combination(nums, 2)
		fmt.Println("combination test", nums, bar)

		for {
			idxs := bar.Next()

			if idxs == nil {
				break
			}
			visitf(idxs)
		}
	}
}
