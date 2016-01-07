using UnityEngine;
using System.Collections;
using Assets;

public class PlayerController : MonoBehaviour {

	// Use this for initialization
	void Start () {
	
	}
	
	// Update is called once per frame
	void Update () {
        var moveVector = new Vector2();
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

        moveVector = moveVector.normalized;
        Debug.Log(moveVector.ToString());
	}

    // Returns if the given key axis is pressed
    bool isPressed(ControlConst cont)
    {
        return Input.GetAxis(cont.Value) > 0;
    }
}
