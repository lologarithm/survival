package messages

import (
	"bytes"
	"encoding/binary"
)

type LoginMessage struct {
	Name string
	Password string
}

func (m *LoginMessage) Serialize(buffer *bytes.Buffer) {
	binary.Write(buffer, binary.LittleEndian, len(m.Name))
	buffer.WriteString(m.Name)
	binary.Write(buffer, binary.LittleEndian, len(m.Password))
	buffer.WriteString(m.Password)
}

func (m *LoginMessage) Deserialize(buffer *bytes.Buffer) {
	var l0_1 int
	binary.Read(buffer, binary.LittleEndian, &l0_1)
	temp0_1 := make([]byte, l0_1)
	buffer.Read(temp0_1)
	m.Name = string(temp0_1)
	var l1_1 int
	binary.Read(buffer, binary.LittleEndian, &l1_1)
	temp1_1 := make([]byte, l1_1)
	buffer.Read(temp1_1)
	m.Password = string(temp1_1)
}

type ListGames struct {
}

func (m *ListGames) Serialize(buffer *bytes.Buffer) {
}

func (m *ListGames) Deserialize(buffer *bytes.Buffer) {
}

type ListGamesResponse struct {
	IDs []uint32
	Names []string
}

func (m *ListGamesResponse) Serialize(buffer *bytes.Buffer) {
	binary.Write(buffer, binary.LittleEndian, len(m.IDs))
	for _, v2 := range m.IDs {
		binary.Write(buffer, binary.LittleEndian, v2)
	}
	binary.Write(buffer, binary.LittleEndian, len(m.Names))
	for _, v2 := range m.Names {
		binary.Write(buffer, binary.LittleEndian, len(v2))
		buffer.WriteString(v2)
	}
}

func (m *ListGamesResponse) Deserialize(buffer *bytes.Buffer) {
	var l0_1 int
	binary.Read(buffer, binary.LittleEndian, &l0_1)
	m.IDs = make([]uint32, l0_1)
	for i := 0; i < l0_1; i++ {
		binary.Read(buffer, binary.LittleEndian, &m.IDs[i])
	}
	var l1_1 int
	binary.Read(buffer, binary.LittleEndian, &l1_1)
	m.Names = make([]string, l1_1)
	for i := 0; i < l1_1; i++ {
		var l0_2 int
		binary.Read(buffer, binary.LittleEndian, &l0_2)
		temp0_2 := make([]byte, l0_2)
		buffer.Read(temp0_2)
		m.Names[i] = string(temp0_2)
	}
}

type CreateGame struct {
	Name string
}

func (m *CreateGame) Serialize(buffer *bytes.Buffer) {
	binary.Write(buffer, binary.LittleEndian, len(m.Name))
	buffer.WriteString(m.Name)
}

func (m *CreateGame) Deserialize(buffer *bytes.Buffer) {
	var l0_1 int
	binary.Read(buffer, binary.LittleEndian, &l0_1)
	temp0_1 := make([]byte, l0_1)
	buffer.Read(temp0_1)
	m.Name = string(temp0_1)
}

type JoinGame struct {
	ID uint32
	CharID uint32
}

func (m *JoinGame) Serialize(buffer *bytes.Buffer) {
	binary.Write(buffer, binary.LittleEndian, m.ID)
	binary.Write(buffer, binary.LittleEndian, m.CharID)
}

func (m *JoinGame) Deserialize(buffer *bytes.Buffer) {
	binary.Read(buffer, binary.LittleEndian, &m.ID)
	binary.Read(buffer, binary.LittleEndian, &m.CharID)
}

type CreateCharacter struct {
	Name string
	Kit byte
}

func (m *CreateCharacter) Serialize(buffer *bytes.Buffer) {
	binary.Write(buffer, binary.LittleEndian, len(m.Name))
	buffer.WriteString(m.Name)
	buffer.WriteByte(m.Kit)
}

func (m *CreateCharacter) Deserialize(buffer *bytes.Buffer) {
	var l0_1 int
	binary.Read(buffer, binary.LittleEndian, &l0_1)
	temp0_1 := make([]byte, l0_1)
	buffer.Read(temp0_1)
	m.Name = string(temp0_1)
	m.Kit, _ = buffer.ReadByte()
}

type DeleteCharacter struct {
	ID int32
}

func (m *DeleteCharacter) Serialize(buffer *bytes.Buffer) {
	binary.Write(buffer, binary.LittleEndian, m.ID)
}

func (m *DeleteCharacter) Deserialize(buffer *bytes.Buffer) {
	binary.Read(buffer, binary.LittleEndian, &m.ID)
}

type MapLoaded struct {
	Tiles [][]byte
	Entities []*Entity
}

func (m *MapLoaded) Serialize(buffer *bytes.Buffer) {
	binary.Write(buffer, binary.LittleEndian, len(m.Tiles))
	for _, v2 := range m.Tiles {
		binary.Write(buffer, binary.LittleEndian, len(v2))
		for _, v3 := range v2 {
			buffer.WriteByte(v3)
		}
	}
	binary.Write(buffer, binary.LittleEndian, len(m.Entities))
	for _, v2 := range m.Entities {
		v2.Serialize(buffer)
	}
}

func (m *MapLoaded) Deserialize(buffer *bytes.Buffer) {
	var l0_1 int
	binary.Read(buffer, binary.LittleEndian, &l0_1)
	m.Tiles = make([][]byte, l0_1)
	for i := 0; i < l0_1; i++ {
		var l0_2 int
		binary.Read(buffer, binary.LittleEndian, &l0_2)
		m.Tiles[i] = make([]byte, l0_2)
		for i := 0; i < l0_2; i++ {
			m.Tiles[i][i], _ = buffer.ReadByte()
		}
	}
	var l1_1 int
	binary.Read(buffer, binary.LittleEndian, &l1_1)
	m.Entities = make([]*Entity, l1_1)
	for i := 0; i < l1_1; i++ {
		m.Entities[i] = new(Entity)
		m.Entities[i].Deserialize(buffer)
	}
}

type Entity struct {
	ID uint32
	HealthPercent byte
	X int32
	Y int32
}

func (m *Entity) Serialize(buffer *bytes.Buffer) {
	binary.Write(buffer, binary.LittleEndian, m.ID)
	buffer.WriteByte(m.HealthPercent)
	binary.Write(buffer, binary.LittleEndian, m.X)
	binary.Write(buffer, binary.LittleEndian, m.Y)
}

func (m *Entity) Deserialize(buffer *bytes.Buffer) {
	binary.Read(buffer, binary.LittleEndian, &m.ID)
	m.HealthPercent, _ = buffer.ReadByte()
	binary.Read(buffer, binary.LittleEndian, &m.X)
	binary.Read(buffer, binary.LittleEndian, &m.Y)
}

type EntityMove struct {
	Direction byte
}

func (m *EntityMove) Serialize(buffer *bytes.Buffer) {
	buffer.WriteByte(m.Direction)
}

func (m *EntityMove) Deserialize(buffer *bytes.Buffer) {
	m.Direction, _ = buffer.ReadByte()
}

type UseAbility struct {
	AbilityID int32
	Target uint32
}

func (m *UseAbility) Serialize(buffer *bytes.Buffer) {
	binary.Write(buffer, binary.LittleEndian, m.AbilityID)
	binary.Write(buffer, binary.LittleEndian, m.Target)
}

func (m *UseAbility) Deserialize(buffer *bytes.Buffer) {
	binary.Read(buffer, binary.LittleEndian, &m.AbilityID)
	binary.Read(buffer, binary.LittleEndian, &m.Target)
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

type EndGame struct {
}

func (m *EndGame) Serialize(buffer *bytes.Buffer) {
}

func (m *EndGame) Deserialize(buffer *bytes.Buffer) {
}

