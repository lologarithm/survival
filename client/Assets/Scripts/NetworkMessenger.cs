using UnityEngine;
using System;
using System.IO;
using System.Net;
using System.Net.Sockets;
using System.Collections.Generic;

public class NetworkMessenger : MonoBehaviour
{
    Socket sending_socket = new Socket(AddressFamily.InterNetwork, SocketType.Dgram, ProtocolType.Udp);
    IPAddress send_to_address;
    IPEndPoint sending_end_point;

    // Caching network state
    private byte[] buff = new byte[1024];
    private List<byte> stored_bytes = new List<byte>();
    private NetMessage current_message = null;

    private Queue<NetMessage> message_queue = new Queue<NetMessage>();

    // Use this for initialization
    void Start()
    {
		Debug.Log("Starting network now!");
        this.stored_bytes = new List<byte>();
        this.send_to_address = IPAddress.Parse("127.0.0.1");
        this.sending_end_point = new IPEndPoint(send_to_address, 24816);
        sending_socket.Connect(this.sending_end_point);

        // 1. Fetch network!
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

	private void sendNetMessage(MsgType t, INet outmsg) {
		NetMessage msg = new NetMessage();
		MemoryStream stream = new MemoryStream();
		BinaryWriter buffer = new BinaryWriter(stream);
		outmsg.Serialize(buffer);
		msg.content = stream.ToArray();
		msg.content_length = (UInt16)msg.content.Length;
		msg.message_type = (byte)t;
		this.sending_socket.Send(msg.MessageBytes());	
	}

	public void CreateAccount(string name, string password) {
		CreateAcct outmsg = new CreateAcct();
		outmsg.Name = name;
		outmsg.Password = password;
		this.sendNetMessage (MsgType.CreateAcct, outmsg);
	}

	public void Login(string name, string password) {
		Login login_msg = new Login();
		login_msg.Name = name;
		login_msg.Password = password;
		this.sendNetMessage (MsgType.Login, login_msg);
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
            //0. Add buffer to all_bytes
            //1. if ( connection.all_bytes.Count > 0 ) - Read int off front (package_size) 
            //2. while ( connection.all_bytes.Count + bytesRead >= package_size )
            //3. add buffer to all_bytes and then queue a message, delete bytes from all_bytes

            byte[] subset_bytes = new byte[bytesRead];
            Array.Copy(this.buff, 0, subset_bytes, 0, bytesRead);
            this.stored_bytes.AddRange(subset_bytes);
            ProcessBytes();
            sending_socket.BeginReceive(this.buff, 0, buff.Length, SocketFlags.None, new AsyncCallback(ReceiveCallback), null);
        }
        else
            CloseConnection();
    }

    private void ProcessBytes()
    {
        byte[] input_bytes = this.stored_bytes.ToArray();
        if (this.current_message == null)
        {
            NetMessage nMsg = NetMessage.fromBytes(input_bytes);
            if (nMsg != null)
            {
                // Check for full content available. If so, time to add this to the processing queue.
				if (nMsg.full_content.Length == nMsg.content_length + NetMessage.DEFAULT_FRAME_LEN) {
					stored_bytes.RemoveRange (0, nMsg.full_content.Length);
					this.message_queue.Enqueue (nMsg);
					// If we have enough bytes to start a new message we call ProcessBytes again.
					if (input_bytes.Length - nMsg.full_content.Length > NetMessage.DEFAULT_FRAME_LEN) {
						ProcessBytes ();
					}
				}
            } else {
                this.current_message = nMsg;
                this.stored_bytes.RemoveRange(0, NetMessage.DEFAULT_FRAME_LEN);
                // Leave this as the this.current_message
            }
        }
        else {
			Debug.Log ("already have a message, lets try to add more bytes!");
            if (this.current_message.loadContent(input_bytes))
            {
                // We have enough bytes!
                stored_bytes.RemoveRange(0, this.current_message.content_length);
                this.message_queue.Enqueue(this.current_message);
                this.current_message = null;
            }
        }
        // We need to wait until later to finish loading!
    }

    // Update is called once per frame
    void Update()
    {
        int loops = this.message_queue.Count;
        for (int i = 0; i < loops; i++)
        {
            NetMessage msg = this.message_queue.Dequeue();
			INet parsedMsg = Messages.Parse (msg.message_type, msg.Content ());

			// Read from message queue and process!
			// Send updates to each object.

			Debug.Log("Got message: " + parsedMsg);
        }
    }

    void CloseConnection()
    {
		if (sending_socket.Connected) {
			sending_socket.Send (new byte[] { 255, 0, 0, 0, 0, 0, 0 });
			sending_socket.Close ();
		}
    }

    void OnApplicationQuit()
    {
        CloseConnection();
    }

	// Awake dont destroy keeps this object in memory even when we load a different scene.
	void Awake () {
		DontDestroyOnLoad(gameObject);
	}
}

public class NetMessage
{
    public const int DEFAULT_FRAME_LEN = 5;

    public byte message_type;
    public Int32 from_player;
    public UInt16 content_length;
    public UInt16 sequence;
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
            newMsg.message_type = bytes[0];
            newMsg.sequence = BitConverter.ToUInt16(bytes, 1);
            newMsg.content_length = BitConverter.ToUInt16(bytes, 3);

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
        if (bytes.Length >= this.content_length+DEFAULT_FRAME_LEN)
        {
			Array.Copy(bytes, 0, this.full_content, 0, DEFAULT_FRAME_LEN + this.content_length);
            return true;
        }

        return false;
    }
}
