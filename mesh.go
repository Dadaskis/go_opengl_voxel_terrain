package main

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type MeshVertex struct {
	position mgl32.Vec3
	color    mgl32.Vec3
	normal   mgl32.Vec3
	UV       mgl32.Vec2
}

type Mesh struct {
	vertices []MeshVertex
	VAO      uint32
}

func (mesh *Mesh) AddVertex(position, color, normal mgl32.Vec3, UV mgl32.Vec2) {
	vertex := MeshVertex{
		position, color, normal, UV,
	}
	mesh.vertices = append(mesh.vertices, vertex)
}

func appendVec3ToArray(array []float32, vec *mgl32.Vec3) []float32 {
	array = append(array, vec[0])
	array = append(array, vec[1])
	array = append(array, vec[2])
	return array
}

func appendVec2ToArray(array []float32, vec *mgl32.Vec2) []float32 {
	array = append(array, vec[0])
	array = append(array, vec[1])
	return array
}

func (mesh *Mesh) UpdateVAO() {
	vertices := []float32{}

	for _, vertex := range mesh.vertices {
		vertices = appendVec3ToArray(vertices, &vertex.position)
		vertices = appendVec3ToArray(vertices, &vertex.color)
		vertices = appendVec3ToArray(vertices, &vertex.normal)
		vertices = appendVec2ToArray(vertices, &vertex.UV)
	}

	// Configure the vertex data
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	mesh.VAO = vao

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	// Position
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 11*4, 0)

	// Color
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, 11*4, 3*4)

	// Normal
	gl.EnableVertexAttribArray(2)
	gl.VertexAttribPointerWithOffset(2, 3, gl.FLOAT, false, 11*4, 6*4)

	// UV
	gl.EnableVertexAttribArray(3)
	gl.VertexAttribPointerWithOffset(3, 2, gl.FLOAT, false, 11*4, 9*4)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
}

func (mesh *Mesh) BindMesh() {
	gl.BindVertexArray(mesh.VAO)
}

func (mesh *Mesh) Render() {
	mesh.BindMesh()
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(mesh.vertices)))
}
