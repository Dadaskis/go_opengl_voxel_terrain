package main

import (
	"fmt"
	"math"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type RenderLoop struct {
	openGLVersion   string
	basicShader     uint32
	triangleVAO     uint32
	clearColor      mgl32.Vec4
	window          *Window
	currentShaderID uint32
	projection      mgl32.Mat4
	camera          mgl32.Mat4
	model           mgl32.Mat4
	cameraPos       mgl32.Vec3
	cameraTarget    mgl32.Vec3
	FOV             float32
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

	shader, err := NewProgramFile("basic")
	if err != nil {
		fmt.Println("Failed to do a shader")
		panic(err)
	}
	loop.basicShader = shader
	loop.triangleVAO = GetTriangleMesh(shader)
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
	projectionUniform := gl.GetUniformLocation(loop.currentShaderID, GLString("projection"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &loop.projection[0])

	cameraUniform := gl.GetUniformLocation(loop.currentShaderID, GLString("camera"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &loop.camera[0])

	modelUniform := gl.GetUniformLocation(loop.currentShaderID, GLString("model"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &loop.model[0])
}

func (loop *RenderLoop) AssignShaderID(shaderID uint32) {
	loop.currentShaderID = shaderID
	gl.UseProgram(shaderID)
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

	loop.AssignShaderID(loop.basicShader)
	gl.BindVertexArray(loop.triangleVAO)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
}
