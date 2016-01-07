using System;
using System.IO;

public class Login {
	public string Name;
	public string Password;

	public void Serialize(BinaryWriter buffer) {
		buffer.Write((Int32)this.Name.Length);
		buffer.Write(this.Name);
		buffer.Write((Int32)this.Password.Length);
		buffer.Write(this.Password);
	}

	public void Deserialize(BinaryReader buffer) {
		int l0_1 = buffer.ReadInt32();
		byte[] temp0_1 = buffer.ReadBytes(l0_1);
		this.Name = System.Text.Encoding.UTF8.GetString(temp0_1);
		int l1_1 = buffer.ReadInt32();
		byte[] temp1_1 = buffer.ReadBytes(l1_1);
		this.Password = System.Text.Encoding.UTF8.GetString(temp1_1);
	}
}

public class ListGames {

	public void Serialize(BinaryWriter buffer) {
	}

	public void Deserialize(BinaryReader buffer) {
	}
}

public class ListGamesResponse {
    public UInt32[] IDs;
    public string[] Names;

    public void Serialize(BinaryWriter buffer)
    {
        buffer.Write((Int32)this.IDs.Length);
        for (int v2 = 0; v2 < this.IDs.Length; v2++)
        {
            buffer.Write(this.IDs[v2]);
        }
        buffer.Write((Int32)this.Names.Length);
        for (int v2 = 0; v2 < this.Names.Length; v2++)
        {
            buffer.Write((Int32)this.Names[v2].Length);
            buffer.Write(this.Names[v2]);
        }
    }

    public void Deserialize(BinaryReader buffer)
    {
        int l0_1 = buffer.ReadInt32();
        this.IDs = new UInt32[l0_1];
        for (int v2 = 0; v2 < l0_1; v2++)
        {
            this.IDs[v2] = buffer.ReadUInt32();
        }
        int l1_1 = buffer.ReadInt32();
        this.Names = new string[l1_1];
        for (int v2 = 0; v2 < l1_1; v2++)
        {
            int l0_2 = buffer.ReadInt32();
            byte[] temp0_2 = buffer.ReadBytes(l0_2);
            this.Names[v2] = System.Text.Encoding.UTF8.GetString(temp0_2);
        }
    }
}

public class CreateGame {
	public string Name;

	public void Serialize(BinaryWriter buffer) {
		buffer.Write((Int32)this.Name.Length);
		buffer.Write(this.Name);
	}

	public void Deserialize(BinaryReader buffer) {
		int l0_1 = buffer.ReadInt32();
		byte[] temp0_1 = buffer.ReadBytes(l0_1);
		this.Name = System.Text.Encoding.UTF8.GetString(temp0_1);
	}
}

public class JoinGame {
	public UInt32 ID;
	public UInt32 CharID;

	public void Serialize(BinaryWriter buffer) {
		buffer.Write(this.ID);
		buffer.Write(this.CharID);
	}

	public void Deserialize(BinaryReader buffer) {
		this.ID = buffer.ReadUInt32();
		this.CharID = buffer.ReadUInt32();
	}
}

public class CreateCharacter {
	public string Name;
	public byte Kit;

	public void Serialize(BinaryWriter buffer) {
		buffer.Write((Int32)this.Name.Length);
		buffer.Write(this.Name);
		buffer.Write(this.Kit);
	}

	public void Deserialize(BinaryReader buffer) {
		int l0_1 = buffer.ReadInt32();
		byte[] temp0_1 = buffer.ReadBytes(l0_1);
		this.Name = System.Text.Encoding.UTF8.GetString(temp0_1);
		this.Kit = buffer.ReadByte();
	}
}

public class DeleteCharacter {
	public Int32 ID;

	public void Serialize(BinaryWriter buffer) {
		buffer.Write(this.ID);
	}

	public void Deserialize(BinaryReader buffer) {
		this.ID = buffer.ReadInt32();
	}
}

public class MapLoaded {
	public byte[][] Tiles;
	public Entity[] Entities;

	public void Serialize(BinaryWriter buffer) {
		buffer.Write((Int32)this.Tiles.Length);
		for (int v2 = 0; v2 < this.Tiles.Length; v2++) {
			buffer.Write((Int32)this.Tiles[v2].Length);
			for (int v3 = 0; v3 < this.Tiles[v2].Length; v3++) {
				buffer.Write(this.Tiles[v2][v3]);
			}
		}
		buffer.Write((Int32)this.Entities.Length);
		for (int v2 = 0; v2 < this.Entities.Length; v2++) {
			this.Entities[v2].Serialize(buffer);
		}
	}

	public void Deserialize(BinaryReader buffer) {
		int l0_1 = buffer.ReadInt32();
		this.Tiles = new byte[l0_1][];
		for (int v2 = 0; v2 < l0_1; v2++) {
			int l0_2 = buffer.ReadInt32();
			this.Tiles[v2] = new byte[l0_2];
			for (int v3 = 0; v3 < l0_2; v3++) {
				this.Tiles[v2][v3] = buffer.ReadByte();
			}
		}
		int l1_1 = buffer.ReadInt32();
		this.Entities = new Entity[l1_1];
		for (int v2 = 0; v2 < l1_1; v2++) {
			this.Entities[v2] = new Entity();
			this.Entities[v2].Deserialize(buffer);
		}
	}
}

public class Entity {
	public UInt32 ID;
	public byte HealthPercent;
	public Int32 X;
	public Int32 Y;

	public void Serialize(BinaryWriter buffer) {
		buffer.Write(this.ID);
		buffer.Write(this.HealthPercent);
		buffer.Write(this.X);
		buffer.Write(this.Y);
	}

	public void Deserialize(BinaryReader buffer) {
		this.ID = buffer.ReadUInt32();
		this.HealthPercent = buffer.ReadByte();
		this.X = buffer.ReadInt32();
		this.Y = buffer.ReadInt32();
	}
}

public class EntityMove {
	public byte Direction;

	public void Serialize(BinaryWriter buffer) {
		buffer.Write(this.Direction);
	}

	public void Deserialize(BinaryReader buffer) {
		this.Direction = buffer.ReadByte();
	}
}

public class UseAbility {
	public Int32 AbilityID;
	public UInt32 Target;

	public void Serialize(BinaryWriter buffer) {
		buffer.Write(this.AbilityID);
		buffer.Write(this.Target);
	}

	public void Deserialize(BinaryReader buffer) {
		this.AbilityID = buffer.ReadInt32();
		this.Target = buffer.ReadUInt32();
	}
}

public class AbilityResult {
	public Entity Target;
	public Int32 Damage;
	public byte State;

	public void Serialize(BinaryWriter buffer) {
		this.Target.Serialize(buffer);
		buffer.Write(this.Damage);
		buffer.Write(this.State);
	}

	public void Deserialize(BinaryReader buffer) {
		this.Target = new Entity();
		this.Target.Deserialize(buffer);
		this.Damage = buffer.ReadInt32();
		this.State = buffer.ReadByte();
	}
}

public class EndGame {

	public void Serialize(BinaryWriter buffer) {
	}

	public void Deserialize(BinaryReader buffer) {
	}
}