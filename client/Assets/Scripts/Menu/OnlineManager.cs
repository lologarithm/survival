using UnityEngine;
using UnityEngine.UI;
using System.Collections;

public class OnlineManager : MonoBehaviour {

	public InputField login;
	public InputField pwd;

	private NetworkMessenger msger;

	public void onLogin() {
		Debug.Log ("loggin in: " + login.text + " " + pwd.text);
		this.msger.Login(login.text, pwd.text);
	}

	public void onCreate() {
	}

	public void Awake() {
		GameObject manager = GameObject.Find("NetworkMessenger");
		this.msger = manager.GetComponent<NetworkMessenger>();
	}
}
