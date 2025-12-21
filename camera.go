// A first-person camera system.
// The Camera struct provides FPS-style movement with mouse look and keyboard controls.
package main

import (
	"math"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

// Camera represents a first-person camera in 3D space.
// It maintains position, orientation, and movement parameters for navigation.
type Camera struct {
	position         mgl32.Vec3 // Current camera position in world space
	front            mgl32.Vec3 // Forward direction vector (normalized)
	up               mgl32.Vec3 // Up direction vector (normalized)
	right            mgl32.Vec3 // Right direction vector (normalized)
	worldUp          mgl32.Vec3 // World's up vector (typically {0,1,0})
	yaw              float32    // Horizontal rotation in degrees (0-360)
	pitch            float32    // Vertical rotation in degrees (-89 to 89)
	movementSpeed    float32    // Base movement speed in units per second
	mouseSensitivity float32    // Mouse look sensitivity multiplier
	FOV              float32    // Field of view in degrees (perspective)
}

// InitializeDefaultValues sets up the camera with standard starting values.
// This includes position above ground (Y=60), looking north, with default speeds.
func (camera *Camera) InitializeDefaultValues() {
	camera.position = mgl32.Vec3{0.0, 60.0, 0.0}
	camera.front = mgl32.Vec3{0.0, 0.0, 1.0}
	camera.up = mgl32.Vec3{0.0, 1.0, 0.0}
	camera.right = mgl32.Vec3{-1.0, 0.0, 0.0}
	camera.worldUp = mgl32.Vec3{0.0, 1.0, 0.0}
	camera.yaw = 275.0 // Points slightly west of north for natural starting view
	camera.pitch = 0.0 // Level horizon
	camera.movementSpeed = 10.0
	camera.mouseSensitivity = 0.1
	camera.FOV = 90.0 // Wide field of view for FPS games
	camera.UpdateCameraVectors()
}

// UpdateCameraVectors recalculates the camera's orientation vectors
// based on current yaw and pitch values. Must be called after changing
// rotation angles to maintain consistent coordinate system.
func (camera *Camera) UpdateCameraVectors() {
	// Calculate new front vector using spherical coordinates
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
	// Recalculate right vector using cross product with world up
	camera.right = camera.front.Cross(camera.worldUp).Normalize()
	// Recalculate up vector to ensure orthonormal basis
	camera.up = camera.right.Cross(camera.front).Normalize()
}

// GetViewMatrix constructs and returns a view matrix for rendering.
// This matrix transforms world coordinates to camera/view space.
func (camera *Camera) GetViewMatrix() mgl32.Mat4 {
	// LookAt creates a view matrix looking from position to position+front
	return mgl32.LookAtV(camera.position, camera.position.Add(camera.front), camera.up)
}

// ProcessMouseMovement updates camera orientation based on mouse input.
// xpos, ypos: Mouse movement deltas (typically in pixels)
// Sensitivity is applied and pitch is clamped to prevent gimbal lock.
func (camera *Camera) ProcessMouseMovement(xpos, ypos float64) {
	// Apply sensitivity scaling to mouse movement
	xpos *= float64(camera.mouseSensitivity)
	ypos *= float64(camera.mouseSensitivity)

	// Update rotation angles (yaw accumulates, pitch subtracts for natural mouse movement)
	camera.yaw += float32(xpos)
	camera.pitch -= float32(ypos)

	// Clamp pitch to prevent camera flipping
	if camera.pitch > 89.0 {
		camera.pitch = 89
	}
	if camera.pitch < -89.0 {
		camera.pitch = -89
	}

	// Update directional vectors with new rotation
	camera.UpdateCameraVectors()
}

// ProcessKeyboard handles camera movement based on keyboard input.
// Uses WASD for horizontal movement, Space/Control for vertical movement.
// deltaTime: Time since last frame (in seconds) for frame-rate independent movement.
func (camera *Camera) ProcessKeyboard(window *Window, deltaTime float64) {
	// Calculate movement distance for this frame
	velocity := camera.movementSpeed * float32(deltaTime)

	// Forward movement (W key)
	if window.windowObj.GetKey(glfw.KeyW) == glfw.Press {
		camera.position = camera.position.Add(camera.front.Mul(velocity))
	}

	// Backward movement (S key)
	if window.windowObj.GetKey(glfw.KeyS) == glfw.Press {
		camera.position = camera.position.Add(camera.front.Mul(-velocity))
	}

	// Left strafe (A key)
	if window.windowObj.GetKey(glfw.KeyA) == glfw.Press {
		camera.position = camera.position.Add(camera.right.Mul(-velocity))
	}

	// Right strafe (D key)
	if window.windowObj.GetKey(glfw.KeyD) == glfw.Press {
		camera.position = camera.position.Add(camera.right.Mul(velocity))
	}

	// Ascend (Space key)
	if window.windowObj.GetKey(glfw.KeySpace) == glfw.Press {
		camera.position = camera.position.Add(camera.up.Mul(velocity))
	}

	// Descend (Left Control key)
	if window.windowObj.GetKey(glfw.KeyLeftControl) == glfw.Press {
		camera.position = camera.position.Add(camera.up.Mul(-velocity))
	}
}
