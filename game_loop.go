// Implements the main game loop and rendering pipeline.
// The GameLoop struct manages OpenGL initialization, shader management, camera controls,
// and the main update/render cycle for the voxel terrain generator.

package main

import (
	"fmt"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// GameLoop is the central coordinator for game systems including rendering,
// input processing, world management, and the main game update cycle.
type GameLoop struct {
	openGLVersion    string     // OpenGL version string retrieved from driver
	basicShader      Shader     // Primary shader program for rendering
	triangleMesh     Mesh       // Simple test mesh (triangle) for debugging/rendering
	clearColor       mgl32.Vec4 // Background clear color (RGBA)
	window           *Window    // Reference to the application window
	currentShader    *Shader    // Currently active shader program
	camera           *Camera    // Main camera for view control
	projection       mgl32.Mat4 // Projection matrix (perspective)
	model            mgl32.Mat4 // Model matrix (world transform)
	cursorPrevPosX   float64    // Previous mouse X position for delta calculation
	cursorPrevPosY   float64    // Previous mouse Y position for delta calculation
	cursorFirstFrame bool       // Flag for ignoring first mouse input frame
	gameWorld        GameWorld  // Main game world containing chunks and entities
	textureAtlas     uint32     // OpenGL texture ID for the block texture atlas
}

// Initialize sets up the game loop with OpenGL, shaders, camera, and world systems.
// window: The application window to render to and receive input from.
func (loop *GameLoop) Initialize(window *Window) {
	// Initialize OpenGL bindings
	if err := gl.Init(); err != nil {
		panic(err)
	}

	// Set default render state values
	loop.clearColor = mgl32.Vec4{0.0, 0.3, 1.0, 1.0} // Sky blue background
	loop.window = window
	loop.camera = &Camera{}
	loop.camera.InitializeDefaultValues()

	// Register mouse callback for camera control
	window.cursorCallbacks = append(window.cursorCallbacks, loop.CursorMove)

	// Log OpenGL version for debugging
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	// Load shader from files ("basic.glsl_vert", "basic.glsl_frag")
	loop.basicShader.LoadFile("basic")

	// Create a simple triangle mesh for testing/debugging
	loop.triangleMesh = GetTriangleMesh()

	// Initialize the game world (chunks, terrain, etc.)
	loop.gameWorld.Initialize()
	loop.gameWorld.currentCamera = loop.camera

	// Load texture atlas containing all block textures
	texture, err := NewTexture("atlas.png")
	if err != nil {
		panic(err)
	}
	loop.textureAtlas = texture
}

// CursorMove handles mouse movement input for camera rotation.
// Called by GLFW when the mouse moves. Calculates delta movement
// and passes it to the camera for look-around functionality.
func (loop *GameLoop) CursorMove(xpos, ypos float64) {
	// Skip first frame to avoid large jump when cursor enters window
	if !loop.cursorFirstFrame {
		loop.cursorFirstFrame = true
		loop.cursorPrevPosX = xpos
		loop.cursorPrevPosY = ypos
		return
	}

	// Calculate delta movement from previous position
	modXPos := xpos - loop.cursorPrevPosX
	modYPos := ypos - loop.cursorPrevPosY

	// Update camera orientation based on mouse movement
	loop.camera.ProcessMouseMovement(modXPos, modYPos)

	// Store current position for next frame's delta calculation
	loop.cursorPrevPosX = xpos
	loop.cursorPrevPosY = ypos
}

// Clear resets the framebuffer and sets up render state for a new frame.
// Sets background color, enables depth testing, and clears color/depth buffers.
func (loop *GameLoop) Clear() {
	// Clear both color and depth buffers
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// Enable depth testing for proper 3D rendering
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS) // Standard depth comparison (near objects obscure far ones)

	// Set clear color (background color)
	gl.ClearColor(
		loop.clearColor[0], loop.clearColor[1],
		loop.clearColor[2], loop.clearColor[3])
}

// UpdateCameraMatrices recalculates projection and model matrices
// based on current window size and camera state.
// Should be called whenever window is resized or camera FOV changes.
func (loop *GameLoop) UpdateCameraMatrices() {
	windowWidth, windowHeight := loop.window.width, loop.window.height

	// Create perspective projection matrix
	// Parameters: FOV, aspect ratio, near clip plane, far clip plane
	loop.projection = mgl32.Perspective(
		mgl32.DegToRad(loop.camera.FOV),
		float32(windowWidth)/float32(windowHeight),
		0.01, 1000.0,
	)

	// Identity model matrix (no world transform applied by default)
	loop.model = mgl32.Ident4()
}

// AssignCameraMatrices uploads the current camera, projection, and model matrices
// to the currently active shader program's uniform variables.
func (loop *GameLoop) AssignCameraMatrices() {
	shader := loop.currentShader

	// Upload matrices to shader uniforms
	shader.UniformSetMat4("projection", &loop.projection)
	cameraMatrix := loop.camera.GetViewMatrix()
	shader.UniformSetMat4("camera", &cameraMatrix)
	shader.UniformSetMat4("model", &loop.model)

	// Upload camera position for lighting calculations
	shader.UniformSetVec3("viewPos", &loop.camera.position)
}

// AssignShader activates a shader program and sets up its camera matrices.
// shader: The shader program to activate for subsequent rendering.
func (loop *GameLoop) AssignShader(shader *Shader) {
	loop.currentShader = shader
	shader.Use()                // Activate the shader program
	loop.AssignCameraMatrices() // Upload camera data
}

// UpdateRoutine is the main game loop function called each frame.
// Handles input processing, state updates, and rendering.
// deltaTime: Time elapsed since last frame (in seconds).
func (loop *GameLoop) UpdateRoutine(deltaTime float64) {
	// Process keyboard input for camera movement
	loop.camera.ProcessKeyboard(loop.window, deltaTime)

	// Clear screen and set up render state
	loop.Clear()

	// Update camera matrices (projection, view)
	loop.UpdateCameraMatrices()

	// Activate the basic shader program
	loop.AssignShader(&loop.basicShader)

	// Bind texture atlas to texture unit 0
	textureUniform := gl.GetUniformLocation(loop.currentShader.ID, GLString("tex"))
	gl.Uniform1i(textureUniform, 0)                  // Set uniform to use texture unit 0
	gl.ActiveTexture(gl.TEXTURE0)                    // Activate texture unit 0
	gl.BindTexture(gl.TEXTURE_2D, loop.textureAtlas) // Bind texture atlas

	// Render test triangle mesh (debug/placeholder)
	loop.triangleMesh.Render()

	// Render the game world (all chunks)
	loop.gameWorld.Render()
}
