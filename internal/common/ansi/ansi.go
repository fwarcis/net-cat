package ansi

import (
	"fmt"
)

func CarriageOffset(x, y int) string {
	xDir, yDir := 'C', 'B'
	col, row := x, y
	if x < 0 {
		xDir = 'D'
		col = -x
	}
	if y < 0 {
		yDir = 'A'
		row = -y
	}
	return fmt.Sprintf(
		"\033[%d%c\033[%d%c",
		row, yDir, col, xDir)
}
