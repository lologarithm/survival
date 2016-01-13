using UnityEngine;
using System.Collections;

public class LevelMaker : MonoBehaviour {
	private NetworkMessenger msger;

	// Use this for initialization
	void Start () {
		GameObject manager = GameObject.Find("NetworkMessenger");
		this.msger = manager.GetComponent<NetworkMessenger>();

		// TODO: Generate dynamic meshes here!
		// we can use this.msger.games[0].entities  for now to draw all entities!
		// 
	}


	// Update is called once per frame
	void Update () {
	
	}
}
