package i3

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/tobiashort/i3-alt-tab/must"
)

type Workspace struct {
	Num     int
	Name    string
	Focused bool
}

func FocusedWorkspace() Workspace {
	cmd := exec.Command("i3-msg", "-t", "get_workspaces")
	data, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(data))
		panic(err)
	}

	var workspaces []Workspace
	must.Do(json.Unmarshal(data, &workspaces))

	for _, workspace := range workspaces {
		if workspace.Focused {
			return workspace
		}
	}

	panic("No workspace focused")
}

func OnWorkspaceChange(focus func(current Workspace, previous Workspace)) {
	cmd := exec.Command("i3-msg", "-t", "subscribe", "-m", "[\"workspace\"]")
	reader := must.Do2(cmd.StdoutPipe())
	buffered := bufio.NewReader(reader)

	go func() {
		for {
			var event struct {
				Change  string
				Current Workspace
				Old     Workspace
			}

			data := must.Do2(buffered.ReadBytes('\n'))
			must.Do(json.Unmarshal(data, &event))

			if event.Change == "focus" {
				event.Current.Focused = true
				focus(event.Current, event.Old)
			}
		}
	}()

	go func() {
		err := cmd.Run()
		panic(fmt.Sprintf("Stream closed. Error: %s", err))
	}()
}

func FocusWorkspace(workspace Workspace) {
	cmd := exec.Command("i3-msg", "workspace", workspace.Name)
	data, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(data))
		panic(err)
	}
}