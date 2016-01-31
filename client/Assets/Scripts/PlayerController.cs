using UnityEngine;
using System.Collections;
using Assets;

public class PlayerController : MonoBehaviour {

    private NetworkMessenger msger;

    // Use this for initialization
    void Start () {
        GameObject managerobj = GameObject.Find("NetworkMessenger");
        this.msger = managerobj.GetComponent<NetworkMessenger>();
	}
	
	// Update is called once per frame
	void Update () {
        Vector2 moveVector = new Vector2(); 
        if (isPressed(ControlConst.MoveUp))
        {
            moveVector = moveVector + new Vector2(0, 1);
        }
        if (isPressed(ControlConst.MoveDown))
        {
            moveVector += new Vector2(0, -1);
        }
        if (isPressed(ControlConst.MoveLeft))
        {
            moveVector += new Vector2(-1, 0);
        }
        if (isPressed(ControlConst.MoveRight))
        {
            moveVector += new Vector2(1, 0);
        }
        if (moveVector.x != 0 || moveVector.y != 0) {
            moveVector = moveVector.normalized;
            Debug.Log(moveVector.ToString());
            msger.MovePlayer(moveVector);
        }
            
        foreach (ControlConst ability in ControlConst.Abilities)
        {
            if (isPressed(ability))
            {
                // do something 
            }
        }
	}

    // Returns if the given key axis is pressed
    bool isPressed(ControlConst cont)
    {
        return Input.GetAxis(cont.Value) > 0;
    }
}
