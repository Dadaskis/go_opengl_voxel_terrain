// Implements mesh data structures and OpenGL rendering operations.
// The Mesh struct handles vertex data storage, VAO/VBO management, and rendering.

package main

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// MeshVertex represents a single vertex in a mesh with position, color, normal, and UV coordinates.
// This is the fundamental data structure for 3D rendering in the engine.
type MeshVertex struct {
	position mgl32.Vec3 // 3D position in model space
	color    mgl32.Vec3 // RGB color values (0.0-1.0)
	normal   mgl32.Vec3 // Surface normal for lighting calculations
	UV       mgl32.Vec2 // Texture coordinates (U, V) in 0.0-1.0 range
}

// Mesh represents a collection of vertices that form a 3D object.
// It manages OpenGL vertex array objects (VAO) and vertex buffer objects (VBO).
type Mesh struct {
	vertices  []MeshVertex // Raw vertex data (CPU-side)
	arrayData []float32    // Flattened vertex data for GPU upload
	VAO       uint32       // OpenGL Vertex Array Object ID
}

// AddVertex appends a new vertex to the mesh.
// position: 3D location of the vertex
// color: RGB color of the vertex
// normal: Surface normal vector (should be normalized)
// UV: Texture coordinates for mapping textures
func (mesh *Mesh) AddVertex(position, color, normal mgl32.Vec3, UV mgl32.Vec2) {
	vertex := MeshVertex{
		position, color, normal, UV,
	}
	mesh.vertices = append(mesh.vertices, vertex)
}

// appendVec3ToArray converts a 3-component vector to a flat float32 array.
// array: Target float32 array to append to
// vec: 3-component vector to flatten
// Returns: Updated array with vector components appended
func appendVec3ToArray(array []float32, vec *mgl32.Vec3) []float32 {
	array = append(array, vec[0])
	array = append(array, vec[1])
	array = append(array, vec[2])
	return array
}

// appendVec2ToArray converts a 2-component vector to a flat float32 array.
// array: Target float32 array to append to
// vec: 2-component vector to flatten
// Returns: Updated array with vector components appended
func appendVec2ToArray(array []float32, vec *mgl32.Vec2) []float32 {
	array = append(array, vec[0])
	array = append(array, vec[1])
	return array
}

// PrepareArrayData converts the mesh's vertex data into a flat float32 array
// suitable for uploading to the GPU via OpenGL buffer objects.
// Each vertex consists of 11 float32 values: position(3), color(3), normal(3), UV(2)
func (mesh *Mesh) PrepareArrayData() {
	vertices := []float32{}

	for _, vertex := range mesh.vertices {
		// Append all vertex attributes in interleaved format
		vertices = appendVec3ToArray(vertices, &vertex.position)
		vertices = appendVec3ToArray(vertices, &vertex.color)
		vertices = appendVec3ToArray(vertices, &vertex.normal)
		vertices = appendVec2ToArray(vertices, &vertex.UV)
	}

	mesh.arrayData = vertices
}

// UpdateVAO creates or updates the OpenGL Vertex Array Object and Vertex Buffer Object
// for this mesh. This should be called after vertex data has been prepared
// and before rendering.
func (mesh *Mesh) UpdateVAO() {
	vertices := mesh.arrayData

	// Create and bind Vertex Array Object (VAO)
	// VAO stores the vertex attribute configuration
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	mesh.VAO = vao

	// Create and bind Vertex Buffer Object (VBO)
	// VBO stores the actual vertex data
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)

	// Upload vertex data to GPU (4 bytes per float32)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	// Configure vertex attribute pointers
	// These tell OpenGL how to interpret the interleaved vertex data

	// Attribute 0: Position (3 floats)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 11*4, 0)

	// Attribute 1: Color (3 floats)
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, 11*4, 3*4)

	// Attribute 2: Normal (3 floats)
	gl.EnableVertexAttribArray(2)
	gl.VertexAttribPointerWithOffset(2, 3, gl.FLOAT, false, 11*4, 6*4)

	// Attribute 3: UV coordinates (2 floats)
	gl.EnableVertexAttribArray(3)
	gl.VertexAttribPointerWithOffset(3, 2, gl.FLOAT, false, 11*4, 9*4)

	// Unbind VBO and VAO (good practice to avoid accidental modifications)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
}

// BindMesh binds this mesh's VAO for rendering.
// Must be called before gl.DrawArrays or gl.DrawElements.
func (mesh *Mesh) BindMesh() {
	gl.BindVertexArray(mesh.VAO)
}

// Render draws the mesh using OpenGL's draw arrays command.
// Assumes the mesh's VAO is properly configured and vertex data is uploaded.
// Uses triangle primitive type - vertices should be in groups of 3.
func (mesh *Mesh) Render() {
	mesh.BindMesh()
	// Draw all vertices as triangles (3 vertices per triangle)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(mesh.vertices)))
}
