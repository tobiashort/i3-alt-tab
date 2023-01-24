package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/t-hg/i3-alt-tab/i3"
)

func main() {
	fmt.Println("Process", os.Getpid())
	i3 := i3.Connect()
	defer i3.Close()
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGUSR1)
	for {
		select {
		case sig := <-signals:
			if sig == syscall.SIGUSR1 {
				i3.SwapWorkspace()
			}
		}
	}
}
