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
	basicShader   uint32
	triangleVAO   uint32
}

func (loop *RenderLoop) Initialize() {
	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

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

func (loop *RenderLoop) UpdateRoutine(deltaTime float64) {
	fmt.Println(deltaTime)

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// Update
	time := glfw.GetTime()

	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.0, 0.3, 1.0, 1.0)

	gl.UseProgram(loop.basicShader)

	windowWidth, windowHeight := window.width, window.height

	projection := mgl32.Perspective(
		mgl32.DegToRad(45.0), float32(windowWidth)/float32(windowHeight), 0.1, 10.0,
	)
	projectionUniform := gl.GetUniformLocation(loop.basicShader, GLString("projection"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	camera := mgl32.LookAtV(
		mgl32.Vec3{
			float32(math.Cos(time) * 2.0), 0, float32(math.Sin(time) * 2.0),
		},
		mgl32.Vec3{0, 0, 0},
		mgl32.Vec3{0, 1, 0},
	)
	cameraUniform := gl.GetUniformLocation(loop.basicShader, GLString("camera"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(loop.basicShader, GLString("model"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	gl.BindVertexArray(loop.triangleVAO)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
}
