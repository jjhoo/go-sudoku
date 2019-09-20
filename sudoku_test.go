package sudoku

import (
	"fmt"
	"testing"
)

func TestGrid1(t *testing.T) {
	grid := "000040700500780020070002006810007900460000051009600078900800010080064009002050000"

	s, err := NewSudoku(grid)
	if err != nil {
		t.Error(err)
	}

	if !s.Solve() {
		t.Error(fmt.Errorf("Sudoku not solved: %v", grid))
	}
	s.PrintGrid()
}

func TestGrid2(t *testing.T) {
	grid := "700600008800030000090000310006740005005806900400092100087000020000060009600008001"

	s, err := NewSudoku(grid)
	if err != nil {
		t.Error(err)
	}

	s.Solve()
	s.PrintGrid()
}

func TestGrid3(t *testing.T) {
	grid := "014600300050000007090840100000400800600050009007009000008016030300000010009008570"

	s, err := NewSudoku(grid)
	if err != nil {
		t.Error(err)
	}

	s.Solve()
	s.PrintGrid()
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
	s.PrintGrid()
}

func TestBadInput1(t *testing.T) {
	// Test invalid characters in input
	grid := "CAT21003009000060000000500080403006007000800500700040003000000020000700800195000"

	s := Sudoku{}

	err := s.initGrid(grid)
	if err == nil {
		t.Error(fmt.Errorf("initGrid was expected to fail"))
	}
}

func TestBadInput2(t *testing.T) {
	// Test short input
	grid := "0092100300900006000000050008040300600700080050070004000300000002000070080019500"

	s := Sudoku{}

	err := s.initGrid(grid)
	if err == nil {
		t.Error(fmt.Errorf("initGrid was expected to fail"))
	}
}

func TestGridString(t *testing.T) {
	grid := "000921003009000060000000500080403006007000800500700040003000000020000700800195000"

	s := Sudoku{}

	err := s.initGrid(grid)
	if err != nil {
		t.Error(err)
	}

	ngrid := s.GetGridString()

	if grid != ngrid {
		t.Error(fmt.Errorf("GetGridString returned different string"))
	}
}

func TestPrintGrid1(t *testing.T) {
	grid := "000921003009000060000000500080403006007000800500700040003000000020000700800195000"

	err := PrintGrid(grid)
	if err != nil {
		t.Error(err)
	}
}

func TestPrintGrid2(t *testing.T) {
	// Short grid length
	grid := "00092100300900006000000050008040300600700080050070004000300000002000070080019500"

	err := PrintGrid(grid)
	if err == nil {
		t.Error(err)
	}
}
