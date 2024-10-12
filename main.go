package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/tobiashort/i3-alt-tab/i3"
)

var previous i3.Workspace
var current i3.Workspace

func main() {
	current = i3.FocusedWorkspace()
	previous = current

	go i3.OnWorkspaceChange(
		// focus
		func(curr i3.Workspace, prev i3.Workspace) {
			previous = prev
			current = curr
		},
	)

	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGUSR1)

	for {
		select {
		case _ = <-signals:
			i3.FocusWorkspace(previous)
		}
	}
}
