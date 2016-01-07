using System;
using System.IO;
using System.Text;

public class LoginMessage {
	public string Name
	public string Password

	public Serialize(BinaryWriter writer) {
		writer.Write(this.Name.length);
		writer.Write(this.Name);
		writer.Write(this.Password.length);
		writer.Write(this.Password);
	}

	public void Deserialize(BinaryReader buffer) {
		var l0_1 int
	binary.Read(buffer, binary.LittleEndian, &l0_1)
	temp0_1 := make([]byte, l0_1)
	buffer.Read(temp0_1)
	this.Name = string(temp0_1)
		var l1_1 int
	binary.Read(buffer, binary.LittleEndian, &l1_1)
	temp1_1 := make([]byte, l1_1)
	buffer.Read(temp1_1)
	this.Password = string(temp1_1)
	}
}

public class ListGames {

	public Serialize(BinaryWriter writer) {
	}

	public void Deserialize(BinaryReader buffer) {
	}
}

public class ListGamesResponse {
	public []uint32 IDs
	public []string Names

	public Serialize(BinaryWriter writer) {
		writer.Write(this.IDs.length);
		for ( int i=0; i < this.IDs.length; i++) {
			writer.Write(this.IDs[i]);
		}
		writer.Write(this.Names.length);
		for ( int i=0; i < this.Names.length; i++) {
			writer.Write(this.Names[i].length);
			writer.Write(this.Names[i]);
		}
	}

	public void Deserialize(BinaryReader buffer) {
		var l0_1 int
	binary.Read(buffer, binary.LittleEndian, &l0_1)
	this.IDs = make([]uint32, l0_1)
	for i := 0; i < l0_1; i++ {
			this.IDs[i] = buffer.ReadUInt32();
	}
		var l1_1 int
	binary.Read(buffer, binary.LittleEndian, &l1_1)
	this.Names = make([]string, l1_1)
	for i := 0; i < l1_1; i++ {
			var l0_2 int
		binary.Read(buffer, binary.LittleEndian, &l0_2)
		temp0_2 := make([]byte, l0_2)
		buffer.Read(temp0_2)
		this.Names[i] = string(temp0_2)
	}
	}
}

public class CreateGame {
	public string Name

	public Serialize(BinaryWriter writer) {
		writer.Write(this.Name.length);
		writer.Write(this.Name);
	}

	public void Deserialize(BinaryReader buffer) {
		var l0_1 int
	binary.Read(buffer, binary.LittleEndian, &l0_1)
	temp0_1 := make([]byte, l0_1)
	buffer.Read(temp0_1)
	this.Name = string(temp0_1)
	}
}

public class JoinGame {
	public uint32 ID
	public uint32 CharID

	public Serialize(BinaryWriter writer) {
		writer.Write(this.ID);
		writer.Write(this.CharID);
	}

	public void Deserialize(BinaryReader buffer) {
		this.ID = buffer.ReadUInt32();
		this.CharID = buffer.ReadUInt32();
	}
}

public class CreateCharacter {
	public string Name
	public byte Kit

	public Serialize(BinaryWriter writer) {
		writer.Write(this.Name.length);
		writer.Write(this.Name);
		writer.Write(this.Kit);
	}

	public void Deserialize(BinaryReader buffer) {
		var l0_1 int
	binary.Read(buffer, binary.LittleEndian, &l0_1)
	temp0_1 := make([]byte, l0_1)
	buffer.Read(temp0_1)
	this.Name = string(temp0_1)
		this.Kit = buffer.ReadByte()
	}
}

public class DeleteCharacter {
	public int32 ID

	public Serialize(BinaryWriter writer) {
		writer.Write(this.ID);
	}

	public void Deserialize(BinaryReader buffer) {
		this.ID = buffer.ReadInt32();
	}
}

public class MapLoaded {
	public [][]byte Tiles
	public []*Entity Entities

	public Serialize(BinaryWriter writer) {
		writer.Write(this.Tiles.length);
		for ( int i=0; i < this.Tiles.length; i++) {
			writer.Write(this.Tiles[i].length);
			for ( int i=0; i < this.Tiles[i].length; i++) {
				writer.Write(this.Tiles[i][i]);
			}
		}
		writer.Write(this.Entities.length);
		for ( int i=0; i < this.Entities.length; i++) {
			this.Entities[i].Serialize(buffer)
		}
	}

	public void Deserialize(BinaryReader buffer) {
		var l0_1 int
	binary.Read(buffer, binary.LittleEndian, &l0_1)
	this.Tiles = make([][]byte, l0_1)
	for i := 0; i < l0_1; i++ {
			var l0_2 int
		binary.Read(buffer, binary.LittleEndian, &l0_2)
		this.Tiles[i] = make([]byte, l0_2)
		for i := 0; i < l0_2; i++ {
				this.Tiles[i][i] = buffer.ReadByte()
		}
	}
		var l1_1 int
	binary.Read(buffer, binary.LittleEndian, &l1_1)
	this.Entities = make([]*Entity, l1_1)
	for i := 0; i < l1_1; i++ {
			this.Entities[i] = new(Entity)
		this.Entities[i].Deserialize(buffer)
	}
	}
}

public class Entity {
	public uint32 ID
	public byte HealthPercent
	public int32 X
	public int32 Y

	public Serialize(BinaryWriter writer) {
		writer.Write(this.ID);
		writer.Write(this.HealthPercent);
		writer.Write(this.X);
		writer.Write(this.Y);
	}

	public void Deserialize(BinaryReader buffer) {
		this.ID = buffer.ReadUInt32();
		this.HealthPercent = buffer.ReadByte()
		this.X = buffer.ReadInt32();
		this.Y = buffer.ReadInt32();
	}
}

public class EntityMove {
	public byte Direction

	public Serialize(BinaryWriter writer) {
		writer.Write(this.Direction);
	}

	public void Deserialize(BinaryReader buffer) {
		this.Direction = buffer.ReadByte()
	}
}

public class UseAbility {
	public int32 AbilityID
	public uint32 Target

	public Serialize(BinaryWriter writer) {
		writer.Write(this.AbilityID);
		writer.Write(this.Target);
	}

	public void Deserialize(BinaryReader buffer) {
		this.AbilityID = buffer.ReadInt32();
		this.Target = buffer.ReadUInt32();
	}
}

public class AbilityResult {
	public *Entity Target
	public int32 Damage
	public byte State

	public Serialize(BinaryWriter writer) {
		this.Target.Serialize(buffer)
		writer.Write(this.Damage);
		writer.Write(this.State);
	}

	public void Deserialize(BinaryReader buffer) {
		this.Target = new(Entity)
	this.Target.Deserialize(buffer)
		this.Damage = buffer.ReadInt32();
		this.State = buffer.ReadByte()
	}
}

public class EndGame {

	public Serialize(BinaryWriter writer) {
	}

	public void Deserialize(BinaryReader buffer) {
	}
}

