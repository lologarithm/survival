using UnityEngine;
using System.Collections;

public class Exit : MonoBehaviour {

    public void ExitGame()
    {
        Debug.Log("Quitting now.");
        Application.Quit();
    }

}
