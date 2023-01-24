package i3

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/t-hg/i3-alt-tab/must"
)

type I3 struct {
	conn          net.Conn
	prevWorkspace uint
	currWorkspace uint
}

type workspace struct {
	Num uint
}

type workspaceEvent struct {
	Change  string
	Current workspace
}

func Connect() *I3 {
	i3Sock := os.Getenv("I3SOCK")
	fmt.Println("I3SOCK", i3Sock)
	conn := must.Do2(net.Dial("unix", i3Sock))
	i3 := I3{
		conn:          conn,
		prevWorkspace: 1,
		currWorkspace: 1,
	}
	go i3.listen()
	i3.subscribe("workspace")
	return &i3
}

func (i3 *I3) Close() error {
	return i3.conn.Close()
}

func (i3 *I3) SwapWorkspace() {
	i3.runCommand(fmt.Sprintf("workspace number %d", i3.prevWorkspace))
}

func (i3 *I3) send(msgType uint32, msg []byte) {
	bs := []byte{}
	bs = append(bs, []byte("i3-ipc")...)
	bs = binary.LittleEndian.AppendUint32(bs, uint32(len(msg)))
	bs = binary.LittleEndian.AppendUint32(bs, msgType)
	bs = append(bs, msg...)
	_ = must.Do2(i3.conn.Write(bs))
}

func (i3 *I3) runCommand(command string) {
	i3.send(0, []byte(command))
}

func (i3 *I3) subscribe(events ...string) {
	i3.send(2, must.Do2(json.Marshal(events)))
}

func (i3 *I3) listen() {
	for {
		header := make([]byte, 6)
		_ = must.Do2(i3.conn.Read(header))
		if bytes.Compare(header, []byte("i3-ipc")) != 0 {
			continue
		}
		length := make([]byte, 4)
		_ = must.Do2(i3.conn.Read(length))
		msgType := make([]byte, 4)
		_ = must.Do2(i3.conn.Read(msgType))
		msg := make([]byte, binary.LittleEndian.Uint32(length))
		_ = must.Do2(i3.conn.Read(msg))
		event := workspaceEvent{}
		err := json.Unmarshal(msg, &event)
		if err != nil {
			continue
		}
		if event.Change == "focus" {
			i3.prevWorkspace = i3.currWorkspace
			i3.currWorkspace = event.Current.Num
			fmt.Println(i3.prevWorkspace, "->", i3.currWorkspace)
		}
	}
}
