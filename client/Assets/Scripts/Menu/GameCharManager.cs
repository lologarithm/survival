using UnityEngine;
using UnityEngine.UI;

public class GameCharManager : MonoBehaviour {
    private NetworkMessenger msger;

    public ScrollRect gamerect;
    public ScrollRect charrect;

    public InputField charname;
    public InputField gamename;

    public GameObject charlistcontent;
    public GameObject toggleprefab;

    public void handleCreateChar()
    {
        this.msger.CreateCharacter(charname.text);
    }

    public void handleCreateGame()
    {
        this.msger.CreateGame(gamename.text);
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
            foreach (Character c in this.msger.characters)
            {
                GameObject exitingtog = GameObject.Find("chartoggle_" + c.Name);
                if (exitingtog != null)
                {
                    ord++;
                    continue;
                }
                GameObject t = Instantiate(toggleprefab);
                t.name = "chartoggle_" + c.Name;
                ((RectTransform)t.transform).position = new Vector3(0, 30 * ord, 0);
                GameObject l = GameObject.Find("chartoggle_" + c.Name + "/Label");
                Text txt = l.GetComponent<Text>();
                txt.text = c.Name;
                t.transform.SetParent(this.charlistcontent.transform, false);
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
