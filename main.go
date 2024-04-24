package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	processes := make([]*exec.Cmd, 0)

loop:
	for {
		time.Sleep(100 * time.Millisecond)
		fmt.Print(">> ")
		inputString, _ := reader.ReadString('\n')
		inputString = strings.TrimSuffix(inputString, "\n")

		args := strings.Split(inputString, " ")
		command := args[0]
		switch command {
		case "exit":
			break loop
		case "cd":
			if len(args) > 1 {
				err := os.Chdir(args[1])
				if err != nil {
					fmt.Println("cd:", err)
				}
			} else {
				fmt.Println("cd: пропущено значение")
			}
		default:
			go func() {
				stopCh := make(chan os.Signal, 1)
				signal.Notify(stopCh, syscall.SIGINT)

				cmd := exec.Command(command, args[1:]...)
				cmd.Stderr = os.Stderr
				cmd.Stdout = os.Stdout

				processes = append(processes, cmd)

				go func() {
					<-stopCh
					for _, pr := range processes {
						pr.Process.Kill()
					}

					processes = []*exec.Cmd{}
				}()

				err := cmd.Run()
				if err == nil {
					for prIndex, pr := range processes {
						if pr == cmd {
							processes = append(processes[:prIndex], processes[prIndex+1:]...)
							break
						}
					}
				}
			}()
		}
	}
}
