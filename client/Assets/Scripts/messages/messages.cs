using System;
using System.IO;
using System.Text;

interface INet {
	void Serialize(BinaryWriter buffer);
	void Deserialize(BinaryReader buffer);
}

enum MsgType : ushort {Unknown=0,Ack=1,Multipart=2,Connected=3,Disconnected=4,CreateAcct=5,CreateAcctResp=6,Login=7,LoginResp=8,CreateChar=9,CreateCharResp=10,DeleteChar=11,Character=12,ListGames=13,ListGamesResp=14,CreateGame=15,CreateGameResp=16,JoinGame=17,GameConnected=18,Entity=19,EntityMove=20,UseAbility=21,AbilityResult=22,EndGame=23}

static class Messages {
// ParseNetMessage accepts input of raw bytes from a NetMessage. Parses and returns a Net message.
public static INet Parse(ushort msgType, byte[] content) {
	INet msg = null;
	MsgType mt = (MsgType)msgType;
	switch (mt)
	{
		case MsgType.Multipart:
			msg = new Multipart();
			break;
		case MsgType.Connected:
			msg = new Connected();
			break;
		case MsgType.Disconnected:
			msg = new Disconnected();
			break;
		case MsgType.CreateAcct:
			msg = new CreateAcct();
			break;
		case MsgType.CreateAcctResp:
			msg = new CreateAcctResp();
			break;
		case MsgType.Login:
			msg = new Login();
			break;
		case MsgType.LoginResp:
			msg = new LoginResp();
			break;
		case MsgType.CreateChar:
			msg = new CreateChar();
			break;
		case MsgType.CreateCharResp:
			msg = new CreateCharResp();
			break;
		case MsgType.DeleteChar:
			msg = new DeleteChar();
			break;
		case MsgType.Character:
			msg = new Character();
			break;
		case MsgType.ListGames:
			msg = new ListGames();
			break;
		case MsgType.ListGamesResp:
			msg = new ListGamesResp();
			break;
		case MsgType.CreateGame:
			msg = new CreateGame();
			break;
		case MsgType.CreateGameResp:
			msg = new CreateGameResp();
			break;
		case MsgType.JoinGame:
			msg = new JoinGame();
			break;
		case MsgType.GameConnected:
			msg = new GameConnected();
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

public class Multipart : INet {
	public ushort ID;
	public uint GroupID;
	public ushort NumParts;
	public byte[] Content;

	public void Serialize(BinaryWriter buffer) {
		buffer.Write(this.ID);
		buffer.Write(this.GroupID);
		buffer.Write(this.NumParts);
		buffer.Write((Int32)this.Content.Length);
		for (int v2 = 0; v2 < this.Content.Length; v2++) {
			buffer.Write(this.Content[v2]);
		}
	}

	public void Deserialize(BinaryReader buffer) {
		this.ID = buffer.ReadUInt16();
		this.GroupID = buffer.ReadUInt32();
		this.NumParts = buffer.ReadUInt16();
		int l3_1 = buffer.ReadInt32();
		this.Content = new byte[l3_1];
		for (int v2 = 0; v2 < l3_1; v2++) {
			this.Content[v2] = buffer.ReadByte();
		}
	}
}

public class Connected : INet {

	public void Serialize(BinaryWriter buffer) {
	}

	public void Deserialize(BinaryReader buffer) {
	}
}

public class Disconnected : INet {

	public void Serialize(BinaryWriter buffer) {
	}

	public void Deserialize(BinaryReader buffer) {
	}
}

public class CreateAcct : INet {
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

public class CreateAcctResp : INet {
	public uint AccountID;
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

public class LoginResp : INet {
	public byte Success;
	public string Name;
	public uint AccountID;
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

public class CreateChar : INet {
	public uint AccountID;
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

public class CreateCharResp : INet {
	public uint AccountID;
	public Character Character;

	public void Serialize(BinaryWriter buffer) {
		buffer.Write(this.AccountID);
		this.Character.Serialize(buffer);
	}

	public void Deserialize(BinaryReader buffer) {
		this.AccountID = buffer.ReadUInt32();
		this.Character = new Character();
		this.Character.Deserialize(buffer);
	}
}

public class DeleteChar : INet {
	public uint ID;

	public void Serialize(BinaryWriter buffer) {
		buffer.Write(this.ID);
	}

	public void Deserialize(BinaryReader buffer) {
		this.ID = buffer.ReadUInt32();
	}
}

public class Character : INet {
	public uint ID;
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

public class ListGamesResp : INet {
	public uint[] IDs;
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
		this.IDs = new uint[l0_1];
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

public class CreateGameResp : INet {
	public string Name;
	public uint ID;
	public ulong Seed;
	public Entity[] Entities;

	public void Serialize(BinaryWriter buffer) {
		buffer.Write((Int32)this.Name.Length);
		buffer.Write(System.Text.Encoding.UTF8.GetBytes(this.Name));
		buffer.Write(this.ID);
		buffer.Write(this.Seed);
		buffer.Write((Int32)this.Entities.Length);
		for (int v2 = 0; v2 < this.Entities.Length; v2++) {
			this.Entities[v2].Serialize(buffer);
		}
	}

	public void Deserialize(BinaryReader buffer) {
		int l0_1 = buffer.ReadInt32();
		byte[] temp0_1 = buffer.ReadBytes(l0_1);
		this.Name = System.Text.Encoding.UTF8.GetString(temp0_1);
		this.ID = buffer.ReadUInt32();
		this.Seed = buffer.ReadUInt64();
		int l3_1 = buffer.ReadInt32();
		this.Entities = new Entity[l3_1];
		for (int v2 = 0; v2 < l3_1; v2++) {
			this.Entities[v2] = new Entity();
			this.Entities[v2].Deserialize(buffer);
		}
	}
}

public class JoinGame : INet {
	public uint ID;
	public uint CharID;

	public void Serialize(BinaryWriter buffer) {
		buffer.Write(this.ID);
		buffer.Write(this.CharID);
	}

	public void Deserialize(BinaryReader buffer) {
		this.ID = buffer.ReadUInt32();
		this.CharID = buffer.ReadUInt32();
	}
}

public class GameConnected : INet {
	public ulong Seed;
	public Entity[] Entities;

	public void Serialize(BinaryWriter buffer) {
		buffer.Write(this.Seed);
		buffer.Write((Int32)this.Entities.Length);
		for (int v2 = 0; v2 < this.Entities.Length; v2++) {
			this.Entities[v2].Serialize(buffer);
		}
	}

	public void Deserialize(BinaryReader buffer) {
		this.Seed = buffer.ReadUInt64();
		int l1_1 = buffer.ReadInt32();
		this.Entities = new Entity[l1_1];
		for (int v2 = 0; v2 < l1_1; v2++) {
			this.Entities[v2] = new Entity();
			this.Entities[v2].Deserialize(buffer);
		}
	}
}

public class Entity : INet {
	public uint ID;
	public ushort EType;
	public ulong Seed;
	public int X;
	public int Y;
	public int Height;
	public int Width;
	public byte HealthPercent;

	public void Serialize(BinaryWriter buffer) {
		buffer.Write(this.ID);
		buffer.Write(this.EType);
		buffer.Write(this.Seed);
		buffer.Write(this.X);
		buffer.Write(this.Y);
		buffer.Write(this.Height);
		buffer.Write(this.Width);
		buffer.Write(this.HealthPercent);
	}

	public void Deserialize(BinaryReader buffer) {
		this.ID = buffer.ReadUInt32();
		this.EType = buffer.ReadUInt16();
		this.Seed = buffer.ReadUInt64();
		this.X = buffer.ReadInt32();
		this.Y = buffer.ReadInt32();
		this.Height = buffer.ReadInt32();
		this.Width = buffer.ReadInt32();
		this.HealthPercent = buffer.ReadByte();
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
	public int AbilityID;
	public uint Target;

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
	public int Damage;
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

