// Implements window management and GLFW integration.
// The Window struct handles window creation, input processing, and the main game loop.

package main

import "github.com/go-gl/glfw/v3.3/glfw"

// Window manages a GLFW window, input callbacks, and the main update loop.
type Window struct {
	width, height   int                      // Current window dimensions in pixels
	windowObj       *glfw.Window             // GLFW window object
	updateCallbacks []func(float64)          // Functions called each frame with deltaTime
	cursorCallbacks []func(float64, float64) // Functions called on mouse movement
}

// Initialize creates and configures a new GLFW window.
// startWidth, startHeight: Initial window dimensions
// titleName: Window title displayed in title bar
func (window *Window) Initialize(startWidth, startHeight int, titleName string) {
	// Initialize GLFW library
	err := glfw.Init()
	if err != nil {
		panic(err)
	}

	// Set window hints (configuration options)
	glfw.WindowHint(glfw.Resizable, glfw.True)   // Allow window resizing
	glfw.WindowHint(glfw.ContextVersionMajor, 3) // OpenGL 3.3
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile) // Core profile (no legacy)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)    // Forward compatibility

	// Create the window
	windowObj, err := glfw.CreateWindow(startWidth, startHeight, titleName, nil, nil)
	if err != nil {
		panic(err)
	}

	window.windowObj = windowObj

	// Create mouse movement callback function
	cursorCallback := func(w *glfw.Window, xpos, ypos float64) {
		// Forward mouse position to all registered callbacks
		for _, callback := range window.cursorCallbacks {
			callback(xpos, ypos)
		}
	}

	// Make this window's OpenGL context current (required for GL operations)
	windowObj.MakeContextCurrent()

	// Register mouse movement callback
	windowObj.SetCursorPosCallback(cursorCallback)

	// Lock cursor to window center (for FPS-style camera control)
	windowObj.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
}

// Terminate cleanly shuts down GLFW and releases resources.
// Should be called when the application exits.
func (window *Window) Terminate() {
	glfw.Terminate()
}

// Update refreshes window state, typically called each frame.
// Updates width and height to current window dimensions.
func (window *Window) Update() {
	window.width, window.height = glfw.GetCurrentContext().GetSize()
}

// ShouldClose checks if the window has been requested to close
// (e.g., user clicked the close button).
func (window *Window) ShouldClose() bool {
	return window.windowObj.ShouldClose()
}

// SwapBuffers swaps the front and back buffers (double buffering).
// This presents the rendered frame to the screen.
func (window *Window) SwapBuffers() {
	window.windowObj.SwapBuffers()
}

// EnterUpdateLoop starts the main game loop.
// This function blocks until the window is closed.
// It calls all registered update callbacks each frame with deltaTime.
func (window *Window) EnterUpdateLoop() {
	previousTime := 0.0
	for !window.ShouldClose() {
		// Update window dimensions
		window.Update()

		// Calculate deltaTime (time since last frame)
		time := glfw.GetTime()
		deltaTime := time - previousTime
		previousTime = time

		// Call all registered update callbacks
		for _, callback := range window.updateCallbacks {
			callback(deltaTime)
		}

		// Present the rendered frame
		window.SwapBuffers()

		// Process pending events (input, window events, etc.)
		glfw.PollEvents()
	}
}
