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
- **Windows**: MinGW-w64 (for CGo support)
- **Linux**: GLFW3 development libraries
- **macOS**: Xcode Command Line Tools

### Platform-Specific Setup

#### Windows
1. **Install MinGW-w64** (required for CGo to work with OpenGL):
   - **Option A - Using MSYS2 (Recommended)**:
     ```bash
     # Download and install MSYS2 from https://www.msys2.org/
     # Then open MSYS2 terminal and run:
     pacman -S mingw-w64-x86_64-gcc
     ```
   - **Option B - Direct download**:
     Download from [MinGW-w64](https://www.mingw-w64.org/downloads/) and install

2. **Add MinGW to PATH**:
   - Add `C:\msys64\mingw64\bin` (or your MinGW installation path) to your system's PATH environment variable
   - Verify installation:
     ```bash
     gcc --version
     ```

#### Linux (Ubuntu/Debian)
```bash
sudo apt update
sudo apt install libgl1-mesa-dev xorg-dev gcc gcc-go
```
*I didn't check it on Ubuntu/Debian, so I can't promise you that it'll work*

#### macOS
```bash
xcode-select --install
brew install glfw
```
*I don't have any devices with MacOS so... I'm not sure it'll work either*

### Installation

1. Clone the repository:
```bash
git clone https://github.com/Dadaskis/go_opengl_voxel_terrain
cd go_opengl_voxel_terrain
```

2. Initialize and download dependencies:
```bash
go mod init go_opengl_voxel_terrain
go mod tidy
```

3. Run the game:
```bash
go run .
```

**Note for Windows**: If you encounter build constraint errors with the OpenGL package, ensure MinGW is properly installed and in your PATH. The error "build constraints exclude all Go files" typically resolves after adding MinGW to PATH.

### Troubleshooting Common Issues

#### "build constraints exclude all Go files" error
- **Cause**: Missing or incorrectly configured MinGW installation
- **Solution**: 
  1. Verify MinGW is installed: `gcc --version`
  2. Check MinGW is in PATH: `echo $PATH`
  3. Ensure CGo is enabled: `go env CGO_ENABLED` should return "1"

#### GLFW-related errors
- **Windows**: Make sure `mingw-w64-x86_64-gcc` is installed via MSYS2
- **Linux**: Install GLFW dev libraries: `sudo apt install libglfw3-dev`
- **macOS**: Install GLFW via Homebrew: `brew install glfw`

#### "undefined: gl.Init" or similar OpenGL errors
- Ensure you're using the correct OpenGL package import:
  ```go
  import "github.com/go-gl/gl/v3.3-core/gl"
  ```
- Ensure that your GPU supports OpenGL 3.3

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
