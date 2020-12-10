package main

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

const (
	Indent               = 20
	DefaultTerminalWidth = 100
)

type ProgressBar struct {
	graph        string
	barWidth     int
	currentValue chan int
	maxValue     int
	doneCh       chan interface{}
}

func NewProgressBar(maxValue int) (*ProgressBar, error) {
	pb := &ProgressBar{}
	pb.graph = "█"

	if maxValue == 0 {
		return nil, fmt.Errorf("unreachable max value for progress bar")
	}
	pb.maxValue = maxValue

	maxWidth, err := pb.getCurrentTerminalWidth()

	if err != nil {
		return nil, err
	}
	pb.barWidth = maxWidth - Indent

	pb.currentValue = make(chan int)
	return pb, nil
}

const (
	TiocgwinszOsx = 1074295912
)

type window struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func (pb *ProgressBar) getCurrentTerminalWidth() (int, error) {
	w := new(window)
	tio := syscall.TIOCGWINSZ
	if runtime.GOOS == "darwin" {
		tio = TiocgwinszOsx
	}
	res, _, err := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(tio),
		uintptr(unsafe.Pointer(w)),
	)
	if err == syscall.ENOTTY {
		return DefaultTerminalWidth, nil
	}
	if int(res) == -1 {
		return 0, err
	}
	return int(w.Col), nil
}

func (pb *ProgressBar) Play() {
	go func() {
		for {
			select {
			case <-pb.doneCh:
				return
			case cv := <-pb.currentValue:
				donePerc := 100 * cv / pb.maxValue

				var doneRepeat int
				if donePerc != 0 {
					doneRepeat = pb.barWidth * donePerc / 100
				}
				if donePerc == 100 {
					doneRepeat = pb.barWidth
				}

				progress := strings.Repeat(pb.graph, doneRepeat)
				clearer := strings.Repeat(" ", pb.barWidth+Indent)
				fmt.Printf("\r%s\r", clearer)
				fmt.Printf("\r[%-"+strconv.Itoa(pb.barWidth)+"s]%1d%% %2d/%d", progress, donePerc, cv, pb.maxValue)

				if cv == pb.maxValue {
					fmt.Println("\nDone!")
					return
				}
			}
		}
	}()
}
