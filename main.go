// Main is the entry point for the voxel terrain application.
// It initializes the main game systems, sets up the window, and starts the game loop.

package main

import (
	"runtime"
)

// Global application instances
var window Window     // Main application window (handles GLFW, input, rendering context)
var gameLoop GameLoop // Central game loop (handles updates, rendering, game state)

// init is called before main() and performs critical initialization.
// It ensures OpenGL/GLFW functions run on the main thread as required by most windowing systems.
func init() {
	// LockOSThread ensures that subsequent code runs on the main OS thread.
	// This is REQUIRED for OpenGL/GLFW operations which must occur on the same thread
	// that created the OpenGL context. Without this, GL calls would crash or fail.
	runtime.LockOSThread()
}

// main is the application entry point.
// Initializes all game systems, sets up the window, and enters the main game loop.
func main() {
	// Initialize the application window with specified dimensions and title
	// 1280x720 is a common HD resolution for games
	window.Initialize(1280, 720, "Voxel Terrain")

	// Initialize the game loop with reference to the window
	// This sets up OpenGL, shaders, camera, textures, and world generation
	gameLoop.Initialize(&window)

	// Register the game loop's update routine as a callback
	// This function will be called every frame to update and render the game
	window.updateCallbacks = append(window.updateCallbacks, gameLoop.UpdateRoutine)

	// Enter the main update/render loop
	// This function blocks until the window is closed
	window.EnterUpdateLoop()
}
