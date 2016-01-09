using UnityEngine;
using System.Collections;

public class PlayParticles : MonoBehaviour {

    public ParticleSystem particles;

    public void Start()
	{
		if (particles != null) {
			this.particles.Stop();
		}
        
    }

    public void StartParticle()
    {
		if (particles != null) {
			this.particles.Play ();
		}
    }

    public void StopParticle()
    {
		if (particles != null) {
			this.particles.Stop ();
		}
    }
}
