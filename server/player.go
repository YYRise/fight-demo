package main

import (
	"fmt"
	"math/rand"
	"net"
)

type Player struct {
	Name    string
	Health  int
	Attack  int
	Defense int
	IsFirst bool
	conn    net.Conn
}

func NewPlayer(playerName string, conn net.Conn) *Player {
	return &Player{
		Name:    playerName,
		Health:  100,
		Attack:  rand.Intn(50) + 50, // 50-100
		Defense: rand.Intn(40) + 10, // 10-50
		IsFirst: false,
		conn:    conn,
	}
}

func (p *Player) String() string {
	return fmt.Sprintf("名字：%s，血量：%d，攻击：%d，防御：%d", p.Name, p.Health, p.Attack, p.Defense)
}

func (p *Player) NoticeJoinRoom(player *Player) {
	p.SendMsg(fmt.Sprintf("玩家：%s加入房间\n", player.Name))
}

func (p *Player) NoticePlayerData(player *Player) {
	p.SendMsg(fmt.Sprintf("你的属性:【%s】\n 对方属性:【%s】\n", p.String(), player.String()))
}

func (p *Player) SendMsg(msg string) {
	fmt.Fprintln(p.conn, msg)
}

func (p *Player) UpdateHealth(damage int) {
	p.Health -= damage
}

func (p *Player) IsAlive() bool {
	return p.Health > 0
}
