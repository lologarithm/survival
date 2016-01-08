using UnityEngine;
using System.Collections;

public class PlayParticles : MonoBehaviour {

    public ParticleSystem particles;

    public void Start()
    {
        this.particles.Stop();
    }

    public void StartParticle()
    {
        this.particles.Play();
    }

    public void StopParticle()
    {
        this.particles.Stop();
    }
}
