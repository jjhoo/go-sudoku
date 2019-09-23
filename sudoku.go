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
//go:generate go run github.com/kulshekhar/fungen -package sudoku -types Pos,Cell,numCount,int8,int

package sudoku

import (
	"fmt"
	"github.com/deckarep/golang-set"
	"unicode"
)

const (
	sudokuBoxes    = 3
	sudokuNumbers  = 9
	sudokuGridSize = sudokuNumbers * sudokuNumbers
)

type Sudoku struct {
	Solved     CellList
	Candidates CellList

	enableLogging bool
	logger        Logger
}

func (s Sudoku) logEliminated(strategy string, cells ...Cell) {
	if s.enableLogging && s.logger != nil {
		s.logger.Eliminated(strategy, cells...)
	}
}

func (s Sudoku) logSolved(strategy string, cells ...Cell) {
	if s.enableLogging && s.logger != nil {
		s.logger.Solved(strategy, cells...)
	}
}

type finderResult struct {
	Solved     CellList
	Eliminated CellList
}

func (s Sudoku) getCell(row, col int8) Cell {
	idx := (row - 1) * sudokuNumbers + (col - 1)
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
	return s.Solved.Filter(func(c Cell) bool {
		return c.Value != 0 && c.Pos.Box == box
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
	return s.Candidates.Filter(func(c Cell) bool {
		return c.Pos.Box == box
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

	for i = 1; i <= sudokuNumbers; i++ {
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
	s := Sudoku{logger: DefaultLogger{}, enableLogging: false}

	err := s.initGrid(grid)
	if err != nil {
		return nil, err
	}

	s.initCandidates()

	return &s, nil
}

func (s *Sudoku) EnableLogging(state bool) {
	s.enableLogging = state
}

func (s *Sudoku) initGrid(grids string) error {
	if len(grids) != sudokuGridSize {
		return fmt.Errorf("Grid has invalid size '%d'", len(grids))
	}

	var row int8 = 1
	var column int8 = 1

	s.Solved = make(CellList, sudokuGridSize)

	var i int
	var c rune
	zero := rune('0')

	for i, c = range grids {
		if !unicode.IsDigit(c) {
			return fmt.Errorf("Invalid rune '%c' in grid", c)
		}

		ascii := int8(c - zero)

		idx := (row - 1) * sudokuNumbers + (column - 1)
		s.Solved[idx] = Cell{}.init(row, column, ascii)

		if (i + 1) % sudokuNumbers == 0 {
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

	for i := 1; i <= sudokuNumbers; i++ {
		for j := 1; j <= sudokuNumbers; j++ {
			for n := 1; n <= sudokuNumbers; n++ {
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
		if (i + 1) % sudokuNumbers == 1 {
			fmt.Print("| ")
		}

		v := cell.Value
		if v == 0 {
			fmt.Printf(".")
		} else {
			fmt.Printf("%d", v)
		}

		if (i + 1) % sudokuNumbers == 0 {
			fmt.Print(" |\n")
		} else {
			fmt.Print(" ")
		}
	}
	fmt.Print("+-------------------+\n")
}

func (s Sudoku) GetGridString() string {
	runes := make([]rune, sudokuGridSize)

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
		idx := (sol.Pos.Row-1)*sudokuNumbers + (sol.Pos.Column - 1)
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
		for i := 1; i <= sudokuNumbers; i++ {
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
		matchedPositions := mapset.NewSet()

		var idxs intList = combs.next()
		if idxs == nil {
			break
		}

		comb := idxs.MapInt8(func(n int) int8 { return nums[n] })

		set1 := mapset.NewSet()
		for _, n := range comb {
			set1.Add(n)
		}

		// fmt.Println("visit comb", comb)
	OUTER:
		for _, pos := range poss {
			cnums := getCellNumbers(pos, cands)
			// fmt.Println("cnums", cnums)

			if len(cnums.Numbers) > limit {
				continue OUTER
			}

			set2 := mapset.NewSet()
			for _, n := range cnums.Numbers {
				set2.Add(n)
			}

			if set2.IsSubset(set1) {
				// fmt.Println("is subset", comb, cnums.Numbers)
				matchedPositions.Add(pos)
			}
		}

		if matchedPositions.Cardinality() == limit {
			nfound := cands.Filter(func(c Cell) bool {
				return !matchedPositions.Contains(c.Pos) && set1.Contains(c.Value)
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

func findHiddenGroupsInSet(limit int, cands CellList) finderResult {
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
		matchedPositions := mapset.NewSet()

		var idxs intList = combs.next()
		if idxs == nil {
			break
		}

		comb := idxs.MapInt8(func(n int) int8 { return unums[n] })

		set1 := mapset.NewSet()
		for _, n := range comb {
			set1.Add(n)
		}

		for _, pos := range poss {
			cnums := getCellNumbers(pos, cands)

			set2 := mapset.NewSet()
			for _, n := range cnums.Numbers {
				set2.Add(n)
			}

			if set2.Intersect(set1).Cardinality() > 0 {
				matchedPositions.Add(pos)
			}
		}

		if matchedPositions.Cardinality() == limit {
			nfound := cands.Filter(func(c Cell) bool {
				// true if position matches but number is not in the combination
				return matchedPositions.Contains(c.Pos) && !set1.Contains(c.Value)
			})
			found = append(found, nfound...)
		}
	}

	return finderResult{Solved: nil, Eliminated: found}
}

func (s *Sudoku) findHiddenGroups2() finderResult {
	return s.finder(func(cells CellList) finderResult {
		found := findHiddenGroupsInSet(2, cells)

		found.Eliminated = uniqueCells(found.Eliminated)

		return found
	})
}

func (s *Sudoku) findHiddenGroups3() finderResult {
	return s.finder(func(cells CellList) finderResult {
		found := findHiddenGroupsInSet(3, cells)

		found.Eliminated = uniqueCells(found.Eliminated)

		return found
	})
}

func (s *Sudoku) findHiddenGroups4() finderResult {
	return s.finder(func(cells CellList) finderResult {
		found := findHiddenGroupsInSet(4, cells)

		found.Eliminated = uniqueCells(found.Eliminated)

		return found
	})
}

func (s *Sudoku) findPointingPairs() finderResult {
	found := CellList{}

	var boxnum int8
	for boxnum = 1; boxnum <= sudokuNumbers; boxnum++ {
		boxCells := s.getCandidateBox(boxnum)
		nums := uniqueNumbers(boxCells)

		for _, n := range nums {
			cells := boxCells.Filter(func(c Cell) bool {
				return c.Value == n
			})

			if len(cells) < 2 {
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
		isSameLine func(Pos, Pos) bool
	}

	fpairs := []pair{
		{s.getCandidateRow, Pos.eqRow},
		{s.getCandidateColumn, Pos.eqColumn},
	}

	for _, fpair := range fpairs {
		var n int8
		for n = 1; n <= sudokuNumbers; n++ {
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
	for key := range interesting {
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

			if w1.sees(pivot) && w1.sees(w2) {
				// fmt.Println("y-wing would be a naked triple", w1, pivot, w2)
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

func (s *Sudoku) findXYZWings() finderResult {
	found := CellList{}
	interesting := make(map[Pos][]int8)

	prev := s.Candidates[0]
	nums := []int8{prev.Value}

	for _, cell := range s.Candidates[1:] {
		if prev.Pos == cell.Pos {
			nums = append(nums, cell.Value)
		} else {
			if len(nums) == 2 || len(nums) == 3 {
				interesting[prev.Pos] = nums
			}

			nums = []int8{cell.Value}
		}
		prev = cell
	}

	if len(nums) == 2 || len(nums) == 3 {
		interesting[prev.Pos] = nums
	}

	poss := []Pos{}
	for key := range interesting {
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

			if w1.sees(pivot) && w1.sees(w2) {
				// fmt.Println("xyz-wing would be a naked triple", w1, pivot, w2)
				continue
			}

			fun := func(p Pos) mapset.Set {
				set := mapset.NewSet()

				for _, n := range interesting[p] {
					set.Add(n)
				}

				return set
			}

			pset := fun(pivot)
			if pset.Cardinality() != 3 {
				continue
			}

			set1 := fun(w1)
			set2 := fun(w2)

			// wings have common number
			nset := set1.Intersect(set2)

			if !(nset.Cardinality() == 1 && pset.Equal(set1.Union(set2))) {
				continue
			}

			// fmt.Println("xyz-wing", w1, pivot, w2, set1, set2, pset)

			common := []int8{}
			for _, n := range nset.ToSlice() {
				common = append(common, n.(int8))
			}
			n := common[0]

			nfound := s.Candidates.Filter(func(c Cell) bool {
				return c.Value == n && c.Pos.sees(pivot) && c.Pos.sees(w1) && c.Pos.sees(w2)
			})

			found = append(found, nfound...)
		}
	}

	return finderResult{Solved: nil, Eliminated: found}
}

func (s *Sudoku) findXWings() finderResult {
	type pair struct {
		linef     func(int8) CellList
		eqLine    func(Pos, Pos) bool
		parallelf func(Pos) CellList
	}

	var i, j int8
	finder := func(fp pair) CellList {
		found := CellList{}

		for i = 1; i <= 8; i++ {
			line1 := fp.linef(i)

			if len(line1) <= 2 {
				continue
			}

			line1Numbers := numbers(line1)
			twos1 := numberCounts(line1Numbers).Filter(
				func(nc numCount) bool {
					return nc.count == 2
				})

			if len(twos1) < 1 {
				continue
			}

			for j = i + 1; j <= sudokuNumbers; j++ {
				line2 := fp.linef(j)
				if len(line2) <= 2 {
					continue
				}

				line2Numbers := numbers(line2)
				twos2 := numberCounts(line2Numbers).Filter(
					func(nc numCount) bool {
						return nc.count == 2
					})

				if len(twos2) < 1 {
					continue
				}

				for _, nc := range twos1 {
					res := twos2.Filter(func(other numCount) bool {
						return nc.num == other.num
					})

					if len(res) != 1 {
						continue
					}

					// fmt.Printf("Check lines %d -- %d\n", i, j)
					// fmt.Println(nc, res[0])

					cells1 := line1.Filter(func(c Cell) bool {
						return nc.num == c.Value
					})

					cells2 := line2.Filter(func(c Cell) bool {
						return nc.num == c.Value
					})

					// fmt.Println(cells1, cells2)

					if fp.eqLine(cells1[0].Pos, cells2[0].Pos) &&
						fp.eqLine(cells1[1].Pos, cells2[1].Pos) {
						found1 := fp.parallelf(cells1[0].Pos).Filter(func(c Cell) bool {
							return nc.num == c.Value &&
								cells1[0].Pos != c.Pos &&
								cells2[0].Pos != c.Pos
						})

						found2 := fp.parallelf(cells1[1].Pos).Filter(func(c Cell) bool {
							return nc.num == c.Value &&
								cells1[1].Pos != c.Pos &&
								cells2[1].Pos != c.Pos
						})

						if len(found1) > 0 {
							found = append(found, found1...)
						}

						if len(found2) > 0 {
							found = append(found, found2...)
						}
					}
				}
			}
		}

		return found
	}

	fpairs := []pair{
		{s.getCandidateRow, Pos.eqColumn,
			func(p Pos) CellList {
				return s.getCandidateColumn(p.Column)
			},
		},
		{s.getCandidateColumn, Pos.eqRow,
			func(p Pos) CellList {
				return s.getCandidateRow(p.Row)
			},
		},
	}

	found := CellList{}

	for _, fp := range fpairs {
		cells := finder(fp)
		found = append(found, cells...)
	}

	return finderResult{Solved: nil, Eliminated: found}
}

func PrintGrid(grid string) error {
	if len(grid) != sudokuGridSize {
		return fmt.Errorf("Grid has invalid size '%d'", len(grid))
	}

	for i, c := range grid {
		fmt.Printf("%c", c)
		if (i + 1) % sudokuNumbers == 0 {
			fmt.Print("\n")
		} else {
			fmt.Print(" ")
		}
	}

	return nil
}

func (s *Sudoku) Solve() bool {
	type finderFunc struct {
		fun  func() finderResult
		name string
	}

	finders := []finderFunc{
		{fun: s.findSinglesSimple, name: "singles (simple)"},
		{fun: s.findSingles, name: "singles"},
		{fun: s.findNakedGroups2, name: "naked pairs"},
		{fun: s.findNakedGroups3, name: "naked triples"},
		{fun: s.findHiddenGroups2, name: "hidden pairs"},
		{fun: s.findHiddenGroups3, name: "hidden triples"},
		{fun: s.findNakedGroups4, name: "naked quads"},
		{fun: s.findHiddenGroups4, name: "hidden quads"},
		{fun: s.findPointingPairs, name: "pointing pairs"},
		{fun: s.findBoxlineReduction, name: "box/line reduction"},
		{fun: s.findXWings, name: "x-wing"},
		{fun: s.findYWings, name: "y-wing"},
		{fun: s.findXYZWings, name: "xyz-wing"},
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

		finder := finders[finderIdx]
		res := finder.fun()

		if len(res.Solved) > 0 {
			s.logSolved(finder.name, res.Solved...)
			s.updateSolved(res.Solved)
		}
		s.validate()

		if len(res.Eliminated) > 0 {
			s.logEliminated(finder.name, res.Eliminated...)
			s.updateCandidates(res.Eliminated)
		}

		if len(s.Candidates) == 0 {
			return true
		}

		if len(res.Solved) != 0 || len(res.Eliminated) != 0 {
			// fmt.Println("progress", len(s.Candidates))
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
