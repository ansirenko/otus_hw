package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const Indent = 20

type ProgressBar struct {
	graph string
	barWidth int
	currentValue chan int
	maxValue int
	doneCh chan interface{}
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

func (pb *ProgressBar) getCurrentTerminalWidth() (int, error) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	b := strings.ReplaceAll(string(out), "\n", "")
	sizes := strings.Split(b, " ")
	width, err := strconv.Atoi(sizes[1])
	if err != nil {
		return 0, err
	}
	return width, nil
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
				clearer := strings.Repeat(" ", int(pb.barWidth + Indent))
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
