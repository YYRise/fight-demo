package main

import (
	"fmt"
	"math/rand"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

var RM *RoomManager

type RoomManager struct {
	RoomID int
	Rooms  map[int]*Room
	sync.RWMutex
}

func init() {
	RM = &RoomManager{
		Rooms: make(map[int]*Room),
	}
}

func (r *RoomManager) CreateRoom() *Room {
	r.Lock()
	defer r.Unlock()
	r.RoomID++
	room := NewRoom(r.RoomID)
	r.Rooms[r.RoomID] = room
	return room
}

func (r *RoomManager) GetRoom(roomID int) *Room {
	r.RLock()
	defer r.RUnlock()
	room, ok := r.Rooms[roomID]
	if !ok {
		return nil
	}
	return room
}

func (r *RoomManager) GetAvailableRoomIDs() []int {
	r.RLock()
	defer r.RUnlock()
	var ids []int
	for id, room := range r.Rooms {
		if (room.Player1 == nil || room.Player2 == nil) && room.IsWaiting() {
			ids = append(ids, id)
		}
	}
	return ids
}

func (r *RoomManager) JoinRoom(roomID int, playerName string, conn net.Conn) error {
	room := r.GetRoom(roomID)
	if room == nil || !room.IsWaiting() {
		return fmt.Errorf("Room:%d not found or not waiting status", roomID)
	}

	newPlayer := NewPlayer(playerName, conn)
	err := room.Join(newPlayer)
	if err != nil {
		return err
	}
	return nil
}
func (r *RoomManager) CheckFight(roomID int) {
	room := r.GetRoom(roomID)
	if room == nil || !room.IsFull() {
		return
	}
	room.Start()
}

const (
	ROOM_STATUS_WAITING int32 = iota
	ROOM_STATUS_FULL
	ROOM_STATUS_FIGHTING
	ROOM_STATUS_FINISHED
)

type Room struct {
	ID      int
	Player1 *Player
	Player2 *Player
	sync.RWMutex
	winner *Player
	Status int32
	Round  int
}

func NewRoom(id int) *Room {
	room := &Room{
		ID:     id,
		Status: ROOM_STATUS_WAITING,
	}
	return room
}

func (r *Room) IsWaiting() bool {
	return atomic.LoadInt32(&r.Status) == ROOM_STATUS_WAITING
}

func (r *Room) IsFull() bool {
	return atomic.LoadInt32(&r.Status) == ROOM_STATUS_FULL
}

func (r *Room) Join(player *Player) error {
	r.Lock()
	defer r.Unlock()
	var other *Player
	if r.Player1 == nil {
		r.Player1, other = player, r.Player2
	} else if r.Player2 == nil {
		r.Player2, other = player, r.Player1
	} else {
		return fmt.Errorf("Room is already full")
	}
	if other != nil {
		r.Status = ROOM_STATUS_FULL
		other.NoticeJoinRoom(player)
		player.NoticePlayerData(other)
		other.NoticePlayerData(player)
	}
	return nil
}

func (r *Room) Start() {
	if !atomic.CompareAndSwapInt32(&r.Status, ROOM_STATUS_FULL, ROOM_STATUS_FIGHTING) {
		return
	}
	r.Player1.IsFirst = rand.Intn(2) == 1
	for {
		r.Round++
		attackPlayer, defendPlayer := r.Player1, r.Player2
		if !attackPlayer.IsFirst {
			attackPlayer, defendPlayer = defendPlayer, attackPlayer
		}

		damage := max(1, attackPlayer.Attack-defendPlayer.Defense)
		defendPlayer.Health -= damage
		msg := fmt.Sprintf("第[%d]回合：%s 攻击 %s, 造成 %d 点伤害", r.Round, attackPlayer.Name, defendPlayer.Name, damage)
		// 向两个玩家发送战斗数据
		attackPlayer.SendMsg(msg)
		defendPlayer.SendMsg(msg)

		if defendPlayer.IsAlive() {
			time.Sleep(1 * time.Second)
			// 切换回合
			attackPlayer.IsFirst = !attackPlayer.IsFirst
			defendPlayer.IsFirst = !defendPlayer.IsFirst
			continue
		}

		r.winner = attackPlayer
		r.Status = ROOM_STATUS_FINISHED
		break
	}
	msg := fmt.Sprintf("战斗结束：%s 获得胜利", r.winner.Name)
	r.Player1.SendMsg(msg)
	r.Player2.SendMsg(msg)
}
