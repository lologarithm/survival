using UnityEngine;
using UnityEngine.UI;

public class GameCharManager : MonoBehaviour {
    private NetworkMessenger msger;

    public ScrollRect gamerect;
    public ScrollRect charrect;

    public InputField charname;
    public InputField gamename;

    public GameObject gamelistcontent;
    public GameObject toggleprefab;

    public void handleCreateGame()
    {
        this.msger.CreateGame(gamename.text);
        UnityEngine.SceneManagement.SceneManager.LoadScene("game_scene");
    }

	public void handleStartGame()
	{
		UnityEngine.SceneManagement.SceneManager.LoadScene("game_scene");
	}

    public void Update()
    {
        if (this.msger != null)
        {
            int ord = 0;
            foreach (GameInstance gi in this.msger.games)
            {
                GameObject exitingtog = GameObject.Find("gtoggle_" + gi.Name);
                if (exitingtog != null)
                {
                    ord++;
                    continue;
                }
                GameObject t = Instantiate(toggleprefab);
                t.name = "gtoggle_" + gi.Name;
                ((RectTransform)t.transform).position = new Vector3(0, 30 * ord, 0);
                GameObject l = GameObject.Find("gtoggle_" + gi.Name + "/Label");
                Text txt = l.GetComponent<Text>();
                txt.text = gi.Name;
                t.transform.SetParent(this.gamelistcontent.transform, false);
                ord++;
            }
        }
    }

    public void Awake()
    {
        GameObject manager = GameObject.Find("NetworkMessenger");
        this.msger = manager.GetComponent<NetworkMessenger>();
    }

}
