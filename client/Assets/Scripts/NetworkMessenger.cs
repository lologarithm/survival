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

	private Queue<NetMessage> message_queue = new Queue<NetMessage>();
	private Dictionary<uint, List<Multipart>> multipart_cache = new Dictionary<uint, List<Multipart>>();

	// Use this for initialization
	void Start()
	{
		Debug.Log("Starting network now!");
		this.send_to_address = IPAddress.Parse("127.0.0.1");
		this.sending_end_point = new IPEndPoint(send_to_address, 24816);
		sending_socket.Connect(this.sending_end_point);

		// 1. Fetch network!
		ListGames outmsg = new ListGames();
		this.sendNetMessage(MsgType.ListGames, outmsg);

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

	private void sendNetMessage(MsgType t, INet outmsg)
	{
		NetMessage msg = new NetMessage();
		MemoryStream stream = new MemoryStream();
		BinaryWriter buffer = new BinaryWriter(stream);
		outmsg.Serialize(buffer);

		if (buffer.BaseStream.Length() + NetMessage.DEFAULT_FRAME_LEN > 512)
		{
			// TODO: Split the messages here!
		}
		else
		{
			msg.content = stream.ToArray();
			msg.content_length = (UInt16)msg.content.Length;
			msg.message_type = (byte)t;
			this.sending_socket.Send(msg.MessageBytes());
		}

	}

	public void CreateAccount(string name, string password)
	{
		CreateAcct outmsg = new CreateAcct();
		outmsg.Name = name;
		outmsg.Password = password;
		this.sendNetMessage(MsgType.CreateAcct, outmsg);
	}

	public void Login(string name, string password)
	{
		Login login_msg = new Login();
		login_msg.Name = name;
		login_msg.Password = password;
		this.sendNetMessage(MsgType.Login, login_msg);
	}

	public void CreateCharacter(string name)
	{
		CreateChar outmsg = new CreateChar();
		outmsg.Name = name;
		outmsg.AccountID = this.accounts[0];
		this.sendNetMessage(MsgType.CreateChar, outmsg);
	}

	public void CreateGame(string name)
	{
		CreateGame outmsg = new CreateGame();
		outmsg.Name = name;
		this.sendNetMessage(MsgType.CreateGame, outmsg);
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
		NetMessage nMsg = NetMessage.fromBytes(input_bytes);
		if (nMsg != null)
		{
			Debug.Log("Got a new netmsg: " + nMsg.message_type + " length: " + nMsg.content_length);
			// Check for full content available. If so, time to add this to the processing queue.
			if (nMsg.full_content != null)
			{
				this.numStored -= nMsg.full_content.Length;
				this.message_queue.Enqueue(nMsg);
				// If we have enough bytes to start a new message we call ProcessBytes again.
				if (input_bytes.Length - nMsg.full_content.Length > NetMessage.DEFAULT_FRAME_LEN)
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
			NetMessage msg = this.message_queue.Dequeue();
			INet parsedMsg = Messages.Parse(msg.message_type, msg.Content());

			// Read from message queue and process!
			// Send updates to each object.
			Debug.Log("Got message: " + parsedMsg);
			switch ((MsgType)msg.message_type)
			{
				case MsgType.Multipart:
					Multipart mpmsg = (Multipart)parsedMsg;
					if (!this.multipart_cache.ContainsKey(mpmsg.GroupID))
					{
						this.multipart_cache[mpmsg.GroupID] = new List<Multipart>(mpmsg.NumParts);
					}
					// 1. If this group doesn't exist, create it
					// 2. Insert message into group
					// 3. Check if all messages exist
					// 4. if so, group up bytes and call 'messages.parse' on the content
					// 5. clean up!
				case MsgType.CreateCharResp:
					characters.Add(((CreateCharResp)parsedMsg).Character);
					break;
				case MsgType.LoginResp:
					LoginResp lr = ((LoginResp)parsedMsg);
					if (lr.Characters.Length > 0)
					{
						characters.AddRange(lr.Characters);
					}
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
					gi.ID = cgr.ID;
					gi.entities = cgr.Entities;
					games.Add(gi);
					Debug.Log("Added game: " + gi.Name);
					break;
			}
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

public class NetMessage
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

	public static NetMessage fromBytes(byte[] bytes)
	{
		NetMessage newMsg = null;
		if (bytes.Length >= DEFAULT_FRAME_LEN)
		{
			newMsg = new NetMessage();
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
