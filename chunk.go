package main

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/ojrac/opensimplex-go"
)

type Chunk struct {
	position mgl32.Vec2
	blocks   [16][16][256]int
	mesh     Mesh
}

func (chunk *Chunk) Generate() {
	noise := opensimplex.New(0)
	blockPos := chunk.position.Mul(16)
	scale := 0.01
	for x := range 16 {
		for y := range 16 {
			height := 50 + noise.Eval2(
				float64(int(blockPos[0])+x)*scale,
				float64(int(blockPos[1])+y)*scale,
			)*30.0
			for z := range int(height) {
				chunk.blocks[x][y][z] = BLOCK_STONE
				if (int(height) - z) < 5 {
					chunk.blocks[x][y][z] = BLOCK_DIRT
				}
				if (int(height) - z) <= 1 {
					chunk.blocks[x][y][z] = BLOCK_GRASS
				}
			}
		}
	}
}

func (chunk *Chunk) UpdateMesh() {
	chunk.mesh = Mesh{}
	blockPos := chunk.position.Mul(16)
	for x := range 16 {
		for y := range 16 {
			for z := range 256 {
				blockID := chunk.blocks[x][y][z]
				if blockID == BLOCK_AIR {
					continue
				}
				vertexPos := mgl32.Vec3{
					float32(int(blockPos[0]) + x),
					float32(z),
					float32(int(blockPos[1]) + y),
				}
				color := mgl32.Vec3{1.0, 1.0, 1.0}

				// SIDE 0

				chunk.mesh.AddVertex(
					vertexPos.Add(mgl32.Vec3{0.0, 0.0, 0.0}), color, mgl32.Vec3{0.0, 0.0, 0.0},
					mgl32.Vec2{
						(blockData[blockID].side0UV[0] + 0) * BLOCK_DATA_UV_SPACE,
						(blockData[blockID].side0UV[1] + 1) * BLOCK_DATA_UV_SPACE,
					},
				)

				chunk.mesh.AddVertex(
					vertexPos.Add(mgl32.Vec3{1.0, 0.0, 0.0}), color, mgl32.Vec3{0.0, 0.0, 0.0},
					mgl32.Vec2{
						(blockData[blockID].side0UV[0] + 1) * BLOCK_DATA_UV_SPACE,
						(blockData[blockID].side0UV[1] + 1) * BLOCK_DATA_UV_SPACE,
					},
				)

				chunk.mesh.AddVertex(
					vertexPos.Add(mgl32.Vec3{0.0, 1.0, 0.0}), color, mgl32.Vec3{0.0, 0.0, 0.0},
					mgl32.Vec2{
						(blockData[blockID].side0UV[0] + 0) * BLOCK_DATA_UV_SPACE,
						(blockData[blockID].side0UV[1] + 0) * BLOCK_DATA_UV_SPACE,
					},
				)

				chunk.mesh.AddVertex(
					vertexPos.Add(mgl32.Vec3{1.0, 1.0, 0.0}), color, mgl32.Vec3{0.0, 0.0, 0.0},
					mgl32.Vec2{
						(blockData[blockID].side0UV[0] + 1) * BLOCK_DATA_UV_SPACE,
						(blockData[blockID].side0UV[1] + 0) * BLOCK_DATA_UV_SPACE,
					},
				)

				chunk.mesh.AddVertex(
					vertexPos.Add(mgl32.Vec3{1.0, 0.0, 0.0}), color, mgl32.Vec3{0.0, 0.0, 0.0},
					mgl32.Vec2{
						(blockData[blockID].side0UV[0] + 1) * BLOCK_DATA_UV_SPACE,
						(blockData[blockID].side0UV[1] + 1) * BLOCK_DATA_UV_SPACE,
					},
				)

				chunk.mesh.AddVertex(
					vertexPos.Add(mgl32.Vec3{0.0, 1.0, 0.0}), color, mgl32.Vec3{0.0, 0.0, 0.0},
					mgl32.Vec2{
						(blockData[blockID].side0UV[0] + 0) * BLOCK_DATA_UV_SPACE,
						(blockData[blockID].side0UV[1] + 0) * BLOCK_DATA_UV_SPACE,
					},
				)
			}
		}
	}
	chunk.mesh.UpdateVAO()
}

func (chunk *Chunk) Render() {
	chunk.mesh.Render()
}
