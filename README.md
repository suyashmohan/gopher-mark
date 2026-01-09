# gopher-mark

A simple fun experiment to benchmark how many gophers can we render when using
game development libraries with go lang. Since most library use CGO to bind to
the underlying C libraries, there will be some limitation. But limitation only
matter if you hit them. So here we are just checking how far can we go.

# Findinds

When using _Mac M1 Pro_ and running game in 1920x1080:

- _Raylib-Go_ was able to keep 60 fps until _40,000_ gophers after which fps
  drops
- _Ebiten_ was able to keep 60 fps until _55,000_ gophers after which fps drops
- ECS version made with Ark, had similar performance
- _SDL3-Go_ was able to keep 60 fps until _110,000_ gophers after which fps drops
- _Raylib (C)_ was able to keep 60 fps until _120,000_ gophers after which fps drops

# Learnings (Abstract from Gemini's explanation)

- Even if a CGO call only takes ~150 nanoseconds, doing it 45,000 times per
  frame adds ~6.75ms of pure CPU overhead just for switching languages. To fix
  this, we need to batch the draw calls. Instead of 45,000 calls to "Draw this
  texture", we make 1 call saying "Draw this mesh 45,000 times at these
  positions".
- Ebiten batches the GPU draw calls but we are still performing 55,000 Go
  function calls (screen.DrawImage) leading to CPU overhead. Suggested solution
  to use DrawTriangles. It changes the equation from 55,000 function calls to 1
  function call
