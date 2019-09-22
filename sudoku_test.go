package sudoku_test

import (
	"gotest.tools/assert"
	"github.com/jjhoo/go-sudoku"

	"testing"
)

func solvableSudoku(t *testing.T, grid string) {
	s, err := sudoku.NewSudoku(grid)
	if err != nil {
		t.Error(err)
	}

	assert.Assert(t, s.Solve(), "Sudoku should have been solved: %v", grid)
	s.PrintGrid()
}

func unsolvableSudoku(t *testing.T, grid string) {
	s, err := sudoku.NewSudoku(grid)
	if err != nil {
		t.Error(err)
	}

	assert.Assert(t, !s.Solve(), "Sudoku should have not been solved: %v", grid)
	s.PrintGrid()
}

func TestGrid1(t *testing.T) {
	grid := "000040700500780020070002006810007900460000051009600078900800010080064009002050000"
	solvableSudoku(t, grid)
}

func TestGrid2(t *testing.T) {
	grid := "700600008800030000090000310006740005005806900400092100087000020000060009600008001"
	solvableSudoku(t, grid)
}

func TestGrid3(t *testing.T) {
	grid := "014600300050000007090840100000400800600050009007009000008016030300000010009008570"
	solvableSudoku(t, grid)
}

func TestGrid4(t *testing.T) {
	grid := "000921003009000060000000500080403006007000800500700040003000000020000700800195000"
	unsolvableSudoku(t, grid)
}

func TestGrid5(t *testing.T) {
	grid := "300000000970010000600583000200000900500621003008000005000435002000090056000000001"
	solvableSudoku(t, grid)
}

func TestBadInput1(t *testing.T) {
	// Test invalid characters in input
	grid := "CAT921003009000060000000500080403006007000800500700040003000000020000700800195000"

	_, err := sudoku.NewSudoku(grid)
	assert.Error(t, err, "Invalid rune 'C' in grid")
}

func TestBadInput2(t *testing.T) {
	// Test short input
	grid := "0092100300900006000000050008040300600700080050070004000300000002000070080019500"

	_, err := sudoku.NewSudoku(grid)
	assert.Error(t, err, "Grid has invalid size '79'")
}

func TestGridString(t *testing.T) {
	grid := "000921003009000060000000500080403006007000800500700040003000000020000700800195000"

	s, err := sudoku.NewSudoku(grid)
	assert.NilError(t, err, "Grid should have been ok: %v", grid)

	ngrid := s.GetGridString()
	assert.Equal(t, grid, ngrid, "GetGridString returned different string: %v != %v", grid, ngrid)
}

func TestPrintGrid1(t *testing.T) {
	grid := "000921003009000060000000500080403006007000800500700040003000000020000700800195000"

	err := sudoku.PrintGrid(grid)
	assert.NilError(t, err)
}

func TestPrintGrid2(t *testing.T) {
	// Short grid length
	grid := "00092100300900006000000050008040300600700080050070004000300000002000070080019500"

	err := sudoku.PrintGrid(grid)
	assert.Error(t, err, "Grid has invalid size '80'")
}
