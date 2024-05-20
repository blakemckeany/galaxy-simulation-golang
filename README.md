# Galaxy (n-body) simulation
## Written in Go using OpenGL

![alt text](<screenshot.png>)



### Optimisations
- Quadtree for faster physics calculations. Reduces the number of calculations per frame from **O(n^2)**, n being the total number of particles, to **O(n log n)**


## Framerate log (Apple M1 Chip):
### Pre optimisation
- 100 particles: 120fps
- 3000 particles: ~18fps
### Quad tree optimisation
- 100 particles: 120fps
- 3000 particles: ~90fps
- 10000 particles: ~38fps
### Shader optimisation
TBC



## Inspiration
[Coding Challenge #98.1: Quadtree - Part 1 / 2 by Coding Train](https://youtu.be/OJxEcs0w_kE?si=41RXxOhwx0NRP2C0)
[N-Body Simulation by a2flo](https://youtu.be/DoLe1c-eokI?si=aGGQCvkAPzL-Xjbu)
[Quadtrees and The Barnes-Hut Algorithm: The Making of a Gravity Simulation by William Y. Feng](https://www.youtube.com/watch?v=tOlKLJ4WmSE)
