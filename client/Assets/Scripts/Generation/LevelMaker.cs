using UnityEngine;
using System.Collections;

public class LevelMaker : MonoBehaviour {
    private NetworkMessenger msger;
	public Transform treePrefab;

    public GameObject player;
    public GameObject playercam;

	// Use this for initialization
	void Start () {
		GameObject manager = GameObject.Find("NetworkMessenger");
		this.msger = manager.GetComponent<NetworkMessenger>();

		// we can use this.msger.games[0].entities  for now to draw all entities!
		for ( int i = 0; i < this.msger.games[0].entities.Length; i++) {
			Entity e = this.msger.games [0].entities [i];
            if (e.EType == 3)
            {
                Transform tree = (Transform)Instantiate(this.treePrefab, new Vector3(e.X, 0, e.Y), Quaternion.identity);
                tree.localScale = new Vector3(e.Height, e.Height, e.Width);
                tree.name = "tree" + e.ID.ToString();
            }
            else if (e.EType == 4)
            {
                // TODO: move player here!
            }
		}
	}


	// Update is called once per frame
	void Update () {
        for ( int i = 0; i < this.msger.games[0].entities.Length; i++) {
            Entity e = this.msger.games[0].entities[i];
            if (e.EType == 4 && (player.transform.position.x != e.X || player.transform.position.z != e.Y))
            {
                Debug.Log("moving player!!");
                player.transform.position = new Vector3(e.X, 0, e.Y);
                playercam.transform.position = new Vector3(e.X-150, playercam.transform.position.y, e.Y);
            }
        }	
	}
}
