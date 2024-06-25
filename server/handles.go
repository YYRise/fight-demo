package main

import (
	"fmt"
	"net"
	"strconv"
)

func HandleAvailableRoom(conn net.Conn) {
	ids := RM.GetAvailableRoomIDs()
	if len(ids) == 0 {
		fmt.Fprintln(conn, "没有可用的房间，请先创建房间")
	} else {
		fmt.Fprintf(conn, "请选择房间%v加入或创建房间\n", ids)
	}
}

func HandleCreateRoom(conn net.Conn) {
	room := RM.CreateRoom()
	fmt.Fprintf(conn, "创建房间成功，房间ID为%v\n", room.ID)
}

func HandleJoinRoom(conn net.Conn, parts []string) {
	if len(parts) < 3 {
		fmt.Fprintln(conn, "参数错误")
		return
	}
	roomID, _ := strconv.Atoi(parts[1])
	playerName := parts[2]
	err := RM.JoinRoom(roomID, playerName, conn)
	if err != nil {
		fmt.Fprintf(conn, "加入房间失败：%v\n", err)
		return
	} else {
		fmt.Fprintln(conn, "加入房间成功")
	}
	RM.CheckFight(roomID)
}
