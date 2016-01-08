using System;
using System.IO;
using System.Text;

interface INet {
	void Serialize(BinaryWriter buffer);
	void Deserialize(BinaryReader buffer);
}

enum MsgType : byte {Unknown=0,CreateAccount=1,CreateAccountResponse=2,Login=3,LoginResponse=4,CreateCharacter=5,DeleteCharacter=6,Character=7,ListGames=8,ListGamesResponse=9,CreateGame=10,JoinGame=11,MapLoaded=12,Entity=13,EntityMove=14,UseAbility=15,AbilityResult=16,EndGame=17}

static class Messages {
// ParseNetMessage accepts input of raw bytes from a NetMessage. Parses and returns a Net message.
public static INet Parse(byte msgType, byte[] content) {
	INet msg = null;
	MsgType mt = (MsgType)msgType;
	switch (mt)
	{
		case MsgType.CreateAccount:
			msg = new CreateAccount();
			break;
		case MsgType.CreateAccountResponse:
			msg = new CreateAccountResponse();
			break;
		case MsgType.Login:
			msg = new Login();
			break;
		case MsgType.LoginResponse:
			msg = new LoginResponse();
			break;
		case MsgType.CreateCharacter:
			msg = new CreateCharacter();
			break;
		case MsgType.DeleteCharacter:
			msg = new DeleteCharacter();
			break;
		case MsgType.Character:
			msg = new Character();
			break;
		case MsgType.ListGames:
			msg = new ListGames();
			break;
		case MsgType.ListGamesResponse:
			msg = new ListGamesResponse();
			break;
		case MsgType.CreateGame:
			msg = new CreateGame();
			break;
		case MsgType.JoinGame:
			msg = new JoinGame();
			break;
		case MsgType.MapLoaded:
			msg = new MapLoaded();
			break;
		case MsgType.Entity:
			msg = new Entity();
			break;
		case MsgType.EntityMove:
			msg = new EntityMove();
			break;
		case MsgType.UseAbility:
			msg = new UseAbility();
			break;
		case MsgType.AbilityResult:
			msg = new AbilityResult();
			break;
		case MsgType.EndGame:
			msg = new EndGame();
			break;
	}
	MemoryStream ms = new MemoryStream(content);
	msg.Deserialize(new BinaryReader(ms));
	return msg;
}
}

public class CreateAccount : INet {
	public string Name;
	public string Password;

	public void Serialize(BinaryWriter buffer) {
		buffer.Write((Int32)this.Name.Length);
		buffer.Write(System.Text.Encoding.UTF8.GetBytes(this.Name));
		buffer.Write((Int32)this.Password.Length);
		buffer.Write(System.Text.Encoding.UTF8.GetBytes(this.Password));
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

public class CreateAccountResponse : INet {
	public UInt32 AccountID;
	public string Name;

	public void Serialize(BinaryWriter buffer) {
		buffer.Write(this.AccountID);
		buffer.Write((Int32)this.Name.Length);
		buffer.Write(System.Text.Encoding.UTF8.GetBytes(this.Name));
	}

	public void Deserialize(BinaryReader buffer) {
		this.AccountID = buffer.ReadUInt32();
		int l1_1 = buffer.ReadInt32();
		byte[] temp1_1 = buffer.ReadBytes(l1_1);
		this.Name = System.Text.Encoding.UTF8.GetString(temp1_1);
	}
}

public class Login : INet {
	public string Name;
	public string Password;

	public void Serialize(BinaryWriter buffer) {
		buffer.Write((Int32)this.Name.Length);
		buffer.Write(System.Text.Encoding.UTF8.GetBytes(this.Name));
		buffer.Write((Int32)this.Password.Length);
		buffer.Write(System.Text.Encoding.UTF8.GetBytes(this.Password));
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

public class LoginResponse : INet {
	public byte Success;
	public string Name;
	public UInt32 AccountID;
	public Character[] Characters;

	public void Serialize(BinaryWriter buffer) {
		buffer.Write(this.Success);
		buffer.Write((Int32)this.Name.Length);
		buffer.Write(System.Text.Encoding.UTF8.GetBytes(this.Name));
		buffer.Write(this.AccountID);
		buffer.Write((Int32)this.Characters.Length);
		for (int v2 = 0; v2 < this.Characters.Length; v2++) {
			this.Characters[v2].Serialize(buffer);
		}
	}

	public void Deserialize(BinaryReader buffer) {
		this.Success = buffer.ReadByte();
		int l1_1 = buffer.ReadInt32();
		byte[] temp1_1 = buffer.ReadBytes(l1_1);
		this.Name = System.Text.Encoding.UTF8.GetString(temp1_1);
		this.AccountID = buffer.ReadUInt32();
		int l3_1 = buffer.ReadInt32();
		this.Characters = new Character[l3_1];
		for (int v2 = 0; v2 < l3_1; v2++) {
			this.Characters[v2] = new Character();
			this.Characters[v2].Deserialize(buffer);
		}
	}
}

public class CreateCharacter : INet {
	public UInt32 AccountID;
	public string Name;
	public byte Kit;

	public void Serialize(BinaryWriter buffer) {
		buffer.Write(this.AccountID);
		buffer.Write((Int32)this.Name.Length);
		buffer.Write(System.Text.Encoding.UTF8.GetBytes(this.Name));
		buffer.Write(this.Kit);
	}

	public void Deserialize(BinaryReader buffer) {
		this.AccountID = buffer.ReadUInt32();
		int l1_1 = buffer.ReadInt32();
		byte[] temp1_1 = buffer.ReadBytes(l1_1);
		this.Name = System.Text.Encoding.UTF8.GetString(temp1_1);
		this.Kit = buffer.ReadByte();
	}
}

public class DeleteCharacter : INet {
	public UInt32 ID;

	public void Serialize(BinaryWriter buffer) {
		buffer.Write(this.ID);
	}

	public void Deserialize(BinaryReader buffer) {
		this.ID = buffer.ReadUInt32();
	}
}

public class Character : INet {
	public UInt32 ID;
	public string Name;

	public void Serialize(BinaryWriter buffer) {
		buffer.Write(this.ID);
		buffer.Write((Int32)this.Name.Length);
		buffer.Write(System.Text.Encoding.UTF8.GetBytes(this.Name));
	}

	public void Deserialize(BinaryReader buffer) {
		this.ID = buffer.ReadUInt32();
		int l1_1 = buffer.ReadInt32();
		byte[] temp1_1 = buffer.ReadBytes(l1_1);
		this.Name = System.Text.Encoding.UTF8.GetString(temp1_1);
	}
}

public class ListGames : INet {

	public void Serialize(BinaryWriter buffer) {
	}

	public void Deserialize(BinaryReader buffer) {
	}
}

public class ListGamesResponse : INet {
	public UInt32[] IDs;
	public string[] Names;

	public void Serialize(BinaryWriter buffer) {
		buffer.Write((Int32)this.IDs.Length);
		for (int v2 = 0; v2 < this.IDs.Length; v2++) {
			buffer.Write(this.IDs[v2]);
		}
		buffer.Write((Int32)this.Names.Length);
		for (int v2 = 0; v2 < this.Names.Length; v2++) {
			buffer.Write((Int32)this.Names[v2].Length);
			buffer.Write(System.Text.Encoding.UTF8.GetBytes(this.Names[v2]));
		}
	}

	public void Deserialize(BinaryReader buffer) {
		int l0_1 = buffer.ReadInt32();
		this.IDs = new UInt32[l0_1];
		for (int v2 = 0; v2 < l0_1; v2++) {
			this.IDs[v2] = buffer.ReadUInt32();
		}
		int l1_1 = buffer.ReadInt32();
		this.Names = new string[l1_1];
		for (int v2 = 0; v2 < l1_1; v2++) {
			int l0_2 = buffer.ReadInt32();
			byte[] temp0_2 = buffer.ReadBytes(l0_2);
			this.Names[v2] = System.Text.Encoding.UTF8.GetString(temp0_2);
		}
	}
}

public class CreateGame : INet {
	public string Name;

	public void Serialize(BinaryWriter buffer) {
		buffer.Write((Int32)this.Name.Length);
		buffer.Write(System.Text.Encoding.UTF8.GetBytes(this.Name));
	}

	public void Deserialize(BinaryReader buffer) {
		int l0_1 = buffer.ReadInt32();
		byte[] temp0_1 = buffer.ReadBytes(l0_1);
		this.Name = System.Text.Encoding.UTF8.GetString(temp0_1);
	}
}

public class JoinGame : INet {
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

public class MapLoaded : INet {
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

public class Entity : INet {
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

public class EntityMove : INet {
	public byte Direction;

	public void Serialize(BinaryWriter buffer) {
		buffer.Write(this.Direction);
	}

	public void Deserialize(BinaryReader buffer) {
		this.Direction = buffer.ReadByte();
	}
}

public class UseAbility : INet {
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

public class AbilityResult : INet {
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

public class EndGame : INet {

	public void Serialize(BinaryWriter buffer) {
	}

	public void Deserialize(BinaryReader buffer) {
	}
}

