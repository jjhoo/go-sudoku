package sudoku

import (
	"fmt"
	"testing"
)

func TestGrid1(t *testing.T) {
	grid := "000040700500780020070002006810007900460000051009600078900800010080064009002050000"

	s := Sudoku{}

	err := s.initGrid(grid)
	if err != nil {
		t.Error(err)
	}

	s.initCandidates()
	if !s.Solve() {
		t.Error(fmt.Errorf("Sudoku not solved: %v", grid))
	}
	s.printGrid()
}

func TestGrid2(t *testing.T) {
	grid := "700600008800030000090000310006740005005806900400092100087000020000060009600008001"

	s := Sudoku{}

	err := s.initGrid(grid)
	if err != nil {
		t.Error(err)
	}

	s.initCandidates()
	if s.Solve() {
		t.Error(fmt.Errorf("Sudoku solved: %v", grid))
	}
	s.printGrid()
}

func TestGrid3(t *testing.T) {
	grid := "014600300050000007090840100000400800600050009007009000008016030300000010009008570"

	s := Sudoku{}

	err := s.initGrid(grid)
	if err != nil {
		t.Error(err)
	}

	s.initCandidates()
	if s.Solve() {
		t.Error(fmt.Errorf("Sudoku solved: %v", grid))
	}
	s.printGrid()
}

func TestGrid4(t *testing.T) {
	grid := "000921003009000060000000500080403006007000800500700040003000000020000700800195000"

	s := Sudoku{}

	err := s.initGrid(grid)
	if err != nil {
		t.Error(err)
	}

	s.initCandidates()
	if s.Solve() {
		t.Error(fmt.Errorf("Sudoku solved: %v", grid))
	}
	s.printGrid()
}
