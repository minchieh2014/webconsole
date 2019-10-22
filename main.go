
package main

import (
    //"fmt"
    //"net"
    //"io"
    "log"
    //"bufio"
	//"errors"
	"os"
	"os/exec"
	//"strconv"
	"strings"
	"runtime"
	"golang.org/x/net/websocket"
)

func getServerAddr() string {
	if len(os.Args) != 2 {
		//log.Println("USEAGE: websocket {serveraddr}")
		//os.exit(-1)
		log.Fatal("USEAGE: websocket {serveraddr}")
	}
	
	return os.Args[1]
}

func closePreCmd(c *exec.Cmd) {
	if c != nil {
		if c.Process != nil {
			c.Process.Kill()
		}
	}
}

func main() {
	log.Println("run at", runtime.GOOS)
	
	url := getServerAddr()
	
	ws, err := websocket.Dial(url, "", "websocket://")
	if err != nil { log.Fatal(err) }
	log.Println(url)
	var n int
	var recvBytes = make([]byte, 1024*10)
	var cmd *exec.Cmd
	for {
		if n, err = ws.Read(recvBytes); err != nil { log.Fatal(err) }
		closePreCmd(cmd)
		if n < 1 { continue }
		msg := string(recvBytes[:n])
		log.Println("1", msg, len(msg))
		msg = strings.TrimSuffix(msg, "\n")
		msg = strings.TrimSuffix(msg, "\r")
		log.Println("2", msg, len(msg))
		arrCommandStr := strings.Fields(msg)
		if len(arrCommandStr) == 0 { continue }
		switch arrCommandStr[0] {
		case "exit":
			return
		case "quit":
			return
		}

		if runtime.GOOS == "windows" {
			//cmd = exec.Command("cmd", "/c", "chcp 65001")
			//cmd.Run()
			//cmd.Wait()
			cmd = exec.Command("cmd", "/c", msg)
		} else {
			cmd = exec.Command("sh", "-c", msg)
		}
		cmd.Stderr = ws
		cmd.Stdout = ws
		go cmd.Run()
	}	
}