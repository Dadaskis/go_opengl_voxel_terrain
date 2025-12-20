package main

import (
	"fmt"
	"math"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type RenderLoop struct {
	openGLVersion string
	basicShader   *Shader
	triangleVAO   uint32
	clearColor    mgl32.Vec4
	window        *Window
	currentShader *Shader
	projection    mgl32.Mat4
	camera        mgl32.Mat4
	model         mgl32.Mat4
	cameraPos     mgl32.Vec3
	cameraTarget  mgl32.Vec3
	FOV           float32
}

func (loop *RenderLoop) Initialize(window *Window) {
	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	loop.clearColor = mgl32.Vec4{0.0, 0.3, 1.0, 1.0}
	loop.window = window
	loop.FOV = 45.0

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	loop.basicShader = &Shader{}
	loop.basicShader.LoadFile("basic")
	loop.triangleVAO = GetTriangleMesh(loop.basicShader.ID)
}

func (loop *RenderLoop) Clear() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(
		loop.clearColor[0], loop.clearColor[1],
		loop.clearColor[2], loop.clearColor[3])
}

func (loop *RenderLoop) UpdateCameraMatrices() {
	windowWidth, windowHeight := loop.window.width, loop.window.height

	loop.projection = mgl32.Perspective(
		mgl32.DegToRad(loop.FOV), float32(windowWidth)/float32(windowHeight), 0.1, 10.0,
	)

	loop.camera = mgl32.LookAtV(
		loop.cameraPos, loop.cameraTarget, mgl32.Vec3{0, 1, 0},
	)

	loop.model = mgl32.Ident4()
}

func (loop *RenderLoop) AssignCameraMatrices() {
	shader := loop.currentShader
	shader.UniformSetMat4("projection", &loop.projection)
	shader.UniformSetMat4("camera", &loop.camera)
	shader.UniformSetMat4("model", &loop.model)
}

func (loop *RenderLoop) AssignShader(shader *Shader) {
	loop.currentShader = shader
	shader.Use()
	loop.AssignCameraMatrices()
}

func (loop *RenderLoop) UpdateRoutine(deltaTime float64) {
	fmt.Println(deltaTime)

	time := glfw.GetTime()
	loop.cameraPos = mgl32.Vec3{
		float32(math.Cos(time) * 2.0), 0, float32(math.Sin(time) * 2.0),
	}
	loop.cameraTarget = mgl32.Vec3{0, 0, 0}

	loop.Clear()
	loop.UpdateCameraMatrices()

	loop.AssignShader(loop.basicShader)
	gl.BindVertexArray(loop.triangleVAO)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
}
