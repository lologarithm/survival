package messages

import (
	"bytes"
	"encoding/binary"
	"log"
)

type Net interface {
	Serialize(*bytes.Buffer)
	Deserialize(*bytes.Buffer)
	Len() int
}

type MessageType uint16

const (
	UnknownMsgType MessageType = iota
	AckMsgType
	MultipartMsgType
	ConnectedMsgType
	DisconnectedMsgType
	CreateAcctMsgType
	CreateAcctRespMsgType
	LoginMsgType
	LoginRespMsgType
	CharacterMsgType
	ListGamesMsgType
	ListGamesRespMsgType
	CreateGameMsgType
	CreateGameRespMsgType
	JoinGameMsgType
	GameConnectedMsgType
	EntityMsgType
	MovePlayerMsgType
	UseAbilityMsgType
	AbilityResultMsgType
	EndGameMsgType
)

// ParseNetMessage accepts input of raw bytes from a NetMessage. Parses and returns a Net message.
func ParseNetMessage(packet Packet, content []byte) Net {
	var msg Net
	switch packet.Frame.MsgType {
	case MultipartMsgType:
		msg = &Multipart{}
	case ConnectedMsgType:
		msg = &Connected{}
	case DisconnectedMsgType:
		msg = &Disconnected{}
	case CreateAcctMsgType:
		msg = &CreateAcct{}
	case CreateAcctRespMsgType:
		msg = &CreateAcctResp{}
	case LoginMsgType:
		msg = &Login{}
	case LoginRespMsgType:
		msg = &LoginResp{}
	case CharacterMsgType:
		msg = &Character{}
	case ListGamesMsgType:
		msg = &ListGames{}
	case ListGamesRespMsgType:
		msg = &ListGamesResp{}
	case CreateGameMsgType:
		msg = &CreateGame{}
	case CreateGameRespMsgType:
		msg = &CreateGameResp{}
	case JoinGameMsgType:
		msg = &JoinGame{}
	case GameConnectedMsgType:
		msg = &GameConnected{}
	case EntityMsgType:
		msg = &Entity{}
	case MovePlayerMsgType:
		msg = &MovePlayer{}
	case UseAbilityMsgType:
		msg = &UseAbility{}
	case AbilityResultMsgType:
		msg = &AbilityResult{}
	case EndGameMsgType:
		msg = &EndGame{}
	default:
		log.Printf("Unknown message type: %d", packet.Frame.MsgType)
		return nil
	}
	msg.Deserialize(bytes.NewBuffer(content))
	return msg
}

type Multipart struct {
	ID uint16
	GroupID uint32
	NumParts uint16
	Content []byte
}

func (m *Multipart) Serialize(buffer *bytes.Buffer) {
	binary.Write(buffer, binary.LittleEndian, m.ID)
	binary.Write(buffer, binary.LittleEndian, m.GroupID)
	binary.Write(buffer, binary.LittleEndian, m.NumParts)
	binary.Write(buffer, binary.LittleEndian, int32(len(m.Content)))
	buffer.Write(m.Content)
}

func (m *Multipart) Deserialize(buffer *bytes.Buffer) {
	binary.Read(buffer, binary.LittleEndian, &m.ID)
	binary.Read(buffer, binary.LittleEndian, &m.GroupID)
	binary.Read(buffer, binary.LittleEndian, &m.NumParts)
	var l3_1 int32
	binary.Read(buffer, binary.LittleEndian, &l3_1)
	m.Content = make([]byte, l3_1)
	for i := 0; i < int(l3_1); i++ {
		m.Content[i], _ = buffer.ReadByte()
	}
}

func (m *Multipart) Len() int {
	mylen := 0
	mylen += 2
	mylen += 4
	mylen += 2
	mylen += 4 + len(m.Content)
	return mylen
}

type Connected struct {
}

func (m *Connected) Serialize(buffer *bytes.Buffer) {
}

func (m *Connected) Deserialize(buffer *bytes.Buffer) {
}

func (m *Connected) Len() int {
	mylen := 0
	return mylen
}

type Disconnected struct {
}

func (m *Disconnected) Serialize(buffer *bytes.Buffer) {
}

func (m *Disconnected) Deserialize(buffer *bytes.Buffer) {
}

func (m *Disconnected) Len() int {
	mylen := 0
	return mylen
}

type CreateAcct struct {
	Name string
	Password string
	CharName string
	DefaultKit byte
}

func (m *CreateAcct) Serialize(buffer *bytes.Buffer) {
	binary.Write(buffer, binary.LittleEndian, int32(len(m.Name)))
	buffer.WriteString(m.Name)
	binary.Write(buffer, binary.LittleEndian, int32(len(m.Password)))
	buffer.WriteString(m.Password)
	binary.Write(buffer, binary.LittleEndian, int32(len(m.CharName)))
	buffer.WriteString(m.CharName)
	buffer.WriteByte(m.DefaultKit)
}

func (m *CreateAcct) Deserialize(buffer *bytes.Buffer) {
	var l0_1 int32
	binary.Read(buffer, binary.LittleEndian, &l0_1)
	temp0_1 := make([]byte, l0_1)
	buffer.Read(temp0_1)
	m.Name = string(temp0_1)
	var l1_1 int32
	binary.Read(buffer, binary.LittleEndian, &l1_1)
	temp1_1 := make([]byte, l1_1)
	buffer.Read(temp1_1)
	m.Password = string(temp1_1)
	var l2_1 int32
	binary.Read(buffer, binary.LittleEndian, &l2_1)
	temp2_1 := make([]byte, l2_1)
	buffer.Read(temp2_1)
	m.CharName = string(temp2_1)
	m.DefaultKit, _ = buffer.ReadByte()
}

func (m *CreateAcct) Len() int {
	mylen := 0
	mylen += 4 + len(m.Name)
	mylen += 4 + len(m.Password)
	mylen += 4 + len(m.CharName)
	mylen += 1
	return mylen
}

type CreateAcctResp struct {
	AccountID uint32
	Name string
}

func (m *CreateAcctResp) Serialize(buffer *bytes.Buffer) {
	binary.Write(buffer, binary.LittleEndian, m.AccountID)
	binary.Write(buffer, binary.LittleEndian, int32(len(m.Name)))
	buffer.WriteString(m.Name)
}

func (m *CreateAcctResp) Deserialize(buffer *bytes.Buffer) {
	binary.Read(buffer, binary.LittleEndian, &m.AccountID)
	var l1_1 int32
	binary.Read(buffer, binary.LittleEndian, &l1_1)
	temp1_1 := make([]byte, l1_1)
	buffer.Read(temp1_1)
	m.Name = string(temp1_1)
}

func (m *CreateAcctResp) Len() int {
	mylen := 0
	mylen += 4
	mylen += 4 + len(m.Name)
	return mylen
}

type Login struct {
	Name string
	Password string
}

func (m *Login) Serialize(buffer *bytes.Buffer) {
	binary.Write(buffer, binary.LittleEndian, int32(len(m.Name)))
	buffer.WriteString(m.Name)
	binary.Write(buffer, binary.LittleEndian, int32(len(m.Password)))
	buffer.WriteString(m.Password)
}

func (m *Login) Deserialize(buffer *bytes.Buffer) {
	var l0_1 int32
	binary.Read(buffer, binary.LittleEndian, &l0_1)
	temp0_1 := make([]byte, l0_1)
	buffer.Read(temp0_1)
	m.Name = string(temp0_1)
	var l1_1 int32
	binary.Read(buffer, binary.LittleEndian, &l1_1)
	temp1_1 := make([]byte, l1_1)
	buffer.Read(temp1_1)
	m.Password = string(temp1_1)
}

func (m *Login) Len() int {
	mylen := 0
	mylen += 4 + len(m.Name)
	mylen += 4 + len(m.Password)
	return mylen
}

type LoginResp struct {
	Success byte
	Name string
	AccountID uint32
	Character *Character
}

func (m *LoginResp) Serialize(buffer *bytes.Buffer) {
	buffer.WriteByte(m.Success)
	binary.Write(buffer, binary.LittleEndian, int32(len(m.Name)))
	buffer.WriteString(m.Name)
	binary.Write(buffer, binary.LittleEndian, m.AccountID)
	m.Character.Serialize(buffer)
}

func (m *LoginResp) Deserialize(buffer *bytes.Buffer) {
	m.Success, _ = buffer.ReadByte()
	var l1_1 int32
	binary.Read(buffer, binary.LittleEndian, &l1_1)
	temp1_1 := make([]byte, l1_1)
	buffer.Read(temp1_1)
	m.Name = string(temp1_1)
	binary.Read(buffer, binary.LittleEndian, &m.AccountID)
	m.Character = new(Character)
	m.Character.Deserialize(buffer)
}

func (m *LoginResp) Len() int {
	mylen := 0
	mylen += 1
	mylen += 4 + len(m.Name)
	mylen += 4
	mylen += m.Character.Len()
	return mylen
}

type Character struct {
	ID uint32
	Name string
}

func (m *Character) Serialize(buffer *bytes.Buffer) {
	binary.Write(buffer, binary.LittleEndian, m.ID)
	binary.Write(buffer, binary.LittleEndian, int32(len(m.Name)))
	buffer.WriteString(m.Name)
}

func (m *Character) Deserialize(buffer *bytes.Buffer) {
	binary.Read(buffer, binary.LittleEndian, &m.ID)
	var l1_1 int32
	binary.Read(buffer, binary.LittleEndian, &l1_1)
	temp1_1 := make([]byte, l1_1)
	buffer.Read(temp1_1)
	m.Name = string(temp1_1)
}

func (m *Character) Len() int {
	mylen := 0
	mylen += 4
	mylen += 4 + len(m.Name)
	return mylen
}

type ListGames struct {
}

func (m *ListGames) Serialize(buffer *bytes.Buffer) {
}

func (m *ListGames) Deserialize(buffer *bytes.Buffer) {
}

func (m *ListGames) Len() int {
	mylen := 0
	return mylen
}

type ListGamesResp struct {
	IDs []uint32
	Names []string
}

func (m *ListGamesResp) Serialize(buffer *bytes.Buffer) {
	binary.Write(buffer, binary.LittleEndian, int32(len(m.IDs)))
	for _, v2 := range m.IDs {
		binary.Write(buffer, binary.LittleEndian, v2)
	}
	binary.Write(buffer, binary.LittleEndian, int32(len(m.Names)))
	for _, v2 := range m.Names {
		binary.Write(buffer, binary.LittleEndian, int32(len(v2)))
		buffer.WriteString(v2)
	}
}

func (m *ListGamesResp) Deserialize(buffer *bytes.Buffer) {
	var l0_1 int32
	binary.Read(buffer, binary.LittleEndian, &l0_1)
	m.IDs = make([]uint32, l0_1)
	for i := 0; i < int(l0_1); i++ {
		binary.Read(buffer, binary.LittleEndian, &m.IDs[i])
	}
	var l1_1 int32
	binary.Read(buffer, binary.LittleEndian, &l1_1)
	m.Names = make([]string, l1_1)
	for i := 0; i < int(l1_1); i++ {
		var l0_2 int32
		binary.Read(buffer, binary.LittleEndian, &l0_2)
		temp0_2 := make([]byte, l0_2)
		buffer.Read(temp0_2)
		m.Names[i] = string(temp0_2)
	}
}

func (m *ListGamesResp) Len() int {
	mylen := 0
	mylen += 4
	for _, v2 := range m.IDs {
	_ = v2
		mylen += 4
	}

	mylen += 4
	for _, v2 := range m.Names {
	_ = v2
		mylen += 4 + len(v2)
	}

	return mylen
}

type CreateGame struct {
	Name string
}

func (m *CreateGame) Serialize(buffer *bytes.Buffer) {
	binary.Write(buffer, binary.LittleEndian, int32(len(m.Name)))
	buffer.WriteString(m.Name)
}

func (m *CreateGame) Deserialize(buffer *bytes.Buffer) {
	var l0_1 int32
	binary.Read(buffer, binary.LittleEndian, &l0_1)
	temp0_1 := make([]byte, l0_1)
	buffer.Read(temp0_1)
	m.Name = string(temp0_1)
}

func (m *CreateGame) Len() int {
	mylen := 0
	mylen += 4 + len(m.Name)
	return mylen
}

type CreateGameResp struct {
	Name string
	Game *GameConnected
}

func (m *CreateGameResp) Serialize(buffer *bytes.Buffer) {
	binary.Write(buffer, binary.LittleEndian, int32(len(m.Name)))
	buffer.WriteString(m.Name)
	m.Game.Serialize(buffer)
}

func (m *CreateGameResp) Deserialize(buffer *bytes.Buffer) {
	var l0_1 int32
	binary.Read(buffer, binary.LittleEndian, &l0_1)
	temp0_1 := make([]byte, l0_1)
	buffer.Read(temp0_1)
	m.Name = string(temp0_1)
	m.Game = new(GameConnected)
	m.Game.Deserialize(buffer)
}

func (m *CreateGameResp) Len() int {
	mylen := 0
	mylen += 4 + len(m.Name)
	mylen += m.Game.Len()
	return mylen
}

type JoinGame struct {
	ID uint32
}

func (m *JoinGame) Serialize(buffer *bytes.Buffer) {
	binary.Write(buffer, binary.LittleEndian, m.ID)
}

func (m *JoinGame) Deserialize(buffer *bytes.Buffer) {
	binary.Read(buffer, binary.LittleEndian, &m.ID)
}

func (m *JoinGame) Len() int {
	mylen := 0
	mylen += 4
	return mylen
}

type GameConnected struct {
	ID uint32
	Seed uint64
	Entities []*Entity
}

func (m *GameConnected) Serialize(buffer *bytes.Buffer) {
	binary.Write(buffer, binary.LittleEndian, m.ID)
	binary.Write(buffer, binary.LittleEndian, m.Seed)
	binary.Write(buffer, binary.LittleEndian, int32(len(m.Entities)))
	for _, v2 := range m.Entities {
		v2.Serialize(buffer)
	}
}

func (m *GameConnected) Deserialize(buffer *bytes.Buffer) {
	binary.Read(buffer, binary.LittleEndian, &m.ID)
	binary.Read(buffer, binary.LittleEndian, &m.Seed)
	var l2_1 int32
	binary.Read(buffer, binary.LittleEndian, &l2_1)
	m.Entities = make([]*Entity, l2_1)
	for i := 0; i < int(l2_1); i++ {
		m.Entities[i] = new(Entity)
		m.Entities[i].Deserialize(buffer)
	}
}

func (m *GameConnected) Len() int {
	mylen := 0
	mylen += 4
	mylen += 8
	mylen += 4
	for _, v2 := range m.Entities {
	_ = v2
		mylen += v2.Len()
	}

	return mylen
}

type Entity struct {
	ID uint32
	EType uint16
	Seed uint64
	X int32
	Y int32
	Height int32
	Width int32
	Angle int16
	HealthPercent byte
}

func (m *Entity) Serialize(buffer *bytes.Buffer) {
	binary.Write(buffer, binary.LittleEndian, m.ID)
	binary.Write(buffer, binary.LittleEndian, m.EType)
	binary.Write(buffer, binary.LittleEndian, m.Seed)
	binary.Write(buffer, binary.LittleEndian, m.X)
	binary.Write(buffer, binary.LittleEndian, m.Y)
	binary.Write(buffer, binary.LittleEndian, m.Height)
	binary.Write(buffer, binary.LittleEndian, m.Width)
	binary.Write(buffer, binary.LittleEndian, m.Angle)
	buffer.WriteByte(m.HealthPercent)
}

func (m *Entity) Deserialize(buffer *bytes.Buffer) {
	binary.Read(buffer, binary.LittleEndian, &m.ID)
	binary.Read(buffer, binary.LittleEndian, &m.EType)
	binary.Read(buffer, binary.LittleEndian, &m.Seed)
	binary.Read(buffer, binary.LittleEndian, &m.X)
	binary.Read(buffer, binary.LittleEndian, &m.Y)
	binary.Read(buffer, binary.LittleEndian, &m.Height)
	binary.Read(buffer, binary.LittleEndian, &m.Width)
	binary.Read(buffer, binary.LittleEndian, &m.Angle)
	m.HealthPercent, _ = buffer.ReadByte()
}

func (m *Entity) Len() int {
	mylen := 0
	mylen += 4
	mylen += 2
	mylen += 8
	mylen += 4
	mylen += 4
	mylen += 4
	mylen += 4
	mylen += 2
	mylen += 1
	return mylen
}

type MovePlayer struct {
	EntityID uint32
	TickID uint32
	Direction uint16
}

func (m *MovePlayer) Serialize(buffer *bytes.Buffer) {
	binary.Write(buffer, binary.LittleEndian, m.EntityID)
	binary.Write(buffer, binary.LittleEndian, m.TickID)
	binary.Write(buffer, binary.LittleEndian, m.Direction)
}

func (m *MovePlayer) Deserialize(buffer *bytes.Buffer) {
	binary.Read(buffer, binary.LittleEndian, &m.EntityID)
	binary.Read(buffer, binary.LittleEndian, &m.TickID)
	binary.Read(buffer, binary.LittleEndian, &m.Direction)
}

func (m *MovePlayer) Len() int {
	mylen := 0
	mylen += 4
	mylen += 4
	mylen += 2
	return mylen
}

type UseAbility struct {
	EntityID uint32
	AbilityID uint32
	TickID uint32
	Target uint32
}

func (m *UseAbility) Serialize(buffer *bytes.Buffer) {
	binary.Write(buffer, binary.LittleEndian, m.EntityID)
	binary.Write(buffer, binary.LittleEndian, m.AbilityID)
	binary.Write(buffer, binary.LittleEndian, m.TickID)
	binary.Write(buffer, binary.LittleEndian, m.Target)
}

func (m *UseAbility) Deserialize(buffer *bytes.Buffer) {
	binary.Read(buffer, binary.LittleEndian, &m.EntityID)
	binary.Read(buffer, binary.LittleEndian, &m.AbilityID)
	binary.Read(buffer, binary.LittleEndian, &m.TickID)
	binary.Read(buffer, binary.LittleEndian, &m.Target)
}

func (m *UseAbility) Len() int {
	mylen := 0
	mylen += 4
	mylen += 4
	mylen += 4
	mylen += 4
	return mylen
}

type AbilityResult struct {
	Target *Entity
	Damage int32
	State byte
}

func (m *AbilityResult) Serialize(buffer *bytes.Buffer) {
	m.Target.Serialize(buffer)
	binary.Write(buffer, binary.LittleEndian, m.Damage)
	buffer.WriteByte(m.State)
}

func (m *AbilityResult) Deserialize(buffer *bytes.Buffer) {
	m.Target = new(Entity)
	m.Target.Deserialize(buffer)
	binary.Read(buffer, binary.LittleEndian, &m.Damage)
	m.State, _ = buffer.ReadByte()
}

func (m *AbilityResult) Len() int {
	mylen := 0
	mylen += m.Target.Len()
	mylen += 4
	mylen += 1
	return mylen
}

type EndGame struct {
}

func (m *EndGame) Serialize(buffer *bytes.Buffer) {
}

func (m *EndGame) Deserialize(buffer *bytes.Buffer) {
}

func (m *EndGame) Len() int {
	mylen := 0
	return mylen
}

