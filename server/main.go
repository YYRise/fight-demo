package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func handleConnection(conn net.Conn) {
	fmt.Println("客户端连接成功")
	defer conn.Close()
	HandleAvailableRoom(conn)

	for {
		data := make([]byte, 1024)
		n, err := conn.Read(data)
		if err != nil {
			if err == io.EOF {
				// 客户端断开连接
				fmt.Println("客户端断开连接")
			} else {
				fmt.Println("读取指令时发生错误:", err)
			}
			break
		}

		msg := strings.TrimSpace(string(data[:n]))
		parts := strings.Split(msg, " ")
		if len(parts) < 1 {
			continue
		}
		command := parts[0]
		fmt.Println("Received:", msg)
		switch command {
		case "CREATE":
			go HandleCreateRoom(conn)
		case "JOIN":
			go HandleJoinRoom(conn, parts)
		default:
			conn.Write([]byte("Unknown command\n"))
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Listening on :8080")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err)
			continue
		}
		go handleConnection(conn)
	}
}
