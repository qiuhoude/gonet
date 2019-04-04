package main

import (
	"fmt"
	"gonet/actor"
	"gonet/base"
	"gonet/message"
	"gonet/network"
)

type (
	EventProcess struct {
		actor.Actor

		Client *network.ClientSocket
		AccountId int64
		PlayerId int64
		AccountName string
		SimId int64
	}

	IEventProcess interface {
		actor.IActor
		LoginGame()
		LoginAccount()
		SendPacket(message.Message)
	}
)

func SendPacket(packet message.Message){
	buff := message.Encode(packet)
	buff = base.SetTcpEnd(buff)
	CLIENT.Send(buff)
}

func (this *EventProcess) SendPacket(packet message.Message){
	buff := message.Encode(packet)
	buff = base.SetTcpEnd(buff)
	this.Client.Send(buff)
}

func (this *EventProcess) PacketFunc(socketid int, buff []byte) bool {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("EventProcess PacketFunc", err)
		}
	}()

	packetId, data := message.Decode(buff)
	packet := message.GetPakcet(packetId)
	if packet == nil{
		return true
	}
	err := message.UnmarshalText(packet, data)
	if err == nil{
		bitstream := base.NewBitStream(make([]byte, 1024), 1024)
		if !message.GetMessagePacket(packet, bitstream) {
			return true
		}
		var io actor.CallIO
		io.Buff = bitstream.GetBuffer()
		io.SocketId = socketid
		this.Send(io)
		return true
	}

	return true
}

func (this *EventProcess) Init(num int) {
	this.Actor.Init(num)
	this.RegisterCall("W_C_SelectPlayerResponse", func(packet *message.W_C_SelectPlayerResponse) {
		this.AccountId = packet.AccountId
		nLen := len(packet.PlayerData)
		//fmt.Println(len(packet.PlayerData), this.AccountId, packet.PlayerData)
		if nLen == 0{
			packet1 := &message.C_W_CreatePlayerRequest{MessageBase:*message.BuildMessageBase( this.AccountId, int(message.SERVICE_WORLDSERVER), "C_W_CreatePlayerRequest"),
				PlayerName:"我是大坏蛋",
				Sex:(0),}
			this.SendPacket(packet1)
		}else{
			this.PlayerId = packet.PlayerData[0].PlayerID
			this.LoginGame()
		}
	})

	this.RegisterCall("W_C_CreatePlayerResponse", func(packet *message.W_C_CreatePlayerResponse) {
		if packet.Error == 0 {
			this.PlayerId = packet.PlayerId
			this.LoginGame()
		}else{//创建失败

		}
	})

	this.RegisterCall("A_C_LoginRequest", func(packet *message.A_C_LoginRequest) {
		if packet.Error == base.ACCOUNT_NOEXIST {
			packet1 := &message.C_A_RegisterRequest{MessageBase:*message.BuildMessageBase( 0, int(message.SERVICE_ACCOUNTSERVER), "C_A_RegisterRequest"),
				AccountName: packet.AccountName, SocketId: 0}
			this.SendPacket(packet1)
		}
	})

	this.RegisterCall("A_C_RegisterResponse", func(packet *message.A_C_RegisterResponse) {
		//注册失败
		if packet.Error != 0 {
		}
	})

	this.RegisterCall("W_C_ChatMessage", func(packet *message.W_C_ChatMessage) {
		fmt.Println("收到【", packet.SenderName, "】发送的消息[", packet.Message+"]")
	})

	this.Actor.Start()
}

func (this *EventProcess)  LoginGame(){
	packet1 := &message.C_W_Game_LoginRequset{MessageBase:*message.BuildMessageBase( this.AccountId, int(message.SERVICE_WORLDSERVER), "C_W_Game_LoginRequset"),
		PlayerId:this.PlayerId,}
	this.SendPacket(packet1)
}

var(
	id int
)

func (this *EventProcess)  LoginAccount() {
	id++
	this.AccountName = fmt.Sprintf("test%d", id)
	//this.AccountName = fmt.Sprintf("test%d", base.RAND().RandI(0, 7000))
	packet1 := &message.C_A_LoginRequest{MessageBase: *message.BuildMessageBase(0, int(message.SERVICE_ACCOUNTSERVER), "C_A_LoginRequest"),
		AccountName: this.AccountName, BuildNo: base.BUILD_NO, SocketId: 0}
	this.SendPacket(packet1)
}

var(
	PACKET *EventProcess
)