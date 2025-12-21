package main

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/ojrac/opensimplex-go"
)

type Chunk struct {
	position    mgl32.Vec2
	blocks      [16][16][256]int
	mesh        Mesh
	isMeshDirty bool
}

func (chunk *Chunk) Generate() {
	go func() {
		noise := opensimplex.New(0)
		blockPos := chunk.position.Mul(16)
		scale := 0.01
		caveScale := 0.04
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
					caveValue := noise.Eval3(
						float64(int(blockPos[0])+x)*caveScale,
						float64(int(blockPos[1])+y)*caveScale,
						float64(z)*caveScale,
					)
					if caveValue > 0.6 {
						chunk.blocks[x][y][z] = BLOCK_AIR
					}
				}
				chunk.blocks[x][y][0] = BLOCK_STONE
			}
		}
		chunk.UpdateMesh()
	}()
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

				isBlockOccupied := func(x, y, z int) bool {
					if x < 0 {
						return false
					}

					if x >= 16 {
						return false
					}

					if y < 0 {
						return false
					}

					if y >= 16 {
						return false
					}

					if z < 0 {
						return false
					}

					if z >= 256 {
						return false
					}

					if chunk.blocks[x][y][z] == BLOCK_AIR {
						return false
					}

					return true
				}

				// SIDE 0
				if !isBlockOccupied(x, y-1, z) {
					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{0.0, 0.0, 0.0}), color, mgl32.Vec3{0.0, 0.0, -1.0},
						mgl32.Vec2{
							(blockData[blockID].side0UV[0] + 0) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side0UV[1] + 1) * BLOCK_DATA_UV_SPACE,
						},
					)

					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{1.0, 0.0, 0.0}), color, mgl32.Vec3{0.0, 0.0, -1.0},
						mgl32.Vec2{
							(blockData[blockID].side0UV[0] + 1) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side0UV[1] + 1) * BLOCK_DATA_UV_SPACE,
						},
					)

					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{0.0, 1.0, 0.0}), color, mgl32.Vec3{0.0, 0.0, -1.0},
						mgl32.Vec2{
							(blockData[blockID].side0UV[0] + 0) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side0UV[1] + 0) * BLOCK_DATA_UV_SPACE,
						},
					)

					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{1.0, 1.0, 0.0}), color, mgl32.Vec3{0.0, 0.0, -1.0},
						mgl32.Vec2{
							(blockData[blockID].side0UV[0] + 1) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side0UV[1] + 0) * BLOCK_DATA_UV_SPACE,
						},
					)

					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{1.0, 0.0, 0.0}), color, mgl32.Vec3{0.0, 0.0, -1.0},
						mgl32.Vec2{
							(blockData[blockID].side0UV[0] + 1) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side0UV[1] + 1) * BLOCK_DATA_UV_SPACE,
						},
					)

					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{0.0, 1.0, 0.0}), color, mgl32.Vec3{0.0, 0.0, -1.0},
						mgl32.Vec2{
							(blockData[blockID].side0UV[0] + 0) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side0UV[1] + 0) * BLOCK_DATA_UV_SPACE,
						},
					)
				}

				// SIDE 1

				if !isBlockOccupied(x-1, y, z) {
					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{0.0, 0.0, 0.0}), color, mgl32.Vec3{-1.0, 0.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].side1UV[0] + 0) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side1UV[1] + 1) * BLOCK_DATA_UV_SPACE,
						},
					)

					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{0.0, 0.0, 1.0}), color, mgl32.Vec3{-1.0, 0.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].side1UV[0] + 1) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side1UV[1] + 1) * BLOCK_DATA_UV_SPACE,
						},
					)

					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{0.0, 1.0, 0.0}), color, mgl32.Vec3{-1.0, 0.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].side1UV[0] + 0) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side1UV[1] + 0) * BLOCK_DATA_UV_SPACE,
						},
					)

					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{0.0, 1.0, 1.0}), color, mgl32.Vec3{-1.0, 0.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].side1UV[0] + 1) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side1UV[1] + 0) * BLOCK_DATA_UV_SPACE,
						},
					)

					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{0.0, 0.0, 1.0}), color, mgl32.Vec3{-1.0, 0.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].side1UV[0] + 1) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side1UV[1] + 1) * BLOCK_DATA_UV_SPACE,
						},
					)

					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{0.0, 1.0, 0.0}), color, mgl32.Vec3{-1.0, 0.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].side1UV[0] + 0) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side1UV[1] + 0) * BLOCK_DATA_UV_SPACE,
						},
					)
				}

				// SIDE 2

				if !isBlockOccupied(x+1, y, z) {
					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{1.0, 0.0, 0.0}), color, mgl32.Vec3{0.0, 0.0, 1.0},
						mgl32.Vec2{
							(blockData[blockID].side2UV[0] + 0) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side2UV[1] + 1) * BLOCK_DATA_UV_SPACE,
						},
					)

					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{1.0, 0.0, 1.0}), color, mgl32.Vec3{0.0, 0.0, 1.0},
						mgl32.Vec2{
							(blockData[blockID].side2UV[0] + 1) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side2UV[1] + 1) * BLOCK_DATA_UV_SPACE,
						},
					)

					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{1.0, 1.0, 0.0}), color, mgl32.Vec3{0.0, 0.0, 1.0},
						mgl32.Vec2{
							(blockData[blockID].side2UV[0] + 0) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side2UV[1] + 0) * BLOCK_DATA_UV_SPACE,
						},
					)

					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{1.0, 1.0, 1.0}), color, mgl32.Vec3{0.0, 0.0, 1.0},
						mgl32.Vec2{
							(blockData[blockID].side2UV[0] + 1) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side2UV[1] + 0) * BLOCK_DATA_UV_SPACE,
						},
					)

					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{1.0, 0.0, 1.0}), color, mgl32.Vec3{0.0, 0.0, 1.0},
						mgl32.Vec2{
							(blockData[blockID].side2UV[0] + 1) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side2UV[1] + 1) * BLOCK_DATA_UV_SPACE,
						},
					)

					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{1.0, 1.0, 0.0}), color, mgl32.Vec3{0.0, 0.0, 1.0},
						mgl32.Vec2{
							(blockData[blockID].side2UV[0] + 0) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side2UV[1] + 0) * BLOCK_DATA_UV_SPACE,
						},
					)
				}

				// SIDE 3

				if !isBlockOccupied(x, y+1, z) {
					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{0.0, 0.0, 1.0}), color, mgl32.Vec3{1.0, 0.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].side3UV[0] + 0) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side3UV[1] + 1) * BLOCK_DATA_UV_SPACE,
						},
					)

					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{1.0, 0.0, 1.0}), color, mgl32.Vec3{1.0, 0.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].side3UV[0] + 1) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side3UV[1] + 1) * BLOCK_DATA_UV_SPACE,
						},
					)

					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{0.0, 1.0, 1.0}), color, mgl32.Vec3{1.0, 0.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].side3UV[0] + 0) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side3UV[1] + 0) * BLOCK_DATA_UV_SPACE,
						},
					)

					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{1.0, 1.0, 1.0}), color, mgl32.Vec3{1.0, 0.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].side3UV[0] + 1) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side3UV[1] + 0) * BLOCK_DATA_UV_SPACE,
						},
					)

					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{1.0, 0.0, 1.0}), color, mgl32.Vec3{1.0, 0.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].side3UV[0] + 1) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side3UV[1] + 1) * BLOCK_DATA_UV_SPACE,
						},
					)

					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{0.0, 1.0, 1.0}), color, mgl32.Vec3{1.0, 0.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].side3UV[0] + 0) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side3UV[1] + 0) * BLOCK_DATA_UV_SPACE,
						},
					)
				}

				// TOP

				if !isBlockOccupied(x, y, z+1) {
					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{0.0, 1.0, 0.0}), color, mgl32.Vec3{0.0, 1.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].topUV[0] + 0) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].topUV[1] + 0) * BLOCK_DATA_UV_SPACE,
						},
					)

					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{0.0, 1.0, 1.0}), color, mgl32.Vec3{0.0, 1.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].topUV[0] + 0) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].topUV[1] + 1) * BLOCK_DATA_UV_SPACE,
						},
					)

					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{1.0, 1.0, 0.0}), color, mgl32.Vec3{0.0, 1.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].topUV[0] + 1) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].topUV[1] + 0) * BLOCK_DATA_UV_SPACE,
						},
					)

					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{1.0, 1.0, 0.0}), color, mgl32.Vec3{0.0, 1.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].topUV[0] + 1) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].topUV[1] + 0) * BLOCK_DATA_UV_SPACE,
						},
					)

					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{0.0, 1.0, 1.0}), color, mgl32.Vec3{0.0, 1.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].topUV[0] + 0) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].topUV[1] + 1) * BLOCK_DATA_UV_SPACE,
						},
					)

					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{1.0, 1.0, 1.0}), color, mgl32.Vec3{0.0, 1.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].topUV[0] + 1) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].topUV[1] + 1) * BLOCK_DATA_UV_SPACE,
						},
					)
				}

				// BOTTOM

				if !isBlockOccupied(x, y, z-1) {
					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{0.0, 0.0, 0.0}), color, mgl32.Vec3{0.0, -1.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].topUV[0] + 0) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].topUV[1] + 0) * BLOCK_DATA_UV_SPACE,
						},
					)

					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{0.0, 0.0, 1.0}), color, mgl32.Vec3{0.0, -1.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].topUV[0] + 0) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].topUV[1] + 1) * BLOCK_DATA_UV_SPACE,
						},
					)

					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{1.0, 0.0, 0.0}), color, mgl32.Vec3{0.0, -1.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].topUV[0] + 1) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].topUV[1] + 0) * BLOCK_DATA_UV_SPACE,
						},
					)

					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{1.0, 0.0, 0.0}), color, mgl32.Vec3{0.0, -1.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].topUV[0] + 1) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].topUV[1] + 0) * BLOCK_DATA_UV_SPACE,
						},
					)

					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{0.0, 0.0, 1.0}), color, mgl32.Vec3{0.0, -1.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].topUV[0] + 0) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].topUV[1] + 1) * BLOCK_DATA_UV_SPACE,
						},
					)

					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{1.0, 0.0, 1.0}), color, mgl32.Vec3{0.0, -1.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].topUV[0] + 1) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].topUV[1] + 1) * BLOCK_DATA_UV_SPACE,
						},
					)
				}
			}
		}
	}

	chunk.mesh.PrepareArrayData()

	chunk.isMeshDirty = true
}

func (chunk *Chunk) Render() {
	if chunk.isMeshDirty == true {
		chunk.mesh.UpdateVAO()
		chunk.isMeshDirty = false
	}
	chunk.mesh.Render()
}
