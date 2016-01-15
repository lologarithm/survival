using UnityEngine;
using System.Collections;

public class LevelMaker : MonoBehaviour {
	private NetworkMessenger msger;

	public Transform treePrefab;

	// Use this for initialization
	void Start () {
		GameObject manager = GameObject.Find("NetworkMessenger");
		this.msger = manager.GetComponent<NetworkMessenger>();

		// we can use this.msger.games[0].entities  for now to draw all entities!
		for ( int i = 0; i < this.msger.games[0].entities.Length; i++) {
			Entity e = this.msger.games [0].entities [i];
			Transform tree = (Transform)Instantiate(this.treePrefab, new Vector3(e.X, 0, e.Y), Quaternion.identity);
			tree.localScale = new Vector3 (e.Height, e.Height, e.Width);
		}
	}


	// Update is called once per frame
	void Update () {
	
	}
}
