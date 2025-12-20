package main

import (
	"fmt"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type RenderLoop struct {
	openGLVersion    string
	basicShader      Shader
	triangleMesh     Mesh
	clearColor       mgl32.Vec4
	window           *Window
	currentShader    *Shader
	camera           *Camera
	projection       mgl32.Mat4
	model            mgl32.Mat4
	cursorPrevPosX   float64
	cursorPrevPosY   float64
	cursorFirstFrame bool
}

func (loop *RenderLoop) Initialize(window *Window) {
	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	loop.clearColor = mgl32.Vec4{0.0, 0.3, 1.0, 1.0}
	loop.window = window
	loop.camera = &Camera{}
	loop.camera.InitializeDefaultValues()
	window.cursorCallbacks = append(window.cursorCallbacks, loop.CursorMove)

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	loop.basicShader = Shader{}
	loop.basicShader.LoadFile("basic")
	loop.triangleMesh = GetTriangleMesh()
}

func (loop *RenderLoop) CursorMove(xpos, ypos float64) {
	if !loop.cursorFirstFrame {
		loop.cursorFirstFrame = true
		loop.cursorPrevPosX = xpos
		loop.cursorPrevPosY = ypos
		return
	}
	modXPos := xpos - loop.cursorPrevPosX
	modYPos := ypos - loop.cursorPrevPosY
	loop.camera.ProcessMouseMovement(modXPos, modYPos)
	loop.cursorPrevPosX = xpos
	loop.cursorPrevPosY = ypos
}

func (loop *RenderLoop) Clear() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.CullFace(gl.BACK)
	gl.ClearColor(
		loop.clearColor[0], loop.clearColor[1],
		loop.clearColor[2], loop.clearColor[3])
}

func (loop *RenderLoop) UpdateCameraMatrices() {
	windowWidth, windowHeight := loop.window.width, loop.window.height

	loop.projection = mgl32.Perspective(
		mgl32.DegToRad(loop.camera.FOV), float32(windowWidth)/float32(windowHeight), 0.1, 10.0,
	)

	loop.model = mgl32.Ident4()
}

func (loop *RenderLoop) AssignCameraMatrices() {
	shader := loop.currentShader
	shader.UniformSetMat4("projection", &loop.projection)
	cameraMatrix := loop.camera.GetViewMatrix()
	shader.UniformSetMat4("camera", &cameraMatrix)
	shader.UniformSetMat4("model", &loop.model)
}

func (loop *RenderLoop) AssignShader(shader *Shader) {
	loop.currentShader = shader
	shader.Use()
	loop.AssignCameraMatrices()
}

func (loop *RenderLoop) UpdateRoutine(deltaTime float64) {
	loop.camera.ProcessKeyboard(loop.window)

	loop.Clear()
	loop.UpdateCameraMatrices()

	loop.AssignShader(&loop.basicShader)
	loop.triangleMesh.Render()
}
