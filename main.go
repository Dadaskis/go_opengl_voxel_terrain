package main

import (
	"runtime"
)

var window Window
var gameLoop GameLoop

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func main() {
	window.Initialize(1280, 720, "Voxel Terrain")
	gameLoop.Initialize(&window)
	window.updateCallbacks = append(window.updateCallbacks, gameLoop.UpdateRoutine)
	window.EnterUpdateLoop()
}
