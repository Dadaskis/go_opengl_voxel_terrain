package main

import (
	"fmt"
	"math"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	position         mgl32.Vec3
	front            mgl32.Vec3
	up               mgl32.Vec3
	right            mgl32.Vec3
	worldUp          mgl32.Vec3
	yaw              float32
	pitch            float32
	movementSpeed    float32
	mouseSensitivity float32
	FOV              float32
}

func (camera *Camera) InitializeDefaultValues() {
	camera.position = mgl32.Vec3{0.0, 0.0, 3.0}
	camera.front = mgl32.Vec3{0.0, 0.0, 1.0}
	camera.up = mgl32.Vec3{0.0, 1.0, 0.0}
	camera.right = mgl32.Vec3{-1.0, 0.0, 0.0}
	camera.worldUp = mgl32.Vec3{0.0, 1.0, 0.0}
	camera.yaw = 275.0
	camera.pitch = 0.0
	camera.movementSpeed = 0.005
	camera.mouseSensitivity = 0.1
	camera.FOV = 90.0
	camera.UpdateCameraVectors()
}

func (camera *Camera) UpdateCameraVectors() {
	front := mgl32.Vec3{}
	front[0] = float32(
		math.Cos(float64(mgl32.DegToRad(camera.yaw))) *
			math.Cos(float64(mgl32.DegToRad(camera.pitch))),
	)
	front[1] = float32(
		math.Sin(float64(mgl32.DegToRad(camera.pitch))),
	)
	front[2] = float32(
		math.Sin(float64(mgl32.DegToRad(camera.yaw))) *
			math.Cos(float64(mgl32.DegToRad(camera.pitch))),
	)

	camera.front = front.Normalize()
	camera.right = front.Cross(camera.worldUp).Normalize()
	camera.up = camera.right.Cross(front).Normalize()
}

func (camera *Camera) GetViewMatrix() mgl32.Mat4 {
	return mgl32.LookAtV(camera.position, camera.position.Add(camera.front), camera.up)
}

func (camera *Camera) ProcessMouseMovement(xpos, ypos float64) {
	xpos *= float64(camera.mouseSensitivity)
	ypos *= float64(camera.mouseSensitivity)
	camera.yaw += float32(xpos)
	camera.pitch -= float32(ypos)
	if camera.pitch > 89.0 {
		camera.pitch = 89
	}
	if camera.pitch < -89.0 {
		camera.pitch = -89
	}
	camera.UpdateCameraVectors()
}

func (camera *Camera) ProcessKeyboard(window *Window) {
	if window.windowObj.GetKey(glfw.KeyW) == glfw.Press {
		camera.position = camera.position.Add(camera.front.Mul(camera.movementSpeed))
	}

	if window.windowObj.GetKey(glfw.KeyS) == glfw.Press {
		camera.position = camera.position.Add(camera.front.Mul(-camera.movementSpeed))
	}

	if window.windowObj.GetKey(glfw.KeyA) == glfw.Press {
		camera.position = camera.position.Add(camera.right.Mul(-camera.movementSpeed))
	}

	if window.windowObj.GetKey(glfw.KeyD) == glfw.Press {
		camera.position = camera.position.Add(camera.right.Mul(camera.movementSpeed))
	}

	if window.windowObj.GetKey(glfw.KeySpace) == glfw.Press {
		camera.position = camera.position.Add(camera.up.Mul(camera.movementSpeed))
	}

	if window.windowObj.GetKey(glfw.KeyLeftControl) == glfw.Press {
		camera.position = camera.position.Add(camera.up.Mul(-camera.movementSpeed))
	}

	fmt.Println(camera.position)
}
