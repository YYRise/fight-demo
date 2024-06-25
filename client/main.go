package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("error dialing:", err)
		return
	}
	defer conn.Close()
	handleResponse(conn)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		readStandardInput(conn)
	}()

	wg.Wait()
}

func readStandardInput(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("输入命令 (CREATE--创建房间, JOIN <roomID> <name>)--加入房间: ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		parts := strings.Split(text, " ")

		if len(parts) < 1 {
			continue
		}

		command := strings.ToUpper(parts[0])
		switch command {
		case "CREATE":
			conn.Write([]byte(command + "\n"))
			handleResponse(conn)
		case "JOIN":
			if len(parts) < 3 {
				fmt.Println("JOIN 参数错误")
				continue
			}
			roomID, _ := strconv.Atoi(parts[1])
			name := parts[2]
			conn.Write([]byte(fmt.Sprintf("%s %d %s\n", command, roomID, name)))
			handleJoin(conn)
		default:
			fmt.Println("无效命令")
		}
	}
}

func handleResponse(conn net.Conn) {
	data := make([]byte, 1024)
	n, err := conn.Read(data)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}
	response := strings.TrimSpace(string(data[:n]))
	fmt.Println("Received:", response)
}

func handleJoin(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("服务器断开连接……")
			} else {
				fmt.Println("读取指令时发生错误:", err)
			}
		}
		text = strings.TrimSpace(text)
		fmt.Println(text)
	}
}
