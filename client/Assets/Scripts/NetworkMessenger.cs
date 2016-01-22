using UnityEngine;
using System;
using System.IO;
using System.Net;
using System.Net.Sockets;
using System.Collections.Generic;

// TODO: This has become way more than just network -- it is now also game state.
// We should separate them and make the parent 'game state manager' that persists across scenes.
public class NetworkMessenger : MonoBehaviour
{
	Socket sending_socket = new Socket(AddressFamily.InterNetwork, SocketType.Dgram, ProtocolType.Udp);
	IPAddress send_to_address;
	IPEndPoint sending_end_point;

	// Caching network state
	private byte[] buff = new byte[8192];
	private byte[] stored_bytes = new byte[8192];
	private int numStored = 0;


	private Queue<NetPacket> message_queue = new Queue<NetPacket>();
	private Dictionary<uint, Multipart[]> multipart_cache = new Dictionary<uint, Multipart[]>();
	private uint multi_groupid = 0;

	// Use this for initialization
	void Start()
	{
		Debug.Log("Starting network now!");
		this.send_to_address = IPAddress.Parse("127.0.0.1");
		this.sending_end_point = new IPEndPoint(send_to_address, 24816);
		sending_socket.Connect(this.sending_end_point);

		// 1. Fetch network!
		ListGames outmsg = new ListGames();
		this.sendNetPacket(MsgType.ListGames, outmsg);

		// Start Receive and a new Accept
		try
		{
			sending_socket.BeginReceive(this.buff, 0, this.buff.Length, SocketFlags.None, new AsyncCallback(ReceiveCallback), null);
		}
		catch (SocketException e)
		{
			// DO something
			System.Console.WriteLine(e.ToString());
		}

	}

	private void sendNetPacket(MsgType t, INet outmsg)
	{
		NetPacket msg = new NetPacket();
		MemoryStream stream = new MemoryStream();
		BinaryWriter buffer = new BinaryWriter(stream);
		outmsg.Serialize(buffer);

		if (buffer.BaseStream.Length + NetPacket.DEFAULT_FRAME_LEN > 512)
		{
			msg.message_type = (byte)MsgType.Multipart;
			//  calculate how many parts we have to split this into
			int maxsize = 512 - (12+NetPacket.DEFAULT_FRAME_LEN);
			int parts = ((int)buffer.BaseStream.Length / maxsize) + 1;
			this.multi_groupid++;
			int bstart = 0;
			for (int i = 0; i < parts; i++) {
				int bend = bstart + maxsize;
				if (i+1 == parts) {
					bend = bstart + (((int)buffer.BaseStream.Length) % maxsize);
				}
				Multipart wrapper = new Multipart();
				wrapper.ID = (ushort)i;
				wrapper.GroupID = this.multi_groupid;
				wrapper.NumParts = (ushort)parts;
				wrapper.Content = new byte[bend-bstart];
				buffer.BaseStream.Read(wrapper.Content, bstart, bend - bstart);

				MemoryStream pstream = new MemoryStream();
				BinaryWriter pbuffer = new BinaryWriter(pstream);
				wrapper.Serialize(pbuffer);

				msg.content = pstream.ToArray();
				msg.content_length = (ushort)pstream.Length;
				this.sending_socket.Send(msg.MessageBytes());
                bstart = bend;
			}
		}
		else
		{
			msg.content = stream.ToArray();
			msg.content_length = (ushort)msg.content.Length;
			msg.message_type = (byte)t;
			this.sending_socket.Send(msg.MessageBytes());
		}

	}

	public void CreateAccount(string name, string password)
	{
		CreateAcct outmsg = new CreateAcct();
		outmsg.Name = name;
		outmsg.Password = password;
        outmsg.CharName = name;
		this.sendNetPacket(MsgType.CreateAcct, outmsg);
	}

	public void Login(string name, string password)
	{
		Login login_msg = new Login();
		login_msg.Name = name;
		login_msg.Password = password;
		this.sendNetPacket(MsgType.Login, login_msg);
	}

    public void CreateGame(string name)
	{
		CreateGame outmsg = new CreateGame();
		outmsg.Name = name;
		this.sendNetPacket(MsgType.CreateGame, outmsg);
	}

	private void ReceiveCallback(IAsyncResult result)
	{
		int bytesRead = 0;
		try
		{
			bytesRead = sending_socket.EndReceive(result);
		}
		catch (SocketException exc)
		{
			CloseConnection();
			Debug.Log("Socket exception: " + exc.SocketErrorCode);
		}
		catch (Exception exc)
		{
			CloseConnection();
			Debug.Log("Exception: " + exc);
		}

		if (bytesRead > 0)
		{
			if (this.stored_bytes.Length < this.numStored + bytesRead)
			{
				byte[] newbuf = new byte[this.stored_bytes.Length * 2];
				Array.Copy(this.stored_bytes, 0, newbuf, 0, this.numStored);
				this.stored_bytes = newbuf;
			}
			Array.Copy(this.buff, 0, this.stored_bytes, this.numStored, bytesRead);
			this.numStored += bytesRead;
			ProcessBytes();
			sending_socket.BeginReceive(this.buff, 0, buff.Length, SocketFlags.None, new AsyncCallback(ReceiveCallback), null);
		}
		else
			CloseConnection();
	}

	private void ProcessBytes()
	{
		byte[] input_bytes = new byte[this.numStored];
		Array.Copy(this.stored_bytes, 0, input_bytes, 0, this.numStored);
		NetPacket nMsg = NetPacket.fromBytes(input_bytes);
		if (nMsg != null)
		{
			Debug.Log("Got a new netmsg: " + nMsg.message_type + " length: " + nMsg.content_length);
			// Check for full content available. If so, time to add this to the processing queue.
			if (nMsg.full_content != null)
			{
				this.numStored -= nMsg.full_content.Length;
				this.message_queue.Enqueue(nMsg);
				// If we have enough bytes to start a new message we call ProcessBytes again.
				if (input_bytes.Length - nMsg.full_content.Length > NetPacket.DEFAULT_FRAME_LEN)
				{
					ProcessBytes();
				}
			}
		}
	}

	public List<Character> characters = new List<Character>();
	public List<GameInstance> games = new List<GameInstance>();
	public List<UInt32> accounts = new List<UInt32>();
	// Update is called once per frame
	void Update()
	{
		int loops = this.message_queue.Count;
		for (int i = 0; i < loops; i++)
		{
			NetPacket msg = this.message_queue.Dequeue();
			this.ParseAndProcess(msg);
		}
	}

	void ParseAndProcess(NetPacket np) {
		INet parsedMsg = Messages.Parse(np.message_type, np.Content());

		// Read from message queue and process!
		// Send updates to each object.
		Debug.Log("Got message: " + parsedMsg);
		switch ((MsgType)np.message_type)
		{
			case MsgType.Multipart:
				Multipart mpmsg = (Multipart)parsedMsg;
				// 1. If this group doesn't exist, create it
				if (!this.multipart_cache.ContainsKey(mpmsg.GroupID))
				{
					this.multipart_cache[mpmsg.GroupID] = new Multipart[mpmsg.NumParts];
				}
				// 2. Insert message into group
				this.multipart_cache[mpmsg.GroupID][mpmsg.ID] = mpmsg;
				// 3. Check if all messages exist
				bool complete = true;
				int totalContent = 0;
				foreach (Multipart m in this.multipart_cache[mpmsg.GroupID])
				{
					if (m == null)
					{
						complete = false;
						break;
					}
					totalContent += m.Content.Length;
				}
				// 4. if so, group up bytes and call 'messages.parse' on the content
				if (complete)
				{
					byte[] content = new byte[totalContent];
					int co = 0;
					foreach (Multipart m in this.multipart_cache[mpmsg.GroupID])
					{
						System.Buffer.BlockCopy(m.Content, 0, content, co, m.Content.Length);
						co += m.Content.Length;
					}
					NetPacket newpacket = NetPacket.fromBytes(content);
					if (newpacket == null)
					{
						Debug.LogError("Multipart message content parsing failed... we done goofed");
					}
					this.ParseAndProcess(newpacket);
				}
				// 5. clean up!
				break;
			case MsgType.LoginResp:
				LoginResp lr = ((LoginResp)parsedMsg);
				characters.Add(lr.Character);
				accounts.Add(lr.AccountID);
				break;
			case MsgType.CreateAcctResp:
				CreateAcctResp car = ((CreateAcctResp)parsedMsg);
				accounts.Add(car.AccountID);
				break;
			case MsgType.ListGamesResp:
				ListGamesResp resp = ((ListGamesResp)parsedMsg);
				for (int j = 0; j < resp.IDs.Length; j++)
				{
					GameInstance ni = new GameInstance();
					ni.ID = resp.IDs[j];
					ni.Name = resp.Names[j];
				}
				break;
			case MsgType.GameConnected:
				GameConnected gc = ((GameConnected)parsedMsg);
				// TODO: handle connecting to a game!
				break;
			case MsgType.CreateGameResp:
				CreateGameResp cgr = ((CreateGameResp)parsedMsg);
				GameInstance gi = new GameInstance();
				gi.Name = cgr.Name;
				gi.ID = cgr.Game.ID;
				gi.entities = cgr.Game.Entities;
				games.Add(gi);
				Debug.Log("Added game: " + gi.Name);
				break;
		}
	}

	void CloseConnection()
	{
		if (sending_socket.Connected)
		{
			// sending_socket.Send (new byte[] { 255, 0, 0, 0, 0, 0, 0 }); // TODO: create a disconnect message.
			sending_socket.Close();
		}
	}

	void OnApplicationQuit()
	{
		CloseConnection();
	}

	// Awake dont destroy keeps this object in memory even when we load a different scene.
	void Awake()
	{
		DontDestroyOnLoad(gameObject);
	}
}

public class GameInstance
{
	public UInt32 ID;
	public string Name;
	public UInt64 Seed;
	public Entity[] entities;
}

public class NetPacket
{
	public const int DEFAULT_FRAME_LEN = 6;

	public ushort message_type;
	public int from_player;
	public ushort content_length;
	public ushort sequence;
	public byte[] content;
	public byte[] full_content;


	public byte[] MessageBytes()
	{
		///byte[] byte_array = new byte[]
		MemoryStream stream = new MemoryStream();
		using (BinaryWriter writer = new BinaryWriter(stream))
		{
			writer.Write(this.message_type);
			writer.Write(sequence);
			writer.Write(content_length);
			writer.Write(content);
		}
		return stream.ToArray();
	}

	public byte[] Content()
	{
		byte[] content = new byte[this.content_length];
		Array.Copy(this.full_content, DEFAULT_FRAME_LEN, content, 0, this.content_length);
		return content;
	}

	public static NetPacket fromBytes(byte[] bytes)
	{
		NetPacket newMsg = null;
		if (bytes.Length >= DEFAULT_FRAME_LEN)
		{
			newMsg = new NetPacket();
			newMsg.message_type = BitConverter.ToUInt16(bytes, 0);
			newMsg.sequence = BitConverter.ToUInt16(bytes, 2);
			newMsg.content_length = BitConverter.ToUInt16(bytes, 4);

			int totalLen = DEFAULT_FRAME_LEN + newMsg.content_length;
			if (bytes.Length >= totalLen)
			{
				newMsg.full_content = new byte[totalLen];
				Array.Copy(bytes, 0, newMsg.full_content, 0, totalLen);
			}
		}
		return newMsg;
	}

	public bool loadContent(byte[] bytes)
	{
		if (bytes.Length >= this.content_length + DEFAULT_FRAME_LEN)
		{
			Array.Copy(bytes, 0, this.full_content, 0, DEFAULT_FRAME_LEN + this.content_length);
			return true;
		}

		return false;
	}
}
