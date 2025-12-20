package main

import (
	"runtime"
)

var window Window
var renderLoop RenderLoop

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func main() {
	window.Initialize(640, 480, "Voxel Terrain")
	renderLoop.Initialize(&window)
	window.updateCallbacks = append(window.updateCallbacks, renderLoop.UpdateRoutine)
	window.EnterUpdateLoop()
}
