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
//go:generate fungen -package sudoku -types Box,Pos,Cell,numCount,int8,int

package sudoku

import (
	"fmt"
	"github.com/deckarep/golang-set"
	"unicode"
)

type Sudoku struct {
	Solved     CellList
	Candidates CellList
}

type finderResult struct {
	Solved     CellList
	Eliminated CellList
}

func (s Sudoku) getCell(row, col int8) Cell {
	idx := (row-1)*9 + (col - 1)
	return s.Solved[idx]
}

func (s Sudoku) getRow(row int8) CellList {
	return s.Solved.Filter(func(c Cell) bool {
		return c.Value != 0 && c.Pos.Row == row
	})
}

func (s Sudoku) getColumn(col int8) CellList {
	return s.Solved.Filter(func(c Cell) bool {
		return c.Value != 0 && c.Pos.Column == col
	})
}

func (s Sudoku) getBox(box int8) CellList {
	b := numToBox(box)
	return s.Solved.Filter(func(c Cell) bool {
		return c.Value != 0 && c.Pos.Box == b
	})
}

func (s Sudoku) getCandidateCell(row, col int8) CellList {
	return s.Candidates.Filter(func(c Cell) bool {
		return c.Pos.Row == row && c.Pos.Column == col
	})
}

func (s Sudoku) getCandidateRow(row int8) CellList {
	return s.Candidates.Filter(func(c Cell) bool {
		return c.Pos.Row == row
	})
}

func (s Sudoku) getCandidateColumn(col int8) CellList {
	return s.Candidates.Filter(func(c Cell) bool {
		return c.Pos.Column == col
	})
}

func (s Sudoku) getCandidateBox(box int8) CellList {
	b := numToBox(box)
	return s.Candidates.Filter(func(c Cell) bool {
		return c.Pos.Box == b
	})
}

func (s Sudoku) getCellNumbers(pos Pos) cellNumbers {
	return getCellNumbers(pos, s.Candidates)
}

func (s *Sudoku) validateSolved() {
	type pair struct {
		desc string
		fun  func(int8) CellList
	}

	var i int8
	pairs := []pair{
		{"row", s.getRow}, {"column", s.getColumn}, {"box", s.getBox},
	}

	for i = 1; i < 10; i++ {
		for _, p := range pairs {
			set := p.fun(i)
			if !validateSet(set) {
				panic(fmt.Sprintf("Invalid %s %d %v", p.desc, i, set))
			}
		}
	}
}

func (s *Sudoku) validate() {
	s.validateSolved()
}

func NewSudoku(grid string) (*Sudoku, error) {
	s := Sudoku{}

	err := s.initGrid(grid)
	if err != nil {
		return nil, err
	}

	s.initCandidates()

	return &s, nil
}

func (s *Sudoku) initGrid(grids string) error {
	if len(grids) != 81 {
		return fmt.Errorf("Grid '%s' has invalid size", grids)
	}

	var row int8 = 1
	var column int8 = 1

	s.Solved = make(CellList, 81)

	var i int
	var c rune
	zero := rune('0')

	for i, c = range grids {
		if !unicode.IsDigit(c) {
			return fmt.Errorf("Invalid rune '%c' in grid", c)
		}

		ascii := int8(c - zero)

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
	s.Candidates = CellList{}

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

		s.Candidates = s.Candidates.Filter(func(c Cell) bool {
			return !(solved.Pos == c.Pos || (solved.Pos.sees(c.Pos) && solved.eqValue(c)))
		})
	}
}

func (s Sudoku) PrintGrid() {
	fmt.Print("+-------------------+\n")
	for i, cell := range s.Solved {
		if (i+1)%9 == 1 {
			fmt.Print("| ")
		}
		v := cell.Value
		if v == 0 {
			fmt.Printf(".")
		} else {
			fmt.Printf("%d", v)
		}
		if (i+1)%9 == 0 {
			fmt.Print(" |\n")
		} else {
			fmt.Print(" ")
		}
	}
	fmt.Print("+-------------------+\n")
}

func (s Sudoku) GetGridString() string {
	runes := make([]rune, 81)

	for i, cell := range s.Solved {
		runes[i] = '0' + rune(cell.Value)
	}
	return string(runes)
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

func (s *Sudoku) updateSolved(solved CellList) {
	for _, sol := range solved {
		idx := (sol.Pos.Row-1)*9 + (sol.Pos.Column - 1)
		s.Solved[idx].Value = sol.Value

		s.Candidates = s.Candidates.Filter(func(c Cell) bool {
			return !(sol.Pos == c.Pos || (sol.Pos.sees(c.Pos) && sol.eqValue(c)))
		})
	}
}

func (s *Sudoku) updateCandidates(eliminated CellList) {
	for _, cell := range eliminated {
		s.Candidates = s.Candidates.Filter(func(c Cell) bool {
			return !(cell.Pos == c.Pos && cell.Value == c.Value)
		})
	}
}

// Simple case where there is only one candidate left for a cell
func (s *Sudoku) findSinglesSimple() finderResult {
	poss := s.ucpos()
	found := CellList{}

	for _, pos := range poss {
		cands := s.Candidates.Filter(func(c Cell) bool {
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

	found := CellList{}
	eliminated := CellList{}

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

// Only one candidate left for a number in row / column / box
func (s *Sudoku) findSingles() finderResult {
	return s.finder(func(cells CellList) finderResult {
		nums := uniqueNumbers(cells)
		// fmt.Println(nums)

		found := CellList{}

		for _, n := range nums {
			ncells := cells.Filter(func(cell Cell) bool {
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

func findNakedGroupsInSet(limit int, cands CellList) finderResult {
	// fmt.Println("naked set", limit)
	poss := ucpos(cands)

	if len(poss) < (limit + 1) {
		return finderResult{Solved: nil, Eliminated: nil}
	}

	nums := numbers(cands)
	sortInt8(nums)

	ncounts := numberCounts(nums)
	unums := ncounts.MapInt8(func(nc numCount) int8 { return nc.num })

	// fmt.Println("counts", cands, ncounts, unums)
	found := CellList{}

	combs := newCombination(len(unums), limit)
	for {
		matches := []cellNumbers{}
		others := []Pos{}

		var idxs intList = combs.next()
		if idxs == nil {
			break
		}

		comb := idxs.MapInt8(func(n int) int8 { return nums[n] })

		// fmt.Println("visit comb", comb)
	OUTER:
		for _, pos := range ucpos(cands) {
			cnums := getCellNumbers(pos, cands)
			// fmt.Println("cnums", cnums)

			if len(cnums.Numbers) > limit {
				others = append(others, cnums.Pos)
				continue OUTER
			}

			set1 := mapset.NewSet()
			for _, n := range comb {
				set1.Add(n)
			}

			set2 := mapset.NewSet()
			for _, n := range cnums.Numbers {
				set2.Add(n)
			}

			if set2.IsSubset(set1) {
				// fmt.Println("is subset", comb, cnums.Numbers)
				matches = append(matches, cnums)
			} else {
				// fmt.Println("is not subset", comb, cnums.Numbers)
				others = append(others, cnums.Pos)
			}
		}

		if len(matches) == limit && len(others) > 0 {
			// fmt.Println("matches", matches, ", others", others)

			nfound := cands.Filter(func(c Cell) bool {
				for _, other := range others {
					if c.Pos != other {
						continue
					}

					// fmt.Println("check", c.Pos)

					for _, n := range comb {
						if n == c.Value {
							// fmt.Println("found", c.Pos, n)
							return true
						}
					}
				}
				return false
			})
			found = append(found, nfound...)
		}
	}

	return finderResult{Solved: nil, Eliminated: found}
}

func (s *Sudoku) findNakedGroups2() finderResult {
	return s.finder(func(cells CellList) finderResult {
		found := findNakedGroupsInSet(2, cells)

		found.Eliminated = uniqueCells(found.Eliminated)

		return found
	})
}

func (s *Sudoku) findNakedGroups3() finderResult {
	return s.finder(func(cells CellList) finderResult {
		found := findNakedGroupsInSet(3, cells)

		found.Eliminated = uniqueCells(found.Eliminated)

		return found
	})
}

func (s *Sudoku) findNakedGroups4() finderResult {
	return s.finder(func(cells CellList) finderResult {
		found := findNakedGroupsInSet(4, cells)

		found.Eliminated = uniqueCells(found.Eliminated)

		return found
	})
}

func (s *Sudoku) findPointingPairs() finderResult {
	found := CellList{}

	var boxnum int8
	for boxnum = 1; boxnum < 10; boxnum++ {
		boxCells := s.getCandidateBox(boxnum)
		nums := uniqueNumbers(boxCells)

		for _, n := range nums {
			cells := boxCells.Filter(func(c Cell) bool {
				return c.Value == n
			})

			if !(len(cells) == 2 || len(cells) == 3) {
				continue
			}

			var others CellList
			if cells[0].inRow(cells[1:]) {
				others = s.getCandidateRow(cells[0].Pos.Row)
			} else if cells[0].inColumn(cells[1:]) {
				others = s.getCandidateColumn(cells[0].Pos.Column)
			} else {
				continue
			}

			nfound := others.Filter(func(c Cell) bool {
				return c.Value == n && !cells[0].Pos.eqBox(c.Pos)
			})

			if len(nfound) > 0 {
				// fmt.Println("pointing pairs", boxnum, n, cells, nfound)
				found = append(found, nfound...)
			}
		}
	}

	return finderResult{Solved: nil, Eliminated: found}
}

func (s *Sudoku) findBoxlineReduction() finderResult {
	found := CellList{}

	type pair struct {
		getCells   func(int8) CellList
		isSameLine func(Pos, ...Pos) bool
	}

	fpairs := []pair{
		{s.getCandidateRow, Pos.eqRow},
		{s.getCandidateColumn, Pos.eqColumn},
	}

	for _, fpair := range fpairs {
		var n int8
		for n = 1; n < 10; n++ {
			cells := fpair.getCells(n)
			if len(cells) < 2 {
				// No panic, even though other finders should have
				// eliminated this case
				continue
			}

			nums := numbers(cells)
			sortInt8(nums)

			ncounts := numberCounts(nums)

			for _, nc := range ncounts {
				if !(nc.count == 2 || nc.count == 3) {
					continue
				}

				// fmt.Println("boxline: good count")

				ncells := cells.Filter(func(c Cell) bool {
					return nc.num == c.Value
				})

				if !(ncells[0].inBox(ncells[1:])) {
					continue
				}

				// fmt.Println("boxline: in same box", ncells)

				// Needs to be same value, same box, different row/col
				nfound := s.Candidates.Filter(func(c Cell) bool {
					return ncells[0].Value == c.Value && ncells[0].Pos.Box == c.Pos.Box && !fpair.isSameLine(ncells[0].Pos, c.Pos)
				})

				// fmt.Println("boxline found", nfound)
				found = append(found, nfound...)
			}
		}
	}

	return finderResult{Solved: nil, Eliminated: found}
}

func (s *Sudoku) findYWings() finderResult {
	found := CellList{}
	interesting := make(map[Pos][]int8)

	prev := s.Candidates[0]
	nums := []int8{prev.Value}

	for _, cell := range s.Candidates[1:] {
		if prev.Pos == cell.Pos {
			nums = append(nums, cell.Value)
		} else {
			if len(nums) == 2 {
				interesting[prev.Pos] = nums
			}

			nums = []int8{cell.Value}
		}
		prev = cell
	}

	if len(nums) == 2 {
		interesting[prev.Pos] = nums
	}

	poss := []Pos{}
	for key, _ := range interesting {
		poss = append(poss, key)
	}

	// fmt.Println("y-wing interesting", interesting)
	combs := newCombination(len(poss), 3)

	for {
		var idxs intList = combs.next()
		if idxs == nil {
			break
		}

		out := make(PosList, len(idxs))

		nums := int8List{}

		for i, n := range idxs {
			out[i] = poss[n]

			nums = append(nums, interesting[poss[n]]...)
		}

		sortInt8(nums)
		nums = dedupeInt8(nums)

		if len(nums) != 3 {
			continue
		}

		// fmt.Println("y-wing 1", out, nums)

		perms := newPermutation(len(out))
		for {
			var idxs intList = perms.next()
			if idxs == nil {
				break
			}

			out2 := idxs.MapPos(func(n int) Pos { return out[n] })

			w1, pivot, w2 := out2[0], out2[1], out2[2]

			// Avoid duplicate wings
			if w2.less(&w1) {
				continue
			}

			if !(pivot.sees(w1) && pivot.sees(w2)) {
				continue
			}

			if w1.eqColumn(pivot, w2) || w1.eqRow(pivot, w2) || w1.eqBox(pivot, w2) {
				// fmt.Println("y-wing is a naked triple", w1, pivot, w2)
				continue
			}

			fun := func(p Pos) mapset.Set {
				set := mapset.NewSet()

				for _, n := range interesting[p] {
					set.Add(n)
				}

				return set
			}

			set1 := fun(w1)
			set2 := fun(w2)
			pset := fun(pivot)

			// wings have common number
			nset := set1.Intersect(set2)

			if !(nset.Cardinality() == 1 &&
				pset.IsSubset(set1.Union(set2).Difference(nset))) {
				continue
			}

			// fmt.Println("y-wing 2", w1, pivot, w2, set1, set2, pset, set1.Union(set2).Union(pset))

			common := []int8{}
			for _, n := range nset.ToSlice() {
				common = append(common, n.(int8))
			}
			n := common[0]

			nfound := s.Candidates.Filter(func(c Cell) bool {
				return c.Value == n && c.Pos.sees(w1) && c.Pos.sees(w2)
			})

			found = append(found, nfound...)
		}
	}

	return finderResult{Solved: nil, Eliminated: found}
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

func (s *Sudoku) Solve() bool {
	finders := []func() finderResult{
		s.findSinglesSimple,
		s.findSingles,
		s.findNakedGroups2,
		s.findNakedGroups3,
		s.findNakedGroups4,
		s.findPointingPairs,
		s.findBoxlineReduction,
		s.findYWings,
	}

	// fmt.Println("begin", len(s.Candidates))
	finderCount := len(finders)
	finderIdx := 0

PROGRESS:
	for finderIdx < finderCount {
		if len(s.Candidates) == 0 {
			return true
		}

		// fmt.Println("Finder", finderIdx)

		res := finders[finderIdx]()

		if len(res.Solved) > 0 {
			fmt.Println("Found", res.Solved)
			s.updateSolved(res.Solved)
		}
		s.validate()

		if len(res.Eliminated) > 0 {
			fmt.Println("Eliminated", res.Eliminated)
			s.updateCandidates(res.Eliminated)
		}

		if len(s.Candidates) == 0 {
			return true
		}

		if len(res.Solved) != 0 || len(res.Eliminated) != 0 {
			fmt.Println("progress", len(s.Candidates))
			finderIdx = 0
			continue PROGRESS
		}
		finderIdx++
	}
	return false
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

func ucpos(cells CellList) []Pos {
	poss := cellPositions(cells)
	poss = dedupePos(poss)

	return poss
}
