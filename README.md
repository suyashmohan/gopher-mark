# gopher-mark

A simple fun experiment to benchmark how many gophers can we render when using game development libraries with go lang. Since most library use CGO to bind to the underlying C libraries, there will be some limitation. But limitation only matter if you hit them. So here we are just checking how far can we go.

# Findinds
When using Mac M1 Pro and running game in 1920x1080:
- Raylib-Go was able to keep 60 fps until 40,000 gophers after which fps drops

# Learnings
- Even if a CGO call only takes ~150 nanoseconds, doing it 45,000 times per frame adds ~6.75ms of pure CPU overhead just for switching languages. To fix this, we need to batch the draw calls. Instead of 45,000 calls to "Draw this texture", we make 1 call saying "Draw this mesh 45,000 times at these positions". (Abstract from Gemini's explanation)
