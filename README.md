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
https://youtu.be/OJxEcs0w_kE?si=41RXxOhwx0NRP2C0