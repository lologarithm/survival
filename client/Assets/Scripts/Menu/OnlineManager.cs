using UnityEngine;
using UnityEngine.UI;
using System.Collections;

public class OnlineManager : MonoBehaviour {

	public InputField login;
	public InputField pwd;

	private NetworkMessenger msger;

	public void onLogin() {
		Debug.Log ("logging in: " + login.text + " " + pwd.text);
		this.msger.Login(login.text, pwd.text);
        UnityEngine.SceneManagement.SceneManager.LoadScene("selection_menu_scene");
    }

	public void onCreate() {
        Debug.Log("Creating account: " + login.text + " " + pwd.text);
        this.msger.CreateAccount(login.text, pwd.text);
        UnityEngine.SceneManagement.SceneManager.LoadScene("selection_menu_scene");
    }

	public void Awake() {
		GameObject manager = GameObject.Find("NetworkMessenger");
        this.msger = manager.GetComponent<NetworkMessenger>();
	}
}
