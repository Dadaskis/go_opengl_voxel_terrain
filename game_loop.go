package main

import (
	"fmt"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type GameLoop struct {
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
	gameWorld        GameWorld
	textureAtlas     uint32
}

func (loop *GameLoop) Initialize(window *Window) {
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

	loop.basicShader.LoadFile("basic")
	loop.triangleMesh = GetTriangleMesh()

	loop.gameWorld.Initialize()

	texture, err := NewTexture("atlas.png")
	if err != nil {
		panic(err)
	}
	loop.textureAtlas = texture
}

func (loop *GameLoop) CursorMove(xpos, ypos float64) {
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

func (loop *GameLoop) Clear() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	//gl.CullFace(gl.BACK)
	gl.ClearColor(
		loop.clearColor[0], loop.clearColor[1],
		loop.clearColor[2], loop.clearColor[3])
}

func (loop *GameLoop) UpdateCameraMatrices() {
	windowWidth, windowHeight := loop.window.width, loop.window.height

	loop.projection = mgl32.Perspective(
		mgl32.DegToRad(loop.camera.FOV), float32(windowWidth)/float32(windowHeight), 0.01, 1000.0,
	)

	loop.model = mgl32.Ident4()
}

func (loop *GameLoop) AssignCameraMatrices() {
	shader := loop.currentShader
	shader.UniformSetMat4("projection", &loop.projection)
	cameraMatrix := loop.camera.GetViewMatrix()
	shader.UniformSetMat4("camera", &cameraMatrix)
	shader.UniformSetMat4("model", &loop.model)
}

func (loop *GameLoop) AssignShader(shader *Shader) {
	loop.currentShader = shader
	shader.Use()
	loop.AssignCameraMatrices()
}

func (loop *GameLoop) UpdateRoutine(deltaTime float64) {
	loop.camera.ProcessKeyboard(loop.window, deltaTime)

	loop.Clear()
	loop.UpdateCameraMatrices()

	loop.AssignShader(&loop.basicShader)
	textureUniform := gl.GetUniformLocation(loop.currentShader.ID, GLString("tex"))
	gl.Uniform1i(textureUniform, 0)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, loop.textureAtlas)
	loop.triangleMesh.Render()
	loop.gameWorld.Render()
}
