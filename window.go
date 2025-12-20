package main

import "github.com/go-gl/glfw/v3.3/glfw"

type Window struct {
	width, height   int
	windowObj       *glfw.Window
	updateCallbacks []func(float64)
	cursorCallbacks []func(float64, float64)
}

func (window *Window) Initialize(startWidth, startHeight int, titleName string) {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	windowObj, err := glfw.CreateWindow(startWidth, startHeight, titleName, nil, nil)

	window.windowObj = windowObj

	if err != nil {
		panic(err)
	}

	cursorCallback := func(w *glfw.Window, xpos, ypos float64) {
		for _, callback := range window.cursorCallbacks {
			callback(xpos, ypos)
		}
	}

	windowObj.MakeContextCurrent()
	windowObj.SetCursorPosCallback(cursorCallback)
	windowObj.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
}

func (window *Window) Terminate() {
	glfw.Terminate()
}

func (window *Window) Update() {
	window.width, window.height = glfw.GetCurrentContext().GetSize()
}

func (window *Window) ShouldClose() bool {
	return window.windowObj.ShouldClose()
}

func (window *Window) SwapBuffers() {
	window.windowObj.SwapBuffers()
}

func (window *Window) EnterUpdateLoop() {
	previousTime := 0.0
	for !window.ShouldClose() {
		window.Update()

		// Update
		time := glfw.GetTime()
		deltaTime := time - previousTime
		previousTime = time

		for _, callback := range window.updateCallbacks {
			callback(deltaTime)
		}

		// Do OpenGL stuff.
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
