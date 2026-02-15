<p align="center">
	<img src="githubImage/example0.png" width="640" alt="Nice image over here...">
</p>

# Voxel Terrain Generator in Go

A Minecraft-style voxel terrain generator built with Go, OpenGL, and GLFW. This project demonstrates procedural world generation, chunk-based rendering, and first-person camera controls in a 3D block world.

## Features

- **Procedural Terrain Generation**: Uses OpenSimplex noise to create realistic terrain with mountains and caves
- **Chunk-Based World**: 16×16×256 block chunks with efficient face culling for optimal rendering
- **First-Person Camera**: Full mouse look and keyboard controls (WASD + Space/Ctrl for vertical movement)
- **Dynamic Loading**: Chunks load and unload based on camera position with background generation
- **Texture Atlas Support**: Multiple block types with different textures per face (grass, dirt, stone)
- **OpenGL 3.3 Core**: Modern rendering pipeline with shaders, VAOs, and VBOs

## Architecture

### Core Components

- **`main.go`**: Entry point, thread locking for OpenGL, main loop initialization
- **`window.go`**: GLFW window management, input callbacks, delta time calculation
- **`game_loop.go`**: Main rendering pipeline, camera updates, shader management
- **`camera.go`**: First-person camera with FPS-style movement and orientation
- **`chunk.go`**: 16×16×256 block container with mesh generation and face culling
- **`game_world.go`**: Chunk management, dynamic loading/unloading, render distance control
- **`mesh.go`**: Vertex data structures, VAO/VBO management, rendering utilities
- **`shader.go`**: GLSL shader compilation, linking, and uniform management
- **`block_data.go`**: Block type definitions and UV texture coordinates
- **`gl_utilities.go`**: OpenGL helpers, texture loading, mesh utilities

## Getting Started

### Prerequisites

- Go 1.21 or higher
- OpenGL 3.3 compatible graphics card
- GLFW3 development libraries

### Installation

1. Clone the repository:
```bash
git clone https://github.com/Dadaskis/go_opengl_voxel_terrain
cd go_opengl_voxel_terrain
```

2. Install dependencies:
```bash
go mod init go_opengl_voxel_terrain
go get github.com/go-gl/gl/v3.3-core/gl
go get github.com/go-gl/glfw/v3.3/glfw
go get github.com/go-gl/mathgl/mgl32
go get github.com/ojrac/opensimplex-go
```

3. Place a texture atlas at `atlas.png` (1024×1024 texture, where a single tile is 64×64)

4. Run the game:
```bash
go run .
```

## Controls

- **W/A/S/D**: Move forward/left/backward/right
- **Space**: Ascend
- **Left Control**: Descend
- **Mouse**: Look around

## Technical Details

### Chunk Generation
- Each chunk generates asynchronously in a goroutine
- 2D noise determines terrain height (50 ±30 blocks)
- 3D noise creates cave systems
- Automatic layering: stone base → dirt (top 5 blocks) → grass (top block)

### Rendering Optimization
- Face culling: only render faces adjacent to air blocks
- Chunk-based render distance (configurable, default 16 chunks in each direction)
- Dirty flag system for mesh updates
- Interleaved vertex attributes for better cache performance

### Block Data System
- Each block type has per-face texture coordinates
- Texture atlas UV mapping (64×64 tiles in 1024×1024 atlas)
- Air blocks skip rendering entirely

### World Management
- Background goroutine monitors camera position every 300ms
- Only loads chunks within render distance
- Maintains map of all loaded chunks for quick lookup

## Project Structure

```
.
├── main.go              # Application entry point
├── window.go            # GLFW window management
├── game_loop.go         # Main rendering loop
├── camera.go            # First-person camera
├── chunk.go             # Block container and mesh generation
├── game_world.go        # World/chunk management
├── mesh.go              # Vertex data and OpenGL buffers
├── shader.go            # Shader compilation
├── block_data.go        # Block type definitions
├── gl_utilities.go      # OpenGL helpers
├── basic.glsl_vert      # Vertex shader
├── basic.glsl_frag      # Fragment shader
└── atlas.png            # Texture atlas
```

## Dependencies

- [go-gl/gl](https://github.com/go-gl/gl): OpenGL bindings
- [go-gl/glfw](https://github.com/go-gl/glfw): GLFW windowing
- [go-gl/mathgl](https://github.com/go-gl/mathgl): Mathematics for 3D graphics
- [ojrac/opensimplex-go](https://github.com/ojrac/opensimplex-go): OpenSimplex noise generation

## Acknowledgments

- Inspired by Minecraft's infinite world generation
- Uses OpenSimplex noise (public domain alternative to Perlin noise)
- Built with modern OpenGL practices and Go's concurrency model
