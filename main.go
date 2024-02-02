package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"
)

type (
	WINNER  string
	FgColor string
	BgColor string
)

// SETTINGS
const (
	// WARNING: DO NOT CHANGE
	WINDOWS               string = "windows"
	WINDOWS_CLEAR_COMMAND string = "cls"
	LINUX                 string = "linux"
	LINUX_CLEAR_COMMAND   string = "clear"

	X      string = "X"   // single byte to draw as X
	O      string = "O"   // single byte to draw as O
	BOARD  string = "|"   // chars to draw as board border
	SPACER string = "---" // chars to draw as board spacers

	// WARNING: DO NOT CHANGE
	ANSI_RESET  FgColor = "\u001B[0m"
	ANSI_RED    FgColor = "\u001B[31m"
	ANSI_GREEN  FgColor = "\u001B[32m"
	ANSI_YELLOW FgColor = "\u001B[33m"
	ANSI_BLUE   FgColor = "\u001B[34m"
	ANSI_PURPLE FgColor = "\u001B[35m"
	ANSI_CYAN   FgColor = "\u001B[36m"

	// WARNING: DO NOT CHANGE
	ANSI_BLACK_BACKGROUND  BgColor = "\u001B[40m"
	ANSI_RED_BACKGROUND    BgColor = "\u001B[41m"
	ANSI_GREEN_BACKGROUND  BgColor = "\u001B[42m"
	ANSI_YELLOW_BACKGROUND BgColor = "\u001B[43m"
	ANSI_BLUE_BACKGROUND   BgColor = "\u001B[44m"
	ANSI_PURPLE_BACKGROUND BgColor = "\u001B[45m"
	ANSI_CYAN_BACKGROUND   BgColor = "\u001B[46m"
	ANSI_WHITE_BACKGROUND  BgColor = "\u001B[47m"

	// duration until game clears screen and exits when someone wins
	END_GAME_PAUSE_DURATION time.Duration = 4 * time.Second

	// colors to draw different entities in
	X_DRAW_COLOR          FgColor = ANSI_RED
	O_DRAW_COLOR          FgColor = ANSI_BLUE
	EMPTY_CELL_DRAW_COLOR FgColor = ANSI_PURPLE
	WINNER_DRAW_COLOR     FgColor = ANSI_GREEN
	BOARD_DRAW_COLOR      FgColor = ANSI_YELLOW
	SPACER_DRAW_COLOR     FgColor = ANSI_YELLOW
	INFO_DRAW_COLOR       FgColor = ANSI_CYAN
	ERROR_DRAW_COLOR      FgColor = ANSI_RED

	// WARNING: CHANGING THE FOLLOWING SETTINGS REQUIRE CODE CHANGES IN ORDER
	// FOR IT TO WORK PROPERLY
	WINNER_NONE WINNER = "N"
	WINNER_DRAW WINNER = "D"
	WINNER_X    WINNER = "X"
	WINNER_O    WINNER = "O"
)

var clear map[string]func()

func DefineClear() {
	clear = make(map[string]func())
	clear[LINUX] = func() {
		cmd := exec.Command(LINUX_CLEAR_COMMAND)
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear[WINDOWS] = func() {
		cmd := exec.Command("cmd", "/c", WINDOWS_CLEAR_COMMAND)
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func colorCell(cell string, index int, winnerTiles []int, size int) FgColor {
	if size > 0 {
		for i := 0; i < 3; i += 1 {
			if winnerTiles[i] == index {
				return WINNER_DRAW_COLOR
			}
		}
	}

	switch cell {
	case X:
		{
			return X_DRAW_COLOR
		}
	case O:
		{
			return O_DRAW_COLOR
		}
	default:
		{
			return EMPTY_CELL_DRAW_COLOR
		}
	}
}

func printBoard(board *[]string, winnerTiles []int, size int) {
	_board := *board
	index := 0

	for i := 0; i < 7; i += 1 {
		if i%2 == 0 {
			fmt.Printf("%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s\n",
				BOARD_DRAW_COLOR, BOARD, ANSI_RESET,
				SPACER_DRAW_COLOR, SPACER, ANSI_RESET,
				BOARD_DRAW_COLOR, BOARD, ANSI_RESET,
				SPACER_DRAW_COLOR, SPACER, ANSI_RESET,
				BOARD_DRAW_COLOR, BOARD, ANSI_RESET,
				SPACER_DRAW_COLOR, SPACER, ANSI_RESET,
				BOARD_DRAW_COLOR, BOARD, ANSI_RESET)
		} else {
			fmt.Printf("%s%s%s %s%s%s %s%s%s %s%s%s %s%s%s %s%s%s %s%s%s\n",
				BOARD_DRAW_COLOR, BOARD, ANSI_RESET,
				colorCell(_board[index], index, winnerTiles, size), _board[index], ANSI_RESET,
				BOARD_DRAW_COLOR, BOARD, ANSI_RESET,
				colorCell(_board[index+1], index+1, winnerTiles, size), _board[index+1], ANSI_RESET,
				BOARD_DRAW_COLOR, BOARD, ANSI_RESET,
				colorCell(_board[index+2], index+2, winnerTiles, size), _board[index+2], ANSI_RESET,
				BOARD_DRAW_COLOR, BOARD, ANSI_RESET)
			index += 3
		}
	}
}

func checkWinner(board *[]string) (WINNER, []int, int) {
	_board := *board
	var line string
	emptyCellCounter := 0
	winnerTiles := make([]int, 3)
	size := 0

	for i := 0; i < 9; i += 1 {
		switch i {
		case 0:
			{
				line = fmt.Sprintf("%s%s%s", _board[0], _board[1], _board[2])
				winnerTiles[0] = 0
				winnerTiles[1] = 1
				winnerTiles[2] = 2
			}
		case 1:
			{
				line = fmt.Sprintf("%s%s%s", _board[3], _board[4], _board[5])
				winnerTiles[0] = 3
				winnerTiles[1] = 4
				winnerTiles[2] = 5
			}
		case 2:
			{
				line = fmt.Sprintf("%s%s%s", _board[6], _board[7], _board[8])
				winnerTiles[0] = 6
				winnerTiles[1] = 7
				winnerTiles[2] = 8
			}
		case 3:
			{
				line = fmt.Sprintf("%s%s%s", _board[0], _board[3], _board[6])
				winnerTiles[0] = 0
				winnerTiles[1] = 3
				winnerTiles[2] = 6
			}
		case 4:
			{
				line = fmt.Sprintf("%s%s%s", _board[1], _board[4], _board[7])
				winnerTiles[0] = 1
				winnerTiles[1] = 4
				winnerTiles[2] = 7
			}
		case 5:
			{
				line = fmt.Sprintf("%s%s%s", _board[2], _board[5], _board[8])
				winnerTiles[0] = 2
				winnerTiles[1] = 5
				winnerTiles[2] = 8
			}
		case 6:
			{
				line = fmt.Sprintf("%s%s%s", _board[0], _board[4], _board[8])
				winnerTiles[0] = 0
				winnerTiles[1] = 4
				winnerTiles[2] = 8
			}
		case 7:
			{
				line = fmt.Sprintf("%s%s%s", _board[2], _board[4], _board[6])
				winnerTiles[0] = 2
				winnerTiles[1] = 4
				winnerTiles[2] = 6
			}
		default:
			{
				line = ""
				winnerTiles[0] = -1
				winnerTiles[1] = -1
				winnerTiles[2] = -1
			}
		}

		switch line {
		case X + X + X:
			{
				size = 3
				return WINNER_X, winnerTiles, size
			}
		case O + O + O:
			{
				size = 3
				return WINNER_O, winnerTiles, size
			}
		}

		if _board[i] == fmt.Sprintf("%d", i+1) {
			emptyCellCounter += 1
		}
	}

	if emptyCellCounter == 0 {
		return WINNER_DRAW, winnerTiles, size
	} else {
		return WINNER_NONE, winnerTiles, size
	}
}

func main() {
	DefineClear()
	clear, clearOk := clear[runtime.GOOS]
	if !clearOk {
		fmt.Printf("%sERROR: clear console for platform not implemented.\nContinuing without it...\n%s", ERROR_DRAW_COLOR, ANSI_RESET)
	}

	winner := WINNER_NONE
	turn := X
	board := make([]string, 9)
	winnerTiles := make([]int, 3)
	size := 0

	for i := 0; i < 9; i += 1 {
		board[i] = fmt.Sprintf("%d", i+1)
	}

	if clearOk {
		clear()
	}
	printBoard(&board, winnerTiles, size)

	for winner == WINNER_NONE {
		fmt.Printf("%sEnter integer to place %s into: ", INFO_DRAW_COLOR, turn)
		var input int
		_, err := fmt.Scanf("%d", &input)
		fmt.Printf("%s", ANSI_RESET)
		if err != nil {
			fmt.Printf("%sERROR: could not read from standard in!\n%s", ERROR_DRAW_COLOR, ANSI_RESET)
			continue
		}

		if input < 1 || input > 9 {
			fmt.Printf("%sERROR: input out of bounds!\n%s", ERROR_DRAW_COLOR, ANSI_RESET)
			continue
		}
		if board[input-1] != fmt.Sprintf("%d", input) {
			fmt.Printf("%sERROR: cell already taken!\n%s", ERROR_DRAW_COLOR, ANSI_RESET)
			continue
		}

		board[input-1] = turn
		switch turn {
		case X:
			{
				turn = O
			}
		case O:
			{
				turn = X
			}
		}

		winner, winnerTiles, size = checkWinner(&board)
		if clearOk {
			clear()
		}
		printBoard(&board, winnerTiles, size)
	}

	fmt.Printf("%sWINNER: %s\n%s", WINNER_DRAW_COLOR, winner, ANSI_RESET)
	time.Sleep(END_GAME_PAUSE_DURATION)
	if clearOk {
		clear()
	}
}
