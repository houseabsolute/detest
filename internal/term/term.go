// +build !windows

package term

import (
	"os"

	"golang.org/x/sys/unix"
)

func Width() int {
	tty, err := os.Open("/dev/tty")
	if err != nil {
		return 0
	}
	ws, err := unix.IoctlGetWinsize(int(tty.Fd()), unix.TIOCGWINSZ)
	if err != nil {
		return 0
	}
	return int(ws.Col)
}
