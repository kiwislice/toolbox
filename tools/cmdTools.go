// tools package 包含各種工具函式
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

// useUtf8 函式會將傳入的 script 前面加上 `chcp 65001 &&`，
// 這是為了在 windows 的 cmd 中，將 code page 切換為 utf-8，
// 避免中文亂碼的問題
func useUtf8(script string) string {
	return "chcp 65001 && " + script
}

// newCmd 函式會將傳入的多個 script 字串，用 `&&` 串接起來，
// 並回傳一個 `*exec.Cmd` 物件
func newCmd(scripts ...string) *exec.Cmd {
	script := strings.Join(scripts, " && ")
	return exec.Command("cmd", "/c", script)
}

// ExecCommand 函式會執行傳入的多個 script，
// 並將 stdout 與 stderr 的輸出，轉換為 utf-8 後，
// 回傳一個 string slice (以 `\r` 分割)
func ExecCommand(scripts ...string) ([]string, error) {
	var buf bytes.Buffer
	cmd := newCmd(scripts...)
	// 將 stdout 與 stderr 導向一個可以將 big5 轉為 utf-8 的 writer
	w := NewBig5ToUtf8Writer(&buf)
	cmd.Stdout = w
	cmd.Stderr = w
	// 開始執行 command
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	// 等待 command 執行完畢
	_, err = cmd.Process.Wait()
	// 將 buffer 轉為 string
	str := buf.String()
	fmt.Println(str)
	// 將 string 以 `\r` 分割為 string slice
	return strings.Split(str, "\r"), err
}

// ExecCommandPrint 函式會執行傳入的 script，
// 並將 stdout 的輸出，即時的印在 console 上，
// 同時也會將輸出的每一行，存到一個 string slice 中並回傳
func ExecCommandPrint(script string) ([]string, error) {
	cmd := newCmd(script)
	// 取得 stdout 的 pipe
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	// 開始執行 command
	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	// 建立一個 string slice，用來存放輸出的每一行
	lines := make([]string, 0, 10)
	// 建立一個 reader，用來讀取 stdout 的輸出，並將 big5 轉為 utf-8
	buf := bufio.NewReader(NewBig5ToUtf8Reader(stdout))
	for {
		// 讀取一行
		n, err := buf.ReadString('\n')
		lines = append(lines, n)
		// fmt.Printf("n = %v err = %v b = %v\n", n, err, b)
		// 即時的印在 console 上
		fmt.Print(n)
		// 如果讀到 EOF，就跳出迴圈
		if err == io.EOF {
			break
		}
	}

	// 等待 command 執行完畢
	err = cmd.Wait()
	// _, err = cmd.Process.Wait()
	return lines, err
}
