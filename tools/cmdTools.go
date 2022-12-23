package tools

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
)

func useUtf8(script string) string {
	return "chcp 65001 && " + script
}

func newCmd(scripts ...string) *exec.Cmd {
	script := strings.Join(scripts, " && ")
	return exec.Command("cmd", "/c", script)
}

func ExecCommand(scripts ...string) ([]string, error) {
	var buf bytes.Buffer
	cmd := newCmd(scripts...)
	w := NewBig5ToUtf8Writer(&buf)
	cmd.Stdout = w
	cmd.Stderr = w
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	_, err = cmd.Process.Wait()
	str := buf.String()
	fmt.Println(str)
	return strings.Split(str, "\r"), err
}

func ExecCommandPrint(script string) ([]string, error) {
	cmd := newCmd(script)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	lines := make([]string, 0, 10)
	buf := bufio.NewReader(NewBig5ToUtf8Reader(stdout))
	for {
		n, err := buf.ReadString('\n')
		lines = append(lines, n)
		// fmt.Printf("n = %v err = %v b = %v\n", n, err, b)
		fmt.Print(n)
		if err == io.EOF {
			break
		}
	}

	err = cmd.Wait()
	// _, err = cmd.Process.Wait()
	return lines, err
}
